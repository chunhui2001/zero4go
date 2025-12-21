package middlewares

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type AccessWrapper struct {
	requiredAuth bool
}

func (w *AccessWrapper) Process(c *gin.Context) {
	if !w.requiredAuth || accessClientsMap == nil || len(accessClientsMap) == 0 {
		c.Next()
		return
	}

	requestUrl := RequestURL(c.Request)
	var accessQuery = requestUrl.Query()

	if !accessQuery.Has(AWSAccessKeyIdFieldKey) {
		AbortAccess(errors.New("ACCESS_KEY_REQUIRED"), c)
		c.Next()
		return
	}

	accessKeyId := accessQuery.Get(AWSAccessKeyIdFieldKey)

	var accessClient = accessClientsMap[accessKeyId]

	if accessClient == nil {
		AbortAccess(errors.New("ACCESS_KEY_INVALID"), c)
		c.Next()

		return
	}

	if !accessClient.Enabled {
		AbortAccess(errors.New("ACCESS_KEY_DISABLED"), c)
		c.Next()

		return
	}

	if _, err := CheckSign(accessKeyId, accessClient.SecretAccessKey, c.Request.Method, requestUrl); err != nil {
		AbortAccess(err, c)
		c.Next()

		return
	}

	c.Next()
}

func Access(requiredAuth bool) gin.HandlerFunc {
	return (&AccessWrapper{requiredAuth}).Process
}

func AbortAccess(err error, c *gin.Context) {
	c.String(401, err.Error())

	if e := c.Error(err); e != nil {
		c.Abort()
		return
	}

	c.Abort()
}
