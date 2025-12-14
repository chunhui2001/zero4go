package bootstrap

import (
	"log"

	"github.com/alecthomas/kong"
	"github.com/chunhui2001/zero4go/pkg/build_info"
	"github.com/chunhui2001/zero4go/pkg/config"
	"github.com/chunhui2001/zero4go/pkg/elasticsearch"
	"github.com/chunhui2001/zero4go/pkg/elasticsearch_openes"
	"github.com/chunhui2001/zero4go/pkg/gkafka"
	"github.com/chunhui2001/zero4go/pkg/gredis"
	"github.com/chunhui2001/zero4go/pkg/gsql"
	"github.com/chunhui2001/zero4go/pkg/gzook"
	"github.com/chunhui2001/zero4go/pkg/http_client"
	"github.com/chunhui2001/zero4go/pkg/logs"
	"github.com/chunhui2001/zero4go/pkg/stdout"
	"github.com/chunhui2001/zero4go/pkg/x"
)

var cli CLI

type CLI struct {
	Env    string `help:"environment name" short:"e"`
	Config string `help:"config file path" short:"c"`
}

func (c *CLI) Run() error {
	log.Println(build_info.INFO.Info())
	x.Info()
	log.Printf("root command running: env=%s", c.Env)

	return nil
}

func cliResolver() {

	if cli.Config != "" {
		config.ConfigurationFolder = cli.Config
	}

	if cli.Env != "" {
		config.EnvName = cli.Env
	}
}

func init() {
	// 设置控制台日志输出
	stdout.SetOutputWriter()

	ctx := kong.Parse(&cli,
		kong.Name("zero4go"),
		kong.Description("Rust clap style CLI in Go using kong."),
	)

	// Run subcommand
	if err := ctx.Run(&cli); err != nil {
		log.Printf("error=%v", err)
	}

	// 设置命令行参数
	cliResolver()

	config.OnLoad()

	logs.InitLog()

	http_client.Init()

	gkafka.Init()
	gredis.Init()
	gsql.Init()
	elasticsearch.Init()
	elasticsearch_openes.Init()
	gzook.Init()
}
