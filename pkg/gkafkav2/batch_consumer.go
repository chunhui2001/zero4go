package gkafkav2

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	. "github.com/chunhui2001/zero4go/pkg/logs"
)

type BatchKafkaConsumer struct {
	*kafka.Consumer
	Topic          string
	msgCh          chan *Msg
	BatchSize      int
	CommitFlushDur time.Duration
	GroupID        string
}

func NewConsumer(broker string, groupId string, topic string, offset string) *BatchKafkaConsumer {
	if _c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          groupId,
		//"auto.offset.reset": "latest",
		//"auto.offset.reset": "earliest",
		"auto.offset.reset": offset,

		// 推荐显式配置
		"enable.auto.commit":    false,
		"session.timeout.ms":    10000,
		"heartbeat.interval.ms": 3000,

		// 批量拉取
		"fetch.min.bytes":           1 * 1024 * 1024, // 1MB
		"fetch.wait.max.ms":         100,
		"max.partition.fetch.bytes": 4 * 1024 * 1024,

		// 本地队列
		"queued.min.messages":        100,
		"queued.max.messages.kbytes": 64 * 1024,
	}); err != nil {
		Log.Errorf("Error creating new consumer: Error=%s", err.Error())

		return nil
	} else {
		if err := _c.SubscribeTopics([]string{topic}, func(c *kafka.Consumer, e kafka.Event) error {
			switch ev := e.(type) {
			case kafka.AssignedPartitions:
				Log.Infof("Partition: Topic=%s, GroupId=%s, Assigned: %+v", topic, groupId, ev.Partitions)

				return c.Assign(ev.Partitions)
			case kafka.RevokedPartitions:
				Log.Infof("Partition: Topic=%s, GroupId=%s, Revoked: %+v", topic, groupId, ev.Partitions)

				return c.Unassign()
			}

			return nil
		}); err != nil {
			Log.Errorf("Error creating new consumer: Error=%s", err.Error())

			return nil
		}

		Log.Infof("创建了一个 kafka 消费者: Topic=%s, GroupId=%s, Broker=%s", topic, groupId, broker)

		return &BatchKafkaConsumer{
			Consumer:       _c,
			Topic:          topic,
			GroupID:        groupId,
			BatchSize:      500,
			CommitFlushDur: 50 * time.Millisecond,
			msgCh:          make(chan *Msg, 5000),
		}
	}
}

func (c *BatchKafkaConsumer) Start(cb func(topic string, groupId string, msgs []*Msg) error) {
	go func() {
		for {
			ev := c.Poll(100)

			if ev == nil {
				time.Sleep(1000 * time.Millisecond)

				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				// c.CommitMessage(e) // 异步 FlushCommit

				c.msgCh <- &Msg{
					Key:       string(e.Key),
					Value:     string(e.Value),
					Partition: e.TopicPartition.Partition,
					Offset:    int64(e.TopicPartition.Offset),
				}
			case kafka.Error:
				Log.Errorf("BatchKafkaConsumer Error: Error%s", e.Error())
			}
		}
	}()

	go func() {
		c.FlushCommit(cb)
	}()
}

func (c *BatchKafkaConsumer) FlushCommit(cb func(topic string, groupId string, msgs []*Msg) error) {
	ticker := time.NewTicker(c.CommitFlushDur)

	defer ticker.Stop()

	_msgs := make([]*Msg, 0, c.BatchSize)

	for {
		select {
		case msg := <-c.msgCh:
			_msgs = append(_msgs, msg)

			if len(_msgs) >= c.BatchSize {
				if c.commit(_msgs, cb) {
					_msgs = _msgs[:0]
				}
			}
		case <-ticker.C:
			if len(_msgs) > 0 {
				if c.commit(_msgs, cb) {
					_msgs = _msgs[:0]
				}
			}
		}
	}
}

func (c *BatchKafkaConsumer) commit(msgs []*Msg, cb func(topic string, groupId string, msgs []*Msg) error) bool {
	if len(msgs) == 0 {
		return true
	}

	type tpKey struct {
		topic     string
		partition int32
	}

	offsets := make(map[tpKey]kafka.Offset)

	for _, m := range msgs {
		_tpKey := tpKey{
			topic:     c.Topic,
			partition: m.Partition,
		}

		// 记录每个分区最后处理的消息偏移量。
		// 每个 partition 只提交一次 offset
		offsets[_tpKey] = kafka.Offset(m.Offset + 1)
	}

	if err := cb(c.Topic, c.GroupID, msgs); err != nil {
		Log.Errorf("Callback failed: Topic=%s, GroupID=%s, Error=%v", c.Topic, c.GroupID, err)

		return false
	}

	var commits []kafka.TopicPartition

	for k, off := range offsets {
		topic := k.topic

		commits = append(commits, kafka.TopicPartition{Topic: &topic, Partition: k.partition, Offset: off})
	}

	if _, err := c.CommitOffsets(commits); err != nil {
		Log.Errorf("CommitOffsets failed: Topic=%s, Error=%v", c.Topic, err)
	}

	return true
}
