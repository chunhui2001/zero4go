package gredis

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"time"

	"github.com/redis/go-redis/v9"

	. "github.com/chunhui2001/zero4go/pkg/logs"
)

var (
	RedisClient redis.UniversalClient
	ctx         context.Context
)

type RedisConf struct {
	Mode           string        `mapstructure:"REDIS_MODE"` // disable, single, sentinel, cluster
	Host           string        `mapstructure:"REDIS_HOST"`
	Addrs          string        `mapstructure:"REDIS_ADDRS"`
	MasterName     string        `mapstructure:"REDIS_MASTER_NAME"`
	Passwd         string        `mapstructure:"REDIS_PASSWORD"`
	Db             int           `mapstructure:"REDIS_DATABASE"`
	MaxIdle        int           `mapstructure:"REDIS_MAX_IDLE"`
	MaxActive      int           `mapstructure:"REDIS_MAX_ACTIVE"`
	IdleTimeout    time.Duration `mapstructure:"REDIS_IDLE_TIMEOUT"`
	RouteByLatency bool          `mapstructure:"REDIS_ROUTE_BY_LATENCY"`
	RouteRandomly  bool          `mapstructure:"REDIS_ROUTE_RANDOMLY"`
	//SubChannels    string        `mapstructure:"REDIS_SUB_CHANNELS"`
	//PrintMessage bool `mapstructure:"REDIS_MESSAGE_PRINT"`
}

var Settings = &RedisConf{
	Mode:           "disable",
	Host:           "127.0.0.1:6379",
	Addrs:          "127.0.0.1:6381,127.0.0.1:6382,127.0.0.1:6383,127.0.0.4:6384,127.0.0.1:6385",
	MasterName:     "redis_master",
	Passwd:         "Cc",
	Db:             0,
	MaxIdle:        30,
	MaxActive:      30,
	IdleTimeout:    time.Second * 20, // 20
	RouteByLatency: false,
	RouteRandomly:  true,
}

func (r *RedisConf) ServerAddrs() string {
	if r.Mode == "standalone" {
		return r.Host
	}

	if r.Mode == "sentinel" {
		return r.Addrs
	}

	if r.Mode == "cluster" {
		return r.Addrs
	}

	return ""
}

func Init() {

	if Settings.Mode == "disable" {
		Log.Infof("Redis-Initialized-Disabled: val=%s", Settings.Mode)

		return
	}

	ctx = context.Background()

	if Settings.Mode == "standalone" {
		if Settings.Passwd != "" {
			RedisClient = redis.NewClient(&redis.Options{
				Addr:     Settings.Host,
				DB:       Settings.Db,
				Password: Settings.Passwd,
			})
		} else {
			RedisClient = redis.NewClient(&redis.Options{
				Addr: Settings.Host,
				DB:   Settings.Db,
			})
		}

		Ping()

		return
	}

	var addrs = strings.Split(Settings.Addrs, ",")

	if Settings.Mode == "sentinel" {
		if Settings.RouteByLatency || Settings.RouteRandomly {
			if Settings.Passwd != "" {
				RedisClient = redis.NewFailoverClusterClient(&redis.FailoverOptions{
					MasterName:     Settings.MasterName,
					SentinelAddrs:  addrs,
					RouteByLatency: Settings.RouteByLatency,
					RouteRandomly:  Settings.RouteRandomly,
					Password:       Settings.Passwd,
				})
			} else {
				RedisClient = redis.NewFailoverClusterClient(&redis.FailoverOptions{
					MasterName:     Settings.MasterName,
					SentinelAddrs:  addrs,
					RouteByLatency: Settings.RouteByLatency,
					RouteRandomly:  Settings.RouteRandomly,
				})
			}
		} else {
			if Settings.Passwd != "" {
				RedisClient = redis.NewFailoverClient(&redis.FailoverOptions{
					MasterName:    Settings.MasterName,
					SentinelAddrs: addrs,
					Password:      Settings.Passwd,
				})
			} else {
				RedisClient = redis.NewFailoverClient(&redis.FailoverOptions{
					MasterName:    Settings.MasterName,
					SentinelAddrs: addrs,
				})
			}
		}

		Ping()

		return
	}

	if Settings.Mode == "cluster" {
		if Settings.RouteByLatency || Settings.RouteRandomly {
			RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:          addrs,
				RouteByLatency: Settings.RouteByLatency,
				RouteRandomly:  Settings.RouteRandomly,
			})
		} else {
			RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: addrs,
			})
		}

		Ping()

		return
	}
}

func Ping() {

	var serverInfo = "N/a"

	if Settings.Mode == "sentinel" {
		serverInfo = fmt.Sprintf("Mode=%s, MasterName=%s, ServerAddrs=%s", Settings.Mode, Settings.MasterName, Settings.ServerAddrs())
	} else if Settings.Mode == "standalone" {
		serverInfo = fmt.Sprintf("Mode=%s, ServerAddrs=%s, DB=%d", Settings.Mode, Settings.ServerAddrs(), Settings.Db)
	} else if Settings.Mode == "cluster" {
		serverInfo = fmt.Sprintf("Mode=%s, ServerAddrs=%s", Settings.Mode, Settings.ServerAddrs())
	}

	info, _ := RedisClient.Info(context.Background(), "server").Result()

	var serverVersion = "N/a"
	var redisVersionRE = regexp.MustCompile(`redis_version:(.+)`)

	match := redisVersionRE.FindAllStringSubmatch(info, -1)

	if len(match) < 1 {
		// could not extract redis version
		// ...
	} else {
		serverVersion = strings.TrimSpace(match[0][1])
	}

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		Log.Error(fmt.Sprintf("Redis-Initialized-Failed: ServerVersion=%s, %s, Error=%v", serverVersion, serverInfo, err))

		return
	}

	Log.Info(fmt.Sprintf("Redis-Initialized-Succeed: ServerVersion=%s, %s", serverVersion, serverInfo))
}
