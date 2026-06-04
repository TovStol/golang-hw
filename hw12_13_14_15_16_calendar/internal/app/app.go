package app

import (
	"context"
	"time"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/logger"
	"github.com/TovStol/hw12_13_14_15_calendar/internal/storage/basic"
	"github.com/TovStol/hw12_13_14_15_calendar/internal/storage/models"
)

type App struct {
	Logger  *logger.Logger
	Storage basic.Storage
}

func New(logger *logger.Logger, storage basic.Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) Create(ctx context.Context, event models.Event) (int64, error) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.Storage.Create(event)
}

func (a *App) Update(ctx context.Context, event models.Event) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	a.Storage.Update(event)
}

func (a *App) DeleteByID(ctx context.Context, eventID int64) (err error) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.Storage.DeleteByID(eventID)
}

func (a *App) FindEventsByDay(ctx context.Context, date time.Time) ([]models.Event, error) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.Storage.FindEventsByDay(date)
}

func (a *App) FindEventsByWeek(ctx context.Context, date time.Time) ([]models.Event, error) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.Storage.FindEventsByWeek(date)
}

func (a *App) FindEventsByMonth(ctx context.Context, date time.Time) ([]models.Event, error) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.Storage.FindEventsByMonth(date)
}
