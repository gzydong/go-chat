package service

import (
	"go-chat/internal/pkg/utils"
	"go-chat/resource"
)

var _ ITemplateService = (*TemplateService)(nil)

type ITemplateService interface {
	CodeTemplate(data map[string]string) (string, error)
}

type TemplateService struct {
}

// CodeTemplate 验证码通知模板
func (t *TemplateService) CodeTemplate(data map[string]string) (string, error) {

	fileContent, err := resource.Template().ReadFile("templates/email/verify_code.tmpl")
	if err != nil {
		return "", err
	}

	return utils.RenderTemplate(fileContent, data)
}
