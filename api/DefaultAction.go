package api

import (
	"fmt"

	"github.com/chunhui2001/zero4go/pkg/gkafka"
	"github.com/chunhui2001/zero4go/pkg/gsql"
	"github.com/chunhui2001/zero4go/pkg/logs"
	. "github.com/chunhui2001/zero4go/pkg/server"
	"github.com/chunhui2001/zero4go/pkg/utils"
)

func ChangeLogHandler(c *RequestContext) {
	body := logs.LogConf{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Failc(400, "bind err: %v", err)

		return
	}

	logs.OnChange(&body)

	c.OK(body)
}

func SendKafkaMessage(c *RequestContext) {
	key := gkafka.KafkaProducer.SendMessageAsync("first-topic", "运行消费者并消费消息")

	c.OK(key)
}

func Insert(c *RequestContext) {
	var list = make([]map[string]interface{}, 0)
	var _map1 = utils.OfMap("ID", 210, "FWaiterID", 11, "FWaiterName", fmt.Sprintf("春辉修改"), "FCreatedAt", utils.NowTimestamp(), "FPriceDeal", "0.3223")
	var _map2 = utils.OfMap("ID", 211, "FWaiterID", 12, "FWaiterName", fmt.Sprintf("春辉修改"), "FCreatedAt", utils.NowTimestamp(), "FPriceDeal", "0.3223")
	var _map3 = utils.OfMap("ID", 212, "FWaiterID", 12, "FWaiterName", fmt.Sprintf("春辉修改"), "FCreatedAt", utils.NowTimestamp(), "FPriceDeal", "0.3223")

	list = append(list, _map1, _map2, _map3)

	newId, err := gsql.Client.Update("order_update_buik.txt", utils.OfMap("orderList", list))

	if err != nil {
		c.Failc(400, "%v", err)

		return
	}

	c.OK(newId)
}
