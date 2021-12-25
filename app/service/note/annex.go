package note

import "go-chat/app/service"

type ArticleAnnexService struct {
	*service.BaseService
}

func NewArticleAnnexService(baseService *service.BaseService) *ArticleAnnexService {
	return &ArticleAnnexService{BaseService: baseService}
}
