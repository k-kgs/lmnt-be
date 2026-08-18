package service

import (
	"encoding/json"
	"fmt"
)

// FieldDef mirrors one entry of a vertical's input_schema (see root .adlc design doc §3.1).
type FieldDef struct {
	Key     string   `json:"key"`
	Type    string   `json:"type"`
	Label   string   `json:"label"`
	Unit    string   `json:"unit,omitempty"`
	Options []string `json:"options,omitempty"`
}

// ValidateAgainstSchema confirms metricData carries every required field from
// schema with a JSON type matching that field's declared type. Intentionally
// a small hand-rolled check, not a full JSON Schema library — the field-type
// vocabulary is fixed and small (number|text|select|photo|exercise_list).
func ValidateAgainstSchema(schemaJSON json.RawMessage, metricData map[string]any) error {
	var fields []FieldDef
	if err := json.Unmarshal(schemaJSON, &fields); err != nil {
		return fmt.Errorf("invalid vertical schema: %w", err)
	}

	for _, f := range fields {
		v, ok := metricData[f.Key]
		if !ok {
			return fmt.Errorf("missing required field %q", f.Key)
		}
		if err := validateFieldType(f, v); err != nil {
			return fmt.Errorf("field %q: %w", f.Key, err)
		}
	}
	return nil
}

func validateFieldType(f FieldDef, v any) error {
	switch f.Type {
	case "number":
		if _, ok := v.(float64); !ok {
			return fmt.Errorf("expected number")
		}
	case "text", "photo":
		s, ok := v.(string)
		if !ok || s == "" {
			return fmt.Errorf("expected non-empty string")
		}
	case "select":
		s, ok := v.(string)
		if !ok {
			return fmt.Errorf("expected string")
		}
		if len(f.Options) > 0 && !contains(f.Options, s) {
			return fmt.Errorf("value %q not in allowed options %v", s, f.Options)
		}
	case "exercise_list":
		if _, ok := v.([]any); !ok {
			return fmt.Errorf("expected an array")
		}
	default:
		return fmt.Errorf("unknown field type %q", f.Type)
	}
	return nil
}

func contains(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}
