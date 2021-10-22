package service

import (
	"fmt"
	"go-chat/testutil"
	"testing"
)

func newTestGroupMemberService() *GroupMemberService {

	return &GroupMemberService{db: testutil.GetDb()}
}

func TestGroupMemberService_isMember(t *testing.T) {
	service := newTestGroupMemberService()

	fmt.Println(service.isMember(87, 3039))
}

func TestGroupMemberService_GetMemberIds(t *testing.T) {
	service := newTestGroupMemberService()

	fmt.Println(service.GetUserGroupIds(2054))
}
