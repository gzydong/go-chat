package testutil

import (
	"fmt"
	"net"
	"testing"
	"time"

	"go-chat/internal/pkg/im/adapter/encoding"
)

func TestTcpServer_Setup(t1 *testing.T) {

	for i := 0; i < 1000; i++ {
		go conn()
	}

	time.Sleep(50 * time.Minute)
}

func conn() {
	conn, err := net.Dial("tcp", "127.0.0.1:9505")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()

	data, _ := encoding.Encode(`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImV4cCI6MTcwMDExNTk0MiwianRpIjoiMjA1NCJ9.2hg4nkwDMflJJs4kqWuNXiizBdGgkbrTNiIM8l84d6E","channel":"chat"}`)
	_, _ = conn.Write(data)

	go func() {
		for {
			msg := `{"event":"heartbeat","content":"ping"}`

			time.Sleep(5 * time.Second)
			data, _ := encoding.Encode(msg)
			_, _ = conn.Write(data)
			time.Sleep(10 * time.Second)
		}
	}()

	time.Sleep(1000 * time.Second)
}
