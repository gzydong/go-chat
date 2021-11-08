package testutil

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWsClient(t *testing.T) {
	for i := 0; i < 1000; i++ {
		go NewClientTest(i)
		time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(120000 * time.Second)
}

func NewClientTest(i int) {
	url := "ws://106.14.177.175:8080/wss/default.io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImV4cCI6NTIzNjE5ODkwNywianRpIjoiMjA1NCJ9.vl5fibjsEXs1U56ZGdLW7aFlsQ6Fm76hNx6CAl7bzeM" // 服务器地址
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := ws.WriteMessage(websocket.BinaryMessage, []byte("ping"))
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second * 1)

			_ = ws.WriteMessage(websocket.BinaryMessage, []byte(strconv.Itoa(i)))

			time.Sleep(time.Second * 5)
		}
	}()

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("receive: ", string(data))
	}
}
