package scenarios

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Report struct {
	Name        string       `json:"name"`
	Status      string       `json:"status"`
	FailureType string       `json:"failure_type,omitempty"`
	StartedAt   time.Time    `json:"started_at"`
	FinishedAt  time.Time    `json:"finished_at"`
	Duration    string       `json:"duration"`
	ClockStart  time.Time    `json:"clock_start"`
	ClockEnd    time.Time    `json:"clock_end"`
	Steps       []StepReport `json:"steps"`
	Errors      []string     `json:"errors,omitempty"`
}

type StepReport struct {
	ID         string         `json:"id"`
	Action     string         `json:"action"`
	Status     string         `json:"status"`
	Clock      time.Time      `json:"clock"`
	StartedAt  time.Time      `json:"started_at"`
	FinishedAt time.Time      `json:"finished_at"`
	Duration   string         `json:"duration"`
	Output     map[string]any `json:"output,omitempty"`
	Error      string         `json:"error,omitempty"`
	Assertion  *Assertion     `json:"assertion,omitempty"`
}

type Assertion struct {
	Target         string         `json:"target"`
	Expected       map[string]any `json:"expected,omitempty"`
	URL            string         `json:"url"`
	Pass           bool           `json:"pass"`
	ResponseStatus int            `json:"response_status,omitempty"`
	ResponseBody   string         `json:"response_body,omitempty"`
	Error          string         `json:"error,omitempty"`
}

func (r Report) ExitCode() int {
	switch r.FailureType {
	case "":
		return ExitPass
	case FailureAssertion:
		return ExitAssertionFailed
	case FailureInvalidConfig:
		return ExitInvalidConfig
	case FailureAppCallback:
		return ExitAppCallbackFailure
	case FailureRunner:
		return ExitRuntimeFailure
	default:
		return ExitAssertionFailed
	}
}

func (r Report) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func (r Report) Markdown() string {
	var b bytes.Buffer
	status := r.Status
	if status == "" {
		status = "unknown"
	}
	fmt.Fprintf(&b, "# Scenario Report: %s\n\n", r.Name)
	fmt.Fprintf(&b, "- Status: `%s`\n", status)
	fmt.Fprintf(&b, "- Exit code: `%d`\n", r.ExitCode())
	if r.FailureType != "" {
		fmt.Fprintf(&b, "- Failure type: `%s`\n", r.FailureType)
	}
	fmt.Fprintf(&b, "- Clock: `%s` to `%s`\n", r.ClockStart.Format(time.RFC3339), r.ClockEnd.Format(time.RFC3339))
	if r.Duration != "" {
		fmt.Fprintf(&b, "- Duration: `%s`\n", r.Duration)
	}
	if len(r.Errors) > 0 {
		b.WriteString("\n## Errors\n\n")
		for _, err := range r.Errors {
			fmt.Fprintf(&b, "- %s\n", err)
		}
	}
	b.WriteString("\n## Steps\n\n")
	b.WriteString("| Step | Action | Status | Clock | Detail |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	for _, step := range r.Steps {
		detail := step.Error
		if detail == "" && step.Assertion != nil {
			detail = assertionDetail(*step.Assertion)
		}
		fmt.Fprintf(&b, "| %s | `%s` | `%s` | `%s` | %s |\n",
			escapeTable(step.ID),
			escapeTable(step.Action),
			escapeTable(step.Status),
			step.Clock.Format(time.RFC3339),
			escapeTable(detail),
		)
	}
	return b.String()
}

func assertionDetail(a Assertion) string {
	if a.Pass {
		return "assertion passed"
	}
	if a.Error != "" {
		return a.Error
	}
	if a.ResponseBody != "" {
		return a.ResponseBody
	}
	return "assertion failed"
}

func escapeTable(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.ReplaceAll(s, "|", "\\|")
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
