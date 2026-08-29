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
	metadataParamRE       = regexp.MustCompile(`^metadata\[[^\]]+\]$`)
	expandParamRE         = regexp.MustCompile(`^expand(\[[^\]]*\])?$`)
	enabledEventsParamRE  = regexp.MustCompile(`^enabled_events(\[[^\]]*\])?$`)
	retryBackoffParamRE   = regexp.MustCompile(`^retry_backoff(\[[^\]]*\])?$`)
	checkoutLineItemRE    = regexp.MustCompile(`^line_items\[(\d+)\]\[(price|quantity)\]$`)
	legacyLineItemParamRE = regexp.MustCompile(`^lineItems\[(\d+)\]\[(price|quantity)\]$`)
	// price_data nested keys for payment-mode one-off line items (Stripe Checkout Sessions).
	checkoutPriceDataRE            = regexp.MustCompile(`^line_items\[(\d+)\]\[price_data\]\[(currency|unit_amount|unit_amount_decimal|tax_behavior|product)\]$`)
	checkoutPriceDataProductDataRE = regexp.MustCompile(`^line_items\[(\d+)\]\[price_data\]\[product_data\]\[(name|description)\]$`)
	checkoutPriceDataProductMetaRE = regexp.MustCompile(`^line_items\[(\d+)\]\[price_data\]\[product_data\]\[metadata\]\[[^\]]+\]$`)
	checkoutPriceDataRecurringRE   = regexp.MustCompile(`^line_items\[(\d+)\]\[price_data\]\[recurring\]\[(interval|interval_count)\]$`)
	// payment_intent_data for payment-mode sessions (nested metadata keys unexpanded by firstValues).
	checkoutPaymentIntentDataRE = regexp.MustCompile(`^payment_intent_data\[(setup_future_usage|description|receipt_email|capture_method)\]$`)
	checkoutPaymentIntentMetaRE = regexp.MustCompile(`^payment_intent_data\[metadata\]\[[^\]]+\]$`)
	// Accept trial_period_days and default_tax_rates in [] / [N] forms
	// (firstValues leaves single-value [] keys unexpanded).
	checkoutSubscriptionDataRE = regexp.MustCompile(`^subscription_data\[(trial_period_days|default_tax_rates)\](\[\d*\])?$`)
	defaultTaxRatesParamRE     = regexp.MustCompile(`^default_tax_rates(\[\d*\])?$`)
	// Item-level tax_rates on subscription_items create: evidence only (not used for totals).
	taxRatesParamRE           = regexp.MustCompile(`^tax_rates(\[\d*\])?$`)
	discountParamRE           = regexp.MustCompile(`^discounts\[\d+\]\[(coupon|promotion_code)\]$`)
	subscriptionItemRE        = regexp.MustCompile(`^items\[(\d+)\]\[(id|price|price_id|quantity)\]$`)
	cancellationDetailsRE     = regexp.MustCompile(`^cancellation_details\[(comment|feedback)\]$`)
	paymentMethodTypesRE      = regexp.MustCompile(`^payment_method_types(\[[^\]]*\])?$`)
	automaticPaymentMethodsRE = regexp.MustCompile(`^automatic_payment_methods\[(enabled)\]$`)
	paymentMethodOptionsRE    = regexp.MustCompile(`^payment_method_options\[.+\]$`)
	accountNestedParamRE      = regexp.MustCompile(`^(business_profile|company|individual|settings|tos_acceptance|controller)\[.+\]$`)
	eventTypeFilterParamRE    = regexp.MustCompile(`^(type|types|event_type|event_types)(\[[^\]]*\])?$`)
	portalFlowDataParamRE     = regexp.MustCompile(`^flow_data(\[[^\]]+\])+$`)
	// Allow empty brackets for Stripe array form applies_to[products][].
	couponAppliesToParamRE      = regexp.MustCompile(`^applies_to(\[[^\]]*\])+$`)
	promotionRestrictionParamRE = regexp.MustCompile(`^restrictions(\[[^\]]+\])+$`)
	schedulePhaseParamRE        = regexp.MustCompile(`^phases\[\d+\]\[(start_date|end_date|iterations|items|plans)\].*$`)
	invoicePreviewItemParamRE   = regexp.MustCompile(`^((subscription_details|subscriptionDetails)\[items\]\[\d+\]\[(id|price|price_id|quantity)\]|(subscription_items|items)\[\d+\]\[(id|price|price_id|quantity)\])$`)
	invoicePaymentSettingsRE    = regexp.MustCompile(`^payment_settings(\[[^\]]+\])+$`)
	portalFeatureParamRE        = regexp.MustCompile(`^features\[(customer_update|invoice_history|payment_method_update|subscription_cancel|subscription_update)\]\[(enabled|mode|proration_behavior|cancellation_reason|allowed_updates|default_allowed_updates)\](\[\d*\])?$`)
)

var stripePaymentMethodTypes = []string{
	"acss_debit",
	"affirm",
	"afterpay_clearpay",
	"alipay",
	"alma",
	"amazon_pay",
	"au_becs_debit",
	"bacs_debit",
	"bancontact",
	"billie",
	"blik",
	"boleto",
	"card",
	"cashapp",
	"crypto",
	"custom",
	"customer_balance",
	"eps",
	"fpx",
	"giropay",
	"grabpay",
	"ideal",
	"kakao_pay",
	"klarna",
	"konbini",
	"kr_card",
	"link",
	"mb_way",
	"mobilepay",
	"multibanco",
	"naver_pay",
	"nz_bank_account",
	"oxxo",
	"p24",
	"pay_by_bank",
	"payco",
	"paynow",
	"paypal",
	"payto",
	"pix",
	"promptpay",
	"revolut_pay",
	"samsung_pay",
	"satispay",
	"sepa_debit",
	"sofort",
	"sunbit",
	"swish",
	"twint",
	"upi",
	"us_bank_account",
	"wechat_pay",
	"zip",
}

