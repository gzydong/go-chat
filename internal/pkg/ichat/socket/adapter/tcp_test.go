package adapter

import (
	"fmt"
	"net"
	"testing"
	"time"

	"go-chat/internal/pkg/ichat/socket/adapter/encoding"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/strutil"
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

type Authorize struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
}

func TestTcp_Client(t1 *testing.T) {
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:9505")
		if err != nil {
			fmt.Println("dial failed, err", err)
			return
		}

		defer conn.Close()

		data, _ := encoding.NewEncode(jsonutil.Marshal(Authorize{
			Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTcxMzY3MzYwOSwiaWF0IjoxNjc3NjczNjA5LCJqdGkiOiIyMDU0In0.HzJjDgq_e3yUORRiHLXRoC16aWjZdbc7MJNwFlWWstA",
			Channel: "chat",
		}))
		_, _ = conn.Write(data)

		go func() {
			for index := 0; index < 100000; index++ {
				msg := fmt.Sprintf(`{"msg_id":"%s","event":"event.talk.text.message","body":{"receiver":{"talk_type":1,"receiver_id":2055},"type":1,"content":"测阿珂神经%d内科","mention":{"all":0,"uids":[]}}}`, strutil.NewMsgId(), index)

				data, _ := encoding.NewEncode([]byte(msg))
				_, err := conn.Write(data)

				if err != nil {
					fmt.Println(err)
					break
				} else {
					fmt.Println("ok")
				}

				time.Sleep(10 * time.Millisecond)
				// index++
			}
		}()

		go func() {

			for {
				data, err := encoding.NewDecode(conn)

				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(string(data))
			}
		}()

		time.Sleep(1 * time.Hour)
	}()

	time.Sleep(50 * time.Minute)
}

// nolint
func TestTcp_Client2(t1 *testing.T) {
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:9505")
		if err != nil {
			fmt.Println("dial failed, err", err)
			return
		}

		defer conn.Close()

		data, _ := encoding.NewEncode(jsonutil.Marshal(Authorize{
			Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6ImltLndlYiIsImV4cCI6MTcxMzUwNDY2MywiaWF0IjoxNjc3NTA0NjYzLCJqdGkiOiIyMDU0In0.qMybbKgqYlDR-5mP7VnoDIV8ex9Hg_tsXu8cTekX7-c",
			Channel: "chat",
		}))
		_, _ = conn.Write(data)

		go func() {
			for {
				msg := fmt.Sprintf(`{"msg_id":"%s","event":"event.talk.text.message","body":{"receiver":{"talk_type":1,"receiver_id":2055},"type":1,"content":"测阿珂神经%d内科","mention":{"all":0,"uids":[]}}}`, strutil.NewMsgId(), 999999)

				data, _ := encoding.NewEncode([]byte(msg))
				conn.Write(data)

				time.Sleep(20 * time.Second)
			}
		}()

		go func() {

			for {
				data, err := encoding.NewDecode(conn)

				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(string(data))
			}
		}()

		time.Sleep(1 * time.Hour)
	}()

	time.Sleep(50 * time.Minute)
}
