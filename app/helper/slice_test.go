package helper

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlice_UniqueSliceString(t *testing.T) {
	arr := []string{
		"1111",
		"2222",
		"1111",
		"3333",
	}

	res := UniqueSliceString(arr)

	assert.Equal(t, true, len(res) == 3)
}

func TestSlice_SliceToMap(t *testing.T) {
	items := make([]map[string]interface{}, 0)

	items = append(items, map[string]interface{}{
		"id":   int64(12),
		"name": "那就开始拿大",
	})

	items = append(items, map[string]interface{}{
		"id":   2,
		"name": "2那就开始拿大",
	})

	arr, err := SliceToMap(items, "id")
	fmt.Println(err)
	fmt.Println(arr)
	fmt.Println(arr[2])
}
