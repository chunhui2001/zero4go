package server

import (
	"strings"

	"github.com/gin-gonic/gin"

	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
	"github.com/chunhui2001/zero4go/pkg/upstream"
)

type SiteConf struct {
	Root      string `mapstructure:"WEB_PAGE_ROOT"`
	Master    string `mapstructure:"WEB_PAGE_MASTER"`
	Extension string `mapstructure:"WEB_PAGE_Extension"`
	LoginUrl  string `mapstructure:"WEB_PAGE_LOGIN"`
	SignUpUrl string `mapstructure:"WEB_PAGE_SIGNUP"`
}

var _SiteConf = &SiteConf{
	Root:      "views",
	Master:    "layouts/master",
	Extension: ".html",
}

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

func (a *Application) Upstream(from string, to string, remotes ...string) {
	a.Any(from+"/*proxyPath", func(c *gin.Context) {
		upstream.Proxy(c, to, remotes...)
	})

	Log.Infoe1().Msgf("ANY, From=%s, To=%s, Remotes=%s", from, to, strings.Join(remotes, ","))
}
