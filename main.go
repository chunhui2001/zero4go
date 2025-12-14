package main

import (
	"math/rand"
	"time"

	"google.golang.org/grpc"

	"github.com/chunhui2001/zero4go/api"
	"github.com/chunhui2001/zero4go/pkg/gtask"
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

		gtask.AddTask("一个示例定时任务", "job1", "0/5 * * * * *", func(key string) {

			var v = rand.Intn(10) + 1

			time.Sleep(time.Duration(v) * time.Second)
		})
	})
}
