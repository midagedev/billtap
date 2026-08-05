package billing

import (
	"strconv"
	"strings"
)

// FormatSubscriptionItemID builds the stable si_<sub>_<n> form used on the API surface.
// The sanitize rules match the historical API sanitizeID (alphanumeric, underscore, hyphen).
func FormatSubscriptionItemID(subscriptionID string, index int) string {
	return "si_" + sanitizeSubscriptionItemToken(subscriptionID+"_"+strconv.Itoa(index))
}

// sanitizeSubscriptionItemToken keeps [A-Za-z0-9_-]; empty input becomes "billtap".
// Shared with the API layer so stored and derived IDs always match.
func sanitizeSubscriptionItemToken(value string) string {
	var b strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "billtap"
	}
	return b.String()
}

// AssignSubscriptionItemIDs fills empty item IDs with stable si_<sub>_<n> values.
// n is the lowest index not already used by the item set, so IDs never shift when
// items are removed. Existing IDs are preserved.
func AssignSubscriptionItemIDs(subscriptionID string, items []LineItem) []LineItem {
	if len(items) == 0 {
		return items
	}
	out := make([]LineItem, len(items))
	copy(out, items)

	used := map[int]struct{}{}
	for _, item := range out {
		id := strings.TrimSpace(item.ID)
		if id == "" {
			continue
		}
		if idx, ok := parseSubscriptionItemIndex(subscriptionID, id); ok {
			used[idx] = struct{}{}
		}
	}

	next := 0
	for i := range out {
		if strings.TrimSpace(out[i].ID) != "" {
			continue
		}
		for {
			if _, taken := used[next]; !taken {
				break
			}
			next++
		}
		out[i].ID = FormatSubscriptionItemID(subscriptionID, next)
		used[next] = struct{}{}
		next++
	}
	return out
}

// parseSubscriptionItemIndex returns n when id equals FormatSubscriptionItemID(sub, n).
// sanitizeSubscriptionItemToken filters per rune and digits survive it, so the formatted
// ID always splits into a fixed prefix plus the decimal index.
func parseSubscriptionItemIndex(subscriptionID, id string) (int, bool) {
	prefix := "si_" + sanitizeSubscriptionItemToken(subscriptionID+"_")
	rest, ok := strings.CutPrefix(id, prefix)
	if !ok || rest == "" {
		return 0, false
	}
	index, err := strconv.Atoi(rest)
	if err != nil || index < 0 {
		return 0, false
	}
	return index, true
}

// ResolvedSubscriptionItemID returns the stored ID when present, otherwise the
// index-derived legacy ID so pre-migration rows stay byte-stable on read.
func ResolvedSubscriptionItemID(subscriptionID string, item LineItem, index int) string {
	if id := strings.TrimSpace(item.ID); id != "" {
		return id
	}
	return FormatSubscriptionItemID(subscriptionID, index)
}
