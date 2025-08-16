package longnet

import (
	"sync/atomic"

	"github.com/bwmarrin/snowflake"
)

type IdGenerator interface {
	// IdGen 获取自增ID
	IdGen() int64
}

type SnowflakeGenerator struct {
	snowflake *snowflake.Node
}

func NewSnowflakeGenerator() *SnowflakeGenerator {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	return &SnowflakeGenerator{
		snowflake: node,
	}
}

func (s *SnowflakeGenerator) IdGen() int64 {
	return s.snowflake.Generate().Int64()
}

type AutoIdGenerator struct {
	lastId int64
}

func NewAutoIdGenerator() *AutoIdGenerator {
	return &AutoIdGenerator{
		lastId: 0,
	}
}

func (a *AutoIdGenerator) IdGen() int64 {
	return atomic.AddInt64(&a.lastId, 1)
}
