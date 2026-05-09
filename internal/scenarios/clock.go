package scenarios

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var durationPart = regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)(ms|d|h|m|s)`)

type Clock struct {
	now time.Time
}

func NewClock(start string) (Clock, error) {
	if strings.TrimSpace(start) == "" {
		return Clock{now: time.Unix(0, 0).UTC()}, nil
	}
	t, err := parseClockStart(start)
	if err != nil {
		return Clock{}, err
	}
	return Clock{now: t.UTC()}, nil
}

func (c Clock) Now() time.Time {
	return c.now
}

func (c *Clock) Advance(raw string) (time.Duration, error) {
	d, err := ParseDuration(raw)
	if err != nil {
		return 0, err
	}
	c.now = c.now.Add(d)
	return d, nil
}

func parseClockStart(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("parse clock.start %q: expected ISO/RFC3339 timestamp", raw)
}

func ParseDuration(raw string) (time.Duration, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, fmt.Errorf("duration is required")
	}
	if d, err := time.ParseDuration(raw); err == nil {
		return d, nil
	}

	clean := strings.ReplaceAll(raw, " ", "")
	matches := durationPart.FindAllStringSubmatchIndex(clean, -1)
	if len(matches) == 0 {
		return 0, fmt.Errorf("parse duration %q: expected values like 3d, 2h, or 30m", raw)
	}

	var total time.Duration
	consumed := 0
	for _, m := range matches {
		if m[0] != consumed {
			return 0, fmt.Errorf("parse duration %q: invalid token near %q", raw, clean[consumed:m[0]])
		}
		valueRaw := clean[m[2]:m[3]]
		unit := strings.ToLower(clean[m[4]:m[5]])
		value, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			return 0, fmt.Errorf("parse duration %q: %w", raw, err)
		}
		var unitDuration time.Duration
		switch unit {
		case "d":
			unitDuration = 24 * time.Hour
		case "h":
			unitDuration = time.Hour
		case "m":
			unitDuration = time.Minute
		case "s":
			unitDuration = time.Second
		case "ms":
			unitDuration = time.Millisecond
		default:
			return 0, fmt.Errorf("parse duration %q: unsupported unit %q", raw, unit)
		}
		total += time.Duration(value * float64(unitDuration))
		consumed = m[1]
	}
	if consumed != len(clean) {
		return 0, fmt.Errorf("parse duration %q: invalid token near %q", raw, clean[consumed:])
	}
	return total, nil
}
