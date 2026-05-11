package api

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hckim/billtap/internal/billing"
)

const (
	stripeCodeParamInvalid = "parameter_invalid"
	stripeCodeParamMissing = "parameter_missing"
	stripeCodeParamUnknown = "parameter_unknown"
)

var (
	metadataParamRE            = regexp.MustCompile(`^metadata\[[^\]]+\]$`)
	expandParamRE              = regexp.MustCompile(`^expand(\[[^\]]*\])?$`)
	enabledEventsParamRE       = regexp.MustCompile(`^enabled_events(\[[^\]]*\])?$`)
	retryBackoffParamRE        = regexp.MustCompile(`^retry_backoff(\[[^\]]*\])?$`)
	checkoutLineItemRE         = regexp.MustCompile(`^line_items\[(\d+)\]\[(price|quantity)\]$`)
	legacyLineItemParamRE      = regexp.MustCompile(`^lineItems\[(\d+)\]\[(price|quantity)\]$`)
	checkoutSubscriptionDataRE = regexp.MustCompile(`^subscription_data\[(trial_period_days)\]$`)
	subscriptionItemRE         = regexp.MustCompile(`^items\[(\d+)\]\[(id|price|price_id|quantity)\]$`)
	cancellationDetailsRE      = regexp.MustCompile(`^cancellation_details\[(comment|feedback)\]$`)
	paymentMethodTypesRE       = regexp.MustCompile(`^payment_method_types(\[[^\]]*\])?$`)
	automaticPaymentMethodsRE  = regexp.MustCompile(`^automatic_payment_methods\[(enabled)\]$`)
	paymentMethodOptionsRE     = regexp.MustCompile(`^payment_method_options\[.+\]$`)
)

type validationError struct {
	Param   string
	Code    string
	Message string
}

func (e *validationError) Error() string {
	return e.Message
}

func missingParam(param string) *validationError {
	return &validationError{
		Param:   param,
		Code:    stripeCodeParamMissing,
		Message: fmt.Sprintf("Missing required param: %s.", param),
	}
}

func unknownParam(param string) *validationError {
	return &validationError{
		Param:   param,
		Code:    stripeCodeParamUnknown,
		Message: fmt.Sprintf("Received unknown parameter: %s.", param),
	}
}

func invalidParam(param string, reason string) *validationError {
	return &validationError{
		Param:   param,
		Code:    stripeCodeParamInvalid,
		Message: fmt.Sprintf("Invalid param: %s. %s", param, reason),
	}
}

type paramSpec struct {
	Allowed       []string
	AllowedRegex  []*regexp.Regexp
	Required      []string
	RequiredAny   [][]string
	Int64Params   []string
	BoolParams    []string
	EnumParams    map[string][]string
	NonNegative   []string
	Positive      []string
	AllowMetadata bool
}

func (p params) validate(spec paramSpec) error {
	allowed := map[string]struct{}{}
	for _, key := range spec.Allowed {
		allowed[key] = struct{}{}
	}

	for key := range p.values {
		if _, ok := allowed[key]; ok {
			continue
		}
		if expandParamRE.MatchString(key) {
			continue
		}
		if spec.AllowMetadata && metadataParamRE.MatchString(key) {
			continue
		}
		matched := false
		for _, pattern := range spec.AllowedRegex {
			if pattern.MatchString(key) {
				matched = true
				break
			}
		}
		if !matched {
			return unknownParam(key)
		}
	}

	for _, key := range spec.Required {
		if !p.has(key) {
			return missingParam(key)
		}
	}
	for _, group := range spec.RequiredAny {
		if !p.hasAny(group...) {
			return missingParam(group[0])
		}
	}
	for _, key := range spec.Int64Params {
		if err := p.validateInt64(key); err != nil {
			return err
		}
	}
	for _, key := range spec.BoolParams {
		if err := p.validateBool(key); err != nil {
			return err
		}
	}
	for key, values := range spec.EnumParams {
		if err := p.validateEnum(key, values); err != nil {
			return err
		}
	}
	for _, key := range spec.NonNegative {
		if err := p.validateMin(key, 0); err != nil {
			return err
		}
	}
	for _, key := range spec.Positive {
		if err := p.validateMin(key, 1); err != nil {
			return err
		}
	}
	return nil
}

func (p params) has(key string) bool {
	value, ok := p.values[key]
	return ok && strings.TrimSpace(value) != ""
}

func (p params) hasAny(keys ...string) bool {
	for _, key := range keys {
		if p.has(key) {
			return true
		}
	}
	return false
}

func (p params) validateInt64(key string) error {
	if !p.has(key) {
		return nil
	}
	if _, err := strconv.ParseInt(p.string(key), 10, 64); err != nil {
		return invalidParam(key, "Expected an integer.")
	}
	return nil
}

