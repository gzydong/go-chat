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

// Decode
// nolint
func Decode(data any, resp any) error {
	switch data.(type) {
	case string:
		return json.Unmarshal([]byte(data.(string)), resp)
	case []byte:
		return json.Unmarshal(data.([]byte), resp)
	default:
		return errors.New("未知类型")
	}
}
