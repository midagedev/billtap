package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/webhooks"
)

type localEvidenceStore struct {
	mu             sync.Mutex
	coupons        map[string]map[string]any
	promotionCodes map[string]map[string]any
	schedules      map[string]map[string]any
	cashBalances   map[string]int64
	cashTxs        map[string][]map[string]any
	disputes       map[string]map[string]any
}

func newLocalEvidenceStore() *localEvidenceStore {
	return &localEvidenceStore{
		coupons:        map[string]map[string]any{},
		promotionCodes: map[string]map[string]any{},
		schedules:      map[string]map[string]any{},
		cashBalances:   map[string]int64{},
		cashTxs:        map[string][]map[string]any{},
		disputes:       map[string]map[string]any{},
	}
}

func (h *Handler) handleCoupons(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateCouponCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		now := time.Now().UTC()
		id := p.string("id")
		if id == "" {
			id = "coupon_" + strconv.FormatInt(now.UnixNano(), 36)
		}
		coupon := map[string]any{
			"id":          id,
			"object":      "coupon",
			"name":        emptyToNil(p.string("name")),
			"duration":    p.stringDefault("duration", "once"),
			"percent_off": nil,
			"amount_off":  nil,
			"currency":    emptyToNil(p.string("currency")),
			"valid":       true,
			"metadata":    nonNilMap(p.metadata()),
			"created":     now.Unix(),
			"livemode":    false,
		}
		if p.has("percent_off") {
			coupon["percent_off"] = p.int64("percent_off")
		}
		if p.has("amount_off") {
			coupon["amount_off"] = p.int64("amount_off")
		}
		h.local.mu.Lock()
		h.local.coupons[id] = coupon
		h.local.mu.Unlock()
		writeJSON(w, http.StatusOK, cloneEvidence(coupon))
	case http.MethodGet:
		h.local.mu.Lock()
		data := evidenceList(h.local.coupons)
		h.local.mu.Unlock()
		writeJSON(w, http.StatusOK, stripeList(r.URL.Path, data))
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleCoupon(w http.ResponseWriter, r *http.Request) {
	id := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/coupons/"), "/")
	if id == "" || strings.Contains(id, "/") {
		h.notFound(w, r)
		return
	}
	h.local.mu.Lock()
	coupon, ok := h.local.coupons[id]
	h.local.mu.Unlock()
	if !ok {
		writeResult(w, nil, billing.ErrNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet, http.MethodPost:
		writeJSON(w, http.StatusOK, cloneEvidence(coupon))
	case http.MethodDelete:
		deleted := map[string]any{"id": id, "object": "coupon", "deleted": true}
		h.local.mu.Lock()
		delete(h.local.coupons, id)
		h.local.mu.Unlock()
		writeJSON(w, http.StatusOK, deleted)
	default:
		h.methodNotAllowed(w, r, "GET, POST, DELETE")
	}
}

func (h *Handler) handlePromotionCodes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validatePromotionCodeCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		now := time.Now().UTC()
		id := p.string("id")
		if id == "" {
			id = "promo_" + strconv.FormatInt(now.UnixNano(), 36)
		}
		code := p.string("code")
		if code == "" {
			code = strings.ToUpper(id)
		}
		promo := map[string]any{
			"id":       id,
			"object":   "promotion_code",
			"code":     code,
			"coupon":   evidenceCouponRef(p.string("coupon")),
			"active":   p.boolDefault("active", true),
			"metadata": nonNilMap(p.metadata()),
			"created":  now.Unix(),
			"livemode": false,
		}
		h.local.mu.Lock()
		h.local.promotionCodes[id] = promo
		h.local.mu.Unlock()
		writeJSON(w, http.StatusOK, cloneEvidence(promo))
	case http.MethodGet:
		h.local.mu.Lock()
		data := evidenceList(h.local.promotionCodes)
		h.local.mu.Unlock()
		writeJSON(w, http.StatusOK, stripeList(r.URL.Path, data))
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handlePromotionCode(w http.ResponseWriter, r *http.Request) {
	id := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/promotion_codes/"), "/")
	if id == "" || strings.Contains(id, "/") {
		h.notFound(w, r)
		return
	}
	h.local.mu.Lock()
	promo, ok := h.local.promotionCodes[id]
	h.local.mu.Unlock()
	if !ok {
		writeResult(w, nil, billing.ErrNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet, http.MethodPost:
		writeJSON(w, http.StatusOK, cloneEvidence(promo))
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func evidenceCouponRef(id string) map[string]any {
	return map[string]any{"id": id, "object": "coupon"}
}

func (h *Handler) handleSubscriptionSchedules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateSubscriptionScheduleCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		subscriptionID := p.first("from_subscription", "subscription")
		subscription, err := h.billing.GetSubscription(r.Context(), subscriptionID)
		if err != nil {
			writeResult(w, nil, err)
			return
		}
		now := time.Now().UTC()
		id := p.string("id")
		if id == "" {
			id = "sub_sched_" + strconv.FormatInt(now.UnixNano(), 36)
		}
		startDate := p.int64Default("phases[0][start_date]", subscription.CurrentPeriodEnd.Unix())
		priceID := p.first("phases[0][items][0][price]", "phases[0][plans][0][price]")
		if priceID == "" {
			writeError(w, http.StatusBadRequest, missingParam("phases[0][items][0][price]"))
			return
		}
		if err := validatePriceExists(h.billing.GetPrice(r.Context(), priceID)); err != nil {
			writeResult(w, nil, err)
			return
		}
		quantity := p.int64Default("phases[0][items][0][quantity]", 1)
		if quantity <= 0 {
			writeError(w, http.StatusBadRequest, invalidParam("phases[0][items][0][quantity]", "Must be at least 1."))
			return
		}
		schedule := map[string]any{
			"id":           id,
			"object":       "subscription_schedule",
			"customer":     subscription.CustomerID,
			"subscription": subscription.ID,
			"status":       "active",
			"start_date":   startDate,
			"end_behavior": p.stringDefault("end_behavior", "release"),
			"phases": []map[string]any{{
				"start_date": startDate,
				"items": []map[string]any{{
					"price":    priceID,
					"quantity": quantity,
				}},
			}},
			"metadata": nonNilMap(p.metadata()),
			"created":  now.Unix(),
			"livemode": false,
		}
		h.local.mu.Lock()
		h.local.schedules[id] = schedule
		h.local.mu.Unlock()
		writeJSON(w, http.StatusOK, cloneEvidence(schedule))
	case http.MethodGet:
		h.local.mu.Lock()
		data := evidenceList(h.local.schedules)
		h.local.mu.Unlock()
		writeJSON(w, http.StatusOK, stripeList(r.URL.Path, data))
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleSubscriptionSchedule(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/subscription_schedules/"), "/")
	id, action, hasAction := strings.Cut(rest, "/")
	if id == "" {
		h.notFound(w, r)
		return
	}
	h.local.mu.Lock()
	schedule, ok := h.local.schedules[id]
	h.local.mu.Unlock()
	if !ok {
		writeResult(w, nil, billing.ErrNotFound)
		return
	}
	if !hasAction {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			h.methodNotAllowed(w, r, "GET, POST")
			return
		}
		writeJSON(w, http.StatusOK, cloneEvidence(schedule))
		return
	}
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	switch action {
	case "cancel":
		schedule["status"] = "canceled"
	case "release":
		schedule["status"] = "released"
	default:
		h.notFound(w, r)
		return
	}
	h.local.mu.Lock()
	h.local.schedules[id] = schedule
	h.local.mu.Unlock()
	writeJSON(w, http.StatusOK, cloneEvidence(schedule))
}

func (h *Handler) applyDueSubscriptionSchedules(r *http.Request, clockID string, at time.Time) []billing.Subscription {
	h.local.mu.Lock()
	schedules := evidenceList(h.local.schedules)
	h.local.mu.Unlock()

	var updated []billing.Subscription
	for _, schedule := range schedules {
		if schedule["status"] != "active" {
			continue
		}
		startDate, _ := asInt64Evidence(schedule["start_date"])
		if startDate == 0 || startDate > at.Unix() {
			continue
		}
		subscriptionID, _ := schedule["subscription"].(string)
		priceID, quantity := schedulePhaseItem(schedule)
		if subscriptionID == "" || priceID == "" {
			continue
		}
		if !h.subscriptionAttachedToClock(r, subscriptionID, clockID) {
			continue
		}
		subscription, err := h.billing.PatchSubscription(r.Context(), subscriptionID, billing.SubscriptionPatch{
			Items:        []billing.LineItem{{PriceID: priceID, Quantity: quantity}},
			ReplaceItems: true,
			Metadata:     map[string]string{"billtap_subscription_schedule": fmt.Sprint(schedule["id"])},
		})
		if err != nil {
			continue
		}
		schedule["status"] = "completed"
		h.local.mu.Lock()
		h.local.schedules[fmt.Sprint(schedule["id"])] = schedule
		h.local.mu.Unlock()
		h.emitSubscriptionWebhook(r, "customer.subscription.updated", subscription, webhooks.SourceAPI)
		updated = append(updated, subscription)
	}
	return updated
}

func (h *Handler) subscriptionAttachedToClock(r *http.Request, subscriptionID string, clockID string) bool {
	clockID = strings.TrimSpace(clockID)
	if clockID == "" {
		return true
	}
	subscription, err := h.billing.GetSubscription(r.Context(), subscriptionID)
	if err != nil {
		return false
	}
	for _, key := range []string{"test_clock", "testClock"} {
		if strings.TrimSpace(subscription.Metadata[key]) == clockID {
			return true
		}
	}
	customer, err := h.billing.GetCustomer(r.Context(), subscription.CustomerID)
	if err != nil {
		return false
	}
	for _, key := range []string{"test_clock", "testClock"} {
		if strings.TrimSpace(customer.Metadata[key]) == clockID {
			return true
		}
	}
	return false
}

func schedulePhaseItem(schedule map[string]any) (string, int64) {
	phases, _ := schedule["phases"].([]map[string]any)
	if len(phases) == 0 {
		return "", 1
	}
	items, _ := phases[0]["items"].([]map[string]any)
	if len(items) == 0 {
		return "", 1
	}
	priceID, _ := items[0]["price"].(string)
	quantity, _ := asInt64Evidence(items[0]["quantity"])
	if quantity <= 0 {
		quantity = 1
	}
	return priceID, quantity
}

func (h *Handler) handleCustomerCashBalance(w http.ResponseWriter, r *http.Request, customerID string) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "GET, POST")
		return
	}
	if _, err := h.billing.GetCustomer(r.Context(), customerID); err != nil {
		writeResult(w, nil, err)
		return
	}
	h.local.mu.Lock()
	available := h.local.cashBalances[customerID]
	h.local.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{
		"object":    "cash_balance",
		"customer":  customerID,
		"available": map[string]int64{"usd": available},
		"livemode":  false,
	})
}

