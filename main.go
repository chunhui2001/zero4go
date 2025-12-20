package main

import (
	"strconv"

	"google.golang.org/grpc"

	"github.com/chunhui2001/zero4go/pkg/gredis"
	. "github.com/chunhui2001/zero4go/pkg/logs"
	. "github.com/chunhui2001/zero4go/pkg/server"
)

func main() {
	Setup(func(r *Application) {
		var q = gredis.RedisBlockingQueue[uint64]{
			Key:      "__kk2",
			MaxCount: 2,
		}

		r.GET("/redis_push", func(c *RequestContext) {
			var v = c.Query("val")

			n, _ := strconv.Atoi(v)

			q.Push(uint64(n))

			Log.Infof("redis_push: val=%d", n)

			c.OK(q)
		})

		r.GET("/redis_pop", func(c *RequestContext) {
			var v = q.Pop(1)

			c.OK(v)
		})
	}).Run(func(grpcServer *grpc.Server) {
		//gtask.AddTask("job", "job1", "0/1 * * * * *", func(key string) {
		//	time.Sleep(15 * time.Second)
		//})
	})
}
