package dao

import "go-chat/internal/model"

type RobotDao struct {
	*BaseDao
}

// FindByLoginRobotInfo 获取登录机器的信息
func (r *RobotDao) FindByLoginRobotInfo() (*model.Robot, error) {

	return &model.Robot{}, nil
}
