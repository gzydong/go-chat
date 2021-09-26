package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	_, err := GeneratePassword([]byte("admin123"))
	assert.NoError(t, err)
}

func TestVerifyPassword(t *testing.T) {
	password := []byte("admin123")

	hash, err := GeneratePassword(password)
	assert.NoError(t, err)

	assert.Equal(t, true, VerifyPassword(password, hash))

	assert.Equal(t, false, VerifyPassword([]byte("testAdmin213"), hash))

	assert.Equal(t, false, VerifyPassword([]byte(""), hash))
}
