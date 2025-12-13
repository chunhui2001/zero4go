
### 当前 Makefile 文件物理路径
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

APP_NAME 	?=zero4go
e 			?=local
#e 			?=production
c 			?=10000
zone 		?=Asia/Shanghai
#zone 		?=UTC
#WSS_HOST	?=ws://127.0.0.1:8080
APP_PORT 	?=8080
GIT_HASH 	?=$(shell git rev-parse HEAD)
COMMITER 	?=$(shell git log -1 --pretty=format:'%ae')
PWD 		?=$(shell pwd)
TIME 		?=$(shell date +%s)
CGO_ENABLED ?=0
GOPROXY 	?=go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct

PROTO_SRC	=./proto
OUT_DIR		=./rpc/gen

### 整理模块
# 确保go.mod与模块中的源代码一致。
# 它添加构建当前模块的包和依赖所必须的任何缺少的模块，删除不提供任何有价值的包的未使用的模块。
# 它也会添加任何缺少的条目至go.mod并删除任何不需要的条目。
# make tidy
tidy:
	go mod tidy

### 显示已安装的模块
# show install utils
list:
	ls -alh `go env GOPATH`/bin

### 安装模块
# make install mod=github.com/codegangsta/gin
install:
	@#$(GOPROXY) && go get github.com/codegangsta/gin
	@#$(GOPROXY) && go install github.com/codegangsta/gin
	@#$(GOPROXY) && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@#$(GOPROXY) && go install github.com/99designs/gqlgen
	$(GOPROXY) && go get -u $(mod)
	@#$(GOPROXY) && go install $(mod)

gqlinit:
	rm -rf graph
	rm -rf gqlgen.yml server.go
	rm -rf
	go clean -cache -modcache -i -r
	go get github.com/99designs/gqlgen
	gqlgen init
	go mod tidy

gqlgen:
	gqlgen generate

### generator code
gen:
	protoc ./proto/*.proto --go_out=.

protoGen:
	protoc \
-I${PROTO_SRC} \
--go_out ${OUT_DIR} --go_opt paths=source_relative \
--go-grpc_out ${OUT_DIR} --go-grpc_opt paths=source_relative \
--grpc-gateway_out ${OUT_DIR} --grpc-gateway_opt paths=source_relative \
${PROTO_SRC}/*.proto
	echo ''
	@ls -l ${OUT_DIR}/

### 下载模块
get:
	go get

### 启动开发程序
# make run e=development 
run:
	rm -rf gin-bin >/dev/null 2>&1
	TZ=$(zone) GIN_ENV=$(e) WORK_DIR=$(PWD) go run .

### 启动调试程序, 当代码变化时自动重启
# make dev
dev:
	TZ=$(zone) GIN_ENV=$(e) GIN_MAPS_TIMESTAMP=$(GIN_MAPS_TIMESTAMP) gin -i --appPort 8080 --port 3000 run main.go

### lint
lint:
	golangci-lint run

### 构建跨平台的可执行程序
Built1:
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GOPROXY) && go build -buildvcs -ldflags "-X main.Name=$(APP_NAME) -X main.Author=$(COMMITER) -X main.Commit=$(GIT_HASH) -X main.Time=$(TIME)" -o ./dist/$(APP_NAME)-darwin-amd64 ./main.go

Built2:
	env GOOS=linux  GOARCH=amd64 CGO_ENABLED=1 $(GOPROXY) && go build -buildvcs -ldflags "-X main.Name=$(APP_NAME) -X main.Author=$(COMMITER) -X main.Commit=$(GIT_HASH) -X main.Time=$(TIME)" -o ./dist/$(APP_NAME)-linux-amd64 ./main.go

Build:
	docker run --platform linux/amd64 --rm -it -v $(PWD):/app:rw --name build_$(APP_NAME) chunhui2001/ubuntu_20.04_dev:golang_1.25 /bin/bash -c 'cd /app && make -f Makefile install Built2' -m 4g

### 通过容器启动
up: rm
	docker-compose -f docker-compose.yml up -d

serve:
	TZ=$(zone) GIN_ENV=$(e) WORK_DIR=$(PWD) ./dist/$(APP_NAME)-darwin-amd64

### 1 = stdout = normal output of a command
### 2 = stderr = error output of a command
### 0 = stdin = input to a command (this isn't usefull for redirecting, more for logging)
# make -i newtag tag=1.1
newtag:
	git tag -d $(tag) >/dev/null 2>&1
	git push --delete origin $(tag) >/dev/null 2>&1
	git tag $(tag)
	git tag -l
	git push origin $(tag)

### 查看程序日志
logs:
	docker logs -f --tail 1000 $(APP_NAME)

### 删除程序容器
rm:
	docker rm -f $(APP_NAME) >/dev/null 2>&1

privateKey:
	@# Key considerations for algorithm "RSA" ≥ 2048-bit
	openssl genrsa -out server.key 2048
	@# Key considerations for algorithm "ECDSA" (X25519 || ≥ secp384r1)
	@# https://safecurves.cr.yp.to/
	@# List ECDSA the supported curves (openssl ecparam -list_curves)
	@#openssl ecparam -genkey -name secp384r1 -out server.key

publicKey:
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

tls:
	openssl s_client -connect 127.0.0.1:8443

### 删除所有缓存的依赖包
# clear modcache
clear:
	go clean --modcache
	rm -rf `go env GOPATH`/bin/$(APP_NAME)
	@#rm -rf `go env GOPATH`/bin/*
	rm -rf dist gin-bin
	docker image prune -a -f

# 随机密码
passwd:
	head -c12 < /dev/random | base64
	@#head -c12 < /dev/urandom | base64

# make ngrok
#ngrok:
#	ngrok start --config ./ngrok.yml $(APP_NAME)

### 性能测试
# make load n=10000 p=info
load:
	@#h2load -n$(n) -c100 -m10 --h1 "http://localhost:4000/$(p)"
	ab -n 10000 -c 10 "http://localhost:8080/info_cache"


# https://diff2html.xyz/
# npm install -g diff2html-cli
diff:
	git diff -U999999 origin/master | diff2html -i stdin -s side -F diff.html
	@#git diff origin/master | diff2html -i stdin -s side -F diff.html

# 👉 ^ 表示它的上一个父提交。
# 这相当于“查看某个提交到底改了什么”。
diff2:
	@#git diff -U999999 $(c1)^ $(c1) | diff2html -i stdin -s side -F diff.html
	git diff $(c1)^ $(c1) | diff2html -i stdin -s side -F diff.html

diff3:
	git diff $(c1) $(c2) | diff2html -i stdin -s side -F diff.html
