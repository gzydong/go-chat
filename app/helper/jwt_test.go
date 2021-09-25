package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go-chat/testutil"
)

func TestGenerateJwtToken(t *testing.T) {
	conf := testutil.GetConfig()
	data, err := GenerateJwtToken(conf, "user", 1)
	assert.NoError(t, err)

	user, err1 := ParseJwtToken(conf, data["token"].(string))
	assert.NoError(t, err1)
	assert.Equal(t, 1, user.UserId)
}
