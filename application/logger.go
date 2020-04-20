package application

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Logger for log middleware
type Logger interface {
	Interceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
	Middleware() gin.HandlerFunc
}

// SkipperFunc check if method need to be skip
type SkipperFunc func(method string) bool

// LogConfig middleware config
type LogConfig struct {
	ProjectName    bool
	ProjectVersion bool

	// Method displays the request path (bool).
	// Defaults to true.
	Path bool

	// used in http request , show http method
	// default true
	Method bool

	// ClientIP
	// default true
	ClientIP bool

	// used in http request , show http status
	Status bool

	// used in http request , show body size
	BodySize bool

	// Req show request info
	// Defaults to true.
	Request bool
	// Response show response einfo
	// Defaults to true.
	Response bool
	// Error show error messagee
	// Defaults to true.
	Error bool

	// Latency show latency
	// Defaulte to true
	Latency bool

	// Runtime
	// Defaulte to true
	Runtime bool

	// Defaults nil
	LogFunc func(now time.Time, projectName, projectVersion, method string, latency time.Duration, runtimeValue map[string]string)

	// Skippers will according method to see if skip the logger middleware
	Skippers []SkipperFunc
}

// DefaultLogConfig default log  miiddleware setting
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		ProjectName:    true,
		ProjectVersion: true,
		Path:           true,
		Method:         true,
		Request:        true,
		Response:       true,
		Error:          true,
		Latency:        true,
		Runtime:        true,
		Status:         true,
		BodySize:       true,
		LogFunc:        nil,
		Skippers:       nil,
	}
}

type loggerImpl struct {
	app    Application
	config *LogConfig
}

// NewLogLogger new log logger
func NewLogLogger(app Application, config *LogConfig) Logger {
	return &loggerImpl{app, config}
}

func (l *loggerImpl) Interceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		// skip the method
		for _, v := range l.config.Skippers {
			if v(info.FullMethod) {
				return handler(ctx, req)
			}
		}

		var path string
		var latency time.Duration
		var startTime, endTime time.Time
		startTime = time.Now()
		resp, err := handler(ctx, req)
		endTime = time.Now()

		line := ""

		if l.config.Runtime {
			md, _ := metadata.FromIncomingContext(ctx)
			for _, v := range l.app.RuntimeKeys() {
				if v.IsLog() {
					line += fmt.Sprintf("%v:%v ", v.GetKeyName(), md[v.GetKeyName()][0])
				}
			}
		}

		if l.config.Latency {
			latency = endTime.Sub(startTime)
			line += fmt.Sprintf("%4v ", latency)
		}

		if l.config.Path {
			path = info.FullMethod
			line += fmt.Sprintf("%v ", path)
		}

		if l.config.Request {
			line += fmt.Sprintf("%v %v ", "Request", req)
		}
		if l.config.Response {
			line += fmt.Sprintf("%v %v ", "Response", resp)
		}
		if l.config.Error && err != nil {
			line += fmt.Sprintf("%v %v ", "Error", err)
		}
		l.app.Logger().Info(line)
		return resp, err
	}
}

func (l *loggerImpl) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// skip the method
		for _, v := range l.config.Skippers {
			if v(c.FullPath()) {
				c.Next()
				return
			}
		}

		var method string
		var path string
		var clientIP string
		var status int
		var bodySize int
		var latency time.Duration
		var startTime, endTime time.Time
		startTime = time.Now()

		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		request, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)

		// request := []byte{}
		c.Next()

		endTime = time.Now()

		line := ""
		if l.config.Runtime {
			runtimeKey := DecodeHTTPRuntimeKey(c, l.app.RuntimeKeys())
			for _, v := range l.app.RuntimeKeys() {
				if v.IsLog() {
					line += fmt.Sprintf("%v:%v ", v.GetKeyName(), runtimeKey[v.GetKeyName()])
				}
			}
		}

		if l.config.Latency {
			latency = endTime.Sub(startTime)
			line += fmt.Sprintf("%4v ", latency)
		}

		if l.config.Path {
			path = c.Request.URL.RequestURI()
			line += fmt.Sprintf("%v ", path)
		}

		if l.config.Method {
			method = c.Request.Method
			line += fmt.Sprintf("%v ", method)
		}

		if l.config.ClientIP {
			clientIP = c.ClientIP()
			line += fmt.Sprintf("%v ", clientIP)
		}

		if l.config.Status {
			status = c.Writer.Status()
			line += fmt.Sprintf("%v ", status)
		}

		if l.config.BodySize {
			bodySize = c.Writer.Size()
			line += fmt.Sprintf("%v ", bodySize)
		}

		if l.config.Request {
			line += fmt.Sprintf("%v ", string(request))
		}

		l.app.Logger().Info(line)
	}
}
