package note

import "go-chat/app/service"

type ArticleService struct {
	*service.BaseService
}

func NewArticleService(baseService *service.BaseService) *ArticleService {
	return &ArticleService{BaseService: baseService}
}
