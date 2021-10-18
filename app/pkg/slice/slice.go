package slice

import (
	"errors"
	"fmt"
	"reflect"
)

// ToMap
func ToMap(arr []map[string]interface{}, field string) (map[int64]map[string]interface{}, error) {
	hashMap := make(map[int64]map[string]interface{}, len(arr))

	for _, data := range arr {
		value, ok := data[field]
		if !ok {
			return nil, errors.New(fmt.Sprintf("%s 字段不存在", field))
		}

		if _, ok := value.(int64); ok {
			hashMap[reflect.ValueOf(value).Int()] = data
		} else {
			return nil, errors.New(fmt.Sprintf("%s 字段非 int64 类型", field))
		}
	}

	return hashMap, nil
}
