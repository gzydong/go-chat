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
	base := &Base{Db: provider.NewMySQLClient(testutil.GetConfig())}
	return base
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

	where := make(map[string]interface{})
	data := make(map[string]interface{})

	where["id IN ?"] = []int{1017, 1018}

	data["motto"] = "tttt"
	data["is_robot"] = 0

	_, err := base.Update(model.User{}, where, data)

	assert.NoError(t, err)
}
