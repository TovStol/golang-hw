package basic

import (
	"time"

	"github.com/TovStol/hw12_13_14_15_16_calendar/internal/storage/models"
)

type Storage interface {
	Create(event models.Event) (ID int64, err error)
	Update(event models.Event)
	DeleteByID(eventID int64) (err error)
	FindEventsByDay(date time.Time) ([]models.Event, error)
	FindEventsByWeek(date time.Time) ([]models.Event, error)
	FindEventsByMonth(date time.Time) ([]models.Event, error)
}