func (p params) validateMin(key string, min int64) error {
	if !p.has(key) {
		return nil
	}
	value, err := strconv.ParseInt(p.string(key), 10, 64)
	if err != nil {
		return invalidParam(key, "Expected an integer.")
	}
	if value < min {
		return invalidParam(key, fmt.Sprintf("Must be at least %d.", min))
	}
	return nil
}

func (p params) validateBool(key string) error {
	if !p.has(key) {
		return nil
	}
	switch p.string(key) {
	case "true", "false", "1", "0":
		return nil
	default:
		return invalidParam(key, "Expected a boolean.")
	}
}

func (p params) validateEnum(key string, allowed []string) error {
	if !p.has(key) {
		return nil
	}
	value := p.string(key)
	for _, candidate := range allowed {
		if value == candidate {
			return nil
		}
	}
	return invalidParam(key, "Expected one of: "+strings.Join(allowed, ", ")+".")
}

func (p params) validateUnixTimestampOrNow(key string) error {
	if !p.has(key) {
		return nil
	}
	if p.string(key) == "now" {
		return nil
	}
	if _, err := strconv.ParseInt(p.string(key), 10, 64); err != nil {
		return invalidParam(key, "Expected a Unix timestamp or now.")
	}
	return nil
}

func validateCustomerCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"id", "email", "name", "test_clock"},
		AllowMetadata: true,
	})
}

func validateCustomerUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"email", "name", "test_clock"},
		AllowMetadata: true,
	})
}

func validateProductCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"id", "name", "description", "active"},
		Required:      []string{"name"},
		BoolParams:    []string{"active"},
		AllowMetadata: true,
	})
}

func validateProductUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"name", "description", "active"},
		BoolParams:    []string{"active"},
		AllowMetadata: true,
	})
}

func validatePriceCreate(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"id",
			"product",
			"product_id",
			"currency",
			"unit_amount",
			"lookup_key",
			"lookupKey",
			"recurring[interval]",
			"recurring_interval",
			"interval",
			"recurring[interval_count]",
			"active",
		},
		Required:    []string{"currency", "unit_amount"},
		RequiredAny: [][]string{{"product", "product_id"}},
		Int64Params: []string{"unit_amount", "recurring[interval_count]"},
		NonNegative: []string{"unit_amount"},
		BoolParams:  []string{"active"},
		EnumParams: map[string][]string{
			"recurring[interval]": {"day", "week", "month", "year"},
			"recurring_interval":  {"day", "week", "month", "year"},
			"interval":            {"day", "week", "month", "year"},
		},
		AllowMetadata: true,
	})
}

func validatePriceUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"product",
			"product_id",
			"currency",
			"unit_amount",
			"lookup_key",
			"lookupKey",
			"recurring[interval]",
			"recurring_interval",
			"interval",
			"recurring[interval_count]",
			"active",
		},
		Int64Params: []string{"unit_amount", "recurring[interval_count]"},
		NonNegative: []string{"unit_amount"},
		BoolParams:  []string{"active"},
		EnumParams: map[string][]string{
			"recurring[interval]": {"day", "week", "month", "year"},
			"recurring_interval":  {"day", "week", "month", "year"},
			"interval":            {"day", "week", "month", "year"},
		},
		AllowMetadata: true,
	})
}

func validateCheckoutSessionCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed:      []string{"customer", "customer_id", "mode", "success_url", "cancel_url", "price", "allow_promotion_codes"},
		AllowedRegex: []*regexp.Regexp{checkoutLineItemRE, legacyLineItemParamRE, checkoutSubscriptionDataRE},
		RequiredAny:  [][]string{{"customer", "customer_id"}},
		Int64Params:  []string{"subscription_data[trial_period_days]"},
		BoolParams:   []string{"allow_promotion_codes"},
		EnumParams:   map[string][]string{"mode": {"subscription"}},
		Positive:     []string{"subscription_data[trial_period_days]"},
	}); err != nil {
		return err
	}
	lineItemIndexes := p.lineItemIndexes()
	if len(lineItemIndexes) == 0 && !p.has("price") {
		return missingParam("line_items")
	}
	for idx := range lineItemIndexes {
		quantityKey := fmt.Sprintf("line_items[%d][quantity]", idx)
		if err := p.validateMin(quantityKey, 1); err != nil {
			return err
		}
		if p.has(quantityKey) && !p.has(fmt.Sprintf("line_items[%d][price]", idx)) {
			return missingParam(fmt.Sprintf("line_items[%d][price]", idx))
		}
	}
	return nil
}

func validateSubscriptionCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed: []string{
			"customer",
			"customer_id",
			"price",
			"price_id",
			"collection_method",
			"days_until_due",
			"cancel_at",
			"billing_cycle_anchor",
			"outcome",
			"test_clock",
		},
		AllowedRegex: []*regexp.Regexp{subscriptionItemRE},
		RequiredAny:  [][]string{{"customer", "customer_id"}},
		Int64Params:  []string{"days_until_due", "cancel_at", "billing_cycle_anchor"},
		Positive:     []string{"days_until_due"},
		EnumParams: map[string][]string{
			"collection_method": {"charge_automatically", "send_invoice"},
		},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	itemIndexes := p.subscriptionItemIndexes()
	if len(itemIndexes) == 0 && !p.hasAny("price", "price_id") {
		return missingParam("items")
	}
	for idx := range itemIndexes {
		quantityKey := fmt.Sprintf("items[%d][quantity]", idx)
		if err := p.validateMin(quantityKey, 1); err != nil {
			return err
		}
		if p.has(quantityKey) && !p.hasAny(fmt.Sprintf("items[%d][price]", idx), fmt.Sprintf("items[%d][price_id]", idx)) {
			return missingParam(fmt.Sprintf("items[%d][price]", idx))
		}
	}
	return nil
}

func validateSubscriptionUpdate(p params) error {
	if err := p.validate(paramSpec{
		Allowed:      []string{"cancel_at_period_end", "proration_behavior", "payment_behavior", "billing_cycle_anchor", "trial_end"},
		AllowedRegex: []*regexp.Regexp{subscriptionItemRE, cancellationDetailsRE},
		BoolParams:   []string{"cancel_at_period_end"},
		EnumParams: map[string][]string{
			"proration_behavior": {"none", "create_prorations", "always_invoice"},
			"payment_behavior":   {"allow_incomplete", "error_if_incomplete", "pending_if_incomplete"},
			"cancellation_details[feedback]": {
				"customer_service",
				"low_quality",
				"missing_features",
				"switched_service",
				"too_complex",
				"too_expensive",
				"unused",
				"other",
			},
		},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	for idx := range p.subscriptionItemIndexes() {
		if err := p.validateMin(fmt.Sprintf("items[%d][quantity]", idx), 1); err != nil {
			return err
		}
	}
	for _, key := range []string{"billing_cycle_anchor", "trial_end"} {
		if err := p.validateUnixTimestampOrNow(key); err != nil {
			return err
		}
	}
	return nil
}

func validateSubscriptionItemCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:     []string{"subscription", "price", "price_id", "quantity"},
		Required:    []string{"subscription"},
		RequiredAny: [][]string{{"price", "price_id"}},
		Int64Params: []string{"quantity"},
		Positive:    []string{"quantity"},
	})
}

func validatePaymentIntentCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed: []string{
			"id",
			"amount",
			"currency",
			"customer",
			"payment_method",
			"confirm",
			"capture_method",
			"outcome",
			"description",
			"receipt_email",
			"setup_future_usage",
			"off_session",
			"automatic_payment_methods[enabled]",
		},
		AllowedRegex: []*regexp.Regexp{paymentMethodTypesRE, paymentMethodOptionsRE},
		Required:     []string{"amount", "currency"},
		Int64Params:  []string{"amount"},
		Positive:     []string{"amount"},
		BoolParams:   []string{"confirm", "off_session", "automatic_payment_methods[enabled]"},
		EnumParams: map[string][]string{
			"capture_method":     {"automatic", "automatic_async", "manual"},
			"setup_future_usage": {"on_session", "off_session"},
		},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	if p.boolDefault("confirm", false) && !p.hasAny("payment_method", "outcome") {
		return missingParam("payment_method")
	}
	return nil
}

func validateInvoicePay(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"forgive",
			"mandate",
			"off_session",
			"paid_out_of_band",
			"payment_method",
			"source",
		},
		BoolParams: []string{"forgive", "off_session", "paid_out_of_band"},
	})
}

func validateRefundCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed: []string{
			"id",
			"charge",
			"payment_intent",
			"invoice",
			"customer",
			"amount",
			"currency",
			"reason",
		},
		RequiredAny:   [][]string{{"charge", "payment_intent", "invoice"}},
		Int64Params:   []string{"amount"},
		Positive:      []string{"amount"},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	if !p.has("amount") {
		return missingParam("amount")
	}
	return nil
}

func validateCreditNoteCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"id", "invoice", "customer", "amount", "currency", "reason"},
		Required:      []string{"invoice", "amount"},
		Int64Params:   []string{"amount"},
		Positive:      []string{"amount"},
		AllowMetadata: true,
	})
}

func validatePaymentIntentConfirm(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"payment_method",
			"outcome",
			"return_url",
			"receipt_email",
			"setup_future_usage",
			"off_session",
			"use_stripe_sdk",
		},
		AllowedRegex: []*regexp.Regexp{paymentMethodOptionsRE},
		BoolParams:   []string{"off_session", "use_stripe_sdk"},
		EnumParams: map[string][]string{
			"setup_future_usage": {"on_session", "off_session"},
		},
	})
}

