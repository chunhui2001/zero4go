package interceptors

import "github.com/gin-gonic/gin"

type AccessInterceptorWrapper struct {
	requiredAuth bool
}

func (w *AccessInterceptorWrapper) Process(c *gin.Context) {

	if w.requiredAuth {

	}

	c.Next()
}

func AccessInterceptor(requiredAuth bool) gin.HandlerFunc {
	return (&AccessInterceptorWrapper{requiredAuth}).Process
}
