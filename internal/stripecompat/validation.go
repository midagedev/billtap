package stripecompat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	ValidationCodeParamInvalid = "parameter_invalid"
	ValidationCodeParamMissing = "parameter_missing"
	ValidationCodeParamUnknown = "parameter_unknown"
)

type OperationValidation struct {
	Method               string
	Path                 string
	OperationID          string
	Params               []ParameterValidation
	AdditionalProperties bool
}

type ParameterValidation struct {
	Name                 string
	Location             string
	Required             bool
	Type                 string
	Enum                 []string
	Children             []ParameterValidation
	AdditionalProperties bool
}

type ValidationError struct {
	Param   string
	Code    string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type ValidationCatalog struct {
	operations map[string]OperationValidation
	ordered    []OperationValidation
}

func NewValidationCatalog(operations []OperationValidation) (ValidationCatalog, error) {
	catalog := ValidationCatalog{operations: map[string]OperationValidation{}}
	for _, operation := range operations {
		operation.Method = strings.ToUpper(strings.TrimSpace(operation.Method))
		operation.Path = NormalizePath(operation.Path)
		if operation.Method == "" {
			return ValidationCatalog{}, fmt.Errorf("validation operation for %s has empty method", operation.Path)
		}
		if operation.Path == "" {
			return ValidationCatalog{}, fmt.Errorf("validation operation for %s has empty path", operation.Method)
		}
		key := RouteKey(operation.Method, operation.Path)
		if _, exists := catalog.operations[key]; exists {
			return ValidationCatalog{}, fmt.Errorf("duplicate validation operation: %s", key)
		}
		copied := cloneOperationValidation(operation)
		catalog.operations[key] = copied
		catalog.ordered = append(catalog.ordered, copied)
	}
	return catalog, nil
}

func MustValidationCatalog(operations []OperationValidation) ValidationCatalog {
	catalog, err := NewValidationCatalog(operations)
	if err != nil {
		panic(err)
	}
	return catalog
}

func DefaultValidationCatalog() ValidationCatalog {
	return MustValidationCatalog(DefaultValidationOperations())
}

func (c ValidationCatalog) Lookup(method string, path string) (OperationValidation, bool) {
	operation, ok := c.operations[RouteKey(method, NormalizePath(path))]
	if ok {
		return cloneOperationValidation(operation), true
	}
	normalizedMethod := strings.ToUpper(strings.TrimSpace(method))
	for _, operation := range c.ordered {
		if operation.Method == normalizedMethod && pathMatches(operation.Path, path) {
			return cloneOperationValidation(operation), true
		}
	}
	return OperationValidation{}, false
}

func (c ValidationCatalog) Operations() []OperationValidation {
	out := make([]OperationValidation, 0, len(c.ordered))
	for _, operation := range c.ordered {
		out = append(out, cloneOperationValidation(operation))
	}
	return out
}

func (c ValidationCatalog) Validate(r *http.Request) *ValidationError {
	operation, ok := c.Lookup(r.Method, r.URL.Path)
	if !ok {
		return nil
	}
	values, err := requestValues(r)
	if err != nil {
		return &ValidationError{
			Param:   "",
			Code:    ValidationCodeParamInvalid,
			Message: "Invalid request body.",
		}
	}
	return validateValues(operation, values)
}

func validateValues(operation OperationValidation, values map[string][]string) *ValidationError {
	params := map[string]ParameterValidation{}
	for _, param := range operation.Params {
		if param.Location == "path" {
			continue
		}
		params[param.Name] = param
	}

	for key, items := range values {
		if len(items) == 0 {
			continue
		}
		root := rootParamName(key)
		param, ok := params[root]
		if !ok {
			return unknownValidationParam(key)
		}
		if err := validateParamValue(param, key, items); err != nil {
			return err
		}
	}

	for _, param := range operation.Params {
		if param.Location == "path" || !param.Required {
			continue
		}
		if !hasRootParam(values, param.Name) {
			return missingValidationParam(param.Name)
		}
	}
	return nil
}

func requestValues(r *http.Request) (map[string][]string, error) {
	values := map[string][]string{}
	mergeValues(values, r.URL.Query())
	if r.Method == http.MethodGet || r.Method == http.MethodDelete || r.Body == nil {
		return values, nil
	}

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewReader(raw))
	if len(bytes.TrimSpace(raw)) == 0 {
		return values, nil
	}

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		decoded := map[string]any{}
		decoder := json.NewDecoder(bytes.NewReader(raw))
		decoder.UseNumber()
		if err := decoder.Decode(&decoded); err != nil {
			return nil, err
		}
		mergeAnyValues(values, "", decoded)
		return values, nil
	}

	form, err := url.ParseQuery(string(raw))
	if err != nil {
		return nil, err
	}
	mergeValues(values, form)
	return values, nil
}

