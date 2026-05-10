package compatibility

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const InventoryVersion = "stripe-api-inventory-v1"

type InventoryOptions struct {
	OpenAPIPath string
	OutputDir   string
	Source      string
	Now         func() time.Time
}

type StripeAPIInventory struct {
	Name              string                    `json:"name"`
	InventoryVersion  string                    `json:"inventory_version"`
	GeneratedAt       time.Time                 `json:"generated_at"`
	Source            string                    `json:"source"`
	OpenAPIVersion    string                    `json:"openapi_version,omitempty"`
	StripeAPIVersion  string                    `json:"stripe_api_version,omitempty"`
	Summary           InventorySummary          `json:"summary"`
	Operations        []StripeOperationCoverage `json:"operations"`
	BilltapOnlyRoutes []BilltapRouteCoverage    `json:"billtap_only_routes"`
}

type InventorySummary struct {
	TotalOperations         int            `json:"total_operations"`
	ImplementedOperations   int            `json:"implemented_operations"`
	InventoryOnlyOperations int            `json:"inventory_only_operations"`
	BilltapOnlyRoutes       int            `json:"billtap_only_routes"`
	ByLevel                 map[string]int `json:"by_level"`
	ByFamily                map[string]int `json:"by_family"`
}

type StripeOperationCoverage struct {
	Family         string   `json:"family"`
	Resource       string   `json:"resource"`
	Path           string   `json:"path"`
	NormalizedPath string   `json:"normalized_path"`
	Method         string   `json:"method"`
	OperationID    string   `json:"operation_id,omitempty"`
	Implemented    bool     `json:"implemented"`
	BilltapLevel   string   `json:"billtap_level"`
	TargetLevel    string   `json:"target_level"`
	Stateful       bool     `json:"stateful"`
	WebhookEvents  []string `json:"webhook_events,omitempty"`
	ScorecardCases []string `json:"scorecard_cases,omitempty"`
	SDKSmoke       []string `json:"sdk_smoke,omitempty"`
	Docs           string   `json:"docs,omitempty"`
	Risks          []string `json:"risks,omitempty"`
}

type BilltapRouteCoverage struct {
	Path         string   `json:"path"`
	Method       string   `json:"method"`
	BilltapLevel string   `json:"billtap_level"`
	Reason       string   `json:"reason"`
	Docs         string   `json:"docs"`
	Risks        []string `json:"risks,omitempty"`
}

type openAPISpec struct {
	OpenAPI string `json:"openapi"`
	Info    struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	} `json:"info"`
	Paths      map[string]map[string]openAPIOperation `json:"paths"`
	Components struct {
		Schemas map[string]openAPISchema `json:"schemas"`
	} `json:"components"`
}

type openAPIOperation struct {
	OperationID string   `json:"operationId"`
	Tags        []string `json:"tags"`
	Deprecated  bool     `json:"deprecated"`
	ResourceID  string   `json:"x-resourceId"`
}

type openAPISchema struct {
	ResourceID       string   `json:"x-resourceId"`
	ExpandableFields []string `json:"x-expandableFields"`
}

type routeCoverage struct {
	Level          string
	Stateful       bool
	Docs           string
	WebhookEvents  []string
	ScorecardCases []string
	SDKSmoke       []string
	Risks          []string
}

