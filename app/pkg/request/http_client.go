package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HttpClient struct {
	debug  bool
	client *http.Client
}

type FileData struct {
	Field    string // 字段名
	FileName string // 字段名
	Content  []byte // 文件字节流
}

func NewHttpClient(client *http.Client) *HttpClient {
	return &HttpClient{client: client}
}

func (c *HttpClient) SetDebug() {
	c.debug = true
}

func (c *HttpClient) Get(url string, params *url.Values) ([]byte, error) {
	if params != nil {
		if strings.Contains(url, "?") {
			url = fmt.Sprintf("%s&%s", url, params.Encode())
		} else {
			url = fmt.Sprintf("%s?%s", url, params.Encode())
		}
	}

	req, _ := http.NewRequest("GET", url, nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if c.debug {
		fmt.Printf("\n[GET] HTTP Request\n")
		fmt.Printf("Request URL : %s\n", url)
		fmt.Printf("Response StatusCode: %d\n", resp.StatusCode)
		fmt.Printf("Response Data: %s\n\n", string(res))
	}

	return res, nil
}

func (c *HttpClient) Post(url string, params *url.Values) ([]byte, error) {
	req, _ := http.NewRequest("POST", url, strings.NewReader(params.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if c.debug {
		fmt.Printf("\n[POST] HTTP Request\n")
		fmt.Printf("Request URL : %s\n", url)
		fmt.Printf("Request Data: %s\n", params.Encode())
		fmt.Printf("Response StatusCode: %d\n", resp.StatusCode)
		fmt.Printf("Response Data: %s\n\n", string(res))
	}

	return res, nil
}

func (c *HttpClient) PostFrom(url string, params *url.Values, files []*FileData) {

}
