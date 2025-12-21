package sitepages

import (
	"net/http"

	"github.com/chunhui2001/zero4go/pkg/server"
	"github.com/gin-gonic/gin"
)

func Index(c *server.RequestContext) {
	c.Render(http.StatusOK, "index", gin.H{
		"content": "This is an Index page...",
	})
}
