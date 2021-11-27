package service

import (
	"fmt"
	"go-chat/testutil"
	"testing"
)

func TestName(t *testing.T) {

	s := NewGroupMemberService(testutil.GetDb())

	arr := s.GetUserGroupIds(2054)

	fmt.Printf("%T", arr)
}
