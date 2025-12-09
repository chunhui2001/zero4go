package logs

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriter(brokers []string, topic string) *KafkaWriter {
	return &KafkaWriter{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      brokers,
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    100,                    // 批量条数
			BatchBytes:   10 * 1024,              // 批量最大字节
			BatchTimeout: 500 * time.Millisecond, // 最长等待时间
		}),
	}
}

// 实现 io.Writer 接口
func (kw *KafkaWriter) Write(p []byte) (n int, err error) {
	// 异步写入
	go func(data []byte) {
		err := kw.writer.WriteMessages(context.Background(), kafka.Message{
			Value: data,
		})

		if err != nil {
			// 可选：打印或收集发送失败日志
			log.Printf("Kafka write error: %v", err)
		}
	}(append([]byte(nil), p...)) // 拷贝 p 避免被覆盖

	return len(p), nil
}
