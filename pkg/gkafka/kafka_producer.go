package gkafka

import (
	"github.com/IBM/sarama"
	"github.com/google/uuid"

	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
)

type KafkaClient struct {
	ProducerSync  sarama.SyncProducer
	ProducerAsync sarama.AsyncProducer
}

func (k KafkaClient) SendMessageAsync(topic string, message string) string {
	key := uuid.New().String()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}

	k.ProducerAsync.Input() <- msg

	return key
}

func (k KafkaClient) SendMessage(topic string, message string) string {
	key := uuid.New().String()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := k.ProducerSync.SendMessage(msg)

	if err != nil {
		Log.Errorf("kafka SendMessage error: Error=%v", err.Error())

		return ""
	}

	Log.Infof("kafka SendMessage success: topic=%s, key=%s, partition=%d, offset=%d", topic, key, partition, offset)

	return key
}
