package compatibility

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hckim/billtap/internal/api"
	"github.com/hckim/billtap/internal/billing"
	"github.com/hckim/billtap/internal/storage"
	"github.com/hckim/billtap/internal/webhooks"
)

const (
	ScorecardVersion = "l3-public-readiness-v4"
	DefaultOutputDir = "dist/compatibility"
)

type Status string

const (
	StatusImported    Status = "imported"
	StatusSkipped     Status = "skipped"
	StatusUnsupported Status = "unsupported"
	StatusMismatch    Status = "mismatch"
	StatusError       Status = "error"
)

var statusDefinitions = map[Status]string{
	StatusImported:    "case ran against Billtap and matched the normalized expectation",
	StatusSkipped:     "case is in the corpus but was not run by this offline scorecard",
	StatusUnsupported: "case documents an unsupported provider behavior from docs/COMPATIBILITY.md",
	StatusMismatch:    "case ran but normalized actual behavior differed from the expectation",
	StatusError:       "case could not run because the scorecard runner or Billtap returned an unexpected internal error",
}

type Options struct {
	OutputDir string
	Now       func() time.Time
	cases     []caseSpec
}

type Scorecard struct {
	Name             string            `json:"name"`
	ScorecardVersion string            `json:"scorecard_version"`
	GeneratedAt      time.Time         `json:"generated_at"`
	StatusMap        map[string]string `json:"status_map"`
	Summary          Summary           `json:"summary"`
	Cases            []CaseResult      `json:"cases"`
}

type Summary struct {
	Total                   int  `json:"total"`
	Imported                int  `json:"imported"`
	Skipped                 int  `json:"skipped"`
	Unsupported             int  `json:"unsupported"`
	Mismatch                int  `json:"mismatch"`
	Error                   int  `json:"error"`
	ReleaseBlocking         int  `json:"release_blocking"`
	ReleaseBlockingMismatch int  `json:"release_blocking_mismatch"`
	ReleaseBlockingError    int  `json:"release_blocking_error"`
	Passed                  bool `json:"passed"`
}

type CaseResult struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Category        string       `json:"category"`
	Level           string       `json:"level"`
	Status          Status       `json:"status"`
	ReleaseBlocking bool         `json:"release_blocking"`
	Reference       string       `json:"reference,omitempty"`
	Reason          string       `json:"reason,omitempty"`
	Expected        *Observation `json:"expected,omitempty"`
	Actual          *Observation `json:"actual,omitempty"`
	ReplayBundle    string       `json:"replay_bundle,omitempty"`
	Duration        string       `json:"duration,omitempty"`

	bundle *ReplayBundle
}

type Observation struct {
	HTTPStatus          int               `json:"http_status,omitempty"`
	Error               *ErrorObservation `json:"error,omitempty"`
	Object              string            `json:"object,omitempty"`
	ObjectStatus        string            `json:"object_status,omitempty"`
	PaymentStatus       string            `json:"payment_status,omitempty"`
	PaymentIntentStatus string            `json:"payment_intent_status,omitempty"`
	PaymentIntentError  *ErrorObservation `json:"payment_intent_error,omitempty"`
}

type ErrorObservation struct {
	Type        string `json:"type,omitempty"`
	Code        string `json:"code,omitempty"`
	Param       string `json:"param,omitempty"`
	DeclineCode string `json:"decline_code,omitempty"`
}

type ReplayBundle struct {
	CaseID      string       `json:"case_id"`
	Status      Status       `json:"status"`
	Reason      string       `json:"reason"`
	GeneratedAt time.Time    `json:"generated_at"`
	Steps       []ReplayStep `json:"steps"`
}

type ReplayStep struct {
	Name     string            `json:"name"`
	Method   string            `json:"method"`
	Path     string            `json:"path"`
	Headers  map[string]string `json:"headers,omitempty"`
	Params   map[string]string `json:"params,omitempty"`
	JSON     any               `json:"json,omitempty"`
	Expected *Observation      `json:"expected,omitempty"`
	Actual   *Observation      `json:"actual,omitempty"`
	Response ReplayResponse    `json:"response"`
}