func (h *Handler) handleCustomerCashBalanceTransactions(w http.ResponseWriter, r *http.Request, customerID string, transactionID string) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	h.local.mu.Lock()
	transactions := append([]map[string]any(nil), h.local.cashTxs[customerID]...)
	h.local.mu.Unlock()
	if transactionID != "" {
		for _, tx := range transactions {
			if tx["id"] == transactionID {
				writeJSON(w, http.StatusOK, cloneEvidence(tx))
				return
			}
		}
		writeResult(w, nil, billing.ErrNotFound)
		return
	}
	writeJSON(w, http.StatusOK, stripeList(r.URL.Path, transactions))
}

func (h *Handler) handleTestHelperCustomer(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/test_helpers/customers/"), "/")
	customerID, action, found := strings.Cut(rest, "/")
	if customerID == "" || !found || action != "fund_cash_balance" {
		h.notFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := validateFundCashBalance(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if _, err := h.billing.GetCustomer(r.Context(), customerID); err != nil {
		writeResult(w, nil, err)
		return
	}
	amount := p.int64("amount")
	now := time.Now().UTC()
	tx := map[string]any{
		"id":         "ccsbtxn_" + strconv.FormatInt(now.UnixNano(), 36),
		"object":     "customer_cash_balance_transaction",
		"customer":   customerID,
		"type":       "funded",
		"net_amount": amount,
		"currency":   strings.ToLower(p.stringDefault("currency", "usd")),
		"created":    now.Unix(),
		"livemode":   false,
	}
	h.local.mu.Lock()
	h.local.cashBalances[customerID] += amount
	h.local.cashTxs[customerID] = append(h.local.cashTxs[customerID], tx)
	h.local.mu.Unlock()
	settled, _ := h.billing.SettleBankTransferPaymentIntents(r.Context(), customerID)
	for _, intent := range settled {
		h.emitPaymentIntentWebhook(r, "payment_intent.succeeded", intent)
	}
	writeJSON(w, http.StatusOK, tx)
}

func (h *Handler) handleDisputeSimulation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := validateDisputeCreate(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	dispute := h.createDispute(r, p.first("charge", "charge_id"), p.int64("amount"), p.stringDefault("reason", "general"))
	writeJSON(w, http.StatusOK, dispute)
}

func (h *Handler) handleChargeSubresource(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/charges/"), "/")
	chargeID, action, found := strings.Cut(rest, "/")
	if chargeID == "" || !found || action != "dispute" {
		h.notFound(w, r)
		return
	}
	if r.Method == http.MethodGet {
		h.local.mu.Lock()
		defer h.local.mu.Unlock()
		for _, dispute := range h.local.disputes {
			if dispute["charge"] == chargeID {
				writeJSON(w, http.StatusOK, cloneEvidence(dispute))
				return
			}
		}
		writeResult(w, nil, billing.ErrNotFound)
		return
	}
	if r.Method != http.MethodPost {
		h.methodNotAllowed(w, r, "GET, POST")
		return
	}
	p, err := parseParams(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := validateDisputeCreate(p); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	dispute := h.createDispute(r, chargeID, p.int64("amount"), p.stringDefault("reason", "general"))
	writeJSON(w, http.StatusOK, dispute)
}

func (h *Handler) handleDisputes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.methodNotAllowed(w, r, "GET")
		return
	}
	h.local.mu.Lock()
	data := evidenceList(h.local.disputes)
	h.local.mu.Unlock()
	writeJSON(w, http.StatusOK, stripeList(r.URL.Path, data))
}

func (h *Handler) handleDispute(w http.ResponseWriter, r *http.Request) {
	rest := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/disputes/"), "/")
	id, action, found := strings.Cut(rest, "/")
	if id == "" {
		h.notFound(w, r)
		return
	}
	h.local.mu.Lock()
	dispute, ok := h.local.disputes[id]
	h.local.mu.Unlock()
	if !ok {
		writeResult(w, nil, billing.ErrNotFound)
		return
	}
	if !found {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			h.methodNotAllowed(w, r, "GET, POST")
			return
		}
		if r.Method == http.MethodPost {
			p, err := parseParams(r)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			evidence := map[string]any{}
			if current, ok := dispute["evidence"].(map[string]any); ok {
				for key, value := range current {
					evidence[key] = value
				}
			}
			for key, value := range p.values {
				if strings.HasPrefix(key, "evidence[") && strings.HasSuffix(key, "]") {
					evidence[strings.TrimSuffix(strings.TrimPrefix(key, "evidence["), "]")] = value
				}
			}
			if status := p.string("status"); status != "" {
				dispute["status"] = status
			}
			if metadata := p.metadata(); metadata != nil {
				dispute["metadata"] = nonNilMap(metadata)
			}
			dispute["evidence"] = evidence
			dispute["evidence_details"] = map[string]any{
				"has_evidence":     len(evidence) > 0,
				"submission_count": 1,
				"past_due":         false,
			}
			h.local.mu.Lock()
			h.local.disputes[id] = dispute
			h.local.mu.Unlock()
			h.emitGenericWebhook(r, "charge.dispute.updated", id, dispute, webhooks.SourceAPI)
		}
		writeJSON(w, http.StatusOK, cloneEvidence(dispute))
		return
	}
	if action != "close" || r.Method != http.MethodPost {
		h.notFound(w, r)
		return
	}
	dispute["status"] = "won"
	dispute["closed_at"] = time.Now().UTC().Unix()
	h.local.mu.Lock()
	h.local.disputes[id] = dispute
	h.local.mu.Unlock()
	h.emitGenericWebhook(r, "charge.dispute.closed", id, dispute, webhooks.SourceAPI)
	writeJSON(w, http.StatusOK, cloneEvidence(dispute))
}

func (h *Handler) createDispute(r *http.Request, chargeID string, amount int64, reason string) map[string]any {
	now := time.Now().UTC()
	if chargeID == "" {
		chargeID = "ch_" + strconv.FormatInt(now.UnixNano(), 36)
	}
	if amount <= 0 {
		amount = 1000
	}
	dispute := map[string]any{
		"id":       "dp_" + strconv.FormatInt(now.UnixNano(), 36),
		"object":   "dispute",
		"charge":   chargeID,
		"amount":   amount,
		"currency": "usd",
		"reason":   reason,
		"status":   "needs_response",
		"created":  now.Unix(),
		"livemode": false,
	}
	h.local.mu.Lock()
	h.local.disputes[fmt.Sprint(dispute["id"])] = dispute
	h.local.mu.Unlock()
	h.emitGenericWebhook(r, "charge.dispute.created", chargeID, dispute, webhooks.SourceAPI)
	return cloneEvidence(dispute)
}

func evidenceList(items map[string]map[string]any) []map[string]any {
	data := make([]map[string]any, 0, len(items))
	for _, item := range items {
		data = append(data, cloneEvidence(item))
	}
	return data
}

func cloneEvidence(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func asInt64Evidence(value any) (int64, bool) {
	switch v := value.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}
