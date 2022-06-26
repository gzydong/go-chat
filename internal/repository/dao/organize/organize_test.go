package organize

import (
	"fmt"
	"testing"

	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/repository/dao"
	"go-chat/testutil"
)

func newTestOrganizeDao(t *testing.T) *OrganizeDao {
	return NewOrganizeDao(dao.NewBaseDao(testutil.GetDb(), testutil.TestRedisClient()))
}

func TestOrganizeDao_FindAll(t *testing.T) {
	all, _ := newTestOrganizeDao(t).FindAll()

	t.Log(jsonutil.Encode(all))
}

func TestOrganizeDao_IsQiyeMember(t *testing.T) {
	isTrue, _ := newTestOrganizeDao(t).IsQiyeMember(2054, 2055, 2)

	fmt.Println(isTrue)
}