type ReplayResponse struct {
	HTTPStatus int    `json:"http_status"`
	JSON       any    `json:"json,omitempty"`
	BodyText   string `json:"body_text,omitempty"`
}

type caseSpec struct {
	ID              string
	Name            string
	Category        string
	Level           string
	ReleaseBlocking bool
	Reference       string
	SkipReason      string
	Unsupported     string
	Steps           []requestSpec
	Expect          Observation
	Run             func(context.Context, *harness) (caseExecution, error)
}

type requestSpec struct {
	Name             string
	Method           string
	Path             string
	Headers          map[string]string
	Params           map[string]string
	JSON             any
	ExpectHTTPStatus int
}

type caseExecution struct {
	Actual Observation
	Steps  []ReplayStep
}

func Generate(ctx context.Context, opts Options) (Scorecard, error) {
	now := opts.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	cases := opts.cases
	if len(cases) == 0 {
		cases = builtinCorpus()
	}

	scorecard := Scorecard{
		Name:             "Billtap Stripe-like compatibility scorecard",
		ScorecardVersion: ScorecardVersion,
		GeneratedAt:      now(),
		StatusMap:        stringStatusDefinitions(),
		Cases:            make([]CaseResult, 0, len(cases)),
	}
	for _, spec := range cases {
		result := runCase(ctx, spec, now)
		scorecard.Summary.add(result)
		scorecard.Cases = append(scorecard.Cases, result)
	}
	scorecard.Summary.Passed = scorecard.Summary.ReleaseBlockingMismatch == 0 && scorecard.Summary.ReleaseBlockingError == 0
	return scorecard, nil
}

func WriteArtifacts(ctx context.Context, opts Options) (Scorecard, error) {
	outputDir := opts.OutputDir
	if outputDir == "" {
		outputDir = DefaultOutputDir
	}
	scorecard, err := Generate(ctx, opts)
	if err != nil {
		return scorecard, err
	}
	if err := os.MkdirAll(filepath.Join(outputDir, "replay-bundles"), 0o755); err != nil {
		return scorecard, err
	}
	body, err := scorecard.JSON()
	if err != nil {
		return scorecard, err
	}
	if err := os.WriteFile(filepath.Join(outputDir, "compatibility-scorecard.json"), body, 0o644); err != nil {
		return scorecard, err
	}
	if err := os.WriteFile(filepath.Join(outputDir, "compatibility-scorecard.md"), []byte(scorecard.Markdown()), 0o644); err != nil {
		return scorecard, err
	}
	for _, result := range scorecard.Cases {
		if result.bundle == nil || result.ReplayBundle == "" {
			continue
		}
		body, err := json.MarshalIndent(result.bundle, "", "  ")
		if err != nil {
			return scorecard, err
		}
		if err := os.WriteFile(filepath.Join(outputDir, result.ReplayBundle), body, 0o644); err != nil {
			return scorecard, err
		}
	}
	return scorecard, nil
}

func (s Scorecard) JSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

