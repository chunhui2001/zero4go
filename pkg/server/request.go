package server

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

type RequestContext struct {
	*gin.Context
}

func (c *RequestContext) UserID() string {
	return c.GetString("user_id")
}

func (c *RequestContext) OK(data any) {
	c.JSON(200, gin.H{"code": 200, "data": data})
}

func (c *RequestContext) Fail(err any) {
	c.JSON(200, gin.H{"code": 400, "msgs": err})
}

func (c *RequestContext) Failf(format string, err ...any) {
	c.JSON(200, gin.H{"code": 400, "msg": fmt.Sprintf(format, err...)})
}

func (c *RequestContext) Failc(code int, format string, err ...any) {
	c.JSON(200, gin.H{"code": code, "msg": fmt.Sprintf(format, err...)})
}

func (c *RequestContext) Text(s any) {

	switch s := s.(type) {
	case float64, float32:
		c.Data(200, "text/plain; charset=utf-8", []byte(fmt.Sprint(s)))
	case string:
		c.Data(200, "text/plain; charset=utf-8", []byte(s))
	case bool:
		c.Data(200, "text/plain; charset=utf-8", []byte(fmt.Sprintf("%t", s)))
	case byte:
		c.Data(200, "text/plain; charset=utf-8", []byte(fmt.Sprintf("%x", s)))
	case []uint8:
		c.Data(200, "text/plain; charset=utf-8", s)
	default:
		c.Data(200, "text/plain; charset=utf-8", []byte(fmt.Sprintf("%d", s)))
	}
}

func Wrap(handlers ...interface{}) []gin.HandlerFunc {
	wrapped := make([]gin.HandlerFunc, 0, len(handlers))

	for _, h := range handlers {

		switch handler := h.(type) {

		case func(*RequestContext):
			// 用户自定义 handler
			wrapped = append(wrapped, func(c *gin.Context) {
				handler(&RequestContext{Context: c})
			})

		case func(*gin.Context):
			wrapped = append(wrapped, handler)

		case gin.HandlerFunc: // == func(*gin.Context)
			wrapped = append(wrapped, handler)

		default:
			panic(fmt.Sprintf("unsupported handler type: %T", h))
		}
	}

	return wrapped
}

func JoinHandlersString(handlers ...interface{}) string {
	var names []string

	for _, h := range handlers {
		names = append(names, extractFuncNames(h)...)
	}

	return strings.Join(names, ", ")
}

func extractFuncNames(h interface{}) []string {
	v := reflect.ValueOf(h)
	t := reflect.TypeOf(h)

	// 如果是函数 —— 直接获取名字
	if t.Kind() == reflect.Func {
		return []string{GetFunctionName(h)}
	}

	// 如果是 slice / array —— 递归处理每个元素
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		var out []string

		for i := 0; i < v.Len(); i++ {
			out = append(out, extractFuncNames(v.Index(i).Interface())...)
		}

		return out
	}

	// 其他类型 —— 忽略
	return []string{}
}

func GetFunctionName(i interface{}) string {
	v := reflect.ValueOf(i)

	if v.Kind() != reflect.Func {
		return ""
	}

	fn := runtime.FuncForPC(v.Pointer())

	if fn == nil {
		return "unknown"
	}

	return fn.Name()
}
