package note

import (
	"go-chat/internal/repository/dao"
)

type ArticleDao struct {
	*dao.BaseDao
}

func NewArticleDao(baseDao *dao.BaseDao) *ArticleDao {
	return &ArticleDao{BaseDao: baseDao}
}
