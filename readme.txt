
# go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
# go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
# git clone https://github.com/googleapis/googleapis.git third_party/googleapis
# grpcurl -plaintext localhost:50051 list
# grpcurl -plaintext -d '{"name":"keesh 阿斯顿发的啥饭"}' localhost:50051 rpc.Greeter/SayHello
