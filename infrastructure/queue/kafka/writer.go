package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Writer struct {
	writer *kafka.Writer
}

func NewWriter(broker string) *Writer {
	return &Writer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    "tweets",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (kw *Writer) WriteMessage(message []byte) error {
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