func (s Scorecard) Markdown() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "# Compatibility Scorecard\n\n")
	fmt.Fprintf(&b, "- Version: `%s`\n", s.ScorecardVersion)
	fmt.Fprintf(&b, "- Generated: `%s`\n", s.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&b, "- Release blocking passed: `%t`\n", s.Summary.Passed)
	fmt.Fprintf(&b, "- Counts: imported `%d`, skipped `%d`, unsupported `%d`, mismatch `%d`, error `%d`\n",
		s.Summary.Imported, s.Summary.Skipped, s.Summary.Unsupported, s.Summary.Mismatch, s.Summary.Error)

	b.WriteString("\n## Status Map\n\n")
	for _, status := range []Status{StatusImported, StatusSkipped, StatusUnsupported, StatusMismatch, StatusError} {
		fmt.Fprintf(&b, "- `%s`: %s\n", status, statusDefinitions[status])
	}

	b.WriteString("\n## Cases\n\n")
	b.WriteString("| Case | Category | Level | Status | Release blocking | Reason |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- |\n")
	for _, result := range s.Cases {
		reason := result.Reason
		if result.ReplayBundle != "" {
			reason = strings.TrimSpace(reason + " replay: " + result.ReplayBundle)
		}
		fmt.Fprintf(&b, "| %s | %s | %s | `%s` | `%t` | %s |\n",
			escapeTable(result.ID),
			escapeTable(result.Category),
			escapeTable(result.Level),
			result.Status,
			result.ReleaseBlocking,
			escapeTable(reason),
		)
	}
	return b.String()
}

func (s *Summary) add(result CaseResult) {
	s.Total++
	if result.ReleaseBlocking {
		s.ReleaseBlocking++
	}
	switch result.Status {
	case StatusImported:
		s.Imported++
	case StatusSkipped:
		s.Skipped++
	case StatusUnsupported:
		s.Unsupported++
	case StatusMismatch:
		s.Mismatch++
		if result.ReleaseBlocking {
			s.ReleaseBlockingMismatch++
		}
	case StatusError:
		s.Error++
		if result.ReleaseBlocking {
			s.ReleaseBlockingError++
		}
	}
}

func runCase(ctx context.Context, spec caseSpec, now func() time.Time) (result CaseResult) {
	started := now()
	result = CaseResult{
		ID:              spec.ID,
		Name:            spec.Name,
		Category:        spec.Category,
		Level:           spec.Level,
		ReleaseBlocking: spec.ReleaseBlocking,
		Reference:       spec.Reference,
	}
	defer func() {
		result.Duration = now().Sub(started).String()
	}()

	if spec.SkipReason != "" {
		result.Status = StatusSkipped
		result.Reason = spec.SkipReason
		return result
	}
	if spec.Unsupported != "" {
		result.Status = StatusUnsupported
		result.Reason = spec.Unsupported
		return result
	}

	h, err := newHarness(ctx)
	if err != nil {
		result.Status = StatusError
		result.Reason = err.Error()
		result.bundle = newReplayBundle(result, now(), nil)
		result.ReplayBundle = replayBundlePath(result.ID)
		return result
	}
	defer h.Close()

	var execution caseExecution
	if spec.Run != nil {
		execution, err = spec.Run(ctx, h)
	} else {
		execution, err = h.runSteps(spec.Steps)
	}
	result.Expected = &spec.Expect
	result.Actual = &execution.Actual
	if err != nil {
		result.Status = StatusError
		result.Reason = err.Error()
		result.bundle = newReplayBundle(result, now(), execution.Steps)
		result.ReplayBundle = replayBundlePath(result.ID)
		return result
	}

	if mismatches := compareObservations(spec.Expect, execution.Actual); len(mismatches) > 0 {
		result.Status = StatusMismatch
		result.Reason = strings.Join(mismatches, "; ")
		result.bundle = newReplayBundle(result, now(), execution.Steps)
		result.ReplayBundle = replayBundlePath(result.ID)
		return result
	}

	result.Status = StatusImported
	return result
}

func newReplayBundle(result CaseResult, generatedAt time.Time, steps []ReplayStep) *ReplayBundle {
	if len(steps) > 0 {
		last := len(steps) - 1
		steps[last].Expected = result.Expected
		steps[last].Actual = result.Actual
	}
	return &ReplayBundle{
		CaseID:      result.ID,
		Status:      result.Status,
		Reason:      result.Reason,
		GeneratedAt: generatedAt,
		Steps:       steps,
	}
}

func replayBundlePath(id string) string {
	return filepath.ToSlash(filepath.Join("replay-bundles", sanitizeFilename(id)+".json"))
}