var stripePortalLocales = []string{
	"auto",
	"bg",
	"cs",
	"da",
	"de",
	"el",
	"en",
	"en-AU",
	"en-CA",
	"en-GB",
	"en-IE",
	"en-IN",
	"en-NZ",
	"en-SG",
	"es",
	"es-419",
	"et",
	"fi",
	"fil",
	"fr",
	"fr-CA",
	"hr",
	"hu",
	"id",
	"it",
	"ja",
	"ko",
	"lt",
	"lv",
	"ms",
	"mt",
	"nb",
	"nl",
	"pl",
	"pt",
	"pt-BR",
	"ro",
	"ru",
	"sk",
	"sl",
	"sv",
	"th",
	"tr",
	"vi",
	"zh",
	"zh-HK",
	"zh-TW",
}

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

func exclusiveParams(param string, names ...string) *validationError {
	return &validationError{
		Param:   param,
		Code:    stripeCodeParamInvalid,
		Message: "You may only specify one of these parameters: " + strings.Join(names, ", "),
	}
}

type paramSpec struct {
	Allowed       []string
	AllowedRegex  []*regexp.Regexp
	Required      []string
	RequiredAny   [][]string
	Int64Params   []string
	FloatParams   []string
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
	for _, key := range spec.FloatParams {
		if err := p.validateFloat64(key); err != nil {
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

func (p params) validateFloat64(key string) error {
	if !p.has(key) {
		return nil
	}
	if _, err := strconv.ParseFloat(p.string(key), 64); err != nil {
		// Same invalid-error shape as Int64Params ("Expected an integer.").
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

// validatePositiveFloat enforces value > 0 for float form params, using the same
// error message shape as Positive/validateMin(min=1).
func (p params) validatePositiveFloat(key string) error {
	if !p.has(key) {
		return nil
	}
	value, err := strconv.ParseFloat(p.string(key), 64)
	if err != nil {
		return invalidParam(key, "Expected an integer.")
	}
	if value <= 0 {
		return invalidParam(key, "Must be at least 1.")
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

func (p params) validateUnixTimestampOrNowOrUnchanged(key string) error {
	if !p.has(key) {
		return nil
	}
	switch p.string(key) {
	case "now", "unchanged":
		return nil
	}
	if _, err := strconv.ParseInt(p.string(key), 10, 64); err != nil {
		return invalidParam(key, "Expected a Unix timestamp, now, or unchanged.")
	}
	return nil
}

func validateCustomerCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"id", "email", "name", "test_clock", "coupon", "promotion_code"},
		AllowedRegex:  []*regexp.Regexp{discountParamRE},
		AllowMetadata: true,
	})
}

func validateCustomerUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"email", "name", "test_clock", "coupon", "promotion_code", "invoice_settings[default_payment_method]"},
		AllowedRegex:  []*regexp.Regexp{discountParamRE},
		AllowMetadata: true,
	})
}

func validateCustomerSearch(p params) error {
	return validateSearch(p)
}

func validateSearch(p params) error {
	return p.validate(paramSpec{
		Allowed:     []string{"query", "limit", "page"},
		Required:    []string{"query"},
		Int64Params: []string{"limit"},
		Positive:    []string{"limit"},
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

func validatePriceSearch(p params) error {
	if err := p.validate(paramSpec{
		Allowed:     []string{"query", "limit", "page"},
		Required:    []string{"query"},
		Int64Params: []string{"limit"},
		Positive:    []string{"limit"},
	}); err != nil {
		return err
	}
	_, err := parsePriceSearchQuery(p.string("query"))
	return err
}

func validateAccountCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed:       []string{"id", "type", "country", "email", "business_type", "default_currency", "details_submitted", "account_token", "external_account"},
		AllowedRegex:  []*regexp.Regexp{capabilityParamPattern, accountNestedParamRE},
		BoolParams:    []string{"details_submitted"},
		EnumParams:    map[string][]string{"type": {"standard", "express", "custom"}, "business_type": {"individual", "company", "non_profit", "government_entity"}},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	return validateCapabilityParams(p)
}

func validateAccountUpdate(p params) error {
	if err := p.validate(paramSpec{
		Allowed:       []string{"email", "business_type", "default_currency", "account_token", "external_account"},
		AllowedRegex:  []*regexp.Regexp{capabilityParamPattern, accountNestedParamRE},
		EnumParams:    map[string][]string{"business_type": {"individual", "company", "non_profit", "government_entity"}},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	return validateCapabilityParams(p)
}

func validateCapabilityParams(p params) error {
	for key := range p.values {
		matches := capabilityParamPattern.FindStringSubmatch(key)
		if len(matches) != 3 {
			continue
		}
		switch matches[2] {
		case "requested":
			if err := p.validateBool(key); err != nil {
				return err
			}
		case "status":
			if err := p.validateEnum(key, []string{"active", "inactive", "pending"}); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateAccountLinkCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:    []string{"account", "refresh_url", "return_url", "type", "collection_options[fields]"},
		Required:   []string{"account", "refresh_url", "return_url", "type"},
		EnumParams: map[string][]string{"type": {"account_onboarding", "account_update"}},
	})
}

func validateAccountSessionCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed:      []string{"account"},
		AllowedRegex: []*regexp.Regexp{accountSessionComponentPattern},
		Required:     []string{"account"},
	}); err != nil {
		return err
	}
	hasComponent := false
	for key := range p.values {
		if accountSessionComponentPattern.MatchString(key) {
			hasComponent = true
			if err := p.validateBool(key); err != nil {
				return err
			}
		}
	}
	if !hasComponent {
		return missingParam("components")
	}
	return nil
}

func validateAccountCapabilityUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed:    []string{"requested", "status"},
		BoolParams: []string{"requested"},
		EnumParams: map[string][]string{"status": {"active", "inactive", "pending"}},
	})
}

