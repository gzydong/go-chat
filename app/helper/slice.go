package helper

import (
	"errors"
	"fmt"
	"reflect"
)

// UniqueSliceString 唯一的字符串切片
func UniqueSliceString(arr []string) []string {
	result := make([]string, 0)

	hash := make(map[string]int)

	for _, value := range arr {
		if _, ok := hash[value]; !ok {
			hash[value] = 0
		}
	}

	for str, _ := range hash {
		result = append(result, str)
	}

	return result
}

// UniqueSliceString 唯一的 Int 切片
func UniqueSliceInt(arr []int) []int {
	result := make([]int, 0)

	return result
}

func SliceToMap(arr []map[string]interface{}, field string) (map[int]map[string]interface{}, error) {
	hashMap := make(map[int]map[string]interface{}, len(arr))

	for _, data := range arr {
		value, ok := data[field]
		if !ok {
			return nil, errors.New(fmt.Sprintf("%s 字段不存在", field))
		}

		if _, ok := value.(int64); ok {
			hashMap[int(reflect.ValueOf(value).Int())] = data
		} else {
			return nil, errors.New(fmt.Sprintf("%s 字段非 int64 类型", field))
		}
	}

	return hashMap, nil
}