func compareObservations(expected Observation, actual Observation) []string {
	var mismatches []string
	if expected.HTTPStatus != 0 && actual.HTTPStatus != expected.HTTPStatus {
		mismatches = append(mismatches, fmt.Sprintf("http_status expected %d got %d", expected.HTTPStatus, actual.HTTPStatus))
	}
	mismatches = append(mismatches, compareError("error", expected.Error, actual.Error)...)
	if expected.Object != "" && actual.Object != expected.Object {
		mismatches = append(mismatches, fmt.Sprintf("object expected %q got %q", expected.Object, actual.Object))
	}
	if expected.ObjectStatus != "" && actual.ObjectStatus != expected.ObjectStatus {
		mismatches = append(mismatches, fmt.Sprintf("object_status expected %q got %q", expected.ObjectStatus, actual.ObjectStatus))
	}
	if expected.PaymentStatus != "" && actual.PaymentStatus != expected.PaymentStatus {
		mismatches = append(mismatches, fmt.Sprintf("payment_status expected %q got %q", expected.PaymentStatus, actual.PaymentStatus))
	}
	if expected.PaymentIntentStatus != "" && actual.PaymentIntentStatus != expected.PaymentIntentStatus {
		mismatches = append(mismatches, fmt.Sprintf("payment_intent_status expected %q got %q", expected.PaymentIntentStatus, actual.PaymentIntentStatus))
	}
	mismatches = append(mismatches, compareError("payment_intent_error", expected.PaymentIntentError, actual.PaymentIntentError)...)
	return mismatches
}

func compareError(label string, expected *ErrorObservation, actual *ErrorObservation) []string {
	if expected == nil {
		return nil
	}
	if actual == nil {
		return []string{label + " expected but missing"}
	}
	var mismatches []string
	if expected.Type != "" && actual.Type != expected.Type {
		mismatches = append(mismatches, fmt.Sprintf("%s.type expected %q got %q", label, expected.Type, actual.Type))
	}
	if expected.Code != "" && actual.Code != expected.Code {
		mismatches = append(mismatches, fmt.Sprintf("%s.code expected %q got %q", label, expected.Code, actual.Code))
	}
	if expected.Param != "" && actual.Param != expected.Param {
		mismatches = append(mismatches, fmt.Sprintf("%s.param expected %q got %q", label, expected.Param, actual.Param))
	}
	if expected.DeclineCode != "" && actual.DeclineCode != expected.DeclineCode {
		mismatches = append(mismatches, fmt.Sprintf("%s.decline_code expected %q got %q", label, expected.DeclineCode, actual.DeclineCode))
	}
	return mismatches
}

type harness struct {
	handler http.Handler
	store   *storage.SQLiteStore
}

func newHarness(ctx context.Context) (*harness, error) {
	store, err := storage.OpenSQLite(ctx, ":memory:")
	if err != nil {
		return nil, err
	}
	h := api.New(api.Options{
		Billing:  billing.NewService(store),
		Webhooks: webhooks.NewService(store),
	})
	return &harness{handler: h, store: store}, nil
}

func (h *harness) Close() {
	if h != nil && h.store != nil {
		_ = h.store.Close()
	}
}

func (h *harness) runSteps(specs []requestSpec) (caseExecution, error) {
	if len(specs) == 0 {
		return caseExecution{}, fmt.Errorf("case has no request steps")
	}
	execution := caseExecution{Steps: make([]ReplayStep, 0, len(specs))}
	for idx, spec := range specs {
		step, err := h.do(spec)
		if err != nil {
			execution.Steps = append(execution.Steps, step)
			return execution, err
		}
		execution.Steps = append(execution.Steps, step)
		execution.Actual = *step.Actual
		if spec.ExpectHTTPStatus != 0 && step.Actual.HTTPStatus != spec.ExpectHTTPStatus {
			return execution, fmt.Errorf("%s returned HTTP %d, want %d", spec.Name, step.Actual.HTTPStatus, spec.ExpectHTTPStatus)
		}
		if idx < len(specs)-1 && spec.ExpectHTTPStatus == 0 && step.Actual.HTTPStatus >= http.StatusBadRequest {
			return execution, fmt.Errorf("%s returned HTTP %d during setup", spec.Name, step.Actual.HTTPStatus)
		}
	}
	return execution, nil
}

