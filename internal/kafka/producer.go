package kafka

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer is a Kafka producer implementation.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer.
func NewProducer() *Producer {
	kafkaURL := os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		kafkaURL = "localhost:9092"
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafkaURL),
		Topic:        "", // Will be set per message
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}

	return &Producer{
		writer: writer,
	}
}

// SendMessage sends a message to a Kafka topic.
func (p *Producer) SendMessage(topic string, message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p.writer.Topic = topic
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("msg_%d", time.Now().UnixNano())),
		Value: message,
	})

	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	return nil
}

// Close closes the Kafka producer.
func (p *Producer) Close() error {
	return p.writer.Close()
}
