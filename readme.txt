
# 📚 验证是否成功发布到 Go Proxy
https://proxy.golang.org/github.com/chunhui2001/zero4go/@v/list

# 查看某个 Go Module 的 最新版本
$ go list -m -versions github.com/chunhui2001/zero4go

# 直连 GitHub
# GOPROXY=direct go list -m github.com/chunhui2001/zero4go@v1.0.0

# 检查你仓库真实存在的版本
$ git ls-remote --tags https://github.com/chunhui2001/zero4go.git

# 快速「一锤定音」诊断命令
$ GOPROXY=direct go get -x github.com/chunhui2001/zero4go@v1.0.0


# 你以后发版一定要记住的铁律
> Go module 的版本号一旦被请求过, 就永远不能“补救”
> commit → tag → push → 再 go get
> ❌ 反过来一次，这个版本号就“污染”了

# 清理 mod cache
$ go clean -modcache

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

$ go clean -modcache && go get -u github.com/chunhui2001/zero4go@latest

$ go get && go run .

# git 修改最后一次提交的 message
$ git commit --amend -m "upgrade redis to v9"

# 真正适合 Go SSR 组件库
https://templ.guide/
> 它是什么？
  编译期组件
  真正的 props
  真正的 children
  真正的组合能力
  输出 纯 HTML
  100% 服务端渲染


# templ 全语言对照表（直接帮你选）
| 语言   | templ 对等选择 | 推荐指数  |
| ---- | ---------- | ----- |
| Go   | templ      | ⭐⭐⭐⭐⭐ |
| Rust | Leptos     | ⭐⭐⭐⭐⭐ |
| Rust | Askama     | ⭐⭐⭐⭐  |
| Java | JTE        | ⭐⭐⭐⭐  |
| Java | Thymeleaf  | ⭐⭐⭐   |
| Node | Astro      | ⭐⭐⭐⭐⭐ |
| Node | React SSR  | ⭐⭐⭐⭐  |

# rbac_model.conf
-------------------------
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act

解释：
  r: 请求（subject, object, action）
  p: 策略（谁可以做什么）
  g: 角色关系
  m: 匹配规则

# policy.csv
-------------------------
p, admin, /admin, GET
p, admin, /admin, POST
p, user, /profile, GET

g, alice, admin
g, bob, user

解释：
  p 表示策略：角色可以访问某个 URL 做某种操作
  g 表示用户属于哪个角色

# 验证权限
$ curl -H "X-User: alice" http://localhost:8080/admin
> 返回: {"message":"welcome admin"}

$ curl -H "X-User: bob" http://localhost:8080/admin
> 返回: {"error":"forbidden"}

$ curl -H "X-User: bob" http://localhost:8080/profile
> 返回: {"message":"profile page"}


