package organize

import (
	"go-chat/internal/dao"
	"go-chat/internal/model"
)

type IOrganizeDao interface {
	FindAll() ([]*UserInfo, error)
	IsQiyeMember(uid ...int) (bool, error)
}

type OrganizeDao struct {
	*dao.BaseDao
}

func NewOrganizeDao(baseDao *dao.BaseDao) *OrganizeDao {
	return &OrganizeDao{BaseDao: baseDao}
}

type UserInfo struct {
	UserId     int    `json:"user_id"`
	Nickname   string `json:"nickname"`
	Gender     int    `json:"gender"`
	Department string `json:"department"`
	Position   string `json:"position"`
}

func (o *OrganizeDao) FindAll() ([]*UserInfo, error) {

	tx := o.Db().Table("organize")
	tx.Select([]string{
		"organize.user_id", "organize.department", "organize.position",
		"users.nickname", "users.gender",
	})
	tx.Joins("left join users on users.id = organize.user_id")

	items := make([]*UserInfo, 0)
	if err := tx.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// IsQiyeMember 判断是否是企业成员
func (o *OrganizeDao) IsQiyeMember(uid ...int) (bool, error) {

	var count int64
	err := o.Db().Model(model.Organize{}).Where("user_id in ?", uid).Count(&count).Error
	if err != nil {
		return false, err
	}

	return int(count) == len(uid), nil
}
