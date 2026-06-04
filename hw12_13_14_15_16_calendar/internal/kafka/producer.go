package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	segmentio "github.com/segmentio/kafka-go"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/storage/models"
)

type KafkaProducer struct {
	writer *segmentio.Writer
}

func NewProducer(brokers []string, topic string) *KafkaProducer {
	w := &segmentio.Writer{
		Addr:         segmentio.TCP(brokers...),
		Topic:        topic,
		Balancer:     &segmentio.LeastBytes{},
		WriteTimeout: 10 * time.Second,
	}
	return &KafkaProducer{writer: w}
}

func (p *KafkaProducer) SendNotification(ctx context.Context, n models.Notification) error {
	data, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("kafka producer marshal: %w", err)
	}
	return p.writer.WriteMessages(ctx, segmentio.Message{
		Key:   []byte(fmt.Sprintf("%d", n.EventID)),
		Value: data,
	})
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