func GenerateInventory(ctx context.Context, opts InventoryOptions) (StripeAPIInventory, error) {
	_ = ctx
	if strings.TrimSpace(opts.OpenAPIPath) == "" {
		return StripeAPIInventory{}, fmt.Errorf("OpenAPI path is required")
	}
	now := opts.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	raw, err := os.ReadFile(opts.OpenAPIPath)
	if err != nil {
		return StripeAPIInventory{}, err
	}
	var spec openAPISpec
	if err := json.Unmarshal(raw, &spec); err != nil {
		return StripeAPIInventory{}, fmt.Errorf("decode OpenAPI spec: %w", err)
	}
	if len(spec.Paths) == 0 {
		return StripeAPIInventory{}, fmt.Errorf("OpenAPI spec has no paths")
	}

	source := opts.Source
	if source == "" {
		source = opts.OpenAPIPath
	}
	inventory := StripeAPIInventory{
		Name:             "Billtap Stripe API compatibility inventory",
		InventoryVersion: InventoryVersion,
		GeneratedAt:      now(),
		Source:           source,
		OpenAPIVersion:   spec.OpenAPI,
		StripeAPIVersion: spec.Info.Version,
		Summary: InventorySummary{
			ByLevel:  map[string]int{},
			ByFamily: map[string]int{},
		},
		BilltapOnlyRoutes: billtapOnlyRoutes(),
	}

	coverage := currentStripeRouteCoverage()
	paths := sortedKeys(spec.Paths)
	for _, path := range paths {
		methods := spec.Paths[path]
		for _, method := range sortedKeys(methods) {
			if !isHTTPMethod(method) {
				continue
			}
			operation := methods[method]
			normalizedPath := normalizeOpenAPIPath(path)
			key := routeKey(method, normalizedPath)
			current, implemented := coverage[key]
			family := inferFamily(path)
			resource := inferResource(path, operation, spec.Components.Schemas)
			level := "L0"
			if implemented {
				level = current.Level
			}
			item := StripeOperationCoverage{
				Family:         family,
				Resource:       resource,
				Path:           path,
				NormalizedPath: normalizedPath,
				Method:         strings.ToUpper(method),
				OperationID:    operation.OperationID,
				Implemented:    implemented,
				BilltapLevel:   level,
				TargetLevel:    targetLevelForFamily(family),
				Stateful:       current.Stateful,
				WebhookEvents:  append([]string(nil), current.WebhookEvents...),
				ScorecardCases: append([]string(nil), current.ScorecardCases...),
				SDKSmoke:       append([]string(nil), current.SDKSmoke...),
				Docs:           current.Docs,
				Risks:          append([]string(nil), current.Risks...),
			}
			if !implemented {
				item.Risks = append(item.Risks, "inventory-only; no Billtap runtime claim")
			}
			inventory.Operations = append(inventory.Operations, item)
			inventory.Summary.TotalOperations++
			inventory.Summary.ByLevel[level]++
			inventory.Summary.ByFamily[family]++
			if implemented {
				inventory.Summary.ImplementedOperations++
			} else {
				inventory.Summary.InventoryOnlyOperations++
			}
		}
	}
	inventory.Summary.BilltapOnlyRoutes = len(inventory.BilltapOnlyRoutes)
	sort.Slice(inventory.Operations, func(i, j int) bool {
		left := inventory.Operations[i]
		right := inventory.Operations[j]
		if left.Path == right.Path {
			return left.Method < right.Method
		}
		return left.Path < right.Path
	})
	return inventory, nil
}

func WriteInventoryArtifacts(ctx context.Context, opts InventoryOptions) (StripeAPIInventory, error) {
	outputDir := opts.OutputDir
	if outputDir == "" {
		outputDir = DefaultOutputDir
	}
	inventory, err := GenerateInventory(ctx, opts)
	if err != nil {
		return inventory, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return inventory, err
	}
	body, err := inventory.JSON()
	if err != nil {
		return inventory, err
	}
	if err := os.WriteFile(filepath.Join(outputDir, "stripe-api-inventory.json"), body, 0o644); err != nil {
		return inventory, err
	}
	if err := os.WriteFile(filepath.Join(outputDir, "stripe-api-inventory.md"), []byte(inventory.Markdown()), 0o644); err != nil {
		return inventory, err
	}
	return inventory, nil
}

func (i StripeAPIInventory) JSON() ([]byte, error) {
	return json.MarshalIndent(i, "", "  ")
}

func (i StripeAPIInventory) Markdown() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "# Stripe API Compatibility Inventory\n\n")
	fmt.Fprintf(&b, "- Version: `%s`\n", i.InventoryVersion)
	fmt.Fprintf(&b, "- Generated: `%s`\n", i.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "- Source: `%s`\n", i.Source)
	if i.StripeAPIVersion != "" {
		fmt.Fprintf(&b, "- Stripe API version: `%s`\n", i.StripeAPIVersion)
	}
	fmt.Fprintf(&b, "- Operations: `%d` total, `%d` implemented, `%d` inventory-only\n",
		i.Summary.TotalOperations, i.Summary.ImplementedOperations, i.Summary.InventoryOnlyOperations)

	b.WriteString("\n## Levels\n\n")
	for _, level := range sortedKeys(i.Summary.ByLevel) {
		fmt.Fprintf(&b, "- `%s`: `%d`\n", level, i.Summary.ByLevel[level])
	}

	b.WriteString("\n## Families\n\n")
	for _, family := range sortedKeys(i.Summary.ByFamily) {
		fmt.Fprintf(&b, "- `%s`: `%d`\n", family, i.Summary.ByFamily[family])
	}

	b.WriteString("\n## Operations\n\n")
	b.WriteString("| Method | Path | Family | Resource | Level | Target | Implemented | Risks |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- | --- | --- |\n")
	for _, op := range i.Operations {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | `%s` | `%s` | `%t` | %s |\n",
			escapeTable(op.Method),
			escapeTable(op.Path),
			escapeTable(op.Family),
			escapeTable(op.Resource),
			op.BilltapLevel,
			op.TargetLevel,
			op.Implemented,
			escapeTable(strings.Join(op.Risks, "; ")),
		)
	}

	if len(i.BilltapOnlyRoutes) > 0 {
		b.WriteString("\n## Billtap-Specific `/v1` Exceptions\n\n")
		b.WriteString("| Method | Path | Level | Reason |\n")
		b.WriteString("| --- | --- | --- | --- |\n")
		for _, route := range i.BilltapOnlyRoutes {
			fmt.Fprintf(&b, "| %s | %s | `%s` | %s |\n",
				escapeTable(route.Method),
				escapeTable(route.Path),
				route.BilltapLevel,
				escapeTable(route.Reason),
			)
		}
	}
	return b.String()
}

