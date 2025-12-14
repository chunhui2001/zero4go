package config

import (
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/chunhui2001/zero4go/pkg/utils"
)

func configRoot() string {
	return filepath.Join(utils.RootDir(), configFolder)
}

func readConfig() *viper.Viper {
	// 读取 viper 配置
	v := ViperConfig()

	// 读取 apollo 配置
	ReadApolloConfig(v)

	return v
}

func configFiles() []string {
	var _enfFile = ".env." + envName

	//var _applicationPropsFile = "application.properties"
	//var _applicationPropsFileEnv = "application-" + env + ".properties"
	//
	//var _applicationYamlFile = "application.properties"
	//var _applicationYamlsFileEnv = "application-" + env + ".properties"

	var configFiles []string

	configFiles = append(configFiles, envDefault)

	//if exists, _ := utils.FileExists(filepath.Join(configRoot(), _enfFile)); exists {

	configFiles = append(configFiles, _enfFile)
	//}

	//if exists, _ := utils.FileExists(filepath.Join(configRoot(), _enfFile)); exists {
	//
	//	configFiles = append(configFiles, _enfFile)
	//}

	return configFiles
}
