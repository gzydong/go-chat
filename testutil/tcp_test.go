package testutil

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"

	"go-chat/internal/pkg/im/tcp"
)

type TcpAuthRequest struct {
	Channel string // 连接业务渠道
	Token   string // 授权token
}

func TestTcpServer_Setup(t1 *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:9505")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()

	go func() {
		msg := `{"event":"login","content":"2056"}`
		data, _ := tcp.Encode(msg)
		_, _ = conn.Write(data)

		for {
			msg := `{"event":"heartbeat","content":"ping"}`

			time.Sleep(5 * time.Second)
			data, _ := tcp.Encode(msg)
			_, _ = conn.Write(data)
			time.Sleep(10 * time.Second)
		}
	}()

	go func() {
		ch := bufio.NewReader(conn)
		for {
			decode, err := tcp.Decode(ch)
			if err != nil {
				return
			}

			fmt.Println(decode)
		}
	}()

	time.Sleep(50 * time.Minute)
}
