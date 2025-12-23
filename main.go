package main

import (
	"google.golang.org/grpc"

	. "github.com/chunhui2001/zero4go/pkg/server" //nolint:staticcheck
	"github.com/chunhui2001/zero4go/sitepages"
)

func main() {
	Setup(func(r *Application) {
		r.GET("/", sitepages.Index)
		r.GET("/index", sitepages.Index)

		r.GET("/admin", func(c *RequestContext) {
			c.OK("get admin")
		})

		r.POST("/admin", func(c *RequestContext) {
			c.OK("post admin")
		})

		r.GET("/profile", func(c *RequestContext) {

			c.OK("profile")
		})
	}).Run(func(grpcServer *grpc.Server) {

	})
}
