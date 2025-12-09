package main

import (
	"google.golang.org/grpc"

	"github.com/chunhui2001/zero4go/api"
	. "github.com/chunhui2001/zero4go/pkg/server"

	"github.com/chunhui2001/zero4go/pkg/gkafka"
)

func main() {
	Setup(func(r *Application) {
		r.GET("/info2", func(c *RequestContext) {
			c.Text("Yeah, your server is running.")
		})

		r.GET("/send_message", api.SendKafkaMessage)

	}).Run(func(grpcServer *grpc.Server) {

		gkafka.CreateConsumer("localhost:9092", "group1", "first-topic", func(key string, message string) {
			//Log.Infof("消费了一个消息: key=%s", key)
		})
	})
}
