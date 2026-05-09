package security

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const MaskedValue = "****"

var (
	urlPattern      = regexp.MustCompile(`https?://[^\s"']+`)
	keyValuePattern = regexp.MustCompile(`(?i)(api[_-]?key|secret|token)=([^&\s]+)`)
)

func MaskSecret(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return MaskedValue
}

func RedactHeaders(headers map[string]string) map[string]string {
	if headers == nil {
		return nil
	}
	out := make(map[string]string, len(headers))
	for key, value := range headers {
		if normalize(key) == "billtapsignature" {
			out[key] = MaskSignatureHeader(value)
			continue
		}
		if IsSensitiveKey(key) {
			out[key] = MaskSecret(value)
			continue
		}
		out[key] = value
	}
	return out
}

func MaskSignatureHeader(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	parts := strings.Split(value, ",")
	for i, part := range parts {
		key, raw, found := strings.Cut(strings.TrimSpace(part), "=")
		if !found {
			continue
		}
		if key == "t" {
			parts[i] = key + "=" + raw
			continue
		}
		parts[i] = key + "=" + MaskedValue
	}
	return strings.Join(parts, ",")
}

func RedactURL(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	parsed, err := url.Parse(value)
	if err != nil {
		return value
	}
	if parsed.User != nil {
		username := parsed.User.Username()
		if _, hasPassword := parsed.User.Password(); hasPassword {
			parsed.User = url.UserPassword(username, MaskedValue)
		}
	}
	query := parsed.Query()
	for key, values := range query {
		if IsSensitiveKey(key) {
			for i := range values {
				values[i] = MaskedValue
			}
			query[key] = values
		}
	}
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func RedactText(value string) string {
	if strings.TrimSpace(value) == "" {
		return value
	}
	var decoded any
	if err := json.Unmarshal([]byte(value), &decoded); err != nil {
		return redactFreeform(value)
	}
	redacted := redactAny("", decoded)
	body, err := json.Marshal(redacted)
	if err != nil {
		return value
	}
	return string(body)
}

func redactFreeform(value string) string {
	value = urlPattern.ReplaceAllStringFunc(value, RedactURL)
	value = keyValuePattern.ReplaceAllString(value, `${1}=****`)
	return value
}

func IsSensitiveKey(key string) bool {
	normalized := normalize(key)
	switch normalized {
	case "authorization", "proxyauthorization", "cookie", "setcookie", "xapikey", "apikey":
		return true
	}
	return strings.Contains(normalized, "secret") ||
		strings.Contains(normalized, "accesstoken") ||
		strings.Contains(normalized, "refreshtoken") ||
		strings.Contains(normalized, "bearertoken")
}

func ContainsCardData(values map[string]string) bool {
	for key := range values {
		if IsCardDataKey(key) {
			return true
		}
	}
	return false
}

func ContainsCardDataAny(value any) bool {
	return containsCardDataAny("", value)
}

func IsCardDataKey(key string) bool {
	normalized := normalize(key)
	if strings.Contains(normalized, "paymentmethoddatacard") {
		return true
	}
	if !strings.Contains(normalized, "card") {
		return false
	}
	for _, marker := range []string{"number", "cvc", "cvv", "expmonth", "expyear", "expiry"} {
		if strings.Contains(normalized, marker) {
			return true
		}
	}
	return false
}

func containsCardDataAny(path string, value any) bool {
	switch v := value.(type) {
	case map[string]any:
		for key, child := range v {
			childPath := key
			if path != "" {
				childPath = path + "." + key
			}
			if IsCardDataKey(childPath) || containsCardDataAny(childPath, child) {
				return true
			}
		}
	case []any:
		for idx, child := range v {
			if containsCardDataAny(fmt.Sprintf("%s.%d", path, idx), child) {
				return true
			}
		}
	}
	return false
}

func redactAny(path string, value any) any {
	switch v := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for key, child := range v {
			childPath := key
			if path != "" {
				childPath = path + "." + key
			}
			if IsSensitiveKey(key) || IsCardDataKey(childPath) {
				out[key] = MaskedValue
				continue
			}
			out[key] = redactAny(childPath, child)
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, child := range v {
			out[i] = redactAny(fmt.Sprintf("%s.%d", path, i), child)
		}
		return out
	default:
		return value
	}
}

func normalize(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer("_", "", "-", "", ".", "", "[", "", "]", "", " ", "")
	return replacer.Replace(value)
}
