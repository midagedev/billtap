package api

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hckim/billtap/internal/billing"
)

// defaultPortalFeatures mirrors the Stripe default portal feature set. The
// hosted portal itself is a local stub, so these drive evidence and echo only.
func defaultPortalFeatures() map[string]any {
	return map[string]any{
		"customer_update": map[string]any{
			"enabled":         true,
			"allowed_updates": []string{"email", "name"},
		},
		"invoice_history": map[string]any{
			"enabled": true,
		},
		"payment_method_update": map[string]any{
			"enabled": true,
		},
		"subscription_cancel": map[string]any{
			"enabled":             true,
			"mode":                "at_period_end",
			"cancellation_reason": nil,
			"proration_behavior":  nil,
		},
		"subscription_update": map[string]any{
			"enabled":                 true,
			"proration_behavior":      "none",
			"default_allowed_updates": []string{"price", "quantity"},
		},
	}
}

func portalFeaturesFromParams(p params, features map[string]any) map[string]any {
	if features == nil {
		features = defaultPortalFeatures()
	}
	feature := func(name string) map[string]any {
		existing, _ := features[name].(map[string]any)
		if existing == nil {
			existing = map[string]any{}
			features[name] = existing
		}
		return existing
	}
	if p.has("features[customer_update][enabled]") {
		feature("customer_update")["enabled"] = p.boolDefault("features[customer_update][enabled]", true)
	}
	if allowed := p.list("features[customer_update][allowed_updates]"); len(allowed) > 0 {
		feature("customer_update")["allowed_updates"] = allowed
	}
	if p.has("features[invoice_history][enabled]") {
		feature("invoice_history")["enabled"] = p.boolDefault("features[invoice_history][enabled]", true)
	}
	if p.has("features[payment_method_update][enabled]") {
		feature("payment_method_update")["enabled"] = p.boolDefault("features[payment_method_update][enabled]", true)
	}
	cancel := feature("subscription_cancel")
	if p.has("features[subscription_cancel][enabled]") {
		cancel["enabled"] = p.boolDefault("features[subscription_cancel][enabled]", true)
	}
	if p.has("features[subscription_cancel][mode]") {
		cancel["mode"] = p.string("features[subscription_cancel][mode]")
	}
	if p.has("features[subscription_cancel][cancellation_reason]") {
		cancel["cancellation_reason"] = emptyToNil(p.string("features[subscription_cancel][cancellation_reason]"))
	}
	if p.has("features[subscription_cancel][proration_behavior]") {
		cancel["proration_behavior"] = emptyToNil(p.string("features[subscription_cancel][proration_behavior]"))
	}
	update := feature("subscription_update")
	if p.has("features[subscription_update][enabled]") {
		update["enabled"] = p.boolDefault("features[subscription_update][enabled]", true)
	}
	if p.has("features[subscription_update][proration_behavior]") {
		update["proration_behavior"] = p.string("features[subscription_update][proration_behavior]")
	}
	if allowed := p.list("features[subscription_update][default_allowed_updates]"); len(allowed) > 0 {
		update["default_allowed_updates"] = allowed
	}
	return features
}

