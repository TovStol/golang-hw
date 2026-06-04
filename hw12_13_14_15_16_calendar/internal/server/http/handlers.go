package internalhttp

import (
	"context"
	"time"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/app"
	"github.com/TovStol/hw12_13_14_15_calendar/internal/storage/models"
)

// CalendarHandlers implements StrictServerInterface using the App layer.
type CalendarHandlers struct {
	app *app.App
}

func NewCalendarHandlers(a *app.App) *CalendarHandlers {
	return &CalendarHandlers{app: a}
}

func (h *CalendarHandlers) CreateEvent(ctx context.Context, input EventInput) (EventResponse, error) {
	event := inputToModel(0, input)
	id, err := h.app.Create(ctx, event)
	if err != nil {
		return EventResponse{}, err
	}
	event.ID = id
	return modelToResponse(event), nil
}

func (h *CalendarHandlers) UpdateEvent(ctx context.Context, id int64, input EventInput) (EventResponse, error) {
	event := inputToModel(id, input)
	h.app.Update(ctx, event)
	return modelToResponse(event), nil
}

func (h *CalendarHandlers) DeleteEvent(ctx context.Context, id int64) error {
	return h.app.DeleteByID(ctx, id)
}

func (h *CalendarHandlers) ListEventsByDay(ctx context.Context, date time.Time) ([]EventResponse, error) {
	events, err := h.app.FindEventsByDay(ctx, date)
	if err != nil {
		return nil, err
	}
	return modelsToResponses(events), nil
}

func (h *CalendarHandlers) ListEventsByWeek(ctx context.Context, date time.Time) ([]EventResponse, error) {
	events, err := h.app.FindEventsByWeek(ctx, date)
	if err != nil {
		return nil, err
	}
	return modelsToResponses(events), nil
}

func (h *CalendarHandlers) ListEventsByMonth(ctx context.Context, date time.Time) ([]EventResponse, error) {
	events, err := h.app.FindEventsByMonth(ctx, date)
	if err != nil {
		return nil, err
	}
	return modelsToResponses(events), nil
}

func inputToModel(id int64, input EventInput) models.Event {
	return models.Event{
		ID:           id,
		Title:        input.Title,
		DateTime:     input.DateTime,
		EndDateTime:  input.EndDateTime,
		Description:  input.Description,
		UserID:       input.UserID,
		NotifyBefore: time.Duration(input.NotifyBefore),
	}
}

func modelToResponse(e models.Event) EventResponse {
	return EventResponse{
		ID:           e.ID,
		Title:        e.Title,
		DateTime:     e.DateTime,
		EndDateTime:  e.EndDateTime,
		Description:  e.Description,
		UserID:       e.UserID,
		NotifyBefore: int64(e.NotifyBefore),
	}
}

func modelsToResponses(events []models.Event) []EventResponse {
	result := make([]EventResponse, 0, len(events))
	for _, e := range events {
		result = append(result, modelToResponse(e))
	}
	return result
}
