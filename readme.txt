
# 📚 验证是否成功发布到 Go Proxy
https://proxy.golang.org/github.com/chunhui2001/zero4go/@v/list

# 查看某个 Go Module 的 最新版本
$ go list -m -versions github.com/gin-contrib/pongo2

# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
# go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
# git clone https://github.com/googleapis/googleapis.git third_party/googleapis
# grpcurl -plaintext localhost:50051 list
# grpcurl -plaintext -d '{"name":"keesh 阿斯顿发的啥饭"}' localhost:50051 rpc.Greeter/SayHello


# 📚 graphql
$ go get github.com/99designs/gqlgen

# 清理旧缓存
$ go clean -modcache

# 升级 x/tools（可选）
$ go get golang.org/x/tools@latest

# 确保你安装的是最新版本：
$ go install github.com/99designs/gqlgen@latest

# 重新生成 gqlgen
$ gqlgen init

# 或者如果已经有旧的 schema，可以直接：
$ gqlgen generate

> 会生成 graph/schema.graphqls、graph/resolver.go 等文件



$ mkdir zero4go01 && cd zero4go01
$ go mod init zero4go01

github.com/chunhui2001/zero4go
$ go get -u github.com/chunhui2001/zero4go@latest

$ go get && go run .

# git 修改最后一次提交的 message
$ git commit --amend -m "upgrade redis to v9"