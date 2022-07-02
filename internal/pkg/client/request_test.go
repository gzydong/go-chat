package client

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func getHttpClient() *RequestClient {
	return NewRequestClient(&http.Client{})
}

func TestHttpClient_GET(t *testing.T) {
	client := getHttpClient()

	params := &url.Values{}
	params.Add("key", "836a9c3ea9bba4e0a4d51bd02fbcc5")
	params.Add("ip", "49.211.168.144")

	client.SetDebug()
	resp, _ := client.Get("http://apis.juhe.cn/ip/ipNew", params)
	fmt.Println(string(resp))
}

// func TestHttpClient_POST(t *testing.T) {
// 	client := getHttpClient()
//
// 	for i := 0; i < 1000; i++ {
// 		go func() {
// 			params := &url.Values{}
// 			params.add("talk_type", "1")
// 			params.add("receiver_id", "2055")
// 			params.add("text", " 那几款撒那看你哪款手机那")
// 			params.add("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJndWFyZCI6ImFwaSIsImV4cCI6NTIzODI4MzQ4MywianRpIjoiMjA1NCJ9.td6MPf_G64XTUWz-wSIb78fBiarKKJA0P3jFjGgde8o")
//
// 			client.SetDebug()
//
// 			resp, _ := client.Post("http://127.0.0.1:9503/api/v1/talk/message/text", params)
//
// 			fmt.Println(string(resp))
// 		}()
// 	}
//
// 	time.Sleep(time.Hour * 1)
//
// }
