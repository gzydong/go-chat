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
	NewArticleHistory,
	NewArticleAnnex,
	NewDepartment,
	NewOrganize,
	NewPosition,
	NewRobot,
	NewSequence,
	NewAdmin,
	NewOAuthUsers,
	NewSysRole,
	NewSysResource,
	NewSysMenu,
	NewSysAdminTotp,
)
