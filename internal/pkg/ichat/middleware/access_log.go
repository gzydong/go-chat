package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func AccessLog() gin.HandlerFunc {

	lgs := logrus.New()
	lgs.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return func(c *gin.Context) {

		writer := responseWriter{c.Writer, bytes.NewBuffer(nil)}
		c.Writer = writer

		access := newAccessLogStore(c)
		if err := access.init(); err != nil {
			c.Abort()
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
			return
		}

		c.Next()

		access.load(writer)

		if c.Request.Method != "OPTIONS" {
			lgs.WithFields(access.data).Info("access_log")
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
	data      map[string]any
}

func newAccessLogStore(c *gin.Context) *AccessLogStore {
	return &AccessLogStore{
		ctx:       c,
		startTime: time.Now(),
		data:      make(map[string]any),
	}
}

func (a *AccessLogStore) init() error {

	hostname, _ := os.Hostname()

	headers := make(map[string]any)
	for k := range a.ctx.Request.Header {
		headers[k] = a.ctx.Request.Header.Get(k)
	}

	body, err := io.ReadAll(a.ctx.Request.Body)
	if err != nil {
		return err
	}

	// 请求日志信息
	a.data = map[string]any{
		"request_id":        a.ctx.Request.Header.Get("X-Request-ID"),
		"request_method":    a.ctx.Request.Method,
		"request_header":    headers,
		"request_uri":       a.ctx.Request.URL.Path,
		"request_body":      string(body),
		"request_time":      a.startTime.Format("2006-01-02 15:04:05"),
		"request_duration":  "",
		"request_query":     urlValuesToMap(a.ctx.Request.URL.Query()),
		"request_body_raw":  "",
		"response_header":   []string{},
		"response_body_raw": "",
		"response_time":     time.Now().Format("2006-01-02 15:04:05"),
		"http_user_agent":   a.ctx.Request.UserAgent(),
		"http_status":       0,
		"host_name":         hostname,
		"server_name":       a.ctx.Request.Host,
		"remote_addr":       a.ctx.RemoteIP(),
	}

	if a.data["request_id"] == "" {
		a.data["request_id"] = strings.ReplaceAll(uuid.New().String(), "-", "")
	}

	if strings.HasPrefix(a.ctx.GetHeader("Content-Type"), "application/json") {
		a.data["request_body"] = &map[string]any{}
		_ = json.Unmarshal(body, a.data["request_body"])
	}

	a.ctx.Request.Body = io.NopCloser(bytes.NewReader(body))

	return nil
}

func (a *AccessLogStore) load(writer responseWriter) {

	headers := make(map[string]any)
	for k := range writer.Header() {
		headers[k] = writer.Header().Get(k)
	}

	a.data["response_header"] = headers
	a.data["response_time"] = time.Now().Format("2006-01-02 15:04:05")
	a.data["request_duration"] = fmt.Sprintf("%.3f", time.Since(a.startTime).Seconds())
	a.data["http_status"] = writer.Status()
	a.data["response_body_raw"] = writer.body.String()

	if strings.HasPrefix(writer.Header().Get("Content-Type"), "application/json") {
		a.data["response_body"] = &map[string]any{}
		_ = json.Unmarshal(writer.body.Bytes(), a.data["response_body"])
		delete(a.data, "response_body_raw")
	}
}

func urlValuesToMap(values url.Values) map[string]any {
	qm := make(map[string]any)
	for k, v := range values {
		if len(v) == 1 {
			qm[k] = v[0]
		} else {
			qm[k] = v
		}
	}
	return qm
}
