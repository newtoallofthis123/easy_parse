package utils

import (
	"encoding/json"
)

// parseValue attempts to parse value v recursively.
// If v is a string that contains valid JSON,
// then it returns the parsed JSON. Otherwise, it returns v.
func parseValue(v any) any {
	switch val := v.(type) {
	case string:
		var parsed any
		if err := json.Unmarshal([]byte(val), &parsed); err == nil {
			return parseValue(parsed)
		}
		return val
	case []any:
		for i, elem := range val {
			val[i] = parseValue(elem)
		}
		return val
	case map[string]any:
		res := make(map[string]any)
		for key, elem := range val {
			res[key] = parseValue(elem)
		}
		return res
	default:
		return val
	}
}

// DecodeAndParse accepts raw JSON data,
// decodes it into a generic map and then recursively
// processes its values.
func DecodeAndParse(jsonData []byte) (any, error) {
	var data any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	parsedData := parseValue(data)
	return parsedData, nil
}
