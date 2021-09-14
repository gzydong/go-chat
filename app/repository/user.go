package repository

import "go-chat/app/model"

type UserRepository struct {
}

// 获取实例
func GetUserInstance() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) First(user *model.User) {

}
