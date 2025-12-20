package boot

import (
	_ "github.com/chunhui2001/zero4go/pkg/config"
	"github.com/chunhui2001/zero4go/pkg/gkafka"
	"github.com/chunhui2001/zero4go/pkg/gredis"
	"github.com/chunhui2001/zero4go/pkg/gsql"
	"github.com/chunhui2001/zero4go/pkg/gzook"
	"github.com/chunhui2001/zero4go/pkg/http_client"
	"github.com/chunhui2001/zero4go/pkg/logs"
	"github.com/chunhui2001/zero4go/pkg/search_elastic"
	"github.com/chunhui2001/zero4go/pkg/search_openes"
)

func init() {
	// 初始化日志
	logs.InitLog()

	http_client.Init()
	gkafka.Init()
	gredis.Init()
	gsql.Init()
	search_elastic.Init()
	search_openes.Init()
	gzook.Init()
}
