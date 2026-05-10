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
			ID:              "customers.create.unknown_param",
			Name:            "Customer create rejects unknown parameter",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create customer with unknown nickname",
				Method: http.MethodPost,
				Path:   "/v1/customers",
				Params: map[string]string{"email": "scorecard@example.test", "nickname": "legacy"},
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
			ID:              "customers.create.expand_accepted",
			Name:            "Customer create accepts Stripe expand parameter as protocol input",
			Category:        "protocol",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/STRIPE_API_COMPATIBILITY_ROADMAP.md#s1-protocol-compatibility-baseline",
			Steps: []requestSpec{{
				Name:   "create customer with ignored expand parameter",
				Method: http.MethodPost,
				Path:   "/v1/customers",
				Params: map[string]string{"email": "expand@example.test", "expand[]": "subscriptions"},
			}},
			Expect: Observation{HTTPStatus: http.StatusOK, Object: "customer"},
		},
		{
			ID:              "customers.update.unknown_param",
			Name:            "Customer update rejects unknown parameter",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "update customer with unknown nickname",
				Method: http.MethodPost,
				Path:   "/v1/customers/cus_missing",
				Params: map[string]string{"nickname": "legacy"},
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
			ID:              "prices.create.invalid_json_amount_type",
			Name:            "Price create rejects JSON decimal amount",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create price with decimal JSON amount",
				Method: http.MethodPost,
				Path:   "/v1/prices",
				JSON: map[string]any{
					"product":     "prod_missing",
					"currency":    "usd",
					"unit_amount": 9.99,
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
			ID:              "prices.create.invalid_interval",
			Name:            "Price create rejects invalid recurring interval",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create price with invalid interval",
				Method: http.MethodPost,
				Path:   "/v1/prices",
				Params: map[string]string{
					"product":             "prod_missing",
					"currency":            "usd",
					"unit_amount":         "4900",
					"recurring[interval]": "decade",
				},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "recurring[interval]",
				},
			},
		},
		{
			ID:              "prices.update.invalid_active",
			Name:            "Price update rejects invalid active flag",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "update price with invalid active flag",
				Method: http.MethodPost,
				Path:   "/v1/prices/price_missing",
				Params: map[string]string{"active": "maybe"},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "active",
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
			ID:              "checkout.sessions.create.invalid_mode",
			Name:            "Checkout session create rejects unsupported mode",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create checkout session with payment mode",
				Method: http.MethodPost,
				Path:   "/v1/checkout/sessions",
				Params: map[string]string{
					"customer":             "cus_missing",
					"mode":                 "payment",
					"line_items[0][price]": "price_missing",
				},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "mode",
				},
			},
		},
		{
			ID:              "checkout.sessions.create.invalid_quantity",
			Name:            "Checkout session create rejects non-positive quantity",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create checkout session with quantity zero",
				Method: http.MethodPost,
				Path:   "/v1/checkout/sessions",
				Params: map[string]string{
					"customer":                "cus_missing",
					"line_items[0][price]":    "price_missing",
					"line_items[0][quantity]": "0",
				},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "line_items[0][quantity]",
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
			ID:              "checkout.sessions.create.java_sdk_optional_params",
			Name:            "Checkout session create accepts Stripe SDK promotion and trial params",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/COMPATIBILITY.md#supported-stripe-like-api-subset",
			Steps: []requestSpec{
				{
					Name:   "create checkout customer",
					Method: http.MethodPost,
					Path:   "/v1/customers",
					Params: map[string]string{
						"id":    "cus_scorecard_checkout_java",
						"email": "scorecard-checkout-java@example.test",
					},
				},
				{
					Name:   "create checkout product",
					Method: http.MethodPost,
					Path:   "/v1/products",
					Params: map[string]string{
						"id":   "prod_scorecard_checkout_java",
						"name": "Scorecard Checkout Java SDK",
					},
				},
				{
					Name:   "create checkout price",
					Method: http.MethodPost,
					Path:   "/v1/prices",
					Params: map[string]string{
						"id":                  "price_scorecard_checkout_java",
						"product":             "prod_scorecard_checkout_java",
						"currency":            "usd",
						"unit_amount":         "9900",
						"recurring[interval]": "month",
					},
				},
				{
					Name:   "create checkout session with Java SDK optional params",
					Method: http.MethodPost,
					Path:   "/v1/checkout/sessions",
					Params: map[string]string{
						"customer":                             "cus_scorecard_checkout_java",
						"mode":                                 "subscription",
						"line_items[0][price]":                 "price_scorecard_checkout_java",
						"line_items[0][quantity]":              "1",
						"allow_promotion_codes":                "true",
						"subscription_data[trial_period_days]": "14",
					},
				},
			},
			Expect: Observation{HTTPStatus: http.StatusOK, Object: "checkout.session", PaymentStatus: "unpaid"},
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
			ID:              "webhook_endpoints.create.invalid_retry_max_attempts",
			Name:            "Webhook endpoint create rejects invalid retry count",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create webhook endpoint with invalid retry count",
				Method: http.MethodPost,
				Path:   "/v1/webhook_endpoints",
				Params: map[string]string{
					"url":                "http://127.0.0.1/webhook",
					"retry_max_attempts": "abc",
				},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "retry_max_attempts",
				},
			},
		},
		{
			ID:              "webhook_endpoints.update.invalid_active",
			Name:            "Webhook endpoint update rejects invalid active flag",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "update webhook endpoint with invalid active flag",
				Method: http.MethodPost,
				Path:   "/v1/webhook_endpoints/we_missing",
				Params: map[string]string{"active": "maybe"},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "active",
				},
			},
		},
		{
			ID:              "subscriptions.create.invalid_quantity",
			Name:            "Subscription create rejects non-positive item quantity",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create subscription with quantity zero",
				Method: http.MethodPost,
				Path:   "/v1/subscriptions",
				Params: map[string]string{
					"customer":           "cus_missing",
					"items[0][price]":    "price_missing",
					"items[0][quantity]": "0",
				},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "items[0][quantity]",
				},
			},
		},
		{
			ID:              "subscriptions.update.invalid_cancel_flag",
			Name:            "Subscription update rejects invalid cancel flag",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "update subscription with invalid cancel flag",
				Method: http.MethodPost,
				Path:   "/v1/subscriptions/sub_missing",
				Params: map[string]string{"cancel_at_period_end": "maybe"},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "cancel_at_period_end",
				},
			},
		},
		{
			ID:              "subscription_items.create.invalid_quantity",
			Name:            "Subscription item create rejects non-positive quantity",
			Category:        "request-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create subscription item with quantity zero",
				Method: http.MethodPost,
				Path:   "/v1/subscription_items",
				Params: map[string]string{
					"subscription": "sub_missing",
					"price":        "price_missing",
					"quantity":     "0",
				},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "quantity",
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
			ID:              "openapi.country_specs.invalid_limit",
			Name:            "OpenAPI-backed fallback validation rejects invalid list limit",
			Category:        "openapi-schema-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/STRIPE_API_COMPATIBILITY_ROADMAP.md#s2-openapi-backed-validation-and-fixture-responses",
			Steps: []requestSpec{{
				Name:   "list country specs with invalid limit",
				Method: http.MethodGet,
				Path:   "/v1/country_specs?limit=not-an-int",
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "limit",
				},
			},
		},
		{
			ID:              "openapi.country_specs.unknown_param",
			Name:            "OpenAPI-backed fallback validation rejects unknown list parameter",
			Category:        "openapi-schema-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/STRIPE_API_COMPATIBILITY_ROADMAP.md#s2-openapi-backed-validation-and-fixture-responses",
			Steps: []requestSpec{{
				Name:   "list country specs with unknown parameter",
				Method: http.MethodGet,
				Path:   "/v1/country_specs?nickname=legacy",
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
			ID:              "openapi.apps.secrets.missing_required",
			Name:            "OpenAPI-backed fallback validation rejects missing required body parameter",
			Category:        "openapi-schema-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/STRIPE_API_COMPATIBILITY_ROADMAP.md#s2-openapi-backed-validation-and-fixture-responses",
			Steps: []requestSpec{{
				Name:   "create app secret without name",
				Method: http.MethodPost,
				Path:   "/v1/apps/secrets",
				Params: map[string]string{
					"payload":     "secret",
					"scope[type]": "account",
				},
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
			ID:              "openapi.apps.secrets.invalid_nested_enum",
			Name:            "OpenAPI-backed fallback validation rejects invalid nested enum",
			Category:        "openapi-schema-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/STRIPE_API_COMPATIBILITY_ROADMAP.md#s2-openapi-backed-validation-and-fixture-responses",
			Steps: []requestSpec{{
				Name:   "create app secret with invalid scope type",
				Method: http.MethodPost,
				Path:   "/v1/apps/secrets",
				Params: map[string]string{
					"name":        "token",
					"payload":     "secret",
					"scope[type]": "workspace",
				},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "scope[type]",
				},
			},
		},
		{
			ID:              "openapi.account_sessions.invalid_deep_nested_boolean",
			Name:            "OpenAPI-backed fallback validation rejects invalid deep nested boolean",
			Category:        "openapi-schema-validation",
			Level:           "L1",
			ReleaseBlocking: true,
			Reference:       "docs/STRIPE_API_COMPATIBILITY_ROADMAP.md#s2-openapi-backed-validation-and-fixture-responses",
			Steps: []requestSpec{{
				Name:   "create account session with invalid nested enabled flag",
				Method: http.MethodPost,
				Path:   "/v1/account_sessions",
				Params: map[string]string{
					"account": "acct_123",
					"components[account_onboarding][enabled]": "maybe",
				},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "components[account_onboarding][enabled]",
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
			Run: runCheckoutPaymentAlias("pm_card_visa_chargeDeclinedInsufficientFunds"),
		},
		{
			ID:              "checkout.complete.generic_decline_alias",
			Name:            "Checkout completion maps generic decline alias to card error",
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
					DeclineCode: "generic_decline",
				},
			},
			Run: runCheckoutPaymentAlias("pm_card_visa_chargeDeclined"),
		},
		{
			ID:              "checkout.complete.customer_payment_method_failed_alias",
			Name:            "Checkout completion maps customer payment method failure alias",
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
					DeclineCode: "do_not_honor",
				},
			},
			Run: runCheckoutPaymentAlias("pm_card_chargeCustomerFail"),
		},
		{
			ID:              "checkout.complete_authentication_required_alias",
			Name:            "Checkout completion maps 3DS alias to requires_action",
			Category:        "error-simulation",
			Level:           "L2",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l2---deterministic-error-simulation",
			Expect: Observation{
				HTTPStatus:          http.StatusOK,
				Object:              "payment_intent",
				PaymentIntentStatus: "requires_action",
				PaymentIntentError: &ErrorObservation{
					Type:        "card_error",
					Code:        "authentication_required",
					DeclineCode: "authentication_required",
				},
			},
			Run: runCheckoutPaymentAlias("pm_card_threeDSecure2Required"),
		},
		{
			ID:              "checkout.complete.expired_card_outcome",
			Name:            "Checkout completion maps expired card outcome",
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
					Code:        "expired_card",
					DeclineCode: "expired_card",
				},
			},
			Run: runCheckoutOutcome("expired_card"),
		},
		{
			ID:              "checkout.complete.incorrect_cvc_outcome",
			Name:            "Checkout completion maps incorrect CVC outcome",
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
					Code:        "incorrect_cvc",
					DeclineCode: "incorrect_cvc",
				},
			},
			Run: runCheckoutOutcome("incorrect_cvc"),
		},
		{
			ID:              "checkout.complete.processing_error_outcome",
			Name:            "Checkout completion maps processing error outcome",
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
					Code:        "processing_error",
					DeclineCode: "processing_error",
				},
			},
			Run: runCheckoutOutcome("processing_error"),
		},
		{
			ID:              "checkout.complete.missing_payment_method_outcome",
			Name:            "Checkout completion maps missing payment method outcome",
			Category:        "error-simulation",
			Level:           "L2",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l2---deterministic-error-simulation",
			Expect: Observation{
				HTTPStatus:          http.StatusOK,
				Object:              "payment_intent",
				PaymentIntentStatus: "requires_payment_method",
				PaymentIntentError: &ErrorObservation{
					Type: "card_error",
					Code: "payment_method_missing",
				},
			},
			Run: runCheckoutOutcome("missing_payment_method"),
		},
		{
			ID:              "payment_intents.create.confirm.succeeds",
			Name:            "Direct PaymentIntent create with confirm succeeds",
			Category:        "stateful-payment-intents",
			Level:           "L3",
			ReleaseBlocking: true,
			Reference:       "docs/COMPATIBILITY.md#supported-stripe-like-api-subset",
			Steps: []requestSpec{{
				Name:   "create and confirm payment intent",
				Method: http.MethodPost,
				Path:   "/v1/payment_intents",
				Params: map[string]string{
					"amount":         "4900",
					"currency":       "usd",
					"payment_method": "pm_card_visa",
					"confirm":        "true",
				},
			}},
			Expect: Observation{HTTPStatus: http.StatusOK, Object: "payment_intent", ObjectStatus: "succeeded", PaymentIntentStatus: "succeeded"},
		},
		{
			ID:              "payment_intents.confirm.card_decline",
			Name:            "Direct PaymentIntent confirm maps card decline alias",
			Category:        "stateful-payment-intents",
			Level:           "L3",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l2---deterministic-error-simulation",
			Steps: []requestSpec{
				{
					Name:             "create payment intent",
					Method:           http.MethodPost,
					Path:             "/v1/payment_intents",
					Params:           map[string]string{"id": "pi_scorecard_direct_decline", "amount": "4900", "currency": "usd"},
					ExpectHTTPStatus: http.StatusOK,
				},
				{
					Name:   "confirm payment intent with decline alias",
					Method: http.MethodPost,
					Path:   "/v1/payment_intents/pi_scorecard_direct_decline/confirm",
					Params: map[string]string{"payment_method": "pm_card_visa_chargeDeclined"},
				},
			},
			Expect: Observation{
				HTTPStatus:          http.StatusOK,
				Object:              "payment_intent",
				ObjectStatus:        "requires_payment_method",
				PaymentIntentStatus: "requires_payment_method",
				PaymentIntentError: &ErrorObservation{
					Type:        "card_error",
					Code:        "card_declined",
					DeclineCode: "generic_decline",
				},
			},
		},
		{
			ID:              "payment_intents.capture.manual.succeeds",
			Name:            "Direct manual-capture PaymentIntent captures successfully",
			Category:        "stateful-payment-intents",
			Level:           "L3",
			ReleaseBlocking: true,
			Reference:       "docs/COMPATIBILITY.md#supported-stripe-like-api-subset",
			Steps: []requestSpec{
				{
					Name:             "create confirmed manual payment intent",
					Method:           http.MethodPost,
					Path:             "/v1/payment_intents",
					Params:           map[string]string{"id": "pi_scorecard_manual_capture", "amount": "4900", "currency": "usd", "capture_method": "manual", "payment_method": "pm_card_visa", "confirm": "true"},
					ExpectHTTPStatus: http.StatusOK,
				},
				{
					Name:   "capture payment intent",
					Method: http.MethodPost,
					Path:   "/v1/payment_intents/pi_scorecard_manual_capture/capture",
					Params: map[string]string{"amount_to_capture": "4900"},
				},
			},
			Expect: Observation{HTTPStatus: http.StatusOK, Object: "payment_intent", ObjectStatus: "succeeded", PaymentIntentStatus: "succeeded"},
		},
		{
			ID:              "payment_intents.confirm.missing_payment_method",
			Name:            "Direct PaymentIntent confirm requires payment method or explicit outcome",
			Category:        "stateful-payment-intents",
			Level:           "L3",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l1---request-validation-parity-for-supported-endpoints",
			Steps: []requestSpec{{
				Name:   "create and confirm payment intent without payment method",
				Method: http.MethodPost,
				Path:   "/v1/payment_intents",
				Params: map[string]string{"amount": "4900", "currency": "usd", "confirm": "true"},
			}},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_missing",
					Param: "payment_method",
				},
			},
		},
		{
			ID:              "payment_intents.cancel.succeeded_rejected",
			Name:            "Direct PaymentIntent cancel rejects succeeded intent",
			Category:        "stateful-payment-intents",
			Level:           "L3",
			ReleaseBlocking: true,
			Reference:       "docs/COMPATIBILITY.md#supported-stripe-like-api-subset",
			Steps: []requestSpec{
				{
					Name:             "create succeeded payment intent",
					Method:           http.MethodPost,
					Path:             "/v1/payment_intents",
					Params:           map[string]string{"id": "pi_scorecard_cancel_terminal", "amount": "4900", "currency": "usd", "payment_method": "pm_card_visa", "confirm": "true"},
					ExpectHTTPStatus: http.StatusOK,
				},
				{
					Name:   "cancel succeeded payment intent",
					Method: http.MethodPost,
					Path:   "/v1/payment_intents/pi_scorecard_cancel_terminal/cancel",
				},
			},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "status",
				},
			},
		},
		{
			ID:              "setup_intents.create.confirm.succeeds",
			Name:            "Direct SetupIntent create with confirm succeeds",
			Category:        "stateful-setup-intents",
			Level:           "L3",
			ReleaseBlocking: true,
			Reference:       "docs/COMPATIBILITY.md#supported-stripe-like-api-subset",
			Steps: []requestSpec{{
				Name:   "create and confirm setup intent",
				Method: http.MethodPost,
				Path:   "/v1/setup_intents",
				Params: map[string]string{
					"payment_method": "pm_card_visa",
					"confirm":        "true",
					"usage":          "off_session",
				},
			}},
			Expect: Observation{HTTPStatus: http.StatusOK, Object: "setup_intent", ObjectStatus: "succeeded"},
		},
		{
			ID:              "setup_intents.confirm.card_decline",
			Name:            "Direct SetupIntent confirm maps card decline alias",
			Category:        "stateful-setup-intents",
			Level:           "L3",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l2---deterministic-error-simulation",
			Steps: []requestSpec{{
				Name:   "create and confirm declined setup intent",
				Method: http.MethodPost,
				Path:   "/v1/setup_intents",
				Params: map[string]string{"payment_method": "pm_card_visa_chargeDeclined", "confirm": "true"},
			}},
			Expect: Observation{HTTPStatus: http.StatusOK, Object: "setup_intent", ObjectStatus: "requires_payment_method"},
		},
		{
			ID:              "setup_intents.confirm.requires_action",
			Name:            "Direct SetupIntent confirm maps requires-action alias",
			Category:        "stateful-setup-intents",
			Level:           "L3",
			ReleaseBlocking: true,
			Reference:       "docs/API_VALIDATION_AND_ERROR_SIMULATION.md#l2---deterministic-error-simulation",
			Steps: []requestSpec{{
				Name:   "create and confirm requires-action setup intent",
				Method: http.MethodPost,
				Path:   "/v1/setup_intents",
				Params: map[string]string{"payment_method": "pm_card_threeDSecure2Required", "confirm": "true"},
			}},
			Expect: Observation{HTTPStatus: http.StatusOK, Object: "setup_intent", ObjectStatus: "requires_action"},
		},
		{
			ID:              "setup_intents.cancel.succeeded_rejected",
			Name:            "Direct SetupIntent cancel rejects succeeded intent",
			Category:        "stateful-setup-intents",
			Level:           "L3",
			ReleaseBlocking: true,
			Reference:       "docs/COMPATIBILITY.md#supported-stripe-like-api-subset",
			Steps: []requestSpec{
				{
					Name:             "create succeeded setup intent",
					Method:           http.MethodPost,
					Path:             "/v1/setup_intents",
					Params:           map[string]string{"id": "seti_scorecard_cancel_terminal", "payment_method": "pm_card_visa", "confirm": "true"},
					ExpectHTTPStatus: http.StatusOK,
				},
				{
					Name:   "cancel succeeded setup intent",
					Method: http.MethodPost,
					Path:   "/v1/setup_intents/seti_scorecard_cancel_terminal/cancel",
				},
			},
			Expect: Observation{
				HTTPStatus: http.StatusBadRequest,
				Error: &ErrorObservation{
					Type:  "invalid_request_error",
					Code:  "parameter_invalid",
					Param: "status",
				},
			},
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
			ID:          "payment_intents.apply_customer_balance.unsupported",
			Name:        "Payment intent customer-balance application is unsupported",
			Category:    "unsupported-provider-behavior",
			Level:       "L3",
			Reference:   "docs/COMPATIBILITY.md#unsupported-stripe-behavior",
			Unsupported: "payment intent customer balance application remains outside the local direct-intent state machine",
		},
	}
}

func runCheckoutPaymentAlias(paymentMethod string) func(context.Context, *harness) (caseExecution, error) {
	return func(ctx context.Context, h *harness) (caseExecution, error) {
		return runCheckoutCompletion(ctx, h, map[string]string{"payment_method": paymentMethod})
	}
}

func runCheckoutOutcome(outcome string) func(context.Context, *harness) (caseExecution, error) {
	return func(ctx context.Context, h *harness) (caseExecution, error) {
		return runCheckoutCompletion(ctx, h, map[string]string{"outcome": outcome})
	}
}

func runCheckoutCompletion(_ context.Context, h *harness, completionParams map[string]string) (caseExecution, error) {
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
		Params:           completionParams,
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
