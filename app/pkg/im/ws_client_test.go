package im

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestWsClient(t *testing.T) {
	for i := 0; i < 500; i++ {
		go NewClientTest(i)
		time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(120000 * time.Second)
}

func NewClientTest(i int) {
	url := "ws://127.0.0.1:8080/wss/socket.io?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsInVzZXJfaWQiOjIwNTQsImV4cCI6MTYzNDM4NzQwOCwiaXNzIjoiZ28tY2hhdCJ9.dMooNhH1K-G3ihFV8_pHzdeHxD4c9Q72BdfSFEFBMTA" //服务器地址
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