func validateExternalAccountCreate(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"id",
			"external_account",
			"token",
			"country",
			"currency",
			"bank_name",
			"routing_number",
			"account_number",
			"account_holder_name",
			"account_holder_type",
			"default_for_currency",
		},
		BoolParams:    []string{"default_for_currency"},
		RequiredAny:   [][]string{{"external_account", "token", "account_number"}},
		AllowMetadata: true,
	})
}

func validateExternalAccountUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"bank_name", "account_holder_name", "account_holder_type", "default_for_currency"},
		BoolParams:    []string{"default_for_currency"},
		AllowMetadata: true,
	})
}

func validateAccountReject(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{"reason"},
		EnumParams: map[string][]string{
			"reason": {"fraud", "terms_of_service", "other"},
		},
	})
}

func validateAccountPersonCreate(p params) error {
	return validateAccountPerson(p, []string{"id"})
}

func validateAccountPersonUpdate(p params) error {
	return validateAccountPerson(p, nil)
}

func validateAccountPerson(p params, allowed []string) error {
	return p.validate(paramSpec{
		Allowed:       allowed,
		AllowedRegex:  []*regexp.Regexp{personDataParamPattern},
		BoolParams:    []string{"relationship[owner]", "relationship[director]", "relationship[executive]", "relationship[representative]"},
		Int64Params:   []string{"dob[day]", "dob[month]", "dob[year]", "relationship[percent_ownership]"},
		NonNegative:   []string{"relationship[percent_ownership]"},
		AllowMetadata: true,
	})
}

func validateTransferCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"id", "amount", "currency", "destination", "source_transaction", "description", "transfer_group"},
		Required:      []string{"amount", "currency", "destination"},
		Int64Params:   []string{"amount"},
		Positive:      []string{"amount"},
		AllowMetadata: true,
	})
}

func validateTransferUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"description"},
		AllowMetadata: true,
	})
}

func validateTransferReversalCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"id", "amount", "refund_application_fee"},
		Int64Params:   []string{"amount"},
		Positive:      []string{"amount"},
		BoolParams:    []string{"refund_application_fee"},
		AllowMetadata: true,
	})
}

func validateTransferReversalUpdate(p params) error {
	return p.validate(paramSpec{AllowMetadata: true})
}

func validatePayoutCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"id", "amount", "currency", "destination", "method", "description", "statement_descriptor"},
		Required:      []string{"amount", "currency"},
		Int64Params:   []string{"amount"},
		Positive:      []string{"amount"},
		EnumParams:    map[string][]string{"method": {"standard", "instant"}},
		AllowMetadata: true,
	})
}

func validatePayoutUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"description"},
		AllowMetadata: true,
	})
}

func validateApplicationFeeRefundCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"id", "amount", "charge"},
		Int64Params:   []string{"amount"},
		Positive:      []string{"amount"},
		AllowMetadata: true,
	})
}

func validateApplicationFeeRefundUpdate(p params) error {
	return p.validate(paramSpec{AllowMetadata: true})
}

func validateCheckoutSessionCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed: []string{
			"customer",
			"customer_id",
			"mode",
			"success_url",
			"cancel_url",
			"price",
			"allow_promotion_codes",
			"coupon",
			"promotion_code",
			"automatic_tax[enabled]",
			"tax_id_collection[enabled]",
			"client_reference_id",
		},
		AllowedRegex: []*regexp.Regexp{
			checkoutLineItemRE,
			legacyLineItemParamRE,
			checkoutPriceDataRE,
			checkoutPriceDataProductDataRE,
			checkoutPriceDataProductMetaRE,
			checkoutPriceDataRecurringRE,
			checkoutPaymentIntentDataRE,
			checkoutPaymentIntentMetaRE,
			checkoutSubscriptionDataRE,
			discountParamRE,
		},
		RequiredAny:   [][]string{{"customer", "customer_id"}},
		AllowMetadata: true,
		Int64Params:   []string{"subscription_data[trial_period_days]"},
		BoolParams:    []string{"allow_promotion_codes", "automatic_tax[enabled]", "tax_id_collection[enabled]"},
		EnumParams: map[string][]string{
			"mode": {"subscription", "payment"},
			"payment_intent_data[setup_future_usage]": {"off_session", "on_session"},
			"payment_intent_data[capture_method]":     {"automatic", "automatic_async", "manual"},
		},
		Positive: []string{"subscription_data[trial_period_days]"},
	}); err != nil {
		return err
	}
	mode := p.stringDefault("mode", "subscription")
	if err := validateCheckoutModeCrossParams(p, mode); err != nil {
		return err
	}
	if mode == "subscription" {
		if err := validateDefaultTaxRatesVsAutomaticTax(p, "subscription_data[default_tax_rates]", p.boolDefault("automatic_tax[enabled]", false)); err != nil {
			return err
		}
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
		priceKey := fmt.Sprintf("line_items[%d][price]", idx)
		hasPrice := p.has(priceKey)
		hasPriceData := p.hasPriceData(idx)
		if mode == "payment" {
			if !hasPrice && !hasPriceData {
				if p.has(quantityKey) {
					return missingParam(priceKey)
				}
				return missingParam(fmt.Sprintf("line_items[%d][price]", idx))
			}
			if hasPriceData {
				if err := validateCheckoutPriceData(p, idx, mode); err != nil {
					return err
				}
			}
			continue
		}
		// subscription mode: price required when quantity is present; price_data stays unknown-style reject.
		if hasPriceData {
			return unknownParam(fmt.Sprintf("line_items[%d][price_data]", idx))
		}
		if p.has(quantityKey) && !hasPrice {
			return missingParam(priceKey)
		}
	}
	return nil
}

