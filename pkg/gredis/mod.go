package gredis

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"time"

	"github.com/go-redis/redis/v8"

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

var ReidsSetting = &RedisConf{
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

	if ReidsSetting.Mode == "disable" {
		Log.Infof("Init redis mode: val=%s", ReidsSetting.Mode)

		return
	}

	ctx = context.Background()

	if ReidsSetting.Mode == "standalone" {
		if ReidsSetting.Passwd != "" {
			RedisClient = redis.NewClient(&redis.Options{
				Addr:     ReidsSetting.Host,
				DB:       ReidsSetting.Db,
				Password: ReidsSetting.Passwd,
			})
		} else {
			RedisClient = redis.NewClient(&redis.Options{
				Addr: ReidsSetting.Host,
				DB:   ReidsSetting.Db,
			})
		}

		Ping()

		return
	}

	var addrs = strings.Split(ReidsSetting.Addrs, ",")

	if ReidsSetting.Mode == "sentinel" {
		if ReidsSetting.RouteByLatency || ReidsSetting.RouteRandomly {
			if ReidsSetting.Passwd != "" {
				RedisClient = redis.NewFailoverClusterClient(&redis.FailoverOptions{
					MasterName:     ReidsSetting.MasterName,
					SentinelAddrs:  addrs,
					RouteByLatency: ReidsSetting.RouteByLatency,
					RouteRandomly:  ReidsSetting.RouteRandomly,
					Password:       ReidsSetting.Passwd,
				})
			} else {
				RedisClient = redis.NewFailoverClusterClient(&redis.FailoverOptions{
					MasterName:     ReidsSetting.MasterName,
					SentinelAddrs:  addrs,
					RouteByLatency: ReidsSetting.RouteByLatency,
					RouteRandomly:  ReidsSetting.RouteRandomly,
				})
			}
		} else {
			if ReidsSetting.Passwd != "" {
				RedisClient = redis.NewFailoverClient(&redis.FailoverOptions{
					MasterName:    ReidsSetting.MasterName,
					SentinelAddrs: addrs,
					Password:      ReidsSetting.Passwd,
				})
			} else {
				RedisClient = redis.NewFailoverClient(&redis.FailoverOptions{
					MasterName:    ReidsSetting.MasterName,
					SentinelAddrs: addrs,
				})
			}
		}

		Ping()

		return
	}

	if ReidsSetting.Mode == "cluster" {
		if ReidsSetting.RouteByLatency || ReidsSetting.RouteRandomly {
			RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:          addrs,
				RouteByLatency: ReidsSetting.RouteByLatency,
				RouteRandomly:  ReidsSetting.RouteRandomly,
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

	if ReidsSetting.Mode == "sentinel" {
		serverInfo = fmt.Sprintf("Mode=%s, MasterName=%s, ServerAddrs=%s", ReidsSetting.Mode, ReidsSetting.MasterName, ReidsSetting.ServerAddrs())
	} else if ReidsSetting.Mode == "standalone" {
		serverInfo = fmt.Sprintf("Mode=%s, ServerAddrs=%s, DB=%d", ReidsSetting.Mode, ReidsSetting.ServerAddrs(), ReidsSetting.Db)
	} else if ReidsSetting.Mode == "cluster" {
		serverInfo = fmt.Sprintf("Mode=%s, ServerAddrs=%s", ReidsSetting.Mode, ReidsSetting.ServerAddrs())
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
		Log.Error(fmt.Sprintf("Redis-Client-Connect-Failed: ServerVersion=%s, %s, Error=%v", serverVersion, serverInfo, err))

		return
	}

	Log.Info(fmt.Sprintf("Redis-Client-Connect-Succeed: ServerVersion=%s, %s", serverVersion, serverInfo))
}
