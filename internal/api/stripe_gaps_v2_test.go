package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/hckim/billtap/internal/billing"
)

func httptestRequest(t *testing.T, method string, path string, formBody string) *http.Request {
	t.Helper()
	var req *http.Request
	if formBody == "" {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, strings.NewReader(formBody))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Host = "billtap.test"
	return req
}

func httptestRecord(t *testing.T, handler http.Handler, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

// The tests in this file reproduce Stripe compatibility gaps found by diffing
// the ds2 platform-server Stripe SDK call surface against billtap. One test
// per gap; see docs/COMPATIBILITY_TRACKING.md conventions for claim wiring.

func gapsV2Customer(t *testing.T, handler http.Handler, email string) billing.Customer {
	t.Helper()
	return postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {email}})
}

func gapsV2Product(t *testing.T, handler http.Handler, name string) billing.Product {
	t.Helper()
	return postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {name}})
}

func gapsV2Price(t *testing.T, handler http.Handler, productID string, unitAmount string, interval string, extra url.Values) billing.Price {
	t.Helper()
	values := url.Values{
		"product":     {productID},
		"currency":    {"usd"},
		"unit_amount": {unitAmount},
	}
	if interval != "" {
		values.Set("recurring[interval]", interval)
	}
	for key, vals := range extra {
		values[key] = vals
	}
	return postForm[billing.Price](t, handler, "/v1/prices", values)
}

type gapsV2SubscriptionSeed struct {
	ID                 string `json:"id"`
	Status             string `json:"status"`
	CurrentPeriodStart int64  `json:"current_period_start"`
	CurrentPeriodEnd   int64  `json:"current_period_end"`
	Items              struct {
		Data []struct {
			ID       string `json:"id"`
			Quantity int64  `json:"quantity"`
		} `json:"data"`
	} `json:"items"`
}

func gapsV2Subscription(t *testing.T, handler http.Handler, customerID string, priceID string, extra url.Values) gapsV2SubscriptionSeed {
	t.Helper()
	values := url.Values{
		"customer":        {customerID},
		"items[0][price]": {priceID},
	}
	for key, vals := range extra {
		values[key] = vals
	}
	return postForm[gapsV2SubscriptionSeed](t, handler, "/v1/subscriptions", values)
}

type gapsV2SubscriptionView struct {
	ID                  string `json:"id"`
	Status              string `json:"status"`
	CancelAtPeriodEnd   bool   `json:"cancel_at_period_end"`
	CancelAt            any    `json:"cancel_at"`
	CurrentPeriodStart  int64  `json:"current_period_start"`
	CurrentPeriodEnd    int64  `json:"current_period_end"`
	CancellationDetails struct {
		Comment  any `json:"comment"`
		Feedback any `json:"feedback"`
	} `json:"cancellation_details"`
	Discounts []struct {
		ID     string `json:"id"`
		Source *struct {
			Type   string `json:"type"`
			Coupon *struct {
				ID     string `json:"id"`
				Object string `json:"object"`
			} `json:"coupon"`
		} `json:"source"`
	} `json:"discounts"`
}

type gapsV2InvoiceView struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Subtotal int64  `json:"subtotal"`
	Total    int64  `json:"total"`
	Parent   *struct {
		SubscriptionDetails *struct {
			Subscription string `json:"subscription"`
		} `json:"subscription_details"`
	} `json:"parent"`
}

