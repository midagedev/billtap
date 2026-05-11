package scenarios

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadFile(path string) (Scenario, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return Scenario{}, err
	}
	return Load(body)
}

func Load(body []byte) (Scenario, error) {
	var scenario Scenario
	dec := yaml.NewDecoder(bytes.NewReader(body))
	dec.KnownFields(true)
	if err := dec.Decode(&scenario); err != nil {
		return Scenario{}, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}
	if err := Validate(scenario); err != nil {
		return Scenario{}, err
	}
	return scenario, nil
}

func Validate(s Scenario) error {
	var problems []string
	if strings.TrimSpace(s.Name) == "" {
		problems = append(problems, "name is required")
	}
	if _, err := NewClock(s.Clock.Start); err != nil {
		problems = append(problems, err.Error())
	}
	if len(s.Steps) == 0 {
		problems = append(problems, "steps is required")
	}

	productIDs := map[string]bool{}
	for idx, product := range s.Catalog.Products {
		if strings.TrimSpace(product.ID) == "" {
			problems = append(problems, fmt.Sprintf("catalog.products[%d].id is required", idx))
		}
		if strings.TrimSpace(product.Name) == "" {
			problems = append(problems, fmt.Sprintf("catalog.products[%d].name is required", idx))
		}
		if productIDs[product.ID] {
			problems = append(problems, fmt.Sprintf("catalog.products[%d].id %q is duplicated", idx, product.ID))
		}
		productIDs[product.ID] = true
	}
	priceIDs := map[string]bool{}
	for idx, price := range s.Catalog.Prices {
		if strings.TrimSpace(price.ID) == "" {
			problems = append(problems, fmt.Sprintf("catalog.prices[%d].id is required", idx))
		}
		if strings.TrimSpace(price.Product) == "" {
			problems = append(problems, fmt.Sprintf("catalog.prices[%d].product is required", idx))
		} else if len(productIDs) > 0 && !productIDs[price.Product] {
			problems = append(problems, fmt.Sprintf("catalog.prices[%d].product %q is not defined", idx, price.Product))
		}
		if strings.TrimSpace(price.Currency) == "" {
			problems = append(problems, fmt.Sprintf("catalog.prices[%d].currency is required", idx))
		}
		if price.UnitAmount < 0 {
			problems = append(problems, fmt.Sprintf("catalog.prices[%d].unitAmount must be non-negative", idx))
		}
		if priceIDs[price.ID] {
			problems = append(problems, fmt.Sprintf("catalog.prices[%d].id %q is duplicated", idx, price.ID))
		}
		priceIDs[price.ID] = true
	}

	stepIDs := map[string]bool{}
	for idx, step := range s.Steps {
		id := strings.TrimSpace(step.ID)
		action := strings.TrimSpace(step.Action)
		if id == "" {
			problems = append(problems, fmt.Sprintf("steps[%d].id is required", idx))
		}
		if action == "" {
			problems = append(problems, fmt.Sprintf("steps[%d].action is required", idx))
		} else if !supportedAction(action) {
			problems = append(problems, fmt.Sprintf("steps[%d].action %q is not supported", idx, action))
		}
		if stepIDs[id] {
			problems = append(problems, fmt.Sprintf("steps[%d].id %q is duplicated", idx, id))
		}
		stepIDs[id] = true

		if action == "app.assert" && strings.TrimSpace(s.App.Assertions.BaseURL) == "" {
			problems = append(problems, fmt.Sprintf("steps[%d] app.assert requires app.assertions.baseUrl", idx))
		}
		if action == "clock.advance" {
			if raw, ok := stringParam(step.Params, "duration"); !ok || strings.TrimSpace(raw) == "" {
				problems = append(problems, fmt.Sprintf("steps[%d] clock.advance requires params.duration", idx))
			} else if _, err := ParseDuration(raw); err != nil {
				problems = append(problems, fmt.Sprintf("steps[%d] %v", idx, err))
			}
		}
	}
	if len(problems) > 0 {
		return &InvalidConfigError{Problems: problems}
	}
	return nil
}

func supportedAction(action string) bool {
	switch action {
	case "customer.create",
		"product.create",
		"price.create",
		"checkout.create",
		"checkout.complete",
		"checkout.cancel",
		"subscription.update",
		"subscription.cancel",
		"subscription.resume",
		"clock.advance",
		"invoice.fail_payment",
		"invoice.retry",
		"webhook.replay",
		"app.assert":
		return true
	default:
		return strings.HasPrefix(action, "saas.") || action == "webhook.deliver_duplicate" || action == "webhook.deliver_out_of_order"
	}
}

func stringParam(params map[string]any, key string) (string, bool) {
	if params == nil {
		return "", false
	}
	value, ok := params[key]
	if !ok {
		return "", false
	}
	switch typed := value.(type) {
	case string:
		return typed, true
	default:
		return fmt.Sprint(typed), true
	}
}
