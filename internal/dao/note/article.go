package note

import "go-chat/internal/dao"

type ArticleDao struct {
	*dao.BaseDao
}

func NewArticleDao(baseDao *dao.BaseDao) *ArticleDao {
	return &ArticleDao{BaseDao: baseDao}
}
