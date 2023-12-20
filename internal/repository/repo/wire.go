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
	NewTalkRecords,
	NewGroupNotice,
	NewTalkSession,
	NewTalkRecordsDelete,
	NewEmoticon,
	NewTalkRecordsVote,
	NewFileSplitUpload,
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
