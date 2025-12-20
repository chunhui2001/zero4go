package cli

import (
	"log"

	"github.com/alecthomas/kong"
	"github.com/chunhui2001/zero4go/pkg/build_info"
	"github.com/chunhui2001/zero4go/pkg/stdout"
	"github.com/chunhui2001/zero4go/pkg/x"
)

var Cli CLI

type CLI struct {
	Env    string `help:"environment name: default:'local'" short:"e"`
	Config string `help:"config file folder, default: 'config'" short:"c"`

	ApolloServer    string `help:"apollo config server" short:"a"`
	ApolloName      string `help:"apollo application name" short:"n"`
	ApolloProfile   string `help:"apollo profile name" short:"p"`
	ApolloNamespace string `help:"apollo namespace, default: 'application.properties,application.yaml'" short:"s"`
}

func (c *CLI) Run() error {
	build_info.INFO.Info()
	x.Info()

	// 设置命令行参数
	cliResolver()

	return nil
}

func init() {
	// 设置控制台日志输出
	stdout.SetOutputWriter()

	ctx := kong.Parse(&Cli,
		kong.Name("zero4go"),
		kong.Description("Rust clap style CLI in Go using kong."),
	)

	// Run subcommand
	if err := ctx.Run(&Cli); err != nil {
		log.Printf("error=%v", err)
	}
}

func cliResolver() {
	if Cli.Env == "" {
		Cli.Env = "local"
	}

	if Cli.Config == "" {
		Cli.Config = "config"
	}

	if Cli.ApolloProfile != "" {
		Cli.Env = Cli.ApolloProfile
	}

	if Cli.ApolloNamespace == "" {
		Cli.ApolloNamespace = "application.properties,application.yaml"
	}

	log.Printf("root command running: env=%s", Cli.Env)
}
