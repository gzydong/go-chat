package utils

import "github.com/bwmarrin/snowflake"

var snowflakeNode, _ = snowflake.NewNode(1)

// GenSnowflakeId 雪花算法生成ID
func GenSnowflakeId() int64 {
	return snowflakeNode.Generate().Int64()
}
