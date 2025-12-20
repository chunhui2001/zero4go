package main

import (
	"time"

	"google.golang.org/grpc"

	"github.com/chunhui2001/zero4go/pkg/gtask"
	. "github.com/chunhui2001/zero4go/pkg/server"
)

func main() {
	Setup(func(r *Application) {

	}).Run(func(grpcServer *grpc.Server) {
		gtask.AddTask("job", "job1", "0/1 * * * * *", func(key string) {
			time.Sleep(15 * time.Second)
		})
	})
}
