package helper

import (
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
		"id":   int64(123),
		"name": "2那就开始拿大",
	})

	_, err := SliceToMap(items, "id")
	assert.NoError(t, err)

	items2 := make([]map[string]interface{}, 0)
	items2 = append(items2, map[string]interface{}{
		"id":   12,
		"name": "那就开始拿大",
	})
	items2 = append(items2, map[string]interface{}{
		"id":   int64(123),
		"name": "2那就开始拿大",
	})
	_, err = SliceToMap(items2, "id")
	assert.Error(t, err)

	items3 := make([]map[string]interface{}, 0)
	items3 = append(items3, map[string]interface{}{
		"id":   int64(111),
		"name": "那就开始拿大",
	})
	items3 = append(items3, map[string]interface{}{
		"id":   int64(123),
		"name": "2那就开始拿大",
	})
	_, err = SliceToMap(items3, "id2")
	assert.Error(t, err)
}
