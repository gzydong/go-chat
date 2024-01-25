package repo

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewSource,
	NewContact,
	NewContactGroup,
	NewGroupMember,
	NewUsers,
	NewGroup,
	NewGroupApply,
	NewTalkRecordGroup,
	NewTalkRecordFriend,
	NewGroupNotice,
	NewTalkSession,
	NewTalkRecordGroupDel,
	NewEmoticon,
	NewGroupVote,
	NewFileUpload,
	NewArticleClass,
	NewArticle,
	NewArticleAnnex,
	NewDepartment,
	NewOrganize,
	NewPosition,
	NewRobot,
	NewSequence,
	NewAdmin,
)