func currentStripeRouteCoverage() map[string]routeCoverage {
	docs := "docs/COMPATIBILITY.md#supported-stripe-like-api-subset"
	out := map[string]routeCoverage{}
	add := func(method string, path string, coverage routeCoverage) {
		if coverage.Docs == "" {
			coverage.Docs = docs
		}
		out[routeKey(method, path)] = coverage
	}

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		add(method, "/v1/customers", routeCoverage{Level: "L3", Stateful: true, SDKSmoke: []string{"stripe-node"}})
		add(method, "/v1/customers/{id}", routeCoverage{Level: "L3", Stateful: true, SDKSmoke: []string{"stripe-node"}})
		add(method, "/v1/products", routeCoverage{Level: "L3", Stateful: true, ScorecardCases: []string{"products.create.success"}, SDKSmoke: []string{"stripe-node"}})
		add(method, "/v1/products/{id}", routeCoverage{Level: "L3", Stateful: true, SDKSmoke: []string{"stripe-node"}})
		add(method, "/v1/prices", routeCoverage{Level: "L3", Stateful: true, ScorecardCases: []string{"prices.create.invalid_json_amount_type"}, SDKSmoke: []string{"stripe-node"}})
		add(method, "/v1/prices/{id}", routeCoverage{Level: "L3", Stateful: true, SDKSmoke: []string{"stripe-node"}})
		add(method, "/v1/subscriptions", routeCoverage{Level: "L3", Stateful: true, WebhookEvents: []string{"customer.subscription.created", "customer.subscription.updated", "customer.subscription.deleted"}})
		add(method, "/v1/subscriptions/{id}", routeCoverage{Level: "L3", Stateful: true, WebhookEvents: []string{"customer.subscription.updated", "customer.subscription.deleted"}})
	}
	add(http.MethodGet, "/v1/products/search", routeCoverage{Level: "L2", Stateful: false, Risks: []string{"metadata equality filters only; no Stripe Search Query Language parity"}})

	add(http.MethodPost, "/v1/checkout/sessions", routeCoverage{Level: "L4", Stateful: true, ScorecardCases: []string{"checkout.sessions.create.java_sdk_optional_params"}, SDKSmoke: []string{"stripe-node"}, Risks: []string{"subscription mode only"}})
	add(http.MethodGet, "/v1/checkout/sessions", routeCoverage{Level: "L4", Stateful: true, SDKSmoke: []string{"stripe-node"}, Risks: []string{"subscription mode only"}})
	add(http.MethodGet, "/v1/checkout/sessions/{id}", routeCoverage{Level: "L4", Stateful: true, SDKSmoke: []string{"stripe-node"}, Risks: []string{"subscription mode only"}})
	add(http.MethodPost, "/v1/billing_portal/sessions", routeCoverage{Level: "L2", Stateful: false, Risks: []string{"portal configuration and full Stripe-hosted portal behavior are not modeled"}})

	add(http.MethodPost, "/v1/subscription_items", routeCoverage{Level: "L3", Stateful: true, ScorecardCases: []string{"subscription_items.create.invalid_quantity"}})
	add(http.MethodDelete, "/v1/subscription_items/{id}", routeCoverage{Level: "L3", Stateful: true})
	add(http.MethodDelete, "/v1/subscriptions/{id}", routeCoverage{Level: "L3", Stateful: true, WebhookEvents: []string{"customer.subscription.deleted"}})

	add(http.MethodGet, "/v1/invoices", routeCoverage{Level: "L3", Stateful: true, SDKSmoke: []string{"stripe-node"}})
	add(http.MethodGet, "/v1/invoices/{id}", routeCoverage{Level: "L3", Stateful: true, SDKSmoke: []string{"stripe-node"}})
	add(http.MethodPost, "/v1/invoices/create_preview", routeCoverage{Level: "L2", Stateful: false, Risks: []string{"zero-value local smoke-test invoice; no full proration model"}})

	add(http.MethodGet, "/v1/payment_intents", routeCoverage{Level: "L3", Stateful: true, SDKSmoke: []string{"stripe-node"}})
	add(http.MethodGet, "/v1/payment_intents/{id}", routeCoverage{Level: "L3", Stateful: true, SDKSmoke: []string{"stripe-node"}})
	add(http.MethodGet, "/v1/payment_methods", routeCoverage{Level: "L2", Stateful: false, Risks: []string{"deterministic sandbox card projection only"}})

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		add(method, "/v1/webhook_endpoints", routeCoverage{Level: "L5", Stateful: true, SDKSmoke: []string{"stripe-node"}})
		add(method, "/v1/webhook_endpoints/{id}", routeCoverage{Level: "L5", Stateful: true, SDKSmoke: []string{"stripe-node"}})
	}
	add(http.MethodDelete, "/v1/webhook_endpoints/{id}", routeCoverage{Level: "L5", Stateful: true, SDKSmoke: []string{"stripe-node"}})
	add(http.MethodGet, "/v1/events", routeCoverage{Level: "L5", Stateful: true, SDKSmoke: []string{"stripe-node"}})
	add(http.MethodGet, "/v1/events/{id}", routeCoverage{Level: "L5", Stateful: true, SDKSmoke: []string{"stripe-node"}})
	return out
}

