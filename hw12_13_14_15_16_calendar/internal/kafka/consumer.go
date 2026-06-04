package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	segmentio "github.com/segmentio/kafka-go"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/storage/models"
)

type KafkaConsumer struct {
	reader *segmentio.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *KafkaConsumer {
	r := segmentio.NewReader(segmentio.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		MaxWait:        1 * time.Second,
		CommitInterval: time.Second,
	})
	return &KafkaConsumer{reader: r}
}

func (c *KafkaConsumer) ReadNotification(ctx context.Context) (models.Notification, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return models.Notification{}, fmt.Errorf("kafka consumer read: %w", err)
	}
	var n models.Notification
	if err := json.Unmarshal(msg.Value, &n); err != nil {
		return models.Notification{}, fmt.Errorf("kafka consumer unmarshal: %w", err)
	}
	return n, nil
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
