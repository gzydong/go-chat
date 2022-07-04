package sliceutil

import (
	"strconv"
	"strings"
)

func UniqueInt(data []int) []int {
	list := make([]int, 0)
	hash := make(map[int]int)

	for _, value := range data {
		if _, ok := hash[value]; !ok {
			list = append(list, value)
			hash[value] = 0
		}
	}

	return list
}

func UniqueInt64(data []int64) []int64 {
	list := make([]int64, 0)
	hash := make(map[int64]int)

	for _, value := range data {
		if _, ok := hash[value]; !ok {
			list = append(list, value)
			hash[value] = 0
		}
	}

	return list
}

func UniqueString(data []string) []string {
	list := make([]string, 0)
	hash := make(map[string]int)

	for _, value := range data {
		if _, ok := hash[value]; !ok {
			list = append(list, value)
			hash[value] = 0
		}
	}

	return list
}

func ParseIds(str string) []int {
	str = strings.TrimSpace(str)
	ids := make([]int, 0)

	if str == "" {
		return ids
	}

	for _, value := range strings.Split(str, ",") {
		if id, err := strconv.Atoi(value); err == nil {
			ids = append(ids, id)
		}
	}

	return ids
}

func IntToIds(items []int) string {
	tmp := make([]string, 0, len(items))

	for _, item := range items {
		tmp = append(tmp, strconv.Itoa(item))
	}

	return strings.Join(tmp, ",")
}
