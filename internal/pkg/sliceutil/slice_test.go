package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice_Unique(t *testing.T) {
	items := []int{1, 2, 5, 3, 6, 1, 8, 99, 22, 3}

	assert.Equal(t, 8, len(Unique(items)))
}

func TestSlice_ParseIds(t *testing.T) {
	assert.Equal(t, 0, len(ParseIds("")))
	assert.Equal(t, 0, len(ParseIds(" ")))
	assert.Equal(t, 3, len(ParseIds("1,2,3")))
	assert.Equal(t, 3, len(ParseIds("3,3,3")))
}

func TestSum(t *testing.T) {
	assert.Equal(t, 28, Sum([]int{1, 2, 3, 4, 5, 6, 7}))
}

func TestToMap(t *testing.T) {
	type Data struct {
		key   string
		value any
	}

	items := make([]*Data, 0)
	items = append(items, &Data{key: "111", value: 111})
	items = append(items, &Data{key: "222", value: 222})
	items = append(items, &Data{key: "333", value: 333})
	items = append(items, &Data{key: "111", value: 444})

	maps := ToMap(items, func(t *Data) string {
		return t.key
	})

	assert.Equal(t, 3, len(maps))
}
