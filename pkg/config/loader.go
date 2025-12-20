package config

import (
	"path/filepath"

	"github.com/chunhui2001/zero4go/pkg/cli"
	"github.com/spf13/viper"

	"github.com/chunhui2001/zero4go/pkg/utils"
)

func configRoot() string {
	return filepath.Join(utils.RootDir(), cli.Cli.Config)
}

func readConfig() *viper.Viper {
	// 读取 viper 配置
	v := ViperConfig()

	// 读取 apollo 配置
	ReadApolloConfig(v)

	return v
}

func configFilesEnv() []string {
	var configFiles []string

	configFiles = append(configFiles, ".env")

	if cli.Cli.Env != "" {
		configFiles = append(configFiles, ".env."+cli.Cli.Env)
	}

	return configFiles
}

func configFilesApplication() []string {
	var _fil1 = "application.properties"
	var _fil3 = "application.yml"
	var _fil4 = "application.yaml"

	var configFiles []string

	configFiles = append(configFiles, _fil1)

	if cli.Cli.Env != "" {
		configFiles = append(configFiles, "application-"+cli.Cli.Env+".properties")
	}

	configFiles = append(configFiles, _fil3, _fil4)

	if cli.Cli.Env != "" {
		configFiles = append(configFiles, "application-"+cli.Cli.Env+".yml", "application-"+cli.Cli.Env+".yaml")
	}

	return configFiles
}
