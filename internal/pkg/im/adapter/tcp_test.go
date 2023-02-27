package adapter

import (
	"fmt"
	"net"
	"testing"
	"time"

	"go-chat/internal/pkg/im/adapter/encoding"
	"go-chat/internal/pkg/logger"
)

func TestTcp_Server(t *testing.T) {
	listener, _ := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 9501))

	defer func() {
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}

		go func() {
			conn, err := NewTcpAdapter(conn)
			if err != nil {
				logger.Errorf("tcp connect error: %s", err.Error())
			}

			for {
				data, err := conn.Read()
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(string(data))
			}
		}()
	}
}

func TestTcp_Client(t1 *testing.T) {
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:9501")
		if err != nil {
			fmt.Println("dial failed, err", err)
			return
		}

		defer conn.Close()

		go func() {
			index := 0
			for {
				msg := fmt.Sprintf(`{"event":"heartbeat%d","content":"ping"}`, index)

				data, _ := encoding.NewEncode([]byte(msg))
				_, _ = conn.Write(data)
				time.Sleep(2 * time.Second)
				index++
			}
		}()

		time.Sleep(1 * time.Hour)
	}()

	time.Sleep(50 * time.Minute)
}
