package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-chat/internal/pkg/jsonutil"
)

type RequestInfo struct {
	HostName        string            `json:"host_name"`
	RemoteAddr      string            `json:"remote_addr"`
	ServerName      string            `json:"server_name"`
	HttpUserAgent   string            `json:"http_user_agent"`
	RequestId       string            `json:"request_id"`
	RequestTime     string            `json:"request_time"`
	RequestMethod   string            `json:"request_method"`
	RequestHeader   map[string]string `json:"request_header"`
	RequestUri      string            `json:"request_uri"`
	RequestQuery    map[string]any    `json:"request_query"`
	RequestBody     string            `json:"request_body"`
	RequestBodyRaw  string            `json:"request_body_raw"`
	ResponseHeader  map[string]string `json:"response_header"`
	ResponseBody    string            `json:"response_body"`
	ResponseBodyRaw string            `json:"response_body_raw"`
	ResponseTime    string            `json:"response_time"`
	HttpStatus      int               `json:"http_status"`
	RequestDuration string            `json:"request_duration"`
	Metadata        map[string]any    `json:"metadata"`
}

type AccessFilterOption struct {
	path string
	fn   func(info *RequestInfo)
}

type AccessFilterRule struct {
	ops           map[string]AccessFilterOption
	excludeRoutes []string
}

func NewAccessFilterRule() *AccessFilterRule {
	return &AccessFilterRule{ops: make(map[string]AccessFilterOption), excludeRoutes: make([]string, 0)}
}

func (a *AccessFilterRule) Exclude(path string) {
	a.excludeRoutes = append(a.excludeRoutes, path)
}

func (a *AccessFilterRule) AddRule(path string, fn func(access *RequestInfo)) {
	a.ops[path] = AccessFilterOption{
		path: path,
		fn:   fn,
	}
}

func (a *AccessFilterRule) filter(access *AccessLogStore) {
	if rule, ok := a.ops[access.info.RequestUri]; ok {
		rule.fn(access.info)
	}
}

func AccessLog(w io.Writer, filterRule *AccessFilterRule) gin.HandlerFunc {
	log := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05.000"))
			}

			return a
		},
	}))

	return func(c *gin.Context) {
		c.Writer = responseWriter{c.Writer, bytes.NewBuffer([]byte{})}

		access := newAccessLogStore(c)
		if err := access.init(); err != nil {
			c.Abort()
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
			return
		}

		c.Next()

		if c.Request.Method != "OPTIONS" {
			access.load()

			if filterRule != nil {
				filterRule.filter(access)
			}

			if slices.Contains(filterRule.excludeRoutes, access.info.RequestUri) {
				return
			}

			access.save(log)
		}
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type AccessLogStore struct {
	ctx       *gin.Context
	startTime time.Time
	info      *RequestInfo
}

func newAccessLogStore(c *gin.Context) *AccessLogStore {
	return &AccessLogStore{
		ctx:       c,
		startTime: time.Now(),
		info:      nil,
	}
}

func (a *AccessLogStore) init() error {
	hostname, _ := os.Hostname()

	headers := make(map[string]string)
	for k := range a.ctx.Request.Header {
		headers[k] = a.ctx.Request.Header.Get(k)
	}

	body, err := io.ReadAll(a.ctx.Request.Body)
	if err != nil {
		return err
	}

	a.ctx.Request.Body = io.NopCloser(bytes.NewReader(body))

	a.info = &RequestInfo{
		HostName:        hostname,
		RemoteAddr:      a.ctx.RemoteIP(),
		ServerName:      a.ctx.Request.Host,
		HttpUserAgent:   a.ctx.Request.UserAgent(),
		RequestId:       a.ctx.Request.Header.Get("X-Request-ID"),
		RequestTime:     a.startTime.Format("2006-01-02 15:04:05"),
		RequestMethod:   a.ctx.Request.Method,
		RequestHeader:   headers,
		RequestUri:      a.ctx.Request.URL.Path,
		RequestQuery:    urlValuesToMap(a.ctx.Request.URL.Query()),
		RequestBody:     string(body),
		ResponseHeader:  make(map[string]string),
		ResponseBody:    "",
		RequestBodyRaw:  "",
		ResponseTime:    "",
		HttpStatus:      0,
		RequestDuration: "",
		Metadata:        make(map[string]any),
	}

	if a.info.RequestId == "" {
		a.info.RequestId = uuid.New().String()
	}

	return nil
}

func (a *AccessLogStore) load() {
	writer := a.ctx.Writer.(responseWriter)

	headers := make(map[string]string)
	for k := range writer.Header() {
		headers[k] = writer.Header().Get(k)
	}

	a.info.ResponseHeader = headers
	a.info.ResponseTime = time.Now().Format("2006-01-02 15:04:05")
	a.info.RequestDuration = fmt.Sprintf("%.3f", time.Since(a.startTime).Seconds())
	a.info.HttpStatus = writer.Status()
	a.info.ResponseBody = writer.body.String()
	a.info.ResponseBodyRaw = a.info.ResponseBody

	session, isOk := a.ctx.Get(JWTSessionConst)
	if isOk {
		a.info.Metadata["uid"] = session.(*JSession).Uid
	}
}

func (a *AccessLogStore) save(log *slog.Logger) {
	data := make(map[string]any)
	if err := jsonutil.Decode(jsonutil.Encode(a.info), &data); err != nil {
		return
	}

	if strings.HasPrefix(a.ctx.GetHeader("Content-Type"), "application/json") {
		var body map[string]any
		_ = json.Unmarshal([]byte(a.info.RequestBody), &body)

		data["request_body"] = body
		delete(data, "request_body_raw")
	} else {
		delete(data, "request_body")
	}

	writer := a.ctx.Writer.(responseWriter)
	if strings.HasPrefix(writer.Header().Get("Content-Type"), "application/json") {
		var body map[string]any
		_ = json.Unmarshal([]byte(a.info.ResponseBody), &body)

		data["response_body"] = body
		delete(data, "response_body_raw")
	} else {
		delete(data, "response_body")
	}

	items := make([]any, 0)
	for k, v := range data {
		items = append(items, k, v)
	}

	log.With(items...).Info("access_log")
}

func urlValuesToMap(values url.Values) map[string]any {
	data := make(map[string]any)
	for k, v := range values {
		if len(v) == 1 {
			data[k] = v[0]
		} else {
			data[k] = v
		}
	}
	return data
}