func (h *Handler) handleBillingPortalConfigurations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateBillingPortalConfigurationCreate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		now := time.Now().UTC()
		id := "bpc_" + strconv.FormatInt(now.UnixNano(), 36)
		h.local.mu.Lock()
		isDefault := len(h.local.portalConfigurations) == 0
		configuration := map[string]any{
			"id":                 id,
			"object":             "billing_portal.configuration",
			"active":             true,
			"is_default":         isDefault,
			"business_profile":   portalBusinessProfileFromParams(p, nil),
			"default_return_url": emptyToNil(p.string("default_return_url")),
			"features":           portalFeaturesFromParams(p, nil),
			"login_page":         portalLoginPageFromParams(p, nil),
			"metadata":           nonNilMap(p.metadata()),
			"created":            now.Unix(),
			"updated":            now.Unix(),
			"livemode":           false,
		}
		err = h.local.saveLocked(kindPortalConfiguration, id, configuration)
		h.local.mu.Unlock()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, cloneEvidence(configuration))
	case http.MethodGet:
		h.local.mu.Lock()
		data := evidenceList(h.local.portalConfigurations)
		h.local.mu.Unlock()
		// Evidence maps iterate in random order; sort so starting_after pagination is stable.
		sort.Slice(data, func(i, j int) bool {
			return fmt.Sprint(data[i]["id"]) < fmt.Sprint(data[j]["id"])
		})
		data = filterPortalConfigurations(data, r)
		writeJSON(w, http.StatusOK, stripeListFromRequest(r, data))
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func (h *Handler) handleBillingPortalConfiguration(w http.ResponseWriter, r *http.Request) {
	id := strings.Trim(strings.TrimPrefix(r.URL.Path, "/v1/billing_portal/configurations/"), "/")
	if id == "" || strings.Contains(id, "/") {
		h.notFound(w, r)
		return
	}
	h.local.mu.Lock()
	configuration, ok := h.local.portalConfigurations[id]
	h.local.mu.Unlock()
	if !ok {
		writeResult(w, nil, billing.ErrNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, cloneEvidence(configuration))
	case http.MethodPost:
		p, err := parseParams(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if err := validateBillingPortalConfigurationUpdate(p); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		h.local.mu.Lock()
		current := h.local.portalConfigurations[id]
		if p.has("active") {
			current["active"] = p.boolDefault("active", true)
		}
		if p.has("business_profile[headline]") || p.has("business_profile[privacy_policy_url]") || p.has("business_profile[terms_of_service_url]") {
			current["business_profile"] = portalBusinessProfileFromParams(p, current["business_profile"])
		}
		if p.has("default_return_url") {
			current["default_return_url"] = emptyToNil(p.string("default_return_url"))
		}
		if p.has("login_page[logo_url]") {
			current["login_page"] = portalLoginPageFromParams(p, current["login_page"])
		}
		if portalFeaturesTouched(p) {
			features, _ := current["features"].(map[string]any)
			current["features"] = portalFeaturesFromParams(p, features)
		}
		if metadata := p.metadata(); metadata != nil {
			merged := map[string]string{}
			if existing, ok := current["metadata"].(map[string]string); ok {
				for key, value := range existing {
					merged[key] = value
				}
			}
			for key, value := range metadata {
				merged[key] = value
			}
			current["metadata"] = nonNilMap(merged)
		}
		current["updated"] = time.Now().UTC().Unix()
		err = h.local.saveLocked(kindPortalConfiguration, id, current)
		h.local.mu.Unlock()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, cloneEvidence(current))
	default:
		h.methodNotAllowed(w, r, "GET, POST")
	}
}

func portalBusinessProfileFromParams(p params, existing any) map[string]any {
	profile := map[string]any{"headline": nil, "privacy_policy_url": nil, "terms_of_service_url": nil}
	if current, ok := existing.(map[string]any); ok {
		for key := range profile {
			profile[key] = current[key]
		}
	}
	if p.has("business_profile[headline]") {
		profile["headline"] = emptyToNil(p.string("business_profile[headline]"))
	}
	if p.has("business_profile[privacy_policy_url]") {
		profile["privacy_policy_url"] = emptyToNil(p.string("business_profile[privacy_policy_url]"))
	}
	if p.has("business_profile[terms_of_service_url]") {
		profile["terms_of_service_url"] = emptyToNil(p.string("business_profile[terms_of_service_url]"))
	}
	return profile
}

func portalLoginPageFromParams(p params, existing any) map[string]any {
	loginPage := map[string]any{"enabled": false, "logo_url": nil}
	if current, ok := existing.(map[string]any); ok {
		for key := range loginPage {
			loginPage[key] = current[key]
		}
	}
	if p.has("login_page[logo_url]") {
		loginPage["logo_url"] = emptyToNil(p.string("login_page[logo_url]"))
	}
	return loginPage
}

func portalFeaturesTouched(p params) bool {
	for key := range p.values {
		if strings.HasPrefix(key, "features[") {
			return true
		}
	}
	return false
}

func filterPortalConfigurations(data []map[string]any, r *http.Request) []map[string]any {
	query := r.URL.Query()
	activeFilter := strings.TrimSpace(query.Get("active"))
	defaultFilter := strings.TrimSpace(query.Get("is_default"))
	if activeFilter == "" && defaultFilter == "" {
		return data
	}
	out := make([]map[string]any, 0, len(data))
	for _, configuration := range data {
		if activeFilter != "" {
			wantActive := activeFilter == "true" || activeFilter == "1"
			active, _ := configuration["active"].(bool)
			if active != wantActive {
				continue
			}
		}
		if defaultFilter != "" {
			wantDefault := defaultFilter == "true" || defaultFilter == "1"
			isDefault, _ := configuration["is_default"].(bool)
			if isDefault != wantDefault {
				continue
			}
		}
		out = append(out, configuration)
	}
	return out
}
