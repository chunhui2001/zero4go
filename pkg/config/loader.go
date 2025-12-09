package config

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/chunhui2001/zero4go/pkg/utils"
)

func configRoot() string {
	return filepath.Join(utils.RootDir(), configFolder)
}

func readConfig() *viper.Viper {
	var configRoot = configRoot()
	var filenames = configFiles()

	v := viper.New()

	var f = func(file string, defaultMaps map[string]interface{}) *viper.Viper {
		v := viper.New()

		for key, value := range defaultMaps {
			v.SetDefault(key, value)
		}

		v.SetConfigName(file)
		v.SetConfigType("env")
		// v.AddConfigPath("/etc/appname/")   // path to look for the config file in
		// v.AddConfigPath("$(home)/.env") // call multiple times to add many search paths
		v.AddConfigPath(configRoot)
		v.AutomaticEnv() // 将读取当前目录下的 .env 配置文件或"环境变量", .env 优先级最高

		err := v.ReadInConfig()

		if err != nil {
			log.Printf("viper Configuration load error: ConfigPath=%s, file=%s, Error=%v", configRoot, file, err)

			return nil
		}

		log.Printf("viper Configuration load succeed: ConfigPath=%s, file=%s", configRoot, file)

		return v
	}

	for _, fname := range filenames {
		if ll := f(fname, v.AllSettings()); ll != nil {
			v = ll
		}
	}

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
