package jsonutil

import jsoniter "github.com/json-iterator/go"

func Encode(value interface{}) string {
	data, _ := jsoniter.MarshalToString(value)
	return data
}

func EncodeToBt(value interface{}) []byte {
	data, _ := jsoniter.Marshal(value)
	return data
}

func Decode(str string, value interface{}) error {
	return jsoniter.UnmarshalFromString(str, value)
}
