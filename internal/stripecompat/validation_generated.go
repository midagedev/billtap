// Code generated from stripe/openapi latest public spec for route and parameter validation; DO NOT EDIT.
package stripecompat

const defaultValidationSource = "stripe/openapi-latest-public"

func DefaultValidationOperations() []OperationValidation {
	operations := make([]OperationValidation, 0, 619)
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/account", OperationID: "GetAccount", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/account_links", OperationID: "PostAccountLinks", Params: []ParameterValidation{
		{Name: "account", Location: "form", Required: true, Type: "string"},
		{Name: "collect", Location: "form", Type: "string", Enum: []string{"currently_due", "eventually_due"}},
		{Name: "collection_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "fields", Location: "form", Type: "string", Enum: []string{"currently_due", "eventually_due"}},
			{Name: "future_requirements", Location: "form", Type: "string", Enum: []string{"include", "omit"}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "refresh_url", Location: "form", Type: "string"},
		{Name: "return_url", Location: "form", Type: "string"},
		{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account_onboarding", "account_update"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/account_sessions", OperationID: "PostAccountSessions", Params: []ParameterValidation{
		{Name: "account", Location: "form", Required: true, Type: "string"},
		{Name: "components", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "account_management", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "account_onboarding", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "balance_report", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "balances", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "disputes_list", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "financial_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "financial_account_transactions", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "instant_payouts_promotion", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "issuing_card", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "issuing_cards_list", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "notification_banner", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "payment_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "payment_disputes", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "payout_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "payout_reconciliation_report", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "payouts", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "payouts_list", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "tax_registrations", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "tax_settings", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/accounts", OperationID: "GetAccounts", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/accounts", OperationID: "PostAccounts", Params: []ParameterValidation{
		{Name: "account_token", Location: "form", Type: "string"},
		{Name: "bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_holder_name", Location: "form", Type: "string"},
			{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
			{Name: "account_number", Location: "form", Required: true, Type: "string"},
			{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "futsu", "savings", "toza"}},
			{Name: "country", Location: "form", Required: true, Type: "string"},
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank_account_ownership_verification", Location: "form", Type: "object"},
			}},
			{Name: "object", Location: "form", Type: "string", Enum: []string{"bank_account"}},
			{Name: "routing_number", Location: "form", Type: "string"},
		}},
		{Name: "business_profile", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "annual_revenue", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "amount", Location: "form", Required: true, Type: "integer"},
				{Name: "currency", Location: "form", Required: true, Type: "string"},
				{Name: "fiscal_year_end", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "estimated_worker_count", Location: "form", Type: "integer"},
			{Name: "mcc", Location: "form", Type: "string"},
			{Name: "minority_owned_business_designation", Location: "form", Type: "array"},
			{Name: "monthly_estimated_revenue", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "amount", Location: "form", Required: true, Type: "integer"},
				{Name: "currency", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "product_description", Location: "form", Type: "string"},
			{Name: "support_address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "support_email", Location: "form", Type: "string"},
			{Name: "support_phone", Location: "form", Type: "string"},
			{Name: "support_url", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "url", Location: "form", Type: "string"},
		}},
		{Name: "business_type", Location: "form", Type: "string", Enum: []string{"company", "government_entity", "individual", "non_profit"}},
		{Name: "capabilities", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "affirm_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "afterpay_clearpay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "alma_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "amazon_pay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "app_distribution", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "au_becs_debit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "bacs_debit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "bancontact_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "bank_transfer_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "billie_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "blik_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "boleto_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "card_issuing", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "card_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "cartes_bancaires_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "cashapp_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "crypto_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "eps_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "fpx_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "gb_bank_transfer_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "giropay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "grabpay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "ideal_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "india_international_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "jcb_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "jp_bank_transfer_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "kakao_pay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "klarna_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "konbini_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "kr_card_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "legacy_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "link_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "mb_way_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "mobilepay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "multibanco_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "mx_bank_transfer_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "naver_pay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "nz_bank_account_becs_debit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "oxxo_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "p24_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "pay_by_bank_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "payco_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "paynow_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "payto_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "pix_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "promptpay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "revolut_pay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "samsung_pay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "satispay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "sepa_bank_transfer_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "sepa_debit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "sofort_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "sunbit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "swish_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "tax_reporting_us_1099_k", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "tax_reporting_us_1099_misc", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "transfers", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "treasury", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "twint_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "upi_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "us_bank_account_ach_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "us_bank_transfer_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "zip_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
		}},
		{Name: "company", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "address_kana", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
			{Name: "address_kanji", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
			{Name: "directors_provided", Location: "form", Type: "boolean"},
			{Name: "directorship_declaration", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Type: "integer"},
				{Name: "ip", Location: "form", Type: "string"},
				{Name: "user_agent", Location: "form", Type: "string"},
			}},
			{Name: "executives_provided", Location: "form", Type: "boolean"},
			{Name: "export_license_id", Location: "form", Type: "string"},
			{Name: "export_purpose_code", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "name_kana", Location: "form", Type: "string"},
			{Name: "name_kanji", Location: "form", Type: "string"},
			{Name: "owners_provided", Location: "form", Type: "boolean"},
			{Name: "ownership_declaration", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Type: "integer"},
				{Name: "ip", Location: "form", Type: "string"},
				{Name: "user_agent", Location: "form", Type: "string"},
			}},
			{Name: "ownership_exemption_reason", Location: "form", Type: "string", Enum: []string{"", "qualified_entity_exceeds_ownership_threshold", "qualifies_as_financial_institution"}},
			{Name: "phone", Location: "form", Type: "string"},
			{Name: "registration_date", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "day", Location: "form", Required: true, Type: "integer"},
				{Name: "month", Location: "form", Required: true, Type: "integer"},
				{Name: "year", Location: "form", Required: true, Type: "integer"},
			}},
			{Name: "registration_number", Location: "form", Type: "string"},
			{Name: "representative_declaration", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Type: "integer"},
				{Name: "ip", Location: "form", Type: "string"},
				{Name: "user_agent", Location: "form", Type: "string"},
			}},
			{Name: "structure", Location: "form", Type: "string", Enum: []string{"", "free_zone_establishment", "free_zone_llc", "government_instrumentality", "governmental_unit", "incorporated_non_profit", "incorporated_partnership", "limited_liability_partnership", "llc", "multi_member_llc", "private_company", "private_corporation", "private_partnership", "public_company", "public_corporation", "public_partnership", "registered_charity", "single_member_llc", "sole_establishment", "sole_proprietorship", "tax_exempt_government_instrumentality", "unincorporated_association", "unincorporated_non_profit", "unincorporated_partnership"}},
			{Name: "tax_id", Location: "form", Type: "string"},
			{Name: "tax_id_registrar", Location: "form", Type: "string"},
			{Name: "vat_id", Location: "form", Type: "string"},
			{Name: "verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "document", Location: "form", Type: "object"},
			}},
		}},
		{Name: "controller", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "fees", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "payer", Location: "form", Type: "string", Enum: []string{"account", "application"}},
			}},
			{Name: "losses", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "payments", Location: "form", Type: "string", Enum: []string{"application", "stripe"}},
			}},
			{Name: "requirement_collection", Location: "form", Type: "string", Enum: []string{"application", "stripe"}},
			{Name: "stripe_dashboard", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Type: "string", Enum: []string{"express", "full", "none"}},
			}},
		}},
		{Name: "country", Location: "form", Type: "string"},
		{Name: "default_currency", Location: "form", Type: "string"},
		{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "bank_account_ownership_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "company_license", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "company_memorandum_of_association", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "company_ministerial_decree", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "company_registration_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "company_tax_id_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "proof_of_address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "proof_of_registration", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
				{Name: "signer", Location: "form", Type: "object"},
			}},
			{Name: "proof_of_ultimate_beneficial_ownership", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
				{Name: "signer", Location: "form", Type: "object"},
			}},
		}},
		{Name: "email", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "external_account", Location: "form", Type: "string"},
		{Name: "groups", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "payments_pricing", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "individual", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "address_kana", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
			{Name: "address_kanji", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
			{Name: "dob", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "day", Location: "form", Required: true, Type: "integer"},
				{Name: "month", Location: "form", Required: true, Type: "integer"},
				{Name: "year", Location: "form", Required: true, Type: "integer"},
			}},
			{Name: "email", Location: "form", Type: "string"},
			{Name: "first_name", Location: "form", Type: "string"},
			{Name: "first_name_kana", Location: "form", Type: "string"},
			{Name: "first_name_kanji", Location: "form", Type: "string"},
			{Name: "full_name_aliases", Location: "form", Enum: []string{""}},
			{Name: "gender", Location: "form", Type: "string"},
			{Name: "id_number", Location: "form", Type: "string"},
			{Name: "id_number_secondary", Location: "form", Type: "string"},
			{Name: "last_name", Location: "form", Type: "string"},
			{Name: "last_name_kana", Location: "form", Type: "string"},
			{Name: "last_name_kanji", Location: "form", Type: "string"},
			{Name: "maiden_name", Location: "form", Type: "string"},
			{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
			{Name: "phone", Location: "form", Type: "string"},
			{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
			{Name: "registered_address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "relationship", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "director", Location: "form", Type: "boolean"},
				{Name: "executive", Location: "form", Type: "boolean"},
				{Name: "owner", Location: "form", Type: "boolean"},
				{Name: "percent_ownership", Location: "form", Enum: []string{""}},
				{Name: "title", Location: "form", Type: "string"},
			}},
			{Name: "ssn_last_4", Location: "form", Type: "string"},
			{Name: "verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "additional_document", Location: "form", Type: "object"},
				{Name: "document", Location: "form", Type: "object"},
			}},
		}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "bacs_debit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "display_name", Location: "form", Type: "string"},
			}},
			{Name: "branding", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "icon", Location: "form", Type: "string"},
				{Name: "logo", Location: "form", Type: "string"},
				{Name: "primary_color", Location: "form", Type: "string"},
				{Name: "secondary_color", Location: "form", Type: "string"},
			}},
			{Name: "card_issuing", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tos_acceptance", Location: "form", Type: "object"},
			}},
			{Name: "card_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "decline_on", Location: "form", Type: "object"},
				{Name: "statement_descriptor_prefix", Location: "form", Type: "string"},
				{Name: "statement_descriptor_prefix_kana", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "statement_descriptor_prefix_kanji", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "invoices", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "hosted_payment_method_save", Location: "form", Type: "string", Enum: []string{"always", "never", "offer"}},
			}},
			{Name: "payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "statement_descriptor", Location: "form", Type: "string"},
				{Name: "statement_descriptor_kana", Location: "form", Type: "string"},
				{Name: "statement_descriptor_kanji", Location: "form", Type: "string"},
			}},
			{Name: "payouts", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "debit_negative_balances", Location: "form", Type: "boolean"},
				{Name: "schedule", Location: "form", Type: "object"},
				{Name: "statement_descriptor", Location: "form", Type: "string"},
			}},
			{Name: "treasury", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tos_acceptance", Location: "form", Type: "object"},
			}},
		}},
		{Name: "tos_acceptance", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "date", Location: "form", Type: "integer"},
			{Name: "ip", Location: "form", Type: "string"},
			{Name: "service_agreement", Location: "form", Type: "string"},
			{Name: "user_agent", Location: "form", Type: "string"},
		}},
		{Name: "type", Location: "form", Type: "string", Enum: []string{"custom", "express", "standard"}},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/accounts/{account}", OperationID: "DeleteAccountsAccount", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/accounts/{account}", OperationID: "GetAccountsAccount", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/accounts/{account}", OperationID: "PostAccountsAccount", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "account_token", Location: "form", Type: "string"},
		{Name: "business_profile", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "annual_revenue", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "amount", Location: "form", Required: true, Type: "integer"},
				{Name: "currency", Location: "form", Required: true, Type: "string"},
				{Name: "fiscal_year_end", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "estimated_worker_count", Location: "form", Type: "integer"},
			{Name: "mcc", Location: "form", Type: "string"},
			{Name: "minority_owned_business_designation", Location: "form", Type: "array"},
			{Name: "monthly_estimated_revenue", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "amount", Location: "form", Required: true, Type: "integer"},
				{Name: "currency", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "product_description", Location: "form", Type: "string"},
			{Name: "support_address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "support_email", Location: "form", Type: "string"},
			{Name: "support_phone", Location: "form", Type: "string"},
			{Name: "support_url", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "url", Location: "form", Type: "string"},
		}},
		{Name: "business_type", Location: "form", Type: "string", Enum: []string{"company", "government_entity", "individual", "non_profit"}},
		{Name: "capabilities", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "affirm_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "afterpay_clearpay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "alma_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "amazon_pay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "app_distribution", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "au_becs_debit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "bacs_debit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "bancontact_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "bank_transfer_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "billie_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "blik_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "boleto_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "card_issuing", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "card_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "cartes_bancaires_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "cashapp_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "crypto_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "eps_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "fpx_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "gb_bank_transfer_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "giropay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "grabpay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "ideal_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "india_international_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "jcb_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "jp_bank_transfer_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "kakao_pay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "klarna_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "konbini_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "kr_card_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "legacy_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "link_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "mb_way_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "mobilepay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "multibanco_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "mx_bank_transfer_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "naver_pay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "nz_bank_account_becs_debit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "oxxo_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "p24_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "pay_by_bank_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "payco_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "paynow_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "payto_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "pix_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "promptpay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "revolut_pay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "samsung_pay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "satispay_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "sepa_bank_transfer_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "sepa_debit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "sofort_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "sunbit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "swish_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "tax_reporting_us_1099_k", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "tax_reporting_us_1099_misc", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "transfers", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "treasury", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "twint_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "upi_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "us_bank_account_ach_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "us_bank_transfer_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
			{Name: "zip_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Type: "boolean"},
			}},
		}},
		{Name: "company", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "address_kana", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
			{Name: "address_kanji", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
			{Name: "directors_provided", Location: "form", Type: "boolean"},
			{Name: "directorship_declaration", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Type: "integer"},
				{Name: "ip", Location: "form", Type: "string"},
				{Name: "user_agent", Location: "form", Type: "string"},
			}},
			{Name: "executives_provided", Location: "form", Type: "boolean"},
			{Name: "export_license_id", Location: "form", Type: "string"},
			{Name: "export_purpose_code", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "name_kana", Location: "form", Type: "string"},
			{Name: "name_kanji", Location: "form", Type: "string"},
			{Name: "owners_provided", Location: "form", Type: "boolean"},
			{Name: "ownership_declaration", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Type: "integer"},
				{Name: "ip", Location: "form", Type: "string"},
				{Name: "user_agent", Location: "form", Type: "string"},
			}},
			{Name: "ownership_exemption_reason", Location: "form", Type: "string", Enum: []string{"", "qualified_entity_exceeds_ownership_threshold", "qualifies_as_financial_institution"}},
			{Name: "phone", Location: "form", Type: "string"},
			{Name: "registration_date", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "day", Location: "form", Required: true, Type: "integer"},
				{Name: "month", Location: "form", Required: true, Type: "integer"},
				{Name: "year", Location: "form", Required: true, Type: "integer"},
			}},
			{Name: "registration_number", Location: "form", Type: "string"},
			{Name: "representative_declaration", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Type: "integer"},
				{Name: "ip", Location: "form", Type: "string"},
				{Name: "user_agent", Location: "form", Type: "string"},
			}},
			{Name: "structure", Location: "form", Type: "string", Enum: []string{"", "free_zone_establishment", "free_zone_llc", "government_instrumentality", "governmental_unit", "incorporated_non_profit", "incorporated_partnership", "limited_liability_partnership", "llc", "multi_member_llc", "private_company", "private_corporation", "private_partnership", "public_company", "public_corporation", "public_partnership", "registered_charity", "single_member_llc", "sole_establishment", "sole_proprietorship", "tax_exempt_government_instrumentality", "unincorporated_association", "unincorporated_non_profit", "unincorporated_partnership"}},
			{Name: "tax_id", Location: "form", Type: "string"},
			{Name: "tax_id_registrar", Location: "form", Type: "string"},
			{Name: "vat_id", Location: "form", Type: "string"},
			{Name: "verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "document", Location: "form", Type: "object"},
			}},
		}},
		{Name: "default_currency", Location: "form", Type: "string"},
		{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "bank_account_ownership_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "company_license", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "company_memorandum_of_association", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "company_ministerial_decree", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "company_registration_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "company_tax_id_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "proof_of_address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "proof_of_registration", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
				{Name: "signer", Location: "form", Type: "object"},
			}},
			{Name: "proof_of_ultimate_beneficial_ownership", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
				{Name: "signer", Location: "form", Type: "object"},
			}},
		}},
		{Name: "email", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "external_account", Location: "form", Type: "string"},
		{Name: "groups", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "payments_pricing", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "individual", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "address_kana", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
			{Name: "address_kanji", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
			{Name: "dob", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "day", Location: "form", Required: true, Type: "integer"},
				{Name: "month", Location: "form", Required: true, Type: "integer"},
				{Name: "year", Location: "form", Required: true, Type: "integer"},
			}},
			{Name: "email", Location: "form", Type: "string"},
			{Name: "first_name", Location: "form", Type: "string"},
			{Name: "first_name_kana", Location: "form", Type: "string"},
			{Name: "first_name_kanji", Location: "form", Type: "string"},
			{Name: "full_name_aliases", Location: "form", Enum: []string{""}},
			{Name: "gender", Location: "form", Type: "string"},
			{Name: "id_number", Location: "form", Type: "string"},
			{Name: "id_number_secondary", Location: "form", Type: "string"},
			{Name: "last_name", Location: "form", Type: "string"},
			{Name: "last_name_kana", Location: "form", Type: "string"},
			{Name: "last_name_kanji", Location: "form", Type: "string"},
			{Name: "maiden_name", Location: "form", Type: "string"},
			{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
			{Name: "phone", Location: "form", Type: "string"},
			{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
			{Name: "registered_address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "relationship", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "director", Location: "form", Type: "boolean"},
				{Name: "executive", Location: "form", Type: "boolean"},
				{Name: "owner", Location: "form", Type: "boolean"},
				{Name: "percent_ownership", Location: "form", Enum: []string{""}},
				{Name: "title", Location: "form", Type: "string"},
			}},
			{Name: "ssn_last_4", Location: "form", Type: "string"},
			{Name: "verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "additional_document", Location: "form", Type: "object"},
				{Name: "document", Location: "form", Type: "object"},
			}},
		}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "bacs_debit_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "display_name", Location: "form", Type: "string"},
			}},
			{Name: "branding", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "icon", Location: "form", Type: "string"},
				{Name: "logo", Location: "form", Type: "string"},
				{Name: "primary_color", Location: "form", Type: "string"},
				{Name: "secondary_color", Location: "form", Type: "string"},
			}},
			{Name: "card_issuing", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tos_acceptance", Location: "form", Type: "object"},
			}},
			{Name: "card_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "decline_on", Location: "form", Type: "object"},
				{Name: "statement_descriptor_prefix", Location: "form", Type: "string"},
				{Name: "statement_descriptor_prefix_kana", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "statement_descriptor_prefix_kanji", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "invoices", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "default_account_tax_ids", Location: "form", Enum: []string{""}},
				{Name: "hosted_payment_method_save", Location: "form", Type: "string", Enum: []string{"always", "never", "offer"}},
			}},
			{Name: "payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "statement_descriptor", Location: "form", Type: "string"},
				{Name: "statement_descriptor_kana", Location: "form", Type: "string"},
				{Name: "statement_descriptor_kanji", Location: "form", Type: "string"},
			}},
			{Name: "payouts", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "debit_negative_balances", Location: "form", Type: "boolean"},
				{Name: "schedule", Location: "form", Type: "object"},
				{Name: "statement_descriptor", Location: "form", Type: "string"},
			}},
			{Name: "treasury", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tos_acceptance", Location: "form", Type: "object"},
			}},
		}},
		{Name: "tos_acceptance", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "date", Location: "form", Type: "integer"},
			{Name: "ip", Location: "form", Type: "string"},
			{Name: "service_agreement", Location: "form", Type: "string"},
			{Name: "user_agent", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/accounts/{account}/bank_accounts", OperationID: "PostAccountsAccountBankAccounts", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_holder_name", Location: "form", Type: "string"},
			{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
			{Name: "account_number", Location: "form", Required: true, Type: "string"},
			{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "futsu", "savings", "toza"}},
			{Name: "country", Location: "form", Required: true, Type: "string"},
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank_account_ownership_verification", Location: "form", Type: "object"},
			}},
			{Name: "object", Location: "form", Type: "string", Enum: []string{"bank_account"}},
			{Name: "routing_number", Location: "form", Type: "string"},
		}},
		{Name: "default_for_currency", Location: "form", Type: "boolean"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "external_account", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/accounts/{account}/bank_accounts/{id}", OperationID: "DeleteAccountsAccountBankAccountsId", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/accounts/{account}/bank_accounts/{id}", OperationID: "GetAccountsAccountBankAccountsId", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/accounts/{account}/bank_accounts/{id}", OperationID: "PostAccountsAccountBankAccountsId", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "account_holder_name", Location: "form", Type: "string"},
		{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"", "company", "individual"}},
		{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "futsu", "savings", "toza"}},
		{Name: "address_city", Location: "form", Type: "string"},
		{Name: "address_country", Location: "form", Type: "string"},
		{Name: "address_line1", Location: "form", Type: "string"},
		{Name: "address_line2", Location: "form", Type: "string"},
		{Name: "address_state", Location: "form", Type: "string"},
		{Name: "address_zip", Location: "form", Type: "string"},
		{Name: "default_for_currency", Location: "form", Type: "boolean"},
		{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "bank_account_ownership_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
		}},
		{Name: "exp_month", Location: "form", Type: "string"},
		{Name: "exp_year", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "name", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/accounts/{account}/capabilities", OperationID: "GetAccountsAccountCapabilities", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/accounts/{account}/capabilities/{capability}", OperationID: "GetAccountsAccountCapabilitiesCapability", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "capability", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/accounts/{account}/capabilities/{capability}", OperationID: "PostAccountsAccountCapabilitiesCapability", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "capability", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "requested", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/accounts/{account}/external_accounts", OperationID: "GetAccountsAccountExternalAccounts", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "object", Location: "query", Type: "string", Enum: []string{"bank_account", "card"}},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/accounts/{account}/external_accounts", OperationID: "PostAccountsAccountExternalAccounts", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_holder_name", Location: "form", Type: "string"},
			{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
			{Name: "account_number", Location: "form", Required: true, Type: "string"},
			{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "futsu", "savings", "toza"}},
			{Name: "country", Location: "form", Required: true, Type: "string"},
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank_account_ownership_verification", Location: "form", Type: "object"},
			}},
			{Name: "object", Location: "form", Type: "string", Enum: []string{"bank_account"}},
			{Name: "routing_number", Location: "form", Type: "string"},
		}},
		{Name: "default_for_currency", Location: "form", Type: "boolean"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "external_account", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/accounts/{account}/external_accounts/{id}", OperationID: "DeleteAccountsAccountExternalAccountsId", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/accounts/{account}/external_accounts/{id}", OperationID: "GetAccountsAccountExternalAccountsId", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/accounts/{account}/external_accounts/{id}", OperationID: "PostAccountsAccountExternalAccountsId", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "account_holder_name", Location: "form", Type: "string"},
		{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"", "company", "individual"}},
		{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "futsu", "savings", "toza"}},
		{Name: "address_city", Location: "form", Type: "string"},
		{Name: "address_country", Location: "form", Type: "string"},
		{Name: "address_line1", Location: "form", Type: "string"},
		{Name: "address_line2", Location: "form", Type: "string"},
		{Name: "address_state", Location: "form", Type: "string"},
		{Name: "address_zip", Location: "form", Type: "string"},
		{Name: "default_for_currency", Location: "form", Type: "boolean"},
		{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "bank_account_ownership_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
		}},
		{Name: "exp_month", Location: "form", Type: "string"},
		{Name: "exp_year", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "name", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/accounts/{account}/login_links", OperationID: "PostAccountsAccountLoginLinks", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/accounts/{account}/people", OperationID: "GetAccountsAccountPeople", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "relationship", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "authorizer", Location: "query", Type: "boolean"},
			{Name: "director", Location: "query", Type: "boolean"},
			{Name: "executive", Location: "query", Type: "boolean"},
			{Name: "legal_guardian", Location: "query", Type: "boolean"},
			{Name: "owner", Location: "query", Type: "boolean"},
			{Name: "representative", Location: "query", Type: "boolean"},
		}},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/accounts/{account}/people", OperationID: "PostAccountsAccountPeople", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "additional_tos_acceptances", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Type: "integer"},
				{Name: "ip", Location: "form", Type: "string"},
				{Name: "user_agent", Location: "form", Type: "string", Enum: []string{""}},
			}},
		}},
		{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
		}},
		{Name: "address_kana", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "address_kanji", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "dob", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "day", Location: "form", Required: true, Type: "integer"},
			{Name: "month", Location: "form", Required: true, Type: "integer"},
			{Name: "year", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "company_authorization", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "passport", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "visa", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
		}},
		{Name: "email", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "first_name", Location: "form", Type: "string"},
		{Name: "first_name_kana", Location: "form", Type: "string"},
		{Name: "first_name_kanji", Location: "form", Type: "string"},
		{Name: "full_name_aliases", Location: "form", Enum: []string{""}},
		{Name: "gender", Location: "form", Type: "string"},
		{Name: "id_number", Location: "form", Type: "string"},
		{Name: "id_number_secondary", Location: "form", Type: "string"},
		{Name: "last_name", Location: "form", Type: "string"},
		{Name: "last_name_kana", Location: "form", Type: "string"},
		{Name: "last_name_kanji", Location: "form", Type: "string"},
		{Name: "maiden_name", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "nationality", Location: "form", Type: "string"},
		{Name: "person_token", Location: "form", Type: "string"},
		{Name: "phone", Location: "form", Type: "string"},
		{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
		{Name: "registered_address", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
		}},
		{Name: "relationship", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "authorizer", Location: "form", Type: "boolean"},
			{Name: "director", Location: "form", Type: "boolean"},
			{Name: "executive", Location: "form", Type: "boolean"},
			{Name: "legal_guardian", Location: "form", Type: "boolean"},
			{Name: "owner", Location: "form", Type: "boolean"},
			{Name: "percent_ownership", Location: "form", Enum: []string{""}},
			{Name: "representative", Location: "form", Type: "boolean"},
			{Name: "title", Location: "form", Type: "string"},
		}},
		{Name: "ssn_last_4", Location: "form", Type: "string"},
		{Name: "us_cfpb_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "ethnicity_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ethnicity", Location: "form", Type: "array"},
				{Name: "ethnicity_other", Location: "form", Type: "string"},
			}},
			{Name: "race_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "race", Location: "form", Type: "array"},
				{Name: "race_other", Location: "form", Type: "string"},
			}},
			{Name: "self_identified_gender", Location: "form", Type: "string"},
		}},
		{Name: "verification", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "additional_document", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "back", Location: "form", Type: "string"},
				{Name: "front", Location: "form", Type: "string"},
			}},
			{Name: "document", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "back", Location: "form", Type: "string"},
				{Name: "front", Location: "form", Type: "string"},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/accounts/{account}/people/{person}", OperationID: "DeleteAccountsAccountPeoplePerson", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "person", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/accounts/{account}/people/{person}", OperationID: "GetAccountsAccountPeoplePerson", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "person", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/accounts/{account}/people/{person}", OperationID: "PostAccountsAccountPeoplePerson", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "person", Location: "path", Required: true, Type: "string"},
		{Name: "additional_tos_acceptances", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Type: "integer"},
				{Name: "ip", Location: "form", Type: "string"},
				{Name: "user_agent", Location: "form", Type: "string", Enum: []string{""}},
			}},
		}},
		{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
		}},
		{Name: "address_kana", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "address_kanji", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "dob", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "day", Location: "form", Required: true, Type: "integer"},
			{Name: "month", Location: "form", Required: true, Type: "integer"},
			{Name: "year", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "company_authorization", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "passport", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "visa", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
		}},
		{Name: "email", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "first_name", Location: "form", Type: "string"},
		{Name: "first_name_kana", Location: "form", Type: "string"},
		{Name: "first_name_kanji", Location: "form", Type: "string"},
		{Name: "full_name_aliases", Location: "form", Enum: []string{""}},
		{Name: "gender", Location: "form", Type: "string"},
		{Name: "id_number", Location: "form", Type: "string"},
		{Name: "id_number_secondary", Location: "form", Type: "string"},
		{Name: "last_name", Location: "form", Type: "string"},
		{Name: "last_name_kana", Location: "form", Type: "string"},
		{Name: "last_name_kanji", Location: "form", Type: "string"},
		{Name: "maiden_name", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "nationality", Location: "form", Type: "string"},
		{Name: "person_token", Location: "form", Type: "string"},
		{Name: "phone", Location: "form", Type: "string"},
		{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
		{Name: "registered_address", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
		}},
		{Name: "relationship", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "authorizer", Location: "form", Type: "boolean"},
			{Name: "director", Location: "form", Type: "boolean"},
			{Name: "executive", Location: "form", Type: "boolean"},
			{Name: "legal_guardian", Location: "form", Type: "boolean"},
			{Name: "owner", Location: "form", Type: "boolean"},
			{Name: "percent_ownership", Location: "form", Enum: []string{""}},
			{Name: "representative", Location: "form", Type: "boolean"},
			{Name: "title", Location: "form", Type: "string"},
		}},
		{Name: "ssn_last_4", Location: "form", Type: "string"},
		{Name: "us_cfpb_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "ethnicity_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ethnicity", Location: "form", Type: "array"},
				{Name: "ethnicity_other", Location: "form", Type: "string"},
			}},
			{Name: "race_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "race", Location: "form", Type: "array"},
				{Name: "race_other", Location: "form", Type: "string"},
			}},
			{Name: "self_identified_gender", Location: "form", Type: "string"},
		}},
		{Name: "verification", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "additional_document", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "back", Location: "form", Type: "string"},
				{Name: "front", Location: "form", Type: "string"},
			}},
			{Name: "document", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "back", Location: "form", Type: "string"},
				{Name: "front", Location: "form", Type: "string"},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/accounts/{account}/persons", OperationID: "GetAccountsAccountPersons", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "relationship", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "authorizer", Location: "query", Type: "boolean"},
			{Name: "director", Location: "query", Type: "boolean"},
			{Name: "executive", Location: "query", Type: "boolean"},
			{Name: "legal_guardian", Location: "query", Type: "boolean"},
			{Name: "owner", Location: "query", Type: "boolean"},
			{Name: "representative", Location: "query", Type: "boolean"},
		}},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/accounts/{account}/persons", OperationID: "PostAccountsAccountPersons", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "additional_tos_acceptances", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Type: "integer"},
				{Name: "ip", Location: "form", Type: "string"},
				{Name: "user_agent", Location: "form", Type: "string", Enum: []string{""}},
			}},
		}},
		{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
		}},
		{Name: "address_kana", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "address_kanji", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "dob", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "day", Location: "form", Required: true, Type: "integer"},
			{Name: "month", Location: "form", Required: true, Type: "integer"},
			{Name: "year", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "company_authorization", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "passport", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "visa", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
		}},
		{Name: "email", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "first_name", Location: "form", Type: "string"},
		{Name: "first_name_kana", Location: "form", Type: "string"},
		{Name: "first_name_kanji", Location: "form", Type: "string"},
		{Name: "full_name_aliases", Location: "form", Enum: []string{""}},
		{Name: "gender", Location: "form", Type: "string"},
		{Name: "id_number", Location: "form", Type: "string"},
		{Name: "id_number_secondary", Location: "form", Type: "string"},
		{Name: "last_name", Location: "form", Type: "string"},
		{Name: "last_name_kana", Location: "form", Type: "string"},
		{Name: "last_name_kanji", Location: "form", Type: "string"},
		{Name: "maiden_name", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "nationality", Location: "form", Type: "string"},
		{Name: "person_token", Location: "form", Type: "string"},
		{Name: "phone", Location: "form", Type: "string"},
		{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
		{Name: "registered_address", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
		}},
		{Name: "relationship", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "authorizer", Location: "form", Type: "boolean"},
			{Name: "director", Location: "form", Type: "boolean"},
			{Name: "executive", Location: "form", Type: "boolean"},
			{Name: "legal_guardian", Location: "form", Type: "boolean"},
			{Name: "owner", Location: "form", Type: "boolean"},
			{Name: "percent_ownership", Location: "form", Enum: []string{""}},
			{Name: "representative", Location: "form", Type: "boolean"},
			{Name: "title", Location: "form", Type: "string"},
		}},
		{Name: "ssn_last_4", Location: "form", Type: "string"},
		{Name: "us_cfpb_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "ethnicity_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ethnicity", Location: "form", Type: "array"},
				{Name: "ethnicity_other", Location: "form", Type: "string"},
			}},
			{Name: "race_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "race", Location: "form", Type: "array"},
				{Name: "race_other", Location: "form", Type: "string"},
			}},
			{Name: "self_identified_gender", Location: "form", Type: "string"},
		}},
		{Name: "verification", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "additional_document", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "back", Location: "form", Type: "string"},
				{Name: "front", Location: "form", Type: "string"},
			}},
			{Name: "document", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "back", Location: "form", Type: "string"},
				{Name: "front", Location: "form", Type: "string"},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/accounts/{account}/persons/{person}", OperationID: "DeleteAccountsAccountPersonsPerson", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "person", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/accounts/{account}/persons/{person}", OperationID: "GetAccountsAccountPersonsPerson", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "person", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/accounts/{account}/persons/{person}", OperationID: "PostAccountsAccountPersonsPerson", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "person", Location: "path", Required: true, Type: "string"},
		{Name: "additional_tos_acceptances", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Type: "integer"},
				{Name: "ip", Location: "form", Type: "string"},
				{Name: "user_agent", Location: "form", Type: "string", Enum: []string{""}},
			}},
		}},
		{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
		}},
		{Name: "address_kana", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "address_kanji", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "dob", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "day", Location: "form", Required: true, Type: "integer"},
			{Name: "month", Location: "form", Required: true, Type: "integer"},
			{Name: "year", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "company_authorization", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "passport", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
			{Name: "visa", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
		}},
		{Name: "email", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "first_name", Location: "form", Type: "string"},
		{Name: "first_name_kana", Location: "form", Type: "string"},
		{Name: "first_name_kanji", Location: "form", Type: "string"},
		{Name: "full_name_aliases", Location: "form", Enum: []string{""}},
		{Name: "gender", Location: "form", Type: "string"},
		{Name: "id_number", Location: "form", Type: "string"},
		{Name: "id_number_secondary", Location: "form", Type: "string"},
		{Name: "last_name", Location: "form", Type: "string"},
		{Name: "last_name_kana", Location: "form", Type: "string"},
		{Name: "last_name_kanji", Location: "form", Type: "string"},
		{Name: "maiden_name", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "nationality", Location: "form", Type: "string"},
		{Name: "person_token", Location: "form", Type: "string"},
		{Name: "phone", Location: "form", Type: "string"},
		{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
		{Name: "registered_address", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
		}},
		{Name: "relationship", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "authorizer", Location: "form", Type: "boolean"},
			{Name: "director", Location: "form", Type: "boolean"},
			{Name: "executive", Location: "form", Type: "boolean"},
			{Name: "legal_guardian", Location: "form", Type: "boolean"},
			{Name: "owner", Location: "form", Type: "boolean"},
			{Name: "percent_ownership", Location: "form", Enum: []string{""}},
			{Name: "representative", Location: "form", Type: "boolean"},
			{Name: "title", Location: "form", Type: "string"},
		}},
		{Name: "ssn_last_4", Location: "form", Type: "string"},
		{Name: "us_cfpb_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "ethnicity_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ethnicity", Location: "form", Type: "array"},
				{Name: "ethnicity_other", Location: "form", Type: "string"},
			}},
			{Name: "race_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "race", Location: "form", Type: "array"},
				{Name: "race_other", Location: "form", Type: "string"},
			}},
			{Name: "self_identified_gender", Location: "form", Type: "string"},
		}},
		{Name: "verification", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "additional_document", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "back", Location: "form", Type: "string"},
				{Name: "front", Location: "form", Type: "string"},
			}},
			{Name: "document", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "back", Location: "form", Type: "string"},
				{Name: "front", Location: "form", Type: "string"},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/accounts/{account}/reject", OperationID: "PostAccountsAccountReject", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "reason", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/apple_pay/domains", OperationID: "GetApplePayDomains", Params: []ParameterValidation{
		{Name: "domain_name", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/apple_pay/domains", OperationID: "PostApplePayDomains", Params: []ParameterValidation{
		{Name: "domain_name", Location: "form", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/apple_pay/domains/{domain}", OperationID: "DeleteApplePayDomainsDomain", Params: []ParameterValidation{
		{Name: "domain", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/apple_pay/domains/{domain}", OperationID: "GetApplePayDomainsDomain", Params: []ParameterValidation{
		{Name: "domain", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/application_fees", OperationID: "GetApplicationFees", Params: []ParameterValidation{
		{Name: "charge", Location: "query", Type: "string"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/application_fees/{fee}/refunds/{id}", OperationID: "GetApplicationFeesFeeRefundsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "fee", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/application_fees/{fee}/refunds/{id}", OperationID: "PostApplicationFeesFeeRefundsId", Params: []ParameterValidation{
		{Name: "fee", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/application_fees/{id}", OperationID: "GetApplicationFeesId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/application_fees/{id}/refund", OperationID: "PostApplicationFeesIdRefund", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "directive", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/application_fees/{id}/refunds", OperationID: "GetApplicationFeesIdRefunds", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/application_fees/{id}/refunds", OperationID: "PostApplicationFeesIdRefunds", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/apps/secrets", OperationID: "GetAppsSecrets", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "scope", Location: "query", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "type", Location: "query", Required: true, Type: "string", Enum: []string{"account", "user"}},
			{Name: "user", Location: "query", Type: "string"},
		}},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/apps/secrets", OperationID: "PostAppsSecrets", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "expires_at", Location: "form", Type: "integer"},
		{Name: "name", Location: "form", Required: true, Type: "string"},
		{Name: "payload", Location: "form", Required: true, Type: "string"},
		{Name: "scope", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "user"}},
			{Name: "user", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/apps/secrets/delete", OperationID: "PostAppsSecretsDelete", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "name", Location: "form", Required: true, Type: "string"},
		{Name: "scope", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "user"}},
			{Name: "user", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/apps/secrets/find", OperationID: "GetAppsSecretsFind", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "name", Location: "query", Required: true, Type: "string"},
		{Name: "scope", Location: "query", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "type", Location: "query", Required: true, Type: "string", Enum: []string{"account", "user"}},
			{Name: "user", Location: "query", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/balance", OperationID: "GetBalance", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/balance/history", OperationID: "GetBalanceHistory", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "currency", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "payout", Location: "query", Type: "string"},
		{Name: "source", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "type", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/balance/history/{id}", OperationID: "GetBalanceHistoryId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/balance_settings", OperationID: "GetBalanceSettings", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/balance_settings", OperationID: "PostBalanceSettings", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "payments", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "debit_negative_balances", Location: "form", Type: "boolean"},
			{Name: "payouts", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "minimum_balance_by_currency", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "schedule", Location: "form", Type: "object"},
				{Name: "statement_descriptor", Location: "form", Type: "string"},
			}},
			{Name: "settlement_timing", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "delay_days_override", Location: "form", Enum: []string{""}},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/balance_transactions", OperationID: "GetBalanceTransactions", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "currency", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "payout", Location: "query", Type: "string"},
		{Name: "source", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "type", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/balance_transactions/{id}", OperationID: "GetBalanceTransactionsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/billing/alerts", OperationID: "GetBillingAlerts", Params: []ParameterValidation{
		{Name: "alert_type", Location: "query", Type: "string", Enum: []string{"usage_threshold"}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "meter", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/alerts", OperationID: "PostBillingAlerts", Params: []ParameterValidation{
		{Name: "alert_type", Location: "form", Required: true, Type: "string", Enum: []string{"usage_threshold"}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "title", Location: "form", Required: true, Type: "string"},
		{Name: "usage_threshold", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "filters", Location: "form", Type: "array"},
			{Name: "gte", Location: "form", Required: true, Type: "integer"},
			{Name: "meter", Location: "form", Required: true, Type: "string"},
			{Name: "recurrence", Location: "form", Required: true, Type: "string", Enum: []string{"one_time"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/billing/alerts/{id}", OperationID: "GetBillingAlertsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/alerts/{id}/activate", OperationID: "PostBillingAlertsIdActivate", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/alerts/{id}/archive", OperationID: "PostBillingAlertsIdArchive", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/alerts/{id}/deactivate", OperationID: "PostBillingAlertsIdDeactivate", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/billing/credit_balance_summary", OperationID: "GetBillingCreditBalanceSummary", Params: []ParameterValidation{
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "filter", Location: "query", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "applicability_scope", Location: "query", Type: "object", Children: []ParameterValidation{
				{Name: "price_type", Location: "query", Type: "string", Enum: []string{"metered"}},
				{Name: "prices", Location: "query", Type: "array"},
			}},
			{Name: "credit_grant", Location: "query", Type: "string"},
			{Name: "type", Location: "query", Required: true, Type: "string", Enum: []string{"applicability_scope", "credit_grant"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/billing/credit_balance_transactions", OperationID: "GetBillingCreditBalanceTransactions", Params: []ParameterValidation{
		{Name: "credit_grant", Location: "query", Type: "string"},
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/billing/credit_balance_transactions/{id}", OperationID: "GetBillingCreditBalanceTransactionsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/billing/credit_grants", OperationID: "GetBillingCreditGrants", Params: []ParameterValidation{
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/credit_grants", OperationID: "PostBillingCreditGrants", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "monetary", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "currency", Location: "form", Required: true, Type: "string"},
				{Name: "value", Location: "form", Required: true, Type: "integer"},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"monetary"}},
		}},
		{Name: "applicability_config", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "scope", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "price_type", Location: "form", Type: "string", Enum: []string{"metered"}},
				{Name: "prices", Location: "form", Type: "array"},
			}},
		}},
		{Name: "category", Location: "form", Type: "string", Enum: []string{"paid", "promotional"}},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "effective_at", Location: "form", Type: "integer"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "expires_at", Location: "form", Type: "integer"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "priority", Location: "form", Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/billing/credit_grants/{id}", OperationID: "GetBillingCreditGrantsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/credit_grants/{id}", OperationID: "PostBillingCreditGrantsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "expires_at", Location: "form", Enum: []string{""}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/credit_grants/{id}/expire", OperationID: "PostBillingCreditGrantsIdExpire", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/credit_grants/{id}/void", OperationID: "PostBillingCreditGrantsIdVoid", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/meter_event_adjustments", OperationID: "PostBillingMeterEventAdjustments", Params: []ParameterValidation{
		{Name: "cancel", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "identifier", Location: "form", Type: "string"},
		}},
		{Name: "event_name", Location: "form", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"cancel"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/meter_events", OperationID: "PostBillingMeterEvents", Params: []ParameterValidation{
		{Name: "event_name", Location: "form", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "identifier", Location: "form", Type: "string"},
		{Name: "payload", Location: "form", Required: true, Type: "object", AdditionalProperties: true},
		{Name: "timestamp", Location: "form", Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/billing/meters", OperationID: "GetBillingMeters", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"active", "inactive"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/meters", OperationID: "PostBillingMeters", Params: []ParameterValidation{
		{Name: "customer_mapping", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "event_payload_key", Location: "form", Required: true, Type: "string"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"by_id"}},
		}},
		{Name: "default_aggregation", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "formula", Location: "form", Required: true, Type: "string", Enum: []string{"count", "last", "sum"}},
		}},
		{Name: "display_name", Location: "form", Required: true, Type: "string"},
		{Name: "event_name", Location: "form", Required: true, Type: "string"},
		{Name: "event_time_window", Location: "form", Type: "string", Enum: []string{"day", "hour"}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "value_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "event_payload_key", Location: "form", Required: true, Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/billing/meters/{id}", OperationID: "GetBillingMetersId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/meters/{id}", OperationID: "PostBillingMetersId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "display_name", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/meters/{id}/deactivate", OperationID: "PostBillingMetersIdDeactivate", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/billing/meters/{id}/event_summaries", OperationID: "GetBillingMetersIdEventSummaries", Params: []ParameterValidation{
		{Name: "customer", Location: "query", Required: true, Type: "string"},
		{Name: "end_time", Location: "query", Required: true, Type: "integer"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "start_time", Location: "query", Required: true, Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "value_grouping_window", Location: "query", Type: "string", Enum: []string{"day", "hour"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing/meters/{id}/reactivate", OperationID: "PostBillingMetersIdReactivate", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/billing_portal/configurations", OperationID: "GetBillingPortalConfigurations", Params: []ParameterValidation{
		{Name: "active", Location: "query", Type: "boolean"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "is_default", Location: "query", Type: "boolean"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing_portal/configurations", OperationID: "PostBillingPortalConfigurations", Params: []ParameterValidation{
		{Name: "business_profile", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "headline", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "privacy_policy_url", Location: "form", Type: "string"},
			{Name: "terms_of_service_url", Location: "form", Type: "string"},
		}},
		{Name: "default_return_url", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "features", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "customer_update", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "allowed_updates", Location: "form", Enum: []string{""}},
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			}},
			{Name: "invoice_history", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			}},
			{Name: "payment_method_update", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "payment_method_configuration", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "subscription_cancel", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "cancellation_reason", Location: "form", Type: "object"},
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "mode", Location: "form", Type: "string", Enum: []string{"at_period_end", "immediately"}},
				{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
			}},
			{Name: "subscription_update", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "billing_cycle_anchor", Location: "form", Type: "string", Enum: []string{"now", "unchanged"}},
				{Name: "default_allowed_updates", Location: "form", Enum: []string{""}},
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "products", Location: "form", Enum: []string{""}},
				{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
				{Name: "schedule_at_period_end", Location: "form", Type: "object"},
				{Name: "trial_update_behavior", Location: "form", Type: "string", Enum: []string{"continue_trial", "end_trial"}},
			}},
		}},
		{Name: "login_page", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "name", Location: "form", Type: "string", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/billing_portal/configurations/{configuration}", OperationID: "GetBillingPortalConfigurationsConfiguration", Params: []ParameterValidation{
		{Name: "configuration", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing_portal/configurations/{configuration}", OperationID: "PostBillingPortalConfigurationsConfiguration", Params: []ParameterValidation{
		{Name: "configuration", Location: "path", Required: true, Type: "string"},
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "business_profile", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "headline", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "privacy_policy_url", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "terms_of_service_url", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "default_return_url", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "features", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "customer_update", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "allowed_updates", Location: "form", Enum: []string{""}},
				{Name: "enabled", Location: "form", Type: "boolean"},
			}},
			{Name: "invoice_history", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			}},
			{Name: "payment_method_update", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "payment_method_configuration", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "subscription_cancel", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "cancellation_reason", Location: "form", Type: "object"},
				{Name: "enabled", Location: "form", Type: "boolean"},
				{Name: "mode", Location: "form", Type: "string", Enum: []string{"at_period_end", "immediately"}},
				{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
			}},
			{Name: "subscription_update", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "billing_cycle_anchor", Location: "form", Type: "string", Enum: []string{"now", "unchanged"}},
				{Name: "default_allowed_updates", Location: "form", Enum: []string{""}},
				{Name: "enabled", Location: "form", Type: "boolean"},
				{Name: "products", Location: "form", Enum: []string{""}},
				{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
				{Name: "schedule_at_period_end", Location: "form", Type: "object"},
				{Name: "trial_update_behavior", Location: "form", Type: "string", Enum: []string{"continue_trial", "end_trial"}},
			}},
		}},
		{Name: "login_page", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "name", Location: "form", Type: "string", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/billing_portal/sessions", OperationID: "PostBillingPortalSessions", Params: []ParameterValidation{
		{Name: "configuration", Location: "form", Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "flow_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "after_completion", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "hosted_confirmation", Location: "form", Type: "object"},
				{Name: "redirect", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"hosted_confirmation", "portal_homepage", "redirect"}},
			}},
			{Name: "subscription_cancel", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "retention", Location: "form", Type: "object"},
				{Name: "subscription", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "subscription_update", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "subscription", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "subscription_update_confirm", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "discounts", Location: "form", Type: "array"},
				{Name: "items", Location: "form", Required: true, Type: "array"},
				{Name: "subscription", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"payment_method_update", "subscription_cancel", "subscription_update", "subscription_update_confirm"}},
		}},
		{Name: "locale", Location: "form", Type: "string", Enum: []string{"auto", "bg", "cs", "da", "de", "el", "en", "en-AU", "en-CA", "en-GB", "en-IE", "en-IN", "en-NZ", "en-SG", "es", "es-419", "et", "fi", "fil", "fr", "fr-CA", "hr", "hu", "id", "it", "ja", "ko", "lt", "lv", "ms", "mt", "nb", "nl", "pl", "pt", "pt-BR", "ro", "ru", "sk", "sl", "sv", "th", "tr", "vi", "zh", "zh-HK", "zh-TW"}},
		{Name: "on_behalf_of", Location: "form", Type: "string"},
		{Name: "return_url", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/charges", OperationID: "GetCharges", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "payment_intent", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "transfer_group", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/charges", OperationID: "PostCharges", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "application_fee", Location: "form", Type: "integer"},
		{Name: "application_fee_amount", Location: "form", Type: "integer"},
		{Name: "capture", Location: "form", Type: "boolean"},
		{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address_city", Location: "form", Type: "string"},
			{Name: "address_country", Location: "form", Type: "string"},
			{Name: "address_line1", Location: "form", Type: "string"},
			{Name: "address_line2", Location: "form", Type: "string"},
			{Name: "address_state", Location: "form", Type: "string"},
			{Name: "address_zip", Location: "form", Type: "string"},
			{Name: "cvc", Location: "form", Type: "string"},
			{Name: "exp_month", Location: "form", Required: true, Type: "integer"},
			{Name: "exp_year", Location: "form", Required: true, Type: "integer"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "number", Location: "form", Required: true, Type: "string"},
			{Name: "object", Location: "form", Type: "string", Enum: []string{"card"}},
		}},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "destination", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Required: true, Type: "string"},
			{Name: "amount", Location: "form", Type: "integer"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "on_behalf_of", Location: "form", Type: "string"},
		{Name: "radar_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "session", Location: "form", Type: "string"},
		}},
		{Name: "receipt_email", Location: "form", Type: "string"},
		{Name: "shipping", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "carrier", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
			{Name: "tracking_number", Location: "form", Type: "string"},
		}},
		{Name: "source", Location: "form", Type: "string"},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
		{Name: "statement_descriptor_suffix", Location: "form", Type: "string"},
		{Name: "transfer_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount", Location: "form", Type: "integer"},
			{Name: "destination", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "transfer_group", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/charges/search", OperationID: "GetChargesSearch", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "page", Location: "query", Type: "string"},
		{Name: "query", Location: "query", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/charges/{charge}", OperationID: "GetChargesCharge", Params: []ParameterValidation{
		{Name: "charge", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/charges/{charge}", OperationID: "PostChargesCharge", Params: []ParameterValidation{
		{Name: "charge", Location: "path", Required: true, Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "fraud_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "user_report", Location: "form", Required: true, Type: "string", Enum: []string{"", "fraudulent", "safe"}},
		}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "receipt_email", Location: "form", Type: "string"},
		{Name: "shipping", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "carrier", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
			{Name: "tracking_number", Location: "form", Type: "string"},
		}},
		{Name: "transfer_group", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/charges/{charge}/capture", OperationID: "PostChargesChargeCapture", Params: []ParameterValidation{
		{Name: "charge", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "application_fee", Location: "form", Type: "integer"},
		{Name: "application_fee_amount", Location: "form", Type: "integer"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "receipt_email", Location: "form", Type: "string"},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
		{Name: "statement_descriptor_suffix", Location: "form", Type: "string"},
		{Name: "transfer_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount", Location: "form", Type: "integer"},
		}},
		{Name: "transfer_group", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/charges/{charge}/dispute", OperationID: "GetChargesChargeDispute", Params: []ParameterValidation{
		{Name: "charge", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/charges/{charge}/dispute", OperationID: "PostChargesChargeDispute", Params: []ParameterValidation{
		{Name: "charge", Location: "path", Required: true, Type: "string"},
		{Name: "evidence", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "access_activity_log", Location: "form", Type: "string"},
			{Name: "billing_address", Location: "form", Type: "string"},
			{Name: "cancellation_policy", Location: "form", Type: "string"},
			{Name: "cancellation_policy_disclosure", Location: "form", Type: "string"},
			{Name: "cancellation_rebuttal", Location: "form", Type: "string"},
			{Name: "customer_communication", Location: "form", Type: "string"},
			{Name: "customer_email_address", Location: "form", Type: "string"},
			{Name: "customer_name", Location: "form", Type: "string"},
			{Name: "customer_purchase_ip", Location: "form", Type: "string"},
			{Name: "customer_signature", Location: "form", Type: "string"},
			{Name: "duplicate_charge_documentation", Location: "form", Type: "string"},
			{Name: "duplicate_charge_explanation", Location: "form", Type: "string"},
			{Name: "duplicate_charge_id", Location: "form", Type: "string"},
			{Name: "enhanced_evidence", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "visa_compelling_evidence_3", Location: "form", Type: "object"},
				{Name: "visa_compliance", Location: "form", Type: "object"},
			}},
			{Name: "product_description", Location: "form", Type: "string"},
			{Name: "receipt", Location: "form", Type: "string"},
			{Name: "refund_policy", Location: "form", Type: "string"},
			{Name: "refund_policy_disclosure", Location: "form", Type: "string"},
			{Name: "refund_refusal_explanation", Location: "form", Type: "string"},
			{Name: "service_date", Location: "form", Type: "string"},
			{Name: "service_documentation", Location: "form", Type: "string"},
			{Name: "shipping_address", Location: "form", Type: "string"},
			{Name: "shipping_carrier", Location: "form", Type: "string"},
			{Name: "shipping_date", Location: "form", Type: "string"},
			{Name: "shipping_documentation", Location: "form", Type: "string"},
			{Name: "shipping_tracking_number", Location: "form", Type: "string"},
			{Name: "uncategorized_file", Location: "form", Type: "string"},
			{Name: "uncategorized_text", Location: "form", Type: "string"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "submit", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/charges/{charge}/dispute/close", OperationID: "PostChargesChargeDisputeClose", Params: []ParameterValidation{
		{Name: "charge", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/charges/{charge}/refund", OperationID: "PostChargesChargeRefund", Params: []ParameterValidation{
		{Name: "charge", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "instructions_email", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "payment_intent", Location: "form", Type: "string"},
		{Name: "reason", Location: "form", Type: "string", Enum: []string{"duplicate", "fraudulent", "requested_by_customer"}},
		{Name: "refund_application_fee", Location: "form", Type: "boolean"},
		{Name: "reverse_transfer", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/charges/{charge}/refunds", OperationID: "GetChargesChargeRefunds", Params: []ParameterValidation{
		{Name: "charge", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/charges/{charge}/refunds", OperationID: "PostChargesChargeRefunds", Params: []ParameterValidation{
		{Name: "charge", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "instructions_email", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "origin", Location: "form", Type: "string", Enum: []string{"customer_balance"}},
		{Name: "payment_intent", Location: "form", Type: "string"},
		{Name: "reason", Location: "form", Type: "string", Enum: []string{"duplicate", "fraudulent", "requested_by_customer"}},
		{Name: "refund_application_fee", Location: "form", Type: "boolean"},
		{Name: "reverse_transfer", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/charges/{charge}/refunds/{refund}", OperationID: "GetChargesChargeRefundsRefund", Params: []ParameterValidation{
		{Name: "charge", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "refund", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/charges/{charge}/refunds/{refund}", OperationID: "PostChargesChargeRefundsRefund", Params: []ParameterValidation{
		{Name: "charge", Location: "path", Required: true, Type: "string"},
		{Name: "refund", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/checkout/sessions", OperationID: "GetCheckoutSessions", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "customer_details", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "email", Location: "query", Required: true, Type: "string"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "payment_intent", Location: "query", Type: "string"},
		{Name: "payment_link", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"complete", "expired", "open"}},
		{Name: "subscription", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/checkout/sessions", OperationID: "PostCheckoutSessions", Params: []ParameterValidation{
		{Name: "adaptive_pricing", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Type: "boolean"},
		}},
		{Name: "after_expiration", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "recovery", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "allow_promotion_codes", Location: "form", Type: "boolean"},
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			}},
		}},
		{Name: "allow_promotion_codes", Location: "form", Type: "boolean"},
		{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "liability", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "billing_address_collection", Location: "form", Type: "string", Enum: []string{"auto", "required"}},
		{Name: "branding_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "background_color", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "border_style", Location: "form", Type: "string", Enum: []string{"", "pill", "rectangular", "rounded"}},
			{Name: "button_color", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "display_name", Location: "form", Type: "string"},
			{Name: "font_family", Location: "form", Type: "string", Enum: []string{"", "be_vietnam_pro", "bitter", "chakra_petch", "default", "hahmlet", "inconsolata", "inter", "lato", "lora", "m_plus_1_code", "montserrat", "noto_sans", "noto_sans_jp", "noto_serif", "nunito", "open_sans", "pridi", "pt_sans", "pt_serif", "raleway", "roboto", "roboto_slab", "source_sans_pro", "titillium_web", "ubuntu_mono", "zen_maru_gothic"}},
			{Name: "icon", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "file", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"file", "url"}},
				{Name: "url", Location: "form", Type: "string"},
			}},
			{Name: "logo", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "file", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"file", "url"}},
				{Name: "url", Location: "form", Type: "string"},
			}},
		}},
		{Name: "cancel_url", Location: "form", Type: "string"},
		{Name: "client_reference_id", Location: "form", Type: "string"},
		{Name: "consent_collection", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "payment_method_reuse_agreement", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "position", Location: "form", Required: true, Type: "string", Enum: []string{"auto", "hidden"}},
			}},
			{Name: "promotions", Location: "form", Type: "string", Enum: []string{"auto", "none"}},
			{Name: "terms_of_service", Location: "form", Type: "string", Enum: []string{"none", "required"}},
		}},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "custom_fields", Location: "form", Type: "array"},
		{Name: "custom_text", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "after_submit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "message", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "shipping_address", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "message", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "submit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "message", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "terms_of_service_acceptance", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "message", Location: "form", Required: true, Type: "string"},
			}},
		}},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "customer_creation", Location: "form", Type: "string", Enum: []string{"always", "if_required"}},
		{Name: "customer_email", Location: "form", Type: "string"},
		{Name: "customer_update", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "string", Enum: []string{"auto", "never"}},
			{Name: "name", Location: "form", Type: "string", Enum: []string{"auto", "never"}},
			{Name: "shipping", Location: "form", Type: "string", Enum: []string{"auto", "never"}},
		}},
		{Name: "discounts", Location: "form", Type: "array"},
		{Name: "excluded_payment_method_types", Location: "form", Type: "array"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "expires_at", Location: "form", Type: "integer"},
		{Name: "integration_identifier", Location: "form", Type: "string"},
		{Name: "invoice_creation", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "invoice_data", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_tax_ids", Location: "form", Enum: []string{""}},
				{Name: "custom_fields", Location: "form", Enum: []string{""}},
				{Name: "description", Location: "form", Type: "string"},
				{Name: "footer", Location: "form", Type: "string"},
				{Name: "issuer", Location: "form", Type: "object"},
				{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
				{Name: "rendering_options", Location: "form", Type: "object", Enum: []string{""}},
			}},
		}},
		{Name: "line_items", Location: "form", Type: "array"},
		{Name: "locale", Location: "form", Type: "string", Enum: []string{"auto", "bg", "cs", "da", "de", "el", "en", "en-GB", "es", "es-419", "et", "fi", "fil", "fr", "fr-CA", "hr", "hu", "id", "it", "ja", "ko", "lt", "lv", "ms", "mt", "nb", "nl", "pl", "pt", "pt-BR", "ro", "ru", "sk", "sl", "sv", "th", "tr", "vi", "zh", "zh-HK", "zh-TW"}},
		{Name: "managed_payments", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Type: "boolean"},
		}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "mode", Location: "form", Type: "string", Enum: []string{"payment", "setup", "subscription"}},
		{Name: "name_collection", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "business", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "optional", Location: "form", Type: "boolean"},
			}},
			{Name: "individual", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "optional", Location: "form", Type: "boolean"},
			}},
		}},
		{Name: "optional_items", Location: "form", Type: "array"},
		{Name: "origin_context", Location: "form", Type: "string", Enum: []string{"mobile_app", "web"}},
		{Name: "payment_intent_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "application_fee_amount", Location: "form", Type: "integer"},
			{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"automatic", "automatic_async", "manual"}},
			{Name: "description", Location: "form", Type: "string"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "on_behalf_of", Location: "form", Type: "string"},
			{Name: "receipt_email", Location: "form", Type: "string"},
			{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"off_session", "on_session"}},
			{Name: "shipping", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Required: true, Type: "object"},
				{Name: "carrier", Location: "form", Type: "string"},
				{Name: "name", Location: "form", Required: true, Type: "string"},
				{Name: "phone", Location: "form", Type: "string"},
				{Name: "tracking_number", Location: "form", Type: "string"},
			}},
			{Name: "statement_descriptor", Location: "form", Type: "string"},
			{Name: "statement_descriptor_suffix", Location: "form", Type: "string"},
			{Name: "transfer_data", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "amount", Location: "form", Type: "integer"},
				{Name: "destination", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "transfer_group", Location: "form", Type: "string"},
		}},
		{Name: "payment_method_collection", Location: "form", Type: "string", Enum: []string{"always", "if_required"}},
		{Name: "payment_method_configuration", Location: "form", Type: "string"},
		{Name: "payment_method_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "allow_redisplay", Location: "form", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
		}},
		{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "currency", Location: "form", Type: "string", Enum: []string{"cad", "usd"}},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant", "microdeposits"}},
			}},
			{Name: "affirm", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "afterpay_clearpay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "alipay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "alma", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
			}},
			{Name: "amazon_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session"}},
			}},
			{Name: "au_becs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "bancontact", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "billie", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
			}},
			{Name: "boleto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "expires_after_days", Location: "form", Type: "integer"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session", "on_session"}},
			}},
			{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
				{Name: "installments", Location: "form", Type: "object"},
				{Name: "request_extended_authorization", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_incremental_authorization", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_multicapture", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_overcapture", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_three_d_secure", Location: "form", Type: "string", Enum: []string{"any", "automatic", "challenge"}},
				{Name: "restrictions", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"off_session", "on_session"}},
				{Name: "statement_descriptor_suffix_kana", Location: "form", Type: "string"},
				{Name: "statement_descriptor_suffix_kanji", Location: "form", Type: "string"},
			}},
			{Name: "cashapp", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session", "on_session"}},
			}},
			{Name: "crypto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "customer_balance", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank_transfer", Location: "form", Type: "object"},
				{Name: "funding_type", Location: "form", Type: "string", Enum: []string{"bank_transfer"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "demo_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session"}},
			}},
			{Name: "eps", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "fpx", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "giropay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "grabpay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "ideal", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "kakao_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session"}},
			}},
			{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
				{Name: "subscriptions", Location: "form", Enum: []string{""}},
			}},
			{Name: "konbini", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "expires_after_days", Location: "form", Type: "integer"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "kr_card", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session"}},
			}},
			{Name: "link", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session"}},
			}},
			{Name: "mobilepay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "multibanco", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "naver_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session"}},
			}},
			{Name: "oxxo", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "expires_after_days", Location: "form", Type: "integer"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "p24", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
				{Name: "tos_shown_and_accepted", Location: "form", Type: "boolean"},
			}},
			{Name: "pay_by_bank", Location: "form", Type: "object"},
			{Name: "payco", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
			}},
			{Name: "paynow", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "paypal", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "preferred_locale", Location: "form", Type: "string", Enum: []string{"cs-CZ", "da-DK", "de-AT", "de-DE", "de-LU", "el-GR", "en-GB", "en-US", "es-ES", "fi-FI", "fr-BE", "fr-FR", "fr-LU", "hu-HU", "it-IT", "nl-BE", "nl-NL", "pl-PL", "pt-PT", "sk-SK", "sv-SE"}},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "risk_correlation_id", Location: "form", Type: "string"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session"}},
			}},
			{Name: "pix", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "amount_includes_iof", Location: "form", Type: "string", Enum: []string{"always", "never"}},
				{Name: "expires_after_seconds", Location: "form", Type: "integer"},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session"}},
			}},
			{Name: "revolut_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session"}},
			}},
			{Name: "samsung_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
			}},
			{Name: "satispay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual"}},
			}},
			{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "sofort", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "swish", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "reference", Location: "form", Type: "string"},
			}},
			{Name: "twint", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "financial_connections", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant"}},
			}},
			{Name: "wechat_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "app_id", Location: "form", Type: "string"},
				{Name: "client", Location: "form", Required: true, Type: "string", Enum: []string{"android", "ios", "web"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
		}},
		{Name: "payment_method_types", Location: "form", Type: "array"},
		{Name: "permissions", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "update_shipping_details", Location: "form", Type: "string", Enum: []string{"client_only", "server_only"}},
		}},
		{Name: "phone_number_collection", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "redirect_on_completion", Location: "form", Type: "string", Enum: []string{"always", "if_required", "never"}},
		{Name: "return_url", Location: "form", Type: "string"},
		{Name: "saved_payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "allow_redisplay_filters", Location: "form", Type: "array"},
			{Name: "payment_method_remove", Location: "form", Type: "string", Enum: []string{"disabled", "enabled"}},
			{Name: "payment_method_save", Location: "form", Type: "string", Enum: []string{"disabled", "enabled"}},
		}},
		{Name: "setup_intent_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "description", Location: "form", Type: "string"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "on_behalf_of", Location: "form", Type: "string"},
		}},
		{Name: "shipping_address_collection", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "allowed_countries", Location: "form", Required: true, Type: "array"},
		}},
		{Name: "shipping_options", Location: "form", Type: "array"},
		{Name: "submit_type", Location: "form", Type: "string", Enum: []string{"auto", "book", "donate", "pay", "subscribe"}},
		{Name: "subscription_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "application_fee_percent", Location: "form", Type: "number"},
			{Name: "billing_cycle_anchor", Location: "form", Type: "integer"},
			{Name: "billing_mode", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "flexible", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"classic", "flexible"}},
			}},
			{Name: "default_tax_rates", Location: "form", Type: "array"},
			{Name: "description", Location: "form", Type: "string"},
			{Name: "invoice_settings", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "issuer", Location: "form", Type: "object"},
			}},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "on_behalf_of", Location: "form", Type: "string"},
			{Name: "pending_invoice_item_interval", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "interval", Location: "form", Required: true, Type: "string", Enum: []string{"day", "month", "week", "year"}},
				{Name: "interval_count", Location: "form", Type: "integer"},
			}},
			{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"create_prorations", "none"}},
			{Name: "transfer_data", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "amount_percent", Location: "form", Type: "number"},
				{Name: "destination", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "trial_end", Location: "form", Type: "integer"},
			{Name: "trial_period_days", Location: "form", Type: "integer"},
			{Name: "trial_settings", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "end_behavior", Location: "form", Required: true, Type: "object"},
			}},
		}},
		{Name: "success_url", Location: "form", Type: "string"},
		{Name: "tax_id_collection", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "required", Location: "form", Type: "string", Enum: []string{"if_supported", "never"}},
		}},
		{Name: "ui_mode", Location: "form", Type: "string", Enum: []string{"elements", "embedded_page", "hosted_page"}},
		{Name: "wallet_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "link", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "display", Location: "form", Type: "string", Enum: []string{"auto", "never"}},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/checkout/sessions/{session}", OperationID: "GetCheckoutSessionsSession", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "session", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/checkout/sessions/{session}", OperationID: "PostCheckoutSessionsSession", Params: []ParameterValidation{
		{Name: "session", Location: "path", Required: true, Type: "string"},
		{Name: "collected_information", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "shipping_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Required: true, Type: "object"},
				{Name: "name", Location: "form", Required: true, Type: "string"},
			}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "line_items", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "shipping_options", Location: "form", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/checkout/sessions/{session}/expire", OperationID: "PostCheckoutSessionsSessionExpire", Params: []ParameterValidation{
		{Name: "session", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/checkout/sessions/{session}/line_items", OperationID: "GetCheckoutSessionsSessionLineItems", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "session", Location: "path", Required: true, Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/climate/orders", OperationID: "GetClimateOrders", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/climate/orders", OperationID: "PostClimateOrders", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "beneficiary", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "public_name", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "metric_tons", Location: "form", Type: "string"},
		{Name: "product", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/climate/orders/{order}", OperationID: "GetClimateOrdersOrder", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "order", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/climate/orders/{order}", OperationID: "PostClimateOrdersOrder", Params: []ParameterValidation{
		{Name: "order", Location: "path", Required: true, Type: "string"},
		{Name: "beneficiary", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "public_name", Location: "form", Required: true, Type: "string", Enum: []string{""}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/climate/orders/{order}/cancel", OperationID: "PostClimateOrdersOrderCancel", Params: []ParameterValidation{
		{Name: "order", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/climate/products", OperationID: "GetClimateProducts", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/climate/products/{product}", OperationID: "GetClimateProductsProduct", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "product", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/climate/suppliers", OperationID: "GetClimateSuppliers", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/climate/suppliers/{supplier}", OperationID: "GetClimateSuppliersSupplier", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "supplier", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/confirmation_tokens/{confirmation_token}", OperationID: "GetConfirmationTokensConfirmationToken", Params: []ParameterValidation{
		{Name: "confirmation_token", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/country_specs", OperationID: "GetCountrySpecs", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/country_specs/{country}", OperationID: "GetCountrySpecsCountry", Params: []ParameterValidation{
		{Name: "country", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/coupons", OperationID: "GetCoupons", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/coupons", OperationID: "PostCoupons", Params: []ParameterValidation{
		{Name: "amount_off", Location: "form", Type: "integer"},
		{Name: "applies_to", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "products", Location: "form", Type: "array"},
		}},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "currency_options", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "duration", Location: "form", Type: "string", Enum: []string{"forever", "once", "repeating"}},
		{Name: "duration_in_months", Location: "form", Type: "integer"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "id", Location: "form", Type: "string"},
		{Name: "max_redemptions", Location: "form", Type: "integer"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "percent_off", Location: "form", Type: "number"},
		{Name: "redeem_by", Location: "form", Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/coupons/{coupon}", OperationID: "DeleteCouponsCoupon", Params: []ParameterValidation{
		{Name: "coupon", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/coupons/{coupon}", OperationID: "GetCouponsCoupon", Params: []ParameterValidation{
		{Name: "coupon", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/coupons/{coupon}", OperationID: "PostCouponsCoupon", Params: []ParameterValidation{
		{Name: "coupon", Location: "path", Required: true, Type: "string"},
		{Name: "currency_options", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "name", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/credit_notes", OperationID: "GetCreditNotes", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "invoice", Location: "query", Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/credit_notes", OperationID: "PostCreditNotes", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "credit_amount", Location: "form", Type: "integer"},
		{Name: "effective_at", Location: "form", Type: "integer"},
		{Name: "email_type", Location: "form", Type: "string", Enum: []string{"credit_note", "none"}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "invoice", Location: "form", Required: true, Type: "string"},
		{Name: "lines", Location: "form", Type: "array"},
		{Name: "memo", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "out_of_band_amount", Location: "form", Type: "integer"},
		{Name: "reason", Location: "form", Type: "string", Enum: []string{"duplicate", "fraudulent", "order_change", "product_unsatisfactory"}},
		{Name: "refund_amount", Location: "form", Type: "integer"},
		{Name: "refunds", Location: "form", Type: "array"},
		{Name: "shipping_cost", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "shipping_rate", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/credit_notes/preview", OperationID: "GetCreditNotesPreview", Params: []ParameterValidation{
		{Name: "amount", Location: "query", Type: "integer"},
		{Name: "credit_amount", Location: "query", Type: "integer"},
		{Name: "effective_at", Location: "query", Type: "integer"},
		{Name: "email_type", Location: "query", Type: "string", Enum: []string{"credit_note", "none"}},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "invoice", Location: "query", Required: true, Type: "string"},
		{Name: "lines", Location: "query", Type: "array"},
		{Name: "memo", Location: "query", Type: "string"},
		{Name: "metadata", Location: "query", Type: "object", AdditionalProperties: true},
		{Name: "out_of_band_amount", Location: "query", Type: "integer"},
		{Name: "reason", Location: "query", Type: "string", Enum: []string{"duplicate", "fraudulent", "order_change", "product_unsatisfactory"}},
		{Name: "refund_amount", Location: "query", Type: "integer"},
		{Name: "refunds", Location: "query", Type: "array"},
		{Name: "shipping_cost", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "shipping_rate", Location: "query", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/credit_notes/preview/lines", OperationID: "GetCreditNotesPreviewLines", Params: []ParameterValidation{
		{Name: "amount", Location: "query", Type: "integer"},
		{Name: "credit_amount", Location: "query", Type: "integer"},
		{Name: "effective_at", Location: "query", Type: "integer"},
		{Name: "email_type", Location: "query", Type: "string", Enum: []string{"credit_note", "none"}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "invoice", Location: "query", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "lines", Location: "query", Type: "array"},
		{Name: "memo", Location: "query", Type: "string"},
		{Name: "metadata", Location: "query", Type: "object", AdditionalProperties: true},
		{Name: "out_of_band_amount", Location: "query", Type: "integer"},
		{Name: "reason", Location: "query", Type: "string", Enum: []string{"duplicate", "fraudulent", "order_change", "product_unsatisfactory"}},
		{Name: "refund_amount", Location: "query", Type: "integer"},
		{Name: "refunds", Location: "query", Type: "array"},
		{Name: "shipping_cost", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "shipping_rate", Location: "query", Type: "string"},
		}},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/credit_notes/{credit_note}/lines", OperationID: "GetCreditNotesCreditNoteLines", Params: []ParameterValidation{
		{Name: "credit_note", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/credit_notes/{id}", OperationID: "GetCreditNotesId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/credit_notes/{id}", OperationID: "PostCreditNotesId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "memo", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/credit_notes/{id}/void", OperationID: "PostCreditNotesIdVoid", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customer_sessions", OperationID: "PostCustomerSessions", Params: []ParameterValidation{
		{Name: "components", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "buy_button", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			}},
			{Name: "customer_sheet", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "mobile_payment_element", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "payment_element", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "features", Location: "form", Type: "object"},
			}},
			{Name: "pricing_table", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			}},
		}},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers", OperationID: "GetCustomers", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "email", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "test_clock", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers", OperationID: "PostCustomers", Params: []ParameterValidation{
		{Name: "address", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
		}},
		{Name: "balance", Location: "form", Type: "integer"},
		{Name: "business_name", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "cash_balance", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "settings", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "reconciliation_mode", Location: "form", Type: "string", Enum: []string{"automatic", "manual", "merchant_default"}},
			}},
		}},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "email", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "individual_name", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "invoice_prefix", Location: "form", Type: "string"},
		{Name: "invoice_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "custom_fields", Location: "form", Enum: []string{""}},
			{Name: "default_payment_method", Location: "form", Type: "string"},
			{Name: "footer", Location: "form", Type: "string"},
			{Name: "rendering_options", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount_tax_display", Location: "form", Type: "string", Enum: []string{"", "exclude_tax", "include_inclusive_tax"}},
				{Name: "template", Location: "form", Type: "string"},
			}},
		}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "next_invoice_sequence", Location: "form", Type: "integer"},
		{Name: "payment_method", Location: "form", Type: "string"},
		{Name: "phone", Location: "form", Type: "string"},
		{Name: "preferred_locales", Location: "form", Type: "array"},
		{Name: "shipping", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
		{Name: "source", Location: "form", Type: "string"},
		{Name: "tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "ip_address", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "validate_location", Location: "form", Type: "string", Enum: []string{"deferred", "immediately"}},
		}},
		{Name: "tax_exempt", Location: "form", Type: "string", Enum: []string{"", "exempt", "none", "reverse"}},
		{Name: "tax_id_data", Location: "form", Type: "array"},
		{Name: "test_clock", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/search", OperationID: "GetCustomersSearch", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "page", Location: "query", Type: "string"},
		{Name: "query", Location: "query", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/customers/{customer}", OperationID: "DeleteCustomersCustomer", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}", OperationID: "GetCustomersCustomer", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}", OperationID: "PostCustomersCustomer", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "address", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
		}},
		{Name: "balance", Location: "form", Type: "integer"},
		{Name: "bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_holder_name", Location: "form", Type: "string"},
			{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
			{Name: "account_number", Location: "form", Required: true, Type: "string"},
			{Name: "country", Location: "form", Required: true, Type: "string"},
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "object", Location: "form", Type: "string", Enum: []string{"bank_account"}},
			{Name: "routing_number", Location: "form", Type: "string"},
		}},
		{Name: "business_name", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address_city", Location: "form", Type: "string"},
			{Name: "address_country", Location: "form", Type: "string"},
			{Name: "address_line1", Location: "form", Type: "string"},
			{Name: "address_line2", Location: "form", Type: "string"},
			{Name: "address_state", Location: "form", Type: "string"},
			{Name: "address_zip", Location: "form", Type: "string"},
			{Name: "cvc", Location: "form", Type: "string"},
			{Name: "exp_month", Location: "form", Required: true, Type: "integer"},
			{Name: "exp_year", Location: "form", Required: true, Type: "integer"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "number", Location: "form", Required: true, Type: "string"},
			{Name: "object", Location: "form", Type: "string", Enum: []string{"card"}},
		}},
		{Name: "cash_balance", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "settings", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "reconciliation_mode", Location: "form", Type: "string", Enum: []string{"automatic", "manual", "merchant_default"}},
			}},
		}},
		{Name: "default_alipay_account", Location: "form", Type: "string"},
		{Name: "default_bank_account", Location: "form", Type: "string"},
		{Name: "default_card", Location: "form", Type: "string"},
		{Name: "default_source", Location: "form", Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "email", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "individual_name", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "invoice_prefix", Location: "form", Type: "string"},
		{Name: "invoice_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "custom_fields", Location: "form", Enum: []string{""}},
			{Name: "default_payment_method", Location: "form", Type: "string"},
			{Name: "footer", Location: "form", Type: "string"},
			{Name: "rendering_options", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount_tax_display", Location: "form", Type: "string", Enum: []string{"", "exclude_tax", "include_inclusive_tax"}},
				{Name: "template", Location: "form", Type: "string"},
			}},
		}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "next_invoice_sequence", Location: "form", Type: "integer"},
		{Name: "phone", Location: "form", Type: "string"},
		{Name: "preferred_locales", Location: "form", Type: "array"},
		{Name: "shipping", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
		{Name: "source", Location: "form", Type: "string"},
		{Name: "tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "ip_address", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "validate_location", Location: "form", Type: "string", Enum: []string{"auto", "deferred", "immediately"}},
		}},
		{Name: "tax_exempt", Location: "form", Type: "string", Enum: []string{"", "exempt", "none", "reverse"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/balance_transactions", OperationID: "GetCustomersCustomerBalanceTransactions", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "invoice", Location: "query", Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/balance_transactions", OperationID: "PostCustomersCustomerBalanceTransactions", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Required: true, Type: "integer"},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/balance_transactions/{transaction}", OperationID: "GetCustomersCustomerBalanceTransactionsTransaction", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "transaction", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/balance_transactions/{transaction}", OperationID: "PostCustomersCustomerBalanceTransactionsTransaction", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "transaction", Location: "path", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/bank_accounts", OperationID: "GetCustomersCustomerBankAccounts", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/bank_accounts", OperationID: "PostCustomersCustomerBankAccounts", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "alipay_account", Location: "form", Type: "string"},
		{Name: "bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_holder_name", Location: "form", Type: "string"},
			{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
			{Name: "account_number", Location: "form", Required: true, Type: "string"},
			{Name: "country", Location: "form", Required: true, Type: "string"},
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "object", Location: "form", Type: "string", Enum: []string{"bank_account"}},
			{Name: "routing_number", Location: "form", Type: "string"},
		}},
		{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address_city", Location: "form", Type: "string"},
			{Name: "address_country", Location: "form", Type: "string"},
			{Name: "address_line1", Location: "form", Type: "string"},
			{Name: "address_line2", Location: "form", Type: "string"},
			{Name: "address_state", Location: "form", Type: "string"},
			{Name: "address_zip", Location: "form", Type: "string"},
			{Name: "cvc", Location: "form", Type: "string"},
			{Name: "exp_month", Location: "form", Required: true, Type: "integer"},
			{Name: "exp_year", Location: "form", Required: true, Type: "integer"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "number", Location: "form", Required: true, Type: "string"},
			{Name: "object", Location: "form", Type: "string", Enum: []string{"card"}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "source", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/customers/{customer}/bank_accounts/{id}", OperationID: "DeleteCustomersCustomerBankAccountsId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/bank_accounts/{id}", OperationID: "GetCustomersCustomerBankAccountsId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/bank_accounts/{id}", OperationID: "PostCustomersCustomerBankAccountsId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "account_holder_name", Location: "form", Type: "string"},
		{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
		{Name: "address_city", Location: "form", Type: "string"},
		{Name: "address_country", Location: "form", Type: "string"},
		{Name: "address_line1", Location: "form", Type: "string"},
		{Name: "address_line2", Location: "form", Type: "string"},
		{Name: "address_state", Location: "form", Type: "string"},
		{Name: "address_zip", Location: "form", Type: "string"},
		{Name: "exp_month", Location: "form", Type: "string"},
		{Name: "exp_year", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "owner", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "email", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/bank_accounts/{id}/verify", OperationID: "PostCustomersCustomerBankAccountsIdVerify", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "amounts", Location: "form", Type: "array"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/cards", OperationID: "GetCustomersCustomerCards", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/cards", OperationID: "PostCustomersCustomerCards", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "alipay_account", Location: "form", Type: "string"},
		{Name: "bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_holder_name", Location: "form", Type: "string"},
			{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
			{Name: "account_number", Location: "form", Required: true, Type: "string"},
			{Name: "country", Location: "form", Required: true, Type: "string"},
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "object", Location: "form", Type: "string", Enum: []string{"bank_account"}},
			{Name: "routing_number", Location: "form", Type: "string"},
		}},
		{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address_city", Location: "form", Type: "string"},
			{Name: "address_country", Location: "form", Type: "string"},
			{Name: "address_line1", Location: "form", Type: "string"},
			{Name: "address_line2", Location: "form", Type: "string"},
			{Name: "address_state", Location: "form", Type: "string"},
			{Name: "address_zip", Location: "form", Type: "string"},
			{Name: "cvc", Location: "form", Type: "string"},
			{Name: "exp_month", Location: "form", Required: true, Type: "integer"},
			{Name: "exp_year", Location: "form", Required: true, Type: "integer"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "number", Location: "form", Required: true, Type: "string"},
			{Name: "object", Location: "form", Type: "string", Enum: []string{"card"}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "source", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/customers/{customer}/cards/{id}", OperationID: "DeleteCustomersCustomerCardsId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/cards/{id}", OperationID: "GetCustomersCustomerCardsId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/cards/{id}", OperationID: "PostCustomersCustomerCardsId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "account_holder_name", Location: "form", Type: "string"},
		{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
		{Name: "address_city", Location: "form", Type: "string"},
		{Name: "address_country", Location: "form", Type: "string"},
		{Name: "address_line1", Location: "form", Type: "string"},
		{Name: "address_line2", Location: "form", Type: "string"},
		{Name: "address_state", Location: "form", Type: "string"},
		{Name: "address_zip", Location: "form", Type: "string"},
		{Name: "exp_month", Location: "form", Type: "string"},
		{Name: "exp_year", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "owner", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "email", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/cash_balance", OperationID: "GetCustomersCustomerCashBalance", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/cash_balance", OperationID: "PostCustomersCustomerCashBalance", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "reconciliation_mode", Location: "form", Type: "string", Enum: []string{"automatic", "manual", "merchant_default"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/cash_balance_transactions", OperationID: "GetCustomersCustomerCashBalanceTransactions", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/cash_balance_transactions/{transaction}", OperationID: "GetCustomersCustomerCashBalanceTransactionsTransaction", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "transaction", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/customers/{customer}/discount", OperationID: "DeleteCustomersCustomerDiscount", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/discount", OperationID: "GetCustomersCustomerDiscount", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/funding_instructions", OperationID: "PostCustomersCustomerFundingInstructions", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "bank_transfer", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "eu_bank_transfer", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "country", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "requested_address_types", Location: "form", Type: "array"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"eu_bank_transfer", "gb_bank_transfer", "jp_bank_transfer", "mx_bank_transfer", "us_bank_transfer"}},
		}},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "funding_type", Location: "form", Required: true, Type: "string", Enum: []string{"bank_transfer"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/payment_methods", OperationID: "GetCustomersCustomerPaymentMethods", Params: []ParameterValidation{
		{Name: "allow_redisplay", Location: "query", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "type", Location: "query", Type: "string", Enum: []string{"acss_debit", "affirm", "afterpay_clearpay", "alipay", "alma", "amazon_pay", "au_becs_debit", "bacs_debit", "bancontact", "billie", "blik", "boleto", "card", "cashapp", "crypto", "custom", "customer_balance", "eps", "fpx", "giropay", "grabpay", "ideal", "kakao_pay", "klarna", "konbini", "kr_card", "link", "mb_way", "mobilepay", "multibanco", "naver_pay", "nz_bank_account", "oxxo", "p24", "pay_by_bank", "payco", "paynow", "paypal", "payto", "pix", "promptpay", "revolut_pay", "samsung_pay", "satispay", "sepa_debit", "sofort", "sunbit", "swish", "twint", "upi", "us_bank_account", "wechat_pay", "zip"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/payment_methods/{payment_method}", OperationID: "GetCustomersCustomerPaymentMethodsPaymentMethod", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "payment_method", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/sources", OperationID: "GetCustomersCustomerSources", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "object", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/sources", OperationID: "PostCustomersCustomerSources", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "alipay_account", Location: "form", Type: "string"},
		{Name: "bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_holder_name", Location: "form", Type: "string"},
			{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
			{Name: "account_number", Location: "form", Required: true, Type: "string"},
			{Name: "country", Location: "form", Required: true, Type: "string"},
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "object", Location: "form", Type: "string", Enum: []string{"bank_account"}},
			{Name: "routing_number", Location: "form", Type: "string"},
		}},
		{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address_city", Location: "form", Type: "string"},
			{Name: "address_country", Location: "form", Type: "string"},
			{Name: "address_line1", Location: "form", Type: "string"},
			{Name: "address_line2", Location: "form", Type: "string"},
			{Name: "address_state", Location: "form", Type: "string"},
			{Name: "address_zip", Location: "form", Type: "string"},
			{Name: "cvc", Location: "form", Type: "string"},
			{Name: "exp_month", Location: "form", Required: true, Type: "integer"},
			{Name: "exp_year", Location: "form", Required: true, Type: "integer"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "number", Location: "form", Required: true, Type: "string"},
			{Name: "object", Location: "form", Type: "string", Enum: []string{"card"}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "source", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/customers/{customer}/sources/{id}", OperationID: "DeleteCustomersCustomerSourcesId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/sources/{id}", OperationID: "GetCustomersCustomerSourcesId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/sources/{id}", OperationID: "PostCustomersCustomerSourcesId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "account_holder_name", Location: "form", Type: "string"},
		{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
		{Name: "address_city", Location: "form", Type: "string"},
		{Name: "address_country", Location: "form", Type: "string"},
		{Name: "address_line1", Location: "form", Type: "string"},
		{Name: "address_line2", Location: "form", Type: "string"},
		{Name: "address_state", Location: "form", Type: "string"},
		{Name: "address_zip", Location: "form", Type: "string"},
		{Name: "exp_month", Location: "form", Type: "string"},
		{Name: "exp_year", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "owner", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "email", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/sources/{id}/verify", OperationID: "PostCustomersCustomerSourcesIdVerify", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "amounts", Location: "form", Type: "array"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/subscriptions", OperationID: "GetCustomersCustomerSubscriptions", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/subscriptions", OperationID: "PostCustomersCustomerSubscriptions", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "add_invoice_items", Location: "form", Type: "array"},
		{Name: "application_fee_percent", Location: "form", Enum: []string{""}},
		{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "liability", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "backdate_start_date", Location: "form", Type: "integer"},
		{Name: "billing_cycle_anchor", Location: "form", Type: "integer"},
		{Name: "billing_thresholds", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "amount_gte", Location: "form", Type: "integer"},
			{Name: "reset_billing_cycle_anchor", Location: "form", Type: "boolean"},
		}},
		{Name: "cancel_at", Location: "form", Enum: []string{"max_period_end", "min_period_end"}},
		{Name: "cancel_at_period_end", Location: "form", Type: "boolean"},
		{Name: "collection_method", Location: "form", Type: "string", Enum: []string{"charge_automatically", "send_invoice"}},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "days_until_due", Location: "form", Type: "integer"},
		{Name: "default_payment_method", Location: "form", Type: "string"},
		{Name: "default_source", Location: "form", Type: "string"},
		{Name: "default_tax_rates", Location: "form", Enum: []string{""}},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "invoice_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_tax_ids", Location: "form", Enum: []string{""}},
			{Name: "issuer", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "items", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "off_session", Location: "form", Type: "boolean"},
		{Name: "payment_behavior", Location: "form", Type: "string", Enum: []string{"allow_incomplete", "default_incomplete", "error_if_incomplete", "pending_if_incomplete"}},
		{Name: "payment_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "acss_debit", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "bancontact", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "card", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "customer_balance", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "konbini", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "payto", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "pix", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "sepa_debit", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "upi", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "us_bank_account", Location: "form", Type: "object", Enum: []string{""}},
			}},
			{Name: "payment_method_types", Location: "form", Enum: []string{""}},
			{Name: "save_default_payment_method", Location: "form", Type: "string", Enum: []string{"off", "on_subscription"}},
		}},
		{Name: "pending_invoice_item_interval", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "interval", Location: "form", Required: true, Type: "string", Enum: []string{"day", "month", "week", "year"}},
			{Name: "interval_count", Location: "form", Type: "integer"},
		}},
		{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
		{Name: "transfer_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount_percent", Location: "form", Type: "number"},
			{Name: "destination", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "trial_end", Location: "form", Enum: []string{"now"}},
		{Name: "trial_from_plan", Location: "form", Type: "boolean"},
		{Name: "trial_period_days", Location: "form", Type: "integer"},
		{Name: "trial_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "end_behavior", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "missing_payment_method", Location: "form", Required: true, Type: "string", Enum: []string{"cancel", "create_invoice", "pause"}},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/customers/{customer}/subscriptions/{subscription_exposed_id}", OperationID: "DeleteCustomersCustomerSubscriptionsSubscriptionExposedId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "subscription_exposed_id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "invoice_now", Location: "form", Type: "boolean"},
		{Name: "prorate", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/subscriptions/{subscription_exposed_id}", OperationID: "GetCustomersCustomerSubscriptionsSubscriptionExposedId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "subscription_exposed_id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/subscriptions/{subscription_exposed_id}", OperationID: "PostCustomersCustomerSubscriptionsSubscriptionExposedId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "subscription_exposed_id", Location: "path", Required: true, Type: "string"},
		{Name: "add_invoice_items", Location: "form", Type: "array"},
		{Name: "application_fee_percent", Location: "form", Enum: []string{""}},
		{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "liability", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "billing_cycle_anchor", Location: "form", Type: "string", Enum: []string{"now", "unchanged"}},
		{Name: "billing_thresholds", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "amount_gte", Location: "form", Type: "integer"},
			{Name: "reset_billing_cycle_anchor", Location: "form", Type: "boolean"},
		}},
		{Name: "cancel_at", Location: "form", Enum: []string{"", "max_period_end", "min_period_end"}},
		{Name: "cancel_at_period_end", Location: "form", Type: "boolean"},
		{Name: "cancellation_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "comment", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "feedback", Location: "form", Type: "string", Enum: []string{"", "customer_service", "low_quality", "missing_features", "other", "switched_service", "too_complex", "too_expensive", "unused"}},
		}},
		{Name: "collection_method", Location: "form", Type: "string", Enum: []string{"charge_automatically", "send_invoice"}},
		{Name: "days_until_due", Location: "form", Type: "integer"},
		{Name: "default_payment_method", Location: "form", Type: "string"},
		{Name: "default_source", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "default_tax_rates", Location: "form", Enum: []string{""}},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "invoice_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_tax_ids", Location: "form", Enum: []string{""}},
			{Name: "issuer", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "items", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "off_session", Location: "form", Type: "boolean"},
		{Name: "pause_collection", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "behavior", Location: "form", Required: true, Type: "string", Enum: []string{"keep_as_draft", "mark_uncollectible", "void"}},
			{Name: "resumes_at", Location: "form", Type: "integer"},
		}},
		{Name: "payment_behavior", Location: "form", Type: "string", Enum: []string{"allow_incomplete", "default_incomplete", "error_if_incomplete", "pending_if_incomplete"}},
		{Name: "payment_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "acss_debit", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "bancontact", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "card", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "customer_balance", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "konbini", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "payto", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "pix", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "sepa_debit", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "upi", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "us_bank_account", Location: "form", Type: "object", Enum: []string{""}},
			}},
			{Name: "payment_method_types", Location: "form", Enum: []string{""}},
			{Name: "save_default_payment_method", Location: "form", Type: "string", Enum: []string{"off", "on_subscription"}},
		}},
		{Name: "pending_invoice_item_interval", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "interval", Location: "form", Required: true, Type: "string", Enum: []string{"day", "month", "week", "year"}},
			{Name: "interval_count", Location: "form", Type: "integer"},
		}},
		{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
		{Name: "proration_date", Location: "form", Type: "integer"},
		{Name: "transfer_data", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "amount_percent", Location: "form", Type: "number"},
			{Name: "destination", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "trial_end", Location: "form", Enum: []string{"now"}},
		{Name: "trial_from_plan", Location: "form", Type: "boolean"},
		{Name: "trial_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "end_behavior", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "missing_payment_method", Location: "form", Required: true, Type: "string", Enum: []string{"cancel", "create_invoice", "pause"}},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/customers/{customer}/subscriptions/{subscription_exposed_id}/discount", OperationID: "DeleteCustomersCustomerSubscriptionsSubscriptionExposedIdDiscount", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "subscription_exposed_id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/subscriptions/{subscription_exposed_id}/discount", OperationID: "GetCustomersCustomerSubscriptionsSubscriptionExposedIdDiscount", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "subscription_exposed_id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/tax_ids", OperationID: "GetCustomersCustomerTaxIds", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/customers/{customer}/tax_ids", OperationID: "PostCustomersCustomerTaxIds", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "type", Location: "form", Required: true, Type: "string"},
		{Name: "value", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/customers/{customer}/tax_ids/{id}", OperationID: "DeleteCustomersCustomerTaxIdsId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/customers/{customer}/tax_ids/{id}", OperationID: "GetCustomersCustomerTaxIdsId", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/disputes", OperationID: "GetDisputes", Params: []ParameterValidation{
		{Name: "charge", Location: "query", Type: "string"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "payment_intent", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/disputes/{dispute}", OperationID: "GetDisputesDispute", Params: []ParameterValidation{
		{Name: "dispute", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/disputes/{dispute}", OperationID: "PostDisputesDispute", Params: []ParameterValidation{
		{Name: "dispute", Location: "path", Required: true, Type: "string"},
		{Name: "evidence", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "access_activity_log", Location: "form", Type: "string"},
			{Name: "billing_address", Location: "form", Type: "string"},
			{Name: "cancellation_policy", Location: "form", Type: "string"},
			{Name: "cancellation_policy_disclosure", Location: "form", Type: "string"},
			{Name: "cancellation_rebuttal", Location: "form", Type: "string"},
			{Name: "customer_communication", Location: "form", Type: "string"},
			{Name: "customer_email_address", Location: "form", Type: "string"},
			{Name: "customer_name", Location: "form", Type: "string"},
			{Name: "customer_purchase_ip", Location: "form", Type: "string"},
			{Name: "customer_signature", Location: "form", Type: "string"},
			{Name: "duplicate_charge_documentation", Location: "form", Type: "string"},
			{Name: "duplicate_charge_explanation", Location: "form", Type: "string"},
			{Name: "duplicate_charge_id", Location: "form", Type: "string"},
			{Name: "enhanced_evidence", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "visa_compelling_evidence_3", Location: "form", Type: "object"},
				{Name: "visa_compliance", Location: "form", Type: "object"},
			}},
			{Name: "product_description", Location: "form", Type: "string"},
			{Name: "receipt", Location: "form", Type: "string"},
			{Name: "refund_policy", Location: "form", Type: "string"},
			{Name: "refund_policy_disclosure", Location: "form", Type: "string"},
			{Name: "refund_refusal_explanation", Location: "form", Type: "string"},
			{Name: "service_date", Location: "form", Type: "string"},
			{Name: "service_documentation", Location: "form", Type: "string"},
			{Name: "shipping_address", Location: "form", Type: "string"},
			{Name: "shipping_carrier", Location: "form", Type: "string"},
			{Name: "shipping_date", Location: "form", Type: "string"},
			{Name: "shipping_documentation", Location: "form", Type: "string"},
			{Name: "shipping_tracking_number", Location: "form", Type: "string"},
			{Name: "uncategorized_file", Location: "form", Type: "string"},
			{Name: "uncategorized_text", Location: "form", Type: "string"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "submit", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/disputes/{dispute}/close", OperationID: "PostDisputesDisputeClose", Params: []ParameterValidation{
		{Name: "dispute", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/entitlements/active_entitlements", OperationID: "GetEntitlementsActiveEntitlements", Params: []ParameterValidation{
		{Name: "customer", Location: "query", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/entitlements/active_entitlements/{id}", OperationID: "GetEntitlementsActiveEntitlementsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/entitlements/features", OperationID: "GetEntitlementsFeatures", Params: []ParameterValidation{
		{Name: "archived", Location: "query", Type: "boolean"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "lookup_key", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/entitlements/features", OperationID: "PostEntitlementsFeatures", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "lookup_key", Location: "form", Required: true, Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "name", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/entitlements/features/{id}", OperationID: "GetEntitlementsFeaturesId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/entitlements/features/{id}", OperationID: "PostEntitlementsFeaturesId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "name", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/ephemeral_keys", OperationID: "PostEphemeralKeys", Params: []ParameterValidation{
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "issuing_card", Location: "form", Type: "string"},
		{Name: "nonce", Location: "form", Type: "string"},
		{Name: "verification_session", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/ephemeral_keys/{key}", OperationID: "DeleteEphemeralKeysKey", Params: []ParameterValidation{
		{Name: "key", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/events", OperationID: "GetEvents", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "delivery_success", Location: "query", Type: "boolean"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "type", Location: "query", Type: "string"},
		{Name: "types", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/events/{id}", OperationID: "GetEventsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/exchange_rates", OperationID: "GetExchangeRates", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/exchange_rates/{rate_id}", OperationID: "GetExchangeRatesRateId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "rate_id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/external_accounts/{id}", OperationID: "PostExternalAccountsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "account_holder_name", Location: "form", Type: "string"},
		{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"", "company", "individual"}},
		{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "futsu", "savings", "toza"}},
		{Name: "address_city", Location: "form", Type: "string"},
		{Name: "address_country", Location: "form", Type: "string"},
		{Name: "address_line1", Location: "form", Type: "string"},
		{Name: "address_line2", Location: "form", Type: "string"},
		{Name: "address_state", Location: "form", Type: "string"},
		{Name: "address_zip", Location: "form", Type: "string"},
		{Name: "default_for_currency", Location: "form", Type: "boolean"},
		{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "bank_account_ownership_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Type: "array"},
			}},
		}},
		{Name: "exp_month", Location: "form", Type: "string"},
		{Name: "exp_year", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "name", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/file_links", OperationID: "GetFileLinks", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "expired", Location: "query", Type: "boolean"},
		{Name: "file", Location: "query", Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/file_links", OperationID: "PostFileLinks", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "expires_at", Location: "form", Type: "integer"},
		{Name: "file", Location: "form", Required: true, Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/file_links/{link}", OperationID: "GetFileLinksLink", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "link", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/file_links/{link}", OperationID: "PostFileLinksLink", Params: []ParameterValidation{
		{Name: "link", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "expires_at", Location: "form", Enum: []string{"", "now"}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/files", OperationID: "GetFiles", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "purpose", Location: "query", Type: "string", Enum: []string{"account_requirement", "additional_verification", "business_icon", "business_logo", "customer_signature", "dispute_evidence", "document_provider_identity_document", "finance_report_run", "financial_account_statement", "identity_document", "identity_document_downloadable", "issuing_regulatory_reporting", "pci_document", "platform_terms_of_service", "selfie", "sigma_scheduled_query", "tax_document_user_upload", "terminal_android_apk", "terminal_reader_splashscreen", "terminal_wifi_certificate", "terminal_wifi_private_key"}},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/files", OperationID: "PostFiles", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "file", Location: "form", Required: true, Type: "string"},
		{Name: "file_link_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "create", Location: "form", Required: true, Type: "boolean"},
			{Name: "expires_at", Location: "form", Type: "integer"},
			{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		}},
		{Name: "purpose", Location: "form", Required: true, Type: "string", Enum: []string{"account_requirement", "additional_verification", "business_icon", "business_logo", "customer_signature", "dispute_evidence", "identity_document", "issuing_regulatory_reporting", "pci_document", "platform_terms_of_service", "tax_document_user_upload", "terminal_android_apk", "terminal_reader_splashscreen", "terminal_wifi_certificate", "terminal_wifi_private_key"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/files/{file}", OperationID: "GetFilesFile", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "file", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/financial_connections/accounts", OperationID: "GetFinancialConnectionsAccounts", Params: []ParameterValidation{
		{Name: "account_holder", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "query", Type: "string"},
			{Name: "customer", Location: "query", Type: "string"},
			{Name: "customer_account", Location: "query", Type: "string"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "session", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/financial_connections/accounts/{account}", OperationID: "GetFinancialConnectionsAccountsAccount", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/financial_connections/accounts/{account}/disconnect", OperationID: "PostFinancialConnectionsAccountsAccountDisconnect", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/financial_connections/accounts/{account}/owners", OperationID: "GetFinancialConnectionsAccountsAccountOwners", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "ownership", Location: "query", Required: true, Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/financial_connections/accounts/{account}/refresh", OperationID: "PostFinancialConnectionsAccountsAccountRefresh", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "features", Location: "form", Required: true, Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/financial_connections/accounts/{account}/subscribe", OperationID: "PostFinancialConnectionsAccountsAccountSubscribe", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "features", Location: "form", Required: true, Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/financial_connections/accounts/{account}/unsubscribe", OperationID: "PostFinancialConnectionsAccountsAccountUnsubscribe", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "features", Location: "form", Required: true, Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/financial_connections/sessions", OperationID: "PostFinancialConnectionsSessions", Params: []ParameterValidation{
		{Name: "account_holder", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Type: "string"},
			{Name: "customer", Location: "form", Type: "string"},
			{Name: "customer_account", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "customer"}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "filters", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_subcategories", Location: "form", Type: "array"},
			{Name: "countries", Location: "form", Type: "array"},
		}},
		{Name: "permissions", Location: "form", Required: true, Type: "array"},
		{Name: "prefetch", Location: "form", Type: "array"},
		{Name: "return_url", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/financial_connections/sessions/{session}", OperationID: "GetFinancialConnectionsSessionsSession", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "session", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/financial_connections/transactions", OperationID: "GetFinancialConnectionsTransactions", Params: []ParameterValidation{
		{Name: "account", Location: "query", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "transacted_at", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "transaction_refresh", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "after", Location: "query", Required: true, Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/financial_connections/transactions/{transaction}", OperationID: "GetFinancialConnectionsTransactionsTransaction", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "transaction", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/forwarding/requests", OperationID: "GetForwardingRequests", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/forwarding/requests", OperationID: "PostForwardingRequests", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "payment_method", Location: "form", Required: true, Type: "string"},
		{Name: "replacements", Location: "form", Required: true, Type: "array"},
		{Name: "request", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "body", Location: "form", Type: "string"},
			{Name: "headers", Location: "form", Type: "array"},
		}},
		{Name: "url", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/forwarding/requests/{id}", OperationID: "GetForwardingRequestsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/identity/verification_reports", OperationID: "GetIdentityVerificationReports", Params: []ParameterValidation{
		{Name: "client_reference_id", Location: "query", Type: "string"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "type", Location: "query", Type: "string", Enum: []string{"document", "id_number"}},
		{Name: "verification_session", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/identity/verification_reports/{report}", OperationID: "GetIdentityVerificationReportsReport", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "report", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/identity/verification_sessions", OperationID: "GetIdentityVerificationSessions", Params: []ParameterValidation{
		{Name: "client_reference_id", Location: "query", Type: "string"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "related_customer", Location: "query", Type: "string"},
		{Name: "related_customer_account", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"canceled", "processing", "requires_input", "verified"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/identity/verification_sessions", OperationID: "PostIdentityVerificationSessions", Params: []ParameterValidation{
		{Name: "client_reference_id", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "document", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "allowed_types", Location: "form", Type: "array"},
				{Name: "require_id_number", Location: "form", Type: "boolean"},
				{Name: "require_live_capture", Location: "form", Type: "boolean"},
				{Name: "require_matching_selfie", Location: "form", Type: "boolean"},
			}},
		}},
		{Name: "provided_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "email", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
		{Name: "related_customer", Location: "form", Type: "string"},
		{Name: "related_customer_account", Location: "form", Type: "string"},
		{Name: "related_person", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Required: true, Type: "string"},
			{Name: "person", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "return_url", Location: "form", Type: "string"},
		{Name: "type", Location: "form", Type: "string", Enum: []string{"document", "id_number"}},
		{Name: "verification_flow", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/identity/verification_sessions/{session}", OperationID: "GetIdentityVerificationSessionsSession", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "session", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/identity/verification_sessions/{session}", OperationID: "PostIdentityVerificationSessionsSession", Params: []ParameterValidation{
		{Name: "session", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "document", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "allowed_types", Location: "form", Type: "array"},
				{Name: "require_id_number", Location: "form", Type: "boolean"},
				{Name: "require_live_capture", Location: "form", Type: "boolean"},
				{Name: "require_matching_selfie", Location: "form", Type: "boolean"},
			}},
		}},
		{Name: "provided_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "email", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
		{Name: "type", Location: "form", Type: "string", Enum: []string{"document", "id_number"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/identity/verification_sessions/{session}/cancel", OperationID: "PostIdentityVerificationSessionsSessionCancel", Params: []ParameterValidation{
		{Name: "session", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/identity/verification_sessions/{session}/redact", OperationID: "PostIdentityVerificationSessionsSessionRedact", Params: []ParameterValidation{
		{Name: "session", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/invoice_payments", OperationID: "GetInvoicePayments", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "invoice", Location: "query", Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "payment", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "payment_intent", Location: "query", Type: "string"},
			{Name: "payment_record", Location: "query", Type: "string"},
			{Name: "type", Location: "query", Required: true, Type: "string", Enum: []string{"payment_intent", "payment_record"}},
		}},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"canceled", "open", "paid"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/invoice_payments/{invoice_payment}", OperationID: "GetInvoicePaymentsInvoicePayment", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "invoice_payment", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/invoice_rendering_templates", OperationID: "GetInvoiceRenderingTemplates", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"active", "archived"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/invoice_rendering_templates/{template}", OperationID: "GetInvoiceRenderingTemplatesTemplate", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "template", Location: "path", Required: true, Type: "string"},
		{Name: "version", Location: "query", Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoice_rendering_templates/{template}/archive", OperationID: "PostInvoiceRenderingTemplatesTemplateArchive", Params: []ParameterValidation{
		{Name: "template", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoice_rendering_templates/{template}/unarchive", OperationID: "PostInvoiceRenderingTemplatesTemplateUnarchive", Params: []ParameterValidation{
		{Name: "template", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/invoiceitems", OperationID: "GetInvoiceitems", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "invoice", Location: "query", Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "pending", Location: "query", Type: "boolean"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoiceitems", OperationID: "PostInvoiceitems", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "discountable", Location: "form", Type: "boolean"},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "invoice", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "period", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "end", Location: "form", Required: true, Type: "integer"},
			{Name: "start", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "price_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "currency", Location: "form", Required: true, Type: "string"},
			{Name: "product", Location: "form", Required: true, Type: "string"},
			{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
			{Name: "unit_amount", Location: "form", Type: "integer"},
			{Name: "unit_amount_decimal", Location: "form", Type: "string"},
		}},
		{Name: "pricing", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "price", Location: "form", Type: "string"},
		}},
		{Name: "quantity", Location: "form", Type: "integer"},
		{Name: "quantity_decimal", Location: "form", Type: "string"},
		{Name: "subscription", Location: "form", Type: "string"},
		{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
		{Name: "tax_code", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "tax_rates", Location: "form", Type: "array"},
		{Name: "unit_amount_decimal", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/invoiceitems/{invoiceitem}", OperationID: "DeleteInvoiceitemsInvoiceitem", Params: []ParameterValidation{
		{Name: "invoiceitem", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/invoiceitems/{invoiceitem}", OperationID: "GetInvoiceitemsInvoiceitem", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "invoiceitem", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoiceitems/{invoiceitem}", OperationID: "PostInvoiceitemsInvoiceitem", Params: []ParameterValidation{
		{Name: "invoiceitem", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "discountable", Location: "form", Type: "boolean"},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "period", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "end", Location: "form", Required: true, Type: "integer"},
			{Name: "start", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "price_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "currency", Location: "form", Required: true, Type: "string"},
			{Name: "product", Location: "form", Required: true, Type: "string"},
			{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
			{Name: "unit_amount", Location: "form", Type: "integer"},
			{Name: "unit_amount_decimal", Location: "form", Type: "string"},
		}},
		{Name: "pricing", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "price", Location: "form", Type: "string"},
		}},
		{Name: "quantity", Location: "form", Type: "integer"},
		{Name: "quantity_decimal", Location: "form", Type: "string"},
		{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
		{Name: "tax_code", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "tax_rates", Location: "form", Enum: []string{""}},
		{Name: "unit_amount_decimal", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/invoices", OperationID: "GetInvoices", Params: []ParameterValidation{
		{Name: "collection_method", Location: "query", Type: "string", Enum: []string{"charge_automatically", "send_invoice"}},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "due_date", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"draft", "open", "paid", "uncollectible", "void"}},
		{Name: "subscription", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoices", OperationID: "PostInvoices", Params: []ParameterValidation{
		{Name: "account_tax_ids", Location: "form", Enum: []string{""}},
		{Name: "application_fee_amount", Location: "form", Type: "integer"},
		{Name: "auto_advance", Location: "form", Type: "boolean"},
		{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "liability", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "automatically_finalizes_at", Location: "form", Type: "integer"},
		{Name: "collection_method", Location: "form", Type: "string", Enum: []string{"charge_automatically", "send_invoice"}},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "custom_fields", Location: "form", Enum: []string{""}},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "days_until_due", Location: "form", Type: "integer"},
		{Name: "default_payment_method", Location: "form", Type: "string"},
		{Name: "default_source", Location: "form", Type: "string"},
		{Name: "default_tax_rates", Location: "form", Type: "array"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "due_date", Location: "form", Type: "integer"},
		{Name: "effective_at", Location: "form", Type: "integer"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "footer", Location: "form", Type: "string"},
		{Name: "from_invoice", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "action", Location: "form", Required: true, Type: "string", Enum: []string{"revision"}},
			{Name: "invoice", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "issuer", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
		}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "number", Location: "form", Type: "string"},
		{Name: "on_behalf_of", Location: "form", Type: "string"},
		{Name: "payment_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "default_mandate", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "acss_debit", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "bancontact", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "card", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "customer_balance", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "konbini", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "payto", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "pix", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "sepa_debit", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "upi", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "us_bank_account", Location: "form", Type: "object", Enum: []string{""}},
			}},
			{Name: "payment_method_types", Location: "form", Enum: []string{""}},
		}},
		{Name: "pending_invoice_items_behavior", Location: "form", Type: "string", Enum: []string{"exclude", "include"}},
		{Name: "rendering", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount_tax_display", Location: "form", Type: "string", Enum: []string{"", "exclude_tax", "include_inclusive_tax"}},
			{Name: "pdf", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "page_size", Location: "form", Type: "string", Enum: []string{"a4", "auto", "letter"}},
			}},
			{Name: "template", Location: "form", Type: "string"},
			{Name: "template_version", Location: "form", Enum: []string{""}},
		}},
		{Name: "shipping_cost", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "shipping_rate", Location: "form", Type: "string"},
			{Name: "shipping_rate_data", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "delivery_estimate", Location: "form", Type: "object"},
				{Name: "display_name", Location: "form", Required: true, Type: "string"},
				{Name: "fixed_amount", Location: "form", Type: "object"},
				{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
				{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
				{Name: "tax_code", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Type: "string", Enum: []string{"fixed_amount"}},
			}},
		}},
		{Name: "shipping_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
		{Name: "subscription", Location: "form", Type: "string"},
		{Name: "transfer_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount", Location: "form", Type: "integer"},
			{Name: "destination", Location: "form", Required: true, Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoices/create_preview", OperationID: "PostInvoicesCreatePreview", Params: []ParameterValidation{
		{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "liability", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "customer_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "shipping", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "address", Location: "form", Required: true, Type: "object"},
				{Name: "name", Location: "form", Required: true, Type: "string"},
				{Name: "phone", Location: "form", Type: "string"},
			}},
			{Name: "tax", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ip_address", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "tax_exempt", Location: "form", Type: "string", Enum: []string{"", "exempt", "none", "reverse"}},
			{Name: "tax_ids", Location: "form", Type: "array"},
		}},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "invoice_items", Location: "form", Type: "array"},
		{Name: "issuer", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
		}},
		{Name: "on_behalf_of", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "preview_mode", Location: "form", Type: "string", Enum: []string{"next", "recurring"}},
		{Name: "schedule", Location: "form", Type: "string"},
		{Name: "schedule_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "billing_mode", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "flexible", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"classic", "flexible"}},
			}},
			{Name: "end_behavior", Location: "form", Type: "string", Enum: []string{"cancel", "release"}},
			{Name: "phases", Location: "form", Type: "array"},
			{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
		}},
		{Name: "subscription", Location: "form", Type: "string"},
		{Name: "subscription_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "billing_cycle_anchor", Location: "form", Enum: []string{"now", "unchanged"}},
			{Name: "billing_mode", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "flexible", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"classic", "flexible"}},
			}},
			{Name: "cancel_at", Location: "form", Enum: []string{"", "max_period_end", "min_period_end"}},
			{Name: "cancel_at_period_end", Location: "form", Type: "boolean"},
			{Name: "cancel_now", Location: "form", Type: "boolean"},
			{Name: "default_tax_rates", Location: "form", Enum: []string{""}},
			{Name: "items", Location: "form", Type: "array"},
			{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
			{Name: "proration_date", Location: "form", Type: "integer"},
			{Name: "resume_at", Location: "form", Type: "string", Enum: []string{"now"}},
			{Name: "start_date", Location: "form", Type: "integer"},
			{Name: "trial_end", Location: "form", Enum: []string{"now"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/invoices/search", OperationID: "GetInvoicesSearch", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "page", Location: "query", Type: "string"},
		{Name: "query", Location: "query", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/invoices/{invoice}", OperationID: "DeleteInvoicesInvoice", Params: []ParameterValidation{
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/invoices/{invoice}", OperationID: "GetInvoicesInvoice", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoices/{invoice}", OperationID: "PostInvoicesInvoice", Params: []ParameterValidation{
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
		{Name: "account_tax_ids", Location: "form", Enum: []string{""}},
		{Name: "application_fee_amount", Location: "form", Type: "integer"},
		{Name: "auto_advance", Location: "form", Type: "boolean"},
		{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "liability", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "automatically_finalizes_at", Location: "form", Type: "integer"},
		{Name: "collection_method", Location: "form", Type: "string", Enum: []string{"charge_automatically", "send_invoice"}},
		{Name: "custom_fields", Location: "form", Enum: []string{""}},
		{Name: "days_until_due", Location: "form", Type: "integer"},
		{Name: "default_payment_method", Location: "form", Type: "string"},
		{Name: "default_source", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "default_tax_rates", Location: "form", Enum: []string{""}},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "due_date", Location: "form", Type: "integer"},
		{Name: "effective_at", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "footer", Location: "form", Type: "string"},
		{Name: "issuer", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
		}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "number", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "on_behalf_of", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "payment_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "default_mandate", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "acss_debit", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "bancontact", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "card", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "customer_balance", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "konbini", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "payto", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "pix", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "sepa_debit", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "upi", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "us_bank_account", Location: "form", Type: "object", Enum: []string{""}},
			}},
			{Name: "payment_method_types", Location: "form", Enum: []string{""}},
		}},
		{Name: "rendering", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount_tax_display", Location: "form", Type: "string", Enum: []string{"", "exclude_tax", "include_inclusive_tax"}},
			{Name: "pdf", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "page_size", Location: "form", Type: "string", Enum: []string{"a4", "auto", "letter"}},
			}},
			{Name: "template", Location: "form", Type: "string"},
			{Name: "template_version", Location: "form", Enum: []string{""}},
		}},
		{Name: "shipping_cost", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "shipping_rate", Location: "form", Type: "string"},
			{Name: "shipping_rate_data", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "delivery_estimate", Location: "form", Type: "object"},
				{Name: "display_name", Location: "form", Required: true, Type: "string"},
				{Name: "fixed_amount", Location: "form", Type: "object"},
				{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
				{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
				{Name: "tax_code", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Type: "string", Enum: []string{"fixed_amount"}},
			}},
		}},
		{Name: "shipping_details", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
		{Name: "transfer_data", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "amount", Location: "form", Type: "integer"},
			{Name: "destination", Location: "form", Required: true, Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoices/{invoice}/add_lines", OperationID: "PostInvoicesInvoiceAddLines", Params: []ParameterValidation{
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "invoice_metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "lines", Location: "form", Required: true, Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoices/{invoice}/attach_payment", OperationID: "PostInvoicesInvoiceAttachPayment", Params: []ParameterValidation{
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "payment_intent", Location: "form", Type: "string"},
		{Name: "payment_record", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoices/{invoice}/finalize", OperationID: "PostInvoicesInvoiceFinalize", Params: []ParameterValidation{
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
		{Name: "auto_advance", Location: "form", Type: "boolean"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/invoices/{invoice}/lines", OperationID: "GetInvoicesInvoiceLines", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoices/{invoice}/lines/{line_item_id}", OperationID: "PostInvoicesInvoiceLinesLineItemId", Params: []ParameterValidation{
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
		{Name: "line_item_id", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "discountable", Location: "form", Type: "boolean"},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "period", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "end", Location: "form", Required: true, Type: "integer"},
			{Name: "start", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "price_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "currency", Location: "form", Required: true, Type: "string"},
			{Name: "product", Location: "form", Type: "string"},
			{Name: "product_data", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "description", Location: "form", Type: "string"},
				{Name: "images", Location: "form", Type: "array"},
				{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
				{Name: "name", Location: "form", Required: true, Type: "string"},
				{Name: "tax_code", Location: "form", Type: "string"},
				{Name: "unit_label", Location: "form", Type: "string"},
			}},
			{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
			{Name: "unit_amount", Location: "form", Type: "integer"},
			{Name: "unit_amount_decimal", Location: "form", Type: "string"},
		}},
		{Name: "pricing", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "price", Location: "form", Type: "string"},
		}},
		{Name: "quantity", Location: "form", Type: "integer"},
		{Name: "quantity_decimal", Location: "form", Type: "string"},
		{Name: "tax_amounts", Location: "form", Enum: []string{""}},
		{Name: "tax_rates", Location: "form", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoices/{invoice}/mark_uncollectible", OperationID: "PostInvoicesInvoiceMarkUncollectible", Params: []ParameterValidation{
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoices/{invoice}/pay", OperationID: "PostInvoicesInvoicePay", Params: []ParameterValidation{
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "forgive", Location: "form", Type: "boolean"},
		{Name: "mandate", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "off_session", Location: "form", Type: "boolean"},
		{Name: "paid_out_of_band", Location: "form", Type: "boolean"},
		{Name: "payment_method", Location: "form", Type: "string"},
		{Name: "source", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoices/{invoice}/remove_lines", OperationID: "PostInvoicesInvoiceRemoveLines", Params: []ParameterValidation{
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "invoice_metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "lines", Location: "form", Required: true, Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoices/{invoice}/send", OperationID: "PostInvoicesInvoiceSend", Params: []ParameterValidation{
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoices/{invoice}/update_lines", OperationID: "PostInvoicesInvoiceUpdateLines", Params: []ParameterValidation{
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "invoice_metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "lines", Location: "form", Required: true, Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/invoices/{invoice}/void", OperationID: "PostInvoicesInvoiceVoid", Params: []ParameterValidation{
		{Name: "invoice", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/authorizations", OperationID: "GetIssuingAuthorizations", Params: []ParameterValidation{
		{Name: "card", Location: "query", Type: "string"},
		{Name: "cardholder", Location: "query", Type: "string"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"closed", "expired", "pending", "reversed"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/authorizations/{authorization}", OperationID: "GetIssuingAuthorizationsAuthorization", Params: []ParameterValidation{
		{Name: "authorization", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/authorizations/{authorization}", OperationID: "PostIssuingAuthorizationsAuthorization", Params: []ParameterValidation{
		{Name: "authorization", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/authorizations/{authorization}/approve", OperationID: "PostIssuingAuthorizationsAuthorizationApprove", Params: []ParameterValidation{
		{Name: "authorization", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/authorizations/{authorization}/decline", OperationID: "PostIssuingAuthorizationsAuthorizationDecline", Params: []ParameterValidation{
		{Name: "authorization", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/cardholders", OperationID: "GetIssuingCardholders", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "email", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "phone_number", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"active", "blocked", "inactive"}},
		{Name: "type", Location: "query", Type: "string", Enum: []string{"company", "individual"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/cardholders", OperationID: "PostIssuingCardholders", Params: []ParameterValidation{
		{Name: "billing", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Required: true, Type: "string"},
				{Name: "country", Location: "form", Required: true, Type: "string"},
				{Name: "line1", Location: "form", Required: true, Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Required: true, Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
		}},
		{Name: "company", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "tax_id", Location: "form", Type: "string"},
		}},
		{Name: "email", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "individual", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "card_issuing", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "user_terms_acceptance", Location: "form", Type: "object"},
			}},
			{Name: "dob", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "day", Location: "form", Required: true, Type: "integer"},
				{Name: "month", Location: "form", Required: true, Type: "integer"},
				{Name: "year", Location: "form", Required: true, Type: "integer"},
			}},
			{Name: "first_name", Location: "form", Type: "string"},
			{Name: "last_name", Location: "form", Type: "string"},
			{Name: "verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "document", Location: "form", Type: "object"},
			}},
		}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "name", Location: "form", Required: true, Type: "string"},
		{Name: "phone_number", Location: "form", Type: "string"},
		{Name: "preferred_locales", Location: "form", Type: "array"},
		{Name: "spending_controls", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "allowed_card_presences", Location: "form", Type: "array"},
			{Name: "allowed_categories", Location: "form", Type: "array"},
			{Name: "allowed_merchant_countries", Location: "form", Type: "array"},
			{Name: "blocked_card_presences", Location: "form", Type: "array"},
			{Name: "blocked_categories", Location: "form", Type: "array"},
			{Name: "blocked_merchant_countries", Location: "form", Type: "array"},
			{Name: "spending_limits", Location: "form", Type: "array"},
			{Name: "spending_limits_currency", Location: "form", Type: "string"},
		}},
		{Name: "status", Location: "form", Type: "string", Enum: []string{"active", "inactive"}},
		{Name: "type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/cardholders/{cardholder}", OperationID: "GetIssuingCardholdersCardholder", Params: []ParameterValidation{
		{Name: "cardholder", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/cardholders/{cardholder}", OperationID: "PostIssuingCardholdersCardholder", Params: []ParameterValidation{
		{Name: "cardholder", Location: "path", Required: true, Type: "string"},
		{Name: "billing", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Required: true, Type: "string"},
				{Name: "country", Location: "form", Required: true, Type: "string"},
				{Name: "line1", Location: "form", Required: true, Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Required: true, Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
		}},
		{Name: "company", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "tax_id", Location: "form", Type: "string"},
		}},
		{Name: "email", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "individual", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "card_issuing", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "user_terms_acceptance", Location: "form", Type: "object"},
			}},
			{Name: "dob", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "day", Location: "form", Required: true, Type: "integer"},
				{Name: "month", Location: "form", Required: true, Type: "integer"},
				{Name: "year", Location: "form", Required: true, Type: "integer"},
			}},
			{Name: "first_name", Location: "form", Type: "string"},
			{Name: "last_name", Location: "form", Type: "string"},
			{Name: "verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "document", Location: "form", Type: "object"},
			}},
		}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "phone_number", Location: "form", Type: "string"},
		{Name: "preferred_locales", Location: "form", Type: "array"},
		{Name: "spending_controls", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "allowed_card_presences", Location: "form", Type: "array"},
			{Name: "allowed_categories", Location: "form", Type: "array"},
			{Name: "allowed_merchant_countries", Location: "form", Type: "array"},
			{Name: "blocked_card_presences", Location: "form", Type: "array"},
			{Name: "blocked_categories", Location: "form", Type: "array"},
			{Name: "blocked_merchant_countries", Location: "form", Type: "array"},
			{Name: "spending_limits", Location: "form", Type: "array"},
			{Name: "spending_limits_currency", Location: "form", Type: "string"},
		}},
		{Name: "status", Location: "form", Type: "string", Enum: []string{"active", "inactive"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/cards", OperationID: "GetIssuingCards", Params: []ParameterValidation{
		{Name: "cardholder", Location: "query", Type: "string"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "exp_month", Location: "query", Type: "integer"},
		{Name: "exp_year", Location: "query", Type: "integer"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "last4", Location: "query", Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "personalization_design", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"active", "canceled", "inactive"}},
		{Name: "type", Location: "query", Type: "string", Enum: []string{"physical", "virtual"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/cards", OperationID: "PostIssuingCards", Params: []ParameterValidation{
		{Name: "cardholder", Location: "form", Type: "string"},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "exp_month", Location: "form", Type: "integer"},
		{Name: "exp_year", Location: "form", Type: "integer"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "financial_account", Location: "form", Type: "string"},
		{Name: "lifecycle_controls", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "cancel_after", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "payment_count", Location: "form", Required: true, Type: "integer"},
			}},
		}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "personalization_design", Location: "form", Type: "string"},
		{Name: "pin", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "encrypted_number", Location: "form", Type: "string"},
		}},
		{Name: "replacement_for", Location: "form", Type: "string"},
		{Name: "replacement_reason", Location: "form", Type: "string", Enum: []string{"damaged", "expired", "lost", "stolen"}},
		{Name: "second_line", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "shipping", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Required: true, Type: "string"},
				{Name: "country", Location: "form", Required: true, Type: "string"},
				{Name: "line1", Location: "form", Required: true, Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Required: true, Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "address_validation", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mode", Location: "form", Required: true, Type: "string", Enum: []string{"disabled", "normalization_only", "validation_and_normalization"}},
			}},
			{Name: "customs", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "eori_number", Location: "form", Type: "string"},
			}},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "phone_number", Location: "form", Type: "string"},
			{Name: "require_signature", Location: "form", Type: "boolean"},
			{Name: "service", Location: "form", Type: "string", Enum: []string{"express", "priority", "standard"}},
			{Name: "type", Location: "form", Type: "string", Enum: []string{"bulk", "individual"}},
		}},
		{Name: "spending_controls", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "allowed_card_presences", Location: "form", Type: "array"},
			{Name: "allowed_categories", Location: "form", Type: "array"},
			{Name: "allowed_merchant_countries", Location: "form", Type: "array"},
			{Name: "blocked_card_presences", Location: "form", Type: "array"},
			{Name: "blocked_categories", Location: "form", Type: "array"},
			{Name: "blocked_merchant_countries", Location: "form", Type: "array"},
			{Name: "spending_limits", Location: "form", Type: "array"},
		}},
		{Name: "status", Location: "form", Type: "string", Enum: []string{"active", "inactive"}},
		{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"physical", "virtual"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/cards/{card}", OperationID: "GetIssuingCardsCard", Params: []ParameterValidation{
		{Name: "card", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/cards/{card}", OperationID: "PostIssuingCardsCard", Params: []ParameterValidation{
		{Name: "card", Location: "path", Required: true, Type: "string"},
		{Name: "cancellation_reason", Location: "form", Type: "string", Enum: []string{"lost", "stolen"}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "personalization_design", Location: "form", Type: "string"},
		{Name: "pin", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "encrypted_number", Location: "form", Type: "string"},
		}},
		{Name: "shipping", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Required: true, Type: "string"},
				{Name: "country", Location: "form", Required: true, Type: "string"},
				{Name: "line1", Location: "form", Required: true, Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Required: true, Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "address_validation", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mode", Location: "form", Required: true, Type: "string", Enum: []string{"disabled", "normalization_only", "validation_and_normalization"}},
			}},
			{Name: "customs", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "eori_number", Location: "form", Type: "string"},
			}},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "phone_number", Location: "form", Type: "string"},
			{Name: "require_signature", Location: "form", Type: "boolean"},
			{Name: "service", Location: "form", Type: "string", Enum: []string{"express", "priority", "standard"}},
			{Name: "type", Location: "form", Type: "string", Enum: []string{"bulk", "individual"}},
		}},
		{Name: "spending_controls", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "allowed_card_presences", Location: "form", Type: "array"},
			{Name: "allowed_categories", Location: "form", Type: "array"},
			{Name: "allowed_merchant_countries", Location: "form", Type: "array"},
			{Name: "blocked_card_presences", Location: "form", Type: "array"},
			{Name: "blocked_categories", Location: "form", Type: "array"},
			{Name: "blocked_merchant_countries", Location: "form", Type: "array"},
			{Name: "spending_limits", Location: "form", Type: "array"},
		}},
		{Name: "status", Location: "form", Type: "string", Enum: []string{"active", "canceled", "inactive"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/disputes", OperationID: "GetIssuingDisputes", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"expired", "lost", "submitted", "unsubmitted", "won"}},
		{Name: "transaction", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/disputes", OperationID: "PostIssuingDisputes", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "evidence", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "canceled", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "canceled_at", Location: "form", Enum: []string{""}},
				{Name: "cancellation_policy_provided", Location: "form", Enum: []string{""}},
				{Name: "cancellation_reason", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "expected_at", Location: "form", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "product_description", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "product_type", Location: "form", Type: "string", Enum: []string{"", "merchandise", "service"}},
				{Name: "return_status", Location: "form", Type: "string", Enum: []string{"", "merchant_rejected", "successful"}},
				{Name: "returned_at", Location: "form", Enum: []string{""}},
			}},
			{Name: "duplicate", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "card_statement", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "cash_receipt", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "check_image", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "original_transaction", Location: "form", Type: "string"},
			}},
			{Name: "fraudulent", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "merchandise_not_as_described", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "received_at", Location: "form", Enum: []string{""}},
				{Name: "return_description", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "return_status", Location: "form", Type: "string", Enum: []string{"", "merchant_rejected", "successful"}},
				{Name: "returned_at", Location: "form", Enum: []string{""}},
			}},
			{Name: "no_valid_authorization", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "not_received", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "expected_at", Location: "form", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "product_description", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "product_type", Location: "form", Type: "string", Enum: []string{"", "merchandise", "service"}},
			}},
			{Name: "other", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "product_description", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "product_type", Location: "form", Type: "string", Enum: []string{"", "merchandise", "service"}},
			}},
			{Name: "reason", Location: "form", Type: "string", Enum: []string{"canceled", "duplicate", "fraudulent", "merchandise_not_as_described", "no_valid_authorization", "not_received", "other", "service_not_as_described"}},
			{Name: "service_not_as_described", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "canceled_at", Location: "form", Enum: []string{""}},
				{Name: "cancellation_reason", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "received_at", Location: "form", Enum: []string{""}},
			}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "transaction", Location: "form", Type: "string"},
		{Name: "treasury", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "received_debit", Location: "form", Required: true, Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/disputes/{dispute}", OperationID: "GetIssuingDisputesDispute", Params: []ParameterValidation{
		{Name: "dispute", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/disputes/{dispute}", OperationID: "PostIssuingDisputesDispute", Params: []ParameterValidation{
		{Name: "dispute", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "evidence", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "canceled", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "canceled_at", Location: "form", Enum: []string{""}},
				{Name: "cancellation_policy_provided", Location: "form", Enum: []string{""}},
				{Name: "cancellation_reason", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "expected_at", Location: "form", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "product_description", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "product_type", Location: "form", Type: "string", Enum: []string{"", "merchandise", "service"}},
				{Name: "return_status", Location: "form", Type: "string", Enum: []string{"", "merchant_rejected", "successful"}},
				{Name: "returned_at", Location: "form", Enum: []string{""}},
			}},
			{Name: "duplicate", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "card_statement", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "cash_receipt", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "check_image", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "original_transaction", Location: "form", Type: "string"},
			}},
			{Name: "fraudulent", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "merchandise_not_as_described", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "received_at", Location: "form", Enum: []string{""}},
				{Name: "return_description", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "return_status", Location: "form", Type: "string", Enum: []string{"", "merchant_rejected", "successful"}},
				{Name: "returned_at", Location: "form", Enum: []string{""}},
			}},
			{Name: "no_valid_authorization", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "not_received", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "expected_at", Location: "form", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "product_description", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "product_type", Location: "form", Type: "string", Enum: []string{"", "merchandise", "service"}},
			}},
			{Name: "other", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "product_description", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "product_type", Location: "form", Type: "string", Enum: []string{"", "merchandise", "service"}},
			}},
			{Name: "reason", Location: "form", Type: "string", Enum: []string{"canceled", "duplicate", "fraudulent", "merchandise_not_as_described", "no_valid_authorization", "not_received", "other", "service_not_as_described"}},
			{Name: "service_not_as_described", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "additional_documentation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "canceled_at", Location: "form", Enum: []string{""}},
				{Name: "cancellation_reason", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "explanation", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "received_at", Location: "form", Enum: []string{""}},
			}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/disputes/{dispute}/submit", OperationID: "PostIssuingDisputesDisputeSubmit", Params: []ParameterValidation{
		{Name: "dispute", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/personalization_designs", OperationID: "GetIssuingPersonalizationDesigns", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "lookup_keys", Location: "query", Type: "array"},
		{Name: "preferences", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "is_default", Location: "query", Type: "boolean"},
			{Name: "is_platform_default", Location: "query", Type: "boolean"},
		}},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"active", "inactive", "rejected", "review"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/personalization_designs", OperationID: "PostIssuingPersonalizationDesigns", Params: []ParameterValidation{
		{Name: "card_logo", Location: "form", Type: "string"},
		{Name: "carrier_text", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "footer_body", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "footer_title", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "header_body", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "header_title", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "lookup_key", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "physical_bundle", Location: "form", Required: true, Type: "string"},
		{Name: "preferences", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "is_default", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "transfer_lookup_key", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/personalization_designs/{personalization_design}", OperationID: "GetIssuingPersonalizationDesignsPersonalizationDesign", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "personalization_design", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/personalization_designs/{personalization_design}", OperationID: "PostIssuingPersonalizationDesignsPersonalizationDesign", Params: []ParameterValidation{
		{Name: "personalization_design", Location: "path", Required: true, Type: "string"},
		{Name: "card_logo", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "carrier_text", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "footer_body", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "footer_title", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "header_body", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "header_title", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "lookup_key", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "name", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "physical_bundle", Location: "form", Type: "string"},
		{Name: "preferences", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "is_default", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "transfer_lookup_key", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/physical_bundles", OperationID: "GetIssuingPhysicalBundles", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"active", "inactive", "review"}},
		{Name: "type", Location: "query", Type: "string", Enum: []string{"custom", "standard"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/physical_bundles/{physical_bundle}", OperationID: "GetIssuingPhysicalBundlesPhysicalBundle", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "physical_bundle", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/settlements/{settlement}", OperationID: "GetIssuingSettlementsSettlement", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "settlement", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/settlements/{settlement}", OperationID: "PostIssuingSettlementsSettlement", Params: []ParameterValidation{
		{Name: "settlement", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/tokens", OperationID: "GetIssuingTokens", Params: []ParameterValidation{
		{Name: "card", Location: "query", Required: true, Type: "string"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"active", "deleted", "requested", "suspended"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/tokens/{token}", OperationID: "GetIssuingTokensToken", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "token", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/tokens/{token}", OperationID: "PostIssuingTokensToken", Params: []ParameterValidation{
		{Name: "token", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "status", Location: "form", Required: true, Type: "string", Enum: []string{"active", "deleted", "suspended"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/transactions", OperationID: "GetIssuingTransactions", Params: []ParameterValidation{
		{Name: "card", Location: "query", Type: "string"},
		{Name: "cardholder", Location: "query", Type: "string"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "type", Location: "query", Type: "string", Enum: []string{"capture", "refund"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/issuing/transactions/{transaction}", OperationID: "GetIssuingTransactionsTransaction", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "transaction", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/issuing/transactions/{transaction}", OperationID: "PostIssuingTransactionsTransaction", Params: []ParameterValidation{
		{Name: "transaction", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/link_account_sessions", OperationID: "PostLinkAccountSessions", Params: []ParameterValidation{
		{Name: "account_holder", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Type: "string"},
			{Name: "customer", Location: "form", Type: "string"},
			{Name: "customer_account", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "customer"}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "filters", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_subcategories", Location: "form", Type: "array"},
			{Name: "countries", Location: "form", Type: "array"},
		}},
		{Name: "permissions", Location: "form", Required: true, Type: "array"},
		{Name: "prefetch", Location: "form", Type: "array"},
		{Name: "return_url", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/link_account_sessions/{session}", OperationID: "GetLinkAccountSessionsSession", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "session", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/linked_accounts", OperationID: "GetLinkedAccounts", Params: []ParameterValidation{
		{Name: "account_holder", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "query", Type: "string"},
			{Name: "customer", Location: "query", Type: "string"},
			{Name: "customer_account", Location: "query", Type: "string"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "session", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/linked_accounts/{account}", OperationID: "GetLinkedAccountsAccount", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/linked_accounts/{account}/disconnect", OperationID: "PostLinkedAccountsAccountDisconnect", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/linked_accounts/{account}/owners", OperationID: "GetLinkedAccountsAccountOwners", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "ownership", Location: "query", Required: true, Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/linked_accounts/{account}/refresh", OperationID: "PostLinkedAccountsAccountRefresh", Params: []ParameterValidation{
		{Name: "account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "features", Location: "form", Required: true, Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/mandates/{mandate}", OperationID: "GetMandatesMandate", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "mandate", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_attempt_records", OperationID: "GetPaymentAttemptRecords", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "payment_record", Location: "query", Required: true, Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_attempt_records/{id}", OperationID: "GetPaymentAttemptRecordsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_intents", OperationID: "GetPaymentIntents", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_intents", OperationID: "PostPaymentIntents", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Required: true, Type: "integer"},
		{Name: "amount_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "discount_amount", Location: "form", Enum: []string{""}},
			{Name: "enforce_arithmetic_validation", Location: "form", Type: "boolean"},
			{Name: "line_items", Location: "form", Enum: []string{""}},
			{Name: "shipping", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount", Location: "form", Enum: []string{""}},
				{Name: "from_postal_code", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "to_postal_code", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "tax", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "total_tax_amount", Location: "form", Required: true, Type: "integer"},
			}},
		}},
		{Name: "application_fee_amount", Location: "form", Type: "integer"},
		{Name: "automatic_payment_methods", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "allow_redirects", Location: "form", Type: "string", Enum: []string{"always", "never"}},
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"automatic", "automatic_async", "manual"}},
		{Name: "confirm", Location: "form", Type: "boolean"},
		{Name: "confirmation_method", Location: "form", Type: "string", Enum: []string{"automatic", "manual"}},
		{Name: "confirmation_token", Location: "form", Type: "string"},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "error_on_requires_action", Location: "form", Type: "boolean"},
		{Name: "excluded_payment_method_types", Location: "form", Type: "array"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "hooks", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "inputs", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tax", Location: "form", Type: "object"},
			}},
		}},
		{Name: "mandate", Location: "form", Type: "string"},
		{Name: "mandate_data", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "customer_acceptance", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "accepted_at", Location: "form", Type: "integer"},
				{Name: "offline", Location: "form", Type: "object"},
				{Name: "online", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"offline", "online"}},
			}},
		}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "off_session", Location: "form", Enum: []string{"one_off", "recurring"}},
		{Name: "on_behalf_of", Location: "form", Type: "string"},
		{Name: "payment_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "customer_reference", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "order_reference", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "payment_method", Location: "form", Type: "string"},
		{Name: "payment_method_configuration", Location: "form", Type: "string"},
		{Name: "payment_method_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "institution_number", Location: "form", Required: true, Type: "string"},
				{Name: "transit_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "affirm", Location: "form", Type: "object"},
			{Name: "afterpay_clearpay", Location: "form", Type: "object"},
			{Name: "alipay", Location: "form", Type: "object"},
			{Name: "allow_redisplay", Location: "form", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
			{Name: "alma", Location: "form", Type: "object"},
			{Name: "amazon_pay", Location: "form", Type: "object"},
			{Name: "au_becs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bsb_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "sort_code", Location: "form", Type: "string"},
			}},
			{Name: "bancontact", Location: "form", Type: "object"},
			{Name: "billie", Location: "form", Type: "object"},
			{Name: "billing_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "email", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "name", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "tax_id", Location: "form", Type: "string"},
			}},
			{Name: "blik", Location: "form", Type: "object"},
			{Name: "boleto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tax_id", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "cashapp", Location: "form", Type: "object"},
			{Name: "crypto", Location: "form", Type: "object"},
			{Name: "customer_balance", Location: "form", Type: "object"},
			{Name: "eps", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"arzte_und_apotheker_bank", "austrian_anadi_bank_ag", "bank_austria", "bankhaus_carl_spangler", "bankhaus_schelhammer_und_schattera_ag", "bawag_psk_ag", "bks_bank_ag", "brull_kallmus_bank_ag", "btv_vier_lander_bank", "capital_bank_grawe_gruppe_ag", "deutsche_bank_ag", "dolomitenbank", "easybank_ag", "erste_bank_und_sparkassen", "hypo_alpeadriabank_international_ag", "hypo_bank_burgenland_aktiengesellschaft", "hypo_noe_lb_fur_niederosterreich_u_wien", "hypo_oberosterreich_salzburg_steiermark", "hypo_tirol_bank_ag", "hypo_vorarlberg_bank_ag", "marchfelder_bank", "oberbank_ag", "raiffeisen_bankengruppe_osterreich", "schoellerbank_ag", "sparda_bank_wien", "volksbank_gruppe", "volkskreditbank_ag", "vr_bank_braunau"}},
			}},
			{Name: "fpx", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Required: true, Type: "string", Enum: []string{"affin_bank", "agrobank", "alliance_bank", "ambank", "bank_islam", "bank_muamalat", "bank_of_china", "bank_rakyat", "bsn", "cimb", "deutsche_bank", "hong_leong_bank", "hsbc", "kfh", "maybank2e", "maybank2u", "ocbc", "pb_enterprise", "public_bank", "rhb", "standard_chartered", "uob"}},
			}},
			{Name: "giropay", Location: "form", Type: "object"},
			{Name: "grabpay", Location: "form", Type: "object"},
			{Name: "ideal", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"abn_amro", "adyen", "asn_bank", "bunq", "buut", "finom", "handelsbanken", "ing", "knab", "mollie", "moneyou", "n26", "nn", "rabobank", "regiobank", "revolut", "sns_bank", "triodos_bank", "van_lanschot", "yoursafe"}},
			}},
			{Name: "interac_present", Location: "form", Type: "object"},
			{Name: "kakao_pay", Location: "form", Type: "object"},
			{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "dob", Location: "form", Type: "object"},
			}},
			{Name: "konbini", Location: "form", Type: "object"},
			{Name: "kr_card", Location: "form", Type: "object"},
			{Name: "link", Location: "form", Type: "object"},
			{Name: "mb_way", Location: "form", Type: "object"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "mobilepay", Location: "form", Type: "object"},
			{Name: "multibanco", Location: "form", Type: "object"},
			{Name: "naver_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "funding", Location: "form", Type: "string", Enum: []string{"card", "points"}},
			}},
			{Name: "nz_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_name", Location: "form", Type: "string"},
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bank_code", Location: "form", Required: true, Type: "string"},
				{Name: "branch_code", Location: "form", Required: true, Type: "string"},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "suffix", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "oxxo", Location: "form", Type: "object"},
			{Name: "p24", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"alior_bank", "bank_millennium", "bank_nowy_bfg_sa", "bank_pekao_sa", "banki_spbdzielcze", "blik", "bnp_paribas", "boz", "citi_handlowy", "credit_agricole", "envelobank", "etransfer_pocztowy24", "getin_bank", "ideabank", "ing", "inteligo", "mbank_mtransfer", "nest_przelew", "noble_pay", "pbac_z_ipko", "plus_bank", "santander_przelew24", "tmobile_usbugi_bankowe", "toyota_bank", "velobank", "volkswagen_bank"}},
			}},
			{Name: "pay_by_bank", Location: "form", Type: "object"},
			{Name: "payco", Location: "form", Type: "object"},
			{Name: "paynow", Location: "form", Type: "object"},
			{Name: "paypal", Location: "form", Type: "object"},
			{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "bsb_number", Location: "form", Type: "string"},
				{Name: "pay_id", Location: "form", Type: "string"},
			}},
			{Name: "pix", Location: "form", Type: "object"},
			{Name: "promptpay", Location: "form", Type: "object"},
			{Name: "radar_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "session", Location: "form", Type: "string"},
			}},
			{Name: "revolut_pay", Location: "form", Type: "object"},
			{Name: "samsung_pay", Location: "form", Type: "object"},
			{Name: "satispay", Location: "form", Type: "object"},
			{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "iban", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "sofort", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "country", Location: "form", Required: true, Type: "string", Enum: []string{"AT", "BE", "DE", "ES", "IT", "NL"}},
			}},
			{Name: "sunbit", Location: "form", Type: "object"},
			{Name: "swish", Location: "form", Type: "object"},
			{Name: "twint", Location: "form", Type: "object"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"acss_debit", "affirm", "afterpay_clearpay", "alipay", "alma", "amazon_pay", "au_becs_debit", "bacs_debit", "bancontact", "billie", "blik", "boleto", "cashapp", "crypto", "customer_balance", "eps", "fpx", "giropay", "grabpay", "ideal", "kakao_pay", "klarna", "konbini", "kr_card", "link", "mb_way", "mobilepay", "multibanco", "naver_pay", "nz_bank_account", "oxxo", "p24", "pay_by_bank", "payco", "paynow", "paypal", "payto", "pix", "promptpay", "revolut_pay", "samsung_pay", "satispay", "sepa_debit", "sofort", "sunbit", "swish", "twint", "upi", "us_bank_account", "wechat_pay", "zip"}},
			{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "savings"}},
				{Name: "financial_connections_account", Location: "form", Type: "string"},
				{Name: "routing_number", Location: "form", Type: "string"},
			}},
			{Name: "wechat_pay", Location: "form", Type: "object"},
			{Name: "zip", Location: "form", Type: "object"},
		}},
		{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant", "microdeposits"}},
			}},
			{Name: "affirm", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "preferred_locale", Location: "form", Type: "string"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "afterpay_clearpay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "alipay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "alma", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "amazon_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "au_becs_debit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "bacs_debit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "bancontact", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "preferred_language", Location: "form", Type: "string", Enum: []string{"de", "en", "fr", "nl"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "billie", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "blik", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "code", Location: "form", Type: "string"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none"}},
			}},
			{Name: "boleto", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "expires_after_days", Location: "form", Type: "integer"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
			}},
			{Name: "card", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "cvc_token", Location: "form", Type: "string"},
				{Name: "installments", Location: "form", Type: "object"},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "network", Location: "form", Type: "string", Enum: []string{"amex", "cartes_bancaires", "diners", "discover", "eftpos_au", "girocard", "interac", "jcb", "link", "mastercard", "unionpay", "unknown", "visa"}},
				{Name: "request_extended_authorization", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_incremental_authorization", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_multicapture", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_overcapture", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_three_d_secure", Location: "form", Type: "string", Enum: []string{"any", "automatic", "challenge"}},
				{Name: "require_cvc_recollection", Location: "form", Type: "boolean"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "statement_descriptor_suffix_kana", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "statement_descriptor_suffix_kanji", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "three_d_secure", Location: "form", Type: "object"},
			}},
			{Name: "card_present", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual", "manual_preferred"}},
				{Name: "request_extended_authorization", Location: "form", Type: "boolean"},
				{Name: "request_incremental_authorization_support", Location: "form", Type: "boolean"},
				{Name: "routing", Location: "form", Type: "object"},
			}},
			{Name: "cashapp", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
			}},
			{Name: "crypto", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "customer_balance", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "bank_transfer", Location: "form", Type: "object"},
				{Name: "funding_type", Location: "form", Type: "string", Enum: []string{"bank_transfer"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "eps", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "fpx", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "giropay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "grabpay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "ideal", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "interac_present", Location: "form", Type: "object", Enum: []string{""}},
			{Name: "kakao_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "klarna", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "on_demand", Location: "form", Type: "object"},
				{Name: "preferred_locale", Location: "form", Type: "string", Enum: []string{"cs-CZ", "da-DK", "de-AT", "de-CH", "de-DE", "el-GR", "en-AT", "en-AU", "en-BE", "en-CA", "en-CH", "en-CZ", "en-DE", "en-DK", "en-ES", "en-FI", "en-FR", "en-GB", "en-GR", "en-IE", "en-IT", "en-NL", "en-NO", "en-NZ", "en-PL", "en-PT", "en-RO", "en-SE", "en-US", "es-ES", "es-US", "fi-FI", "fr-BE", "fr-CA", "fr-CH", "fr-FR", "it-CH", "it-IT", "nb-NO", "nl-BE", "nl-NL", "pl-PL", "pt-PT", "ro-RO", "sv-FI", "sv-SE"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session", "on_session"}},
				{Name: "subscriptions", Location: "form", Enum: []string{""}},
			}},
			{Name: "konbini", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "confirmation_number", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "expires_after_days", Location: "form", Enum: []string{""}},
				{Name: "expires_at", Location: "form", Enum: []string{""}},
				{Name: "product_description", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "kr_card", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "link", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "mb_way", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "mobilepay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "multibanco", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "naver_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "nz_bank_account", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "oxxo", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "expires_after_days", Location: "form", Type: "integer"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "p24", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
				{Name: "tos_shown_and_accepted", Location: "form", Type: "boolean"},
			}},
			{Name: "pay_by_bank", Location: "form", Type: "object", Enum: []string{""}},
			{Name: "payco", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "paynow", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "paypal", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "preferred_locale", Location: "form", Type: "string", Enum: []string{"cs-CZ", "da-DK", "de-AT", "de-DE", "de-LU", "el-GR", "en-GB", "en-US", "es-ES", "fi-FI", "fr-BE", "fr-FR", "fr-LU", "hu-HU", "it-IT", "nl-BE", "nl-NL", "pl-PL", "pt-PT", "sk-SK", "sv-SE"}},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "risk_correlation_id", Location: "form", Type: "string"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "payto", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "pix", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount_includes_iof", Location: "form", Type: "string", Enum: []string{"always", "never"}},
				{Name: "expires_after_seconds", Location: "form", Type: "integer"},
				{Name: "expires_at", Location: "form", Type: "integer"},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session"}},
			}},
			{Name: "promptpay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "revolut_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "samsung_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "satispay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "sepa_debit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "sofort", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "preferred_language", Location: "form", Type: "string", Enum: []string{"", "de", "en", "es", "fr", "it", "nl", "pl"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "swish", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "reference", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "twint", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "upi", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "financial_connections", Location: "form", Type: "object"},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "networks", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
				{Name: "transaction_purpose", Location: "form", Type: "string", Enum: []string{"", "goods", "other", "services", "unspecified"}},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant", "microdeposits"}},
			}},
			{Name: "wechat_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "app_id", Location: "form", Type: "string"},
				{Name: "client", Location: "form", Type: "string", Enum: []string{"android", "ios", "web"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "zip", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
		}},
		{Name: "payment_method_types", Location: "form", Type: "array"},
		{Name: "radar_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "session", Location: "form", Type: "string"},
		}},
		{Name: "receipt_email", Location: "form", Type: "string"},
		{Name: "return_url", Location: "form", Type: "string"},
		{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"off_session", "on_session"}},
		{Name: "shipping", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "carrier", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
			{Name: "tracking_number", Location: "form", Type: "string"},
		}},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
		{Name: "statement_descriptor_suffix", Location: "form", Type: "string"},
		{Name: "transfer_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount", Location: "form", Type: "integer"},
			{Name: "destination", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "transfer_group", Location: "form", Type: "string"},
		{Name: "use_stripe_sdk", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_intents/search", OperationID: "GetPaymentIntentsSearch", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "page", Location: "query", Type: "string"},
		{Name: "query", Location: "query", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_intents/{intent}", OperationID: "GetPaymentIntentsIntent", Params: []ParameterValidation{
		{Name: "client_secret", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "intent", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_intents/{intent}", OperationID: "PostPaymentIntentsIntent", Params: []ParameterValidation{
		{Name: "intent", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "amount_details", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "discount_amount", Location: "form", Enum: []string{""}},
			{Name: "enforce_arithmetic_validation", Location: "form", Type: "boolean"},
			{Name: "line_items", Location: "form", Enum: []string{""}},
			{Name: "shipping", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount", Location: "form", Enum: []string{""}},
				{Name: "from_postal_code", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "to_postal_code", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "tax", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "total_tax_amount", Location: "form", Required: true, Type: "integer"},
			}},
		}},
		{Name: "application_fee_amount", Location: "form", Enum: []string{""}},
		{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"automatic", "automatic_async", "manual"}},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "excluded_payment_method_types", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "hooks", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "inputs", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tax", Location: "form", Type: "object"},
			}},
		}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "payment_details", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "customer_reference", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "order_reference", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "payment_method", Location: "form", Type: "string"},
		{Name: "payment_method_configuration", Location: "form", Type: "string"},
		{Name: "payment_method_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "institution_number", Location: "form", Required: true, Type: "string"},
				{Name: "transit_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "affirm", Location: "form", Type: "object"},
			{Name: "afterpay_clearpay", Location: "form", Type: "object"},
			{Name: "alipay", Location: "form", Type: "object"},
			{Name: "allow_redisplay", Location: "form", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
			{Name: "alma", Location: "form", Type: "object"},
			{Name: "amazon_pay", Location: "form", Type: "object"},
			{Name: "au_becs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bsb_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "sort_code", Location: "form", Type: "string"},
			}},
			{Name: "bancontact", Location: "form", Type: "object"},
			{Name: "billie", Location: "form", Type: "object"},
			{Name: "billing_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "email", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "name", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "tax_id", Location: "form", Type: "string"},
			}},
			{Name: "blik", Location: "form", Type: "object"},
			{Name: "boleto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tax_id", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "cashapp", Location: "form", Type: "object"},
			{Name: "crypto", Location: "form", Type: "object"},
			{Name: "customer_balance", Location: "form", Type: "object"},
			{Name: "eps", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"arzte_und_apotheker_bank", "austrian_anadi_bank_ag", "bank_austria", "bankhaus_carl_spangler", "bankhaus_schelhammer_und_schattera_ag", "bawag_psk_ag", "bks_bank_ag", "brull_kallmus_bank_ag", "btv_vier_lander_bank", "capital_bank_grawe_gruppe_ag", "deutsche_bank_ag", "dolomitenbank", "easybank_ag", "erste_bank_und_sparkassen", "hypo_alpeadriabank_international_ag", "hypo_bank_burgenland_aktiengesellschaft", "hypo_noe_lb_fur_niederosterreich_u_wien", "hypo_oberosterreich_salzburg_steiermark", "hypo_tirol_bank_ag", "hypo_vorarlberg_bank_ag", "marchfelder_bank", "oberbank_ag", "raiffeisen_bankengruppe_osterreich", "schoellerbank_ag", "sparda_bank_wien", "volksbank_gruppe", "volkskreditbank_ag", "vr_bank_braunau"}},
			}},
			{Name: "fpx", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Required: true, Type: "string", Enum: []string{"affin_bank", "agrobank", "alliance_bank", "ambank", "bank_islam", "bank_muamalat", "bank_of_china", "bank_rakyat", "bsn", "cimb", "deutsche_bank", "hong_leong_bank", "hsbc", "kfh", "maybank2e", "maybank2u", "ocbc", "pb_enterprise", "public_bank", "rhb", "standard_chartered", "uob"}},
			}},
			{Name: "giropay", Location: "form", Type: "object"},
			{Name: "grabpay", Location: "form", Type: "object"},
			{Name: "ideal", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"abn_amro", "adyen", "asn_bank", "bunq", "buut", "finom", "handelsbanken", "ing", "knab", "mollie", "moneyou", "n26", "nn", "rabobank", "regiobank", "revolut", "sns_bank", "triodos_bank", "van_lanschot", "yoursafe"}},
			}},
			{Name: "interac_present", Location: "form", Type: "object"},
			{Name: "kakao_pay", Location: "form", Type: "object"},
			{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "dob", Location: "form", Type: "object"},
			}},
			{Name: "konbini", Location: "form", Type: "object"},
			{Name: "kr_card", Location: "form", Type: "object"},
			{Name: "link", Location: "form", Type: "object"},
			{Name: "mb_way", Location: "form", Type: "object"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "mobilepay", Location: "form", Type: "object"},
			{Name: "multibanco", Location: "form", Type: "object"},
			{Name: "naver_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "funding", Location: "form", Type: "string", Enum: []string{"card", "points"}},
			}},
			{Name: "nz_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_name", Location: "form", Type: "string"},
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bank_code", Location: "form", Required: true, Type: "string"},
				{Name: "branch_code", Location: "form", Required: true, Type: "string"},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "suffix", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "oxxo", Location: "form", Type: "object"},
			{Name: "p24", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"alior_bank", "bank_millennium", "bank_nowy_bfg_sa", "bank_pekao_sa", "banki_spbdzielcze", "blik", "bnp_paribas", "boz", "citi_handlowy", "credit_agricole", "envelobank", "etransfer_pocztowy24", "getin_bank", "ideabank", "ing", "inteligo", "mbank_mtransfer", "nest_przelew", "noble_pay", "pbac_z_ipko", "plus_bank", "santander_przelew24", "tmobile_usbugi_bankowe", "toyota_bank", "velobank", "volkswagen_bank"}},
			}},
			{Name: "pay_by_bank", Location: "form", Type: "object"},
			{Name: "payco", Location: "form", Type: "object"},
			{Name: "paynow", Location: "form", Type: "object"},
			{Name: "paypal", Location: "form", Type: "object"},
			{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "bsb_number", Location: "form", Type: "string"},
				{Name: "pay_id", Location: "form", Type: "string"},
			}},
			{Name: "pix", Location: "form", Type: "object"},
			{Name: "promptpay", Location: "form", Type: "object"},
			{Name: "radar_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "session", Location: "form", Type: "string"},
			}},
			{Name: "revolut_pay", Location: "form", Type: "object"},
			{Name: "samsung_pay", Location: "form", Type: "object"},
			{Name: "satispay", Location: "form", Type: "object"},
			{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "iban", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "sofort", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "country", Location: "form", Required: true, Type: "string", Enum: []string{"AT", "BE", "DE", "ES", "IT", "NL"}},
			}},
			{Name: "sunbit", Location: "form", Type: "object"},
			{Name: "swish", Location: "form", Type: "object"},
			{Name: "twint", Location: "form", Type: "object"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"acss_debit", "affirm", "afterpay_clearpay", "alipay", "alma", "amazon_pay", "au_becs_debit", "bacs_debit", "bancontact", "billie", "blik", "boleto", "cashapp", "crypto", "customer_balance", "eps", "fpx", "giropay", "grabpay", "ideal", "kakao_pay", "klarna", "konbini", "kr_card", "link", "mb_way", "mobilepay", "multibanco", "naver_pay", "nz_bank_account", "oxxo", "p24", "pay_by_bank", "payco", "paynow", "paypal", "payto", "pix", "promptpay", "revolut_pay", "samsung_pay", "satispay", "sepa_debit", "sofort", "sunbit", "swish", "twint", "upi", "us_bank_account", "wechat_pay", "zip"}},
			{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "savings"}},
				{Name: "financial_connections_account", Location: "form", Type: "string"},
				{Name: "routing_number", Location: "form", Type: "string"},
			}},
			{Name: "wechat_pay", Location: "form", Type: "object"},
			{Name: "zip", Location: "form", Type: "object"},
		}},
		{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant", "microdeposits"}},
			}},
			{Name: "affirm", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "preferred_locale", Location: "form", Type: "string"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "afterpay_clearpay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "alipay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "alma", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "amazon_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "au_becs_debit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "bacs_debit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "bancontact", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "preferred_language", Location: "form", Type: "string", Enum: []string{"de", "en", "fr", "nl"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "billie", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "blik", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "code", Location: "form", Type: "string"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none"}},
			}},
			{Name: "boleto", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "expires_after_days", Location: "form", Type: "integer"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
			}},
			{Name: "card", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "cvc_token", Location: "form", Type: "string"},
				{Name: "installments", Location: "form", Type: "object"},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "network", Location: "form", Type: "string", Enum: []string{"amex", "cartes_bancaires", "diners", "discover", "eftpos_au", "girocard", "interac", "jcb", "link", "mastercard", "unionpay", "unknown", "visa"}},
				{Name: "request_extended_authorization", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_incremental_authorization", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_multicapture", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_overcapture", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_three_d_secure", Location: "form", Type: "string", Enum: []string{"any", "automatic", "challenge"}},
				{Name: "require_cvc_recollection", Location: "form", Type: "boolean"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "statement_descriptor_suffix_kana", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "statement_descriptor_suffix_kanji", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "three_d_secure", Location: "form", Type: "object"},
			}},
			{Name: "card_present", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual", "manual_preferred"}},
				{Name: "request_extended_authorization", Location: "form", Type: "boolean"},
				{Name: "request_incremental_authorization_support", Location: "form", Type: "boolean"},
				{Name: "routing", Location: "form", Type: "object"},
			}},
			{Name: "cashapp", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
			}},
			{Name: "crypto", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "customer_balance", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "bank_transfer", Location: "form", Type: "object"},
				{Name: "funding_type", Location: "form", Type: "string", Enum: []string{"bank_transfer"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "eps", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "fpx", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "giropay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "grabpay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "ideal", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "interac_present", Location: "form", Type: "object", Enum: []string{""}},
			{Name: "kakao_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "klarna", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "on_demand", Location: "form", Type: "object"},
				{Name: "preferred_locale", Location: "form", Type: "string", Enum: []string{"cs-CZ", "da-DK", "de-AT", "de-CH", "de-DE", "el-GR", "en-AT", "en-AU", "en-BE", "en-CA", "en-CH", "en-CZ", "en-DE", "en-DK", "en-ES", "en-FI", "en-FR", "en-GB", "en-GR", "en-IE", "en-IT", "en-NL", "en-NO", "en-NZ", "en-PL", "en-PT", "en-RO", "en-SE", "en-US", "es-ES", "es-US", "fi-FI", "fr-BE", "fr-CA", "fr-CH", "fr-FR", "it-CH", "it-IT", "nb-NO", "nl-BE", "nl-NL", "pl-PL", "pt-PT", "ro-RO", "sv-FI", "sv-SE"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session", "on_session"}},
				{Name: "subscriptions", Location: "form", Enum: []string{""}},
			}},
			{Name: "konbini", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "confirmation_number", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "expires_after_days", Location: "form", Enum: []string{""}},
				{Name: "expires_at", Location: "form", Enum: []string{""}},
				{Name: "product_description", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "kr_card", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "link", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "mb_way", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "mobilepay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "multibanco", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "naver_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "nz_bank_account", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "oxxo", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "expires_after_days", Location: "form", Type: "integer"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "p24", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
				{Name: "tos_shown_and_accepted", Location: "form", Type: "boolean"},
			}},
			{Name: "pay_by_bank", Location: "form", Type: "object", Enum: []string{""}},
			{Name: "payco", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "paynow", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "paypal", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "preferred_locale", Location: "form", Type: "string", Enum: []string{"cs-CZ", "da-DK", "de-AT", "de-DE", "de-LU", "el-GR", "en-GB", "en-US", "es-ES", "fi-FI", "fr-BE", "fr-FR", "fr-LU", "hu-HU", "it-IT", "nl-BE", "nl-NL", "pl-PL", "pt-PT", "sk-SK", "sv-SE"}},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "risk_correlation_id", Location: "form", Type: "string"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "payto", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "pix", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount_includes_iof", Location: "form", Type: "string", Enum: []string{"always", "never"}},
				{Name: "expires_after_seconds", Location: "form", Type: "integer"},
				{Name: "expires_at", Location: "form", Type: "integer"},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session"}},
			}},
			{Name: "promptpay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "revolut_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "samsung_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "satispay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "sepa_debit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "sofort", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "preferred_language", Location: "form", Type: "string", Enum: []string{"", "de", "en", "es", "fr", "it", "nl", "pl"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "swish", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "reference", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "twint", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "upi", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "financial_connections", Location: "form", Type: "object"},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "networks", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
				{Name: "transaction_purpose", Location: "form", Type: "string", Enum: []string{"", "goods", "other", "services", "unspecified"}},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant", "microdeposits"}},
			}},
			{Name: "wechat_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "app_id", Location: "form", Type: "string"},
				{Name: "client", Location: "form", Type: "string", Enum: []string{"android", "ios", "web"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "zip", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
		}},
		{Name: "payment_method_types", Location: "form", Type: "array"},
		{Name: "receipt_email", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "off_session", "on_session"}},
		{Name: "shipping", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "carrier", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
			{Name: "tracking_number", Location: "form", Type: "string"},
		}},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
		{Name: "statement_descriptor_suffix", Location: "form", Type: "string"},
		{Name: "transfer_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount", Location: "form", Type: "integer"},
		}},
		{Name: "transfer_group", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_intents/{intent}/amount_details_line_items", OperationID: "GetPaymentIntentsIntentAmountDetailsLineItems", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "intent", Location: "path", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_intents/{intent}/apply_customer_balance", OperationID: "PostPaymentIntentsIntentApplyCustomerBalance", Params: []ParameterValidation{
		{Name: "intent", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_intents/{intent}/cancel", OperationID: "PostPaymentIntentsIntentCancel", Params: []ParameterValidation{
		{Name: "intent", Location: "path", Required: true, Type: "string"},
		{Name: "cancellation_reason", Location: "form", Type: "string", Enum: []string{"abandoned", "duplicate", "fraudulent", "requested_by_customer"}},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_intents/{intent}/capture", OperationID: "PostPaymentIntentsIntentCapture", Params: []ParameterValidation{
		{Name: "intent", Location: "path", Required: true, Type: "string"},
		{Name: "amount_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "discount_amount", Location: "form", Enum: []string{""}},
			{Name: "enforce_arithmetic_validation", Location: "form", Type: "boolean"},
			{Name: "line_items", Location: "form", Enum: []string{""}},
			{Name: "shipping", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount", Location: "form", Enum: []string{""}},
				{Name: "from_postal_code", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "to_postal_code", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "tax", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "total_tax_amount", Location: "form", Required: true, Type: "integer"},
			}},
		}},
		{Name: "amount_to_capture", Location: "form", Type: "integer"},
		{Name: "application_fee_amount", Location: "form", Type: "integer"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "final_capture", Location: "form", Type: "boolean"},
		{Name: "hooks", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "inputs", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tax", Location: "form", Type: "object"},
			}},
		}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "payment_details", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "customer_reference", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "order_reference", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
		{Name: "statement_descriptor_suffix", Location: "form", Type: "string"},
		{Name: "transfer_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount", Location: "form", Type: "integer"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_intents/{intent}/confirm", OperationID: "PostPaymentIntentsIntentConfirm", Params: []ParameterValidation{
		{Name: "intent", Location: "path", Required: true, Type: "string"},
		{Name: "amount_details", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "discount_amount", Location: "form", Enum: []string{""}},
			{Name: "enforce_arithmetic_validation", Location: "form", Type: "boolean"},
			{Name: "line_items", Location: "form", Enum: []string{""}},
			{Name: "shipping", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount", Location: "form", Enum: []string{""}},
				{Name: "from_postal_code", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "to_postal_code", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "tax", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "total_tax_amount", Location: "form", Required: true, Type: "integer"},
			}},
		}},
		{Name: "amount_to_confirm", Location: "form", Type: "integer"},
		{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"automatic", "automatic_async", "manual"}},
		{Name: "client_secret", Location: "form", Type: "string"},
		{Name: "confirmation_token", Location: "form", Type: "string"},
		{Name: "error_on_requires_action", Location: "form", Type: "boolean"},
		{Name: "excluded_payment_method_types", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "hooks", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "inputs", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tax", Location: "form", Type: "object"},
			}},
		}},
		{Name: "mandate", Location: "form", Type: "string"},
		{Name: "mandate_data", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "customer_acceptance", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "accepted_at", Location: "form", Type: "integer"},
				{Name: "offline", Location: "form", Type: "object"},
				{Name: "online", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"offline", "online"}},
			}},
		}},
		{Name: "off_session", Location: "form", Enum: []string{"one_off", "recurring"}},
		{Name: "payment_details", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "customer_reference", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "order_reference", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "payment_method", Location: "form", Type: "string"},
		{Name: "payment_method_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "institution_number", Location: "form", Required: true, Type: "string"},
				{Name: "transit_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "affirm", Location: "form", Type: "object"},
			{Name: "afterpay_clearpay", Location: "form", Type: "object"},
			{Name: "alipay", Location: "form", Type: "object"},
			{Name: "allow_redisplay", Location: "form", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
			{Name: "alma", Location: "form", Type: "object"},
			{Name: "amazon_pay", Location: "form", Type: "object"},
			{Name: "au_becs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bsb_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "sort_code", Location: "form", Type: "string"},
			}},
			{Name: "bancontact", Location: "form", Type: "object"},
			{Name: "billie", Location: "form", Type: "object"},
			{Name: "billing_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "email", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "name", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "tax_id", Location: "form", Type: "string"},
			}},
			{Name: "blik", Location: "form", Type: "object"},
			{Name: "boleto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tax_id", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "cashapp", Location: "form", Type: "object"},
			{Name: "crypto", Location: "form", Type: "object"},
			{Name: "customer_balance", Location: "form", Type: "object"},
			{Name: "eps", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"arzte_und_apotheker_bank", "austrian_anadi_bank_ag", "bank_austria", "bankhaus_carl_spangler", "bankhaus_schelhammer_und_schattera_ag", "bawag_psk_ag", "bks_bank_ag", "brull_kallmus_bank_ag", "btv_vier_lander_bank", "capital_bank_grawe_gruppe_ag", "deutsche_bank_ag", "dolomitenbank", "easybank_ag", "erste_bank_und_sparkassen", "hypo_alpeadriabank_international_ag", "hypo_bank_burgenland_aktiengesellschaft", "hypo_noe_lb_fur_niederosterreich_u_wien", "hypo_oberosterreich_salzburg_steiermark", "hypo_tirol_bank_ag", "hypo_vorarlberg_bank_ag", "marchfelder_bank", "oberbank_ag", "raiffeisen_bankengruppe_osterreich", "schoellerbank_ag", "sparda_bank_wien", "volksbank_gruppe", "volkskreditbank_ag", "vr_bank_braunau"}},
			}},
			{Name: "fpx", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Required: true, Type: "string", Enum: []string{"affin_bank", "agrobank", "alliance_bank", "ambank", "bank_islam", "bank_muamalat", "bank_of_china", "bank_rakyat", "bsn", "cimb", "deutsche_bank", "hong_leong_bank", "hsbc", "kfh", "maybank2e", "maybank2u", "ocbc", "pb_enterprise", "public_bank", "rhb", "standard_chartered", "uob"}},
			}},
			{Name: "giropay", Location: "form", Type: "object"},
			{Name: "grabpay", Location: "form", Type: "object"},
			{Name: "ideal", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"abn_amro", "adyen", "asn_bank", "bunq", "buut", "finom", "handelsbanken", "ing", "knab", "mollie", "moneyou", "n26", "nn", "rabobank", "regiobank", "revolut", "sns_bank", "triodos_bank", "van_lanschot", "yoursafe"}},
			}},
			{Name: "interac_present", Location: "form", Type: "object"},
			{Name: "kakao_pay", Location: "form", Type: "object"},
			{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "dob", Location: "form", Type: "object"},
			}},
			{Name: "konbini", Location: "form", Type: "object"},
			{Name: "kr_card", Location: "form", Type: "object"},
			{Name: "link", Location: "form", Type: "object"},
			{Name: "mb_way", Location: "form", Type: "object"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "mobilepay", Location: "form", Type: "object"},
			{Name: "multibanco", Location: "form", Type: "object"},
			{Name: "naver_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "funding", Location: "form", Type: "string", Enum: []string{"card", "points"}},
			}},
			{Name: "nz_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_name", Location: "form", Type: "string"},
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bank_code", Location: "form", Required: true, Type: "string"},
				{Name: "branch_code", Location: "form", Required: true, Type: "string"},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "suffix", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "oxxo", Location: "form", Type: "object"},
			{Name: "p24", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"alior_bank", "bank_millennium", "bank_nowy_bfg_sa", "bank_pekao_sa", "banki_spbdzielcze", "blik", "bnp_paribas", "boz", "citi_handlowy", "credit_agricole", "envelobank", "etransfer_pocztowy24", "getin_bank", "ideabank", "ing", "inteligo", "mbank_mtransfer", "nest_przelew", "noble_pay", "pbac_z_ipko", "plus_bank", "santander_przelew24", "tmobile_usbugi_bankowe", "toyota_bank", "velobank", "volkswagen_bank"}},
			}},
			{Name: "pay_by_bank", Location: "form", Type: "object"},
			{Name: "payco", Location: "form", Type: "object"},
			{Name: "paynow", Location: "form", Type: "object"},
			{Name: "paypal", Location: "form", Type: "object"},
			{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "bsb_number", Location: "form", Type: "string"},
				{Name: "pay_id", Location: "form", Type: "string"},
			}},
			{Name: "pix", Location: "form", Type: "object"},
			{Name: "promptpay", Location: "form", Type: "object"},
			{Name: "radar_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "session", Location: "form", Type: "string"},
			}},
			{Name: "revolut_pay", Location: "form", Type: "object"},
			{Name: "samsung_pay", Location: "form", Type: "object"},
			{Name: "satispay", Location: "form", Type: "object"},
			{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "iban", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "sofort", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "country", Location: "form", Required: true, Type: "string", Enum: []string{"AT", "BE", "DE", "ES", "IT", "NL"}},
			}},
			{Name: "sunbit", Location: "form", Type: "object"},
			{Name: "swish", Location: "form", Type: "object"},
			{Name: "twint", Location: "form", Type: "object"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"acss_debit", "affirm", "afterpay_clearpay", "alipay", "alma", "amazon_pay", "au_becs_debit", "bacs_debit", "bancontact", "billie", "blik", "boleto", "cashapp", "crypto", "customer_balance", "eps", "fpx", "giropay", "grabpay", "ideal", "kakao_pay", "klarna", "konbini", "kr_card", "link", "mb_way", "mobilepay", "multibanco", "naver_pay", "nz_bank_account", "oxxo", "p24", "pay_by_bank", "payco", "paynow", "paypal", "payto", "pix", "promptpay", "revolut_pay", "samsung_pay", "satispay", "sepa_debit", "sofort", "sunbit", "swish", "twint", "upi", "us_bank_account", "wechat_pay", "zip"}},
			{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "savings"}},
				{Name: "financial_connections_account", Location: "form", Type: "string"},
				{Name: "routing_number", Location: "form", Type: "string"},
			}},
			{Name: "wechat_pay", Location: "form", Type: "object"},
			{Name: "zip", Location: "form", Type: "object"},
		}},
		{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant", "microdeposits"}},
			}},
			{Name: "affirm", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "preferred_locale", Location: "form", Type: "string"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "afterpay_clearpay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "alipay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "alma", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "amazon_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "au_becs_debit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "bacs_debit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "bancontact", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "preferred_language", Location: "form", Type: "string", Enum: []string{"de", "en", "fr", "nl"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "billie", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "blik", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "code", Location: "form", Type: "string"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none"}},
			}},
			{Name: "boleto", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "expires_after_days", Location: "form", Type: "integer"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
			}},
			{Name: "card", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "cvc_token", Location: "form", Type: "string"},
				{Name: "installments", Location: "form", Type: "object"},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "network", Location: "form", Type: "string", Enum: []string{"amex", "cartes_bancaires", "diners", "discover", "eftpos_au", "girocard", "interac", "jcb", "link", "mastercard", "unionpay", "unknown", "visa"}},
				{Name: "request_extended_authorization", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_incremental_authorization", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_multicapture", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_overcapture", Location: "form", Type: "string", Enum: []string{"if_available", "never"}},
				{Name: "request_three_d_secure", Location: "form", Type: "string", Enum: []string{"any", "automatic", "challenge"}},
				{Name: "require_cvc_recollection", Location: "form", Type: "boolean"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "statement_descriptor_suffix_kana", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "statement_descriptor_suffix_kanji", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "three_d_secure", Location: "form", Type: "object"},
			}},
			{Name: "card_present", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"manual", "manual_preferred"}},
				{Name: "request_extended_authorization", Location: "form", Type: "boolean"},
				{Name: "request_incremental_authorization_support", Location: "form", Type: "boolean"},
				{Name: "routing", Location: "form", Type: "object"},
			}},
			{Name: "cashapp", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
			}},
			{Name: "crypto", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "customer_balance", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "bank_transfer", Location: "form", Type: "object"},
				{Name: "funding_type", Location: "form", Type: "string", Enum: []string{"bank_transfer"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "eps", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "fpx", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "giropay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "grabpay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "ideal", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "interac_present", Location: "form", Type: "object", Enum: []string{""}},
			{Name: "kakao_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "klarna", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "on_demand", Location: "form", Type: "object"},
				{Name: "preferred_locale", Location: "form", Type: "string", Enum: []string{"cs-CZ", "da-DK", "de-AT", "de-CH", "de-DE", "el-GR", "en-AT", "en-AU", "en-BE", "en-CA", "en-CH", "en-CZ", "en-DE", "en-DK", "en-ES", "en-FI", "en-FR", "en-GB", "en-GR", "en-IE", "en-IT", "en-NL", "en-NO", "en-NZ", "en-PL", "en-PT", "en-RO", "en-SE", "en-US", "es-ES", "es-US", "fi-FI", "fr-BE", "fr-CA", "fr-CH", "fr-FR", "it-CH", "it-IT", "nb-NO", "nl-BE", "nl-NL", "pl-PL", "pt-PT", "ro-RO", "sv-FI", "sv-SE"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session", "on_session"}},
				{Name: "subscriptions", Location: "form", Enum: []string{""}},
			}},
			{Name: "konbini", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "confirmation_number", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "expires_after_days", Location: "form", Enum: []string{""}},
				{Name: "expires_at", Location: "form", Enum: []string{""}},
				{Name: "product_description", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "kr_card", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "link", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "mb_way", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "mobilepay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "multibanco", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "naver_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "nz_bank_account", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "oxxo", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "expires_after_days", Location: "form", Type: "integer"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "p24", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
				{Name: "tos_shown_and_accepted", Location: "form", Type: "boolean"},
			}},
			{Name: "pay_by_bank", Location: "form", Type: "object", Enum: []string{""}},
			{Name: "payco", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "paynow", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "paypal", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "preferred_locale", Location: "form", Type: "string", Enum: []string{"cs-CZ", "da-DK", "de-AT", "de-DE", "de-LU", "el-GR", "en-GB", "en-US", "es-ES", "fi-FI", "fr-BE", "fr-FR", "fr-LU", "hu-HU", "it-IT", "nl-BE", "nl-NL", "pl-PL", "pt-PT", "sk-SK", "sv-SE"}},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "risk_correlation_id", Location: "form", Type: "string"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "payto", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "pix", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount_includes_iof", Location: "form", Type: "string", Enum: []string{"always", "never"}},
				{Name: "expires_after_seconds", Location: "form", Type: "integer"},
				{Name: "expires_at", Location: "form", Type: "integer"},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none", "off_session"}},
			}},
			{Name: "promptpay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "revolut_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "samsung_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "satispay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"", "manual"}},
			}},
			{Name: "sepa_debit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
			}},
			{Name: "sofort", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "preferred_language", Location: "form", Type: "string", Enum: []string{"", "de", "en", "es", "fr", "it", "nl", "pl"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session"}},
			}},
			{Name: "swish", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "reference", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "twint", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "upi", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "financial_connections", Location: "form", Type: "object"},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "networks", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
				{Name: "target_date", Location: "form", Type: "string"},
				{Name: "transaction_purpose", Location: "form", Type: "string", Enum: []string{"", "goods", "other", "services", "unspecified"}},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant", "microdeposits"}},
			}},
			{Name: "wechat_pay", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "app_id", Location: "form", Type: "string"},
				{Name: "client", Location: "form", Type: "string", Enum: []string{"android", "ios", "web"}},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
			{Name: "zip", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"none"}},
			}},
		}},
		{Name: "payment_method_types", Location: "form", Type: "array"},
		{Name: "radar_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "session", Location: "form", Type: "string"},
		}},
		{Name: "receipt_email", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "return_url", Location: "form", Type: "string"},
		{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "off_session", "on_session"}},
		{Name: "shipping", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "carrier", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
			{Name: "tracking_number", Location: "form", Type: "string"},
		}},
		{Name: "use_stripe_sdk", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_intents/{intent}/increment_authorization", OperationID: "PostPaymentIntentsIntentIncrementAuthorization", Params: []ParameterValidation{
		{Name: "intent", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Required: true, Type: "integer"},
		{Name: "amount_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "discount_amount", Location: "form", Enum: []string{""}},
			{Name: "enforce_arithmetic_validation", Location: "form", Type: "boolean"},
			{Name: "line_items", Location: "form", Enum: []string{""}},
			{Name: "shipping", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount", Location: "form", Enum: []string{""}},
				{Name: "from_postal_code", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "to_postal_code", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "tax", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "total_tax_amount", Location: "form", Required: true, Type: "integer"},
			}},
		}},
		{Name: "application_fee_amount", Location: "form", Type: "integer"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "hooks", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "inputs", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tax", Location: "form", Type: "object"},
			}},
		}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "payment_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "customer_reference", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "order_reference", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
		{Name: "transfer_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount", Location: "form", Type: "integer"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_intents/{intent}/verify_microdeposits", OperationID: "PostPaymentIntentsIntentVerifyMicrodeposits", Params: []ParameterValidation{
		{Name: "intent", Location: "path", Required: true, Type: "string"},
		{Name: "amounts", Location: "form", Type: "array"},
		{Name: "client_secret", Location: "form", Type: "string"},
		{Name: "descriptor_code", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_links", OperationID: "GetPaymentLinks", Params: []ParameterValidation{
		{Name: "active", Location: "query", Type: "boolean"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_links", OperationID: "PostPaymentLinks", Params: []ParameterValidation{
		{Name: "after_completion", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "hosted_confirmation", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "custom_message", Location: "form", Type: "string"},
			}},
			{Name: "redirect", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "url", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"hosted_confirmation", "redirect"}},
		}},
		{Name: "allow_promotion_codes", Location: "form", Type: "boolean"},
		{Name: "application_fee_amount", Location: "form", Type: "integer"},
		{Name: "application_fee_percent", Location: "form", Type: "number"},
		{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "liability", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "billing_address_collection", Location: "form", Type: "string", Enum: []string{"auto", "required"}},
		{Name: "consent_collection", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "payment_method_reuse_agreement", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "position", Location: "form", Required: true, Type: "string", Enum: []string{"auto", "hidden"}},
			}},
			{Name: "promotions", Location: "form", Type: "string", Enum: []string{"auto", "none"}},
			{Name: "terms_of_service", Location: "form", Type: "string", Enum: []string{"none", "required"}},
		}},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "custom_fields", Location: "form", Type: "array"},
		{Name: "custom_text", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "after_submit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "message", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "shipping_address", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "message", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "submit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "message", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "terms_of_service_acceptance", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "message", Location: "form", Required: true, Type: "string"},
			}},
		}},
		{Name: "customer_creation", Location: "form", Type: "string", Enum: []string{"always", "if_required"}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "inactive_message", Location: "form", Type: "string"},
		{Name: "invoice_creation", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "invoice_data", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_tax_ids", Location: "form", Enum: []string{""}},
				{Name: "custom_fields", Location: "form", Enum: []string{""}},
				{Name: "description", Location: "form", Type: "string"},
				{Name: "footer", Location: "form", Type: "string"},
				{Name: "issuer", Location: "form", Type: "object"},
				{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "rendering_options", Location: "form", Type: "object", Enum: []string{""}},
			}},
		}},
		{Name: "line_items", Location: "form", Required: true, Type: "array"},
		{Name: "managed_payments", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Type: "boolean"},
		}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "name_collection", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "business", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "optional", Location: "form", Type: "boolean"},
			}},
			{Name: "individual", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "optional", Location: "form", Type: "boolean"},
			}},
		}},
		{Name: "on_behalf_of", Location: "form", Type: "string"},
		{Name: "optional_items", Location: "form", Type: "array"},
		{Name: "payment_intent_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "capture_method", Location: "form", Type: "string", Enum: []string{"automatic", "automatic_async", "manual"}},
			{Name: "description", Location: "form", Type: "string"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"off_session", "on_session"}},
			{Name: "statement_descriptor", Location: "form", Type: "string"},
			{Name: "statement_descriptor_suffix", Location: "form", Type: "string"},
			{Name: "transfer_group", Location: "form", Type: "string"},
		}},
		{Name: "payment_method_collection", Location: "form", Type: "string", Enum: []string{"always", "if_required"}},
		{Name: "payment_method_types", Location: "form", Type: "array"},
		{Name: "phone_number_collection", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "restrictions", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "completed_sessions", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "limit", Location: "form", Required: true, Type: "integer"},
			}},
		}},
		{Name: "shipping_address_collection", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "allowed_countries", Location: "form", Required: true, Type: "array"},
		}},
		{Name: "shipping_options", Location: "form", Type: "array"},
		{Name: "submit_type", Location: "form", Type: "string", Enum: []string{"auto", "book", "donate", "pay", "subscribe"}},
		{Name: "subscription_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "description", Location: "form", Type: "string"},
			{Name: "invoice_settings", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "issuer", Location: "form", Type: "object"},
			}},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "trial_period_days", Location: "form", Type: "integer"},
			{Name: "trial_settings", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "end_behavior", Location: "form", Required: true, Type: "object"},
			}},
		}},
		{Name: "tax_id_collection", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "required", Location: "form", Type: "string", Enum: []string{"if_supported", "never"}},
		}},
		{Name: "transfer_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount", Location: "form", Type: "integer"},
			{Name: "destination", Location: "form", Required: true, Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_links/{payment_link}", OperationID: "GetPaymentLinksPaymentLink", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "payment_link", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_links/{payment_link}", OperationID: "PostPaymentLinksPaymentLink", Params: []ParameterValidation{
		{Name: "payment_link", Location: "path", Required: true, Type: "string"},
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "after_completion", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "hosted_confirmation", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "custom_message", Location: "form", Type: "string"},
			}},
			{Name: "redirect", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "url", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"hosted_confirmation", "redirect"}},
		}},
		{Name: "allow_promotion_codes", Location: "form", Type: "boolean"},
		{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "liability", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "billing_address_collection", Location: "form", Type: "string", Enum: []string{"auto", "required"}},
		{Name: "custom_fields", Location: "form", Enum: []string{""}},
		{Name: "custom_text", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "after_submit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "message", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "shipping_address", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "message", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "submit", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "message", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "terms_of_service_acceptance", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "message", Location: "form", Required: true, Type: "string"},
			}},
		}},
		{Name: "customer_creation", Location: "form", Type: "string", Enum: []string{"always", "if_required"}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "inactive_message", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "invoice_creation", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "invoice_data", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_tax_ids", Location: "form", Enum: []string{""}},
				{Name: "custom_fields", Location: "form", Enum: []string{""}},
				{Name: "description", Location: "form", Type: "string"},
				{Name: "footer", Location: "form", Type: "string"},
				{Name: "issuer", Location: "form", Type: "object"},
				{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "rendering_options", Location: "form", Type: "object", Enum: []string{""}},
			}},
		}},
		{Name: "line_items", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "name_collection", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "business", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "optional", Location: "form", Type: "boolean"},
			}},
			{Name: "individual", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "optional", Location: "form", Type: "boolean"},
			}},
		}},
		{Name: "optional_items", Location: "form", Enum: []string{""}},
		{Name: "payment_intent_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "description", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
			{Name: "statement_descriptor", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "statement_descriptor_suffix", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "transfer_group", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "payment_method_collection", Location: "form", Type: "string", Enum: []string{"always", "if_required"}},
		{Name: "payment_method_types", Location: "form", Enum: []string{""}},
		{Name: "phone_number_collection", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "restrictions", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "completed_sessions", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "limit", Location: "form", Required: true, Type: "integer"},
			}},
		}},
		{Name: "shipping_address_collection", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "allowed_countries", Location: "form", Required: true, Type: "array"},
		}},
		{Name: "submit_type", Location: "form", Type: "string", Enum: []string{"auto", "book", "donate", "pay", "subscribe"}},
		{Name: "subscription_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "invoice_settings", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "issuer", Location: "form", Type: "object"},
			}},
			{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
			{Name: "trial_period_days", Location: "form", Enum: []string{""}},
			{Name: "trial_settings", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "end_behavior", Location: "form", Required: true, Type: "object"},
			}},
		}},
		{Name: "tax_id_collection", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "required", Location: "form", Type: "string", Enum: []string{"if_supported", "never"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_links/{payment_link}/line_items", OperationID: "GetPaymentLinksPaymentLinkLineItems", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "payment_link", Location: "path", Required: true, Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_method_configurations", OperationID: "GetPaymentMethodConfigurations", Params: []ParameterValidation{
		{Name: "application", Location: "query", Type: "string", Enum: []string{""}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_method_configurations", OperationID: "PostPaymentMethodConfigurations", Params: []ParameterValidation{
		{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "affirm", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "afterpay_clearpay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "alipay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "alma", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "amazon_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "apple_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "apple_pay_later", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "au_becs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "bancontact", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "billie", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "blik", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "boleto", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "cartes_bancaires", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "cashapp", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "crypto", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "customer_balance", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "eps", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "fpx", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "fr_meal_voucher_conecs", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "giropay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "google_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "grabpay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "ideal", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "jcb", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "kakao_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "konbini", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "kr_card", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "link", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "mb_way", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "mobilepay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "multibanco", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "naver_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "nz_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "oxxo", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "p24", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "parent", Location: "form", Type: "string"},
		{Name: "pay_by_bank", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "payco", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "paynow", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "paypal", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "pix", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "promptpay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "revolut_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "samsung_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "satispay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "sofort", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "sunbit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "swish", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "twint", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "wechat_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "zip", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_method_configurations/{configuration}", OperationID: "GetPaymentMethodConfigurationsConfiguration", Params: []ParameterValidation{
		{Name: "configuration", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_method_configurations/{configuration}", OperationID: "PostPaymentMethodConfigurationsConfiguration", Params: []ParameterValidation{
		{Name: "configuration", Location: "path", Required: true, Type: "string"},
		{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "affirm", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "afterpay_clearpay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "alipay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "alma", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "amazon_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "apple_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "apple_pay_later", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "au_becs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "bancontact", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "billie", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "blik", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "boleto", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "cartes_bancaires", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "cashapp", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "crypto", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "customer_balance", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "eps", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "fpx", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "fr_meal_voucher_conecs", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "giropay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "google_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "grabpay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "ideal", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "jcb", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "kakao_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "konbini", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "kr_card", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "link", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "mb_way", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "mobilepay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "multibanco", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "naver_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "nz_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "oxxo", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "p24", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "pay_by_bank", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "payco", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "paynow", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "paypal", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "pix", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "promptpay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "revolut_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "samsung_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "satispay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "sofort", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "sunbit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "swish", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "twint", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "wechat_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
		{Name: "zip", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "display_preference", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preference", Location: "form", Type: "string", Enum: []string{"none", "off", "on"}},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_method_domains", OperationID: "GetPaymentMethodDomains", Params: []ParameterValidation{
		{Name: "domain_name", Location: "query", Type: "string"},
		{Name: "enabled", Location: "query", Type: "boolean"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_method_domains", OperationID: "PostPaymentMethodDomains", Params: []ParameterValidation{
		{Name: "domain_name", Location: "form", Required: true, Type: "string"},
		{Name: "enabled", Location: "form", Type: "boolean"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_method_domains/{payment_method_domain}", OperationID: "GetPaymentMethodDomainsPaymentMethodDomain", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "payment_method_domain", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_method_domains/{payment_method_domain}", OperationID: "PostPaymentMethodDomainsPaymentMethodDomain", Params: []ParameterValidation{
		{Name: "payment_method_domain", Location: "path", Required: true, Type: "string"},
		{Name: "enabled", Location: "form", Type: "boolean"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_method_domains/{payment_method_domain}/validate", OperationID: "PostPaymentMethodDomainsPaymentMethodDomainValidate", Params: []ParameterValidation{
		{Name: "payment_method_domain", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_methods", OperationID: "GetPaymentMethods", Params: []ParameterValidation{
		{Name: "allow_redisplay", Location: "query", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "type", Location: "query", Type: "string", Enum: []string{"acss_debit", "affirm", "afterpay_clearpay", "alipay", "alma", "amazon_pay", "au_becs_debit", "bacs_debit", "bancontact", "billie", "blik", "boleto", "card", "cashapp", "crypto", "custom", "customer_balance", "eps", "fpx", "giropay", "grabpay", "ideal", "kakao_pay", "klarna", "konbini", "kr_card", "link", "mb_way", "mobilepay", "multibanco", "naver_pay", "nz_bank_account", "oxxo", "p24", "pay_by_bank", "payco", "paynow", "paypal", "payto", "pix", "promptpay", "revolut_pay", "samsung_pay", "satispay", "sepa_debit", "sofort", "sunbit", "swish", "twint", "upi", "us_bank_account", "wechat_pay", "zip"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_methods", OperationID: "PostPaymentMethods", Params: []ParameterValidation{
		{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_number", Location: "form", Required: true, Type: "string"},
			{Name: "institution_number", Location: "form", Required: true, Type: "string"},
			{Name: "transit_number", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "affirm", Location: "form", Type: "object"},
		{Name: "afterpay_clearpay", Location: "form", Type: "object"},
		{Name: "alipay", Location: "form", Type: "object"},
		{Name: "allow_redisplay", Location: "form", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
		{Name: "alma", Location: "form", Type: "object"},
		{Name: "amazon_pay", Location: "form", Type: "object"},
		{Name: "au_becs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_number", Location: "form", Required: true, Type: "string"},
			{Name: "bsb_number", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_number", Location: "form", Type: "string"},
			{Name: "sort_code", Location: "form", Type: "string"},
		}},
		{Name: "bancontact", Location: "form", Type: "object"},
		{Name: "billie", Location: "form", Type: "object"},
		{Name: "billing_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "email", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "name", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "tax_id", Location: "form", Type: "string"},
		}},
		{Name: "blik", Location: "form", Type: "object"},
		{Name: "boleto", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "tax_id", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "cvc", Location: "form", Type: "string"},
			{Name: "exp_month", Location: "form", Required: true, Type: "integer"},
			{Name: "exp_year", Location: "form", Required: true, Type: "integer"},
			{Name: "networks", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preferred", Location: "form", Type: "string", Enum: []string{"cartes_bancaires", "mastercard", "visa"}},
			}},
			{Name: "number", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "cashapp", Location: "form", Type: "object"},
		{Name: "crypto", Location: "form", Type: "object"},
		{Name: "custom", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "type", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_balance", Location: "form", Type: "object"},
		{Name: "eps", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "bank", Location: "form", Type: "string", Enum: []string{"arzte_und_apotheker_bank", "austrian_anadi_bank_ag", "bank_austria", "bankhaus_carl_spangler", "bankhaus_schelhammer_und_schattera_ag", "bawag_psk_ag", "bks_bank_ag", "brull_kallmus_bank_ag", "btv_vier_lander_bank", "capital_bank_grawe_gruppe_ag", "deutsche_bank_ag", "dolomitenbank", "easybank_ag", "erste_bank_und_sparkassen", "hypo_alpeadriabank_international_ag", "hypo_bank_burgenland_aktiengesellschaft", "hypo_noe_lb_fur_niederosterreich_u_wien", "hypo_oberosterreich_salzburg_steiermark", "hypo_tirol_bank_ag", "hypo_vorarlberg_bank_ag", "marchfelder_bank", "oberbank_ag", "raiffeisen_bankengruppe_osterreich", "schoellerbank_ag", "sparda_bank_wien", "volksbank_gruppe", "volkskreditbank_ag", "vr_bank_braunau"}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "fpx", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "bank", Location: "form", Required: true, Type: "string", Enum: []string{"affin_bank", "agrobank", "alliance_bank", "ambank", "bank_islam", "bank_muamalat", "bank_of_china", "bank_rakyat", "bsn", "cimb", "deutsche_bank", "hong_leong_bank", "hsbc", "kfh", "maybank2e", "maybank2u", "ocbc", "pb_enterprise", "public_bank", "rhb", "standard_chartered", "uob"}},
		}},
		{Name: "giropay", Location: "form", Type: "object"},
		{Name: "grabpay", Location: "form", Type: "object"},
		{Name: "ideal", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "bank", Location: "form", Type: "string", Enum: []string{"abn_amro", "adyen", "asn_bank", "bunq", "buut", "finom", "handelsbanken", "ing", "knab", "mollie", "moneyou", "n26", "nn", "rabobank", "regiobank", "revolut", "sns_bank", "triodos_bank", "van_lanschot", "yoursafe"}},
		}},
		{Name: "interac_present", Location: "form", Type: "object"},
		{Name: "kakao_pay", Location: "form", Type: "object"},
		{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "dob", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "day", Location: "form", Required: true, Type: "integer"},
				{Name: "month", Location: "form", Required: true, Type: "integer"},
				{Name: "year", Location: "form", Required: true, Type: "integer"},
			}},
		}},
		{Name: "konbini", Location: "form", Type: "object"},
		{Name: "kr_card", Location: "form", Type: "object"},
		{Name: "link", Location: "form", Type: "object"},
		{Name: "mb_way", Location: "form", Type: "object"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "mobilepay", Location: "form", Type: "object"},
		{Name: "multibanco", Location: "form", Type: "object"},
		{Name: "naver_pay", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "funding", Location: "form", Type: "string", Enum: []string{"card", "points"}},
		}},
		{Name: "nz_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_holder_name", Location: "form", Type: "string"},
			{Name: "account_number", Location: "form", Required: true, Type: "string"},
			{Name: "bank_code", Location: "form", Required: true, Type: "string"},
			{Name: "branch_code", Location: "form", Required: true, Type: "string"},
			{Name: "reference", Location: "form", Type: "string"},
			{Name: "suffix", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "oxxo", Location: "form", Type: "object"},
		{Name: "p24", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "bank", Location: "form", Type: "string", Enum: []string{"alior_bank", "bank_millennium", "bank_nowy_bfg_sa", "bank_pekao_sa", "banki_spbdzielcze", "blik", "bnp_paribas", "boz", "citi_handlowy", "credit_agricole", "envelobank", "etransfer_pocztowy24", "getin_bank", "ideabank", "ing", "inteligo", "mbank_mtransfer", "nest_przelew", "noble_pay", "pbac_z_ipko", "plus_bank", "santander_przelew24", "tmobile_usbugi_bankowe", "toyota_bank", "velobank", "volkswagen_bank"}},
		}},
		{Name: "pay_by_bank", Location: "form", Type: "object"},
		{Name: "payco", Location: "form", Type: "object"},
		{Name: "payment_method", Location: "form", Type: "string"},
		{Name: "paynow", Location: "form", Type: "object"},
		{Name: "paypal", Location: "form", Type: "object"},
		{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_number", Location: "form", Type: "string"},
			{Name: "bsb_number", Location: "form", Type: "string"},
			{Name: "pay_id", Location: "form", Type: "string"},
		}},
		{Name: "pix", Location: "form", Type: "object"},
		{Name: "promptpay", Location: "form", Type: "object"},
		{Name: "radar_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "session", Location: "form", Type: "string"},
		}},
		{Name: "revolut_pay", Location: "form", Type: "object"},
		{Name: "samsung_pay", Location: "form", Type: "object"},
		{Name: "satispay", Location: "form", Type: "object"},
		{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "iban", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "sofort", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "country", Location: "form", Required: true, Type: "string", Enum: []string{"AT", "BE", "DE", "ES", "IT", "NL"}},
		}},
		{Name: "sunbit", Location: "form", Type: "object"},
		{Name: "swish", Location: "form", Type: "object"},
		{Name: "twint", Location: "form", Type: "object"},
		{Name: "type", Location: "form", Type: "string", Enum: []string{"acss_debit", "affirm", "afterpay_clearpay", "alipay", "alma", "amazon_pay", "au_becs_debit", "bacs_debit", "bancontact", "billie", "blik", "boleto", "card", "cashapp", "crypto", "custom", "customer_balance", "eps", "fpx", "giropay", "grabpay", "ideal", "kakao_pay", "klarna", "konbini", "kr_card", "link", "mb_way", "mobilepay", "multibanco", "naver_pay", "nz_bank_account", "oxxo", "p24", "pay_by_bank", "payco", "paynow", "paypal", "payto", "pix", "promptpay", "revolut_pay", "samsung_pay", "satispay", "sepa_debit", "sofort", "sunbit", "swish", "twint", "upi", "us_bank_account", "wechat_pay", "zip"}},
		{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "mandate_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "amount", Location: "form", Type: "integer"},
				{Name: "amount_type", Location: "form", Type: "string", Enum: []string{"fixed", "maximum"}},
				{Name: "description", Location: "form", Type: "string"},
				{Name: "end_date", Location: "form", Type: "integer"},
			}},
		}},
		{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
			{Name: "account_number", Location: "form", Type: "string"},
			{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "savings"}},
			{Name: "financial_connections_account", Location: "form", Type: "string"},
			{Name: "routing_number", Location: "form", Type: "string"},
		}},
		{Name: "wechat_pay", Location: "form", Type: "object"},
		{Name: "zip", Location: "form", Type: "object"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_methods/{payment_method}", OperationID: "GetPaymentMethodsPaymentMethod", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "payment_method", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_methods/{payment_method}", OperationID: "PostPaymentMethodsPaymentMethod", Params: []ParameterValidation{
		{Name: "payment_method", Location: "path", Required: true, Type: "string"},
		{Name: "allow_redisplay", Location: "form", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
		{Name: "billing_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "email", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "name", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "tax_id", Location: "form", Type: "string"},
		}},
		{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "exp_month", Location: "form", Type: "integer"},
			{Name: "exp_year", Location: "form", Type: "integer"},
			{Name: "networks", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preferred", Location: "form", Type: "string", Enum: []string{"", "cartes_bancaires", "mastercard", "visa"}},
			}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_number", Location: "form", Type: "string"},
			{Name: "bsb_number", Location: "form", Type: "string"},
			{Name: "pay_id", Location: "form", Type: "string"},
		}},
		{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
			{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "savings"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_methods/{payment_method}/attach", OperationID: "PostPaymentMethodsPaymentMethodAttach", Params: []ParameterValidation{
		{Name: "payment_method", Location: "path", Required: true, Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_methods/{payment_method}/detach", OperationID: "PostPaymentMethodsPaymentMethodDetach", Params: []ParameterValidation{
		{Name: "payment_method", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_records/report_payment", OperationID: "PostPaymentRecordsReportPayment", Params: []ParameterValidation{
		{Name: "amount_requested", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "currency", Location: "form", Required: true, Type: "string"},
			{Name: "value", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "customer_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "customer", Location: "form", Type: "string"},
			{Name: "email", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
		{Name: "customer_presence", Location: "form", Type: "string", Enum: []string{"off_session", "on_session"}},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "failed", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "failed_at", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "guaranteed", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "guaranteed_at", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "initiated_at", Location: "form", Required: true, Type: "integer"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "outcome", Location: "form", Type: "string", Enum: []string{"failed", "guaranteed"}},
		{Name: "payment_method_details", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "billing_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object"},
				{Name: "email", Location: "form", Type: "string"},
				{Name: "name", Location: "form", Type: "string"},
				{Name: "phone", Location: "form", Type: "string"},
			}},
			{Name: "custom", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "display_name", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Type: "string"},
			}},
			{Name: "payment_method", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Type: "string", Enum: []string{"custom"}},
		}},
		{Name: "processor_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "custom", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "payment_reference", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"custom"}},
		}},
		{Name: "shipping_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payment_records/{id}", OperationID: "GetPaymentRecordsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_records/{id}/report_payment_attempt", OperationID: "PostPaymentRecordsIdReportPaymentAttempt", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "failed", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "failed_at", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "guaranteed", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "guaranteed_at", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "initiated_at", Location: "form", Required: true, Type: "integer"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "outcome", Location: "form", Type: "string", Enum: []string{"failed", "guaranteed"}},
		{Name: "payment_method_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "billing_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object"},
				{Name: "email", Location: "form", Type: "string"},
				{Name: "name", Location: "form", Type: "string"},
				{Name: "phone", Location: "form", Type: "string"},
			}},
			{Name: "custom", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "display_name", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Type: "string"},
			}},
			{Name: "payment_method", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Type: "string", Enum: []string{"custom"}},
		}},
		{Name: "shipping_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_records/{id}/report_payment_attempt_canceled", OperationID: "PostPaymentRecordsIdReportPaymentAttemptCanceled", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "canceled_at", Location: "form", Required: true, Type: "integer"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_records/{id}/report_payment_attempt_failed", OperationID: "PostPaymentRecordsIdReportPaymentAttemptFailed", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "failed_at", Location: "form", Required: true, Type: "integer"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_records/{id}/report_payment_attempt_guaranteed", OperationID: "PostPaymentRecordsIdReportPaymentAttemptGuaranteed", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "guaranteed_at", Location: "form", Required: true, Type: "integer"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_records/{id}/report_payment_attempt_informational", OperationID: "PostPaymentRecordsIdReportPaymentAttemptInformational", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "customer_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "customer", Location: "form", Type: "string"},
			{Name: "email", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
		{Name: "description", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "shipping_details", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payment_records/{id}/report_refund", OperationID: "PostPaymentRecordsIdReportRefund", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "currency", Location: "form", Required: true, Type: "string"},
			{Name: "value", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "initiated_at", Location: "form", Type: "integer"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "outcome", Location: "form", Required: true, Type: "string", Enum: []string{"refunded"}},
		{Name: "processor_details", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "custom", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "refund_reference", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"custom"}},
		}},
		{Name: "refunded", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "refunded_at", Location: "form", Required: true, Type: "integer"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payouts", OperationID: "GetPayouts", Params: []ParameterValidation{
		{Name: "arrival_date", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "destination", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payouts", OperationID: "PostPayouts", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Required: true, Type: "integer"},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "destination", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "method", Location: "form", Type: "string", Enum: []string{"instant", "standard"}},
		{Name: "payout_method", Location: "form", Type: "string"},
		{Name: "source_type", Location: "form", Type: "string", Enum: []string{"bank_account", "card", "fpx"}},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/payouts/{payout}", OperationID: "GetPayoutsPayout", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "payout", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payouts/{payout}", OperationID: "PostPayoutsPayout", Params: []ParameterValidation{
		{Name: "payout", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payouts/{payout}/cancel", OperationID: "PostPayoutsPayoutCancel", Params: []ParameterValidation{
		{Name: "payout", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/payouts/{payout}/reverse", OperationID: "PostPayoutsPayoutReverse", Params: []ParameterValidation{
		{Name: "payout", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/plans", OperationID: "GetPlans", Params: []ParameterValidation{
		{Name: "active", Location: "query", Type: "boolean"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "product", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/plans", OperationID: "PostPlans", Params: []ParameterValidation{
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "amount_decimal", Location: "form", Type: "string"},
		{Name: "billing_scheme", Location: "form", Type: "string", Enum: []string{"per_unit", "tiered"}},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "id", Location: "form", Type: "string"},
		{Name: "interval", Location: "form", Required: true, Type: "string", Enum: []string{"day", "month", "week", "year"}},
		{Name: "interval_count", Location: "form", Type: "integer"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "meter", Location: "form", Type: "string"},
		{Name: "nickname", Location: "form", Type: "string"},
		{Name: "product", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "active", Location: "form", Type: "boolean"},
			{Name: "id", Location: "form", Type: "string"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "statement_descriptor", Location: "form", Type: "string"},
			{Name: "tax_code", Location: "form", Type: "string"},
			{Name: "unit_label", Location: "form", Type: "string"},
		}},
		{Name: "tiers", Location: "form", Type: "array"},
		{Name: "tiers_mode", Location: "form", Type: "string", Enum: []string{"graduated", "volume"}},
		{Name: "transform_usage", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "divide_by", Location: "form", Required: true, Type: "integer"},
			{Name: "round", Location: "form", Required: true, Type: "string", Enum: []string{"down", "up"}},
		}},
		{Name: "trial_period_days", Location: "form", Type: "integer"},
		{Name: "usage_type", Location: "form", Type: "string", Enum: []string{"licensed", "metered"}},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/plans/{plan}", OperationID: "DeletePlansPlan", Params: []ParameterValidation{
		{Name: "plan", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/plans/{plan}", OperationID: "GetPlansPlan", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "plan", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/plans/{plan}", OperationID: "PostPlansPlan", Params: []ParameterValidation{
		{Name: "plan", Location: "path", Required: true, Type: "string"},
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "nickname", Location: "form", Type: "string"},
		{Name: "product", Location: "form", Type: "string"},
		{Name: "trial_period_days", Location: "form", Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/prices", OperationID: "GetPrices", Params: []ParameterValidation{
		{Name: "active", Location: "query", Type: "boolean"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "currency", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "lookup_keys", Location: "query", Type: "array"},
		{Name: "product", Location: "query", Type: "string"},
		{Name: "recurring", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "interval", Location: "query", Type: "string", Enum: []string{"day", "month", "week", "year"}},
			{Name: "meter", Location: "query", Type: "string"},
			{Name: "usage_type", Location: "query", Type: "string", Enum: []string{"licensed", "metered"}},
		}},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "type", Location: "query", Type: "string", Enum: []string{"one_time", "recurring"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/prices", OperationID: "PostPrices", Params: []ParameterValidation{
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "billing_scheme", Location: "form", Type: "string", Enum: []string{"per_unit", "tiered"}},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "currency_options", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "custom_unit_amount", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "maximum", Location: "form", Type: "integer"},
			{Name: "minimum", Location: "form", Type: "integer"},
			{Name: "preset", Location: "form", Type: "integer"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "lookup_key", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "nickname", Location: "form", Type: "string"},
		{Name: "product", Location: "form", Type: "string"},
		{Name: "product_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "active", Location: "form", Type: "boolean"},
			{Name: "id", Location: "form", Type: "string"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "statement_descriptor", Location: "form", Type: "string"},
			{Name: "tax_code", Location: "form", Type: "string"},
			{Name: "unit_label", Location: "form", Type: "string"},
		}},
		{Name: "recurring", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "interval", Location: "form", Required: true, Type: "string", Enum: []string{"day", "month", "week", "year"}},
			{Name: "interval_count", Location: "form", Type: "integer"},
			{Name: "meter", Location: "form", Type: "string"},
			{Name: "usage_type", Location: "form", Type: "string", Enum: []string{"licensed", "metered"}},
		}},
		{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
		{Name: "tiers", Location: "form", Type: "array"},
		{Name: "tiers_mode", Location: "form", Type: "string", Enum: []string{"graduated", "volume"}},
		{Name: "transfer_lookup_key", Location: "form", Type: "boolean"},
		{Name: "transform_quantity", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "divide_by", Location: "form", Required: true, Type: "integer"},
			{Name: "round", Location: "form", Required: true, Type: "string", Enum: []string{"down", "up"}},
		}},
		{Name: "unit_amount", Location: "form", Type: "integer"},
		{Name: "unit_amount_decimal", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/prices/search", OperationID: "GetPricesSearch", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "page", Location: "query", Type: "string"},
		{Name: "query", Location: "query", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/prices/{price}", OperationID: "GetPricesPrice", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "price", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/prices/{price}", OperationID: "PostPricesPrice", Params: []ParameterValidation{
		{Name: "price", Location: "path", Required: true, Type: "string"},
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "currency_options", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "lookup_key", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "nickname", Location: "form", Type: "string"},
		{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
		{Name: "transfer_lookup_key", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/products", OperationID: "GetProducts", Params: []ParameterValidation{
		{Name: "active", Location: "query", Type: "boolean"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "ids", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "shippable", Location: "query", Type: "boolean"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "url", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/products", OperationID: "PostProducts", Params: []ParameterValidation{
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "default_price_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "currency", Location: "form", Required: true, Type: "string"},
			{Name: "currency_options", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "custom_unit_amount", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "maximum", Location: "form", Type: "integer"},
				{Name: "minimum", Location: "form", Type: "integer"},
				{Name: "preset", Location: "form", Type: "integer"},
			}},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "recurring", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "interval", Location: "form", Required: true, Type: "string", Enum: []string{"day", "month", "week", "year"}},
				{Name: "interval_count", Location: "form", Type: "integer"},
			}},
			{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
			{Name: "unit_amount", Location: "form", Type: "integer"},
			{Name: "unit_amount_decimal", Location: "form", Type: "string"},
		}},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "id", Location: "form", Type: "string"},
		{Name: "images", Location: "form", Type: "array"},
		{Name: "marketing_features", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "name", Location: "form", Required: true, Type: "string"},
		{Name: "package_dimensions", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "height", Location: "form", Required: true, Type: "number"},
			{Name: "length", Location: "form", Required: true, Type: "number"},
			{Name: "weight", Location: "form", Required: true, Type: "number"},
			{Name: "width", Location: "form", Required: true, Type: "number"},
		}},
		{Name: "shippable", Location: "form", Type: "boolean"},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
		{Name: "tax_code", Location: "form", Type: "string"},
		{Name: "unit_label", Location: "form", Type: "string"},
		{Name: "url", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/products/search", OperationID: "GetProductsSearch", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "page", Location: "query", Type: "string"},
		{Name: "query", Location: "query", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/products/{id}", OperationID: "DeleteProductsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/products/{id}", OperationID: "GetProductsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/products/{id}", OperationID: "PostProductsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "default_price", Location: "form", Type: "string"},
		{Name: "description", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "images", Location: "form", Enum: []string{""}},
		{Name: "marketing_features", Location: "form", Enum: []string{""}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "package_dimensions", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "height", Location: "form", Required: true, Type: "number"},
			{Name: "length", Location: "form", Required: true, Type: "number"},
			{Name: "weight", Location: "form", Required: true, Type: "number"},
			{Name: "width", Location: "form", Required: true, Type: "number"},
		}},
		{Name: "shippable", Location: "form", Type: "boolean"},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
		{Name: "tax_code", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "unit_label", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "url", Location: "form", Type: "string", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/products/{product}/features", OperationID: "GetProductsProductFeatures", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "product", Location: "path", Required: true, Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/products/{product}/features", OperationID: "PostProductsProductFeatures", Params: []ParameterValidation{
		{Name: "product", Location: "path", Required: true, Type: "string"},
		{Name: "entitlement_feature", Location: "form", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/products/{product}/features/{id}", OperationID: "DeleteProductsProductFeaturesId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "product", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/products/{product}/features/{id}", OperationID: "GetProductsProductFeaturesId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "product", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/promotion_codes", OperationID: "GetPromotionCodes", Params: []ParameterValidation{
		{Name: "active", Location: "query", Type: "boolean"},
		{Name: "code", Location: "query", Type: "string"},
		{Name: "coupon", Location: "query", Type: "string"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/promotion_codes", OperationID: "PostPromotionCodes", Params: []ParameterValidation{
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "code", Location: "form", Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "expires_at", Location: "form", Type: "integer"},
		{Name: "max_redemptions", Location: "form", Type: "integer"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "promotion", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "coupon", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"coupon"}},
		}},
		{Name: "restrictions", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "currency_options", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "first_time_transaction", Location: "form", Type: "boolean"},
			{Name: "minimum_amount", Location: "form", Type: "integer"},
			{Name: "minimum_amount_currency", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/promotion_codes/{promotion_code}", OperationID: "GetPromotionCodesPromotionCode", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "promotion_code", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/promotion_codes/{promotion_code}", OperationID: "PostPromotionCodesPromotionCode", Params: []ParameterValidation{
		{Name: "promotion_code", Location: "path", Required: true, Type: "string"},
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "restrictions", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "currency_options", Location: "form", Type: "object", AdditionalProperties: true},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/quotes", OperationID: "GetQuotes", Params: []ParameterValidation{
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"accepted", "canceled", "draft", "open"}},
		{Name: "test_clock", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/quotes", OperationID: "PostQuotes", Params: []ParameterValidation{
		{Name: "application_fee_amount", Location: "form", Enum: []string{""}},
		{Name: "application_fee_percent", Location: "form", Enum: []string{""}},
		{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "liability", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "collection_method", Location: "form", Type: "string", Enum: []string{"charge_automatically", "send_invoice"}},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "default_tax_rates", Location: "form", Enum: []string{""}},
		{Name: "description", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "expires_at", Location: "form", Type: "integer"},
		{Name: "footer", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "from_quote", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "is_revision", Location: "form", Type: "boolean"},
			{Name: "quote", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "header", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "invoice_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "days_until_due", Location: "form", Type: "integer"},
			{Name: "issuer", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "line_items", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "on_behalf_of", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "subscription_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "billing_mode", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "flexible", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"classic", "flexible"}},
			}},
			{Name: "description", Location: "form", Type: "string"},
			{Name: "effective_date", Location: "form", Enum: []string{"", "current_period_end"}},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "trial_period_days", Location: "form", Enum: []string{""}},
		}},
		{Name: "test_clock", Location: "form", Type: "string"},
		{Name: "transfer_data", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "amount", Location: "form", Type: "integer"},
			{Name: "amount_percent", Location: "form", Type: "number"},
			{Name: "destination", Location: "form", Required: true, Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/quotes/{quote}", OperationID: "GetQuotesQuote", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "quote", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/quotes/{quote}", OperationID: "PostQuotesQuote", Params: []ParameterValidation{
		{Name: "quote", Location: "path", Required: true, Type: "string"},
		{Name: "application_fee_amount", Location: "form", Enum: []string{""}},
		{Name: "application_fee_percent", Location: "form", Enum: []string{""}},
		{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "liability", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "collection_method", Location: "form", Type: "string", Enum: []string{"charge_automatically", "send_invoice"}},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "default_tax_rates", Location: "form", Enum: []string{""}},
		{Name: "description", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "expires_at", Location: "form", Type: "integer"},
		{Name: "footer", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "header", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "invoice_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "days_until_due", Location: "form", Type: "integer"},
			{Name: "issuer", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "line_items", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "on_behalf_of", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "subscription_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "description", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "effective_date", Location: "form", Enum: []string{"", "current_period_end"}},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "trial_period_days", Location: "form", Enum: []string{""}},
		}},
		{Name: "transfer_data", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "amount", Location: "form", Type: "integer"},
			{Name: "amount_percent", Location: "form", Type: "number"},
			{Name: "destination", Location: "form", Required: true, Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/quotes/{quote}/accept", OperationID: "PostQuotesQuoteAccept", Params: []ParameterValidation{
		{Name: "quote", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/quotes/{quote}/cancel", OperationID: "PostQuotesQuoteCancel", Params: []ParameterValidation{
		{Name: "quote", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/quotes/{quote}/computed_upfront_line_items", OperationID: "GetQuotesQuoteComputedUpfrontLineItems", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "quote", Location: "path", Required: true, Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/quotes/{quote}/finalize", OperationID: "PostQuotesQuoteFinalize", Params: []ParameterValidation{
		{Name: "quote", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "expires_at", Location: "form", Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/quotes/{quote}/line_items", OperationID: "GetQuotesQuoteLineItems", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "quote", Location: "path", Required: true, Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/quotes/{quote}/pdf", OperationID: "GetQuotesQuotePdf", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "quote", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/radar/early_fraud_warnings", OperationID: "GetRadarEarlyFraudWarnings", Params: []ParameterValidation{
		{Name: "charge", Location: "query", Type: "string"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "payment_intent", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/radar/early_fraud_warnings/{early_fraud_warning}", OperationID: "GetRadarEarlyFraudWarningsEarlyFraudWarning", Params: []ParameterValidation{
		{Name: "early_fraud_warning", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/radar/payment_evaluations", OperationID: "PostRadarPaymentEvaluations", Params: []ParameterValidation{
		{Name: "client_device_metadata_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "radar_session", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "customer_details", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "customer", Location: "form", Type: "string"},
			{Name: "customer_account", Location: "form", Type: "string"},
			{Name: "email", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "payment_details", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "amount", Location: "form", Required: true, Type: "integer"},
			{Name: "currency", Location: "form", Required: true, Type: "string"},
			{Name: "description", Location: "form", Type: "string"},
			{Name: "money_movement_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "card", Location: "form", Type: "object"},
				{Name: "money_movement_type", Location: "form", Required: true, Type: "string", Enum: []string{"card"}},
			}},
			{Name: "payment_method_details", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "billing_details", Location: "form", Type: "object"},
				{Name: "payment_method", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "shipping_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object"},
				{Name: "name", Location: "form", Type: "string"},
				{Name: "phone", Location: "form", Type: "string"},
			}},
			{Name: "statement_descriptor", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/radar/value_list_items", OperationID: "GetRadarValueListItems", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "value", Location: "query", Type: "string"},
		{Name: "value_list", Location: "query", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/radar/value_list_items", OperationID: "PostRadarValueListItems", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "value", Location: "form", Required: true, Type: "string"},
		{Name: "value_list", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/radar/value_list_items/{item}", OperationID: "DeleteRadarValueListItemsItem", Params: []ParameterValidation{
		{Name: "item", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/radar/value_list_items/{item}", OperationID: "GetRadarValueListItemsItem", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "item", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/radar/value_lists", OperationID: "GetRadarValueLists", Params: []ParameterValidation{
		{Name: "alias", Location: "query", Type: "string"},
		{Name: "contains", Location: "query", Type: "string"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/radar/value_lists", OperationID: "PostRadarValueLists", Params: []ParameterValidation{
		{Name: "alias", Location: "form", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "item_type", Location: "form", Type: "string", Enum: []string{"account", "card_bin", "card_fingerprint", "case_sensitive_string", "country", "crypto_fingerprint", "customer_id", "email", "ip_address", "sepa_debit_fingerprint", "string", "us_bank_account_fingerprint"}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "name", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/radar/value_lists/{value_list}", OperationID: "DeleteRadarValueListsValueList", Params: []ParameterValidation{
		{Name: "value_list", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/radar/value_lists/{value_list}", OperationID: "GetRadarValueListsValueList", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "value_list", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/radar/value_lists/{value_list}", OperationID: "PostRadarValueListsValueList", Params: []ParameterValidation{
		{Name: "value_list", Location: "path", Required: true, Type: "string"},
		{Name: "alias", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "name", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/refunds", OperationID: "GetRefunds", Params: []ParameterValidation{
		{Name: "charge", Location: "query", Type: "string"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "payment_intent", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/refunds", OperationID: "PostRefunds", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "charge", Location: "form", Type: "string"},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "instructions_email", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "origin", Location: "form", Type: "string", Enum: []string{"customer_balance"}},
		{Name: "payment_intent", Location: "form", Type: "string"},
		{Name: "reason", Location: "form", Type: "string", Enum: []string{"duplicate", "fraudulent", "requested_by_customer"}},
		{Name: "refund_application_fee", Location: "form", Type: "boolean"},
		{Name: "reverse_transfer", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/refunds/{refund}", OperationID: "GetRefundsRefund", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "refund", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/refunds/{refund}", OperationID: "PostRefundsRefund", Params: []ParameterValidation{
		{Name: "refund", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/refunds/{refund}/cancel", OperationID: "PostRefundsRefundCancel", Params: []ParameterValidation{
		{Name: "refund", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/reporting/report_runs", OperationID: "GetReportingReportRuns", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/reporting/report_runs", OperationID: "PostReportingReportRuns", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "parameters", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "columns", Location: "form", Type: "array"},
			{Name: "connected_account", Location: "form", Type: "string"},
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "interval_end", Location: "form", Type: "integer"},
			{Name: "interval_start", Location: "form", Type: "integer"},
			{Name: "payout", Location: "form", Type: "string"},
			{Name: "reporting_category", Location: "form", Type: "string", Enum: []string{"advance", "advance_funding", "anticipation_repayment", "charge", "charge_failure", "climate_order_purchase", "climate_order_refund", "connect_collection_transfer", "connect_reserved_funds", "contribution", "dispute", "dispute_reversal", "fee", "financing_paydown", "financing_paydown_reversal", "financing_payout", "financing_payout_reversal", "issuing_authorization_hold", "issuing_authorization_release", "issuing_dispute", "issuing_transaction", "network_cost", "other_adjustment", "partial_capture_reversal", "payout", "payout_reversal", "platform_earning", "platform_earning_refund", "refund", "refund_failure", "risk_reserved_funds", "tax", "topup", "topup_reversal", "transfer", "transfer_reversal", "unreconciled_customer_funds"}},
			{Name: "timezone", Location: "form", Type: "string"},
		}},
		{Name: "report_type", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/reporting/report_runs/{report_run}", OperationID: "GetReportingReportRunsReportRun", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "report_run", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/reporting/report_types", OperationID: "GetReportingReportTypes", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/reporting/report_types/{report_type}", OperationID: "GetReportingReportTypesReportType", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "report_type", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/reviews", OperationID: "GetReviews", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/reviews/{review}", OperationID: "GetReviewsReview", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "review", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/reviews/{review}/approve", OperationID: "PostReviewsReviewApprove", Params: []ParameterValidation{
		{Name: "review", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/setup_attempts", OperationID: "GetSetupAttempts", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "setup_intent", Location: "query", Required: true, Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/setup_intents", OperationID: "GetSetupIntents", Params: []ParameterValidation{
		{Name: "attach_to_self", Location: "query", Type: "boolean"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "payment_method", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/setup_intents", OperationID: "PostSetupIntents", Params: []ParameterValidation{
		{Name: "attach_to_self", Location: "form", Type: "boolean"},
		{Name: "automatic_payment_methods", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "allow_redirects", Location: "form", Type: "string", Enum: []string{"always", "never"}},
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "confirm", Location: "form", Type: "boolean"},
		{Name: "confirmation_token", Location: "form", Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "excluded_payment_method_types", Location: "form", Type: "array"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "flow_directions", Location: "form", Type: "array"},
		{Name: "mandate_data", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "customer_acceptance", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "accepted_at", Location: "form", Type: "integer"},
				{Name: "offline", Location: "form", Type: "object"},
				{Name: "online", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"offline", "online"}},
			}},
		}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "on_behalf_of", Location: "form", Type: "string"},
		{Name: "payment_method", Location: "form", Type: "string"},
		{Name: "payment_method_configuration", Location: "form", Type: "string"},
		{Name: "payment_method_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "institution_number", Location: "form", Required: true, Type: "string"},
				{Name: "transit_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "affirm", Location: "form", Type: "object"},
			{Name: "afterpay_clearpay", Location: "form", Type: "object"},
			{Name: "alipay", Location: "form", Type: "object"},
			{Name: "allow_redisplay", Location: "form", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
			{Name: "alma", Location: "form", Type: "object"},
			{Name: "amazon_pay", Location: "form", Type: "object"},
			{Name: "au_becs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bsb_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "sort_code", Location: "form", Type: "string"},
			}},
			{Name: "bancontact", Location: "form", Type: "object"},
			{Name: "billie", Location: "form", Type: "object"},
			{Name: "billing_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "email", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "name", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "tax_id", Location: "form", Type: "string"},
			}},
			{Name: "blik", Location: "form", Type: "object"},
			{Name: "boleto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tax_id", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "cashapp", Location: "form", Type: "object"},
			{Name: "crypto", Location: "form", Type: "object"},
			{Name: "customer_balance", Location: "form", Type: "object"},
			{Name: "eps", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"arzte_und_apotheker_bank", "austrian_anadi_bank_ag", "bank_austria", "bankhaus_carl_spangler", "bankhaus_schelhammer_und_schattera_ag", "bawag_psk_ag", "bks_bank_ag", "brull_kallmus_bank_ag", "btv_vier_lander_bank", "capital_bank_grawe_gruppe_ag", "deutsche_bank_ag", "dolomitenbank", "easybank_ag", "erste_bank_und_sparkassen", "hypo_alpeadriabank_international_ag", "hypo_bank_burgenland_aktiengesellschaft", "hypo_noe_lb_fur_niederosterreich_u_wien", "hypo_oberosterreich_salzburg_steiermark", "hypo_tirol_bank_ag", "hypo_vorarlberg_bank_ag", "marchfelder_bank", "oberbank_ag", "raiffeisen_bankengruppe_osterreich", "schoellerbank_ag", "sparda_bank_wien", "volksbank_gruppe", "volkskreditbank_ag", "vr_bank_braunau"}},
			}},
			{Name: "fpx", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Required: true, Type: "string", Enum: []string{"affin_bank", "agrobank", "alliance_bank", "ambank", "bank_islam", "bank_muamalat", "bank_of_china", "bank_rakyat", "bsn", "cimb", "deutsche_bank", "hong_leong_bank", "hsbc", "kfh", "maybank2e", "maybank2u", "ocbc", "pb_enterprise", "public_bank", "rhb", "standard_chartered", "uob"}},
			}},
			{Name: "giropay", Location: "form", Type: "object"},
			{Name: "grabpay", Location: "form", Type: "object"},
			{Name: "ideal", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"abn_amro", "adyen", "asn_bank", "bunq", "buut", "finom", "handelsbanken", "ing", "knab", "mollie", "moneyou", "n26", "nn", "rabobank", "regiobank", "revolut", "sns_bank", "triodos_bank", "van_lanschot", "yoursafe"}},
			}},
			{Name: "interac_present", Location: "form", Type: "object"},
			{Name: "kakao_pay", Location: "form", Type: "object"},
			{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "dob", Location: "form", Type: "object"},
			}},
			{Name: "konbini", Location: "form", Type: "object"},
			{Name: "kr_card", Location: "form", Type: "object"},
			{Name: "link", Location: "form", Type: "object"},
			{Name: "mb_way", Location: "form", Type: "object"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "mobilepay", Location: "form", Type: "object"},
			{Name: "multibanco", Location: "form", Type: "object"},
			{Name: "naver_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "funding", Location: "form", Type: "string", Enum: []string{"card", "points"}},
			}},
			{Name: "nz_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_name", Location: "form", Type: "string"},
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bank_code", Location: "form", Required: true, Type: "string"},
				{Name: "branch_code", Location: "form", Required: true, Type: "string"},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "suffix", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "oxxo", Location: "form", Type: "object"},
			{Name: "p24", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"alior_bank", "bank_millennium", "bank_nowy_bfg_sa", "bank_pekao_sa", "banki_spbdzielcze", "blik", "bnp_paribas", "boz", "citi_handlowy", "credit_agricole", "envelobank", "etransfer_pocztowy24", "getin_bank", "ideabank", "ing", "inteligo", "mbank_mtransfer", "nest_przelew", "noble_pay", "pbac_z_ipko", "plus_bank", "santander_przelew24", "tmobile_usbugi_bankowe", "toyota_bank", "velobank", "volkswagen_bank"}},
			}},
			{Name: "pay_by_bank", Location: "form", Type: "object"},
			{Name: "payco", Location: "form", Type: "object"},
			{Name: "paynow", Location: "form", Type: "object"},
			{Name: "paypal", Location: "form", Type: "object"},
			{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "bsb_number", Location: "form", Type: "string"},
				{Name: "pay_id", Location: "form", Type: "string"},
			}},
			{Name: "pix", Location: "form", Type: "object"},
			{Name: "promptpay", Location: "form", Type: "object"},
			{Name: "radar_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "session", Location: "form", Type: "string"},
			}},
			{Name: "revolut_pay", Location: "form", Type: "object"},
			{Name: "samsung_pay", Location: "form", Type: "object"},
			{Name: "satispay", Location: "form", Type: "object"},
			{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "iban", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "sofort", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "country", Location: "form", Required: true, Type: "string", Enum: []string{"AT", "BE", "DE", "ES", "IT", "NL"}},
			}},
			{Name: "sunbit", Location: "form", Type: "object"},
			{Name: "swish", Location: "form", Type: "object"},
			{Name: "twint", Location: "form", Type: "object"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"acss_debit", "affirm", "afterpay_clearpay", "alipay", "alma", "amazon_pay", "au_becs_debit", "bacs_debit", "bancontact", "billie", "blik", "boleto", "cashapp", "crypto", "customer_balance", "eps", "fpx", "giropay", "grabpay", "ideal", "kakao_pay", "klarna", "konbini", "kr_card", "link", "mb_way", "mobilepay", "multibanco", "naver_pay", "nz_bank_account", "oxxo", "p24", "pay_by_bank", "payco", "paynow", "paypal", "payto", "pix", "promptpay", "revolut_pay", "samsung_pay", "satispay", "sepa_debit", "sofort", "sunbit", "swish", "twint", "upi", "us_bank_account", "wechat_pay", "zip"}},
			{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "savings"}},
				{Name: "financial_connections_account", Location: "form", Type: "string"},
				{Name: "routing_number", Location: "form", Type: "string"},
			}},
			{Name: "wechat_pay", Location: "form", Type: "object"},
			{Name: "zip", Location: "form", Type: "object"},
		}},
		{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "currency", Location: "form", Type: "string", Enum: []string{"cad", "usd"}},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant", "microdeposits"}},
			}},
			{Name: "amazon_pay", Location: "form", Type: "object"},
			{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "network", Location: "form", Type: "string", Enum: []string{"amex", "cartes_bancaires", "diners", "discover", "eftpos_au", "girocard", "interac", "jcb", "link", "mastercard", "unionpay", "unknown", "visa"}},
				{Name: "request_three_d_secure", Location: "form", Type: "string", Enum: []string{"any", "automatic", "challenge"}},
				{Name: "three_d_secure", Location: "form", Type: "object"},
			}},
			{Name: "card_present", Location: "form", Type: "object"},
			{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "currency", Location: "form", Type: "string"},
				{Name: "on_demand", Location: "form", Type: "object"},
				{Name: "preferred_locale", Location: "form", Type: "string", Enum: []string{"cs-CZ", "da-DK", "de-AT", "de-CH", "de-DE", "el-GR", "en-AT", "en-AU", "en-BE", "en-CA", "en-CH", "en-CZ", "en-DE", "en-DK", "en-ES", "en-FI", "en-FR", "en-GB", "en-GR", "en-IE", "en-IT", "en-NL", "en-NO", "en-NZ", "en-PL", "en-PT", "en-RO", "en-SE", "en-US", "es-ES", "es-US", "fi-FI", "fr-BE", "fr-CA", "fr-CH", "fr-FR", "it-CH", "it-IT", "nb-NO", "nl-BE", "nl-NL", "pl-PL", "pt-PT", "ro-RO", "sv-FI", "sv-SE"}},
				{Name: "subscriptions", Location: "form", Enum: []string{""}},
			}},
			{Name: "link", Location: "form", Type: "object"},
			{Name: "paypal", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "billing_agreement_id", Location: "form", Type: "string"},
			}},
			{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "pix", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "financial_connections", Location: "form", Type: "object"},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "networks", Location: "form", Type: "object"},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant", "microdeposits"}},
			}},
		}},
		{Name: "payment_method_types", Location: "form", Type: "array"},
		{Name: "return_url", Location: "form", Type: "string"},
		{Name: "single_use", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount", Location: "form", Required: true, Type: "integer"},
			{Name: "currency", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "usage", Location: "form", Type: "string", Enum: []string{"off_session", "on_session"}},
		{Name: "use_stripe_sdk", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/setup_intents/{intent}", OperationID: "GetSetupIntentsIntent", Params: []ParameterValidation{
		{Name: "client_secret", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "intent", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/setup_intents/{intent}", OperationID: "PostSetupIntentsIntent", Params: []ParameterValidation{
		{Name: "intent", Location: "path", Required: true, Type: "string"},
		{Name: "attach_to_self", Location: "form", Type: "boolean"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "excluded_payment_method_types", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "flow_directions", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "payment_method", Location: "form", Type: "string"},
		{Name: "payment_method_configuration", Location: "form", Type: "string"},
		{Name: "payment_method_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "institution_number", Location: "form", Required: true, Type: "string"},
				{Name: "transit_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "affirm", Location: "form", Type: "object"},
			{Name: "afterpay_clearpay", Location: "form", Type: "object"},
			{Name: "alipay", Location: "form", Type: "object"},
			{Name: "allow_redisplay", Location: "form", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
			{Name: "alma", Location: "form", Type: "object"},
			{Name: "amazon_pay", Location: "form", Type: "object"},
			{Name: "au_becs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bsb_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "sort_code", Location: "form", Type: "string"},
			}},
			{Name: "bancontact", Location: "form", Type: "object"},
			{Name: "billie", Location: "form", Type: "object"},
			{Name: "billing_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "email", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "name", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "tax_id", Location: "form", Type: "string"},
			}},
			{Name: "blik", Location: "form", Type: "object"},
			{Name: "boleto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tax_id", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "cashapp", Location: "form", Type: "object"},
			{Name: "crypto", Location: "form", Type: "object"},
			{Name: "customer_balance", Location: "form", Type: "object"},
			{Name: "eps", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"arzte_und_apotheker_bank", "austrian_anadi_bank_ag", "bank_austria", "bankhaus_carl_spangler", "bankhaus_schelhammer_und_schattera_ag", "bawag_psk_ag", "bks_bank_ag", "brull_kallmus_bank_ag", "btv_vier_lander_bank", "capital_bank_grawe_gruppe_ag", "deutsche_bank_ag", "dolomitenbank", "easybank_ag", "erste_bank_und_sparkassen", "hypo_alpeadriabank_international_ag", "hypo_bank_burgenland_aktiengesellschaft", "hypo_noe_lb_fur_niederosterreich_u_wien", "hypo_oberosterreich_salzburg_steiermark", "hypo_tirol_bank_ag", "hypo_vorarlberg_bank_ag", "marchfelder_bank", "oberbank_ag", "raiffeisen_bankengruppe_osterreich", "schoellerbank_ag", "sparda_bank_wien", "volksbank_gruppe", "volkskreditbank_ag", "vr_bank_braunau"}},
			}},
			{Name: "fpx", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Required: true, Type: "string", Enum: []string{"affin_bank", "agrobank", "alliance_bank", "ambank", "bank_islam", "bank_muamalat", "bank_of_china", "bank_rakyat", "bsn", "cimb", "deutsche_bank", "hong_leong_bank", "hsbc", "kfh", "maybank2e", "maybank2u", "ocbc", "pb_enterprise", "public_bank", "rhb", "standard_chartered", "uob"}},
			}},
			{Name: "giropay", Location: "form", Type: "object"},
			{Name: "grabpay", Location: "form", Type: "object"},
			{Name: "ideal", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"abn_amro", "adyen", "asn_bank", "bunq", "buut", "finom", "handelsbanken", "ing", "knab", "mollie", "moneyou", "n26", "nn", "rabobank", "regiobank", "revolut", "sns_bank", "triodos_bank", "van_lanschot", "yoursafe"}},
			}},
			{Name: "interac_present", Location: "form", Type: "object"},
			{Name: "kakao_pay", Location: "form", Type: "object"},
			{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "dob", Location: "form", Type: "object"},
			}},
			{Name: "konbini", Location: "form", Type: "object"},
			{Name: "kr_card", Location: "form", Type: "object"},
			{Name: "link", Location: "form", Type: "object"},
			{Name: "mb_way", Location: "form", Type: "object"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "mobilepay", Location: "form", Type: "object"},
			{Name: "multibanco", Location: "form", Type: "object"},
			{Name: "naver_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "funding", Location: "form", Type: "string", Enum: []string{"card", "points"}},
			}},
			{Name: "nz_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_name", Location: "form", Type: "string"},
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bank_code", Location: "form", Required: true, Type: "string"},
				{Name: "branch_code", Location: "form", Required: true, Type: "string"},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "suffix", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "oxxo", Location: "form", Type: "object"},
			{Name: "p24", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"alior_bank", "bank_millennium", "bank_nowy_bfg_sa", "bank_pekao_sa", "banki_spbdzielcze", "blik", "bnp_paribas", "boz", "citi_handlowy", "credit_agricole", "envelobank", "etransfer_pocztowy24", "getin_bank", "ideabank", "ing", "inteligo", "mbank_mtransfer", "nest_przelew", "noble_pay", "pbac_z_ipko", "plus_bank", "santander_przelew24", "tmobile_usbugi_bankowe", "toyota_bank", "velobank", "volkswagen_bank"}},
			}},
			{Name: "pay_by_bank", Location: "form", Type: "object"},
			{Name: "payco", Location: "form", Type: "object"},
			{Name: "paynow", Location: "form", Type: "object"},
			{Name: "paypal", Location: "form", Type: "object"},
			{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "bsb_number", Location: "form", Type: "string"},
				{Name: "pay_id", Location: "form", Type: "string"},
			}},
			{Name: "pix", Location: "form", Type: "object"},
			{Name: "promptpay", Location: "form", Type: "object"},
			{Name: "radar_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "session", Location: "form", Type: "string"},
			}},
			{Name: "revolut_pay", Location: "form", Type: "object"},
			{Name: "samsung_pay", Location: "form", Type: "object"},
			{Name: "satispay", Location: "form", Type: "object"},
			{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "iban", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "sofort", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "country", Location: "form", Required: true, Type: "string", Enum: []string{"AT", "BE", "DE", "ES", "IT", "NL"}},
			}},
			{Name: "sunbit", Location: "form", Type: "object"},
			{Name: "swish", Location: "form", Type: "object"},
			{Name: "twint", Location: "form", Type: "object"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"acss_debit", "affirm", "afterpay_clearpay", "alipay", "alma", "amazon_pay", "au_becs_debit", "bacs_debit", "bancontact", "billie", "blik", "boleto", "cashapp", "crypto", "customer_balance", "eps", "fpx", "giropay", "grabpay", "ideal", "kakao_pay", "klarna", "konbini", "kr_card", "link", "mb_way", "mobilepay", "multibanco", "naver_pay", "nz_bank_account", "oxxo", "p24", "pay_by_bank", "payco", "paynow", "paypal", "payto", "pix", "promptpay", "revolut_pay", "samsung_pay", "satispay", "sepa_debit", "sofort", "sunbit", "swish", "twint", "upi", "us_bank_account", "wechat_pay", "zip"}},
			{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "savings"}},
				{Name: "financial_connections_account", Location: "form", Type: "string"},
				{Name: "routing_number", Location: "form", Type: "string"},
			}},
			{Name: "wechat_pay", Location: "form", Type: "object"},
			{Name: "zip", Location: "form", Type: "object"},
		}},
		{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "currency", Location: "form", Type: "string", Enum: []string{"cad", "usd"}},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant", "microdeposits"}},
			}},
			{Name: "amazon_pay", Location: "form", Type: "object"},
			{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "network", Location: "form", Type: "string", Enum: []string{"amex", "cartes_bancaires", "diners", "discover", "eftpos_au", "girocard", "interac", "jcb", "link", "mastercard", "unionpay", "unknown", "visa"}},
				{Name: "request_three_d_secure", Location: "form", Type: "string", Enum: []string{"any", "automatic", "challenge"}},
				{Name: "three_d_secure", Location: "form", Type: "object"},
			}},
			{Name: "card_present", Location: "form", Type: "object"},
			{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "currency", Location: "form", Type: "string"},
				{Name: "on_demand", Location: "form", Type: "object"},
				{Name: "preferred_locale", Location: "form", Type: "string", Enum: []string{"cs-CZ", "da-DK", "de-AT", "de-CH", "de-DE", "el-GR", "en-AT", "en-AU", "en-BE", "en-CA", "en-CH", "en-CZ", "en-DE", "en-DK", "en-ES", "en-FI", "en-FR", "en-GB", "en-GR", "en-IE", "en-IT", "en-NL", "en-NO", "en-NZ", "en-PL", "en-PT", "en-RO", "en-SE", "en-US", "es-ES", "es-US", "fi-FI", "fr-BE", "fr-CA", "fr-CH", "fr-FR", "it-CH", "it-IT", "nb-NO", "nl-BE", "nl-NL", "pl-PL", "pt-PT", "ro-RO", "sv-FI", "sv-SE"}},
				{Name: "subscriptions", Location: "form", Enum: []string{""}},
			}},
			{Name: "link", Location: "form", Type: "object"},
			{Name: "paypal", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "billing_agreement_id", Location: "form", Type: "string"},
			}},
			{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "pix", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "financial_connections", Location: "form", Type: "object"},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "networks", Location: "form", Type: "object"},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant", "microdeposits"}},
			}},
		}},
		{Name: "payment_method_types", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/setup_intents/{intent}/cancel", OperationID: "PostSetupIntentsIntentCancel", Params: []ParameterValidation{
		{Name: "intent", Location: "path", Required: true, Type: "string"},
		{Name: "cancellation_reason", Location: "form", Type: "string", Enum: []string{"abandoned", "duplicate", "requested_by_customer"}},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/setup_intents/{intent}/confirm", OperationID: "PostSetupIntentsIntentConfirm", Params: []ParameterValidation{
		{Name: "intent", Location: "path", Required: true, Type: "string"},
		{Name: "client_secret", Location: "form", Type: "string"},
		{Name: "confirmation_token", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "mandate_data", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "customer_acceptance", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "accepted_at", Location: "form", Type: "integer"},
				{Name: "offline", Location: "form", Type: "object"},
				{Name: "online", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"offline", "online"}},
			}},
		}},
		{Name: "payment_method", Location: "form", Type: "string"},
		{Name: "payment_method_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "institution_number", Location: "form", Required: true, Type: "string"},
				{Name: "transit_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "affirm", Location: "form", Type: "object"},
			{Name: "afterpay_clearpay", Location: "form", Type: "object"},
			{Name: "alipay", Location: "form", Type: "object"},
			{Name: "allow_redisplay", Location: "form", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
			{Name: "alma", Location: "form", Type: "object"},
			{Name: "amazon_pay", Location: "form", Type: "object"},
			{Name: "au_becs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bsb_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "sort_code", Location: "form", Type: "string"},
			}},
			{Name: "bancontact", Location: "form", Type: "object"},
			{Name: "billie", Location: "form", Type: "object"},
			{Name: "billing_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "email", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "name", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "tax_id", Location: "form", Type: "string"},
			}},
			{Name: "blik", Location: "form", Type: "object"},
			{Name: "boleto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tax_id", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "cashapp", Location: "form", Type: "object"},
			{Name: "crypto", Location: "form", Type: "object"},
			{Name: "customer_balance", Location: "form", Type: "object"},
			{Name: "eps", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"arzte_und_apotheker_bank", "austrian_anadi_bank_ag", "bank_austria", "bankhaus_carl_spangler", "bankhaus_schelhammer_und_schattera_ag", "bawag_psk_ag", "bks_bank_ag", "brull_kallmus_bank_ag", "btv_vier_lander_bank", "capital_bank_grawe_gruppe_ag", "deutsche_bank_ag", "dolomitenbank", "easybank_ag", "erste_bank_und_sparkassen", "hypo_alpeadriabank_international_ag", "hypo_bank_burgenland_aktiengesellschaft", "hypo_noe_lb_fur_niederosterreich_u_wien", "hypo_oberosterreich_salzburg_steiermark", "hypo_tirol_bank_ag", "hypo_vorarlberg_bank_ag", "marchfelder_bank", "oberbank_ag", "raiffeisen_bankengruppe_osterreich", "schoellerbank_ag", "sparda_bank_wien", "volksbank_gruppe", "volkskreditbank_ag", "vr_bank_braunau"}},
			}},
			{Name: "fpx", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Required: true, Type: "string", Enum: []string{"affin_bank", "agrobank", "alliance_bank", "ambank", "bank_islam", "bank_muamalat", "bank_of_china", "bank_rakyat", "bsn", "cimb", "deutsche_bank", "hong_leong_bank", "hsbc", "kfh", "maybank2e", "maybank2u", "ocbc", "pb_enterprise", "public_bank", "rhb", "standard_chartered", "uob"}},
			}},
			{Name: "giropay", Location: "form", Type: "object"},
			{Name: "grabpay", Location: "form", Type: "object"},
			{Name: "ideal", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"abn_amro", "adyen", "asn_bank", "bunq", "buut", "finom", "handelsbanken", "ing", "knab", "mollie", "moneyou", "n26", "nn", "rabobank", "regiobank", "revolut", "sns_bank", "triodos_bank", "van_lanschot", "yoursafe"}},
			}},
			{Name: "interac_present", Location: "form", Type: "object"},
			{Name: "kakao_pay", Location: "form", Type: "object"},
			{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "dob", Location: "form", Type: "object"},
			}},
			{Name: "konbini", Location: "form", Type: "object"},
			{Name: "kr_card", Location: "form", Type: "object"},
			{Name: "link", Location: "form", Type: "object"},
			{Name: "mb_way", Location: "form", Type: "object"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "mobilepay", Location: "form", Type: "object"},
			{Name: "multibanco", Location: "form", Type: "object"},
			{Name: "naver_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "funding", Location: "form", Type: "string", Enum: []string{"card", "points"}},
			}},
			{Name: "nz_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_name", Location: "form", Type: "string"},
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bank_code", Location: "form", Required: true, Type: "string"},
				{Name: "branch_code", Location: "form", Required: true, Type: "string"},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "suffix", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "oxxo", Location: "form", Type: "object"},
			{Name: "p24", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"alior_bank", "bank_millennium", "bank_nowy_bfg_sa", "bank_pekao_sa", "banki_spbdzielcze", "blik", "bnp_paribas", "boz", "citi_handlowy", "credit_agricole", "envelobank", "etransfer_pocztowy24", "getin_bank", "ideabank", "ing", "inteligo", "mbank_mtransfer", "nest_przelew", "noble_pay", "pbac_z_ipko", "plus_bank", "santander_przelew24", "tmobile_usbugi_bankowe", "toyota_bank", "velobank", "volkswagen_bank"}},
			}},
			{Name: "pay_by_bank", Location: "form", Type: "object"},
			{Name: "payco", Location: "form", Type: "object"},
			{Name: "paynow", Location: "form", Type: "object"},
			{Name: "paypal", Location: "form", Type: "object"},
			{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "bsb_number", Location: "form", Type: "string"},
				{Name: "pay_id", Location: "form", Type: "string"},
			}},
			{Name: "pix", Location: "form", Type: "object"},
			{Name: "promptpay", Location: "form", Type: "object"},
			{Name: "radar_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "session", Location: "form", Type: "string"},
			}},
			{Name: "revolut_pay", Location: "form", Type: "object"},
			{Name: "samsung_pay", Location: "form", Type: "object"},
			{Name: "satispay", Location: "form", Type: "object"},
			{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "iban", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "sofort", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "country", Location: "form", Required: true, Type: "string", Enum: []string{"AT", "BE", "DE", "ES", "IT", "NL"}},
			}},
			{Name: "sunbit", Location: "form", Type: "object"},
			{Name: "swish", Location: "form", Type: "object"},
			{Name: "twint", Location: "form", Type: "object"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"acss_debit", "affirm", "afterpay_clearpay", "alipay", "alma", "amazon_pay", "au_becs_debit", "bacs_debit", "bancontact", "billie", "blik", "boleto", "cashapp", "crypto", "customer_balance", "eps", "fpx", "giropay", "grabpay", "ideal", "kakao_pay", "klarna", "konbini", "kr_card", "link", "mb_way", "mobilepay", "multibanco", "naver_pay", "nz_bank_account", "oxxo", "p24", "pay_by_bank", "payco", "paynow", "paypal", "payto", "pix", "promptpay", "revolut_pay", "samsung_pay", "satispay", "sepa_debit", "sofort", "sunbit", "swish", "twint", "upi", "us_bank_account", "wechat_pay", "zip"}},
			{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "savings"}},
				{Name: "financial_connections_account", Location: "form", Type: "string"},
				{Name: "routing_number", Location: "form", Type: "string"},
			}},
			{Name: "wechat_pay", Location: "form", Type: "object"},
			{Name: "zip", Location: "form", Type: "object"},
		}},
		{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "currency", Location: "form", Type: "string", Enum: []string{"cad", "usd"}},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant", "microdeposits"}},
			}},
			{Name: "amazon_pay", Location: "form", Type: "object"},
			{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "network", Location: "form", Type: "string", Enum: []string{"amex", "cartes_bancaires", "diners", "discover", "eftpos_au", "girocard", "interac", "jcb", "link", "mastercard", "unionpay", "unknown", "visa"}},
				{Name: "request_three_d_secure", Location: "form", Type: "string", Enum: []string{"any", "automatic", "challenge"}},
				{Name: "three_d_secure", Location: "form", Type: "object"},
			}},
			{Name: "card_present", Location: "form", Type: "object"},
			{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "currency", Location: "form", Type: "string"},
				{Name: "on_demand", Location: "form", Type: "object"},
				{Name: "preferred_locale", Location: "form", Type: "string", Enum: []string{"cs-CZ", "da-DK", "de-AT", "de-CH", "de-DE", "el-GR", "en-AT", "en-AU", "en-BE", "en-CA", "en-CH", "en-CZ", "en-DE", "en-DK", "en-ES", "en-FI", "en-FR", "en-GB", "en-GR", "en-IE", "en-IT", "en-NL", "en-NO", "en-NZ", "en-PL", "en-PT", "en-RO", "en-SE", "en-US", "es-ES", "es-US", "fi-FI", "fr-BE", "fr-CA", "fr-CH", "fr-FR", "it-CH", "it-IT", "nb-NO", "nl-BE", "nl-NL", "pl-PL", "pt-PT", "ro-RO", "sv-FI", "sv-SE"}},
				{Name: "subscriptions", Location: "form", Enum: []string{""}},
			}},
			{Name: "link", Location: "form", Type: "object"},
			{Name: "paypal", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "billing_agreement_id", Location: "form", Type: "string"},
			}},
			{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "pix", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"", "none", "off_session", "on_session"}},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "financial_connections", Location: "form", Type: "object"},
				{Name: "mandate_options", Location: "form", Type: "object"},
				{Name: "networks", Location: "form", Type: "object"},
				{Name: "verification_method", Location: "form", Type: "string", Enum: []string{"automatic", "instant", "microdeposits"}},
			}},
		}},
		{Name: "return_url", Location: "form", Type: "string"},
		{Name: "use_stripe_sdk", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/setup_intents/{intent}/verify_microdeposits", OperationID: "PostSetupIntentsIntentVerifyMicrodeposits", Params: []ParameterValidation{
		{Name: "intent", Location: "path", Required: true, Type: "string"},
		{Name: "amounts", Location: "form", Type: "array"},
		{Name: "client_secret", Location: "form", Type: "string"},
		{Name: "descriptor_code", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/shipping_rates", OperationID: "GetShippingRates", Params: []ParameterValidation{
		{Name: "active", Location: "query", Type: "boolean"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "currency", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/shipping_rates", OperationID: "PostShippingRates", Params: []ParameterValidation{
		{Name: "delivery_estimate", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "maximum", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "unit", Location: "form", Required: true, Type: "string", Enum: []string{"business_day", "day", "hour", "month", "week"}},
				{Name: "value", Location: "form", Required: true, Type: "integer"},
			}},
			{Name: "minimum", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "unit", Location: "form", Required: true, Type: "string", Enum: []string{"business_day", "day", "hour", "month", "week"}},
				{Name: "value", Location: "form", Required: true, Type: "integer"},
			}},
		}},
		{Name: "display_name", Location: "form", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "fixed_amount", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount", Location: "form", Required: true, Type: "integer"},
			{Name: "currency", Location: "form", Required: true, Type: "string"},
			{Name: "currency_options", Location: "form", Type: "object", AdditionalProperties: true},
		}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
		{Name: "tax_code", Location: "form", Type: "string"},
		{Name: "type", Location: "form", Type: "string", Enum: []string{"fixed_amount"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/shipping_rates/{shipping_rate_token}", OperationID: "GetShippingRatesShippingRateToken", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "shipping_rate_token", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/shipping_rates/{shipping_rate_token}", OperationID: "PostShippingRatesShippingRateToken", Params: []ParameterValidation{
		{Name: "shipping_rate_token", Location: "path", Required: true, Type: "string"},
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "fixed_amount", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "currency_options", Location: "form", Type: "object", AdditionalProperties: true},
		}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/sigma/saved_queries/{id}", OperationID: "PostSigmaSavedQueriesId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "sql", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/sigma/scheduled_query_runs", OperationID: "GetSigmaScheduledQueryRuns", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/sigma/scheduled_query_runs/{scheduled_query_run}", OperationID: "GetSigmaScheduledQueryRunsScheduledQueryRun", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "scheduled_query_run", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/sources", OperationID: "PostSources", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "flow", Location: "form", Type: "string", Enum: []string{"code_verification", "none", "receiver", "redirect"}},
		{Name: "mandate", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acceptance", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Type: "integer"},
				{Name: "ip", Location: "form", Type: "string"},
				{Name: "offline", Location: "form", Type: "object"},
				{Name: "online", Location: "form", Type: "object"},
				{Name: "status", Location: "form", Required: true, Type: "string", Enum: []string{"accepted", "pending", "refused", "revoked"}},
				{Name: "type", Location: "form", Type: "string", Enum: []string{"offline", "online"}},
				{Name: "user_agent", Location: "form", Type: "string"},
			}},
			{Name: "amount", Location: "form", Enum: []string{""}},
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "interval", Location: "form", Type: "string", Enum: []string{"one_time", "scheduled", "variable"}},
			{Name: "notification_method", Location: "form", Type: "string", Enum: []string{"deprecated_none", "email", "manual", "none", "stripe_email"}},
		}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "original_source", Location: "form", Type: "string"},
		{Name: "owner", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "email", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
		{Name: "receiver", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "refund_attributes_method", Location: "form", Type: "string", Enum: []string{"email", "manual", "none"}},
		}},
		{Name: "redirect", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "return_url", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "source_order", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "items", Location: "form", Type: "array"},
			{Name: "shipping", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Required: true, Type: "object"},
				{Name: "carrier", Location: "form", Type: "string"},
				{Name: "name", Location: "form", Type: "string"},
				{Name: "phone", Location: "form", Type: "string"},
				{Name: "tracking_number", Location: "form", Type: "string"},
			}},
		}},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
		{Name: "token", Location: "form", Type: "string"},
		{Name: "type", Location: "form", Type: "string"},
		{Name: "usage", Location: "form", Type: "string", Enum: []string{"reusable", "single_use"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/sources/{source}", OperationID: "GetSourcesSource", Params: []ParameterValidation{
		{Name: "client_secret", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "source", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/sources/{source}", OperationID: "PostSourcesSource", Params: []ParameterValidation{
		{Name: "source", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "mandate", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acceptance", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Type: "integer"},
				{Name: "ip", Location: "form", Type: "string"},
				{Name: "offline", Location: "form", Type: "object"},
				{Name: "online", Location: "form", Type: "object"},
				{Name: "status", Location: "form", Required: true, Type: "string", Enum: []string{"accepted", "pending", "refused", "revoked"}},
				{Name: "type", Location: "form", Type: "string", Enum: []string{"offline", "online"}},
				{Name: "user_agent", Location: "form", Type: "string"},
			}},
			{Name: "amount", Location: "form", Enum: []string{""}},
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "interval", Location: "form", Type: "string", Enum: []string{"one_time", "scheduled", "variable"}},
			{Name: "notification_method", Location: "form", Type: "string", Enum: []string{"deprecated_none", "email", "manual", "none", "stripe_email"}},
		}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "owner", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "email", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
		}},
		{Name: "source_order", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "items", Location: "form", Type: "array"},
			{Name: "shipping", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Required: true, Type: "object"},
				{Name: "carrier", Location: "form", Type: "string"},
				{Name: "name", Location: "form", Type: "string"},
				{Name: "phone", Location: "form", Type: "string"},
				{Name: "tracking_number", Location: "form", Type: "string"},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/sources/{source}/mandate_notifications/{mandate_notification}", OperationID: "GetSourcesSourceMandateNotificationsMandateNotification", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "mandate_notification", Location: "path", Required: true, Type: "string"},
		{Name: "source", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/sources/{source}/source_transactions", OperationID: "GetSourcesSourceSourceTransactions", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "source", Location: "path", Required: true, Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/sources/{source}/source_transactions/{source_transaction}", OperationID: "GetSourcesSourceSourceTransactionsSourceTransaction", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "source", Location: "path", Required: true, Type: "string"},
		{Name: "source_transaction", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/sources/{source}/verify", OperationID: "PostSourcesSourceVerify", Params: []ParameterValidation{
		{Name: "source", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "values", Location: "form", Required: true, Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/subscription_items", OperationID: "GetSubscriptionItems", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "subscription", Location: "query", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/subscription_items", OperationID: "PostSubscriptionItems", Params: []ParameterValidation{
		{Name: "billing_thresholds", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "usage_gte", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "payment_behavior", Location: "form", Type: "string", Enum: []string{"allow_incomplete", "default_incomplete", "error_if_incomplete", "pending_if_incomplete"}},
		{Name: "price", Location: "form", Type: "string"},
		{Name: "price_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "currency", Location: "form", Required: true, Type: "string"},
			{Name: "product", Location: "form", Required: true, Type: "string"},
			{Name: "recurring", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "interval", Location: "form", Required: true, Type: "string", Enum: []string{"day", "month", "week", "year"}},
				{Name: "interval_count", Location: "form", Type: "integer"},
			}},
			{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
			{Name: "unit_amount", Location: "form", Type: "integer"},
			{Name: "unit_amount_decimal", Location: "form", Type: "string"},
		}},
		{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
		{Name: "proration_date", Location: "form", Type: "integer"},
		{Name: "quantity", Location: "form", Type: "integer"},
		{Name: "subscription", Location: "form", Required: true, Type: "string"},
		{Name: "tax_rates", Location: "form", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/subscription_items/{item}", OperationID: "DeleteSubscriptionItemsItem", Params: []ParameterValidation{
		{Name: "item", Location: "path", Required: true, Type: "string"},
		{Name: "clear_usage", Location: "form", Type: "boolean"},
		{Name: "payment_behavior", Location: "form", Type: "string", Enum: []string{"allow_incomplete", "default_incomplete", "error_if_incomplete", "pending_if_incomplete"}},
		{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
		{Name: "proration_date", Location: "form", Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/subscription_items/{item}", OperationID: "GetSubscriptionItemsItem", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "item", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/subscription_items/{item}", OperationID: "PostSubscriptionItemsItem", Params: []ParameterValidation{
		{Name: "item", Location: "path", Required: true, Type: "string"},
		{Name: "billing_thresholds", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "usage_gte", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "off_session", Location: "form", Type: "boolean"},
		{Name: "payment_behavior", Location: "form", Type: "string", Enum: []string{"allow_incomplete", "default_incomplete", "error_if_incomplete", "pending_if_incomplete"}},
		{Name: "price", Location: "form", Type: "string"},
		{Name: "price_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "currency", Location: "form", Required: true, Type: "string"},
			{Name: "product", Location: "form", Required: true, Type: "string"},
			{Name: "recurring", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "interval", Location: "form", Required: true, Type: "string", Enum: []string{"day", "month", "week", "year"}},
				{Name: "interval_count", Location: "form", Type: "integer"},
			}},
			{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "unspecified"}},
			{Name: "unit_amount", Location: "form", Type: "integer"},
			{Name: "unit_amount_decimal", Location: "form", Type: "string"},
		}},
		{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
		{Name: "proration_date", Location: "form", Type: "integer"},
		{Name: "quantity", Location: "form", Type: "integer"},
		{Name: "tax_rates", Location: "form", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/subscription_schedules", OperationID: "GetSubscriptionSchedules", Params: []ParameterValidation{
		{Name: "canceled_at", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "completed_at", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "released_at", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "scheduled", Location: "query", Type: "boolean"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/subscription_schedules", OperationID: "PostSubscriptionSchedules", Params: []ParameterValidation{
		{Name: "billing_mode", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "flexible", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "proration_discounts", Location: "form", Type: "string", Enum: []string{"included", "itemized"}},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"classic", "flexible"}},
		}},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "default_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "application_fee_percent", Location: "form", Type: "number"},
			{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "liability", Location: "form", Type: "object"},
			}},
			{Name: "billing_cycle_anchor", Location: "form", Type: "string", Enum: []string{"automatic", "phase_start"}},
			{Name: "billing_thresholds", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount_gte", Location: "form", Type: "integer"},
				{Name: "reset_billing_cycle_anchor", Location: "form", Type: "boolean"},
			}},
			{Name: "collection_method", Location: "form", Type: "string", Enum: []string{"charge_automatically", "send_invoice"}},
			{Name: "default_payment_method", Location: "form", Type: "string"},
			{Name: "description", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "invoice_settings", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_tax_ids", Location: "form", Enum: []string{""}},
				{Name: "days_until_due", Location: "form", Type: "integer"},
				{Name: "issuer", Location: "form", Type: "object"},
			}},
			{Name: "on_behalf_of", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "transfer_data", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount_percent", Location: "form", Type: "number"},
				{Name: "destination", Location: "form", Required: true, Type: "string"},
			}},
		}},
		{Name: "end_behavior", Location: "form", Type: "string", Enum: []string{"cancel", "none", "release", "renew"}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "from_subscription", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "phases", Location: "form", Type: "array"},
		{Name: "start_date", Location: "form", Enum: []string{"now"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/subscription_schedules/{schedule}", OperationID: "GetSubscriptionSchedulesSchedule", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "schedule", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/subscription_schedules/{schedule}", OperationID: "PostSubscriptionSchedulesSchedule", Params: []ParameterValidation{
		{Name: "schedule", Location: "path", Required: true, Type: "string"},
		{Name: "default_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "application_fee_percent", Location: "form", Type: "number"},
			{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
				{Name: "liability", Location: "form", Type: "object"},
			}},
			{Name: "billing_cycle_anchor", Location: "form", Type: "string", Enum: []string{"automatic", "phase_start"}},
			{Name: "billing_thresholds", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount_gte", Location: "form", Type: "integer"},
				{Name: "reset_billing_cycle_anchor", Location: "form", Type: "boolean"},
			}},
			{Name: "collection_method", Location: "form", Type: "string", Enum: []string{"charge_automatically", "send_invoice"}},
			{Name: "default_payment_method", Location: "form", Type: "string"},
			{Name: "description", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "invoice_settings", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_tax_ids", Location: "form", Enum: []string{""}},
				{Name: "days_until_due", Location: "form", Type: "integer"},
				{Name: "issuer", Location: "form", Type: "object"},
			}},
			{Name: "on_behalf_of", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "transfer_data", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "amount_percent", Location: "form", Type: "number"},
				{Name: "destination", Location: "form", Required: true, Type: "string"},
			}},
		}},
		{Name: "end_behavior", Location: "form", Type: "string", Enum: []string{"cancel", "none", "release", "renew"}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "phases", Location: "form", Type: "array"},
		{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/subscription_schedules/{schedule}/cancel", OperationID: "PostSubscriptionSchedulesScheduleCancel", Params: []ParameterValidation{
		{Name: "schedule", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "invoice_now", Location: "form", Type: "boolean"},
		{Name: "prorate", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/subscription_schedules/{schedule}/release", OperationID: "PostSubscriptionSchedulesScheduleRelease", Params: []ParameterValidation{
		{Name: "schedule", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "preserve_cancel_date", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/subscriptions", OperationID: "GetSubscriptions", Params: []ParameterValidation{
		{Name: "automatic_tax", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "query", Required: true, Type: "boolean"},
		}},
		{Name: "collection_method", Location: "query", Type: "string", Enum: []string{"charge_automatically", "send_invoice"}},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "current_period_end", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "current_period_start", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "customer_account", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "price", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"active", "all", "canceled", "ended", "incomplete", "incomplete_expired", "past_due", "paused", "trialing", "unpaid"}},
		{Name: "test_clock", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/subscriptions", OperationID: "PostSubscriptions", Params: []ParameterValidation{
		{Name: "add_invoice_items", Location: "form", Type: "array"},
		{Name: "application_fee_percent", Location: "form", Enum: []string{""}},
		{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "liability", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "backdate_start_date", Location: "form", Type: "integer"},
		{Name: "billing_cycle_anchor", Location: "form", Type: "integer"},
		{Name: "billing_cycle_anchor_config", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "day_of_month", Location: "form", Required: true, Type: "integer"},
			{Name: "hour", Location: "form", Type: "integer"},
			{Name: "minute", Location: "form", Type: "integer"},
			{Name: "month", Location: "form", Type: "integer"},
			{Name: "second", Location: "form", Type: "integer"},
		}},
		{Name: "billing_mode", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "flexible", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "proration_discounts", Location: "form", Type: "string", Enum: []string{"included", "itemized"}},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"classic", "flexible"}},
		}},
		{Name: "billing_thresholds", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "amount_gte", Location: "form", Type: "integer"},
			{Name: "reset_billing_cycle_anchor", Location: "form", Type: "boolean"},
		}},
		{Name: "cancel_at", Location: "form", Enum: []string{"max_period_end", "min_period_end"}},
		{Name: "cancel_at_period_end", Location: "form", Type: "boolean"},
		{Name: "collection_method", Location: "form", Type: "string", Enum: []string{"charge_automatically", "send_invoice"}},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_account", Location: "form", Type: "string"},
		{Name: "days_until_due", Location: "form", Type: "integer"},
		{Name: "default_payment_method", Location: "form", Type: "string"},
		{Name: "default_source", Location: "form", Type: "string"},
		{Name: "default_tax_rates", Location: "form", Enum: []string{""}},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "invoice_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_tax_ids", Location: "form", Enum: []string{""}},
			{Name: "issuer", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "items", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "off_session", Location: "form", Type: "boolean"},
		{Name: "on_behalf_of", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "payment_behavior", Location: "form", Type: "string", Enum: []string{"allow_incomplete", "default_incomplete", "error_if_incomplete", "pending_if_incomplete"}},
		{Name: "payment_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "acss_debit", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "bancontact", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "card", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "customer_balance", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "konbini", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "payto", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "pix", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "sepa_debit", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "upi", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "us_bank_account", Location: "form", Type: "object", Enum: []string{""}},
			}},
			{Name: "payment_method_types", Location: "form", Enum: []string{""}},
			{Name: "save_default_payment_method", Location: "form", Type: "string", Enum: []string{"off", "on_subscription"}},
		}},
		{Name: "pending_invoice_item_interval", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "interval", Location: "form", Required: true, Type: "string", Enum: []string{"day", "month", "week", "year"}},
			{Name: "interval_count", Location: "form", Type: "integer"},
		}},
		{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
		{Name: "transfer_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount_percent", Location: "form", Type: "number"},
			{Name: "destination", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "trial_end", Location: "form", Enum: []string{"now"}},
		{Name: "trial_from_plan", Location: "form", Type: "boolean"},
		{Name: "trial_period_days", Location: "form", Type: "integer"},
		{Name: "trial_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "end_behavior", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "missing_payment_method", Location: "form", Required: true, Type: "string", Enum: []string{"cancel", "create_invoice", "pause"}},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/subscriptions/search", OperationID: "GetSubscriptionsSearch", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "page", Location: "query", Type: "string"},
		{Name: "query", Location: "query", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/subscriptions/{subscription_exposed_id}", OperationID: "DeleteSubscriptionsSubscriptionExposedId", Params: []ParameterValidation{
		{Name: "subscription_exposed_id", Location: "path", Required: true, Type: "string"},
		{Name: "cancellation_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "comment", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "feedback", Location: "form", Type: "string", Enum: []string{"", "customer_service", "low_quality", "missing_features", "other", "switched_service", "too_complex", "too_expensive", "unused"}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "invoice_now", Location: "form", Type: "boolean"},
		{Name: "prorate", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/subscriptions/{subscription_exposed_id}", OperationID: "GetSubscriptionsSubscriptionExposedId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "subscription_exposed_id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/subscriptions/{subscription_exposed_id}", OperationID: "PostSubscriptionsSubscriptionExposedId", Params: []ParameterValidation{
		{Name: "subscription_exposed_id", Location: "path", Required: true, Type: "string"},
		{Name: "add_invoice_items", Location: "form", Type: "array"},
		{Name: "application_fee_percent", Location: "form", Enum: []string{""}},
		{Name: "automatic_tax", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
			{Name: "liability", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "billing_cycle_anchor", Location: "form", Type: "string", Enum: []string{"now", "unchanged"}},
		{Name: "billing_thresholds", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "amount_gte", Location: "form", Type: "integer"},
			{Name: "reset_billing_cycle_anchor", Location: "form", Type: "boolean"},
		}},
		{Name: "cancel_at", Location: "form", Enum: []string{"", "max_period_end", "min_period_end"}},
		{Name: "cancel_at_period_end", Location: "form", Type: "boolean"},
		{Name: "cancellation_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "comment", Location: "form", Type: "string", Enum: []string{""}},
			{Name: "feedback", Location: "form", Type: "string", Enum: []string{"", "customer_service", "low_quality", "missing_features", "other", "switched_service", "too_complex", "too_expensive", "unused"}},
		}},
		{Name: "collection_method", Location: "form", Type: "string", Enum: []string{"charge_automatically", "send_invoice"}},
		{Name: "days_until_due", Location: "form", Type: "integer"},
		{Name: "default_payment_method", Location: "form", Type: "string"},
		{Name: "default_source", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "default_tax_rates", Location: "form", Enum: []string{""}},
		{Name: "description", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "discounts", Location: "form", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "invoice_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_tax_ids", Location: "form", Enum: []string{""}},
			{Name: "issuer", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "self"}},
			}},
		}},
		{Name: "items", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "off_session", Location: "form", Type: "boolean"},
		{Name: "on_behalf_of", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "pause_collection", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "behavior", Location: "form", Required: true, Type: "string", Enum: []string{"keep_as_draft", "mark_uncollectible", "void"}},
			{Name: "resumes_at", Location: "form", Type: "integer"},
		}},
		{Name: "payment_behavior", Location: "form", Type: "string", Enum: []string{"allow_incomplete", "default_incomplete", "error_if_incomplete", "pending_if_incomplete"}},
		{Name: "payment_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "acss_debit", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "bancontact", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "card", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "customer_balance", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "konbini", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "payto", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "pix", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "sepa_debit", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "upi", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "us_bank_account", Location: "form", Type: "object", Enum: []string{""}},
			}},
			{Name: "payment_method_types", Location: "form", Enum: []string{""}},
			{Name: "save_default_payment_method", Location: "form", Type: "string", Enum: []string{"off", "on_subscription"}},
		}},
		{Name: "pending_invoice_item_interval", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "interval", Location: "form", Required: true, Type: "string", Enum: []string{"day", "month", "week", "year"}},
			{Name: "interval_count", Location: "form", Type: "integer"},
		}},
		{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
		{Name: "proration_date", Location: "form", Type: "integer"},
		{Name: "transfer_data", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "amount_percent", Location: "form", Type: "number"},
			{Name: "destination", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "trial_end", Location: "form", Enum: []string{"now"}},
		{Name: "trial_from_plan", Location: "form", Type: "boolean"},
		{Name: "trial_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "end_behavior", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "missing_payment_method", Location: "form", Required: true, Type: "string", Enum: []string{"cancel", "create_invoice", "pause"}},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/subscriptions/{subscription_exposed_id}/discount", OperationID: "DeleteSubscriptionsSubscriptionExposedIdDiscount", Params: []ParameterValidation{
		{Name: "subscription_exposed_id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/subscriptions/{subscription}/migrate", OperationID: "PostSubscriptionsSubscriptionMigrate", Params: []ParameterValidation{
		{Name: "subscription", Location: "path", Required: true, Type: "string"},
		{Name: "billing_mode", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "flexible", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "proration_discounts", Location: "form", Type: "string", Enum: []string{"included", "itemized"}},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"flexible"}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/subscriptions/{subscription}/resume", OperationID: "PostSubscriptionsSubscriptionResume", Params: []ParameterValidation{
		{Name: "subscription", Location: "path", Required: true, Type: "string"},
		{Name: "billing_cycle_anchor", Location: "form", Type: "string", Enum: []string{"now", "unchanged"}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "proration_behavior", Location: "form", Type: "string", Enum: []string{"always_invoice", "create_prorations", "none"}},
		{Name: "proration_date", Location: "form", Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax/associations/find", OperationID: "GetTaxAssociationsFind", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "payment_intent", Location: "query", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/tax/calculations", OperationID: "PostTaxCalculations", Params: []ParameterValidation{
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "customer_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "country", Location: "form", Required: true, Type: "string"},
				{Name: "line1", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "line2", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "postal_code", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "state", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "address_source", Location: "form", Type: "string", Enum: []string{"billing", "shipping"}},
			{Name: "ip_address", Location: "form", Type: "string"},
			{Name: "tax_ids", Location: "form", Type: "array"},
			{Name: "taxability_override", Location: "form", Type: "string", Enum: []string{"customer_exempt", "none", "reverse_charge"}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "line_items", Location: "form", Required: true, Type: "array"},
		{Name: "ship_from_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "country", Location: "form", Required: true, Type: "string"},
				{Name: "line1", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "line2", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "postal_code", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "state", Location: "form", Type: "string", Enum: []string{""}},
			}},
		}},
		{Name: "shipping_cost", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount", Location: "form", Type: "integer"},
			{Name: "shipping_rate", Location: "form", Type: "string"},
			{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive"}},
			{Name: "tax_code", Location: "form", Type: "string"},
		}},
		{Name: "tax_date", Location: "form", Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax/calculations/{calculation}", OperationID: "GetTaxCalculationsCalculation", Params: []ParameterValidation{
		{Name: "calculation", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax/calculations/{calculation}/line_items", OperationID: "GetTaxCalculationsCalculationLineItems", Params: []ParameterValidation{
		{Name: "calculation", Location: "path", Required: true, Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax/registrations", OperationID: "GetTaxRegistrations", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"active", "all", "expired", "scheduled"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/tax/registrations", OperationID: "PostTaxRegistrations", Params: []ParameterValidation{
		{Name: "active_from", Location: "form", Required: true, Enum: []string{"now"}},
		{Name: "country", Location: "form", Required: true, Type: "string"},
		{Name: "country_options", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "ae", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "al", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "am", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "ao", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "at", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "au", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "aw", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "az", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "ba", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "bb", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "bd", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "be", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "bf", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "bg", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "bh", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "bj", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "bs", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "by", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "ca", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "province_standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"province_standard", "simplified", "standard"}},
			}},
			{Name: "cd", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "ch", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "cl", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "cm", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "co", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "cr", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "cv", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "cy", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "cz", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "de", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "dk", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "ec", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "ee", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "eg", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "es", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "et", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "fi", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "fr", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "gb", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "ge", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "gn", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "gr", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "hr", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "hu", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "id", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "ie", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "in", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "is", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "it", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "jp", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "ke", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "kg", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "kh", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "kr", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "kz", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "la", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "lk", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "lt", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "lu", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "lv", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "ma", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "md", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "me", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "mk", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "mr", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "mt", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "mx", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "my", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "ng", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "nl", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "no", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "np", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "nz", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "om", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "pe", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "ph", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "pl", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "pt", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "ro", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "rs", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "ru", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "sa", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "se", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "sg", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "si", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "sk", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ioss", "oss_non_union", "oss_union", "standard"}},
			}},
			{Name: "sn", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "sr", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "th", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "tj", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "tr", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "tw", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "tz", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "ua", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "ug", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "us", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "local_amusement_tax", Location: "form", Type: "object"},
				{Name: "local_lease_tax", Location: "form", Type: "object"},
				{Name: "state", Location: "form", Required: true, Type: "string"},
				{Name: "state_sales_tax", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"local_amusement_tax", "local_lease_tax", "state_communications_tax", "state_retail_delivery_fee", "state_sales_tax"}},
			}},
			{Name: "uy", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "uz", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "vn", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "za", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
			{Name: "zm", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"simplified"}},
			}},
			{Name: "zw", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "standard", Location: "form", Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"standard"}},
			}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "expires_at", Location: "form", Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax/registrations/{id}", OperationID: "GetTaxRegistrationsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/tax/registrations/{id}", OperationID: "PostTaxRegistrationsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "active_from", Location: "form", Enum: []string{"now"}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "expires_at", Location: "form", Enum: []string{"", "now"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax/settings", OperationID: "GetTaxSettings", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/tax/settings", OperationID: "PostTaxSettings", Params: []ParameterValidation{
		{Name: "defaults", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "tax_behavior", Location: "form", Type: "string", Enum: []string{"exclusive", "inclusive", "inferred_by_currency"}},
			{Name: "tax_code", Location: "form", Type: "string"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "head_office", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/tax/transactions/create_from_calculation", OperationID: "PostTaxTransactionsCreateFromCalculation", Params: []ParameterValidation{
		{Name: "calculation", Location: "form", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "posted_at", Location: "form", Type: "integer"},
		{Name: "reference", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/tax/transactions/create_reversal", OperationID: "PostTaxTransactionsCreateReversal", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "flat_amount", Location: "form", Type: "integer"},
		{Name: "line_items", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "mode", Location: "form", Required: true, Type: "string", Enum: []string{"full", "partial"}},
		{Name: "original_transaction", Location: "form", Required: true, Type: "string"},
		{Name: "reference", Location: "form", Required: true, Type: "string"},
		{Name: "shipping_cost", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "amount", Location: "form", Required: true, Type: "integer"},
			{Name: "amount_tax", Location: "form", Required: true, Type: "integer"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax/transactions/{transaction}", OperationID: "GetTaxTransactionsTransaction", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "transaction", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax/transactions/{transaction}/line_items", OperationID: "GetTaxTransactionsTransactionLineItems", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "transaction", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax_codes", OperationID: "GetTaxCodes", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax_codes/{id}", OperationID: "GetTaxCodesId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax_ids", OperationID: "GetTaxIds", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "owner", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "query", Type: "string"},
			{Name: "customer", Location: "query", Type: "string"},
			{Name: "customer_account", Location: "query", Type: "string"},
			{Name: "type", Location: "query", Required: true, Type: "string", Enum: []string{"account", "application", "customer", "self"}},
		}},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/tax_ids", OperationID: "PostTaxIds", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "owner", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Type: "string"},
			{Name: "customer", Location: "form", Type: "string"},
			{Name: "customer_account", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account", "application", "customer", "self"}},
		}},
		{Name: "type", Location: "form", Required: true, Type: "string"},
		{Name: "value", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/tax_ids/{id}", OperationID: "DeleteTaxIdsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax_ids/{id}", OperationID: "GetTaxIdsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax_rates", OperationID: "GetTaxRates", Params: []ParameterValidation{
		{Name: "active", Location: "query", Type: "boolean"},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "inclusive", Location: "query", Type: "boolean"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/tax_rates", OperationID: "PostTaxRates", Params: []ParameterValidation{
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "country", Location: "form", Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "display_name", Location: "form", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "inclusive", Location: "form", Required: true, Type: "boolean"},
		{Name: "jurisdiction", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "percentage", Location: "form", Required: true, Type: "number"},
		{Name: "state", Location: "form", Type: "string"},
		{Name: "tax_type", Location: "form", Type: "string", Enum: []string{"amusement_tax", "communications_tax", "gst", "hst", "igst", "jct", "lease_tax", "pst", "qst", "retail_delivery_fee", "rst", "sales_tax", "service_tax", "vat"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tax_rates/{tax_rate}", OperationID: "GetTaxRatesTaxRate", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "tax_rate", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/tax_rates/{tax_rate}", OperationID: "PostTaxRatesTaxRate", Params: []ParameterValidation{
		{Name: "tax_rate", Location: "path", Required: true, Type: "string"},
		{Name: "active", Location: "form", Type: "boolean"},
		{Name: "country", Location: "form", Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "display_name", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "jurisdiction", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "state", Location: "form", Type: "string"},
		{Name: "tax_type", Location: "form", Type: "string", Enum: []string{"amusement_tax", "communications_tax", "gst", "hst", "igst", "jct", "lease_tax", "pst", "qst", "retail_delivery_fee", "rst", "sales_tax", "service_tax", "vat"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/terminal/configurations", OperationID: "GetTerminalConfigurations", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "is_account_default", Location: "query", Type: "boolean"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/configurations", OperationID: "PostTerminalConfigurations", Params: []ParameterValidation{
		{Name: "bbpos_wisepad3", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "splashscreen", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "bbpos_wisepos_e", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "splashscreen", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "cellular", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "offline", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "reboot_window", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "end_hour", Location: "form", Required: true, Type: "integer"},
			{Name: "start_hour", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "stripe_s700", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "splashscreen", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "stripe_s710", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "splashscreen", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "tipping", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "aed", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "aud", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "cad", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "chf", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "czk", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "dkk", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "eur", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "gbp", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "gip", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "hkd", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "huf", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "jpy", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "mxn", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "myr", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "nok", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "nzd", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "pln", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "ron", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "sek", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "sgd", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "usd", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
		}},
		{Name: "verifone_p400", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "splashscreen", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "wifi", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "enterprise_eap_peap", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ca_certificate_file", Location: "form", Type: "string"},
				{Name: "password", Location: "form", Required: true, Type: "string"},
				{Name: "ssid", Location: "form", Required: true, Type: "string"},
				{Name: "username", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "enterprise_eap_tls", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ca_certificate_file", Location: "form", Type: "string"},
				{Name: "client_certificate_file", Location: "form", Required: true, Type: "string"},
				{Name: "private_key_file", Location: "form", Required: true, Type: "string"},
				{Name: "private_key_file_password", Location: "form", Type: "string"},
				{Name: "ssid", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "personal_psk", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "password", Location: "form", Required: true, Type: "string"},
				{Name: "ssid", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"enterprise_eap_peap", "enterprise_eap_tls", "personal_psk"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/terminal/configurations/{configuration}", OperationID: "DeleteTerminalConfigurationsConfiguration", Params: []ParameterValidation{
		{Name: "configuration", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/terminal/configurations/{configuration}", OperationID: "GetTerminalConfigurationsConfiguration", Params: []ParameterValidation{
		{Name: "configuration", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/configurations/{configuration}", OperationID: "PostTerminalConfigurationsConfiguration", Params: []ParameterValidation{
		{Name: "configuration", Location: "path", Required: true, Type: "string"},
		{Name: "bbpos_wisepad3", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "splashscreen", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "bbpos_wisepos_e", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "splashscreen", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "cellular", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "offline", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "enabled", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "reboot_window", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "end_hour", Location: "form", Required: true, Type: "integer"},
			{Name: "start_hour", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "stripe_s700", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "splashscreen", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "stripe_s710", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "splashscreen", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "tipping", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "aed", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "aud", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "cad", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "chf", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "czk", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "dkk", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "eur", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "gbp", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "gip", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "hkd", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "huf", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "jpy", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "mxn", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "myr", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "nok", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "nzd", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "pln", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "ron", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "sek", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "sgd", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
			{Name: "usd", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fixed_amounts", Location: "form", Type: "array"},
				{Name: "percentages", Location: "form", Type: "array"},
				{Name: "smart_tip_threshold", Location: "form", Type: "integer"},
			}},
		}},
		{Name: "verifone_p400", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "splashscreen", Location: "form", Type: "string", Enum: []string{""}},
		}},
		{Name: "wifi", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
			{Name: "enterprise_eap_peap", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ca_certificate_file", Location: "form", Type: "string"},
				{Name: "password", Location: "form", Required: true, Type: "string"},
				{Name: "ssid", Location: "form", Required: true, Type: "string"},
				{Name: "username", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "enterprise_eap_tls", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ca_certificate_file", Location: "form", Type: "string"},
				{Name: "client_certificate_file", Location: "form", Required: true, Type: "string"},
				{Name: "private_key_file", Location: "form", Required: true, Type: "string"},
				{Name: "private_key_file_password", Location: "form", Type: "string"},
				{Name: "ssid", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "personal_psk", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "password", Location: "form", Required: true, Type: "string"},
				{Name: "ssid", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"enterprise_eap_peap", "enterprise_eap_tls", "personal_psk"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/connection_tokens", OperationID: "PostTerminalConnectionTokens", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "location", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/terminal/locations", OperationID: "GetTerminalLocations", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/locations", OperationID: "PostTerminalLocations", Params: []ParameterValidation{
		{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Required: true, Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
		}},
		{Name: "address_kana", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "address_kanji", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "configuration_overrides", Location: "form", Type: "string"},
		{Name: "display_name", Location: "form", Type: "string"},
		{Name: "display_name_kana", Location: "form", Type: "string"},
		{Name: "display_name_kanji", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "phone", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/terminal/locations/{location}", OperationID: "DeleteTerminalLocationsLocation", Params: []ParameterValidation{
		{Name: "location", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/terminal/locations/{location}", OperationID: "GetTerminalLocationsLocation", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "location", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/locations/{location}", OperationID: "PostTerminalLocationsLocation", Params: []ParameterValidation{
		{Name: "location", Location: "path", Required: true, Type: "string"},
		{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
		}},
		{Name: "address_kana", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "address_kanji", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "configuration_overrides", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "display_name", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "display_name_kana", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "display_name_kanji", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/onboarding_links", OperationID: "PostTerminalOnboardingLinks", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "link_options", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "apple_terms_and_conditions", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "allow_relinking", Location: "form", Type: "boolean"},
				{Name: "merchant_display_name", Location: "form", Required: true, Type: "string"},
			}},
		}},
		{Name: "link_type", Location: "form", Required: true, Type: "string", Enum: []string{"apple_terms_and_conditions"}},
		{Name: "on_behalf_of", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/terminal/readers", OperationID: "GetTerminalReaders", Params: []ParameterValidation{
		{Name: "device_type", Location: "query", Type: "string", Enum: []string{"bbpos_chipper2x", "bbpos_wisepad3", "bbpos_wisepos_e", "mobile_phone_reader", "simulated_stripe_s700", "simulated_stripe_s710", "simulated_wisepos_e", "stripe_m2", "stripe_s700", "stripe_s710", "verifone_P400"}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "location", Location: "query", Type: "string"},
		{Name: "serial_number", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"offline", "online"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/readers", OperationID: "PostTerminalReaders", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "label", Location: "form", Type: "string"},
		{Name: "location", Location: "form", Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "registration_code", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/terminal/readers/{reader}", OperationID: "DeleteTerminalReadersReader", Params: []ParameterValidation{
		{Name: "reader", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/terminal/readers/{reader}", OperationID: "GetTerminalReadersReader", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "reader", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/readers/{reader}", OperationID: "PostTerminalReadersReader", Params: []ParameterValidation{
		{Name: "reader", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "label", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/readers/{reader}/cancel_action", OperationID: "PostTerminalReadersReaderCancelAction", Params: []ParameterValidation{
		{Name: "reader", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/readers/{reader}/collect_inputs", OperationID: "PostTerminalReadersReaderCollectInputs", Params: []ParameterValidation{
		{Name: "reader", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "inputs", Location: "form", Required: true, Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/readers/{reader}/collect_payment_method", OperationID: "PostTerminalReadersReaderCollectPaymentMethod", Params: []ParameterValidation{
		{Name: "reader", Location: "path", Required: true, Type: "string"},
		{Name: "collect_config", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "allow_redisplay", Location: "form", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
			{Name: "enable_customer_cancellation", Location: "form", Type: "boolean"},
			{Name: "skip_tipping", Location: "form", Type: "boolean"},
			{Name: "tipping", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "amount_eligible", Location: "form", Type: "integer"},
			}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "payment_intent", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/readers/{reader}/confirm_payment_intent", OperationID: "PostTerminalReadersReaderConfirmPaymentIntent", Params: []ParameterValidation{
		{Name: "reader", Location: "path", Required: true, Type: "string"},
		{Name: "confirm_config", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "return_url", Location: "form", Type: "string"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "payment_intent", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/readers/{reader}/process_payment_intent", OperationID: "PostTerminalReadersReaderProcessPaymentIntent", Params: []ParameterValidation{
		{Name: "reader", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "payment_intent", Location: "form", Required: true, Type: "string"},
		{Name: "process_config", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "allow_redisplay", Location: "form", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
			{Name: "enable_customer_cancellation", Location: "form", Type: "boolean"},
			{Name: "return_url", Location: "form", Type: "string"},
			{Name: "skip_tipping", Location: "form", Type: "boolean"},
			{Name: "tipping", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "amount_eligible", Location: "form", Type: "integer"},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/readers/{reader}/process_setup_intent", OperationID: "PostTerminalReadersReaderProcessSetupIntent", Params: []ParameterValidation{
		{Name: "reader", Location: "path", Required: true, Type: "string"},
		{Name: "allow_redisplay", Location: "form", Required: true, Type: "string", Enum: []string{"always", "limited", "unspecified"}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "process_config", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enable_customer_cancellation", Location: "form", Type: "boolean"},
		}},
		{Name: "setup_intent", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/readers/{reader}/refund_payment", OperationID: "PostTerminalReadersReaderRefundPayment", Params: []ParameterValidation{
		{Name: "reader", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "charge", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "payment_intent", Location: "form", Type: "string"},
		{Name: "refund_application_fee", Location: "form", Type: "boolean"},
		{Name: "refund_payment_config", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "enable_customer_cancellation", Location: "form", Type: "boolean"},
		}},
		{Name: "reverse_transfer", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/readers/{reader}/set_reader_display", OperationID: "PostTerminalReadersReaderSetReaderDisplay", Params: []ParameterValidation{
		{Name: "reader", Location: "path", Required: true, Type: "string"},
		{Name: "cart", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "currency", Location: "form", Required: true, Type: "string"},
			{Name: "line_items", Location: "form", Required: true, Type: "array"},
			{Name: "tax", Location: "form", Type: "integer"},
			{Name: "total", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"cart"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/terminal/refunds", OperationID: "PostTerminalRefunds", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "charge", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "payment_intent", Location: "form", Type: "string"},
		{Name: "reason", Location: "form", Type: "string", Enum: []string{"duplicate", "fraudulent", "requested_by_customer"}},
		{Name: "refund_application_fee", Location: "form", Type: "boolean"},
		{Name: "reverse_transfer", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/confirmation_tokens", OperationID: "PostTestHelpersConfirmationTokens", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "payment_method", Location: "form", Type: "string"},
		{Name: "payment_method_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acss_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "institution_number", Location: "form", Required: true, Type: "string"},
				{Name: "transit_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "affirm", Location: "form", Type: "object"},
			{Name: "afterpay_clearpay", Location: "form", Type: "object"},
			{Name: "alipay", Location: "form", Type: "object"},
			{Name: "allow_redisplay", Location: "form", Type: "string", Enum: []string{"always", "limited", "unspecified"}},
			{Name: "alma", Location: "form", Type: "object"},
			{Name: "amazon_pay", Location: "form", Type: "object"},
			{Name: "au_becs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bsb_number", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "bacs_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "sort_code", Location: "form", Type: "string"},
			}},
			{Name: "bancontact", Location: "form", Type: "object"},
			{Name: "billie", Location: "form", Type: "object"},
			{Name: "billing_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "email", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "name", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "tax_id", Location: "form", Type: "string"},
			}},
			{Name: "blik", Location: "form", Type: "object"},
			{Name: "boleto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "tax_id", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "cashapp", Location: "form", Type: "object"},
			{Name: "crypto", Location: "form", Type: "object"},
			{Name: "customer_balance", Location: "form", Type: "object"},
			{Name: "eps", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"arzte_und_apotheker_bank", "austrian_anadi_bank_ag", "bank_austria", "bankhaus_carl_spangler", "bankhaus_schelhammer_und_schattera_ag", "bawag_psk_ag", "bks_bank_ag", "brull_kallmus_bank_ag", "btv_vier_lander_bank", "capital_bank_grawe_gruppe_ag", "deutsche_bank_ag", "dolomitenbank", "easybank_ag", "erste_bank_und_sparkassen", "hypo_alpeadriabank_international_ag", "hypo_bank_burgenland_aktiengesellschaft", "hypo_noe_lb_fur_niederosterreich_u_wien", "hypo_oberosterreich_salzburg_steiermark", "hypo_tirol_bank_ag", "hypo_vorarlberg_bank_ag", "marchfelder_bank", "oberbank_ag", "raiffeisen_bankengruppe_osterreich", "schoellerbank_ag", "sparda_bank_wien", "volksbank_gruppe", "volkskreditbank_ag", "vr_bank_braunau"}},
			}},
			{Name: "fpx", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Required: true, Type: "string", Enum: []string{"affin_bank", "agrobank", "alliance_bank", "ambank", "bank_islam", "bank_muamalat", "bank_of_china", "bank_rakyat", "bsn", "cimb", "deutsche_bank", "hong_leong_bank", "hsbc", "kfh", "maybank2e", "maybank2u", "ocbc", "pb_enterprise", "public_bank", "rhb", "standard_chartered", "uob"}},
			}},
			{Name: "giropay", Location: "form", Type: "object"},
			{Name: "grabpay", Location: "form", Type: "object"},
			{Name: "ideal", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"abn_amro", "adyen", "asn_bank", "bunq", "buut", "finom", "handelsbanken", "ing", "knab", "mollie", "moneyou", "n26", "nn", "rabobank", "regiobank", "revolut", "sns_bank", "triodos_bank", "van_lanschot", "yoursafe"}},
			}},
			{Name: "interac_present", Location: "form", Type: "object"},
			{Name: "kakao_pay", Location: "form", Type: "object"},
			{Name: "klarna", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "dob", Location: "form", Type: "object"},
			}},
			{Name: "konbini", Location: "form", Type: "object"},
			{Name: "kr_card", Location: "form", Type: "object"},
			{Name: "link", Location: "form", Type: "object"},
			{Name: "mb_way", Location: "form", Type: "object"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "mobilepay", Location: "form", Type: "object"},
			{Name: "multibanco", Location: "form", Type: "object"},
			{Name: "naver_pay", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "funding", Location: "form", Type: "string", Enum: []string{"card", "points"}},
			}},
			{Name: "nz_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_name", Location: "form", Type: "string"},
				{Name: "account_number", Location: "form", Required: true, Type: "string"},
				{Name: "bank_code", Location: "form", Required: true, Type: "string"},
				{Name: "branch_code", Location: "form", Required: true, Type: "string"},
				{Name: "reference", Location: "form", Type: "string"},
				{Name: "suffix", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "oxxo", Location: "form", Type: "object"},
			{Name: "p24", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bank", Location: "form", Type: "string", Enum: []string{"alior_bank", "bank_millennium", "bank_nowy_bfg_sa", "bank_pekao_sa", "banki_spbdzielcze", "blik", "bnp_paribas", "boz", "citi_handlowy", "credit_agricole", "envelobank", "etransfer_pocztowy24", "getin_bank", "ideabank", "ing", "inteligo", "mbank_mtransfer", "nest_przelew", "noble_pay", "pbac_z_ipko", "plus_bank", "santander_przelew24", "tmobile_usbugi_bankowe", "toyota_bank", "velobank", "volkswagen_bank"}},
			}},
			{Name: "pay_by_bank", Location: "form", Type: "object"},
			{Name: "payco", Location: "form", Type: "object"},
			{Name: "paynow", Location: "form", Type: "object"},
			{Name: "paypal", Location: "form", Type: "object"},
			{Name: "payto", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "bsb_number", Location: "form", Type: "string"},
				{Name: "pay_id", Location: "form", Type: "string"},
			}},
			{Name: "pix", Location: "form", Type: "object"},
			{Name: "promptpay", Location: "form", Type: "object"},
			{Name: "radar_options", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "session", Location: "form", Type: "string"},
			}},
			{Name: "revolut_pay", Location: "form", Type: "object"},
			{Name: "samsung_pay", Location: "form", Type: "object"},
			{Name: "satispay", Location: "form", Type: "object"},
			{Name: "sepa_debit", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "iban", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "sofort", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "country", Location: "form", Required: true, Type: "string", Enum: []string{"AT", "BE", "DE", "ES", "IT", "NL"}},
			}},
			{Name: "sunbit", Location: "form", Type: "object"},
			{Name: "swish", Location: "form", Type: "object"},
			{Name: "twint", Location: "form", Type: "object"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"acss_debit", "affirm", "afterpay_clearpay", "alipay", "alma", "amazon_pay", "au_becs_debit", "bacs_debit", "bancontact", "billie", "blik", "boleto", "cashapp", "crypto", "customer_balance", "eps", "fpx", "giropay", "grabpay", "ideal", "kakao_pay", "klarna", "konbini", "kr_card", "link", "mb_way", "mobilepay", "multibanco", "naver_pay", "nz_bank_account", "oxxo", "p24", "pay_by_bank", "payco", "paynow", "paypal", "payto", "pix", "promptpay", "revolut_pay", "samsung_pay", "satispay", "sepa_debit", "sofort", "sunbit", "swish", "twint", "upi", "us_bank_account", "wechat_pay", "zip"}},
			{Name: "upi", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "mandate_options", Location: "form", Type: "object"},
			}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "savings"}},
				{Name: "financial_connections_account", Location: "form", Type: "string"},
				{Name: "routing_number", Location: "form", Type: "string"},
			}},
			{Name: "wechat_pay", Location: "form", Type: "object"},
			{Name: "zip", Location: "form", Type: "object"},
		}},
		{Name: "payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "installments", Location: "form", Type: "object"},
			}},
		}},
		{Name: "return_url", Location: "form", Type: "string"},
		{Name: "setup_future_usage", Location: "form", Type: "string", Enum: []string{"off_session", "on_session"}},
		{Name: "shipping", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "name", Location: "form", Required: true, Type: "string"},
			{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/customers/{customer}/fund_cash_balance", OperationID: "PostTestHelpersCustomersCustomerFundCashBalance", Params: []ParameterValidation{
		{Name: "customer", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Required: true, Type: "integer"},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "reference", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/authorizations", OperationID: "PostTestHelpersIssuingAuthorizations", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "amount_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "atm_fee", Location: "form", Type: "integer"},
			{Name: "cashback_amount", Location: "form", Type: "integer"},
		}},
		{Name: "authorization_method", Location: "form", Type: "string", Enum: []string{"chip", "contactless", "keyed_in", "online", "swipe"}},
		{Name: "card", Location: "form", Required: true, Type: "string"},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "fleet", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "cardholder_prompt_data", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "driver_id", Location: "form", Type: "string"},
				{Name: "odometer", Location: "form", Type: "integer"},
				{Name: "unspecified_id", Location: "form", Type: "string"},
				{Name: "user_id", Location: "form", Type: "string"},
				{Name: "vehicle_number", Location: "form", Type: "string"},
			}},
			{Name: "purchase_type", Location: "form", Type: "string", Enum: []string{"fuel_and_non_fuel_purchase", "fuel_purchase", "non_fuel_purchase"}},
			{Name: "reported_breakdown", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fuel", Location: "form", Type: "object"},
				{Name: "non_fuel", Location: "form", Type: "object"},
				{Name: "tax", Location: "form", Type: "object"},
			}},
			{Name: "service_type", Location: "form", Type: "string", Enum: []string{"full_service", "non_fuel_transaction", "self_service"}},
		}},
		{Name: "fraud_disputability_likelihood", Location: "form", Type: "string", Enum: []string{"neutral", "unknown", "very_likely", "very_unlikely"}},
		{Name: "fuel", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "industry_product_code", Location: "form", Type: "string"},
			{Name: "quantity_decimal", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Type: "string", Enum: []string{"diesel", "other", "unleaded_plus", "unleaded_regular", "unleaded_super"}},
			{Name: "unit", Location: "form", Type: "string", Enum: []string{"charging_minute", "imperial_gallon", "kilogram", "kilowatt_hour", "liter", "other", "pound", "us_gallon"}},
			{Name: "unit_cost_decimal", Location: "form", Type: "string"},
		}},
		{Name: "is_amount_controllable", Location: "form", Type: "boolean"},
		{Name: "merchant_amount", Location: "form", Type: "integer"},
		{Name: "merchant_currency", Location: "form", Type: "string"},
		{Name: "merchant_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "category", Location: "form", Type: "string"},
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "network_id", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "terminal_id", Location: "form", Type: "string"},
			{Name: "url", Location: "form", Type: "string"},
		}},
		{Name: "network_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "acquiring_institution_id", Location: "form", Type: "string"},
		}},
		{Name: "risk_assessment", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "card_testing_risk", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "invalid_account_number_decline_rate_past_hour", Location: "form", Type: "integer"},
				{Name: "invalid_credentials_decline_rate_past_hour", Location: "form", Type: "integer"},
				{Name: "level", Location: "form", Required: true, Type: "string", Enum: []string{"elevated", "highest", "low", "normal", "not_assessed", "unknown"}},
			}},
			{Name: "fraud_risk", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "level", Location: "form", Required: true, Type: "string", Enum: []string{"elevated", "highest", "low", "normal", "not_assessed", "unknown"}},
				{Name: "score", Location: "form", Type: "number"},
			}},
			{Name: "merchant_dispute_risk", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "dispute_rate", Location: "form", Type: "integer"},
				{Name: "level", Location: "form", Required: true, Type: "string", Enum: []string{"elevated", "highest", "low", "normal", "not_assessed", "unknown"}},
			}},
		}},
		{Name: "verification_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address_line1_check", Location: "form", Type: "string", Enum: []string{"match", "mismatch", "not_provided"}},
			{Name: "address_postal_code_check", Location: "form", Type: "string", Enum: []string{"match", "mismatch", "not_provided"}},
			{Name: "authentication_exemption", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "claimed_by", Location: "form", Required: true, Type: "string", Enum: []string{"acquirer", "issuer"}},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"low_value_transaction", "transaction_risk_analysis", "unknown"}},
			}},
			{Name: "cvc_check", Location: "form", Type: "string", Enum: []string{"match", "mismatch", "not_provided"}},
			{Name: "expiry_check", Location: "form", Type: "string", Enum: []string{"match", "mismatch", "not_provided"}},
			{Name: "three_d_secure", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "result", Location: "form", Required: true, Type: "string", Enum: []string{"attempt_acknowledged", "authenticated", "failed", "required"}},
			}},
		}},
		{Name: "wallet", Location: "form", Type: "string", Enum: []string{"apple_pay", "google_pay", "samsung_pay"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/authorizations/{authorization}/capture", OperationID: "PostTestHelpersIssuingAuthorizationsAuthorizationCapture", Params: []ParameterValidation{
		{Name: "authorization", Location: "path", Required: true, Type: "string"},
		{Name: "capture_amount", Location: "form", Type: "integer"},
		{Name: "close_authorization", Location: "form", Type: "boolean"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "purchase_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "fleet", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "cardholder_prompt_data", Location: "form", Type: "object"},
				{Name: "purchase_type", Location: "form", Type: "string", Enum: []string{"fuel_and_non_fuel_purchase", "fuel_purchase", "non_fuel_purchase"}},
				{Name: "reported_breakdown", Location: "form", Type: "object"},
				{Name: "service_type", Location: "form", Type: "string", Enum: []string{"full_service", "non_fuel_transaction", "self_service"}},
			}},
			{Name: "flight", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "departure_at", Location: "form", Type: "integer"},
				{Name: "passenger_name", Location: "form", Type: "string"},
				{Name: "refundable", Location: "form", Type: "boolean"},
				{Name: "segments", Location: "form", Type: "array"},
				{Name: "travel_agency", Location: "form", Type: "string"},
			}},
			{Name: "fuel", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "industry_product_code", Location: "form", Type: "string"},
				{Name: "quantity_decimal", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Type: "string", Enum: []string{"diesel", "other", "unleaded_plus", "unleaded_regular", "unleaded_super"}},
				{Name: "unit", Location: "form", Type: "string", Enum: []string{"charging_minute", "imperial_gallon", "kilogram", "kilowatt_hour", "liter", "other", "pound", "us_gallon"}},
				{Name: "unit_cost_decimal", Location: "form", Type: "string"},
			}},
			{Name: "lodging", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "check_in_at", Location: "form", Type: "integer"},
				{Name: "nights", Location: "form", Type: "integer"},
			}},
			{Name: "receipt", Location: "form", Type: "array"},
			{Name: "reference", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/authorizations/{authorization}/expire", OperationID: "PostTestHelpersIssuingAuthorizationsAuthorizationExpire", Params: []ParameterValidation{
		{Name: "authorization", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/authorizations/{authorization}/finalize_amount", OperationID: "PostTestHelpersIssuingAuthorizationsAuthorizationFinalizeAmount", Params: []ParameterValidation{
		{Name: "authorization", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "final_amount", Location: "form", Required: true, Type: "integer"},
		{Name: "fleet", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "cardholder_prompt_data", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "driver_id", Location: "form", Type: "string"},
				{Name: "odometer", Location: "form", Type: "integer"},
				{Name: "unspecified_id", Location: "form", Type: "string"},
				{Name: "user_id", Location: "form", Type: "string"},
				{Name: "vehicle_number", Location: "form", Type: "string"},
			}},
			{Name: "purchase_type", Location: "form", Type: "string", Enum: []string{"fuel_and_non_fuel_purchase", "fuel_purchase", "non_fuel_purchase"}},
			{Name: "reported_breakdown", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fuel", Location: "form", Type: "object"},
				{Name: "non_fuel", Location: "form", Type: "object"},
				{Name: "tax", Location: "form", Type: "object"},
			}},
			{Name: "service_type", Location: "form", Type: "string", Enum: []string{"full_service", "non_fuel_transaction", "self_service"}},
		}},
		{Name: "fuel", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "industry_product_code", Location: "form", Type: "string"},
			{Name: "quantity_decimal", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Type: "string", Enum: []string{"diesel", "other", "unleaded_plus", "unleaded_regular", "unleaded_super"}},
			{Name: "unit", Location: "form", Type: "string", Enum: []string{"charging_minute", "imperial_gallon", "kilogram", "kilowatt_hour", "liter", "other", "pound", "us_gallon"}},
			{Name: "unit_cost_decimal", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/authorizations/{authorization}/fraud_challenges/respond", OperationID: "PostTestHelpersIssuingAuthorizationsAuthorizationFraudChallengesRespond", Params: []ParameterValidation{
		{Name: "authorization", Location: "path", Required: true, Type: "string"},
		{Name: "confirmed", Location: "form", Required: true, Type: "boolean"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/authorizations/{authorization}/increment", OperationID: "PostTestHelpersIssuingAuthorizationsAuthorizationIncrement", Params: []ParameterValidation{
		{Name: "authorization", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "increment_amount", Location: "form", Required: true, Type: "integer"},
		{Name: "is_amount_controllable", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/authorizations/{authorization}/reverse", OperationID: "PostTestHelpersIssuingAuthorizationsAuthorizationReverse", Params: []ParameterValidation{
		{Name: "authorization", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "reverse_amount", Location: "form", Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/cards/{card}/shipping/deliver", OperationID: "PostTestHelpersIssuingCardsCardShippingDeliver", Params: []ParameterValidation{
		{Name: "card", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/cards/{card}/shipping/fail", OperationID: "PostTestHelpersIssuingCardsCardShippingFail", Params: []ParameterValidation{
		{Name: "card", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/cards/{card}/shipping/return", OperationID: "PostTestHelpersIssuingCardsCardShippingReturn", Params: []ParameterValidation{
		{Name: "card", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/cards/{card}/shipping/ship", OperationID: "PostTestHelpersIssuingCardsCardShippingShip", Params: []ParameterValidation{
		{Name: "card", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/cards/{card}/shipping/submit", OperationID: "PostTestHelpersIssuingCardsCardShippingSubmit", Params: []ParameterValidation{
		{Name: "card", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/personalization_designs/{personalization_design}/activate", OperationID: "PostTestHelpersIssuingPersonalizationDesignsPersonalizationDesignActivate", Params: []ParameterValidation{
		{Name: "personalization_design", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/personalization_designs/{personalization_design}/deactivate", OperationID: "PostTestHelpersIssuingPersonalizationDesignsPersonalizationDesignDeactivate", Params: []ParameterValidation{
		{Name: "personalization_design", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/personalization_designs/{personalization_design}/reject", OperationID: "PostTestHelpersIssuingPersonalizationDesignsPersonalizationDesignReject", Params: []ParameterValidation{
		{Name: "personalization_design", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "rejection_reasons", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "card_logo", Location: "form", Type: "array"},
			{Name: "carrier_text", Location: "form", Type: "array"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/settlements", OperationID: "PostTestHelpersIssuingSettlements", Params: []ParameterValidation{
		{Name: "bin", Location: "form", Required: true, Type: "string"},
		{Name: "clearing_date", Location: "form", Required: true, Type: "integer"},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "interchange_fees_amount", Location: "form", Type: "integer"},
		{Name: "net_total_amount", Location: "form", Required: true, Type: "integer"},
		{Name: "network", Location: "form", Type: "string", Enum: []string{"maestro", "visa"}},
		{Name: "network_settlement_identifier", Location: "form", Type: "string"},
		{Name: "transaction_amount", Location: "form", Type: "integer"},
		{Name: "transaction_count", Location: "form", Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/settlements/{settlement}/complete", OperationID: "PostTestHelpersIssuingSettlementsSettlementComplete", Params: []ParameterValidation{
		{Name: "settlement", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/transactions/create_force_capture", OperationID: "PostTestHelpersIssuingTransactionsCreateForceCapture", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Required: true, Type: "integer"},
		{Name: "card", Location: "form", Required: true, Type: "string"},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "merchant_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "category", Location: "form", Type: "string"},
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "network_id", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "terminal_id", Location: "form", Type: "string"},
			{Name: "url", Location: "form", Type: "string"},
		}},
		{Name: "purchase_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "fleet", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "cardholder_prompt_data", Location: "form", Type: "object"},
				{Name: "purchase_type", Location: "form", Type: "string", Enum: []string{"fuel_and_non_fuel_purchase", "fuel_purchase", "non_fuel_purchase"}},
				{Name: "reported_breakdown", Location: "form", Type: "object"},
				{Name: "service_type", Location: "form", Type: "string", Enum: []string{"full_service", "non_fuel_transaction", "self_service"}},
			}},
			{Name: "flight", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "departure_at", Location: "form", Type: "integer"},
				{Name: "passenger_name", Location: "form", Type: "string"},
				{Name: "refundable", Location: "form", Type: "boolean"},
				{Name: "segments", Location: "form", Type: "array"},
				{Name: "travel_agency", Location: "form", Type: "string"},
			}},
			{Name: "fuel", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "industry_product_code", Location: "form", Type: "string"},
				{Name: "quantity_decimal", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Type: "string", Enum: []string{"diesel", "other", "unleaded_plus", "unleaded_regular", "unleaded_super"}},
				{Name: "unit", Location: "form", Type: "string", Enum: []string{"charging_minute", "imperial_gallon", "kilogram", "kilowatt_hour", "liter", "other", "pound", "us_gallon"}},
				{Name: "unit_cost_decimal", Location: "form", Type: "string"},
			}},
			{Name: "lodging", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "check_in_at", Location: "form", Type: "integer"},
				{Name: "nights", Location: "form", Type: "integer"},
			}},
			{Name: "receipt", Location: "form", Type: "array"},
			{Name: "reference", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/transactions/create_unlinked_refund", OperationID: "PostTestHelpersIssuingTransactionsCreateUnlinkedRefund", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Required: true, Type: "integer"},
		{Name: "card", Location: "form", Required: true, Type: "string"},
		{Name: "currency", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "merchant_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "category", Location: "form", Type: "string"},
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "network_id", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "terminal_id", Location: "form", Type: "string"},
			{Name: "url", Location: "form", Type: "string"},
		}},
		{Name: "purchase_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "fleet", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "cardholder_prompt_data", Location: "form", Type: "object"},
				{Name: "purchase_type", Location: "form", Type: "string", Enum: []string{"fuel_and_non_fuel_purchase", "fuel_purchase", "non_fuel_purchase"}},
				{Name: "reported_breakdown", Location: "form", Type: "object"},
				{Name: "service_type", Location: "form", Type: "string", Enum: []string{"full_service", "non_fuel_transaction", "self_service"}},
			}},
			{Name: "flight", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "departure_at", Location: "form", Type: "integer"},
				{Name: "passenger_name", Location: "form", Type: "string"},
				{Name: "refundable", Location: "form", Type: "boolean"},
				{Name: "segments", Location: "form", Type: "array"},
				{Name: "travel_agency", Location: "form", Type: "string"},
			}},
			{Name: "fuel", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "industry_product_code", Location: "form", Type: "string"},
				{Name: "quantity_decimal", Location: "form", Type: "string"},
				{Name: "type", Location: "form", Type: "string", Enum: []string{"diesel", "other", "unleaded_plus", "unleaded_regular", "unleaded_super"}},
				{Name: "unit", Location: "form", Type: "string", Enum: []string{"charging_minute", "imperial_gallon", "kilogram", "kilowatt_hour", "liter", "other", "pound", "us_gallon"}},
				{Name: "unit_cost_decimal", Location: "form", Type: "string"},
			}},
			{Name: "lodging", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "check_in_at", Location: "form", Type: "integer"},
				{Name: "nights", Location: "form", Type: "integer"},
			}},
			{Name: "receipt", Location: "form", Type: "array"},
			{Name: "reference", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/issuing/transactions/{transaction}/refund", OperationID: "PostTestHelpersIssuingTransactionsTransactionRefund", Params: []ParameterValidation{
		{Name: "transaction", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "refund_amount", Location: "form", Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/refunds/{refund}/expire", OperationID: "PostTestHelpersRefundsRefundExpire", Params: []ParameterValidation{
		{Name: "refund", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/terminal/readers/{reader}/present_payment_method", OperationID: "PostTestHelpersTerminalReadersReaderPresentPaymentMethod", Params: []ParameterValidation{
		{Name: "reader", Location: "path", Required: true, Type: "string"},
		{Name: "amount_tip", Location: "form", Type: "integer"},
		{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "cvc", Location: "form", Type: "string"},
			{Name: "exp_month", Location: "form", Required: true, Type: "integer"},
			{Name: "exp_year", Location: "form", Required: true, Type: "integer"},
			{Name: "number", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "card_present", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "number", Location: "form", Type: "string"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "interac_present", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "number", Location: "form", Type: "string"},
		}},
		{Name: "type", Location: "form", Type: "string", Enum: []string{"card", "card_present", "interac_present"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/terminal/readers/{reader}/succeed_input_collection", OperationID: "PostTestHelpersTerminalReadersReaderSucceedInputCollection", Params: []ParameterValidation{
		{Name: "reader", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "skip_non_required_inputs", Location: "form", Type: "string", Enum: []string{"all", "none"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/terminal/readers/{reader}/timeout_input_collection", OperationID: "PostTestHelpersTerminalReadersReaderTimeoutInputCollection", Params: []ParameterValidation{
		{Name: "reader", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/test_helpers/test_clocks", OperationID: "GetTestHelpersTestClocks", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/test_clocks", OperationID: "PostTestHelpersTestClocks", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "frozen_time", Location: "form", Required: true, Type: "integer"},
		{Name: "name", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/test_helpers/test_clocks/{test_clock}", OperationID: "DeleteTestHelpersTestClocksTestClock", Params: []ParameterValidation{
		{Name: "test_clock", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/test_helpers/test_clocks/{test_clock}", OperationID: "GetTestHelpersTestClocksTestClock", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "test_clock", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/test_clocks/{test_clock}/advance", OperationID: "PostTestHelpersTestClocksTestClockAdvance", Params: []ParameterValidation{
		{Name: "test_clock", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "frozen_time", Location: "form", Required: true, Type: "integer"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/treasury/inbound_transfers/{id}/fail", OperationID: "PostTestHelpersTreasuryInboundTransfersIdFail", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "failure_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "code", Location: "form", Type: "string", Enum: []string{"account_closed", "account_frozen", "bank_account_restricted", "bank_ownership_changed", "debit_not_authorized", "incorrect_account_holder_address", "incorrect_account_holder_name", "incorrect_account_holder_tax_id", "insufficient_funds", "invalid_account_number", "invalid_currency", "no_account", "other"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/treasury/inbound_transfers/{id}/return", OperationID: "PostTestHelpersTreasuryInboundTransfersIdReturn", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/treasury/inbound_transfers/{id}/succeed", OperationID: "PostTestHelpersTreasuryInboundTransfersIdSucceed", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/treasury/outbound_payments/{id}", OperationID: "PostTestHelpersTreasuryOutboundPaymentsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "tracking_details", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "ach", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "trace_id", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ach", "us_domestic_wire"}},
			{Name: "us_domestic_wire", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "chips", Location: "form", Type: "string"},
				{Name: "imad", Location: "form", Type: "string"},
				{Name: "omad", Location: "form", Type: "string"},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/treasury/outbound_payments/{id}/fail", OperationID: "PostTestHelpersTreasuryOutboundPaymentsIdFail", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/treasury/outbound_payments/{id}/post", OperationID: "PostTestHelpersTreasuryOutboundPaymentsIdPost", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/treasury/outbound_payments/{id}/return", OperationID: "PostTestHelpersTreasuryOutboundPaymentsIdReturn", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "returned_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "code", Location: "form", Type: "string", Enum: []string{"account_closed", "account_frozen", "bank_account_restricted", "bank_ownership_changed", "declined", "incorrect_account_holder_name", "invalid_account_number", "invalid_currency", "no_account", "other"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/treasury/outbound_transfers/{outbound_transfer}", OperationID: "PostTestHelpersTreasuryOutboundTransfersOutboundTransfer", Params: []ParameterValidation{
		{Name: "outbound_transfer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "tracking_details", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "ach", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "trace_id", Location: "form", Required: true, Type: "string"},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"ach", "us_domestic_wire"}},
			{Name: "us_domestic_wire", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "chips", Location: "form", Type: "string"},
				{Name: "imad", Location: "form", Type: "string"},
				{Name: "omad", Location: "form", Type: "string"},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/treasury/outbound_transfers/{outbound_transfer}/fail", OperationID: "PostTestHelpersTreasuryOutboundTransfersOutboundTransferFail", Params: []ParameterValidation{
		{Name: "outbound_transfer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/treasury/outbound_transfers/{outbound_transfer}/post", OperationID: "PostTestHelpersTreasuryOutboundTransfersOutboundTransferPost", Params: []ParameterValidation{
		{Name: "outbound_transfer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/treasury/outbound_transfers/{outbound_transfer}/return", OperationID: "PostTestHelpersTreasuryOutboundTransfersOutboundTransferReturn", Params: []ParameterValidation{
		{Name: "outbound_transfer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "returned_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "code", Location: "form", Type: "string", Enum: []string{"account_closed", "account_frozen", "bank_account_restricted", "bank_ownership_changed", "declined", "incorrect_account_holder_name", "invalid_account_number", "invalid_currency", "no_account", "other"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/treasury/received_credits", OperationID: "PostTestHelpersTreasuryReceivedCredits", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Required: true, Type: "integer"},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "financial_account", Location: "form", Required: true, Type: "string"},
		{Name: "initiating_payment_method_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"us_bank_account"}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_name", Location: "form", Type: "string"},
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "routing_number", Location: "form", Type: "string"},
			}},
		}},
		{Name: "network", Location: "form", Required: true, Type: "string", Enum: []string{"ach", "us_domestic_wire"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/test_helpers/treasury/received_debits", OperationID: "PostTestHelpersTreasuryReceivedDebits", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Required: true, Type: "integer"},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "financial_account", Location: "form", Required: true, Type: "string"},
		{Name: "initiating_payment_method_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"us_bank_account"}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_name", Location: "form", Type: "string"},
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "routing_number", Location: "form", Type: "string"},
			}},
		}},
		{Name: "network", Location: "form", Required: true, Type: "string", Enum: []string{"ach"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/tokens", OperationID: "PostTokens", Params: []ParameterValidation{
		{Name: "account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "business_type", Location: "form", Type: "string", Enum: []string{"company", "government_entity", "individual", "non_profit"}},
			{Name: "company", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object"},
				{Name: "address_kana", Location: "form", Type: "object"},
				{Name: "address_kanji", Location: "form", Type: "object"},
				{Name: "directors_provided", Location: "form", Type: "boolean"},
				{Name: "directorship_declaration", Location: "form", Type: "object"},
				{Name: "executives_provided", Location: "form", Type: "boolean"},
				{Name: "export_license_id", Location: "form", Type: "string"},
				{Name: "export_purpose_code", Location: "form", Type: "string"},
				{Name: "name", Location: "form", Type: "string"},
				{Name: "name_kana", Location: "form", Type: "string"},
				{Name: "name_kanji", Location: "form", Type: "string"},
				{Name: "owners_provided", Location: "form", Type: "boolean"},
				{Name: "ownership_declaration", Location: "form", Type: "object"},
				{Name: "ownership_declaration_shown_and_signed", Location: "form", Type: "boolean"},
				{Name: "ownership_exemption_reason", Location: "form", Type: "string", Enum: []string{"", "qualified_entity_exceeds_ownership_threshold", "qualifies_as_financial_institution"}},
				{Name: "phone", Location: "form", Type: "string"},
				{Name: "registration_date", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "registration_number", Location: "form", Type: "string"},
				{Name: "representative_declaration", Location: "form", Type: "object"},
				{Name: "structure", Location: "form", Type: "string", Enum: []string{"", "free_zone_establishment", "free_zone_llc", "government_instrumentality", "governmental_unit", "incorporated_non_profit", "incorporated_partnership", "limited_liability_partnership", "llc", "multi_member_llc", "private_company", "private_corporation", "private_partnership", "public_company", "public_corporation", "public_partnership", "registered_charity", "single_member_llc", "sole_establishment", "sole_proprietorship", "tax_exempt_government_instrumentality", "unincorporated_association", "unincorporated_non_profit", "unincorporated_partnership"}},
				{Name: "tax_id", Location: "form", Type: "string"},
				{Name: "tax_id_registrar", Location: "form", Type: "string"},
				{Name: "vat_id", Location: "form", Type: "string"},
				{Name: "verification", Location: "form", Type: "object"},
			}},
			{Name: "individual", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object"},
				{Name: "address_kana", Location: "form", Type: "object"},
				{Name: "address_kanji", Location: "form", Type: "object"},
				{Name: "dob", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "email", Location: "form", Type: "string"},
				{Name: "first_name", Location: "form", Type: "string"},
				{Name: "first_name_kana", Location: "form", Type: "string"},
				{Name: "first_name_kanji", Location: "form", Type: "string"},
				{Name: "full_name_aliases", Location: "form", Enum: []string{""}},
				{Name: "gender", Location: "form", Type: "string"},
				{Name: "id_number", Location: "form", Type: "string"},
				{Name: "id_number_secondary", Location: "form", Type: "string"},
				{Name: "last_name", Location: "form", Type: "string"},
				{Name: "last_name_kana", Location: "form", Type: "string"},
				{Name: "last_name_kanji", Location: "form", Type: "string"},
				{Name: "maiden_name", Location: "form", Type: "string"},
				{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "phone", Location: "form", Type: "string"},
				{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
				{Name: "registered_address", Location: "form", Type: "object"},
				{Name: "relationship", Location: "form", Type: "object"},
				{Name: "ssn_last_4", Location: "form", Type: "string"},
				{Name: "verification", Location: "form", Type: "object"},
			}},
			{Name: "tos_shown_and_accepted", Location: "form", Type: "boolean"},
		}},
		{Name: "bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account_holder_name", Location: "form", Type: "string"},
			{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
			{Name: "account_number", Location: "form", Required: true, Type: "string"},
			{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "futsu", "savings", "toza"}},
			{Name: "country", Location: "form", Required: true, Type: "string"},
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "payment_method", Location: "form", Type: "string"},
			{Name: "routing_number", Location: "form", Type: "string"},
		}},
		{Name: "card", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "address_city", Location: "form", Type: "string"},
			{Name: "address_country", Location: "form", Type: "string"},
			{Name: "address_line1", Location: "form", Type: "string"},
			{Name: "address_line2", Location: "form", Type: "string"},
			{Name: "address_state", Location: "form", Type: "string"},
			{Name: "address_zip", Location: "form", Type: "string"},
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "cvc", Location: "form", Type: "string"},
			{Name: "exp_month", Location: "form", Required: true, Type: "string"},
			{Name: "exp_year", Location: "form", Required: true, Type: "string"},
			{Name: "name", Location: "form", Type: "string"},
			{Name: "networks", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "preferred", Location: "form", Type: "string", Enum: []string{"cartes_bancaires", "mastercard", "visa"}},
			}},
			{Name: "number", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "cvc_update", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "cvc", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "person", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "additional_tos_acceptances", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account", Location: "form", Type: "object"},
			}},
			{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "address_kana", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
			{Name: "address_kanji", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
			{Name: "dob", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "day", Location: "form", Required: true, Type: "integer"},
				{Name: "month", Location: "form", Required: true, Type: "integer"},
				{Name: "year", Location: "form", Required: true, Type: "integer"},
			}},
			{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "company_authorization", Location: "form", Type: "object"},
				{Name: "passport", Location: "form", Type: "object"},
				{Name: "visa", Location: "form", Type: "object"},
			}},
			{Name: "email", Location: "form", Type: "string"},
			{Name: "first_name", Location: "form", Type: "string"},
			{Name: "first_name_kana", Location: "form", Type: "string"},
			{Name: "first_name_kanji", Location: "form", Type: "string"},
			{Name: "full_name_aliases", Location: "form", Enum: []string{""}},
			{Name: "gender", Location: "form", Type: "string"},
			{Name: "id_number", Location: "form", Type: "string"},
			{Name: "id_number_secondary", Location: "form", Type: "string"},
			{Name: "last_name", Location: "form", Type: "string"},
			{Name: "last_name_kana", Location: "form", Type: "string"},
			{Name: "last_name_kanji", Location: "form", Type: "string"},
			{Name: "maiden_name", Location: "form", Type: "string"},
			{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
			{Name: "nationality", Location: "form", Type: "string"},
			{Name: "phone", Location: "form", Type: "string"},
			{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
			{Name: "registered_address", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
			}},
			{Name: "relationship", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "authorizer", Location: "form", Type: "boolean"},
				{Name: "director", Location: "form", Type: "boolean"},
				{Name: "executive", Location: "form", Type: "boolean"},
				{Name: "legal_guardian", Location: "form", Type: "boolean"},
				{Name: "owner", Location: "form", Type: "boolean"},
				{Name: "percent_ownership", Location: "form", Enum: []string{""}},
				{Name: "representative", Location: "form", Type: "boolean"},
				{Name: "title", Location: "form", Type: "string"},
			}},
			{Name: "ssn_last_4", Location: "form", Type: "string"},
			{Name: "us_cfpb_data", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ethnicity_details", Location: "form", Type: "object"},
				{Name: "race_details", Location: "form", Type: "object"},
				{Name: "self_identified_gender", Location: "form", Type: "string"},
			}},
			{Name: "verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "additional_document", Location: "form", Type: "object"},
				{Name: "document", Location: "form", Type: "object"},
			}},
		}},
		{Name: "pii", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "id_number", Location: "form", Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/tokens/{token}", OperationID: "GetTokensToken", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "token", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/topups", OperationID: "GetTopups", Params: []ParameterValidation{
		{Name: "amount", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"canceled", "failed", "pending", "succeeded"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/topups", OperationID: "PostTopups", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Required: true, Type: "integer"},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "source", Location: "form", Type: "string"},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
		{Name: "transfer_group", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/topups/{topup}", OperationID: "GetTopupsTopup", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "topup", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/topups/{topup}", OperationID: "PostTopupsTopup", Params: []ParameterValidation{
		{Name: "topup", Location: "path", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/topups/{topup}/cancel", OperationID: "PostTopupsTopupCancel", Params: []ParameterValidation{
		{Name: "topup", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/transfers", OperationID: "GetTransfers", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "destination", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "transfer_group", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/transfers", OperationID: "PostTransfers", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "destination", Location: "form", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "source_transaction", Location: "form", Type: "string"},
		{Name: "source_type", Location: "form", Type: "string", Enum: []string{"bank_account", "card", "fpx"}},
		{Name: "transfer_group", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/transfers/{id}/reversals", OperationID: "GetTransfersIdReversals", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/transfers/{id}/reversals", OperationID: "PostTransfersIdReversals", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "amount", Location: "form", Type: "integer"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "refund_application_fee", Location: "form", Type: "boolean"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/transfers/{transfer}", OperationID: "GetTransfersTransfer", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "transfer", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/transfers/{transfer}", OperationID: "PostTransfersTransfer", Params: []ParameterValidation{
		{Name: "transfer", Location: "path", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/transfers/{transfer}/reversals/{id}", OperationID: "GetTransfersTransferReversalsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "transfer", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/transfers/{transfer}/reversals/{id}", OperationID: "PostTransfersTransferReversalsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "transfer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/credit_reversals", OperationID: "GetTreasuryCreditReversals", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "financial_account", Location: "query", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "received_credit", Location: "query", Type: "string"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"canceled", "posted", "processing"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/treasury/credit_reversals", OperationID: "PostTreasuryCreditReversals", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "received_credit", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/credit_reversals/{credit_reversal}", OperationID: "GetTreasuryCreditReversalsCreditReversal", Params: []ParameterValidation{
		{Name: "credit_reversal", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/debit_reversals", OperationID: "GetTreasuryDebitReversals", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "financial_account", Location: "query", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "received_debit", Location: "query", Type: "string"},
		{Name: "resolution", Location: "query", Type: "string", Enum: []string{"lost", "won"}},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"canceled", "completed", "processing"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/treasury/debit_reversals", OperationID: "PostTreasuryDebitReversals", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "received_debit", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/debit_reversals/{debit_reversal}", OperationID: "GetTreasuryDebitReversalsDebitReversal", Params: []ParameterValidation{
		{Name: "debit_reversal", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/financial_accounts", OperationID: "GetTreasuryFinancialAccounts", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"closed", "open"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/treasury/financial_accounts", OperationID: "PostTreasuryFinancialAccounts", Params: []ParameterValidation{
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "features", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "card_issuing", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Required: true, Type: "boolean"},
			}},
			{Name: "deposit_insurance", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Required: true, Type: "boolean"},
			}},
			{Name: "financial_addresses", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "aba", Location: "form", Type: "object"},
			}},
			{Name: "inbound_transfers", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ach", Location: "form", Type: "object"},
			}},
			{Name: "intra_stripe_flows", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Required: true, Type: "boolean"},
			}},
			{Name: "outbound_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ach", Location: "form", Type: "object"},
				{Name: "us_domestic_wire", Location: "form", Type: "object"},
			}},
			{Name: "outbound_transfers", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ach", Location: "form", Type: "object"},
				{Name: "us_domestic_wire", Location: "form", Type: "object"},
			}},
		}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "nickname", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "platform_restrictions", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "inbound_flows", Location: "form", Type: "string", Enum: []string{"restricted", "unrestricted"}},
			{Name: "outbound_flows", Location: "form", Type: "string", Enum: []string{"restricted", "unrestricted"}},
		}},
		{Name: "supported_currencies", Location: "form", Required: true, Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/financial_accounts/{financial_account}", OperationID: "GetTreasuryFinancialAccountsFinancialAccount", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "financial_account", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/treasury/financial_accounts/{financial_account}", OperationID: "PostTreasuryFinancialAccountsFinancialAccount", Params: []ParameterValidation{
		{Name: "financial_account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "features", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "card_issuing", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Required: true, Type: "boolean"},
			}},
			{Name: "deposit_insurance", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Required: true, Type: "boolean"},
			}},
			{Name: "financial_addresses", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "aba", Location: "form", Type: "object"},
			}},
			{Name: "inbound_transfers", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ach", Location: "form", Type: "object"},
			}},
			{Name: "intra_stripe_flows", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Required: true, Type: "boolean"},
			}},
			{Name: "outbound_payments", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ach", Location: "form", Type: "object"},
				{Name: "us_domestic_wire", Location: "form", Type: "object"},
			}},
			{Name: "outbound_transfers", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "ach", Location: "form", Type: "object"},
				{Name: "us_domestic_wire", Location: "form", Type: "object"},
			}},
		}},
		{Name: "forwarding_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "financial_account", Location: "form", Type: "string"},
			{Name: "payment_method", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"financial_account", "payment_method"}},
		}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "nickname", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "platform_restrictions", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "inbound_flows", Location: "form", Type: "string", Enum: []string{"restricted", "unrestricted"}},
			{Name: "outbound_flows", Location: "form", Type: "string", Enum: []string{"restricted", "unrestricted"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/treasury/financial_accounts/{financial_account}/close", OperationID: "PostTreasuryFinancialAccountsFinancialAccountClose", Params: []ParameterValidation{
		{Name: "financial_account", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "forwarding_settings", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "financial_account", Location: "form", Type: "string"},
			{Name: "payment_method", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"financial_account", "payment_method"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/financial_accounts/{financial_account}/features", OperationID: "GetTreasuryFinancialAccountsFinancialAccountFeatures", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "financial_account", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/treasury/financial_accounts/{financial_account}/features", OperationID: "PostTreasuryFinancialAccountsFinancialAccountFeatures", Params: []ParameterValidation{
		{Name: "financial_account", Location: "path", Required: true, Type: "string"},
		{Name: "card_issuing", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "requested", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "deposit_insurance", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "requested", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "financial_addresses", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "aba", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Required: true, Type: "boolean"},
			}},
		}},
		{Name: "inbound_transfers", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "ach", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Required: true, Type: "boolean"},
			}},
		}},
		{Name: "intra_stripe_flows", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "requested", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "outbound_payments", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "ach", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Required: true, Type: "boolean"},
			}},
			{Name: "us_domestic_wire", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Required: true, Type: "boolean"},
			}},
		}},
		{Name: "outbound_transfers", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "ach", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Required: true, Type: "boolean"},
			}},
			{Name: "us_domestic_wire", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "requested", Location: "form", Required: true, Type: "boolean"},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/inbound_transfers", OperationID: "GetTreasuryInboundTransfers", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "financial_account", Location: "query", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"canceled", "failed", "processing", "succeeded"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/treasury/inbound_transfers", OperationID: "PostTreasuryInboundTransfers", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Required: true, Type: "integer"},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "financial_account", Location: "form", Required: true, Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "origin_payment_method", Location: "form", Required: true, Type: "string"},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/inbound_transfers/{id}", OperationID: "GetTreasuryInboundTransfersId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/treasury/inbound_transfers/{inbound_transfer}/cancel", OperationID: "PostTreasuryInboundTransfersInboundTransferCancel", Params: []ParameterValidation{
		{Name: "inbound_transfer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/outbound_payments", OperationID: "GetTreasuryOutboundPayments", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "customer", Location: "query", Type: "string"},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "financial_account", Location: "query", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"canceled", "failed", "posted", "processing", "returned"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/treasury/outbound_payments", OperationID: "PostTreasuryOutboundPayments", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Required: true, Type: "integer"},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "customer", Location: "form", Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "destination_payment_method", Location: "form", Type: "string"},
		{Name: "destination_payment_method_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "billing_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object", Enum: []string{""}},
				{Name: "email", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "name", Location: "form", Type: "string", Enum: []string{""}},
				{Name: "phone", Location: "form", Type: "string", Enum: []string{""}},
			}},
			{Name: "financial_account", Location: "form", Type: "string"},
			{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"financial_account", "us_bank_account"}},
			{Name: "us_bank_account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "account_holder_type", Location: "form", Type: "string", Enum: []string{"company", "individual"}},
				{Name: "account_number", Location: "form", Type: "string"},
				{Name: "account_type", Location: "form", Type: "string", Enum: []string{"checking", "savings"}},
				{Name: "financial_connections_account", Location: "form", Type: "string"},
				{Name: "routing_number", Location: "form", Type: "string"},
			}},
		}},
		{Name: "destination_payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "us_bank_account", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "network", Location: "form", Type: "string", Enum: []string{"ach", "us_domestic_wire"}},
			}},
		}},
		{Name: "end_user_details", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "ip_address", Location: "form", Type: "string"},
			{Name: "present", Location: "form", Required: true, Type: "boolean"},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "financial_account", Location: "form", Required: true, Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/outbound_payments/{id}", OperationID: "GetTreasuryOutboundPaymentsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/treasury/outbound_payments/{id}/cancel", OperationID: "PostTreasuryOutboundPaymentsIdCancel", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/outbound_transfers", OperationID: "GetTreasuryOutboundTransfers", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "financial_account", Location: "query", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"canceled", "failed", "posted", "processing", "returned"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/treasury/outbound_transfers", OperationID: "PostTreasuryOutboundTransfers", Params: []ParameterValidation{
		{Name: "amount", Location: "form", Required: true, Type: "integer"},
		{Name: "currency", Location: "form", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "destination_payment_method", Location: "form", Type: "string"},
		{Name: "destination_payment_method_data", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "financial_account", Location: "form", Type: "string"},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"financial_account"}},
		}},
		{Name: "destination_payment_method_options", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "us_bank_account", Location: "form", Type: "object", Enum: []string{""}, Children: []ParameterValidation{
				{Name: "network", Location: "form", Type: "string", Enum: []string{"ach", "us_domestic_wire"}},
			}},
		}},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "financial_account", Location: "form", Required: true, Type: "string"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "statement_descriptor", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/outbound_transfers/{outbound_transfer}", OperationID: "GetTreasuryOutboundTransfersOutboundTransfer", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "outbound_transfer", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/treasury/outbound_transfers/{outbound_transfer}/cancel", OperationID: "PostTreasuryOutboundTransfersOutboundTransferCancel", Params: []ParameterValidation{
		{Name: "outbound_transfer", Location: "path", Required: true, Type: "string"},
		{Name: "expand", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/received_credits", OperationID: "GetTreasuryReceivedCredits", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "financial_account", Location: "query", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "linked_flows", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "source_flow_type", Location: "query", Required: true, Type: "string", Enum: []string{"credit_reversal", "other", "outbound_payment", "outbound_transfer", "payout"}},
		}},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"failed", "succeeded"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/received_credits/{id}", OperationID: "GetTreasuryReceivedCreditsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/received_debits", OperationID: "GetTreasuryReceivedDebits", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "financial_account", Location: "query", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"failed", "succeeded"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/received_debits/{id}", OperationID: "GetTreasuryReceivedDebitsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/transaction_entries", OperationID: "GetTreasuryTransactionEntries", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "effective_at", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "financial_account", Location: "query", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "order_by", Location: "query", Type: "string", Enum: []string{"created", "effective_at"}},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "transaction", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/transaction_entries/{id}", OperationID: "GetTreasuryTransactionEntriesId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/transactions", OperationID: "GetTreasuryTransactions", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "integer"},
			{Name: "gte", Location: "query", Type: "integer"},
			{Name: "lt", Location: "query", Type: "integer"},
			{Name: "lte", Location: "query", Type: "integer"},
		}},
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "financial_account", Location: "query", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "order_by", Location: "query", Type: "string", Enum: []string{"created", "posted_at"}},
		{Name: "starting_after", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"open", "posted", "void"}},
		{Name: "status_transitions", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "posted_at", Location: "query", Type: "object", Children: []ParameterValidation{
				{Name: "gt", Location: "query", Type: "integer"},
				{Name: "gte", Location: "query", Type: "integer"},
				{Name: "lt", Location: "query", Type: "integer"},
				{Name: "lte", Location: "query", Type: "integer"},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/treasury/transactions/{id}", OperationID: "GetTreasuryTransactionsId", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/webhook_endpoints", OperationID: "GetWebhookEndpoints", Params: []ParameterValidation{
		{Name: "ending_before", Location: "query", Type: "string"},
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "starting_after", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/webhook_endpoints", OperationID: "PostWebhookEndpoints", Params: []ParameterValidation{
		{Name: "api_version", Location: "form", Type: "string"},
		{Name: "connect", Location: "form", Type: "boolean"},
		{Name: "description", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "enabled_events", Location: "form", Required: true, Type: "array"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "url", Location: "form", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v1/webhook_endpoints/{webhook_endpoint}", OperationID: "DeleteWebhookEndpointsWebhookEndpoint", Params: []ParameterValidation{
		{Name: "webhook_endpoint", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v1/webhook_endpoints/{webhook_endpoint}", OperationID: "GetWebhookEndpointsWebhookEndpoint", Params: []ParameterValidation{
		{Name: "expand", Location: "query", Type: "array"},
		{Name: "webhook_endpoint", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v1/webhook_endpoints/{webhook_endpoint}", OperationID: "PostWebhookEndpointsWebhookEndpoint", Params: []ParameterValidation{
		{Name: "webhook_endpoint", Location: "path", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string", Enum: []string{""}},
		{Name: "disabled", Location: "form", Type: "boolean"},
		{Name: "enabled_events", Location: "form", Type: "array"},
		{Name: "expand", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", Enum: []string{""}},
		{Name: "url", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/billing/meter_event_adjustments", OperationID: "PostV2BillingMeterEventAdjustments", Params: []ParameterValidation{
		{Name: "cancel", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "identifier", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "event_name", Location: "form", Required: true, Type: "string"},
		{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"cancel"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/billing/meter_event_session", OperationID: "PostV2BillingMeterEventSession"})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/billing/meter_event_stream", OperationID: "PostV2BillingMeterEventStream", Params: []ParameterValidation{
		{Name: "events", Location: "form", Required: true, Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/billing/meter_events", OperationID: "PostV2BillingMeterEvents", Params: []ParameterValidation{
		{Name: "event_name", Location: "form", Required: true, Type: "string"},
		{Name: "identifier", Location: "form", Type: "string"},
		{Name: "payload", Location: "form", Required: true, Type: "object", AdditionalProperties: true},
		{Name: "timestamp", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v2/commerce/product_catalog/imports", OperationID: "GetV2CommerceProductCatalogImports", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "string"},
		{Name: "created_gt", Location: "query", Type: "string"},
		{Name: "created_gte", Location: "query", Type: "string"},
		{Name: "created_lt", Location: "query", Type: "string"},
		{Name: "created_lte", Location: "query", Type: "string"},
		{Name: "feed_type", Location: "query", Type: "string", Enum: []string{"inventory", "pricing", "product"}},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "page", Location: "query", Type: "string"},
		{Name: "status", Location: "query", Type: "string", Enum: []string{"awaiting_upload", "failed", "processing", "succeeded", "succeeded_with_errors"}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/commerce/product_catalog/imports", OperationID: "PostV2CommerceProductCatalogImports", Params: []ParameterValidation{
		{Name: "feed_type", Location: "form", Required: true, Type: "string", Enum: []string{"inventory", "pricing", "product"}},
		{Name: "metadata", Location: "form", Required: true, Type: "object", AdditionalProperties: true},
		{Name: "mode", Location: "form", Required: true, Type: "string", Enum: []string{"replace", "upsert"}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v2/commerce/product_catalog/imports/{id}", OperationID: "GetV2CommerceProductCatalogImportsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/core/account_links", OperationID: "PostV2CoreAccountLinks", Params: []ParameterValidation{
		{Name: "account", Location: "form", Required: true, Type: "string"},
		{Name: "use_case", Location: "form", Required: true, Type: "object", Children: []ParameterValidation{
			{Name: "account_onboarding", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "collection_options", Location: "form", Type: "object"},
				{Name: "configurations", Location: "form", Required: true, Type: "array"},
				{Name: "refresh_url", Location: "form", Required: true, Type: "string"},
				{Name: "return_url", Location: "form", Type: "string"},
			}},
			{Name: "account_update", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "collection_options", Location: "form", Type: "object"},
				{Name: "configurations", Location: "form", Required: true, Type: "array"},
				{Name: "refresh_url", Location: "form", Required: true, Type: "string"},
				{Name: "return_url", Location: "form", Type: "string"},
			}},
			{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"account_onboarding", "account_update"}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/core/account_tokens", OperationID: "PostV2CoreAccountTokens", Params: []ParameterValidation{
		{Name: "contact_email", Location: "form", Type: "string"},
		{Name: "contact_phone", Location: "form", Type: "string"},
		{Name: "display_name", Location: "form", Type: "string"},
		{Name: "identity", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "attestations", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "directorship_declaration", Location: "form", Type: "object"},
				{Name: "ownership_declaration", Location: "form", Type: "object"},
				{Name: "persons_provided", Location: "form", Type: "object"},
				{Name: "representative_declaration", Location: "form", Type: "object"},
				{Name: "terms_of_service", Location: "form", Type: "object"},
			}},
			{Name: "business_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object"},
				{Name: "annual_revenue", Location: "form", Type: "object"},
				{Name: "documents", Location: "form", Type: "object"},
				{Name: "estimated_worker_count", Location: "form", Type: "integer"},
				{Name: "id_numbers", Location: "form", Type: "array"},
				{Name: "monthly_estimated_revenue", Location: "form", Type: "object"},
				{Name: "phone", Location: "form", Type: "string"},
				{Name: "registered_name", Location: "form", Type: "string"},
				{Name: "registration_date", Location: "form", Type: "object"},
				{Name: "script_addresses", Location: "form", Type: "object"},
				{Name: "script_names", Location: "form", Type: "object"},
				{Name: "structure", Location: "form", Type: "string", Enum: []string{"cooperative", "free_zone_establishment", "free_zone_llc", "government_instrumentality", "governmental_unit", "incorporated_association", "incorporated_non_profit", "incorporated_partnership", "limited_liability_partnership", "llc", "multi_member_llc", "private_company", "private_corporation", "private_partnership", "public_company", "public_corporation", "public_listed_corporation", "public_partnership", "registered_charity", "single_member_llc", "sole_establishment", "sole_proprietorship", "tax_exempt_government_instrumentality", "trust", "unincorporated_association", "unincorporated_non_profit", "unincorporated_partnership"}},
			}},
			{Name: "entity_type", Location: "form", Type: "string", Enum: []string{"company", "government_entity", "individual", "non_profit"}},
			{Name: "individual", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "additional_addresses", Location: "form", Type: "array"},
				{Name: "additional_names", Location: "form", Type: "array"},
				{Name: "address", Location: "form", Type: "object"},
				{Name: "date_of_birth", Location: "form", Type: "object"},
				{Name: "documents", Location: "form", Type: "object"},
				{Name: "email", Location: "form", Type: "string"},
				{Name: "given_name", Location: "form", Type: "string"},
				{Name: "id_numbers", Location: "form", Type: "array"},
				{Name: "legal_gender", Location: "form", Type: "string", Enum: []string{"female", "male"}},
				{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
				{Name: "nationalities", Location: "form", Type: "array"},
				{Name: "phone", Location: "form", Type: "string"},
				{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
				{Name: "relationship", Location: "form", Type: "object"},
				{Name: "script_addresses", Location: "form", Type: "object"},
				{Name: "script_names", Location: "form", Type: "object"},
				{Name: "surname", Location: "form", Type: "string"},
			}},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v2/core/account_tokens/{id}", OperationID: "GetV2CoreAccountTokensId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v2/core/accounts", OperationID: "GetV2CoreAccounts", Params: []ParameterValidation{
		{Name: "applied_configurations", Location: "query", Type: "array"},
		{Name: "closed", Location: "query", Type: "boolean"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "page", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/core/accounts", OperationID: "PostV2CoreAccounts", Params: []ParameterValidation{
		{Name: "account_token", Location: "form", Type: "string"},
		{Name: "configuration", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "customer", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "automatic_indirect_tax", Location: "form", Type: "object"},
				{Name: "billing", Location: "form", Type: "object"},
				{Name: "capabilities", Location: "form", Type: "object"},
				{Name: "shipping", Location: "form", Type: "object"},
				{Name: "test_clock", Location: "form", Type: "string"},
			}},
			{Name: "merchant", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "bacs_debit_payments", Location: "form", Type: "object"},
				{Name: "branding", Location: "form", Type: "object"},
				{Name: "capabilities", Location: "form", Type: "object"},
				{Name: "card_payments", Location: "form", Type: "object"},
				{Name: "konbini_payments", Location: "form", Type: "object"},
				{Name: "mcc", Location: "form", Type: "string"},
				{Name: "script_statement_descriptor", Location: "form", Type: "object"},
				{Name: "statement_descriptor", Location: "form", Type: "object"},
				{Name: "support", Location: "form", Type: "object"},
			}},
			{Name: "recipient", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "capabilities", Location: "form", Type: "object"},
			}},
		}},
		{Name: "contact_email", Location: "form", Type: "string"},
		{Name: "contact_phone", Location: "form", Type: "string"},
		{Name: "dashboard", Location: "form", Type: "string", Enum: []string{"express", "full", "none"}},
		{Name: "defaults", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "locales", Location: "form", Type: "array"},
			{Name: "profile", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "business_url", Location: "form", Type: "string"},
				{Name: "doing_business_as", Location: "form", Type: "string"},
				{Name: "product_description", Location: "form", Type: "string"},
			}},
			{Name: "responsibilities", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fees_collector", Location: "form", Required: true, Type: "string", Enum: []string{"application", "application_custom", "application_express", "stripe"}},
				{Name: "losses_collector", Location: "form", Required: true, Type: "string", Enum: []string{"application", "stripe"}},
			}},
		}},
		{Name: "display_name", Location: "form", Type: "string"},
		{Name: "identity", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "attestations", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "directorship_declaration", Location: "form", Type: "object"},
				{Name: "ownership_declaration", Location: "form", Type: "object"},
				{Name: "persons_provided", Location: "form", Type: "object"},
				{Name: "representative_declaration", Location: "form", Type: "object"},
				{Name: "terms_of_service", Location: "form", Type: "object"},
			}},
			{Name: "business_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object"},
				{Name: "annual_revenue", Location: "form", Type: "object"},
				{Name: "documents", Location: "form", Type: "object"},
				{Name: "estimated_worker_count", Location: "form", Type: "integer"},
				{Name: "id_numbers", Location: "form", Type: "array"},
				{Name: "monthly_estimated_revenue", Location: "form", Type: "object"},
				{Name: "phone", Location: "form", Type: "string"},
				{Name: "registered_name", Location: "form", Type: "string"},
				{Name: "registration_date", Location: "form", Type: "object"},
				{Name: "script_addresses", Location: "form", Type: "object"},
				{Name: "script_names", Location: "form", Type: "object"},
				{Name: "structure", Location: "form", Type: "string", Enum: []string{"cooperative", "free_zone_establishment", "free_zone_llc", "government_instrumentality", "governmental_unit", "incorporated_association", "incorporated_non_profit", "incorporated_partnership", "limited_liability_partnership", "llc", "multi_member_llc", "private_company", "private_corporation", "private_partnership", "public_company", "public_corporation", "public_listed_corporation", "public_partnership", "registered_charity", "single_member_llc", "sole_establishment", "sole_proprietorship", "tax_exempt_government_instrumentality", "trust", "unincorporated_association", "unincorporated_non_profit", "unincorporated_partnership"}},
			}},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "entity_type", Location: "form", Type: "string", Enum: []string{"company", "government_entity", "individual", "non_profit"}},
			{Name: "individual", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "additional_addresses", Location: "form", Type: "array"},
				{Name: "additional_names", Location: "form", Type: "array"},
				{Name: "address", Location: "form", Type: "object"},
				{Name: "date_of_birth", Location: "form", Type: "object"},
				{Name: "documents", Location: "form", Type: "object"},
				{Name: "email", Location: "form", Type: "string"},
				{Name: "given_name", Location: "form", Type: "string"},
				{Name: "id_numbers", Location: "form", Type: "array"},
				{Name: "legal_gender", Location: "form", Type: "string", Enum: []string{"female", "male"}},
				{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
				{Name: "nationalities", Location: "form", Type: "array"},
				{Name: "phone", Location: "form", Type: "string"},
				{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
				{Name: "relationship", Location: "form", Type: "object"},
				{Name: "script_addresses", Location: "form", Type: "object"},
				{Name: "script_names", Location: "form", Type: "object"},
				{Name: "surname", Location: "form", Type: "string"},
			}},
		}},
		{Name: "include", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/core/accounts/{account_id}/person_tokens", OperationID: "PostV2CoreAccountsAccountIdPersonTokens", Params: []ParameterValidation{
		{Name: "account_id", Location: "path", Required: true, Type: "string"},
		{Name: "additional_addresses", Location: "form", Type: "array"},
		{Name: "additional_names", Location: "form", Type: "array"},
		{Name: "additional_terms_of_service", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "shown_and_accepted", Location: "form", Type: "boolean"},
			}},
		}},
		{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "date_of_birth", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "day", Location: "form", Required: true, Type: "integer"},
			{Name: "month", Location: "form", Required: true, Type: "integer"},
			{Name: "year", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "company_authorization", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Required: true, Type: "array"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"files"}},
			}},
			{Name: "passport", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Required: true, Type: "array"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"files"}},
			}},
			{Name: "primary_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "front_back", Location: "form", Required: true, Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"front_back"}},
			}},
			{Name: "secondary_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "front_back", Location: "form", Required: true, Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"front_back"}},
			}},
			{Name: "visa", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Required: true, Type: "array"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"files"}},
			}},
		}},
		{Name: "email", Location: "form", Type: "string"},
		{Name: "given_name", Location: "form", Type: "string"},
		{Name: "id_numbers", Location: "form", Type: "array"},
		{Name: "legal_gender", Location: "form", Type: "string", Enum: []string{"female", "male"}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "nationalities", Location: "form", Type: "array"},
		{Name: "phone", Location: "form", Type: "string"},
		{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
		{Name: "relationship", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "authorizer", Location: "form", Type: "boolean"},
			{Name: "director", Location: "form", Type: "boolean"},
			{Name: "executive", Location: "form", Type: "boolean"},
			{Name: "legal_guardian", Location: "form", Type: "boolean"},
			{Name: "owner", Location: "form", Type: "boolean"},
			{Name: "percent_ownership", Location: "form", Type: "string"},
			{Name: "representative", Location: "form", Type: "boolean"},
			{Name: "title", Location: "form", Type: "string"},
		}},
		{Name: "script_addresses", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "kana", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
			{Name: "kanji", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
		}},
		{Name: "script_names", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "kana", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "given_name", Location: "form", Type: "string"},
				{Name: "surname", Location: "form", Type: "string"},
			}},
			{Name: "kanji", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "given_name", Location: "form", Type: "string"},
				{Name: "surname", Location: "form", Type: "string"},
			}},
		}},
		{Name: "surname", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v2/core/accounts/{account_id}/person_tokens/{id}", OperationID: "GetV2CoreAccountsAccountIdPersonTokensId", Params: []ParameterValidation{
		{Name: "account_id", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v2/core/accounts/{account_id}/persons", OperationID: "GetV2CoreAccountsAccountIdPersons", Params: []ParameterValidation{
		{Name: "account_id", Location: "path", Required: true, Type: "string"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "page", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/core/accounts/{account_id}/persons", OperationID: "PostV2CoreAccountsAccountIdPersons", Params: []ParameterValidation{
		{Name: "account_id", Location: "path", Required: true, Type: "string"},
		{Name: "additional_addresses", Location: "form", Type: "array"},
		{Name: "additional_names", Location: "form", Type: "array"},
		{Name: "additional_terms_of_service", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Required: true, Type: "string"},
				{Name: "ip", Location: "form", Required: true, Type: "string"},
				{Name: "user_agent", Location: "form", Type: "string"},
			}},
		}},
		{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Required: true, Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "date_of_birth", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "day", Location: "form", Required: true, Type: "integer"},
			{Name: "month", Location: "form", Required: true, Type: "integer"},
			{Name: "year", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "company_authorization", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Required: true, Type: "array"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"files"}},
			}},
			{Name: "passport", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Required: true, Type: "array"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"files"}},
			}},
			{Name: "primary_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "front_back", Location: "form", Required: true, Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"front_back"}},
			}},
			{Name: "secondary_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "front_back", Location: "form", Required: true, Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"front_back"}},
			}},
			{Name: "visa", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Required: true, Type: "array"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"files"}},
			}},
		}},
		{Name: "email", Location: "form", Type: "string"},
		{Name: "given_name", Location: "form", Type: "string"},
		{Name: "id_numbers", Location: "form", Type: "array"},
		{Name: "legal_gender", Location: "form", Type: "string", Enum: []string{"female", "male"}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "nationalities", Location: "form", Type: "array"},
		{Name: "person_token", Location: "form", Type: "string"},
		{Name: "phone", Location: "form", Type: "string"},
		{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
		{Name: "relationship", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "authorizer", Location: "form", Type: "boolean"},
			{Name: "director", Location: "form", Type: "boolean"},
			{Name: "executive", Location: "form", Type: "boolean"},
			{Name: "legal_guardian", Location: "form", Type: "boolean"},
			{Name: "owner", Location: "form", Type: "boolean"},
			{Name: "percent_ownership", Location: "form", Type: "string"},
			{Name: "representative", Location: "form", Type: "boolean"},
			{Name: "title", Location: "form", Type: "string"},
		}},
		{Name: "script_addresses", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "kana", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Required: true, Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
			{Name: "kanji", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Required: true, Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
		}},
		{Name: "script_names", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "kana", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "given_name", Location: "form", Type: "string"},
				{Name: "surname", Location: "form", Type: "string"},
			}},
			{Name: "kanji", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "given_name", Location: "form", Type: "string"},
				{Name: "surname", Location: "form", Type: "string"},
			}},
		}},
		{Name: "surname", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v2/core/accounts/{account_id}/persons/{id}", OperationID: "DeleteV2CoreAccountsAccountIdPersonsId", Params: []ParameterValidation{
		{Name: "account_id", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v2/core/accounts/{account_id}/persons/{id}", OperationID: "GetV2CoreAccountsAccountIdPersonsId", Params: []ParameterValidation{
		{Name: "account_id", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/core/accounts/{account_id}/persons/{id}", OperationID: "PostV2CoreAccountsAccountIdPersonsId", Params: []ParameterValidation{
		{Name: "account_id", Location: "path", Required: true, Type: "string"},
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "additional_addresses", Location: "form", Type: "array"},
		{Name: "additional_names", Location: "form", Type: "array"},
		{Name: "additional_terms_of_service", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "account", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "date", Location: "form", Type: "string"},
				{Name: "ip", Location: "form", Type: "string"},
				{Name: "user_agent", Location: "form", Type: "string"},
			}},
		}},
		{Name: "address", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "city", Location: "form", Type: "string"},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "line1", Location: "form", Type: "string"},
			{Name: "line2", Location: "form", Type: "string"},
			{Name: "postal_code", Location: "form", Type: "string"},
			{Name: "state", Location: "form", Type: "string"},
			{Name: "town", Location: "form", Type: "string"},
		}},
		{Name: "date_of_birth", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "day", Location: "form", Required: true, Type: "integer"},
			{Name: "month", Location: "form", Required: true, Type: "integer"},
			{Name: "year", Location: "form", Required: true, Type: "integer"},
		}},
		{Name: "documents", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "company_authorization", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Required: true, Type: "array"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"files"}},
			}},
			{Name: "passport", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Required: true, Type: "array"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"files"}},
			}},
			{Name: "primary_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "front_back", Location: "form", Required: true, Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"front_back"}},
			}},
			{Name: "secondary_verification", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "front_back", Location: "form", Required: true, Type: "object"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"front_back"}},
			}},
			{Name: "visa", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "files", Location: "form", Required: true, Type: "array"},
				{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"files"}},
			}},
		}},
		{Name: "email", Location: "form", Type: "string"},
		{Name: "given_name", Location: "form", Type: "string"},
		{Name: "id_numbers", Location: "form", Type: "array"},
		{Name: "legal_gender", Location: "form", Type: "string", Enum: []string{"female", "male"}},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "nationalities", Location: "form", Type: "array"},
		{Name: "person_token", Location: "form", Type: "string"},
		{Name: "phone", Location: "form", Type: "string"},
		{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
		{Name: "relationship", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "authorizer", Location: "form", Type: "boolean"},
			{Name: "director", Location: "form", Type: "boolean"},
			{Name: "executive", Location: "form", Type: "boolean"},
			{Name: "legal_guardian", Location: "form", Type: "boolean"},
			{Name: "owner", Location: "form", Type: "boolean"},
			{Name: "percent_ownership", Location: "form", Type: "string"},
			{Name: "representative", Location: "form", Type: "boolean"},
			{Name: "title", Location: "form", Type: "string"},
		}},
		{Name: "script_addresses", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "kana", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
			{Name: "kanji", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "city", Location: "form", Type: "string"},
				{Name: "country", Location: "form", Type: "string"},
				{Name: "line1", Location: "form", Type: "string"},
				{Name: "line2", Location: "form", Type: "string"},
				{Name: "postal_code", Location: "form", Type: "string"},
				{Name: "state", Location: "form", Type: "string"},
				{Name: "town", Location: "form", Type: "string"},
			}},
		}},
		{Name: "script_names", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "kana", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "given_name", Location: "form", Type: "string"},
				{Name: "surname", Location: "form", Type: "string"},
			}},
			{Name: "kanji", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "given_name", Location: "form", Type: "string"},
				{Name: "surname", Location: "form", Type: "string"},
			}},
		}},
		{Name: "surname", Location: "form", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v2/core/accounts/{id}", OperationID: "GetV2CoreAccountsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "include", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/core/accounts/{id}", OperationID: "PostV2CoreAccountsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "account_token", Location: "form", Type: "string"},
		{Name: "configuration", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "customer", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "applied", Location: "form", Type: "boolean"},
				{Name: "automatic_indirect_tax", Location: "form", Type: "object"},
				{Name: "billing", Location: "form", Type: "object"},
				{Name: "capabilities", Location: "form", Type: "object"},
				{Name: "shipping", Location: "form", Type: "object"},
				{Name: "test_clock", Location: "form", Type: "string"},
			}},
			{Name: "merchant", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "applied", Location: "form", Type: "boolean"},
				{Name: "bacs_debit_payments", Location: "form", Type: "object"},
				{Name: "branding", Location: "form", Type: "object"},
				{Name: "capabilities", Location: "form", Type: "object"},
				{Name: "card_payments", Location: "form", Type: "object"},
				{Name: "konbini_payments", Location: "form", Type: "object"},
				{Name: "mcc", Location: "form", Type: "string"},
				{Name: "script_statement_descriptor", Location: "form", Type: "object"},
				{Name: "statement_descriptor", Location: "form", Type: "object"},
				{Name: "support", Location: "form", Type: "object"},
			}},
			{Name: "recipient", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "applied", Location: "form", Type: "boolean"},
				{Name: "capabilities", Location: "form", Type: "object"},
			}},
		}},
		{Name: "contact_email", Location: "form", Type: "string"},
		{Name: "contact_phone", Location: "form", Type: "string"},
		{Name: "dashboard", Location: "form", Type: "string", Enum: []string{"express", "full", "none"}},
		{Name: "defaults", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "currency", Location: "form", Type: "string"},
			{Name: "locales", Location: "form", Type: "array"},
			{Name: "profile", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "business_url", Location: "form", Type: "string"},
				{Name: "doing_business_as", Location: "form", Type: "string"},
				{Name: "product_description", Location: "form", Type: "string"},
			}},
			{Name: "responsibilities", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "fees_collector", Location: "form", Required: true, Type: "string", Enum: []string{"application", "application_custom", "application_express", "stripe"}},
				{Name: "losses_collector", Location: "form", Required: true, Type: "string", Enum: []string{"application", "stripe"}},
			}},
		}},
		{Name: "display_name", Location: "form", Type: "string"},
		{Name: "identity", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "attestations", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "directorship_declaration", Location: "form", Type: "object"},
				{Name: "ownership_declaration", Location: "form", Type: "object"},
				{Name: "persons_provided", Location: "form", Type: "object"},
				{Name: "representative_declaration", Location: "form", Type: "object"},
				{Name: "terms_of_service", Location: "form", Type: "object"},
			}},
			{Name: "business_details", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "address", Location: "form", Type: "object"},
				{Name: "annual_revenue", Location: "form", Type: "object"},
				{Name: "documents", Location: "form", Type: "object"},
				{Name: "estimated_worker_count", Location: "form", Type: "integer"},
				{Name: "id_numbers", Location: "form", Type: "array"},
				{Name: "monthly_estimated_revenue", Location: "form", Type: "object"},
				{Name: "phone", Location: "form", Type: "string"},
				{Name: "registered_name", Location: "form", Type: "string"},
				{Name: "registration_date", Location: "form", Type: "object"},
				{Name: "script_addresses", Location: "form", Type: "object"},
				{Name: "script_names", Location: "form", Type: "object"},
				{Name: "structure", Location: "form", Type: "string", Enum: []string{"cooperative", "free_zone_establishment", "free_zone_llc", "government_instrumentality", "governmental_unit", "incorporated_association", "incorporated_non_profit", "incorporated_partnership", "limited_liability_partnership", "llc", "multi_member_llc", "private_company", "private_corporation", "private_partnership", "public_company", "public_corporation", "public_listed_corporation", "public_partnership", "registered_charity", "single_member_llc", "sole_establishment", "sole_proprietorship", "tax_exempt_government_instrumentality", "trust", "unincorporated_association", "unincorporated_non_profit", "unincorporated_partnership"}},
			}},
			{Name: "country", Location: "form", Type: "string"},
			{Name: "entity_type", Location: "form", Type: "string", Enum: []string{"company", "government_entity", "individual", "non_profit"}},
			{Name: "individual", Location: "form", Type: "object", Children: []ParameterValidation{
				{Name: "additional_addresses", Location: "form", Type: "array"},
				{Name: "additional_names", Location: "form", Type: "array"},
				{Name: "address", Location: "form", Type: "object"},
				{Name: "date_of_birth", Location: "form", Type: "object"},
				{Name: "documents", Location: "form", Type: "object"},
				{Name: "email", Location: "form", Type: "string"},
				{Name: "given_name", Location: "form", Type: "string"},
				{Name: "id_numbers", Location: "form", Type: "array"},
				{Name: "legal_gender", Location: "form", Type: "string", Enum: []string{"female", "male"}},
				{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
				{Name: "nationalities", Location: "form", Type: "array"},
				{Name: "phone", Location: "form", Type: "string"},
				{Name: "political_exposure", Location: "form", Type: "string", Enum: []string{"existing", "none"}},
				{Name: "relationship", Location: "form", Type: "object"},
				{Name: "script_addresses", Location: "form", Type: "object"},
				{Name: "script_names", Location: "form", Type: "object"},
				{Name: "surname", Location: "form", Type: "string"},
			}},
		}},
		{Name: "include", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/core/accounts/{id}/close", OperationID: "PostV2CoreAccountsIdClose", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "applied_configurations", Location: "form", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v2/core/event_destinations", OperationID: "GetV2CoreEventDestinations", Params: []ParameterValidation{
		{Name: "include", Location: "query", Type: "array"},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "page", Location: "query", Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/core/event_destinations", OperationID: "PostV2CoreEventDestinations", Params: []ParameterValidation{
		{Name: "amazon_eventbridge", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "aws_account_id", Location: "form", Required: true, Type: "string"},
			{Name: "aws_region", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "azure_event_grid", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "azure_region", Location: "form", Required: true, Type: "string"},
			{Name: "azure_resource_group_name", Location: "form", Required: true, Type: "string"},
			{Name: "azure_subscription_id", Location: "form", Required: true, Type: "string"},
		}},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "enabled_events", Location: "form", Required: true, Type: "array"},
		{Name: "event_payload", Location: "form", Required: true, Type: "string", Enum: []string{"snapshot", "thin"}},
		{Name: "events_from", Location: "form", Type: "array"},
		{Name: "include", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "name", Location: "form", Required: true, Type: "string"},
		{Name: "snapshot_api_version", Location: "form", Type: "string"},
		{Name: "type", Location: "form", Required: true, Type: "string", Enum: []string{"amazon_eventbridge", "azure_event_grid", "webhook_endpoint"}},
		{Name: "webhook_endpoint", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "url", Location: "form", Required: true, Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "DELETE", Path: "/v2/core/event_destinations/{id}", OperationID: "DeleteV2CoreEventDestinationsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v2/core/event_destinations/{id}", OperationID: "GetV2CoreEventDestinationsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "include", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/core/event_destinations/{id}", OperationID: "PostV2CoreEventDestinationsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
		{Name: "description", Location: "form", Type: "string"},
		{Name: "enabled_events", Location: "form", Type: "array"},
		{Name: "include", Location: "form", Type: "array"},
		{Name: "metadata", Location: "form", Type: "object", AdditionalProperties: true},
		{Name: "name", Location: "form", Type: "string"},
		{Name: "webhook_endpoint", Location: "form", Type: "object", Children: []ParameterValidation{
			{Name: "url", Location: "form", Required: true, Type: "string"},
		}},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/core/event_destinations/{id}/disable", OperationID: "PostV2CoreEventDestinationsIdDisable", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/core/event_destinations/{id}/enable", OperationID: "PostV2CoreEventDestinationsIdEnable", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "POST", Path: "/v2/core/event_destinations/{id}/ping", OperationID: "PostV2CoreEventDestinationsIdPing", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v2/core/events", OperationID: "GetV2CoreEvents", Params: []ParameterValidation{
		{Name: "created", Location: "query", Type: "object", Children: []ParameterValidation{
			{Name: "gt", Location: "query", Type: "string"},
			{Name: "gte", Location: "query", Type: "string"},
			{Name: "lt", Location: "query", Type: "string"},
			{Name: "lte", Location: "query", Type: "string"},
		}},
		{Name: "limit", Location: "query", Type: "integer"},
		{Name: "object_id", Location: "query", Type: "string"},
		{Name: "page", Location: "query", Type: "string"},
		{Name: "types", Location: "query", Type: "array"},
	}})
	operations = append(operations, OperationValidation{Method: "GET", Path: "/v2/core/events/{id}", OperationID: "GetV2CoreEventsId", Params: []ParameterValidation{
		{Name: "id", Location: "path", Required: true, Type: "string"},
	}})
	return operations
}
