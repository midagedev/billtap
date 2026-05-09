package fixtures

import (
	"fmt"
	"strings"

	"github.com/hckim/billtap/internal/billing"
)

func evaluateExpectation(snapshot Snapshot, expectation Expectation) AssertionResult {
	target := normalizeTarget(expectation.Target)
	result := AssertionResult{
		Target:   target,
		Expected: expectationMap(expectation),
	}

	switch target {
	case "customer":
		result.Matched = countCustomers(snapshot.Customers, expectation)
	case "product":
		result.Matched = countProducts(snapshot.Products, expectation)
	case "price":
		result.Matched = countPrices(snapshot.Prices, expectation)
	case "checkout_session":
		result.Matched = countCheckoutSessions(snapshot.CheckoutSessions, expectation)
	case "subscription":
		result.Matched = countSubscriptions(snapshot.Subscriptions, expectation)
	case "invoice":
		result.Matched = countInvoices(snapshot.Invoices, expectation)
	case "payment_intent":
		result.Matched = countPaymentIntents(snapshot.PaymentIntents, expectation)
	case "timeline":
		result.Matched = countTimeline(snapshot.Timeline, expectation)
	default:
		result.Pass = false
		result.Message = fmt.Sprintf("unsupported target %q", expectation.Target)
		return result
	}

	result.Pass = expectationPasses(result.Matched, expectation)
	result.Message = expectationMessage(result.Matched, expectation, result.Pass)
	return result
}

func countCustomers(items []billing.Customer, expectation Expectation) int {
	count := 0
	for _, item := range items {
		if expectation.ID != "" && item.ID != expectation.ID {
			continue
		}
		if expectation.Customer != "" && item.ID != expectation.Customer {
			continue
		}
		if expectation.Email != "" && item.Email != expectation.Email {
			continue
		}
		if !metadataContains(item.Metadata, expectation.Metadata) {
			continue
		}
		count++
	}
	return count
}

func countProducts(items []billing.Product, expectation Expectation) int {
	count := 0
	for _, item := range items {
		id := firstNonEmpty(expectation.ID, expectation.Product)
		if id != "" && item.ID != id {
			continue
		}
		if !metadataContains(item.Metadata, expectation.Metadata) {
			continue
		}
		count++
	}
	return count
}

func countPrices(items []billing.Price, expectation Expectation) int {
	count := 0
	for _, item := range items {
		id := firstNonEmpty(expectation.ID, expectation.Price)
		if id != "" && item.ID != id {
			continue
		}
		if expectation.Product != "" && item.ProductID != expectation.Product {
			continue
		}
		if expectation.LookupKey != "" && item.LookupKey != expectation.LookupKey && item.Metadata["lookup_key"] != expectation.LookupKey {
			continue
		}
		if !metadataContains(item.Metadata, expectation.Metadata) {
			continue
		}
		count++
	}
	return count
}

func countCheckoutSessions(items []billing.CheckoutSession, expectation Expectation) int {
	count := 0
	for _, item := range items {
		if expectation.ID != "" && item.ID != expectation.ID {
			continue
		}
		if expectation.Customer != "" && item.CustomerID != expectation.Customer {
			continue
		}
		if expectation.Status != "" && item.Status != expectation.Status {
			continue
		}
		count++
	}
	return count
}

func countSubscriptions(items []billing.Subscription, expectation Expectation) int {
	count := 0
	for _, item := range items {
		if expectation.ID != "" && item.ID != expectation.ID {
			continue
		}
		if expectation.Customer != "" && item.CustomerID != expectation.Customer {
			continue
		}
		if expectation.Status != "" && item.Status != expectation.Status {
			continue
		}
		if expectation.Price != "" && !subscriptionHasPrice(item, expectation.Price) {
			continue
		}
		if expectation.Quantity != nil && !subscriptionHasQuantity(item, *expectation.Quantity) {
			continue
		}
		if !metadataContains(item.Metadata, expectation.Metadata) {
			continue
		}
		count++
	}
	return count
}

func countInvoices(items []billing.Invoice, expectation Expectation) int {
	count := 0
	for _, item := range items {
		if expectation.ID != "" && item.ID != expectation.ID {
			continue
		}
		if expectation.Customer != "" && item.CustomerID != expectation.Customer {
			continue
		}
		if expectation.Status != "" && item.Status != expectation.Status {
			continue
		}
		count++
	}
	return count
}

