package helper

import jsoniter "github.com/json-iterator/go"

func JsonEncode(str string, value interface{}) error {
	return jsoniter.UnmarshalFromString(str, value)
}

func JsonDecode(value interface{}) string {
	content, _ := jsoniter.MarshalToString(value)
	return content
}
