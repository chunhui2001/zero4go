package server

import (
	"github.com/gin-gonic/gin"

	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
)

func (a *Application) POST(relativePath string, handlers ...interface{}) gin.IRoutes {
	Log.Infoe1().Msgf("%s, Path=%s, Handlers=%s", "POST", relativePath, JoinHandlersString(handlers))

	return a.Engine.POST(relativePath, Wrap(handlers...)...)
}

func (a *Application) GET(relativePath string, handlers ...interface{}) gin.IRoutes {
	Log.Infoe1().Msgf("%s, Path=%s, Handlers=%s", "GET", relativePath, JoinHandlersString(handlers))

	return a.Engine.GET(relativePath, Wrap(handlers...)...)
}

func (a *Application) PUT(relativePath string, handlers ...interface{}) gin.IRoutes {
	Log.Infoe1().Msgf("%s, Path=%s, Handlers=%s", "PUT", relativePath, JoinHandlersString(handlers))

	return a.Engine.PUT(relativePath, Wrap(handlers...)...)
}

func (a *Application) DELETE(relativePath string, handlers ...interface{}) gin.IRoutes {
	Log.Infoe1().Msgf("%s, Path=%s, Handlers=%s", "DELETE", relativePath, JoinHandlersString(handlers))

	return a.Engine.DELETE(relativePath, Wrap(handlers...)...)
}
