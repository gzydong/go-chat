package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

const (
	JsonContentType = "application/json; charset=utf-8"
	FormContentType = "application/x-www-form-urlencoded; charset=utf-8"
	XmlContentType  = "application/xml; charset=utf-8"
)

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

// nolint
func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

type TestRequest struct {
	*http.Request
	path       string
	handlers   []gin.HandlerFunc
	beforeHook func(req *http.Request)
}

func NewTestRequest(url string, handlers ...gin.HandlerFunc) *TestRequest {
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	return &TestRequest{
		path:     url,
		Request:  req,
		handlers: handlers,
	}
}

func NewTestRequestWithPath(path, url string, handlers ...gin.HandlerFunc) *TestRequest {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return &TestRequest{
		path:     path,
		Request:  req,
		handlers: handlers,
	}
}

func (r *TestRequest) Form(values *url.Values) (*TestResponse, error) {
	return r.Post(FormContentType, strings.NewReader(values.Encode()))
}

func (r *TestRequest) JSON(data interface{}) (*TestResponse, error) {
	var body []byte
	switch data := data.(type) {
	case []byte:
		body = data
	case string:
		body = []byte(data)
	default:
		body, _ = json.Marshal(data)
	}

	return r.Post(JsonContentType, bytes.NewReader(body))
}

func (r *TestRequest) XML(body string) (*TestResponse, error) {
	return r.Post(XmlContentType, strings.NewReader(body))
}

func (r *TestRequest) Get() (*TestResponse, error) {
	return r.doRequest(http.MethodGet, "", nil)
}

func (r *TestRequest) Post(contentType string, body io.Reader) (*TestResponse, error) {
	return r.doRequest(http.MethodPost, contentType, body)
}

func (r *TestRequest) BeforeHook(fn func(req *http.Request)) {
	r.beforeHook = fn
}

func (r *TestRequest) doRequest(method string, contentType string, body io.Reader) (*TestResponse, error) {
	var err error
	var bb []byte

	req := r.Request.Clone(context.Background())
	req.Method = method
	if body != nil {
		bb, err = ioutil.ReadAll(body)
		if err != nil {
			return nil, err
		}

		req.Body = ioutil.NopCloser(bytes.NewReader(bb))
	}

	if method == http.MethodPost {
		req.Header.Set("Content-Type", contentType)
	}

	r.Request = req

	// send before hook
	if r.beforeHook != nil {
		r.beforeHook(r.Request)
	}

	w := newCloseNotifyingRecorder()

	fmt.Println("------------------------------ request start ------------------------------")
	fmt.Printf("=>request url:%s\n", req.URL.RequestURI())
	fmt.Printf("=>request header:\n%v\n", indentJson(req.Header))
	if method == http.MethodPost {
		fmt.Printf("=>request body:\n%s\n", string(bb))
	}
	fmt.Println("------------------------------ request end --------------------------------")

	h := gin.Default()

	routePath := r.path
	i := strings.Index(routePath, "?")
	if i != -1 {
		routePath = routePath[0:i]
	}

	h.Handle(method, routePath, r.handlers...)

	h.ServeHTTP(w, req)
	resp := &TestResponse{
		w.ResponseRecorder,
	}

	fmt.Println("------------------------------ response start -----------------------------")
	fmt.Printf("<=response headers:\n%v\n", indentJson(resp.Result().Header))
	fmt.Printf("<=response body:\n%s\n", resp.GetBody())
	fmt.Println("------------------------------ response end -------------------------------")
	return resp, nil
}

type TestResponse struct {
	*httptest.ResponseRecorder
}

func (r *TestResponse) GetBody() []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		r.Body = bytes.NewBuffer(body)
	}
	return body
}

func (r *TestResponse) GetBodyString() string {
	return string(r.GetBody())
}

func (r *TestResponse) GetJsonPath(path string) gjson.Result {
	return gjson.GetBytes(r.GetBody(), path)
}

func (r *TestResponse) GetJson() gjson.Result {
	return gjson.ParseBytes(r.GetBody())
}

func indentJson(data interface{}) string {
	v, _ := json.MarshalIndent(data, "", "\t")
	return string(v)
}
