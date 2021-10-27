package dao

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-chat/app/model"
	"go-chat/provider"
	"go-chat/testutil"
	"testing"
)

func newBase() *Base {
	base := &Base{db: provider.MysqlConnect(testutil.GetConfig())}
	return base
}

func TestBase_Create(t *testing.T) {
	base := newBase()

	var user model.User

	user.Nickname = "test"
	user.Mobile = "18721312319"
	user.CreatedAt = "2020-10-24 00:00:00"
	user.UpdatedAt = "2020-10-24 00:00:00"

	err := base.Create(&user)

	fmt.Printf("%#v\n", user)
	fmt.Printf("error :%s\n", err)
}

func TestBase_FindByIds(t *testing.T) {
	base := newBase()

	var items []model.User

	ok, err := base.FindByIds(&items, []int{2054, 231231}, "*")

	fmt.Printf("%#v\n", items)
	fmt.Println(ok)
	fmt.Printf("length :%d", len(items))
	fmt.Printf("error :%s\n", err)
}

func TestBase_Update(t *testing.T) {
	base := newBase()

	where := make(map[string]interface{}, 0)
	data := make(map[string]interface{}, 0)

	where["id IN ?"] = []int{1017, 1018}

	data["motto"] = "tttt"
	data["is_robot"] = 0

	_, err := base.Update(model.User{}, where, data)

	assert.NoError(t, err)
}
