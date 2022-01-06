package encrypt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBase64Decode(t *testing.T) {

}

func TestBase64Encode(t *testing.T) {

}

func TestHashPassword(t *testing.T) {
	pwd, err := HashPassword("admin123sdgsdsgds")

	assert.NoError(t, err)
	assert.Equal(t, true, len(pwd) == 60)
}

func TestMd5(t *testing.T) {
	assert.Equal(t, "c069d1fbc7bf8e994de8299110e68bc5", Md5("s6hqzp6j0kdfzh4n_cjq6b180000gn"))
}

func TestVerifyPassword(t *testing.T) {
	pwd, err := HashPassword("admin123")

	assert.NoError(t, err)

	assert.Equal(t, true, VerifyPassword(pwd, "admin123"))
	assert.Equal(t, false, VerifyPassword(pwd, "admin1234453"))
}
