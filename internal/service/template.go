package service

import (
	"go-chat/internal/pkg/utils"
	"go-chat/resource"
)

type TemplateService struct {
}

func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

// CodeTemplate 验证码通知模板
func (t *TemplateService) CodeTemplate(data map[string]string) (string, error) {

	fileContent, err := resource.Template().ReadFile("templates/email/verify_code.tmpl")
	if err != nil {
		return "", err
	}

	return utils.RenderTemplate(fileContent, data)
}
