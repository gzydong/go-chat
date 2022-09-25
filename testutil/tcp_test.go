package testutil

import (
	"fmt"
	"net"
	"testing"
	"time"

	"go-chat/internal/pkg/im/tcp"
)

func TestTcpServer_Setup(t1 *testing.T) {

	for i := 0; i < 1000; i++ {
		go conn()
	}

	time.Sleep(50 * time.Minute)
}

func conn() {
	conn, err := net.Dial("tcp", "106.14.177.175:9505")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()

	data, _ := tcp.Encode(`{"event":"authorize","content":"2056"}`)
	_, _ = conn.Write(data)

	go func() {
		for {
			msg := `{"event":"heartbeat","content":"ping"}`

			time.Sleep(5 * time.Second)
			data, _ := tcp.Encode(msg)
			_, _ = conn.Write(data)
			time.Sleep(10 * time.Second)
		}
	}()

	time.Sleep(1000 * time.Second)
}
