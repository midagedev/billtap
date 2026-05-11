package compatibility

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hckim/billtap/internal/stripecompat"
)

const InventoryVersion = "stripe-api-inventory-v2"

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
	TotalOperations           int              `json:"total_operations"`
	ImplementedOperations     int              `json:"implemented_operations"`
	InventoryOnlyOperations   int              `json:"inventory_only_operations"`
	SchemaValidatedOperations int              `json:"schema_validated_operations"`
	BilltapOnlyRoutes         int              `json:"billtap_only_routes"`
	ImplementedPercent        float64          `json:"implemented_percent"`
	SchemaValidatedPercent    float64          `json:"schema_validated_percent"`
	ByLevel                   map[string]int   `json:"by_level"`
	ByFamily                  map[string]int   `json:"by_family"`
	Families                  []FamilyCoverage `json:"families"`
}

type FamilyCoverage struct {
	Family                  string         `json:"family"`
	Priority                string         `json:"priority"`
	TargetLevel             string         `json:"target_level"`
	TotalOperations         int            `json:"total_operations"`
	ImplementedOperations   int            `json:"implemented_operations"`
	InventoryOnlyOperations int            `json:"inventory_only_operations"`
	ImplementedPercent      float64        `json:"implemented_percent"`
	ByLevel                 map[string]int `json:"by_level"`
	NextMilestone           string         `json:"next_milestone"`
}

type StripeOperationCoverage struct {
	Family          string   `json:"family"`
	Resource        string   `json:"resource"`
	Path            string   `json:"path"`
	NormalizedPath  string   `json:"normalized_path"`
	Method          string   `json:"method"`
	OperationID     string   `json:"operation_id,omitempty"`
	Implemented     bool     `json:"implemented"`
	SchemaValidated bool     `json:"schema_validated"`
	BilltapLevel    string   `json:"billtap_level"`
	TargetLevel     string   `json:"target_level"`
	Stateful        bool     `json:"stateful"`
	WebhookEvents   []string `json:"webhook_events,omitempty"`
	ScorecardCases  []string `json:"scorecard_cases,omitempty"`
	SDKSmoke        []string `json:"sdk_smoke,omitempty"`
	Docs            string   `json:"docs,omitempty"`
	Risks           []string `json:"risks,omitempty"`
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
	OperationID string              `json:"operationId"`
	Tags        []string            `json:"tags"`
	Deprecated  bool                `json:"deprecated"`
	ResourceID  string              `json:"x-resourceId"`
	Parameters  []openAPIParameter  `json:"parameters"`
	RequestBody *openAPIRequestBody `json:"requestBody"`
}

type openAPIParameter struct {
	Name   string          `json:"name"`
	In     string          `json:"in"`
	Schema json.RawMessage `json:"schema"`
}

type openAPIRequestBody struct {
	Content map[string]openAPIRequestMedia `json:"content"`
}

