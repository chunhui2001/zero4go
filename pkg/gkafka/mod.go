package gkafka

import (
	"strings"

	"github.com/IBM/sarama"

	. "github.com/chunhui2001/zero4go/pkg/logs" //nolint:staticcheck
)

type KafkaConf struct {
	Enable           bool   `mapstructure:"KAFKA_ENABLE" json:"kafka_enable"`
	BootstrapServers string `mapstructure:"BOOTSTRAP_SERVERS" json:"bootstrap_servers"`
	MessageTimeoutMs uint32 `mapstructure:"MESSAGE_TIMEOUT_MS" json:"message_timeout_ms"`
	//Topic            string `mapstructure:"TOPIC" json:"topic"` // topic
}

var Settings = &KafkaConf{
	Enable:           true,
	BootstrapServers: "localhost:9092",
	MessageTimeoutMs: 5000,
}

func readKey(key sarama.Encoder) string {
	kb, _ := key.Encode()

	return string(kb)
}

var KafkaProducer *KafkaClient

func Init() {
	if !Settings.Enable {
		Log.Infof("Kafka-Disabled: val=%t", Settings.Enable)
		return
	}

	brokers := strings.Split(Settings.BootstrapServers, ",")

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

	if err != nil {
		Log.Errorf("Kafka-Failed: bootstrap_servers=%s, Error=%v", Settings.BootstrapServers, err.Error())

		return
	}

	producerAsync, err := sarama.NewAsyncProducer(brokers, config) // 异步生产者

	if err != nil {
		Log.Errorf("Kafka-Failed: bootstrap_servers=%s, Error=%v", Settings.BootstrapServers, err.Error())

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

	Log.Infof("Kafka-Succeed: bootstrap_servers=%s", Settings.BootstrapServers)
}
