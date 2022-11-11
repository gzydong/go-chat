package repo

import (
	"go-chat/internal/repository/model"
)

type Robot struct {
	*Base
}

func NewRobot(baseDao *Base) *Robot {
	return &Robot{Base: baseDao}
}

// FindLoginRobot 获取登录机器的信息
func (repo *Robot) FindLoginRobot() (*model.Robot, error) {

	robot := &model.Robot{}

	err := repo.db.Where("type = ? and status = ?", 1, 0).First(robot).Error
	if err != nil {
		return nil, err
	}

	return robot, nil
}