type openAPIRequestMedia struct {
	Schema json.RawMessage `json:"schema"`
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
	validationCatalog := stripecompat.DefaultValidationCatalog()
	familyStats := map[string]*FamilyCoverage{}
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
			_, hasBundledValidation := validationCatalog.Lookup(method, normalizedPath)
			schemaValidated := hasBundledValidation && operationHasValidationSurface(operation)
			family := inferFamily(path)
			resource := inferResource(path, operation, spec.Components.Schemas)
			level := "L0"
			if implemented {
				level = current.Level
			}
			item := StripeOperationCoverage{
				Family:          family,
				Resource:        resource,
				Path:            path,
				NormalizedPath:  normalizedPath,
				Method:          strings.ToUpper(method),
				OperationID:     operation.OperationID,
				Implemented:     implemented,
				SchemaValidated: schemaValidated,
				BilltapLevel:    level,
				TargetLevel:     targetLevelForFamily(family),
				Stateful:        current.Stateful,
				WebhookEvents:   append([]string(nil), current.WebhookEvents...),
				ScorecardCases:  append([]string(nil), current.ScorecardCases...),
				SDKSmoke:        append([]string(nil), current.SDKSmoke...),
				Docs:            current.Docs,
				Risks:           append([]string(nil), current.Risks...),
			}
			if !implemented {
				item.Risks = append(item.Risks, "inventory-only; no Billtap runtime claim")
			}
			inventory.Operations = append(inventory.Operations, item)
			inventory.Summary.TotalOperations++
			if schemaValidated {
				inventory.Summary.SchemaValidatedOperations++
			}
			inventory.Summary.ByLevel[level]++
			inventory.Summary.ByFamily[family]++
			stats := familyStats[family]
			if stats == nil {
				stats = &FamilyCoverage{
					Family:        family,
					Priority:      priorityForFamily(family),
					TargetLevel:   targetLevelForFamily(family),
					ByLevel:       map[string]int{},
					NextMilestone: nextMilestoneForFamily(family),
				}
				familyStats[family] = stats
			}
			stats.TotalOperations++
			stats.ByLevel[level]++
			if implemented {
				inventory.Summary.ImplementedOperations++
				stats.ImplementedOperations++
			} else {
				inventory.Summary.InventoryOnlyOperations++
				stats.InventoryOnlyOperations++
			}
		}
	}
	inventory.Summary.BilltapOnlyRoutes = len(inventory.BilltapOnlyRoutes)
	inventory.Summary.ImplementedPercent = percentage(inventory.Summary.ImplementedOperations, inventory.Summary.TotalOperations)
	inventory.Summary.SchemaValidatedPercent = percentage(inventory.Summary.SchemaValidatedOperations, inventory.Summary.TotalOperations)
	inventory.Summary.Families = sortedFamilyCoverage(familyStats)
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
	fmt.Fprintf(&b, "- Operations: `%d` total, `%d` implemented, `%d` inventory-only, `%.1f%%` implemented\n",
		i.Summary.TotalOperations, i.Summary.ImplementedOperations, i.Summary.InventoryOnlyOperations, i.Summary.ImplementedPercent)
	fmt.Fprintf(&b, "- OpenAPI validation catalog: `%d` route schemas, `%.1f%%` schema-visible\n",
		i.Summary.SchemaValidatedOperations, i.Summary.SchemaValidatedPercent)

	b.WriteString("\n## Levels\n\n")
	for _, level := range sortedKeys(i.Summary.ByLevel) {
		fmt.Fprintf(&b, "- `%s`: `%d`\n", level, i.Summary.ByLevel[level])
	}

	b.WriteString("\n## Family Coverage\n\n")
	b.WriteString("| Priority | Family | Target | Total | Implemented | Inventory-only | Coverage | Next milestone |\n")
	b.WriteString("| --- | --- | --- | ---: | ---: | ---: | ---: | --- |\n")
	for _, family := range i.Summary.Families {
		fmt.Fprintf(&b, "| `%s` | `%s` | `%s` | %d | %d | %d | %.1f%% | %s |\n",
			escapeTable(family.Priority),
			escapeTable(family.Family),
			escapeTable(family.TargetLevel),
			family.TotalOperations,
			family.ImplementedOperations,
			family.InventoryOnlyOperations,
			family.ImplementedPercent,
			escapeTable(family.NextMilestone),
		)
	}

	b.WriteString("\n## Operations\n\n")
	b.WriteString("| Method | Path | Family | Resource | Level | Target | Implemented | Schema validated | Risks |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- | --- | --- | --- |\n")
	for _, op := range i.Operations {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | `%s` | `%s` | `%t` | `%t` | %s |\n",
			escapeTable(op.Method),
			escapeTable(op.Path),
			escapeTable(op.Family),
			escapeTable(op.Resource),
			op.BilltapLevel,
			op.TargetLevel,
			op.Implemented,
			op.SchemaValidated,
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
	out := map[string]routeCoverage{}
	for _, claim := range stripecompat.DefaultRegistry().Claims() {
		out[routeKey(claim.Method, claim.Path)] = routeCoverage{
			Level:          claim.Level,
			Stateful:       claim.Stateful,
			Docs:           claim.Docs,
			WebhookEvents:  append([]string(nil), claim.WebhookEvents...),
			ScorecardCases: append([]string(nil), claim.ScorecardCases...),
			SDKSmoke:       append([]string(nil), claim.SDKSmoke...),
			Risks:          append([]string(nil), claim.Risks...),
		}
	}
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
	return stripecompat.RouteKey(method, normalizedPath)
}

