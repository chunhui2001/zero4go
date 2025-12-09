package server

import (
	"fmt"

	"github.com/gin-gonic/gin"

	. "github.com/chunhui2001/zero4go/pkg/logs"
)

func (a *Application) POST(relativePath string, handlers ...interface{}) gin.IRoutes {
	Log.Infof("%s, Path=%s, Handlers=%s", fmt.Sprintf("%6s", "POST"), relativePath, JoinHandlersString(handlers))

	return a.Engine.POST(relativePath, Wrap(handlers...)...)
}

func (a *Application) GET(relativePath string, handlers ...interface{}) gin.IRoutes {
	Log.Infof("%s, Path=%s, Handlers=%s", fmt.Sprintf("%6s", "GET"), relativePath, JoinHandlersString(handlers))

	return a.Engine.GET(relativePath, Wrap(handlers...)...)
}

func (a *Application) PUT(relativePath string, handlers ...interface{}) gin.IRoutes {
	Log.Infof("%s, Path=%s, Handlers=%s", fmt.Sprintf("%6s", "PUT"), relativePath, JoinHandlersString(handlers))

	return a.Engine.PUT(relativePath, Wrap(handlers...)...)
}

func (a *Application) DELETE(relativePath string, handlers ...interface{}) gin.IRoutes {
	Log.Infof("%s, Path=%s, Handlers=%s", fmt.Sprintf("%6s", "DELETE"), relativePath, JoinHandlersString(handlers))

	return a.Engine.DELETE(relativePath, Wrap(handlers...)...)
}
