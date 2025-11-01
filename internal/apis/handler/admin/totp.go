package admin

import (
	"context"
	"fmt"

	"github.com/pquerna/otp/totp"
	"github.com/samber/lo"
	"go-chat/api/pb/admin/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/errorx"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/encrypt/rsautil"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm/clause"
)

var _ admin.ITotpHandler = (*Totp)(nil)

type Totp struct {
	Rsa              rsautil.IRsa
	AdminRepo        *repo.Admin
	SysAdminTotpRepo *repo.SysAdminTotp
}

func (t *Totp) Status(ctx context.Context, in *admin.TotpStatusRequest) (*admin.TotpStatusResponse, error) {
	uid := middleware.FormContextAuthId[entity.AdminClaims](ctx)

	info, err := t.SysAdminTotpRepo.FindByWhere(ctx, "admin_id = ?", uid)
	if err != nil && !utils.IsSqlNoRows(err) {
		return nil, err
	}

	if info == nil {
		return &admin.TotpStatusResponse{Enable: "N"}, nil
	}

	return &admin.TotpStatusResponse{Enable: info.IsEnabled}, nil
}

func (t *Totp) Close(ctx context.Context, in *admin.TotpCloseRequest) (*admin.TotpCloseResponse, error) {
	uid := middleware.FormContextAuthId[entity.AdminClaims](ctx)

	info, err := t.SysAdminTotpRepo.FindByWhere(ctx, "admin_id = ?", uid)
	if err != nil {
		return nil, err
	}

	if info.IsEnabled == "N" {
		return &admin.TotpCloseResponse{}, nil
	}

	_, err = t.SysAdminTotpRepo.UpdateByWhere(ctx, map[string]any{
		"is_enabled": "N",
	}, "id = ?", info.Id)
	if err != nil {
		return nil, err
	}

	return &admin.TotpCloseResponse{}, nil
}

func (t *Totp) Init(ctx context.Context, in *admin.TotpInitRequest) (*admin.TotpInitResponse, error) {
	uid := middleware.FormContextAuthId[entity.AdminClaims](ctx)

	adminInfo, err := t.AdminRepo.FindById(ctx, uid)
	if err != nil {
		return nil, err
	}

	password, err := t.Rsa.Decrypt(in.Password)
	if err != nil {
		return nil, err
	}

	if !adminInfo.VerifyPassword(string(password)) {
		return nil, errorx.NewInvalidParams("密码填写错误")
	}

	resp, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "LumenIM-ADMIN",
		AccountName: adminInfo.Email,
	})

	if err != nil {
		return nil, err
	}

	return &admin.TotpInitResponse{
		QrcodeUri: resp.URL(),
		Secret:    resp.Secret(),
	}, nil
}

func (t *Totp) Qrcode(ctx context.Context, in *admin.TotpQrcodeRequest) (*admin.TotpQrcodeResponse, error) {
	uid := middleware.FormContextAuthId[entity.AdminClaims](ctx)

	info, err := t.SysAdminTotpRepo.FindByWhere(ctx, "admin_id = ?", uid)
	if err != nil {
		return nil, err
	}

	if info.IsEnabled == "N" {
		return nil, errorx.New(400, "此账号未开启双因子认证")
	}

	adminInfo, err := t.AdminRepo.FindById(ctx, uid)
	if err != nil {
		return nil, err
	}

	uri := "otpauth://totp/LumenIM-ADMIN:%s?algorithm=SHA1&digits=6&issuer=LumenIM-ADMIN&period=30&secret=%s"

	return &admin.TotpQrcodeResponse{
		QrcodeUri: fmt.Sprintf(uri, adminInfo.Email, info.Secret),
	}, nil
}

func (t *Totp) Submit(ctx context.Context, in *admin.TotpSubmitRequest) (*admin.TotpSubmitResponse, error) {
	uid := middleware.FormContextAuthId[entity.AdminClaims](ctx)

	adminInfo, err := t.AdminRepo.FindById(ctx, uid)
	if err != nil {
		return nil, err
	}

	if !totp.Validate(in.Code, in.Session) {
		return nil, errorx.New(400, "验证码填写错误")
	}

	codes := make([]string, 0)
	for i := 0; i < 12; i++ {
		codes = append(codes, lo.RandomString(6, lo.UpperCaseLettersCharset))
	}

	err = t.SysAdminTotpRepo.Db.WithContext(ctx).Clauses(clause.OnConflict{
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
		return nil, err
	}

	return &admin.TotpSubmitResponse{
		OneTimeCode: codes,
	}, nil
}
