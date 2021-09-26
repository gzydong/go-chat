package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go-chat/connect"
	"go-chat/testutil"
)

func TestUserRepository_FindByMobile(t *testing.T) {
	conf := testutil.GetConfig()
	db := connect.MysqlConnect(conf)
	userRep := UserRepository{DB: db}
	u, err := userRep.FindByMobile("123")
	assert.Error(t, err)
	assert.Nil(t, u)
	u, err = userRep.FindByMobile("18457673247")
	assert.NoError(t, err)
	assert.Equal(t, "18457673247", u.Nickname)
}
