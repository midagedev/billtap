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
	enabledEventsParamRE  = regexp.MustCompile(`^enabled_events(\[[^\]]*\])?$`)
	retryBackoffParamRE   = regexp.MustCompile(`^retry_backoff(\[[^\]]*\])?$`)
	checkoutLineItemRE    = regexp.MustCompile(`^line_items\[(\d+)\]\[(price|quantity)\]$`)
	legacyLineItemParamRE = regexp.MustCompile(`^lineItems\[(\d+)\]\[(price|quantity)\]$`)
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

func validateCustomerCreate(p params) error {
	return p.validate(paramSpec{
		Allowed:       []string{"id", "email", "name"},
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
		Required:      []string{"currency", "unit_amount"},
		RequiredAny:   [][]string{{"product", "product_id"}},
		Int64Params:   []string{"unit_amount", "recurring[interval_count]"},
		NonNegative:   []string{"unit_amount"},
		BoolParams:    []string{"active"},
		EnumParams:    map[string][]string{"recurring[interval]": {"day", "week", "month", "year"}},
		AllowMetadata: true,
	})
}

func validateCheckoutSessionCreate(p params) error {
	if err := p.validate(paramSpec{
		Allowed:      []string{"customer", "customer_id", "mode", "success_url", "cancel_url", "price"},
		AllowedRegex: []*regexp.Regexp{checkoutLineItemRE, legacyLineItemParamRE},
		RequiredAny:  [][]string{{"customer", "customer_id"}},
		EnumParams:   map[string][]string{"mode": {"subscription"}},
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
