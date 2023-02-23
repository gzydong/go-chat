package jsonutil

import (
	"errors"

	jsoniter "github.com/json-iterator/go"
)

func Encode(value any) string {
	data, _ := jsoniter.MarshalToString(value)
	return data
}

func Marshal(value any) []byte {
	data, _ := jsoniter.Marshal(value)
	return data
}

// nolint
func Decode(data any, resp any) error {
	switch data.(type) {
	case string:
		return jsoniter.UnmarshalFromString(data.(string), resp)
	case []byte:
		return jsoniter.Unmarshal(data.([]byte), resp)
	default:
		return errors.New("未知类型")
	}
}
