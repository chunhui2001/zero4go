package bootstrap

import (
	"github.com/chunhui2001/zero4go/pkg/config"
	"github.com/chunhui2001/zero4go/pkg/gkafka"
	"github.com/chunhui2001/zero4go/pkg/gredis"
	"github.com/chunhui2001/zero4go/pkg/http_client"
	"github.com/chunhui2001/zero4go/pkg/logs"
	"github.com/chunhui2001/zero4go/pkg/mysqlg"
)

func init() {
	config.OnLoad()

	logs.InitLog()
	http_client.Init()
	gkafka.InitKafka()

	gredis.Init()
	mysqlg.Init()
}
