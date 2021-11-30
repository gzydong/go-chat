package testutil

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWsClient(t *testing.T) {
	for i := 0; i < 100; i++ {
		go NewClientTest(i)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(120000 * time.Second)
}

func NewClientTest(i int) {
	url := "ws://im-serve.local-admin.com/wss/default.io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImV4cCI6NTIzODI4MjE0MCwianRpIjoiNDEzNSJ9.BEpQL_YR4JM7qfpP4CskM043HQg9hvIH7dkQCK8VTfw" // 服务器地址
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := ws.WriteMessage(websocket.BinaryMessage, []byte(`{"event":"heartbeat","data":"ping"}`))
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second * 10)
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
