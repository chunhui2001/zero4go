package gkafka

import (
	"context"
	"log"
	"strings"

	"github.com/IBM/sarama"

	. "github.com/chunhui2001/zero4go/pkg/logs"
)

// 🎯 Bonus：消费者读取 Key 和 Value（封装函数）
func ReadMessage(msg *sarama.ConsumerMessage) (string, string) {
	if msg.Key == nil {
		return "", string(msg.Value)
	}

	return string(msg.Key), string(msg.Value)
}

type ConsumerHandler struct {
	Brokers string
	GroupId string
	Topic   string
	Handler func(string, string)
}

func (h ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	Log.Infof("Creating consumer group: brokers=%s, topic=%s, groupId=%s", h.Brokers, h.Topic, h.GroupId)

	return nil
}

func (h ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	Log.Infof("Cleanup consumer group: brokers=%s, topic=%s, groupId=%s", h.Brokers, h.Topic, h.GroupId)

	return nil
}

func (h ConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		key, val := ReadMessage(msg)

		h.Handler(key, val)

		Log.Infof("消费了一条消息[OK]: Topic=%s, Key=%s, groupId=%s, Value=%s", h.Topic, key, h.GroupId, val)

		// 手动标记 offset（非常重要）
		sess.MarkMessage(msg, "")
	}

	return nil
}

// 🔥 消费者组最佳实践（生产环境）
// ✔ 1. 自动提交 offset → 不推荐: 容易出现重复消息。
// ✔ 2. 使用 sess.MarkMessage() 手动提交: sess.MarkMessage(msg, "")
// ✔ 3. 消费失败建议写入 重试 topic（DLQ）: main-topic → retry-topic → dlq-topic
// 我能帮你生成完整的 Kafka 重试队列 示例（生产级）。
func CreateConsumer(brokers string, groupId string, topic string, handler func(string, string)) error {
	topics := []string{topic}

	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), groupId, config)

	if err != nil {
		Log.Errorf("Error creating group: brokers=%s, topic=%s, groupId=%s, Error=%v", brokers, topic, groupId, err)

		return err
	}

	_handler := ConsumerHandler{
		Brokers: brokers,
		GroupId: groupId,
		Topic:   topic,
		Handler: handler,
	}

	ctx := context.Background()

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, topics, _handler); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
		}
	}()

	// Log.Infof("Consumer started...")

	return nil
}
