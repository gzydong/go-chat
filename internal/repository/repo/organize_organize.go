package repo

import (
	"context"

	"github.com/gzydong/go-chat/internal/pkg/core"
	"github.com/gzydong/go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Organize struct {
	core.Repo[model.Organize]
}

func NewOrganize(db *gorm.DB) *Organize {
	return &Organize{Repo: core.NewRepo[model.Organize](db)}
}

type UserInfo struct {
	UserId     int    `json:"user_id"`
	Nickname   string `json:"nickname"`
	Gender     int    `json:"gender"`
	Avatar     string `json:"avatar"`
	Department int    `json:"department"`
	Position   string `json:"position"`
}

func (o *Organize) List() ([]*UserInfo, error) {

	tx := o.Repo.Db.Table("organize")
	tx.Select([]string{
		"organize.user_id", "organize.dept_id as department", "organize.position_id as position",
		"users.nickname", "users.gender", "users.avatar",
	})
	tx.Joins("left join users on users.id = organize.user_id")

	items := make([]*UserInfo, 0)
	if err := tx.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// IsQiyeMember 判断是否是企业成员
func (o *Organize) IsQiyeMember(ctx context.Context, uid ...int) (bool, error) {

	count, err := o.Repo.FindCount(ctx, "user_id in ?", uid)
	if err != nil {
		return false, err
	}

	return int(count) == len(uid), nil
}

func (o *Organize) GetMemberIds(ctx context.Context) ([]int64, error) {
	var ids []int64
	if err := o.Repo.Db.WithContext(ctx).Table("organize").Pluck("user_id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

type GroupCount struct {
	DeptId int32 `gorm:"column:dept_id" json:"dept_id"`
	Count  int32 `gorm:"column:count" json:"count"`
}

func (o *Organize) DepartmentGroupCount(ctx context.Context) ([]*GroupCount, error) {
	var resp []*GroupCount
	err := o.Db.WithContext(ctx).Raw("select dept_id,count(*) as count from organize where dept_id group by dept_id").Scan(&resp).Error
	if err != nil {
		return nil, err
	}

	return resp, nil
}
