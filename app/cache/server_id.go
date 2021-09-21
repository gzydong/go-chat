package cache

import (
	"strconv"
	"time"
)

type ServerRunID struct {
}

const (
	ServerRunIdKey = "server_ids"

	// 运行检测超时时间（单位秒）
	ServerOverTime = 35
)

func NewServerRun() *ServerRunID {
	return new(ServerRunID)
}

func (s *ServerRunID) SetServerRunId(server string, time int64) {
	Rdb.HSet(ServerRunIdKey, server, time)
}

// GetServerRunIdAll 获取指定状态的运行ID
// status 状态[1:运行中;2:已超时;3:全部]
func (s ServerRunID) GetServerRunIdAll(status int) []string {
	result, err := Rdb.HGetAll(ServerRunIdKey).Result()

	slice := make([]string, 0)

	t := time.Now().Unix()
	if err == nil {
		for key, val := range result {
			value, err := strconv.Atoi(val)

			if err != nil {
				continue
			}

			switch status {
			case 1:
				if t-int64(value) >= ServerOverTime {
					continue
				}
				break
			case 2:
				if t-int64(value) < ServerOverTime {
					continue
				}

				break
			}

			slice = append(slice, key)
		}
	}

	return slice
}
