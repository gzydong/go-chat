package testutil

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
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

func TestName(t *testing.T) {

	for i := 0; i < 100; i++ {

		index := i
		go func() {
			client := &http.Client{}
			var data = strings.NewReader(fmt.Sprintf(`talk_type=1&receiver_id=2055&text=%d那几款撒那看你哪款手机那`, index))
			req, err := http.NewRequest("POST", "http://127.0.0.1:9503/api/v1/talk/message/text", data)
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Set("User-Agent", "Apipost client Runtime/+https://www.apipost.cn/")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImlzcyI6Imx1bWVuLWltIiwiZXhwIjoxNzA0MjI2MzQxLCJpYXQiOjE2NjgyMjYzNDEsImp0aSI6IjIwNTQifQ.m9a6zPy7dBSsVyUsR6IKQVRUAzb9R1vW-R2Jb9yOcT8")
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			bodyText, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", bodyText)
		}()
	}

	time.Sleep(10 * time.Second)

}
