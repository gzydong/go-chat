package socket

import (
	"github.com/bwmarrin/snowflake"
)

type IdGenerator interface {
	// IdGen 获取自增ID
	IdGen() int64
}

var defaultIdGenerator IdGenerator

func init() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	defaultIdGenerator = &SnowflakeGenerator{
		snowflake: node,
	}
}

type SnowflakeGenerator struct {
	snowflake *snowflake.Node
}

func (s *SnowflakeGenerator) IdGen() int64 {
	return s.snowflake.Generate().Int64()
}
