package compatibility

import (
	"context"
	"fmt"
	"net/http"
)

func builtinCorpus() []caseSpec {
	return []caseSpec{
		{
			ID:              "products.create.success",
			Name:            "Create a supported product",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/COMPATIBILITY.md#supported-stripe-like-api-subset",
			Steps: []requestSpec{{
				Name:   "create product",
				Method: http.MethodPost,
				Path:   "/v1/products",
				Params: map[string]string{"name": "Scorecard Product"},
			}},
			Expect: Observation{HTTPStatus: http.StatusOK, Object: "product", ObjectStatus: ""},
		},
		{
			ID:              "products.create.missing_name",
			Name:            "Product create missing required name",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create product without name",
				Method: http.MethodPost,
				Path:   "/v1/products",
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_missing",
					Param: "name",
				},
			},
		},
		{
			ID:              "products.create.unknown_param",
			Name:            "Product create rejects unknown parameter",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create product with unknown nickname",
				Method: http.MethodPost,
				Path:   "/v1/products",
				Params: map[string]string{"name": "Scorecard Product", "nickname": "legacy"},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_unknown",
					Param: "nickname",
				},
			},
		},
		{
			ID:              "prices.create.invalid_amount_type",
			Name:            "Price create rejects invalid amount type",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create price with non-integer amount",
				Method: http.MethodPost,
				Path:   "/v1/prices",
				Params: map[string]string{
					"product":     "prod_missing",
					"currency":    "usd",
					"unit_amount": "not-an-int",
				},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "unit_amount",
				},
			},
		},
		{
			ID:              "prices.create.missing_product",
			Name:            "Price create requires a product",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create price without product",
				Method: http.MethodPost,
				Path:   "/v1/prices",
				Params: map[string]string{
					"currency":    "usd",
					"unit_amount": "4900",
				},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_missing",
					Param: "product",
				},
			},
		},
		{
			ID:              "checkout.sessions.create.missing_line_items",
			Name:            "Checkout session create requires line items",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create checkout session without line items",
				Method: http.MethodPost,
				Path:   "/v1/checkout/sessions",
				Params: map[string]string{"customer": "cus_missing"},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_missing",
					Param: "line_items",
				},
			},
		},
		{
			ID:              "billing_portal.sessions.create.missing_customer",
			Name:            "Billing portal session requires a customer",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create billing portal session without customer",
				Method: http.MethodPost,
				Path:   "/v1/billing_portal/sessions",
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_missing",
					Param: "customer",
				},
			},
		},
		{
			ID:              "customers.create.idempotency_mismatch",
			Name:            "Idempotency key parameter mismatch returns conflict",
			Category:        "idempotency",
			Level:           "L2",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l2---deterministic-error-simulation",
			Steps: []requestSpec{
				{
					Name:             "create customer with idempotency key",
					Method:           http.MethodPost,
					Path:             "/v1/customers",
					Headers:          map[string]string{"Idempotency-Key": "scorecard-idem-customer"},
					Params:           map[string]string{"email": "first@example.test"},
					ExpectHTTPStatus: http.StatusOK,
				},
				{
					Name:    "reuse idempotency key with different params",
					Method:  http.MethodPost,
					Path:    "/v1/customers",
					Headers: map[string]string{"Idempotency-Key": "scorecard-idem-customer"},
					Params:  map[string]string{"email": "second@example.test"},
				},
			},
			Expect: Observation{
				HTTPStatus: http.StatusConflict,
				Error: &ErrorObservation{
					Type: "idempotency_error",
					Code: "idempotency_key_in_use",
				},
			},
		},
		{
			ID:              "checkout.complete.insufficient_funds_alias",
			Name:            "Checkout completion maps Stripe test PaymentMethod alias to card error",
			Category:        "error-simulation",
			Level:           "L2",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l2---deterministic-error-simulation",
			Expect: Observation{
				HTTPStatus:          http.StatusOK,
				Object:              "payment_intent",
				PaymentIntentStatus: "requires_payment_method",
				PaymentIntentError: &ErrorObservation{
					Type:        "card_error",
					Code:        "card_declined",
					DeclineCode: "insufficient_funds",
				},
			},
			Run: runInsufficientFundsAlias,
		},
		{
			ID:         "stripe_mock.oracle.route_parameter_lane",
			Name:       "stripe-mock route and parameter oracle lane",
			Category:   "differential-oracle",
			Level:      "L3",
			Reference:  "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l3---compatibility-scorecard-and-differential-harness",
			SkipReason: "external oracle intentionally skipped in the offline scorecard; Billtap-owned expectations still run locally",
		},
		{
			ID:          "payment_intents.confirm.unsupported",
			Name:        "Direct payment intent confirm is unsupported",
			Category:    "unsupported-provider-behavior",
			Level:       "L3",
			Reference:   "docs/COMPATIBILITY.md#unsupported-stripe-behavior",
			Unsupported: "direct payment intent create and confirm endpoints are documented as unsupported",
		},
	}
}