func validateCheckoutModeCrossParams(p params, mode string) error {
	hasPaymentIntentData := false
	hasSubscriptionData := false
	for key := range p.values {
		if strings.HasPrefix(key, "payment_intent_data[") {
			hasPaymentIntentData = true
		}
		if strings.HasPrefix(key, "subscription_data[") {
			hasSubscriptionData = true
		}
	}
	if mode == "subscription" && hasPaymentIntentData {
		return invalidParam("payment_intent_data", "You cannot use payment_intent_data in subscription mode.")
	}
	if mode == "payment" && hasSubscriptionData {
		return invalidParam("subscription_data", "You cannot use subscription_data in payment mode.")
	}
	return nil
}

func validateCheckoutPriceData(p params, idx int, mode string) error {
	prefix := fmt.Sprintf("line_items[%d][price_data]", idx)
	currencyKey := prefix + "[currency]"
	if !p.has(currencyKey) {
		return missingParam(currencyKey)
	}
	unitAmountKey := prefix + "[unit_amount]"
	unitAmountDecimalKey := prefix + "[unit_amount_decimal]"
	hasUnitAmount := p.has(unitAmountKey)
	hasUnitAmountDecimal := p.has(unitAmountDecimalKey)
	if hasUnitAmount == hasUnitAmountDecimal {
		if hasUnitAmount {
			return invalidParam(unitAmountKey, "Only one of unit_amount or unit_amount_decimal may be provided.")
		}
		return missingParam(unitAmountKey)
	}
	if hasUnitAmount {
		if err := p.validateMin(unitAmountKey, 0); err != nil {
			return err
		}
	}
	if hasUnitAmountDecimal {
		if _, err := parseUnitAmountDecimal(p.string(unitAmountDecimalKey)); err != nil {
			return invalidParam(unitAmountDecimalKey, "Must be an integer amount in the smallest currency unit.")
		}
	}
	if err := p.validateEnum(prefix+"[tax_behavior]", []string{"exclusive", "inclusive", "unspecified"}); err != nil {
		return err
	}
	productKey := prefix + "[product]"
	productDataNameKey := prefix + "[product_data][name]"
	hasProduct := p.has(productKey)
	hasProductData := p.has(productDataNameKey) || p.hasPriceDataProductData(idx)
	if hasProduct == hasProductData {
		if hasProduct {
			return invalidParam(productKey, "Only one of product or product_data may be provided.")
		}
		return missingParam(productDataNameKey)
	}
	if hasProductData && !p.has(productDataNameKey) {
		return missingParam(productDataNameKey)
	}
	if p.has(prefix+"[recurring][interval]") || p.has(prefix+"[recurring][interval_count]") {
		if mode == "payment" {
			return invalidParam(prefix+"[recurring]", "Recurring price_data is not supported in payment mode.")
		}
	}
	if p.has(prefix + "[recurring][interval_count]") {
		if err := p.validateMin(prefix+"[recurring][interval_count]", 1); err != nil {
			return err
		}
	}
	return nil
}

func (p params) hasPriceData(idx int) bool {
	prefix := fmt.Sprintf("line_items[%d][price_data]", idx)
	for key := range p.values {
		if strings.HasPrefix(key, prefix+"[") {
			return true
		}
	}
	return false
}

func (p params) hasPriceDataProductData(idx int) bool {
	prefix := fmt.Sprintf("line_items[%d][price_data][product_data]", idx)
	for key := range p.values {
		if strings.HasPrefix(key, prefix+"[") {
			return true
		}
	}
	return false
}

