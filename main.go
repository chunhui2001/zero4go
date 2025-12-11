package main

import (
	"google.golang.org/grpc"

	"github.com/chunhui2001/zero4go/api"
	. "github.com/chunhui2001/zero4go/pkg/server"
)

func main() {
	Setup(func(r *Application) {
		r.GET("/info2", func(c *RequestContext) {
			c.Text("Yeah, your server is running.")
		})

		r.GET("/send_message", api.SendKafkaMessage)
		r.GET("/insert", api.Insert)

	}).Run(func(grpcServer *grpc.Server) {

	})
}
