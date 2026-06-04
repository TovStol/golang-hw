package app

import (
	"context"
	"testing"
	"time"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/TovStol/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/TovStol/hw12_13_14_15_calendar/internal/storage/models"
)

func newTestApp() *App {
	logg := logger.New("error", "")
	storage := memorystorage.New()
	return New(logg, storage)
}

func newEvent(title string, dt time.Time) models.Event {
	return models.Event{
		Title:       title,
		DateTime:    dt,
		EndDateTime: dt.Add(time.Hour),
		UserID:      1,
	}
}

func TestApp_Create(t *testing.T) {
	a := newTestApp()
	ctx := context.Background()

	id, err := a.Create(ctx, newEvent("Meeting", time.Now()))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if id != 0 {
		t.Errorf("Create() id = %d, want 0", id)
	}

	id2, err := a.Create(ctx, newEvent("Standup", time.Now()))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if id2 != 1 {
		t.Errorf("Create() id = %d, want 1", id2)
	}
}

func TestApp_Update(t *testing.T) {
	a := newTestApp()
	ctx := context.Background()

	dt := time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC)
	id, _ := a.Create(ctx, newEvent("Original", dt))

	updated := models.Event{
		ID:          id,
		Title:       "Updated",
		DateTime:    dt,
		EndDateTime: dt.Add(time.Hour),
		UserID:      1,
	}
	a.Update(ctx, updated)

	events, err := a.FindEventsByDay(ctx, dt)
	if err != nil {
		t.Fatalf("FindEventsByDay() error = %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("FindEventsByDay() len = %d, want 1", len(events))
	}
	if events[0].Title != "Updated" {
		t.Errorf("Title = %q, want 'Updated'", events[0].Title)
	}
}

func TestApp_DeleteByID(t *testing.T) {
	a := newTestApp()
	ctx := context.Background()

	dt := time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC)
	id, _ := a.Create(ctx, newEvent("To delete", dt))

	if err := a.DeleteByID(ctx, id); err != nil {
		t.Fatalf("DeleteByID() error = %v", err)
	}

	events, _ := a.FindEventsByDay(ctx, dt)
	if len(events) != 0 {
		t.Errorf("after DeleteByID, FindEventsByDay len = %d, want 0", len(events))
	}
}

func TestApp_FindEventsByDay(t *testing.T) {
	a := newTestApp()
	ctx := context.Background()

	day := time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC)
	a.Create(ctx, newEvent("E1", day.Add(9*time.Hour)))
	a.Create(ctx, newEvent("E2", day.Add(14*time.Hour)))
	a.Create(ctx, newEvent("E3", day.Add(25*time.Hour))) // next day

	events, err := a.FindEventsByDay(ctx, day)
	if err != nil {
		t.Fatalf("FindEventsByDay() error = %v", err)
	}
	if len(events) != 2 {
		t.Errorf("FindEventsByDay() len = %d, want 2", len(events))
	}
}

func TestApp_FindEventsByWeek(t *testing.T) {
	a := newTestApp()
	ctx := context.Background()

	mon := time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)  // week 24
	fri := time.Date(2026, 6, 12, 10, 0, 0, 0, time.UTC) // same week
	nxt := time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC) // next week

	a.Create(ctx, newEvent("Mon", mon))
	a.Create(ctx, newEvent("Fri", fri))
	a.Create(ctx, newEvent("NextWeek", nxt))

	events, err := a.FindEventsByWeek(ctx, mon)
	if err != nil {
		t.Fatalf("FindEventsByWeek() error = %v", err)
	}
	if len(events) != 2 {
		t.Errorf("FindEventsByWeek() len = %d, want 2", len(events))
	}
}

func TestApp_FindEventsByMonth(t *testing.T) {
	a := newTestApp()
	ctx := context.Background()

	jun1 := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	jun30 := time.Date(2026, 6, 30, 10, 0, 0, 0, time.UTC)
	jul1 := time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC)

	a.Create(ctx, newEvent("Jun 1", jun1))
	a.Create(ctx, newEvent("Jun 30", jun30))
	a.Create(ctx, newEvent("Jul 1", jul1))

	events, err := a.FindEventsByMonth(ctx, jun1)
	if err != nil {
		t.Fatalf("FindEventsByMonth() error = %v", err)
	}
	if len(events) != 2 {
		t.Errorf("FindEventsByMonth() len = %d, want 2", len(events))
	}
	for _, e := range events {
		if e.DateTime.Month() != time.June {
			t.Errorf("event %q not in June", e.Title)
		}
	}
}

func TestApp_FindEventsByDay_Empty(t *testing.T) {
	a := newTestApp()
	ctx := context.Background()

	events, err := a.FindEventsByDay(ctx, time.Now())
	if err != nil {
		t.Fatalf("FindEventsByDay() error = %v", err)
	}
	if len(events) != 0 {
		t.Errorf("FindEventsByDay() on empty app len = %d, want 0", len(events))
	}
}
