package bootstrap

import (
	"github.com/chunhui2001/zero4go/pkg/config"
	"github.com/chunhui2001/zero4go/pkg/logs"
)

func init() {
	// 读取配置文件
	config.OnLoad()
	// 初始化日志
	logs.InitLog()
}
