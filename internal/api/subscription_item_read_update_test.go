package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hckim/billtap/internal/billing"
)

type itemReadResponse struct {
	ID           string `json:"id"`
	Object       string `json:"object"`
	Subscription string `json:"subscription"`
	Quantity     int64  `json:"quantity"`
	Price        struct {
		ID string `json:"id"`
	} `json:"price"`
	Metadata map[string]string `json:"metadata"`
}

func setupItemReadFixtures(t *testing.T, handler http.Handler) (customer billing.Customer, first, second billing.Price, subscription prorationSubResponse) {
	t.Helper()
	customer = postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"items-read@example.test"}})
	product := postForm[billing.Product](t, handler, "/v1/products", url.Values{"name": {"Item Read Plan"}})
	first = postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"5000"},
		"recurring[interval]": {"month"},
	})
	second = postForm[billing.Price](t, handler, "/v1/prices", url.Values{
		"product":             {product.ID},
		"currency":            {"usd"},
		"unit_amount":         {"7000"},
		"recurring[interval]": {"month"},
	})
	subscription = postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":           {customer.ID},
		"items[0][price]":    {first.ID},
		"items[0][quantity]": {"1"},
	})
	return customer, first, second, subscription
}

func TestSubscriptionItemListAndRetrieve(t *testing.T) {
	handler := newTestHandler(t)
	_, first, _, subscription := setupItemReadFixtures(t, handler)
	secondItem := postForm[itemReadResponse](t, handler, "/v1/subscription_items", url.Values{
		"subscription": {subscription.ID},
		"price":        {first.ID},
		"quantity":     {"2"},
	})

	listed := getJSON[struct {
		Data []itemReadResponse `json:"data"`
	}](t, handler, "/v1/subscription_items")
	if len(listed.Data) != 2 {
		t.Fatalf("list items = %d, want 2", len(listed.Data))
	}
	for _, item := range listed.Data {
		if item.Subscription != subscription.ID || item.Price.ID != first.ID {
			t.Fatalf("listed item = %#v, want subscription %s price %s", item, subscription.ID, first.ID)
		}
	}

	filtered := getJSON[struct {
		Data []itemReadResponse `json:"data"`
	}](t, handler, "/v1/subscription_items?subscription=sub_other")
	if len(filtered.Data) != 0 {
		t.Fatalf("filtered items = %d, want 0 for other subscription", len(filtered.Data))
	}

	fetched := getJSON[itemReadResponse](t, handler, "/v1/subscription_items/"+secondItem.ID)
	if fetched.ID != secondItem.ID || fetched.Quantity != 2 || fetched.Subscription != subscription.ID {
		t.Fatalf("retrieved item = %#v, want created item", fetched)
	}

	missingReq := httptest.NewRequest(http.MethodGet, "/v1/subscription_items/si_missing", nil)
	missingRec := httptest.NewRecorder()
	handler.ServeHTTP(missingRec, missingReq)
	if missingRec.Code != http.StatusNotFound {
		t.Fatalf("missing item status = %d body = %s, want 404", missingRec.Code, missingRec.Body.String())
	}

	// List validation: unknown filter params are rejected.
	badReq := httptest.NewRequest(http.MethodGet, "/v1/subscription_items?price="+first.ID, nil)
	badRec := httptest.NewRecorder()
	handler.ServeHTTP(badRec, badReq)
	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("unknown list filter status = %d body = %s, want 400", badRec.Code, badRec.Body.String())
	}
}

