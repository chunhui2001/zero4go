package config

import (
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/chunhui2001/zero4go/pkg/utils"
)

func configRoot() string {
	return filepath.Join(utils.RootDir(), ConfigurationFolder)
}

func readConfig() *viper.Viper {
	// 读取 viper 配置
	v := ViperConfig()

	// 读取 apollo 配置
	ReadApolloConfig(v)

	return v
}

func configFilesEnv() []string {
	var _enfFile = ".env." + EnvName

	var configFiles []string

	configFiles = append(configFiles, envDefault)
	configFiles = append(configFiles, _enfFile)
	
	return configFiles
}

func configFilesAppliction() []string {
	var _fil1 = "application.properties"
	var _fil2 = "application-" + EnvName + ".properties"

	var _fil3 = "application.yml"
	var _fil4 = "application.yaml"
	var _fil5 = "application-" + EnvName + ".yml"
	var _fil6 = "application-" + EnvName + ".yaml"

	var configFiles []string

	configFiles = append(configFiles, _fil1, _fil2, _fil3, _fil4, _fil5, _fil6)

	return configFiles
}