func (h *harness) do(spec requestSpec) (ReplayStep, error) {
	method := spec.Method
	if method == "" {
		method = http.MethodPost
	}
	var body *strings.Reader
	contentType := ""
	if spec.JSON != nil {
		raw, err := json.Marshal(spec.JSON)
		if err != nil {
			return ReplayStep{}, err
		}
		body = strings.NewReader(string(raw))
		contentType = "application/json"
	} else {
		values := url.Values{}
		for key, value := range spec.Params {
			values.Set(key, value)
		}
		body = strings.NewReader(values.Encode())
		if method == http.MethodPost {
			contentType = "application/x-www-form-urlencoded"
		}
	}

	req := httptest.NewRequest(method, spec.Path, body)
	req.Host = "billtap.local"
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for key, value := range spec.Headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	h.handler.ServeHTTP(rec, req)

	rawBody := rec.Body.Bytes()
	responseJSON := decodeJSONAny(rawBody)
	actual := normalizeObservation(rec.Code, responseJSON)
	step := ReplayStep{
		Name:     spec.Name,
		Method:   method,
		Path:     spec.Path,
		Headers:  copyStringMap(spec.Headers),
		Params:   copyStringMap(spec.Params),
		JSON:     spec.JSON,
		Actual:   &actual,
		Response: ReplayResponse{HTTPStatus: rec.Code, JSON: responseJSON},
	}
	if responseJSON == nil && len(bytes.TrimSpace(rawBody)) > 0 {
		step.Response.BodyText = strings.TrimSpace(string(rawBody))
	}
	return step, nil
}

func normalizeObservation(status int, decoded any) Observation {
	obs := Observation{HTTPStatus: status}
	m, ok := decoded.(map[string]any)
	if !ok {
		return obs
	}
	if errMap, ok := m["error"].(map[string]any); ok {
		obs.Error = normalizeError(errMap)
	}
	obs.Object = stringField(m, "object")
	obs.ObjectStatus = stringField(m, "status")
	obs.PaymentStatus = stringField(m, "payment_status")
	obs.PaymentIntentStatus = stringField(m, "status")
	if piErr, ok := m["last_payment_error"].(map[string]any); ok {
		obs.PaymentIntentError = normalizeError(piErr)
	}
	return obs
}

func normalizeError(m map[string]any) *ErrorObservation {
	return &ErrorObservation{
		Type:        stringField(m, "type"),
		Code:        stringField(m, "code"),
		Param:       stringField(m, "param"),
		DeclineCode: stringField(m, "decline_code"),
	}
}

func decodeJSONAny(body []byte) any {
	if len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	var decoded any
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil
	}
	return decoded
}

func stringField(m map[string]any, key string) string {
	value, _ := m[key].(string)
	return value
}

func copyStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func stringStatusDefinitions() map[string]string {
	out := make(map[string]string, len(statusDefinitions))
	for status, meaning := range statusDefinitions {
		out[string(status)] = meaning
	}
	return out
}

func escapeTable(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.ReplaceAll(s, "|", "\\|")
}

func sanitizeFilename(id string) string {
	var b strings.Builder
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	if b.Len() == 0 {
		return "case"
	}
	return b.String()
}

func sortedCaseIDs(results []CaseResult) []string {
	ids := make([]string, 0, len(results))
	for _, result := range results {
		ids = append(ids, result.ID)
	}
	sort.Strings(ids)
	return ids
}
