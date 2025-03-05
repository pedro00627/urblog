package infrastructure

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriter(broker string) *KafkaWriter {
	return &KafkaWriter{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{broker},
			Topic:   "tweets",
		}),
	}
}

func (kw *KafkaWriter) WriteMessage(message []byte) error {
	err := kw.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key"),
			Value: message,
		},
	)
	if err != nil {
		log.Println("Error writing message to Kafka:", err)
		return err
	}
	return nil
}
