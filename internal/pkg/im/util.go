package im

import (
	"sync"
)

// 获取 sync.Map 切片
func maps(num int) []*sync.Map {
	if num <= 0 {
		num = 1
	}

	items := make([]*sync.Map, 0, num)

	for i := 0; i < num; i++ {
		items = append(items, &sync.Map{})
	}

	return items
}
