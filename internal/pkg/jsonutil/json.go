package jsonutil

import (
	"encoding/json"
	"errors"
)

func Encode(value any) string {
	data, _ := json.Marshal(value)
	return string(data)
}

func Marshal(value any) []byte {
	data, _ := json.Marshal(value)
	return data
}

// Unmarshal
// nolint
func Unmarshal(data any, v any) error {
	switch data.(type) {
	case string:
		return json.Unmarshal([]byte(data.(string)), v)
	case []byte:
		return json.Unmarshal(data.([]byte), v)
	default:
		return errors.New("未知类型")
	}
}
