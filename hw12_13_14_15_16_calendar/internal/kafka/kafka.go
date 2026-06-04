package kafka

import (
	"context"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/storage/models"
)

type Producer interface {
	SendNotification(ctx context.Context, n models.Notification) error
	Close() error
}

type Consumer interface {
	ReadNotification(ctx context.Context) (models.Notification, error)
	Close() error
}