// parseUnitAmountDecimal accepts integer minor-unit strings only (no fractional part).
// Stripe allows up to 12 decimal places; billtap simplifies: non-integral values are 400.
func parseUnitAmountDecimal(raw string) (int64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, fmt.Errorf("empty")
	}
	if strings.ContainsAny(raw, "eE") {
		return 0, fmt.Errorf("scientific notation")
	}
	if strings.Contains(raw, ".") {
		parts := strings.SplitN(raw, ".", 2)
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid decimal")
		}
		frac := strings.TrimRight(parts[1], "0")
		if frac != "" {
			return 0, fmt.Errorf("non-integer decimal")
		}
		raw = parts[0]
		if raw == "" || raw == "-" || raw == "+" {
			raw = "0"
		}
	}
	return strconv.ParseInt(raw, 10, 64)
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
			"coupon",
			"promotion_code",
			"automatic_tax[enabled]",
			"proration_behavior",
		},
		AllowedRegex: []*regexp.Regexp{subscriptionItemRE, discountParamRE, defaultTaxRatesParamRE, invoicePaymentSettingsRE},
		RequiredAny:  [][]string{{"customer", "customer_id"}},
		Int64Params:  []string{"days_until_due", "cancel_at", "billing_cycle_anchor"},
		BoolParams:   []string{"automatic_tax[enabled]"},
		Positive:     []string{"days_until_due"},
		EnumParams: map[string][]string{
			"collection_method":  {"charge_automatically", "send_invoice"},
			"proration_behavior": {"none", "create_prorations", "always_invoice"},
		},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	if err := validateDefaultTaxRatesVsAutomaticTax(p, "default_tax_rates", p.boolDefault("automatic_tax[enabled]", false)); err != nil {
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
		Allowed: []string{
			"cancel_at_period_end",
			"cancel_at",
			"pause_collection",
			"pause_collection[behavior]",
			"pause_collection[resumes_at]",
			"proration_behavior",
			"proration_date",
			"payment_behavior",
			"billing_cycle_anchor",
			"trial_end",
			"coupon",
			"promotion_code",
			// Emptyable clear form: default_tax_rates="" (single empty string).
			"default_tax_rates",
		},
		AllowedRegex: []*regexp.Regexp{subscriptionItemRE, cancellationDetailsRE, discountParamRE, defaultTaxRatesParamRE},
		BoolParams:   []string{"cancel_at_period_end"},
		Int64Params:  []string{"pause_collection[resumes_at]", "proration_date", "cancel_at"},
		EnumParams: map[string][]string{
			"pause_collection[behavior]": {"void", "keep_as_draft", "mark_uncollectible"},
			"proration_behavior":         {"none", "create_prorations", "always_invoice"},
			"payment_behavior": {
				"allow_incomplete",
				"default_incomplete",
				"error_if_incomplete",
				"pending_if_incomplete",
			},
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
	// billing_cycle_anchor accepts now | unchanged | unix timestamp (Stripe dual form).
	if err := p.validateUnixTimestampOrNowOrUnchanged("billing_cycle_anchor"); err != nil {
		return err
	}
	if err := p.validateUnixTimestampOrNow("trial_end"); err != nil {
		return err
	}
	return nil
}

func validateSubscriptionResume(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"billing_cycle_anchor",
			"proration_behavior",
			"proration_date",
		},
		Int64Params: []string{"proration_date"},
		EnumParams: map[string][]string{
			"billing_cycle_anchor": {"now", "unchanged"},
			"proration_behavior":   {"none", "create_prorations", "always_invoice"},
		},
	})
}

func validateSubscriptionItemCreate(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"subscription",
			"price",
			"price_id",
			"quantity",
			"proration_behavior",
			"proration_date",
		},
		AllowedRegex: []*regexp.Regexp{taxRatesParamRE},
		Required:     []string{"subscription"},
		RequiredAny:  [][]string{{"price", "price_id"}},
		Int64Params:  []string{"quantity", "proration_date"},
		Positive:     []string{"quantity"},
		EnumParams: map[string][]string{
			"proration_behavior": {"none", "create_prorations", "always_invoice"},
		},
		AllowMetadata: true,
	})
}

func validateSubscriptionItemDelete(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"proration_behavior",
			"proration_date",
			"clear_usage",
		},
		Int64Params: []string{"proration_date"},
		BoolParams:  []string{"clear_usage"},
		EnumParams: map[string][]string{
			"proration_behavior": {"none", "create_prorations", "always_invoice"},
		},
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
			"billtap_outcome",
			"deferred_outcome",
			"payment_intent_outcome",
			"description",
			"return_url",
			"receipt_email",
			"setup_future_usage",
			"off_session",
			"use_stripe_sdk",
			"billtap_next_action_type",
			"automatic_payment_methods[enabled]",
		},
		AllowedRegex: []*regexp.Regexp{paymentMethodTypesRE, paymentMethodOptionsRE},
		Required:     []string{"amount", "currency"},
		Int64Params:  []string{"amount"},
		Positive:     []string{"amount"},
		BoolParams:   []string{"confirm", "off_session", "use_stripe_sdk", "automatic_payment_methods[enabled]"},
		EnumParams: map[string][]string{
			"capture_method":           {"automatic", "automatic_async", "manual"},
			"setup_future_usage":       {"on_session", "off_session"},
			"billtap_next_action_type": {"use_stripe_sdk", "redirect_to_url"},
		},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	if p.boolDefault("confirm", false) && !p.hasAny("payment_method", "outcome", "customer") && !hasPaymentIntentDeferredOutcome(p) {
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
			"outcome",
			"paid_out_of_band",
			"payment_method",
			"source",
			"billtap_outcome",
		},
		BoolParams: []string{"forgive", "off_session", "paid_out_of_band"},
	})
}

func validateInvoiceCreate(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"id",
			"customer",
			"currency",
			"collection_method",
			"default_payment_method",
			"description",
			"auto_advance",
			"pending_invoice_items_behavior",
			"subscription",
			"days_until_due",
		},
		AllowedRegex: []*regexp.Regexp{invoicePaymentSettingsRE},
		Required:     []string{"customer"},
		BoolParams:   []string{"auto_advance"},
		Int64Params:  []string{"days_until_due"},
		Positive:     []string{"days_until_due"},
		EnumParams: map[string][]string{
			"collection_method":              {"charge_automatically", "send_invoice"},
			"pending_invoice_items_behavior": {"exclude", "include"},
		},
		AllowMetadata: true,
	})
}

func validateInvoiceItemCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed: []string{
			"id",
			"customer",
			"invoice",
			"subscription",
			"amount",
			"currency",
			"description",
			"pricing[price]",
			"quantity",
		},
		Required:      []string{"customer"},
		Int64Params:   []string{"amount", "quantity"},
		NonNegative:   []string{"quantity"},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	hasAmount := p.has("amount")
	hasPricing := p.has("pricing[price]")
	hasQuantity := p.has("quantity")
	if hasAmount && hasQuantity {
		return exclusiveParams("quantity", "amount", "quantity")
	}
	if hasAmount && hasPricing {
		return exclusiveParams("pricing", "amount", "pricing")
	}
	if !hasAmount && !hasPricing {
		return missingParam("amount")
	}
	if hasAmount {
		if !p.has("currency") {
			return missingParam("currency")
		}
		if p.int64("amount") == 0 {
			return invalidParam("amount", "Must be non-zero.")
		}
	}
	return nil
}

func validateInvoiceFinalize(p params) error {
	return p.validate(paramSpec{
		Allowed:    []string{"auto_advance"},
		BoolParams: []string{"auto_advance"},
	})
}

func validateInvoiceSend(p params) error {
	// Stripe InvoiceSendInvoiceParams only accepts expand (handled globally).
	return p.validate(paramSpec{})
}

func validateInvoiceVoid(p params) error {
	return p.validate(paramSpec{})
}

func validateInvoiceMarkUncollectible(p params) error {
	return p.validate(paramSpec{})
}

func validateCheckoutSessionExpire(p params) error {
	return p.validate(paramSpec{})
}

func validateInvoicePreview(p params) error {
	if err := p.validate(paramSpec{
		Allowed: []string{
			"customer",
			"subscription",
			"schedule",
			"currency",
			"price",
			"proration_behavior",
			"proration_date",
			"subscription_details[proration_behavior]",
			"subscription_details[proration_date]",
			"subscription_details[billing_cycle_anchor]",
			"subscription_details[trial_end]",
			"subscriptionDetails[prorationBehavior]",
			"subscriptionDetails[prorationDate]",
			"subscriptionDetails[billingCycleAnchor]",
			"subscriptionDetails[trialEnd]",
			"preview_mode",
			"coupon",
			"promotion_code",
		},
		AllowedRegex: []*regexp.Regexp{invoicePreviewItemParamRE, expandParamRE, discountParamRE},
		Int64Params:  []string{"proration_date", "subscription_details[proration_date]", "subscriptionDetails[prorationDate]"},
		EnumParams: map[string][]string{
			"proration_behavior":                       {"always_invoice", "create_prorations", "none"},
			"subscription_details[proration_behavior]": {"always_invoice", "create_prorations", "none"},
			"subscriptionDetails[prorationBehavior]":   {"always_invoice", "create_prorations", "none"},
			"preview_mode":                             {"next", "recurring"},
		},
	}); err != nil {
		return err
	}
	for _, key := range []string{
		"subscription_details[billing_cycle_anchor]",
		"subscriptionDetails[billingCycleAnchor]",
		"subscription_details[trial_end]",
		"subscriptionDetails[trialEnd]",
	} {
		if err := p.validateUnixTimestampOrNow(key); err != nil {
			return err
		}
	}
	return nil
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
			"status",
		},
		RequiredAny: [][]string{{"charge", "payment_intent", "invoice"}},
		Int64Params: []string{"amount"},
		Positive:    []string{"amount"},
		EnumParams: map[string][]string{
			"status": {"pending", "succeeded", "failed", "canceled"},
		},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	return nil
}

func validateRefundUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"metadata",
			"status",
		},
		EnumParams: map[string][]string{
			"status": {"pending", "succeeded", "failed", "canceled"},
		},
		AllowMetadata: true,
	})
}

func validateCreditNoteCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed: []string{
			"id",
			"invoice",
			"customer",
			"amount",
			"currency",
			"reason",
			"status",
			"memo",
			"out_of_band_amount",
			"refund_amount",
		},
		Required:      []string{"invoice"},
		Int64Params:   []string{"amount", "out_of_band_amount", "refund_amount"},
		Positive:      []string{"amount"},
		NonNegative:   []string{"out_of_band_amount", "refund_amount"},
		AllowMetadata: true,
		EnumParams: map[string][]string{
			"status": {"issued", "void"},
		},
	}); err != nil {
		return err
	}
	if !p.has("amount") && !p.has("out_of_band_amount") && !p.has("refund_amount") {
		return missingParam("amount")
	}
	return nil
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
			"billtap_next_action_type",
		},
		AllowedRegex: []*regexp.Regexp{paymentMethodOptionsRE},
		BoolParams:   []string{"off_session", "use_stripe_sdk"},
		EnumParams: map[string][]string{
			"setup_future_usage":       {"on_session", "off_session"},
			"billtap_next_action_type": {"use_stripe_sdk", "redirect_to_url"},
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

func validatePaymentIntentOutcomeUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed:  []string{"outcome"},
		Required: []string{"outcome"},
	})
}

