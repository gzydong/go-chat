package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	str, err := GeneratePassword([]byte("admin123"))
	t.Logf("%s\n", str)
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
