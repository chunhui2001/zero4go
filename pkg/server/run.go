package server

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/chunhui2001/zero4go/pkg/config"
	"github.com/chunhui2001/zero4go/pkg/favicon"
	"github.com/chunhui2001/zero4go/pkg/interceptors"
	"github.com/chunhui2001/zero4go/pkg/utils"

	_ "github.com/chunhui2001/zero4go/pkg/bootstrap"

	. "github.com/chunhui2001/zero4go/pkg"
	. "github.com/chunhui2001/zero4go/pkg/logs"

	"github.com/chunhui2001/zero4go/graph"
	"github.com/chunhui2001/zero4go/pkg/graphql"
	"github.com/chunhui2001/zero4go/rpc"
	pb "github.com/chunhui2001/zero4go/rpc/gen"
)

func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler(endpoint string) gin.HandlerFunc {
	h := graphql.Playground("GraphQL playground", endpoint)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func Setup(f func(*Application)) *Application {
	gin.SetMode(gin.ReleaseMode)

	r := &Application{Engine: gin.New()}

	r.Use(gin.Recovery())

	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".pdf", ".mp4", ".ico"})))
	r.Use(static.Serve("/RichMedias", static.LocalFile(filepath.Join(utils.RootDir(), "./static"), false)))
	r.Use(favicon.Favicon())

	r.Use(interceptors.AccessLog("/favicon.ico", "/static"))

	if config.AppSetting.GraphQLEnable {
		r.POST(config.AppSetting.GraphQLServerURI, graphqlHandler())
		r.GET(config.AppSetting.GraphQLPlaygroundURI, playgroundHandler(config.AppSetting.GraphQLServerURI))
	}

	// routers
	r.GET("/info", func(c *RequestContext) {
		c.Text("Yeah, your server is running.")
	})

	// customer http router
	f(r)

	r.NoRoute(Wrap(func(c *RequestContext) {
		if c.Request.RequestURI == "/favicon.ico" {
			c.Next()
		} else {
			c.Text("404 page not found.")
		}
	})...)

	return r
}

func (a *Application) Run(f func(*grpc.Server)) {

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryLoggingInterceptor),
		grpc.StreamInterceptor(interceptors.StreamLoggingInterceptor),
	)

	pb.RegisterGreeterServer(grpcServer, &rpc.GreeterServer{})

	// customer grpc service
	f(grpcServer)

	// ⭐ 开启 Reflection
	reflection.Register(grpcServer)

	srv := &http.Server{
		Addr:        config.AppSetting.AppPort,
		Handler:     a,
		IdleTimeout: 5 * time.Second,
	}

	l, err := net.Listen("tcp", config.AppSetting.AppPort)

	if err != nil {
		Log.Info("Application Run Failed: ErrorMessage=" + err.Error())
		os.Exit(1)

		return
	}

	lis, err := net.Listen("tcp", config.AppSetting.RpcPort)

	if err != nil {
		Log.Fatal(err.Error())
	}

	AddShutDownHook(func() {
		Log.Info("shutting down server")

		grpcServer.Stop()

		if err := srv.Shutdown(context.Background()); err != nil {
			Log.Info("shutting down server-err")
		} else {
			Log.Info("shutting down server-done")
		}
	})

	// gRPC server
	go func() {
		Log.Infof("gRPC server listening on %s", config.AppSetting.RpcPort)

		if err := grpcServer.Serve(lis); err != nil {
			Log.Fatal(err.Error())
		}
	}()

	// httpserver
	go func() {
		Log.Infof("RootDir: val=%s", utils.RootDir())
		Log.Infof("TempDir: val=%s", utils.TempDir())
		Log.Infof("Congratulations! Your server startup successfully, Listening and serving HTTP on %s", config.AppSetting.AppPort)

		if err := srv.Serve(l); err != nil {
			Log.Info(err.Error())
		}
	}()

	WaitShutDown()
}