func hasPaymentIntentDeferredOutcome(p params) bool {
	return p.hasAny(
		"billtap_outcome",
		"deferred_outcome",
		"payment_intent_outcome",
		"metadata["+billing.MetadataPaymentIntentOutcome+"]",
		"metadata[billtap_outcome]",
	)
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
	if err := p.validate(paramSpec{
		Allowed:      []string{"customer", "customer_id", "customer_account", "return_url", "configuration", "locale", "on_behalf_of"},
		AllowedRegex: []*regexp.Regexp{portalFlowDataParamRE},
		RequiredAny:  [][]string{{"customer", "customer_id"}},
		EnumParams: map[string][]string{
			"flow_data[type]":                   {"payment_method_update", "subscription_cancel", "subscription_update", "subscription_update_confirm"},
			"flow_data[after_completion][type]": {"hosted_confirmation", "portal_homepage", "redirect"},
			"locale":                            stripePortalLocales,
		},
	}); err != nil {
		return err
	}
	flowDataPresent := false
	for key := range p.values {
		if strings.HasPrefix(key, "flow_data[") {
			flowDataPresent = true
			break
		}
	}
	if flowDataPresent && !p.has("flow_data[type]") {
		return missingParam("flow_data[type]")
	}
	switch p.string("flow_data[type]") {
	case "subscription_cancel":
		if !p.has("flow_data[subscription_cancel][subscription]") {
			return missingParam("flow_data[subscription_cancel][subscription]")
		}
	case "subscription_update":
		if !p.has("flow_data[subscription_update][subscription]") {
			return missingParam("flow_data[subscription_update][subscription]")
		}
	case "subscription_update_confirm":
		if !p.has("flow_data[subscription_update_confirm][subscription]") {
			return missingParam("flow_data[subscription_update_confirm][subscription]")
		}
		if !p.has("flow_data[subscription_update_confirm][items]") && !p.has("flow_data[subscription_update_confirm][items][0][price]") {
			return missingParam("flow_data[subscription_update_confirm][items]")
		}
	}
	if p.string("flow_data[after_completion][type]") == "redirect" && !p.has("flow_data[after_completion][redirect][return_url]") {
		return missingParam("flow_data[after_completion][redirect][return_url]")
	}
	return nil
}

func validateBillingPortalConfigurationCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed: []string{
			"business_profile[headline]",
			"business_profile[privacy_policy_url]",
			"business_profile[terms_of_service_url]",
			"default_return_url",
			"login_page[logo_url]",
		},
		AllowedRegex: []*regexp.Regexp{portalFeatureParamRE},
		EnumParams: map[string][]string{
			"features[subscription_cancel][mode]":               {"at_period_end", "immediately"},
			"features[subscription_cancel][proration_behavior]": {"always_invoice", "create_prorations", "none"},
			"features[subscription_update][proration_behavior]": {"always_invoice", "create_prorations", "none"},
		},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	return validatePortalFeatureUpdateLists(p)
}

func validateBillingPortalConfigurationUpdate(p params) error {
	if err := p.validate(paramSpec{
		Allowed: []string{
			"active",
			"business_profile[headline]",
			"business_profile[privacy_policy_url]",
			"business_profile[terms_of_service_url]",
			"default_return_url",
			"login_page[logo_url]",
		},
		AllowedRegex: []*regexp.Regexp{portalFeatureParamRE},
		BoolParams:   []string{"active"},
		EnumParams: map[string][]string{
			"features[subscription_cancel][mode]":               {"at_period_end", "immediately"},
			"features[subscription_cancel][proration_behavior]": {"always_invoice", "create_prorations", "none"},
			"features[subscription_update][proration_behavior]": {"always_invoice", "create_prorations", "none"},
		},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	return validatePortalFeatureUpdateLists(p)
}

func validatePortalFeatureUpdateLists(p params) error {
	for _, allowed := range p.list("features[customer_update][allowed_updates]") {
		if allowed != "email" && allowed != "name" {
			return invalidParam("features[customer_update][allowed_updates]", "Invalid allowed_updates value. Allowed values are email, name.")
		}
	}
	for _, allowed := range p.list("features[subscription_update][default_allowed_updates]") {
		switch allowed {
		case "price", "quantity":
		default:
			return invalidParam("features[subscription_update][default_allowed_updates]", "Invalid default_allowed_updates value. Allowed values are price, quantity.")
		}
	}
	return nil
}

func validateSubscriptionItemList(p params) error {
	return p.validate(paramSpec{
		Allowed:     []string{"ending_before", "expand", "limit", "starting_after", "subscription"},
		Int64Params: []string{"limit"},
		Positive:    []string{"limit"},
	})
}

func validateSubscriptionItemUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"price",
			"price_id",
			"quantity",
			"proration_behavior",
			"proration_date",
		},
		AllowedRegex: []*regexp.Regexp{taxRatesParamRE},
		EnumParams: map[string][]string{
			"proration_behavior": {"always_invoice", "create_prorations", "none"},
		},
		Int64Params:   []string{"quantity"},
		Positive:      []string{"quantity"},
		AllowMetadata: true,
	})
}

func validatePaymentMethodList(p params) error {
	return p.validate(paramSpec{
		Allowed:     []string{"allow_redisplay", "customer", "customer_account", "ending_before", "limit", "starting_after", "type"},
		Int64Params: []string{"limit"},
		Positive:    []string{"limit"},
		EnumParams: map[string][]string{
			"allow_redisplay": {"always", "limited", "unspecified"},
			"type":            stripePaymentMethodTypes,
		},
	})
}

func validatePaymentMethodCreate(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"type",
			"customer",
			"billing_details[email]",
			"billing_details[name]",
			"card[token]",
		},
		EnumParams: map[string][]string{
			"type": {"card"},
		},
		AllowMetadata: true,
	})
}

func validatePaymentMethodAttach(p params) error {
	return p.validate(paramSpec{
		Allowed:  []string{"customer"},
		Required: []string{"customer"},
	})
}

func validatePaymentMethodUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"billing_details[email]",
			"billing_details[name]",
		},
		AllowMetadata: true,
	})
}

func validateTaxRateCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed: []string{
			"display_name",
			"percentage",
			"inclusive",
			"active",
			"country",
			"state",
			"jurisdiction",
			"description",
		},
		Required:      []string{"display_name", "percentage", "inclusive"},
		BoolParams:    []string{"inclusive", "active"},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	if _, err := strconv.ParseFloat(p.string("percentage"), 64); err != nil {
		return invalidParam("percentage", "Expected a number.")
	}
	return nil
}

