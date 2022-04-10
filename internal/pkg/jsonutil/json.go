package jsonutil

import jsoniter "github.com/json-iterator/go"

func Decode(str string, value interface{}) error {
	return jsoniter.UnmarshalFromString(str, value)
}

func Encode(value interface{}) string {
	content, _ := jsoniter.MarshalToString(value)
	return content
}

func EncodeByte(value interface{}) (content []byte) {
	content, _ = jsoniter.Marshal(value)
	return
}
