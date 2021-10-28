package testutil

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestWsClient(t *testing.T) {
	for i := 0; i < 100; i++ {
		go NewClientTest(i)
		time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(120000 * time.Second)
}

func NewClientTest(i int) {
	url := "ws://127.0.0.1:8080/wss/socket.io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImV4cCI6MjQ5ODY1MzMzNCwianRpIjoiMjA1NCJ9.F8tKnShU6IXpN9-OzJJa6ZI7f29z1KkcqmUKiJ55MIc" //服务器地址
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
