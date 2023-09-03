package jsonutil

import (
	"errors"

	"github.com/bytedance/sonic"
)

func Encode(value any) string {
	data, _ := sonic.MarshalString(value)
	return data
}

func Marshal(value any) []byte {
	data, _ := sonic.Marshal(value)
	return data
}

// nolint
func Decode(data any, resp any) error {
	switch data.(type) {
	case string:
		return sonic.UnmarshalString(data.(string), resp)
	case []byte:
		return sonic.Unmarshal(data.([]byte), resp)
	default:
		return errors.New("未知类型")
	}
}
