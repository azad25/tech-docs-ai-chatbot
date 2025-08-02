package kafka

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer is a Kafka producer for sending messages
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer
func NewProducer() *Producer {
	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		kafkaURL = "localhost:9092"
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafkaURL),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
	}

	return &Producer{
		writer: writer,
	}
}

// SendMessage sends a message to the specified Kafka topic
func (p *Producer) SendMessage(topic string, message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: message,
	})

	if err != nil {
		return fmt.Errorf("failed to send message to topic %s: %w", topic, err)
	}

	return nil
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	return p.writer.Close()
}