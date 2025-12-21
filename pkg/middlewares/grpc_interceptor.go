package middlewares

import (
	"strings"

	"context"
	"time"

	"google.golang.org/grpc"

	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
)

// UnaryLoggingInterceptor Unary interceptor
func UnaryLoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()

	// 调用真正的 RPC 方法
	resp, err = handler(ctx, req)

	Log.Infof("gRPC call: %s completed in %s",
		info.FullMethod, time.Since(start))

	return resp, err
}

func StreamLoggingInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	if strings.HasPrefix(info.FullMethod, "/grpc.reflection.v1.ServerReflection/") {
		return handler(srv, ss) // 忽略打印
	}

	start := time.Now()

	err := handler(srv, ss)

	Log.Infof("gRPC stream call: %s finished in %s",
		info.FullMethod, time.Since(start))

	return err
}
