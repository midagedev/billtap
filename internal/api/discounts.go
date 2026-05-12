package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/webhooks"
)

func (h *Handler) discountsFromParamsOrCustomer(r *http.Request, p params, customer billing.Customer) ([]billing.Discount, error) {
	discounts, err := h.discountsFromParams(p)
	if err != nil || len(discounts) > 0 {
		return discounts, err
	}
	return billing.DiscountsFromMetadata(customer.Metadata), nil
}

func hasDiscountParams(p params) bool {
	if p.hasAny("coupon", "promotion_code", "discounts[0][coupon]", "discounts[0][promotion_code]") {
		return true
	}
	for key := range p.values {
		if strings.HasPrefix(key, "discounts[") {
			return true
		}
	}
	return false
}

func (h *Handler) discountsFromParams(p params) ([]billing.Discount, error) {
	couponID := p.first("discounts[0][coupon]", "coupon")
	promotionCodeID := p.first("discounts[0][promotion_code]", "promotion_code")
	if couponID != "" && promotionCodeID != "" {
		return nil, invalidParam("discounts[0]", "Specify either coupon or promotion_code, not both.")
	}
	if promotionCodeID != "" {
		discount, err := h.discountFromPromotionCode(promotionCodeID)
		if err != nil {
			return nil, err
		}
		return []billing.Discount{discount}, nil
	}
	if couponID != "" {
		discount, err := h.discountFromCoupon(couponID, "")
		if err != nil {
			return nil, err
		}
		return []billing.Discount{discount}, nil
	}
	return nil, nil
}

func (h *Handler) discountFromPromotionCode(idOrCode string) (billing.Discount, error) {
	idOrCode = strings.TrimSpace(idOrCode)
	h.local.mu.Lock()
	defer h.local.mu.Unlock()
	var promo map[string]any
	for id, candidate := range h.local.promotionCodes {
		if id == idOrCode || strings.EqualFold(fmt.Sprint(candidate["code"]), idOrCode) {
			promo = candidate
			break
		}
	}
	if promo == nil {
		return billing.Discount{}, billing.ErrNotFound
	}
	if active, ok := promo["active"].(bool); ok && !active {
		return billing.Discount{}, fmt.Errorf("%w: promotion_code is inactive", billing.ErrInvalidInput)
	}
	couponID := ""
	if coupon, ok := promo["coupon"].(map[string]any); ok {
		couponID = fmt.Sprint(coupon["id"])
	}
	discount, err := discountFromCouponEvidence(couponID, h.local.coupons[couponID])
	if err != nil {
		return billing.Discount{}, err
	}
	discount.PromotionCodeID = fmt.Sprint(promo["id"])
	discount.Metadata = stringMapFromEvidence(promo["metadata"])
	return discount, nil
}

func (h *Handler) discountFromCoupon(id string, promotionCodeID string) (billing.Discount, error) {
	id = strings.TrimSpace(id)
	h.local.mu.Lock()
	defer h.local.mu.Unlock()
	discount, err := discountFromCouponEvidence(id, h.local.coupons[id])
	if err != nil {
		return billing.Discount{}, err
	}
	discount.PromotionCodeID = promotionCodeID
	return discount, nil
}

func discountFromCouponEvidence(id string, coupon map[string]any) (billing.Discount, error) {
	if coupon == nil {
		return billing.Discount{}, billing.ErrNotFound
	}
	if valid, ok := coupon["valid"].(bool); ok && !valid {
		return billing.Discount{}, fmt.Errorf("%w: coupon is not valid", billing.ErrInvalidInput)
	}
	percentOff, _ := asInt64Evidence(coupon["percent_off"])
	amountOff, _ := asInt64Evidence(coupon["amount_off"])
	createdAt := time.Now().UTC()
	if created, ok := asInt64Evidence(coupon["created"]); ok && created > 0 {
		createdAt = time.Unix(created, 0).UTC()
	}
	return billing.Discount{
		CouponID:   id,
		PercentOff: percentOff,
		AmountOff:  amountOff,
		Currency:   strings.ToLower(evidenceString(coupon["currency"])),
		Duration:   evidenceString(coupon["duration"]),
		Metadata:   stringMapFromEvidence(coupon["metadata"]),
		CreatedAt:  createdAt,
	}, nil
}

func evidenceString(value any) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(value))
}

func stringMapFromEvidence(value any) map[string]string {
	raw, ok := value.(map[string]string)
	if ok {
		return raw
	}
	anyMap, ok := value.(map[string]any)
	if !ok || len(anyMap) == 0 {
		return nil
	}
	out := map[string]string{}
	for key, value := range anyMap {
		out[key] = fmt.Sprint(value)
	}
	return out
}

