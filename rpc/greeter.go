package rpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/chunhui2001/zero4go/rpc/gen"

)

// gRPC server implementation
type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

// grpcurl -plaintext -d '{"name":"keesh 阿斯顿发的啥饭"}' localhost:50051 rpc.Greeter/SayHello
func (s *GreeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	//return &pb.HelloReply{Message: "Hello " + req.Name}, nil

	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}

// grpcurl -plaintext localhost:50051 rpc.Greeter/SayInfo
func (s *GreeterServer) SayInfo(ctx context.Context, req *emptypb.Empty) (*pb.InfoReply, error) {
	
	return &pb.InfoReply{Message: "Yeah, your server is running"}, nil
}
