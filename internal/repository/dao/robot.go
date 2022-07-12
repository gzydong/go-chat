package dao

import (
	"go-chat/internal/repository/model"
)

type RobotDao struct {
	*BaseDao
}

func NewRobotDao(baseDao *BaseDao) *RobotDao {
	return &RobotDao{BaseDao: baseDao}
}

// FindLoginRobot 获取登录机器的信息
func (r *RobotDao) FindLoginRobot() (*model.Robot, error) {

	robot := &model.Robot{}

	err := r.db.Where("type = ? and status = ?", 1, 0).First(robot).Error
	if err != nil {
		return nil, err
	}

	return robot, nil
}
