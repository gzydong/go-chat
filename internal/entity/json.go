package entity

import "go-chat/internal/pkg/jsonutil"

type JsonText map[string]interface{}

func (j JsonText) Json() string {
	return jsonutil.Encode(j)
}
