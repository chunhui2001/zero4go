package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/chunhui2001/zero4go/pkg/gkafka"
	"github.com/chunhui2001/zero4go/pkg/gredis"
	"github.com/chunhui2001/zero4go/pkg/gsql"
	"github.com/chunhui2001/zero4go/pkg/gzook"
	"github.com/chunhui2001/zero4go/pkg/http_client"
	"github.com/chunhui2001/zero4go/pkg/logs"
	"github.com/chunhui2001/zero4go/pkg/search_elastic"
	"github.com/chunhui2001/zero4go/pkg/search_openes"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type AppConf struct {
	Env                  string `mapstructure:"GIN_ENV"`
	AppName              string `mapstructure:"APP_NAME"`
	AppPort              string `mapstructure:"APP_PORT"`
	TimeZone             string `mapstructure:"APP_TIMEZONE"`
	NodeId               int64  `mapstructure:"NODE_ID"`
	RpcPort              string `mapstructure:"RPC_PORT"`
	GraphQLEnable        bool   `mapstructure:"GRAPHQL_ENABLE"`
	GraphQLServerURI     string `mapstructure:"GRAPHQL_SERVER_URI"`
	GraphQLPlaygroundURI string `mapstructure:"GRAPHQL_PLAYGROUND_URI"`
}

var AppSetting = &AppConf{
	Env:      "local",
	AppName:  "zero4go",
	AppPort:  "0.0.0.0:8080",
	TimeZone: "Asia/Shanghai",
	NodeId:   0,
	RpcPort:  "0.0.0.0:50051",
}

var viperConfig *viper.Viper

var onChangeFuncs []func()

// RegisterOnChange 注册回调函数
func RegisterOnChange(f func()) {
	onChangeFuncs = append(onChangeFuncs, f)
}

// TriggerOnChange 主动触发回调
func TriggerOnChange() {
	for _, f := range onChangeFuncs {
		f()
	}
}

func init() {
	//RegisterOnChange(loadConfig)

	// 5️⃣ 监听文件变化
	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config changed:", e.Name)
	})

	loadConfig()

	// TriggerOnChange()
}

func Viper() *viper.Viper {
	return viperConfig
}

func loadConfig() {
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

		if err := v1.Unmarshal(gkafka.Settings); err != nil {
			log.Printf("viper parse KafkaConf error: configRoot=%s, errorMessage=%v", configRoot(), err)

			os.Exit(3)
		}

		if err := v1.Unmarshal(http_client.Settings); err != nil {
			log.Printf("viper parse HttpConf error: configRoot=%s, errorMessage=%v", configRoot(), err)

			os.Exit(3)
		}

		if err := v1.Unmarshal(gredis.Settings); err != nil {
			log.Printf("viper parse RedisConf error: configRoot=%s, errorMessage=%v", configRoot(), err)

			os.Exit(3)
		}

		if err := v1.Unmarshal(gsql.Settings); err != nil {
			log.Printf("viper parse MySQLConf error: configRoot=%s, errorMessage=%v", configRoot(), err)

			os.Exit(3)
		}

		if err := v1.Unmarshal(search_elastic.Settings); err != nil {
			log.Printf("viper parse ESConf error: configRoot=%s, errorMessage=%v", configRoot(), err)

			os.Exit(3)
		}

		if err := v1.Unmarshal(search_openes.Settings); err != nil {
			log.Printf("viper parse ESConf error: configRoot=%s, errorMessage=%v", configRoot(), err)

			os.Exit(3)
		}

		if err := v1.Unmarshal(gzook.Settings); err != nil {
			log.Printf("viper parse ZookConf error: configRoot=%s, errorMessage=%v", configRoot(), err)

			os.Exit(3)
		}

		if _c := v1.Get("MySQLDataSource"); _c != nil {
			var raw = make([]map[string]any, 0)

			for _, c := range _c.([]interface{}) {

				raw = append(raw, c.(map[string]any))
			}

			var _databases []gsql.MySQLConf

			for _, m := range raw {
				_databases = append(_databases, gsql.MySQLConf{
					Name:     m["MYSQL_NAME"].(string),
					Enable:   m["MYSQL_ENABLE"].(bool),
					Server:   m["MYSQL_SERVER"].(string),
					Database: m["MYSQL_DATABASE"].(string),
					User:     m["MYSQL_USER_NAME"].(string),
					Passwd:   m["MYSQL_PASSWD"].(string),
					Location: m["MYSQL_MAPPER_LOCATION"].(string),
				})
			}

			gsql.Databases = _databases
		}
	}
}

func GetConfig(key string) any {
	var keyPath = strings.Split(strings.ToLower(key), ".")

	for i, k := range keyPath {
		if viperConfig.Get(k) != nil {
			var val = viperConfig.Get(k)

			switch val.(type) {
			case float64, float32:
				return val
			case string:
				return val
			case bool:
				return val
			case byte:
				return val
			case []uint8:
				return val
			case map[string]any:
				if i == len(keyPath)-1 {
					return viperConfig.Get(k)
				}

				return GetByPathAny(strings.Join(keyPath[i+1:], "."), val)
			default:
				if i == len(keyPath)-1 {
					return viperConfig.Get(k)
				}

				return GetConfig(strings.Join(keyPath[i+1:], "."))
			}
		}
	}

	return viperConfig.Get(key)
}

func GetByPathAny(path string, v any) any {
	parts := strings.Split(path, ".")

	cur := v

	for _, p := range parts {
		switch node := cur.(type) {
		case map[string]any:
			val, ok := node[p]

			if !ok {
				return nil
			}

			cur = val

		case []any:
			idx, err := strconv.Atoi(p)

			if err != nil || idx < 0 || idx >= len(node) {
				return nil
			}

			cur = node[idx]

		default:
			return nil
		}
	}

	return cur
}

func Configurations() map[string]any {
	return viperConfig.AllSettings()
}
