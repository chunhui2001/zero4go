package config

import (
	"log"
	"os"

	"github.com/chunhui2001/zero4go/pkg/stdout"

	"github.com/alecthomas/kong"
	"github.com/spf13/viper"

	"github.com/chunhui2001/zero4go/pkg/gkafka"
	"github.com/chunhui2001/zero4go/pkg/gredis"
	"github.com/chunhui2001/zero4go/pkg/http_client"
	"github.com/chunhui2001/zero4go/pkg/logs"
	"github.com/chunhui2001/zero4go/pkg/mysqlg"
)

type CLI struct {
	Env    string `help:"environment name" short:"e"`
	Config string `help:"config file path" short:"c"`
}

func (c *CLI) Run() error {
	log.Printf("root command running env=%s", c.Env)

	return nil
}

type AppConf struct {
	Env      string `mapstructure:"GIN_ENV"`
	AppName  string `mapstructure:"APP_NAME"`
	AppPort  string `mapstructure:"APP_PORT"`
	TimeZone string `mapstructure:"APP_TIMEZONE"`
	NodeId   int64  `mapstructure:"NODE_ID"`
	RpcPort  string `mapstructure:"RPC_PORT"`
}

var AppSetting = &AppConf{
	Env:      "local",
	AppName:  "zero4go",
	AppPort:  "0.0.0.0:8080",
	TimeZone: "Asia/Shanghai",
	NodeId:   0,
	RpcPort:  "0.0.0.0:50051",
}

var envName = os.Getenv("GIN_ENV")

var configFolder = "config"
var envDefault = ".env"
var viperConfig *viper.Viper

var cli CLI

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
}

func cliResolver() {

	if cli.Config != "" {
		configFolder = cli.Config
	}

	if cli.Env != "" {
		envName = cli.Env
	}
}

func OnLoad() {
	// 解析命令行参数
	cliResolver()

	// 读取配置
	if v1 := readConfig(); v1 != nil {
		viperConfig = v1

		if err := v1.Unmarshal(AppSetting); err != nil {
			log.Printf("viper parse AppConf error: configRoot=%s, errorMessage=%v", configRoot(), err)

			os.Exit(3)
		}

		if err := v1.Unmarshal(logs.LogSetting); err != nil {
			log.Printf("viper parse LogConf error: configRoot=%s, errorMessage=%v", configRoot(), err)

			os.Exit(3)
		}

		if err := v1.Unmarshal(gkafka.KafkaSetting); err != nil {
			log.Printf("viper parse KafkaConf error: configRoot=%s, errorMessage=%v", configRoot(), err)

			os.Exit(3)
		}

		if err := v1.Unmarshal(http_client.HttpSetting); err != nil {
			log.Printf("viper parse HttpConf error: configRoot=%s, errorMessage=%v", configRoot(), err)

			os.Exit(3)
		}

		if err := v1.Unmarshal(gredis.ReidsSetting); err != nil {
			log.Printf("viper parse RedisConf error: configRoot=%s, errorMessage=%v", configRoot(), err)

			os.Exit(3)
		}

		if err := v1.Unmarshal(mysqlg.MysqlSetting); err != nil {
			log.Printf("viper parse MySqlConf error: configRoot=%s, errorMessage=%v", configRoot(), err)

			os.Exit(3)
		}
	}
}

// GetEnv returns an environment variable or a default value if not present
func GetEnv(key, defaultValue string) string {

	value := os.Getenv(key)

	if value != "" {
		return value
	}

	value = viperConfig.GetString(key)

	if value != "" {
		return value
	}

	return defaultValue
}
