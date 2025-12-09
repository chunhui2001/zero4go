package gkafka

import (
	"strings"

	"github.com/IBM/sarama"
	"github.com/google/uuid"

	. "github.com/chunhui2001/zero4go/pkg/logs"
)

var KafkaProducer *KafkaClient

func InitKafka() {

	brokers := strings.Split(KafkaSetting.BootstrapServers, ",")

	config := sarama.NewConfig()

	// ---------- 最关键的生产环境配置 ----------
	config.Net.MaxOpenRequests = 1                          // 幂等要求
	config.Producer.RequiredAcks = sarama.WaitForAll        // 等待所有副本确认
	config.Producer.Idempotent = true                       // 幂等生产者（生产环境必须）
	config.Producer.Retry.Max = 5                           // 重试
	config.Producer.Return.Successes = true                 // 同步生产必须开启
	config.Producer.Compression = sarama.CompressionSnappy  // 压缩提高吞吐
	config.Producer.Partitioner = sarama.NewHashPartitioner // 默认分区策略

	producerSync, err := sarama.NewSyncProducer(brokers, config)
	producerAsync, err := sarama.NewAsyncProducer(brokers, config) // 异步生产者

	if err != nil {
		Log.Errorf("kafka init failed: bootstrap_servers=%s, Error=%v", KafkaSetting.BootstrapServers, err.Error())

		return
	}

	// 读取结果的 goroutine
	go func() {
		for {
			select {
			case suc := <-producerAsync.Successes():
				Log.Infof("发送了一条消息[OK]: Topic=%s, Key=%s, offset=%d, partition=%d", suc.Topic, readKey(suc.Key), suc.Offset, suc.Partition)
			case err := <-producerAsync.Errors():
				Log.Errorf("发送了一条消息[ERR]: Error=%v", err)
			}
		}
	}()

	KafkaProducer = &KafkaClient{
		ProducerSync:  producerSync,
		ProducerAsync: producerAsync,
	}

	Log.Infof("kafka init success: bootstrap_servers=%s", KafkaSetting.BootstrapServers)
}

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
