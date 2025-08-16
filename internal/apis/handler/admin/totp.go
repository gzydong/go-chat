package admin

import (
	"fmt"

	"github.com/pquerna/otp/totp"
	"github.com/samber/lo"
	"go-chat/api/pb/admin/v1"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/core/errorx"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/encrypt/rsautil"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm/clause"
)

type Totp struct {
	Rsa              rsautil.IRsa
	AdminRepo        *repo.Admin
	SysAdminTotpRepo *repo.SysAdminTotp
}

func (t *Totp) Status(ctx *core.Context) error {
	info, err := t.SysAdminTotpRepo.FindByWhere(ctx.GetContext(), "admin_id = ?", ctx.AuthId())
	if err != nil && !utils.IsSqlNoRows(err) {
		return ctx.Error(err)
	}

	if info == nil {
		return ctx.Success(admin.TotpStatusResponse{Enable: "N"})
	}

	return ctx.Success(admin.TotpStatusResponse{Enable: info.IsEnabled})
}

func (t *Totp) Close(ctx *core.Context) error {
	info, err := t.SysAdminTotpRepo.FindByWhere(ctx.GetContext(), "admin_id = ?", ctx.AuthId())
	if err != nil {
		return ctx.Error(err)
	}

	if info.IsEnabled == "N" {
		return ctx.Success(nil)
	}

	_, err = t.SysAdminTotpRepo.UpdateByWhere(ctx.GetContext(), map[string]any{
		"is_enabled": "N",
	}, "id = ?", info.Id)
	if err != nil {
		return err
	}

	return ctx.Success(nil)
}

func (t *Totp) Init(ctx *core.Context) error {
	var in admin.TotpInitRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	adminInfo, err := t.AdminRepo.FindById(ctx.GetContext(), ctx.AuthId())
	if err != nil {
		return err
	}

	password, err := t.Rsa.Decrypt(in.Password)
	if err != nil {
		return ctx.Error(err)
	}

	if !adminInfo.VerifyPassword(string(password)) {
		return ctx.InvalidParams("密码填写错误!")
	}

	resp, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "LumenIM-ADMIN",
		AccountName: adminInfo.Email,
	})

	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(admin.TotpInitResponse{
		QrcodeUri: resp.URL(),
		Secret:    resp.Secret(),
	})
}

func (t *Totp) Qrcode(ctx *core.Context) error {
	info, err := t.SysAdminTotpRepo.FindByWhere(ctx.GetContext(), "admin_id = ?", ctx.AuthId())
	if err != nil {
		return ctx.Error(err)
	}

	if info.IsEnabled == "N" {
		return errorx.New(400, "此账号未开启双因子认证")
	}

	adminInfo, err := t.AdminRepo.FindById(ctx.GetContext(), ctx.AuthId())
	if err != nil {
		return err
	}

	uri := "otpauth://totp/LumenIM-ADMIN:%s?algorithm=SHA1&digits=6&issuer=LumenIM-ADMIN&period=30&secret=%s"

	return ctx.Success(admin.TotpQrcodeResponse{
		QrcodeUri: fmt.Sprintf(uri, adminInfo.Email, info.Secret),
	})
}

func (t *Totp) Submit(ctx *core.Context) error {
	var in admin.TotpSubmitRequest
	if err := ctx.ShouldBindProto(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	adminInfo, err := t.AdminRepo.FindById(ctx.GetContext(), ctx.AuthId())
	if err != nil {
		return err
	}

	if !totp.Validate(in.Code, in.Session) {
		return errorx.New(400, "验证码填写错误")
	}

	codes := make([]string, 0)
	for i := 0; i < 12; i++ {
		codes = append(codes, lo.RandomString(6, lo.UpperCaseLettersCharset))
	}

	err = t.SysAdminTotpRepo.Db.WithContext(ctx.GetContext()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "admin_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_enabled", "secret", "one_time_code", "one_time_code", "updated_at"}),
	}).Create(&model.SysAdminTotp{
		AdminId: adminInfo.Id,
		Secret:  in.Session,
		OneTimeCode: lo.Map(codes, func(item string, _ int) string {
			return encrypt.HashPassword(item)
		}),
		IsEnabled: "Y",
	}).Error

	if err != nil {
		return err
	}

	return ctx.Success(admin.TotpSubmitResponse{
		OneTimeCode: codes,
	})
}
