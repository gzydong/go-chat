package repo

import (
	"github.com/google/wire"
	note3 "go-chat/internal/repository/repo/note"
	organize3 "go-chat/internal/repository/repo/organize"
)

var ProviderSet = wire.NewSet(
	NewSource,
	NewContact,
	NewContactGroup,
	NewGroupMember,
	NewUsers,
	NewGroup,
	NewGroupApply,
	NewTalkRecords,
	NewGroupNotice,
	NewTalkSession,
	NewEmoticon,
	NewTalkRecordsVote,
	NewFileSplitUpload,
	note3.NewArticleClass,
	note3.NewArticleAnnex,
	organize3.NewDepartment,
	organize3.NewOrganize,
	organize3.NewPosition,
	NewRobot,
	NewSequence,
	NewAdmin,
)