func validateParamValue(param ParameterValidation, key string, values []string) *ValidationError {
	return validateParamPath(param, key, paramPathSegments(key), values)
}

func validateScalarParamValue(param ParameterValidation, key string, values []string) *ValidationError {
	switch param.Type {
	case "", "any", "string", "object":
	case "integer":
		for _, value := range values {
			if _, err := strconv.ParseInt(value, 10, 64); err != nil {
				return invalidValidationParam(key, "Expected an integer.")
			}
		}
	case "boolean":
		for _, value := range values {
			switch value {
			case "true", "false", "1", "0":
			default:
				return invalidValidationParam(key, "Expected a boolean.")
			}
		}
	case "array":
		return nil
	default:
		return nil
	}
	if len(param.Enum) > 0 {
		allowed := map[string]struct{}{}
		for _, value := range param.Enum {
			allowed[value] = struct{}{}
		}
		for _, value := range values {
			if _, ok := allowed[value]; !ok {
				return invalidValidationParam(key, "Expected one of: "+strings.Join(param.Enum, ", ")+".")
			}
		}
	}
	return nil
}

func validateParamPath(param ParameterValidation, key string, segments []string, values []string) *ValidationError {
	if len(segments) <= 1 {
		return validateScalarParamValue(param, key, values)
	}
	if param.Type == "array" {
		return nil
	}
	if param.Type != "object" && len(param.Children) == 0 {
		return unknownValidationParam(key)
	}
	childName := segments[1]
	for _, child := range param.Children {
		if child.Name == childName {
			return validateParamPath(child, key, segments[1:], values)
		}
	}
	if param.AdditionalProperties {
		return nil
	}
	return unknownValidationParam(key)
}

func paramPathSegments(key string) []string {
	segments := []string{}
	var b strings.Builder
	for _, r := range key {
		switch r {
		case '[':
			if b.Len() > 0 {
				segments = append(segments, b.String())
				b.Reset()
			}
		case ']':
			if b.Len() > 0 {
				segments = append(segments, b.String())
				b.Reset()
			}
		default:
			b.WriteRune(r)
		}
	}
	if b.Len() > 0 {
		segments = append(segments, b.String())
	}
	if len(segments) == 0 {
		return []string{key}
	}
	return segments
}

func mergeValues(dst map[string][]string, src url.Values) {
	for key, values := range src {
		dst[key] = append(dst[key], values...)
	}
}

func mergeAnyValues(dst map[string][]string, prefix string, value any) {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			next := key
			if prefix != "" {
				next = prefix + "[" + key + "]"
			}
			mergeAnyValues(dst, next, child)
		}
	case []any:
		for i, child := range typed {
			next := prefix + "[" + strconv.Itoa(i) + "]"
			mergeAnyValues(dst, next, child)
		}
	case json.Number:
		dst[prefix] = append(dst[prefix], typed.String())
	case bool:
		dst[prefix] = append(dst[prefix], strconv.FormatBool(typed))
	case nil:
		dst[prefix] = append(dst[prefix], "")
	default:
		dst[prefix] = append(dst[prefix], fmt.Sprint(typed))
	}
}

func hasRootParam(values map[string][]string, name string) bool {
	for key, items := range values {
		if len(items) == 0 {
			continue
		}
		if rootParamName(key) == name {
			for _, item := range items {
				if strings.TrimSpace(item) != "" {
					return true
				}
			}
		}
	}
	return false
}

func rootParamName(key string) string {
	if idx := strings.IndexByte(key, '['); idx >= 0 {
		return key[:idx]
	}
	return key
}

func missingValidationParam(param string) *ValidationError {
	return &ValidationError{
		Param:   param,
		Code:    ValidationCodeParamMissing,
		Message: fmt.Sprintf("Missing required param: %s.", param),
	}
}

func unknownValidationParam(param string) *ValidationError {
	return &ValidationError{
		Param:   param,
		Code:    ValidationCodeParamUnknown,
		Message: fmt.Sprintf("Received unknown parameter: %s.", param),
	}
}

func invalidValidationParam(param string, reason string) *ValidationError {
	return &ValidationError{
		Param:   param,
		Code:    ValidationCodeParamInvalid,
		Message: fmt.Sprintf("Invalid param: %s. %s", param, reason),
	}
}

func cloneOperationValidation(operation OperationValidation) OperationValidation {
	operation.Params = cloneParameterValidations(operation.Params)
	return operation
}

func cloneParameterValidations(params []ParameterValidation) []ParameterValidation {
	out := make([]ParameterValidation, len(params))
	for i, param := range params {
		out[i] = param
		out[i].Enum = append([]string(nil), param.Enum...)
		out[i].Children = cloneParameterValidations(param.Children)
	}
	return out
}