func countPaymentIntents(items []billing.PaymentIntent, expectation Expectation) int {
	count := 0
	for _, item := range items {
		if expectation.ID != "" && item.ID != expectation.ID {
			continue
		}
		if expectation.Customer != "" && item.CustomerID != expectation.Customer {
			continue
		}
		if expectation.Status != "" && item.Status != expectation.Status {
			continue
		}
		count++
	}
	return count
}

func countTimeline(items []billing.TimelineEntry, expectation Expectation) int {
	count := 0
	for _, item := range items {
		if expectation.ID != "" && item.ID != expectation.ID {
			continue
		}
		if expectation.Customer != "" && item.CustomerID != expectation.Customer {
			continue
		}
		if expectation.Status != "" && item.Action != expectation.Status {
			continue
		}
		if !metadataContains(item.Data, expectation.Metadata) {
			continue
		}
		count++
	}
	return count
}

func expectationPasses(matched int, expectation Expectation) bool {
	if expectation.Count != nil {
		return matched == *expectation.Count
	}
	if expectation.CountAtLeast != nil {
		return matched >= *expectation.CountAtLeast
	}
	exists := true
	if expectation.Exists != nil {
		exists = *expectation.Exists
	}
	if exists {
		return matched > 0
	}
	return matched == 0
}

func expectationMessage(matched int, expectation Expectation, pass bool) string {
	if pass {
		return fmt.Sprintf("matched %d object(s)", matched)
	}
	if expectation.Count != nil {
		return fmt.Sprintf("matched %d object(s), expected exactly %d", matched, *expectation.Count)
	}
	if expectation.CountAtLeast != nil {
		return fmt.Sprintf("matched %d object(s), expected at least %d", matched, *expectation.CountAtLeast)
	}
	exists := true
	if expectation.Exists != nil {
		exists = *expectation.Exists
	}
	if exists {
		return fmt.Sprintf("matched %d object(s), expected at least one", matched)
	}
	return fmt.Sprintf("matched %d object(s), expected none", matched)
}

func expectationMap(expectation Expectation) map[string]any {
	out := map[string]any{"target": expectation.Target}
	if expectation.ID != "" {
		out["id"] = expectation.ID
	}
	if expectation.Customer != "" {
		out["customer"] = expectation.Customer
	}
	if expectation.Email != "" {
		out["email"] = expectation.Email
	}
	if expectation.Product != "" {
		out["product"] = expectation.Product
	}
	if expectation.Price != "" {
		out["price"] = expectation.Price
	}
	if expectation.LookupKey != "" {
		out["lookupKey"] = expectation.LookupKey
	}
	if expectation.Status != "" {
		out["status"] = expectation.Status
	}
	if len(expectation.Metadata) > 0 {
		out["metadata"] = expectation.Metadata
	}
	if expectation.Exists != nil {
		out["exists"] = *expectation.Exists
	}
	if expectation.Count != nil {
		out["count"] = *expectation.Count
	}
	if expectation.CountAtLeast != nil {
		out["countAtLeast"] = *expectation.CountAtLeast
	}
	if expectation.Quantity != nil {
		out["quantity"] = *expectation.Quantity
	}
	return out
}

func normalizeTarget(target string) string {
	target = strings.ToLower(strings.TrimSpace(target))
	target = strings.ReplaceAll(target, "-", "_")
	target = strings.ReplaceAll(target, ".", "_")
	switch target {
	case "customers":
		return "customer"
	case "products":
		return "product"
	case "prices":
		return "price"
	case "checkout", "checkout_sessions", "checkout_session":
		return "checkout_session"
	case "subscriptions":
		return "subscription"
	case "invoices":
		return "invoice"
	case "payment_intents", "paymentintent":
		return "payment_intent"
	case "timeline_entries", "timeline_entry":
		return "timeline"
	default:
		return target
	}
}

func metadataContains(actual map[string]string, expected map[string]string) bool {
	for key, value := range expected {
		if actual[key] != value {
			return false
		}
	}
	return true
}

func subscriptionHasPrice(item billing.Subscription, priceID string) bool {
	for _, lineItem := range item.Items {
		if lineItem.PriceID == priceID {
			return true
		}
	}
	return false
}

func subscriptionHasQuantity(item billing.Subscription, quantity int64) bool {
	for _, lineItem := range item.Items {
		if lineItem.Quantity == quantity {
			return true
		}
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
