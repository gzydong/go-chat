package sliceutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice_UniqueString(t *testing.T) {
	items := []string{"11", "22", "33", "22", "55", "11", "66"}

	assert.Equal(t, 5, len(UniqueString(items)))
}

func TestSlice_UniqueInt(t *testing.T) {
	items := []int{1, 2, 5, 3, 6, 1, 8, 99, 22, 3}

	assert.Equal(t, 8, len(UniqueInt(items)))
}

func TestSlice_UniqueInt64(t *testing.T) {
	items := []int64{1, 2, 5, 3, 6, 1, 8, 99, 22, 3}

	assert.Equal(t, 8, len(UniqueInt64(items)))
}

func TestSlice_ParseIds(t *testing.T) {
	assert.Equal(t, 0, len(ParseIds("")))
	assert.Equal(t, 0, len(ParseIds(" ")))
	assert.Equal(t, 3, len(ParseIds("1,2,3")))
	assert.Equal(t, 3, len(ParseIds("3,3,3")))
}
