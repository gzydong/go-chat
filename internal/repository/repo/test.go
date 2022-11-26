package repo

import (
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Test struct {
	ichat.Base[model.Article]
}

func NewTest(db *gorm.DB) *Test {
	return &Test{Base: ichat.Base[model.Article]{Db: db}}
}

func (t *Test) name() {}
