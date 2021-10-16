package repository

import (
	"fmt"
	"go-chat/app/model"
	"go-chat/connect"
	"go-chat/testutil"
	"testing"
)

func newBase() *Base {
	base := &Base{db: connect.MysqlConnect(testutil.GetConfig())}
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
