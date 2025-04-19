package repo

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type SequenceType int32

const (
	SequenceTypeUser  = 1
	SequenceTypeGroup = 2
)

type Sequence struct {
	db    *gorm.DB
	cache *cache.Sequence
	core.Repo[model.Sequence]
	snowflake *snowflake.Node
}

func NewSequence(db *gorm.DB, cache *cache.Sequence) *Sequence {

	res := cache.Redis().Incr(context.Background(), "snowflake_work_node")

	node, err := snowflake.NewNode(res.Val() % 100)
	if err != nil {
		panic(err)
	}

	return &Sequence{db: db, cache: cache, Repo: core.NewRepo[model.Sequence](db), snowflake: node}
}

// Get 获取会话间的时序ID
func (s *Sequence) Get(ctx context.Context, seqType SequenceType, sourceId int32) int64 {
	return s.snowflake.Generate().Int64()
}

// BatchGet 批量获取会话间的时序ID
func (s *Sequence) BatchGet(ctx context.Context, seqType SequenceType, sourceId int32, num int) []int64 {
	ids := make([]int64, 0)

	for i := 0; i < num; i++ {
		ids = append(ids, s.snowflake.Generate().Int64())
	}

	return ids
}
