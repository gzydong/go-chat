package entity

import "go-chat/app/pkg/jsonutil"

type JsonText map[string]interface{}

func (j JsonText) Json() string {
	return jsonutil.JsonEncode(j)
}
