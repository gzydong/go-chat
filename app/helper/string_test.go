package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenValidateCode(t *testing.T) {
	assert.Equal(t, 0, len(GenValidateCode(0)))
	assert.Equal(t, 6, len(GenValidateCode(6)))
	assert.Equal(t, 10, len(GenValidateCode(10)))
}

func TestGetRandomString(t *testing.T) {
	assert.Equal(t, 0, len(GetRandomString(0)))
	assert.Equal(t, 6, len(GetRandomString(6)))
	assert.Equal(t, 10, len(GetRandomString(10)))
}