// Gap: /v1/prices ignored lookup_keys[], so ds2's one-time extra-export price
// resolution could pick an arbitrary active one-time price.
func TestPricesListLookupKeysFilter(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "lookup@example.test")
	_ = customer
	productA := gapsV2Product(t, handler, "Lookup A")
	productB := gapsV2Product(t, handler, "Lookup B")
	priceA := gapsV2Price(t, handler, productA.ID, "1500", "", url.Values{"lookup_key": {"tenant_extra_export"}})
	priceB := gapsV2Price(t, handler, productB.ID, "2500", "", url.Values{"lookup_key": {"tenant_extra_export_v2"}})

	status, body := getStatus(t, handler, "/v1/prices?lookup_keys[0]=tenant_extra_export&active=true&type=one_time")
	if status != http.StatusOK {
		t.Fatalf("status=%d body=%s, want 200", status, body)
	}
	var list struct {
		Data []billing.Price `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(list.Data) != 1 || list.Data[0].ID != priceA.ID {
		t.Fatalf("prices = %#v, want only %s (priceB=%s)", list.Data, priceA.ID, priceB.ID)
	}

	// Multiple keys return both, and absent params return everything.
	status, body = getStatus(t, handler, "/v1/prices?lookup_keys[0]=tenant_extra_export&lookup_keys[1]=tenant_extra_export_v2&type=one_time")
	if err := json.Unmarshal([]byte(body), &list); err != nil {
		t.Fatalf("decode two-key list: %v", err)
	}
	if status != http.StatusOK || len(list.Data) != 2 {
		t.Fatalf("status=%d data=%d, want 2", status, len(list.Data))
	}
}

// Gap: checkout-completed subscription periods were hard-coded to +1 month,
// so yearly plans renewed after one month.
func TestCheckoutSubscriptionPeriodFollowsPriceInterval(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "yearly@example.test")
	product := gapsV2Product(t, handler, "Yearly")
	price := gapsV2Price(t, handler, product.ID, "99000", "year", nil)

	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":             {customer.ID},
		"mode":                 {"subscription"},
		"line_items[0][price]": {price.ID},
	})
	completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{"outcome": "payment_succeeded"})
	var subscription billing.Subscription
	if err := json.Unmarshal(completion["subscription"], &subscription); err != nil {
		t.Fatalf("decode subscription: %v", err)
	}
	days := subscription.CurrentPeriodEnd.Sub(subscription.CurrentPeriodStart).Hours() / 24
	if days < 364 || days > 366 {
		t.Fatalf("period days = %.1f, want ~365 for a yearly price", days)
	}
}

// Gap: subscription update accepted trial_end=now but left the subscription
// trialing; the trialing→paid upgrade confirm path depends on it.
func TestSubscriptionUpdateTrialEndNowEndsTrial(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "trial@example.test")
	product := gapsV2Product(t, handler, "Trial Plan")
	price := gapsV2Price(t, handler, product.ID, "4900", "month", nil)

	session := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                             {customer.ID},
		"line_items[0][price]":                 {price.ID},
		"subscription_data[trial_period_days]": {"14"},
	})
	completion := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session.ID+"/complete", map[string]string{"outcome": "payment_succeeded"})
	var seeded billing.Subscription
	if err := json.Unmarshal(completion["subscription"], &seeded); err != nil {
		t.Fatalf("decode subscription: %v", err)
	}
	if seeded.Status != "trialing" {
		t.Fatalf("seeded status = %q, want trialing", seeded.Status)
	}

	updated := postForm[gapsV2SubscriptionView](t, handler, "/v1/subscriptions/"+seeded.ID, url.Values{
		"items[0][id]":       {seeded.Items[0].ID},
		"items[0][quantity]": {"1"},
		"trial_end":          {"now"},
		"proration_behavior": {"none"},
	})
	if updated.Status != "active" {
		t.Fatalf("status after trial_end=now = %q, want active", updated.Status)
	}

	// Plain patch path (no item override) also ends the trial.
	session2 := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":                             {customer.ID},
		"line_items[0][price]":                 {price.ID},
		"subscription_data[trial_period_days]": {"7"},
	})
	completion2 := postJSON[map[string]json.RawMessage](t, handler, "/api/checkout/sessions/"+session2.ID+"/complete", map[string]string{"outcome": "payment_succeeded"})
	var seeded2 billing.Subscription
	if err := json.Unmarshal(completion2["subscription"], &seeded2); err != nil {
		t.Fatalf("decode subscription2: %v", err)
	}
	updated2 := postForm[gapsV2SubscriptionView](t, handler, "/v1/subscriptions/"+seeded2.ID, url.Values{
		"trial_end": {"now"},
	})
	if updated2.Status != "active" {
		t.Fatalf("plain-path status after trial_end=now = %q, want active", updated2.Status)
	}
}

// gapsV2TestClock creates a frozen test clock so subscription periods can be
// seeded deterministically (checkout completion runs at the frozen time).
func gapsV2TestClock(t *testing.T, handler http.Handler, frozen time.Time) struct{ ID string } {
	t.Helper()
	clock := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/test_helpers/test_clocks", url.Values{
		"frozen_time": {fmt.Sprintf("%d", frozen.Unix())},
	})
	return struct{ ID string }{ID: clock.ID}
}

// Gap: subscription list ignored current_period_end[gte]/[lt], so the ds2
// reminder and seat-sync sweeps processed every subscription in range-less fashion.
func TestSubscriptionsListCurrentPeriodEndRange(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "range@example.test")
	product := gapsV2Product(t, handler, "Range Plan")
	price := gapsV2Price(t, handler, product.ID, "1900", "month", nil)

	now := time.Now().UTC()
	clock1 := gapsV2TestClock(t, handler, now.Add(-24*time.Hour))
	clock2 := gapsV2TestClock(t, handler, now.Add(-48*time.Hour))
	first := gapsV2Subscription(t, handler, customer.ID, price.ID, url.Values{"test_clock": {clock1.ID}})
	second := gapsV2Subscription(t, handler, customer.ID, price.ID, url.Values{"test_clock": {clock2.ID}})
	if first.ID == second.ID {
		t.Fatalf("seeded subscriptions share id %s", first.ID)
	}
	if first.CurrentPeriodEnd <= second.CurrentPeriodEnd {
		t.Fatalf("seeded periods not ordered: first=%d second=%d", first.CurrentPeriodEnd, second.CurrentPeriodEnd)
	}
	low := second.CurrentPeriodEnd
	mid := first.CurrentPeriodEnd

	var list struct {
		Data []gapsV2SubscriptionView `json:"data"`
	}
	body := getJSON[json.RawMessage](t, handler, fmt.Sprintf("/v1/subscriptions?customer=%s&current_period_end[gte]=%d&current_period_end[lt]=%d", customer.ID, low, mid))
	if err := json.Unmarshal(body, &list); err != nil {
		t.Fatalf("decode gte/lt list: %v", err)
	}
	if len(list.Data) != 1 || list.Data[0].ID != second.ID {
		t.Fatalf("gte/lt data = %#v, want only %s", list.Data, second.ID)
	}

	body = getJSON[json.RawMessage](t, handler, fmt.Sprintf("/v1/subscriptions?customer=%s&current_period_end[gt]=%d", customer.ID, second.CurrentPeriodEnd))
	if err := json.Unmarshal(body, &list); err != nil {
		t.Fatalf("decode gt list: %v", err)
	}
	if len(list.Data) != 1 || list.Data[0].ID != first.ID {
		t.Fatalf("gt data = %#v, want only %s", list.Data, first.ID)
	}

	// Empty bounds return everything for the customer.
	body = getJSON[json.RawMessage](t, handler, fmt.Sprintf("/v1/subscriptions?customer=%s", customer.ID))
	if err := json.Unmarshal(body, &list); err != nil {
		t.Fatalf("decode unfiltered list: %v", err)
	}
	if len(list.Data) != 2 {
		t.Fatalf("unfiltered data = %d, want 2", len(list.Data))
	}
}

// Gap: updateCancelAt (BO reserved cancellation) sent cancel_at without
// cancel_at_period_end and was first rejected, then dropped by the serializer gate.
func TestSubscriptionUpdateCancelAtStandalone(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "reserved@example.test")
	product := gapsV2Product(t, handler, "Reserved Plan")
	price := gapsV2Price(t, handler, product.ID, "2900", "month", nil)
	subscription := gapsV2Subscription(t, handler, customer.ID, price.ID, nil)

	cancelAt := time.Now().UTC().Add(72 * time.Hour).Truncate(time.Second)
	updated := postForm[gapsV2SubscriptionView](t, handler, "/v1/subscriptions/"+subscription.ID, url.Values{
		"cancel_at": {fmt.Sprintf("%d", cancelAt.Unix())},
	})
	if updated.CancelAtPeriodEnd {
		t.Fatalf("cancel_at_period_end = true, want untouched false")
	}
	seconds, ok := updated.CancelAt.(float64)
	if !ok || int64(seconds) != cancelAt.Unix() {
		t.Fatalf("cancel_at = %#v, want %d", updated.CancelAt, cancelAt.Unix())
	}
}

// Gap: StripeParamBuilder.buildSubscriptionCreateParams always sends
// proration_behavior; subscription create rejected it with parameter_unknown.
func TestSubscriptionCreateAcceptsProrationBehavior(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "create-proration@example.test")
	product := gapsV2Product(t, handler, "Create Proration")
	price := gapsV2Price(t, handler, product.ID, "3900", "month", nil)

	subscription := gapsV2Subscription(t, handler, customer.ID, price.ID, url.Values{
		"proration_behavior": {"none"},
		"collection_method":  {"send_invoice"},
		"days_until_due":     {"1"},
		"cancel_at":          {fmt.Sprintf("%d", time.Now().UTC().Add(240*time.Hour).Unix())},
	})
	if subscription.ID == "" {
		t.Fatalf("subscription = %#v, want created with proration_behavior accepted", subscription)
	}
}

// Gap: invoiceitems rejected pricing[price]+quantity and required amount, and
// items without an invoice could not be stored as pending. The extra-export
// card-on-file path (prepareOneTimeInvoice) and the BO out-of-band seeding both
// depend on these forms.
func TestInvoiceItemPricingPriceAndPendingSweep(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "items@example.test")
	product := gapsV2Product(t, handler, "Extra Export")
	price := gapsV2Price(t, handler, product.ID, "1500", "", nil)

	// Pricing form derives amount = unit_amount × quantity and currency from the price.
	item := postForm[struct {
		ID       string `json:"id"`
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
		Invoice  any    `json:"invoice"`
	}](t, handler, "/v1/invoiceitems", url.Values{
		"customer":       {customer.ID},
		"pricing[price]": {price.ID},
		"quantity":       {"2"},
		"description":    {"Extra export pack"},
	})
	if item.Amount != 3000 || item.Currency != "usd" {
		t.Fatalf("item = %#v, want amount 3000 usd", item)
	}
	if item.Invoice != nil {
		t.Fatalf("item.invoice = %#v, want null pending item", item.Invoice)
	}

	// Invoice create sweeps the pending item by default (include).
	invoice := postForm[gapsV2InvoiceView](t, handler, "/v1/invoices", url.Values{
		"customer": {customer.ID},
	})
	if invoice.Subtotal != 3000 || invoice.Total != 3000 {
		t.Fatalf("swept invoice = %#v, want totals 3000", invoice)
	}

	// exclude skips the sweep.
	pending := postForm[billing.InvoiceItem](t, handler, "/v1/invoiceitems", url.Values{
		"customer": {customer.ID},
		"amount":   {"700"},
		"currency": {"usd"},
	})
	excluded := postForm[gapsV2InvoiceView](t, handler, "/v1/invoices", url.Values{
		"customer":                       {customer.ID},
		"pending_invoice_items_behavior": {"exclude"},
	})
	if excluded.Subtotal != 0 {
		t.Fatalf("exclude invoice subtotal = %d, want 0", excluded.Subtotal)
	}
	status, body := getStatus(t, handler, "/v1/invoiceitems/"+pending.ID)
	if status != http.StatusOK {
		t.Fatalf("pending item lookup status=%d body=%s", status, body)
	}

	// The amount path still attaches directly to an invoice.
	attached := postForm[billing.InvoiceItem](t, handler, "/v1/invoiceitems", url.Values{
		"customer": {customer.ID},
		"invoice":  {excluded.ID},
		"amount":   {"250"},
		"currency": {"usd"},
	})
	if attached.InvoiceID != excluded.ID || attached.Amount != 250 {
		t.Fatalf("attached item = %#v, want invoice %s amount 250", attached, excluded.ID)
	}
}

// Gap: invoice create rejected the subscription param, so manually created
// invoices never linked to their subscription (payment-history joins ran empty).
func TestInvoiceCreateWithSubscriptionParam(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "inv-sub@example.test")
	product := gapsV2Product(t, handler, "Linked Plan")
	price := gapsV2Price(t, handler, product.ID, "5900", "month", nil)
	subscription := gapsV2Subscription(t, handler, customer.ID, price.ID, nil)

	invoice := postForm[gapsV2InvoiceView](t, handler, "/v1/invoices", url.Values{
		"customer":     {customer.ID},
		"subscription": {subscription.ID},
	})
	if invoice.Parent == nil || invoice.Parent.SubscriptionDetails == nil || invoice.Parent.SubscriptionDetails.Subscription != subscription.ID {
		t.Fatalf("invoice parent = %#v, want subscription_details.subscription %s", invoice.Parent, subscription.ID)
	}

	var list struct {
		Data []gapsV2InvoiceView `json:"data"`
	}
	body := getJSON[json.RawMessage](t, handler, "/v1/invoices?subscription="+subscription.ID)
	if err := json.Unmarshal(body, &list); err != nil {
		t.Fatalf("decode invoice list: %v", err)
	}
	found := false
	for _, existing := range list.Data {
		if existing.ID == invoice.ID {
			found = true
		}
	}
	if !found {
		t.Fatalf("invoice list = %#v, want to contain %s", list.Data, invoice.ID)
	}

	status, bodyStr := postFormStatus(t, handler, "/v1/invoices", url.Values{
		"customer":     {customer.ID},
		"subscription": {"sub_does_not_exist"},
	})
	if status != http.StatusNotFound {
		t.Fatalf("unknown subscription status=%d body=%s, want 404", status, bodyStr)
	}
}

// Gap: customer update rejected invoice_settings[default_payment_method]; the
// extra-export card transition (setSucceededPaymentMethodAsCustomerDefault) 400'd.
func TestCustomerInvoiceSettingsDefaultPaymentMethod(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "dpm@example.test")

	updated := postForm[struct {
		InvoiceSettings struct {
			DefaultPaymentMethod any `json:"default_payment_method"`
		} `json:"invoice_settings"`
	}](t, handler, "/v1/customers/"+customer.ID, url.Values{
		"invoice_settings[default_payment_method]": {"pm_saved_card"},
	})
	if updated.InvoiceSettings.DefaultPaymentMethod != "pm_saved_card" {
		t.Fatalf("invoice_settings = %#v, want pm_saved_card", updated.InvoiceSettings)
	}

	var methods struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	body := getJSON[json.RawMessage](t, handler, "/v1/payment_methods?customer="+customer.ID+"&type=card")
	if err := json.Unmarshal(body, &methods); err != nil {
		t.Fatalf("decode payment methods: %v", err)
	}
	if len(methods.Data) == 0 || methods.Data[0].ID != "pm_saved_card" {
		t.Fatalf("payment methods = %#v, want pm_saved_card first", methods.Data)
	}

	// Empty string clears the default (Stripe unset form).
	cleared := postForm[struct {
		InvoiceSettings struct {
			DefaultPaymentMethod any `json:"default_payment_method"`
		} `json:"invoice_settings"`
	}](t, handler, "/v1/customers/"+customer.ID, url.Values{
		"invoice_settings[default_payment_method]": {""},
	})
	if cleared.InvoiceSettings.DefaultPaymentMethod != nil {
		t.Fatalf("cleared invoice_settings = %#v, want null", cleared.InvoiceSettings)
	}
}

// Gap: /v1/invoices/{id}/void and /mark_uncollectible were 404; the extra-export
// compensation path (voidInvoice) and SDS §4.0.3.8 (uncollectible) need them.
func TestInvoiceVoidAndMarkUncollectible(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "void@example.test")
	product := gapsV2Product(t, handler, "Voidable")
	_ = gapsV2Price(t, handler, product.ID, "1200", "", nil)

	newInvoice := func() string {
		invoice := postForm[gapsV2InvoiceView](t, handler, "/v1/invoices", url.Values{
			"customer":                       {customer.ID},
			"pending_invoice_items_behavior": {"exclude"},
		})
		postForm[billing.InvoiceItem](t, handler, "/v1/invoiceitems", url.Values{
			"customer": {customer.ID},
			"invoice":  {invoice.ID},
			"amount":   {"1200"},
			"currency": {"usd"},
		})
		return invoice.ID
	}

	voidable := newInvoice()
	postForm[map[string]any](t, handler, "/v1/invoices/"+voidable+"/finalize", url.Values{})
	invalid := newInvoice() // draft: void must refuse
	paidID := newInvoice()
	postForm[map[string]any](t, handler, "/v1/invoices/"+paidID+"/finalize", url.Values{})
	postForm[map[string]any](t, handler, "/v1/invoices/"+paidID+"/pay", url.Values{"paid_out_of_band": {"true"}})

	voided := postForm[gapsV2InvoiceView](t, handler, "/v1/invoices/"+voidable+"/void", url.Values{})
	if voided.Status != "void" {
		t.Fatalf("voided status = %q, want void", voided.Status)
	}

	if status, body := postFormStatus(t, handler, "/v1/invoices/"+invalid+"/void", url.Values{}); status != http.StatusBadRequest {
		t.Fatalf("draft void status=%d body=%s, want 400", status, body)
	}
	if status, body := postFormStatus(t, handler, "/v1/invoices/"+paidID+"/void", url.Values{}); status != http.StatusBadRequest {
		t.Fatalf("paid void status=%d body=%s, want 400", status, body)
	}

	uncollectibleID := newInvoice()
	postForm[map[string]any](t, handler, "/v1/invoices/"+uncollectibleID+"/finalize", url.Values{})
	uncollectible := postForm[gapsV2InvoiceView](t, handler, "/v1/invoices/"+uncollectibleID+"/mark_uncollectible", url.Values{})
	if uncollectible.Status != "uncollectible" {
		t.Fatalf("uncollectible status = %q, want uncollectible", uncollectible.Status)
	}
	if status, body := postFormStatus(t, handler, "/v1/invoices/"+newInvoice()+"/mark_uncollectible", url.Values{}); status != http.StatusBadRequest {
		t.Fatalf("draft mark_uncollectible status=%d body=%s, want 400", status, body)
	}
}

// Gap: POST /v1/checkout/sessions/{id}/expire was 404; the orphan-session
// compensation path (expireCheckoutSession) needs it.
func TestCheckoutSessionExpire(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "expire@example.test")
	product := gapsV2Product(t, handler, "Expirable")
	price := gapsV2Price(t, handler, product.ID, "800", "", nil)

	open := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":             {customer.ID},
		"mode":                 {"payment"},
		"line_items[0][price]": {price.ID},
	})
	expired := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions/"+open.ID+"/expire", url.Values{})
	if expired.Status != "expired" {
		t.Fatalf("status after expire = %q, want expired", expired.Status)
	}
	// Expired sessions are terminal: expiring again is resource_missing.
	if status, _ := postFormStatus(t, handler, "/v1/checkout/sessions/"+open.ID+"/expire", url.Values{}); status != http.StatusNotFound {
		t.Fatalf("re-expire status=%d, want 404", status)
	}

	completed := postForm[billing.CheckoutSession](t, handler, "/v1/checkout/sessions", url.Values{
		"customer":             {customer.ID},
		"mode":                 {"payment"},
		"line_items[0][price]": {price.ID},
	})
	postJSON[map[string]any](t, handler, "/api/checkout/sessions/"+completed.ID+"/complete", map[string]string{"outcome": "payment_succeeded"})
	if status, _ := postFormStatus(t, handler, "/v1/checkout/sessions/"+completed.ID+"/expire", url.Values{}); status != http.StatusNotFound {
		t.Fatalf("completed expire status=%d, want 404", status)
	}
}

// Gap: prices create rejected transfer_lookup_key; catalog convergence moves a
// lookup key to a new price and the old holder must lose it.
func TestPriceTransferLookupKey(t *testing.T) {
	handler := newTestHandler(t)
	product := gapsV2Product(t, handler, "Transfer")
	first := gapsV2Price(t, handler, product.ID, "5000", "month", url.Values{"lookup_key": {"tenant_plan_pro_monthly"}})

	second := postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"6000"},
		"recurring[interval]": {"month"},
		"lookup_key":          {"tenant_plan_pro_monthly"},
		"transfer_lookup_key": {"true"},
	})
	if second.LookupKey != "tenant_plan_pro_monthly" {
		t.Fatalf("new price lookup_key = %q, want tenant_plan_pro_monthly", second.LookupKey)
	}
	reloadedFirst := getJSON[billing.Price](t, handler, "/v1/prices/"+first.ID)
	if reloadedFirst.LookupKey != "" {
		t.Fatalf("old price lookup_key = %q, want cleared", reloadedFirst.LookupKey)
	}
}

// Gap: credit notes rejected out_of_band_amount/memo; the BO out-of-band refund
// path requires both, plus refunds staying empty and no balance transaction.
func TestCreditNoteOutOfBandAmountAndMemo(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "oob@example.test")

	invoice := postForm[gapsV2InvoiceView](t, handler, "/v1/invoices", url.Values{
		"customer":                       {customer.ID},
		"pending_invoice_items_behavior": {"exclude"},
	})
	postForm[billing.InvoiceItem](t, handler, "/v1/invoiceitems", url.Values{
		"customer": {customer.ID},
		"invoice":  {invoice.ID},
		"amount":   {"4000"},
		"currency": {"usd"},
	})
	postForm[map[string]any](t, handler, "/v1/invoices/"+invoice.ID+"/finalize", url.Values{})
	postForm[map[string]any](t, handler, "/v1/invoices/"+invoice.ID+"/pay", url.Values{"paid_out_of_band": {"true"}})

	note := postForm[struct {
		ID                         string `json:"id"`
		OutOfBandAmount            int64  `json:"out_of_band_amount"`
		Memo                       string `json:"memo"`
		CustomerBalanceTransaction any    `json:"customer_balance_transaction"`
		Refunds                    struct {
			Data []any `json:"data"`
		} `json:"refunds"`
		Status string `json:"status"`
	}](t, handler, "/v1/credit_notes", url.Values{
		"invoice":            {invoice.ID},
		"amount":             {"4000"},
		"out_of_band_amount": {"4000"},
		"memo":               {"BO refund"},
		"reason":             {"order_change"},
	})
	if note.OutOfBandAmount != 4000 || note.Memo != "BO refund" {
		t.Fatalf("note = %#v, want out_of_band_amount 4000 and memo", note)
	}
	if note.CustomerBalanceTransaction != nil || len(note.Refunds.Data) != 0 {
		t.Fatalf("note balance/refunds = %#v, want both empty for out-of-band", note)
	}
	if note.Status != "issued" {
		t.Fatalf("note status = %q, want issued", note.Status)
	}
}

// Gap: discount objects lacked source.coupon, so stripe-java 31.1.0
// (Discount.getSource().getCouponObject()) read null for coupon-backed discounts.
func TestSubscriptionDiscountSourceCouponObject(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "coupon@example.test")
	product := gapsV2Product(t, handler, "Coupon Plan")
	price := gapsV2Price(t, handler, product.ID, "3000", "month", nil)

	coupon := postForm[map[string]any](t, handler, "/v1/coupons", url.Values{
		"percent_off": {"10"},
		"duration":    {"forever"},
	})
	couponID, _ := coupon["id"].(string)

	subscription := gapsV2Subscription(t, handler, customer.ID, price.ID, url.Values{
		"coupon": {couponID},
	})
	view := getJSON[gapsV2SubscriptionView](t, handler, "/v1/subscriptions/"+subscription.ID)
	if len(view.Discounts) != 1 {
		t.Fatalf("discounts = %#v, want one", view.Discounts)
	}
	discount := view.Discounts[0]
	if discount.Source == nil || discount.Source.Coupon == nil || discount.Source.Coupon.ID != couponID || discount.Source.Coupon.Object != "coupon" {
		t.Fatalf("discount = %#v, want source.coupon %s", discount, couponID)
	}
}

// Gap: billing_cycle_anchor=now only reset the cycle on always_invoice; Stripe
// applies the anchor reset in every proration mode (ds2 bank-transfer path uses none).
func TestSubscriptionUpdateBillingCycleAnchorNowWithNone(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "anchor-none@example.test")
	product := gapsV2Product(t, handler, "Anchor Plan")
	price := gapsV2Price(t, handler, product.ID, "3300", "month", nil)
	// Seed the period 20 days in the past so the reset is observable.
	clock := gapsV2TestClock(t, handler, time.Now().UTC().Add(-20*24*time.Hour))
	subscription := gapsV2Subscription(t, handler, customer.ID, price.ID, url.Values{"test_clock": {clock.ID}})

	before := getJSON[gapsV2SubscriptionView](t, handler, "/v1/subscriptions/"+subscription.ID)
	updated := postForm[gapsV2SubscriptionView](t, handler, "/v1/subscriptions/"+subscription.ID, url.Values{
		"items[0][id]":         {subscription.Items.Data[0].ID},
		"items[0][quantity]":   {"1"},
		"proration_behavior":   {"none"},
		"billing_cycle_anchor": {"now"},
	})
	if updated.CurrentPeriodStart <= before.CurrentPeriodStart+86400 {
		t.Fatalf("period start %d not reset from %d", updated.CurrentPeriodStart, before.CurrentPeriodStart)
	}
	days := time.Unix(updated.CurrentPeriodEnd, 0).Sub(time.Unix(updated.CurrentPeriodStart, 0)).Hours() / 24
	if days < 28 || days > 31 {
		t.Fatalf("reset period = %.1f days, want one month", days)
	}
}

// Gap: immediate DELETE dropped cancellation_details[feedback|comment]; BO
// screens redisplay the collected reason from the subscription response.
func TestSubscriptionDeleteKeepsCancellationDetails(t *testing.T) {
	handler := newTestHandler(t)
	customer := gapsV2Customer(t, handler, "cancel-detail@example.test")
	product := gapsV2Product(t, handler, "Cancel Detail Plan")
	price := gapsV2Price(t, handler, product.ID, "2100", "month", nil)
	subscription := gapsV2Subscription(t, handler, customer.ID, price.ID, nil)

	req := httptestRequest(t, http.MethodDelete, "/v1/subscriptions/"+subscription.ID, url.Values{
		"cancellation_details[feedback]": {"too_expensive"},
		"cancellation_details[comment]":  {"Switching to annual elsewhere"},
	}.Encode())
	rec := httptestRecord(t, handler, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s, want 200", rec.Code, rec.Body.String())
	}
	var view gapsV2SubscriptionView
	if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
		t.Fatalf("decode canceled subscription: %v", err)
	}
	if view.Status != "canceled" {
		t.Fatalf("status = %q, want canceled", view.Status)
	}
	if view.CancellationDetails.Feedback != "too_expensive" || view.CancellationDetails.Comment != "Switching to annual elsewhere" {
		t.Fatalf("cancellation_details = %#v, want feedback+comment preserved", view.CancellationDetails)
	}
}
