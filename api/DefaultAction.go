package api

import (
	"github.com/chunhui2001/zero4go/pkg/gkafka"
	"github.com/chunhui2001/zero4go/pkg/logs"
	. "github.com/chunhui2001/zero4go/pkg/server"
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
