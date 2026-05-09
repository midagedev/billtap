package fixtures

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

func LoadPack(body []byte, contentType string) (Pack, error) {
	var pack Pack
	if isYAMLContentType(contentType) {
		dec := yaml.NewDecoder(bytes.NewReader(body))
		dec.KnownFields(true)
		if err := dec.Decode(&pack); err != nil {
			return Pack{}, fmt.Errorf("decode fixture YAML: %w", err)
		}
		return pack, nil
	}
	if err := json.Unmarshal(body, &pack); err != nil {
		return Pack{}, fmt.Errorf("decode fixture JSON: %w", err)
	}
	return pack, nil
}

func LoadAssertionRequest(body []byte, contentType string) (AssertionRequest, error) {
	var req AssertionRequest
	if isYAMLContentType(contentType) {
		dec := yaml.NewDecoder(bytes.NewReader(body))
		dec.KnownFields(true)
		if err := dec.Decode(&req); err != nil {
			return AssertionRequest{}, fmt.Errorf("decode assertion YAML: %w", err)
		}
		return req, nil
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return AssertionRequest{}, fmt.Errorf("decode assertion JSON: %w", err)
	}
	return req, nil
}

func isYAMLContentType(contentType string) bool {
	return bytes.Contains([]byte(contentType), []byte("yaml")) || bytes.Contains([]byte(contentType), []byte("yml"))
}
