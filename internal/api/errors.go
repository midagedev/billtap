package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/fixtures"
	"github.com/hckim/billtap/internal/webhooks"
)

const (
	stripeErrorAPI         = "api_error"
	stripeErrorCard        = "card_error"
	stripeErrorIdempotency = "idempotency_error"
	stripeErrorInvalidReq  = "invalid_request_error"
)

type stripeAPIError struct {
	Type        string `json:"type"`
	Message     string `json:"message"`
	Param       string `json:"param,omitempty"`
	Code        string `json:"code,omitempty"`
	DeclineCode string `json:"decline_code,omitempty"`
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeStripeError(w, status, stripeErrorFor(status, err))
}

func writeStripeError(w http.ResponseWriter, status int, apiErr stripeAPIError) {
	writeJSON(w, status, map[string]any{
		"error": apiErr,
	})
}

func stripeErrorFor(status int, err error) stripeAPIError {
	message := publicErrorMessage(err)
	apiErr := stripeAPIError{
		Type:    stripeErrorType(status),
		Message: message,
	}

	switch {
	case status == http.StatusNotFound || errors.Is(err, billing.ErrNotFound) || errors.Is(err, webhooks.ErrNotFound):
		apiErr.Code = "resource_missing"
	case status == http.StatusMethodNotAllowed:
		apiErr.Code = "method_not_allowed"
	case status == http.StatusConflict:
		apiErr.Type = stripeErrorIdempotency
		apiErr.Code = "idempotency_key_in_use"
	case strings.Contains(message, "real card data"):
		apiErr.Code = "real_card_data_not_allowed"
	default:
		apiErr.Param, apiErr.Code = inferValidationErrorFields(message)
	}

	return apiErr
}

func stripeErrorType(status int) string {
	switch {
	case status == http.StatusPaymentRequired:
		return stripeErrorCard
	case status == http.StatusConflict:
		return stripeErrorIdempotency
	case status >= 500:
		return stripeErrorAPI
	default:
		return stripeErrorInvalidReq
	}
}

func publicErrorMessage(err error) string {
	if err == nil {
		return "An unknown error occurred."
	}
	message := err.Error()
	for _, prefix := range []string{
		billing.ErrInvalidInput.Error() + ": ",
		billing.ErrUnsupportedOutcome.Error() + ": ",
		billing.ErrNotFound.Error() + ": ",
		fixtures.ErrInvalidFixture.Error() + ": ",
		webhooks.ErrInvalidInput.Error() + ": ",
		webhooks.ErrNotFound.Error() + ": ",
	} {
		message = strings.TrimPrefix(message, prefix)
	}
	return message
}

func inferValidationErrorFields(message string) (string, string) {
	switch {
	case message == "at least one line item is required":
		return "line_items", "parameter_missing"
	case strings.HasSuffix(message, " is required"):
		param := strings.TrimSuffix(message, " is required")
		param = strings.TrimPrefix(param, "product ")
		return stripeParamName(param), "parameter_missing"
	case strings.Contains(message, " must be "):
		param := strings.SplitN(message, " must be ", 2)[0]
		return stripeParamName(param), "parameter_invalid"
	default:
		return "", ""
	}
}

func stripeParamName(param string) string {
	param = strings.TrimSpace(param)
	param = strings.ReplaceAll(param, "].", "][")
	param = strings.ReplaceAll(param, ".", "][")
	if strings.Contains(param, "][") && !strings.HasSuffix(param, "]") {
		param += "]"
	}
	return param
}

func methodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", allow)
	writeStripeError(w, http.StatusMethodNotAllowed, stripeAPIError{
		Type:    stripeErrorInvalidReq,
		Message: "Method not allowed.",
		Code:    "method_not_allowed",
	})
}