func validatePaymentIntentCapture(p params) error {
	return p.validate(paramSpec{
		Allowed:     []string{"amount_to_capture"},
		Int64Params: []string{"amount_to_capture"},
		Positive:    []string{"amount_to_capture"},
	})
}

func validatePaymentIntentCancel(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{"cancellation_reason"},
		EnumParams: map[string][]string{
			"cancellation_reason": {"duplicate", "fraudulent", "requested_by_customer", "abandoned"},
		},
	})
}

func validateSetupIntentCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed: []string{
			"id",
			"customer",
			"payment_method",
			"confirm",
			"usage",
			"outcome",
			"description",
			"return_url",
			"automatic_payment_methods[enabled]",
		},
		AllowedRegex: []*regexp.Regexp{paymentMethodTypesRE, paymentMethodOptionsRE},
		BoolParams:   []string{"confirm", "automatic_payment_methods[enabled]"},
		EnumParams: map[string][]string{
			"usage": {"on_session", "off_session"},
		},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	if p.boolDefault("confirm", false) && !p.hasAny("payment_method", "outcome") {
		return missingParam("payment_method")
	}
	return nil
}

func validateSetupIntentConfirm(p params) error {
	return p.validate(paramSpec{
		Allowed:      []string{"payment_method", "outcome", "return_url", "use_stripe_sdk"},
		AllowedRegex: []*regexp.Regexp{paymentMethodOptionsRE},
		BoolParams:   []string{"use_stripe_sdk"},
	})
}

func validateSetupIntentCancel(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{"cancellation_reason"},
		EnumParams: map[string][]string{
			"cancellation_reason": {"duplicate", "fraudulent", "requested_by_customer", "abandoned"},
		},
	})
}

func validateTestClockCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed:     []string{"id", "name", "frozen_time", "frozenTime"},
		RequiredAny: [][]string{{"frozen_time", "frozenTime"}},
	}); err != nil {
		return err
	}
	return nil
}

func validateTestClockAdvance(p params) error {
	return p.validate(paramSpec{
		Allowed:     []string{"frozen_time", "frozenTime"},
		RequiredAny: [][]string{{"frozen_time", "frozenTime"}},
	})
}

func validateBillingPortalSessionCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:     []string{"customer", "customer_id", "return_url"},
		RequiredAny: [][]string{{"customer", "customer_id"}},
	})
}

func validateWebhookEndpointCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:      []string{"url", "secret", "retry_max_attempts"},
		AllowedRegex: []*regexp.Regexp{enabledEventsParamRE, retryBackoffParamRE},
		Required:     []string{"url"},
		Int64Params:  []string{"retry_max_attempts"},
		Positive:     []string{"retry_max_attempts"},
	})
}

func validateWebhookEndpointUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed:      []string{"url", "secret", "retry_max_attempts", "active", "enabled"},
		AllowedRegex: []*regexp.Regexp{enabledEventsParamRE, retryBackoffParamRE},
		Int64Params:  []string{"retry_max_attempts"},
		Positive:     []string{"retry_max_attempts"},
		BoolParams:   []string{"active", "enabled"},
	})
}

func (p params) lineItemIndexes() map[int]struct{} {
	indexes := map[int]struct{}{}
	for key := range p.values {
		for _, pattern := range []*regexp.Regexp{checkoutLineItemRE, legacyLineItemParamRE} {
			matches := pattern.FindStringSubmatch(key)
			if len(matches) != 3 {
				continue
			}
			idx, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}
			indexes[idx] = struct{}{}
		}
	}
	return indexes
}

func (p params) subscriptionItemIndexes() map[int]struct{} {
	indexes := map[int]struct{}{}
	for key := range p.values {
		matches := subscriptionItemRE.FindStringSubmatch(key)
		if len(matches) != 3 {
			continue
		}
		idx, err := strconv.Atoi(matches[1])
		if err != nil {
			continue
		}
		indexes[idx] = struct{}{}
	}
	return indexes
}

func validateProductExists(product billing.Product, err error) error {
	if err != nil {
		return err
	}
	if strings.TrimSpace(product.ID) == "" {
		return billing.ErrNotFound
	}
	return nil
}

func validatePriceExists(price billing.Price, err error) error {
	if err != nil {
		return err
	}
	if strings.TrimSpace(price.ID) == "" {
		return billing.ErrNotFound
	}
	return nil
}

func validateCustomerExists(customer billing.Customer, err error) error {
	if err != nil {
		return err
	}
	if strings.TrimSpace(customer.ID) == "" {
		return billing.ErrNotFound
	}
	return nil
}