func runInsufficientFundsAlias(_ context.Context, h *harness) (caseExecution, error) {
	execution := caseExecution{}
	appendStep := func(spec requestSpec) (ReplayStep, error) {
		step, err := h.do(spec)
		execution.Steps = append(execution.Steps, step)
		if err != nil {
			return step, err
		}
		if spec.ExpectHTTPStatus != 0 && step.Actual.HTTPStatus != spec.ExpectHTTPStatus {
			return step, fmt.Errorf("%s returned HTTP %d, want %d", spec.Name, step.Actual.HTTPStatus, spec.ExpectHTTPStatus)
		}
		return step, nil
	}

	customer, err := appendStep(requestSpec{
		Name:             "create customer",
		Method:           http.MethodPost,
		Path:             "/v1/customers",
		Params:           map[string]string{"email": "scorecard@example.test"},
		ExpectHTTPStatus: http.StatusOK,
	})
	if err != nil {
		return execution, err
	}
	customerID := responseString(customer.Response.JSON, "id")
	if customerID == "" {
		return execution, fmt.Errorf("create customer did not return id")
	}

	product, err := appendStep(requestSpec{
		Name:             "create product",
		Method:           http.MethodPost,
		Path:             "/v1/products",
		Params:           map[string]string{"name": "Scorecard Pro"},
		ExpectHTTPStatus: http.StatusOK,
	})
	if err != nil {
		return execution, err
	}
	productID := responseString(product.Response.JSON, "id")
	if productID == "" {
		return execution, fmt.Errorf("create product did not return id")
	}

	price, err := appendStep(requestSpec{
		Name:             "create price",
		Method:           http.MethodPost,
		Path:             "/v1/prices",
		Params:           map[string]string{"product": productID, "currency": "usd", "unit_amount": "4900"},
		ExpectHTTPStatus: http.StatusOK,
	})
	if err != nil {
		return execution, err
	}
	priceID := responseString(price.Response.JSON, "id")
	if priceID == "" {
		return execution, fmt.Errorf("create price did not return id")
	}

	session, err := appendStep(requestSpec{
		Name:             "create checkout session",
		Method:           http.MethodPost,
		Path:             "/v1/checkout/sessions",
		Params:           map[string]string{"customer": customerID, "line_items[0][price]": priceID},
		ExpectHTTPStatus: http.StatusOK,
	})
	if err != nil {
		return execution, err
	}
	sessionID := responseString(session.Response.JSON, "id")
	if sessionID == "" {
		return execution, fmt.Errorf("create checkout session did not return id")
	}

	completion, err := appendStep(requestSpec{
		Name:             "complete checkout with insufficient funds alias",
		Method:           http.MethodPost,
		Path:             "/api/checkout/sessions/" + sessionID + "/complete",
		Params:           map[string]string{"payment_method": "pm_card_visa_chargeDeclinedInsufficientFunds"},
		ExpectHTTPStatus: http.StatusOK,
	})
	if err != nil {
		return execution, err
	}
	paymentIntentID := nestedResponseString(completion.Response.JSON, "payment_intent", "id")
	if paymentIntentID == "" {
		return execution, fmt.Errorf("checkout completion did not return payment_intent.id")
	}

	paymentIntent, err := appendStep(requestSpec{
		Name:   "retrieve payment intent projection",
		Method: http.MethodGet,
		Path:   "/v1/payment_intents/" + paymentIntentID,
	})
	if err != nil {
		return execution, err
	}
	execution.Actual = *paymentIntent.Actual
	return execution, nil
}

func responseString(decoded any, key string) string {
	m, _ := decoded.(map[string]any)
	value, _ := m[key].(string)
	return value
}

func nestedResponseString(decoded any, key string, nested string) string {
	m, _ := decoded.(map[string]any)
	child, _ := m[key].(map[string]any)
	value, _ := child[nested].(string)
	return value
}
