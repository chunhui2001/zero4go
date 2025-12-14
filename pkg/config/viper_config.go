package config

import (
	"log"

	"github.com/spf13/viper"
)

func ViperConfig() *viper.Viper {
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
			log.Printf("Viper Configuration load error: ConfigPath=%s, file=%s, Error=%v", configRoot, file, err)

			return nil
		}

		log.Printf("Viper Configuration loaded: ConfigPath=%s, file=%s", configRoot, file)

		return v
	}

	for _, name := range filenames {
		if ll := f(name, v.AllSettings()); ll != nil {
			v = ll
		}
	}

	return v
}
