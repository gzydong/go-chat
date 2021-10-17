package helper

import jsoniter "github.com/json-iterator/go"

func JsonDecode(str string, value interface{}) error {
	return jsoniter.UnmarshalFromString(str, value)
}

func JsonEncode(value interface{}) string {
	content, _ := jsoniter.MarshalToString(value)
	return content
}
