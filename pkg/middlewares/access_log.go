package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
)

var DefaultLogFormatter = func(param gin.LogFormatterParams) string {

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	return fmt.Sprintf("Access %s \"%s %s %s %d %s\"",
		param.ClientIP,
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		//param.Request.UserAgent(),
	)
}

func init() {

}

func AccessLog(skips ...string) gin.HandlerFunc {
	return Print(gin.LoggerConfig{
		SkipPaths: skips,
	})
}

func Print(conf gin.LoggerConfig) gin.HandlerFunc {

	notlogged := conf.SkipPaths

	var skip map[string]bool

	if length := len(notlogged); length > 0 {
		skip = make(map[string]bool, length)

		for _, path := range notlogged {
			skip[path] = true
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			Log.Info(DefaultLogFormatter(LogParam(c, c.Writer.Status(), path, start)))
		}
	}
}

func LogParam(c *gin.Context, code int, path string, start time.Time) gin.LogFormatterParams {

	raw := c.Request.URL.RawQuery

	param := gin.LogFormatterParams{
		Request: c.Request,
		Keys:    c.Keys,
	}

	// Stop timer
	param.TimeStamp = time.Now()
	param.Latency = param.TimeStamp.Sub(start)

	param.ClientIP = c.Request.RemoteAddr
	param.Method = c.Request.Method
	param.StatusCode = code
	param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

	param.BodySize = c.Writer.Size()

	if raw != "" {
		path = path + "?" + raw
	}

	param.Path = path

	return param
}