func billtapOnlyRoutes() []BilltapRouteCoverage {
	return []BilltapRouteCoverage{
		{
			Method:       http.MethodPost,
			Path:         "/v1/checkout/sessions/{id}/complete",
			BilltapLevel: "L4",
			Reason:       "sandbox checkout completion endpoint used by local tests and hosted Billtap checkout UI",
			Docs:         "docs/COMPATIBILITY.md#billtap-apis",
			Risks:        []string{"not a Stripe API endpoint"},
		},
	}
}

func routeKey(method string, normalizedPath string) string {
	return strings.ToUpper(method) + " " + normalizedPath
}

var pathParamPattern = regexp.MustCompile(`\{[^}/]+\}`)

func normalizeOpenAPIPath(path string) string {
	return pathParamPattern.ReplaceAllString(path, "{id}")
}

func isHTTPMethod(method string) bool {
	switch strings.ToLower(method) {
	case "get", "post", "delete", "put", "patch":
		return true
	default:
		return false
	}
}

func inferFamily(path string) string {
	switch {
	case strings.Contains(path, "/checkout/"):
		return "checkout"
	case strings.Contains(path, "/billing_portal/"):
		return "billing_portal"
	case strings.Contains(path, "/webhook_endpoints") || strings.Contains(path, "/events"):
		return "webhooks"
	case strings.Contains(path, "/subscriptions") || strings.Contains(path, "/subscription_items") || strings.Contains(path, "/invoices"):
		return "billing"
	case strings.Contains(path, "/payment_intents") || strings.Contains(path, "/setup_intents") || strings.Contains(path, "/payment_methods") || strings.Contains(path, "/charges"):
		return "payments"
	case strings.Contains(path, "/refunds") || strings.Contains(path, "/disputes") || strings.Contains(path, "/balance_transactions") || strings.Contains(path, "/credit_notes"):
		return "payment_history"
	case strings.Contains(path, "/accounts") || strings.Contains(path, "/account_links") || strings.Contains(path, "/transfers") || strings.Contains(path, "/payouts") || strings.Contains(path, "/application_fees"):
		return "connect"
	case strings.Contains(path, "/products") || strings.Contains(path, "/prices") || strings.Contains(path, "/coupons") || strings.Contains(path, "/promotion_codes") || strings.Contains(path, "/tax"):
		return "catalog"
	case strings.Contains(path, "/customers"):
		return "customers"
	default:
		return "auxiliary"
	}
}

func inferResource(path string, operation openAPIOperation, schemas map[string]openAPISchema) string {
	if operation.ResourceID != "" {
		return operation.ResourceID
	}
	for _, segment := range strings.Split(strings.Trim(path, "/"), "/") {
		if segment == "" || segment == "v1" || segment == "v2" || strings.HasPrefix(segment, "{") {
			continue
		}
		return singularize(segment)
	}
	for _, schema := range schemas {
		if schema.ResourceID != "" {
			return schema.ResourceID
		}
	}
	return "unknown"
}

func singularize(value string) string {
	value = strings.ReplaceAll(value, "-", "_")
	if strings.HasSuffix(value, "ies") {
		return strings.TrimSuffix(value, "ies") + "y"
	}
	if strings.HasSuffix(value, "s") {
		return strings.TrimSuffix(value, "s")
	}
	return value
}

func targetLevelForFamily(family string) string {
	switch family {
	case "checkout", "webhooks", "billing":
		return "L4-L6"
	case "payments", "payment_history", "customers", "catalog":
		return "L3-L6"
	case "billing_portal":
		return "L2-L5"
	case "connect":
		return "L2-L5"
	default:
		return "L1-L2"
	}
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
