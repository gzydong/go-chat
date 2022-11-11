package organize

import (
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type IOrganizeDao interface {
	FindAll() ([]*UserInfo, error)
	IsQiyeMember(uid ...int) (bool, error)
}

type Organize struct {
	*repo.Base
}

func NewOrganize(baseDao *repo.Base) *Organize {
	return &Organize{Base: baseDao}
}

type UserInfo struct {
	UserId     int    `json:"user_id"`
	Nickname   string `json:"nickname"`
	Gender     int    `json:"gender"`
	Department string `json:"department"`
	Position   string `json:"position"`
}

func (repo *Organize) FindAll() ([]*UserInfo, error) {

	tx := repo.Db().Table("organize")
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
func (repo *Organize) IsQiyeMember(uid ...int) (bool, error) {

	var count int64
	err := repo.Db().Model(model.Organize{}).Where("user_id in ?", uid).Count(&count).Error
	if err != nil {
		return false, err
	}

	return int(count) == len(uid), nil
}
