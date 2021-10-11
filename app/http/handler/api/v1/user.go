package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/entity"
	"go-chat/app/helper"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/model"
	"go-chat/app/repository"
	"go-chat/app/service"
)

type User struct {
	UserRepo   *repository.UserRepository
	SmsService *service.SmsService
}

// Detail 个人用户信息
func (u *User) Detail(c *gin.Context) {
	user, _ := u.UserRepo.FindById(c.GetInt("__user_id__"))

	response.Success(c, gin.H{
		"detail": user,
	})
}

// ChangeDetail 修改个人用户信息
func (u *User) ChangeDetail(c *gin.Context) {
	params := &request.ChangeDetailRequest{}
	if err := c.Bind(params); err != nil {
		response.InvalidParams(c, err)
		return
	}

	_, _ = u.UserRepo.Update(&model.User{ID: c.GetInt("__user_id__")}, map[string]interface{}{
		"nickname": params.Nickname,
		"avatar":   params.Avatar,
		"gender":   params.Gender,
		"motto":    params.Profile,
	})

	response.Success(c, gin.H{}, "个人信息修改成功！")
}

// ChangePassword 修改密码接口
func (u *User) ChangePassword(c *gin.Context) {
	params := &request.ChangePasswordRequest{}
	if err := c.Bind(params); err != nil {
		response.InvalidParams(c, err)
		return
	}

	user, _ := u.UserRepo.FindById(c.GetInt("__user_id__"))
	if !helper.VerifyPassword([]byte(params.OldPassword), []byte(user.Password)) {
		response.BusinessError(c, "密码填写错误！")
		return
	}

	// 生成 hash 密码
	hash, _ := helper.GeneratePassword([]byte(params.NewPassword))

	_, err := u.UserRepo.Update(&model.User{ID: user.ID}, map[string]interface{}{
		"password": hash,
	})

	if err != nil {
		response.BusinessError(c, "密码修改失败！")
		return
	}

	response.Success(c, gin.H{}, "密码修改成功！")
}

// ChangeMobile 修改手机号接口
func (u *User) ChangeMobile(c *gin.Context) {
	params := &request.ChangeMobileRequest{}
	if err := c.Bind(params); err != nil {
		response.InvalidParams(c, err)
		return
	}

	if !u.SmsService.CheckSmsCode(c.Request.Context(), entity.SmsChangeAccountChannel, params.Mobile, params.SmsCode) {
		response.BusinessError(c, "短信验证码填写错误！")
		return
	}

	user, _ := u.UserRepo.FindById(c.GetInt("__user_id__"))

	if user.Mobile != params.Mobile {
		response.BusinessError(c, "手机号与原手机号一致无需修改！")
		return
	}

	if !helper.VerifyPassword([]byte(params.Password), []byte(user.Password)) {
		response.BusinessError(c, "账号密码填写错误！")
		return
	}

	_, err := u.UserRepo.Update(&model.User{ID: user.ID}, map[string]interface{}{
		"mobile": params.Mobile,
	})

	if err != nil {
		response.BusinessError(c, "手机号修改失败！")
		return
	}

	response.Success(c, gin.H{}, "手机号修改成功！")
}

// ChangeEmail 修改邮箱接口
func (u *User) ChangeEmail(c *gin.Context) {
	params := &request.ChangeEmailRequest{}
	if err := c.Bind(params); err != nil {
		response.InvalidParams(c, err)
		return
	}

	// todo 1.验证邮件激活码是否正确
}
