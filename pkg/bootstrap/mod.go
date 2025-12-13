package bootstrap

import (
	"github.com/chunhui2001/zero4go/pkg/config"
	"github.com/chunhui2001/zero4go/pkg/elasticsearch"
	"github.com/chunhui2001/zero4go/pkg/elasticsearch_openes"
	"github.com/chunhui2001/zero4go/pkg/gkafka"
	"github.com/chunhui2001/zero4go/pkg/gredis"
	"github.com/chunhui2001/zero4go/pkg/gsql"
	"github.com/chunhui2001/zero4go/pkg/gzook"
	"github.com/chunhui2001/zero4go/pkg/http_client"
	"github.com/chunhui2001/zero4go/pkg/logs"
)

func init() {
	config.OnLoad()

	logs.InitLog()
	http_client.Init()

	gkafka.Init()
	gredis.Init()
	gsql.Init()
	elasticsearch.Init()
	elasticsearch_openes.Init()
	gzook.Init()
}
