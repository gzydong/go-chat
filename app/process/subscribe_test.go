package process

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	m := map[string]interface{}{
		"group_id": 10,
		"uids":     []int{1, 2, 3},
	}

	if _, ok := m["group_id"].(int); !ok {
		return
	}

	fmt.Println(m["group_id"])
}
