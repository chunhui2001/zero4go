package main

import (
	"google.golang.org/grpc"

	. "github.com/chunhui2001/zero4go/pkg/server" //nolint:staticcheck
)

func main() {
	Setup(func(r *Application) {

	}).Run(func(grpcServer *grpc.Server) {

	})
}