func validateTaxRateUpdate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"active", "display_name", "description"},
		BoolParams:    []string{"active"},
		AllowMetadata: true,
	})
}

func validateCustomerTaxIDCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:  []string{"type", "value"},
		Required: []string{"type", "value"},
	})
}

func validateCouponCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed: []string{
			"id",
			"name",
			"duration",
			"percent_off",
			"amount_off",
			"currency",
			"duration_in_months",
			"redeem_by",
			"max_redemptions",
		},
		AllowedRegex: []*regexp.Regexp{couponAppliesToParamRE},
		// percent_off is a Stripe number (may be fractional, e.g. 12.5).
		Int64Params: []string{"amount_off", "duration_in_months", "redeem_by", "max_redemptions"},
		FloatParams: []string{"percent_off"},
		Positive:    []string{"amount_off", "duration_in_months", "max_redemptions"},
		EnumParams: map[string][]string{
			"duration": {"forever", "once", "repeating"},
		},
		AllowMetadata: true,
	}); err != nil {
		return err
	}
	if !p.hasAny("percent_off", "amount_off") {
		return missingParam("percent_off")
	}
	if p.has("percent_off") {
		if err := p.validatePositiveFloat("percent_off"); err != nil {
			return err
		}
		if p.float64("percent_off") > 100 {
			return invalidParam("percent_off", "Must be at most 100.")
		}
	}
	if p.has("amount_off") && !p.has("currency") {
		return missingParam("currency")
	}
	duration := p.stringDefault("duration", "once")
	if duration == "repeating" {
		if !p.has("duration_in_months") {
			return missingParam("duration_in_months")
		}
	} else if p.has("duration_in_months") {
		return invalidParam("duration_in_months", "Can only be set when duration is repeating.")
	}
	return nil
}

func validatePromotionCodeCreate(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"id",
			"coupon",
			"code",
			"active",
			"customer",
			"expires_at",
			"max_redemptions",
		},
		AllowedRegex:  []*regexp.Regexp{promotionRestrictionParamRE},
		Required:      []string{"coupon"},
		Int64Params:   []string{"expires_at", "max_redemptions"},
		Positive:      []string{"max_redemptions"},
		BoolParams:    []string{"active", "restrictions[first_time_transaction]"},
		AllowMetadata: true,
	})
}

func validateSubscriptionScheduleCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:      []string{"id", "from_subscription", "subscription", "end_behavior"},
		AllowedRegex: []*regexp.Regexp{schedulePhaseParamRE},
		RequiredAny:  [][]string{{"from_subscription", "subscription"}},
		Int64Params:  []string{"phases[0][start_date]", "phases[0][items][0][quantity]"},
		Positive:     []string{"phases[0][items][0][quantity]"},
		EnumParams: map[string][]string{
			"end_behavior": {"release", "cancel", "none", "renew"},
		},
		AllowMetadata: true,
	})
}

func validateFundCashBalance(p params) error {
	return p.validate(paramSpec{
		Allowed:     []string{"amount", "currency", "reference"},
		Required:    []string{"amount"},
		Int64Params: []string{"amount"},
		Positive:    []string{"amount"},
	})
}

func validateDisputeCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:     []string{"charge", "charge_id", "amount", "currency", "reason"},
		Int64Params: []string{"amount"},
		Positive:    []string{"amount"},
		EnumParams: map[string][]string{
			"reason": {"bank_cannot_process", "check_returned", "credit_not_processed", "customer_initiated", "debit_not_authorized", "duplicate", "fraudulent", "general", "incorrect_account_details", "insufficient_funds", "product_not_received", "product_unacceptable", "subscription_canceled", "unrecognized"},
		},
		AllowMetadata: true,
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

func validateHistoricalReplay(p params) error {
	return p.validate(paramSpec{
		Allowed: []string{
			"since",
			"until",
			"created_after",
			"createdAfter",
			"created_before",
			"createdBefore",
			"limit",
			"force",
		},
		AllowedRegex: []*regexp.Regexp{eventTypeFilterParamRE},
		Int64Params:  []string{"limit"},
		Positive:     []string{"limit"},
		BoolParams:   []string{"force"},
	})
}

func (p params) lineItemIndexes() map[int]struct{} {
	indexes := map[int]struct{}{}
	patterns := []*regexp.Regexp{
		checkoutLineItemRE,
		legacyLineItemParamRE,
		checkoutPriceDataRE,
		checkoutPriceDataProductDataRE,
		checkoutPriceDataProductMetaRE,
		checkoutPriceDataRecurringRE,
	}
	for key := range p.values {
		for _, pattern := range patterns {
			matches := pattern.FindStringSubmatch(key)
			if len(matches) < 2 {
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

// validateDefaultTaxRatesVsAutomaticTax rejects simultaneous automatic_tax and default_tax_rates.
// Presence of any non-empty default_tax_rates ID is enough; empty clear values are ignored.
func validateDefaultTaxRatesVsAutomaticTax(p params, prefix string, automaticTax bool) error {
	if !automaticTax {
		return nil
	}
	ids := p.defaultTaxRateIDs(prefix)
	if len(ids) == 0 {
		return nil
	}
	return &validationError{
		Param:   "default_tax_rates",
		Code:    stripeCodeParamInvalid,
		Message: "You cannot specify both automatic_tax[enabled]=true and default_tax_rates.",
	}
}

func missingTaxRate(id string) *validationError {
	return &validationError{
		Param:   "default_tax_rates",
		Code:    "resource_missing",
		Message: fmt.Sprintf("No such tax rate: '%s'", id),
	}
}
