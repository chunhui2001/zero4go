package gkafka

import (
	"github.com/IBM/sarama"
)

type KafkaConf struct {
	Enable           bool   `mapstructure:"KAFKA_ENABLE" json:"kafka_enable"`
	BootstrapServers string `mapstructure:"BOOTSTRAP_SERVERS" json:"bootstrap_servers"`
	MessageTimeoutMs uint32 `mapstructure:"MESSAGE_TIMEOUT_MS" json:"message_timeout_ms"`
	//Topic            string `mapstructure:"TOPIC" json:"topic"` // topic
}

var KafkaSetting = &KafkaConf{
	BootstrapServers: "localhost:9092",
	MessageTimeoutMs: 5000,
}

func readKey(key sarama.Encoder) string {
	kb, _ := key.Encode()

	return string(kb)
}
