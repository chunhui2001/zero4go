package server

import (
	"context"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/chunhui2001/zero4go/pkg/config"
	"github.com/chunhui2001/zero4go/pkg/favicon"
	"github.com/chunhui2001/zero4go/pkg/middlewares"
	"github.com/chunhui2001/zero4go/pkg/utils"

	_ "github.com/chunhui2001/zero4go/pkg/boot"

	. "github.com/chunhui2001/zero4go/pkg/logs"   //nolint:staticcheck
	. "github.com/chunhui2001/zero4go/pkg/single" //nolint:staticcheck

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

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

func Setup(f func(*Application)) *Application {
	gin.SetMode(gin.ReleaseMode)

	r := &Application{Engine: gin.New()}

	r.HTMLRender = ginview.New(goview.Config{
		Root:      filepath.Join(utils.RootDir(), _SiteConf.Root),
		Extension: _SiteConf.Extension,
		Master:    _SiteConf.Master,
		Partials:  []string{"partials/ad"},
		Funcs: template.FuncMap{
			"timeString": func(b uint32) string {
				return time.Unix(int64(b), 0).Format("2006-01-02T15:04:05Z07:00")
			},
		},
		DisableCache: true,
	})

	rateLimitMiddleWare := ratelimit.RateLimiter(
		ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
			Rate:  time.Second,
			Limit: 10,
		}),
		&ratelimit.Options{
			ErrorHandler: errorHandler,
			KeyFunc: func(c *gin.Context) string {
				return c.ClientIP()
			},
		})

	r.Use(gin.Recovery())
	r.Use(rateLimitMiddleWare)

	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".pdf", ".mp4", ".ico"})))
	r.Use(static.Serve("/RichMedias", static.LocalFile(filepath.Join(utils.RootDir(), "./static"), false)))
	r.Use(favicon.Favicon())

	r.Use(middlewares.AccessLog("/favicon.ico", "/static"))

	if config.AppSetting.GraphQLEnable {
		r.POST(config.AppSetting.GraphQLServerURI, graphqlHandler())
		r.GET(config.AppSetting.GraphQLPlaygroundURI, playgroundHandler(config.AppSetting.GraphQLServerURI))
	}

	// routers
	r.GET("/info", func(c *RequestContext) {
		c.Text("Yeah, your server is running.")
	})

	r.Upstream("/index2", "/index", "http://127.0.0.1:8080")

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
		grpc.UnaryInterceptor(middlewares.UnaryLoggingInterceptor),
		grpc.StreamInterceptor(middlewares.StreamLoggingInterceptor),
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
		Log.Infof("%s: %s", fmt.Sprintf("%-20s", "gRPC Port"), config.AppSetting.RpcPort)

		if err := grpcServer.Serve(lis); err != nil {
			Log.Fatal(err.Error())
		}
	}()

	// httpserver
	go func() {
		name, offset := utils.DateOffsets()

		Log.Infof("%s: %s", fmt.Sprintf("%-20s", "RootDir"), utils.RootDir())
		Log.Infof("%s: %s", fmt.Sprintf("%-20s", "TempDir"), utils.TempDir())
		Log.Infof("%s: %s", fmt.Sprintf("%-20s", "Zone Name"), name)
		Log.Infof("%s: %s", fmt.Sprintf("%-20s", "Datetime Offset"), offset)
		Log.Infof("%s: %s", fmt.Sprintf("%-20s", "Http Address"), config.AppSetting.AppPort)
		Log.Infof("%s! %s", fmt.Sprintf("%-20s", "Congratulations"), "Your server startup and running ~")

		if err := srv.Serve(l); err != nil {
			Log.Info(err.Error())
		}
	}()

	WaitShutDown()
}
