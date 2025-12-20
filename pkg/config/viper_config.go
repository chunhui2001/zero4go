package config

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/chunhui2001/zero4go/pkg/utils"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func ViperConfig() *viper.Viper {
	var configRoot = configRoot()
	var filenames1 = configFilesEnv()

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

		if err := v.ReadInConfig(); err != nil {
			log.Printf("Viper Configuration load error: ConfigPath=%s, file=%s, Error=%v", configRoot, file, err)

			return nil
		}

		log.Printf("Viper Configuration loaded: ConfigPath=%s, file=%s", configRoot, file)

		return v
	}

	for _, name := range filenames1 {
		if ll := f(name, v.AllSettings()); ll != nil {
			v = ll
		}
	}

	var filenames2 = configFilesApplication()

	ff := func(str string) (key, value string, ok bool) {
		idx := strings.Index(str, "=")

		if idx == -1 {
			return "", "", false
		}

		return str[:idx], str[idx+1:], true
	}

	for _, name := range filenames2 {
		var currFile = filepath.Join(configRoot, "", name)

		if b, _ := utils.FileExists(currFile); !b {
			continue
		}

		log.Printf("Viper Configuration loaded: ConfigPath=%s, file=%s", configRoot, name)

		if strings.HasSuffix(name, ".properties") {
			var fileLines = utils.ReadAllLines(currFile)
			var _map = map[string]interface{}{}

			for _, line := range fileLines {
				if k, v, ok := ff(line); ok {
					_map[k] = strings.TrimSpace(v)
				}
			}

			for key, val := range _map {
				v.SetDefault(strings.TrimSpace(key), strings.TrimSpace(val.(string)))
			}
		} else {
			// yaml
			var yamlBytes = utils.ReadFile(currFile)
			var _map map[string]any

			if err := yaml.Unmarshal(yamlBytes, &_map); err != nil {
				panic(err)
			}

			for key, val := range _map {
				v.SetDefault(strings.TrimSpace(key), val)
			}
		}
	}

	return v
}