func normalizeOpenAPIPath(path string) string {
	return stripecompat.NormalizePath(path)
}

func isHTTPMethod(method string) bool {
	switch strings.ToLower(method) {
	case "get", "post", "delete", "put", "patch":
		return true
	default:
		return false
	}
}

func operationHasValidationSurface(operation openAPIOperation) bool {
	for _, param := range operation.Parameters {
		if strings.TrimSpace(param.Name) != "" || len(param.Schema) > 0 {
			return true
		}
	}
	if operation.RequestBody == nil {
		return false
	}
	for _, media := range operation.RequestBody.Content {
		if len(media.Schema) > 0 {
			return true
		}
	}
	return false
}

func inferFamily(path string) string {
	switch {
	case isConnectPath(path):
		return "connect"
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
	case strings.Contains(path, "/products") || strings.Contains(path, "/prices") || strings.Contains(path, "/coupons") || strings.Contains(path, "/promotion_codes") || strings.Contains(path, "/tax"):
		return "catalog"
	case strings.Contains(path, "/customers"):
		return "customers"
	default:
		return "auxiliary"
	}
}

func isConnectPath(path string) bool {
	for _, prefix := range []string{
		"/v1/account",
		"/v1/accounts",
		"/v1/account_links",
		"/v1/account_sessions",
		"/v1/application_fees",
		"/v1/transfers",
		"/v1/payouts",
		"/v2/core/accounts",
	} {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	return false
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

func priorityForFamily(family string) string {
	switch family {
	case "checkout", "billing_portal", "webhooks", "billing":
		return "P0"
	case "customers", "catalog", "payments", "payment_history", "connect":
		return "P1"
	default:
		return "P3"
	}
}

func nextMilestoneForFamily(family string) string {
	switch family {
	case "checkout":
		return "Close subscription checkout gaps and add SDK/adoption smoke for optional params."
	case "billing_portal":
		return "Add portal configuration fixtures and subscription/payment-method update scenarios."
	case "webhooks":
		return "Expand connected-account routing, thin event fixtures, and replay evidence."
	case "billing":
		return "Add renewal, trial, dunning, subscription schedule, coupon, and credit-note scenarios."
	case "customers":
		return "Add OpenAPI-backed validation, search/list parity, and payment source fixtures."
	case "catalog":
		return "Add coupon, promotion code, tax-rate, and product/price search validation."
	case "payments":
		return "Add PaymentIntent and SetupIntent create/confirm/capture/cancel state machines."
	case "payment_history":
		return "Add charge, refund, balance transaction, dispute, and payment history evidence."
	case "connect":
		return "Add account self/delete and people/person fixtures to close the remaining Connect inventory routes."
	default:
		return "Keep inventory visible and add schema/fixture smoke only when adoption requires it."
	}
}

func sortedFamilyCoverage(stats map[string]*FamilyCoverage) []FamilyCoverage {
	keys := make([]string, 0, len(stats))
	for key := range stats {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		left := stats[keys[i]]
		right := stats[keys[j]]
		if priorityRank(left.Priority) == priorityRank(right.Priority) {
			return left.Family < right.Family
		}
		return priorityRank(left.Priority) < priorityRank(right.Priority)
	})

	families := make([]FamilyCoverage, 0, len(keys))
	for _, key := range keys {
		family := *stats[key]
		family.ImplementedPercent = percentage(family.ImplementedOperations, family.TotalOperations)
		families = append(families, family)
	}
	return families
}

func priorityRank(priority string) int {
	switch priority {
	case "P0":
		return 0
	case "P1":
		return 1
	case "P2":
		return 2
	case "P3":
		return 3
	default:
		return 99
	}
}

func percentage(part int, total int) float64 {
	if total == 0 {
		return 0
	}
	return math.Round((float64(part)/float64(total))*1000) / 10
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
