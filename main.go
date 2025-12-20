package main

import (
	iter "github.com/chunhui2001/zero4go/pkg/interceptors"
	"github.com/chunhui2001/zero4go/pkg/utils"
	"google.golang.org/grpc"

	"github.com/chunhui2001/zero4go/pkg/config"
	. "github.com/chunhui2001/zero4go/pkg/server" //nolint:staticcheck
)

func main() {
	Setup(func(r *Application) {
		r.GET("/configuration", iter.Access(true), func(r *RequestContext) {
			var v = config.GetConfig("Access-Clients").(map[string]interface{})
			var lll = make(map[string]iter.AccessClients)

			for key, val := range v {
				if vv, err := utils.ToStructByTag[iter.AccessClients](val, "yaml"); err == nil {
					lll[key] = vv
				}
			}

			r.OK(lll)
		})
	}).Run(func(grpcServer *grpc.Server) {

	})
}
