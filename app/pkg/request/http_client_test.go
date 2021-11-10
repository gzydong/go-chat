package request

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func getHttpClient() *HttpClient {
	return NewHttpClient(&http.Client{})
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

func TestHttpClient_POST(t *testing.T) {
	client := getHttpClient()

	params := &url.Values{}
	params.Add("mobile", "18798276809")
	params.Add("password", "12dvsds")
	params.Add("platform", "web")

	client.SetDebug()

	resp, _ := client.Post("http://im-serve.gzydong.club/api/v1/auth/login", params)

	fmt.Println(string(resp))
}