func TestSubscriptionItemUpdateQuantityAndPrice(t *testing.T) {
	handler := newTestHandler(t)
	_, first, second, subscription := setupItemReadFixtures(t, handler)
	itemID := subscription.Items.Data[0].ID

	// Unknown params are rejected.
	status, body := postFormStatus(t, handler, "/v1/subscription_items/"+itemID, url.Values{
		"payment_behavior": {"pending_if_incomplete"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("unknown update param status = %d body = %s, want 400", status, body)
	}
	// quantity must stay positive.
	status, body = postFormStatus(t, handler, "/v1/subscription_items/"+itemID, url.Values{
		"quantity": {"0"},
	})
	if status != http.StatusBadRequest {
		t.Fatalf("zero quantity status = %d body = %s, want 400", status, body)
	}
	// Missing item is 404.
	status, body = postFormStatus(t, handler, "/v1/subscription_items/si_missing", url.Values{
		"quantity": {"2"},
	})
	if status != http.StatusNotFound {
		t.Fatalf("missing item update status = %d body = %s, want 404", status, body)
	}

	updated := postForm[itemReadResponse](t, handler, "/v1/subscription_items/"+itemID, url.Values{
		"quantity":           {"3"},
		"metadata[k]":        {"v"},
		"proration_behavior": {"none"},
	})
	if updated.ID != itemID || updated.Quantity != 3 || updated.Price.ID != first.ID {
		t.Fatalf("updated item = %#v, want same id quantity 3 price %s", updated, first.ID)
	}
	if updated.Metadata["k"] != "v" {
		t.Fatalf("updated item metadata = %#v, want k=v", updated.Metadata)
	}

	// Price swap keeps the item id and subscription shape.
	swapped := postForm[itemReadResponse](t, handler, "/v1/subscription_items/"+itemID, url.Values{
		"price": {second.ID},
	})
	if swapped.ID != itemID || swapped.Price.ID != second.ID {
		t.Fatalf("swapped item = %#v, want same id price %s", swapped, second.ID)
	}
	subscriptionAfter := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+subscription.ID)
	if len(subscriptionAfter.Items.Data) != 1 || subscriptionAfter.Items.Data[0].ID != itemID {
		t.Fatalf("subscription after item update = %#v, want single stable item", subscriptionAfter.Items.Data)
	}
	if subscriptionAfter.Items.Data[0].Price.ID != second.ID {
		t.Fatalf("subscription item price after swap = %s, want %s", subscriptionAfter.Items.Data[0].Price.ID, second.ID)
	}
}

func TestSubscriptionItemUpdateAlwaysInvoiceBillsDelta(t *testing.T) {
	handler := newTestHandler(t)
	_, _, _, subscription := setupItemReadFixtures(t, handler)
	itemID := subscription.Items.Data[0].ID

	updated := postForm[itemReadResponse](t, handler, "/v1/subscription_items/"+itemID, url.Values{
		"quantity":           {"2"},
		"proration_behavior": {"always_invoice"},
		"proration_date":     {fmt.Sprint(subscription.CurrentPeriodStart)},
	})
	if updated.Quantity != 2 || updated.ID != itemID {
		t.Fatalf("always_invoice item update = %#v, want quantity 2 stable id", updated)
	}
	after := getJSON[prorationSubResponse](t, handler, "/v1/subscriptions/"+subscription.ID)
	if after.LatestInvoice == subscription.LatestInvoice {
		t.Fatalf("latest_invoice unchanged after always_invoice update: %s", after.LatestInvoice)
	}
}

func TestNestedCustomerSubscriptionDiscountRoutes(t *testing.T) {
	handler := newTestHandler(t)
	customer, first, _, _ := setupItemReadFixtures(t, handler)
	coupon := postForm[struct {
		ID string `json:"id"`
	}](t, handler, "/v1/coupons", url.Values{
		"percent_off": {"10"},
		"duration":    {"forever"},
	})

	discounted := postForm[prorationSubResponse](t, handler, "/v1/subscriptions", url.Values{
		"customer":             {customer.ID},
		"items[0][price]":      {first.ID},
		"items[0][quantity]":   {"1"},
		"discounts[0][coupon]": {coupon.ID},
	})
	nestedPath := "/v1/customers/" + customer.ID + "/subscriptions/" + discounted.ID + "/discount"

	discount := getJSON[struct {
		Object string `json:"object"`
		Coupon struct {
			ID string `json:"id"`
		} `json:"coupon"`
		Subscription string `json:"subscription"`
	}](t, handler, nestedPath)
	if discount.Object != "discount" || discount.Coupon.ID != coupon.ID || discount.Subscription != discounted.ID {
		t.Fatalf("nested discount = %#v, want %s coupon on %s", discount, coupon.ID, discounted.ID)
	}

	// The nested route is customer-scoped: another customer sees 404.
	other := postForm[billing.Customer](t, handler, "/v1/customers", url.Values{"email": {"other@example.test"}})
	scopedReq := httptest.NewRequest(http.MethodGet, "/v1/customers/"+other.ID+"/subscriptions/"+discounted.ID+"/discount", nil)
	scopedRec := httptest.NewRecorder()
	handler.ServeHTTP(scopedRec, scopedReq)
	if scopedRec.Code != http.StatusNotFound {
		t.Fatalf("other-customer nested discount status = %d body = %s, want 404", scopedRec.Code, scopedRec.Body.String())
	}

	deleted := deleteForm[struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Deleted bool   `json:"deleted"`
	}](t, handler, nestedPath, url.Values{})
	if !deleted.Deleted || deleted.Object != "discount" {
		t.Fatalf("nested discount delete = %#v, want deleted discount", deleted)
	}
	afterReq := httptest.NewRequest(http.MethodGet, nestedPath, nil)
	afterRec := httptest.NewRecorder()
	handler.ServeHTTP(afterRec, afterReq)
	if afterRec.Code != http.StatusNotFound {
		t.Fatalf("nested discount after delete status = %d body = %s, want 404", afterRec.Code, afterRec.Body.String())
	}
}