func (h *Handler) handleCustomerDiscount(w http.ResponseWriter, r *http.Request, customerID string) {
	customer, err := h.billing.GetCustomer(r.Context(), customerID)
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	discounts := billing.DiscountsFromMetadata(customer.Metadata)
	if len(discounts) == 0 {
		writeResult(w, nil, billing.ErrNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, h.stripeDiscount(discounts[0], customer.ID, "", ""))
	case http.MethodDelete:
		customer.Metadata = billing.ClearDiscountMetadata(copyStringMap(customer.Metadata))
		updated, err := h.billing.UpdateCustomer(r.Context(), customer.ID, billing.Customer{Metadata: customer.Metadata})
		if err == nil {
			h.emitGenericWebhook(r, "customer.discount.deleted", discounts[0].ID, h.stripeDiscount(discounts[0], updated.ID, "", ""), webhooks.SourceAPI)
		}
		writeResult(w, stripeDeleted(discounts[0].ID, "discount"), err)
	default:
		h.methodNotAllowed(w, r, "GET, DELETE")
	}
}

func (h *Handler) handleSubscriptionDiscount(w http.ResponseWriter, r *http.Request, subscriptionID string) {
	subscription, err := h.billing.GetSubscription(r.Context(), subscriptionID)
	if err != nil {
		writeResult(w, nil, err)
		return
	}
	discounts := billing.DiscountsFromMetadata(subscription.Metadata)
	if len(discounts) == 0 {
		writeResult(w, nil, billing.ErrNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, h.stripeDiscount(discounts[0], subscription.CustomerID, subscription.ID, ""))
	case http.MethodDelete:
		metadata := billing.ClearDiscountMetadata(copyStringMap(subscription.Metadata))
		updated, err := h.billing.PatchSubscription(r.Context(), subscription.ID, billing.SubscriptionPatch{
			Metadata:        metadata,
			TimelineSource:  "api",
			TimelineMessage: "Subscription discount deleted",
		})
		if err == nil {
			h.emitGenericWebhook(r, "customer.discount.deleted", discounts[0].ID, h.stripeDiscount(discounts[0], updated.CustomerID, updated.ID, ""), webhooks.SourceAPI)
			h.emitSubscriptionWebhook(r, "customer.subscription.updated", updated, webhooks.SourceAPI)
		}
		writeResult(w, stripeDeleted(discounts[0].ID, "discount"), err)
	default:
		h.methodNotAllowed(w, r, "GET, DELETE")
	}
}

func copyStringMap(in map[string]string) map[string]string {
	if in == nil {
		return map[string]string{}
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func (h *Handler) invoicePreviewDiscounts(ctx context.Context, p params, subscription billing.Subscription, customerID string) ([]billing.Discount, error) {
	discounts, err := h.discountsFromParams(p)
	if err != nil || len(discounts) > 0 {
		return discounts, err
	}
	if subscription.ID != "" {
		discounts = billing.DiscountsFromMetadata(subscription.Metadata)
		if len(discounts) > 0 {
			return discounts, nil
		}
	}
	if customerID == "" {
		return nil, nil
	}
	customer, err := h.billing.GetCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}
	return billing.DiscountsFromMetadata(customer.Metadata), nil
}

func (h *Handler) stripeDiscount(discount billing.Discount, customerID string, subscriptionID string, invoiceID string) map[string]any {
	if discount.Object == "" {
		discount.Object = "discount"
	}
	if discount.ID == "" {
		discount.ID = discountID(discount)
	}
	coupon := map[string]any{"id": discount.CouponID, "object": "coupon"}
	if h != nil && discount.CouponID != "" {
		h.local.mu.Lock()
		if stored, ok := h.local.coupons[discount.CouponID]; ok {
			coupon = cloneEvidence(stored)
		}
		h.local.mu.Unlock()
	}
	return map[string]any{
		"id":             discount.ID,
		"object":         "discount",
		"coupon":         coupon,
		"customer":       emptyToNil(customerID),
		"promotion_code": emptyToNil(discount.PromotionCodeID),
		"subscription":   emptyToNil(subscriptionID),
		"invoice":        emptyToNil(invoiceID),
		"invoice_item":   nil,
		"start":          unix(discount.CreatedAt),
		"end":            nil,
	}
}

func firstDiscountObject(h *Handler, discounts []billing.Discount, customerID string, subscriptionID string, invoiceID string) any {
	if len(discounts) == 0 {
		return nil
	}
	return h.stripeDiscount(discounts[0], customerID, subscriptionID, invoiceID)
}

func discountObjects(h *Handler, discounts []billing.Discount, customerID string, subscriptionID string, invoiceID string) []map[string]any {
	out := make([]map[string]any, 0, len(discounts))
	for _, discount := range discounts {
		out = append(out, h.stripeDiscount(discount, customerID, subscriptionID, invoiceID))
	}
	return out
}

func discountAmounts(discounts []billing.Discount, amount int64) []map[string]any {
	if len(discounts) == 0 || amount <= 0 {
		return []map[string]any{}
	}
	return []map[string]any{{
		"amount":   amount,
		"discount": discountID(discounts[0]),
	}}
}

func discountID(discount billing.Discount) string {
	if discount.ID != "" {
		return discount.ID
	}
	source := firstNonEmptyString(discount.PromotionCodeID, discount.CouponID, "discount")
	return "di_" + sanitizeID(source)
}
