package logger

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessLog struct {
	Method     string
	Url        string
	ReqBody    string
	RespBody   string
	Duration   string
	StatusCode int
}
type MiddlewareBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc: fn,
	}
}
func (m *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		//访问al
		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}
		if m.allowReqBody && ctx.Request.Body != nil {
			//body读取完就没了(流)
			body, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body)) //读回body
			if len(body) > 1024 {
				body = body[:1024]
			}
			al.ReqBody = string(body)
		}
		start := time.Now()
		if m.allowRespBody {
			ctx.Writer = responseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}
		defer func() {
			al.Duration = time.Since(start).String()
			m.loggerFunc(ctx, al)
		}()
		//执行业务逻辑
		ctx.Next()
		//m.loggerFunc(ctx, al)
	}
}

func (m *MiddlewareBuilder) AllowRepBody() *MiddlewareBuilder {
	m.allowRespBody = false
	return m
}
func (m *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	m.allowRespBody = false
	return m
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (r responseWriter) WriteHeader(statusCode int) {
	r.al.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r responseWriter) Write(data []byte) (int, error) {
	r.al.RespBody = string(data)
	return r.ResponseWriter.Write(data)
}

func (r responseWriter) WriteString(data string) (int, error) {
	r.al.RespBody = data
	return r.ResponseWriter.WriteString(data)
}
