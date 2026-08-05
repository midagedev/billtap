package billing

import "time"

// ProrationFactor returns the unused fraction of a billing period as
// (remaining, total) seconds. zero, false when the period is unusable
// (zero bounds, already ended, inverted).
func ProrationFactor(periodStart, periodEnd, at time.Time) (remaining int64, total int64, ok bool) {
	if periodStart.IsZero() || periodEnd.IsZero() {
		return 0, 0, false
	}
	if !periodEnd.After(periodStart) {
		return 0, 0, false
	}
	if !periodEnd.After(at) {
		return 0, 0, false
	}
	total = periodEnd.Unix() - periodStart.Unix()
	remaining = periodEnd.Unix() - at.Unix()
	if total <= 0 || remaining <= 0 {
		return 0, 0, false
	}
	return remaining, total, true
}

// ProrateDelta scales an amount delta by the remaining period fraction
// (truncating integer division, multiply first — same as Stripe / invoicePreview).
func ProrateDelta(delta int64, remaining, total int64) int64 {
	if total <= 0 || remaining <= 0 {
		return 0
	}
	return delta * remaining / total
}
