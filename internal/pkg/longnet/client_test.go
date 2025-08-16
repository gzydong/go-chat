package longnet

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTcpClient(t *testing.T) {

	t.Skip()

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZXRhZGF0YSI6eyJ1c2VyX2lkIjoyMDU0fSwiaXNzIjoid2ViIiwiZXhwIjoxNzUxMDgyMTQ1LCJpYXQiOjE3NTA0ODIxNDUsImp0aSI6ImI2MTRmMzk0YzM3MjRmZjBhMmI2NmIwNzI5YWI5MWRmIn0.o-WJA60wefRbCWcmxJRq1Jso9g1dz9bfkfYAqchLKUw"

	client := NewTcpClient("127.0.0.1:9507")
	//client.SetCompress(NewSnappyCompress())
	client.SetOnMessage(func(c IClient, data []byte) {
		fmt.Println("收到的消息===>", string(data))
	})

	err := client.Connect(token)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("连接成功")
	}

	for i := 0; i < 10; i++ {
		_ = client.Write([]byte(fmt.Sprintf(`{"event":"message","body":{"text":"tcp 消息 %d"}}`, i)))
	}

	time.Sleep(3 * time.Second)
	client.Close()
}
