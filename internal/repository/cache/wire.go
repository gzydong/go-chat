package cache

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewCaptchaStorage,
	NewClientStorage,
	NewContactRemark,
	NewRedisLock,
	NewMessageStorage,
	NewRelation,
	NewRoomStorage,
	NewSequence,
	NewTokenSessionStorage,
	NewSidStorage,
	NewSmsStorage,
	NewVote,
	NewUnreadStorage,
	NewGroupApplyStorage,
)
