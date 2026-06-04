//go:build integration

package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/server/http"
)

const (
	calendarURL = "http://127.0.0.1:8081"
	pgDSN       = "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=idm_tests sslmode=disable"
)

// helpers

func mustJSON(t *testing.T, v any) *bytes.Reader {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	return bytes.NewReader(b)
}

func createEvent(t *testing.T, input internalhttp.EventInput) internalhttp.EventResponse {
	t.Helper()
	resp, err := http.Post(calendarURL+"/events", "application/json", mustJSON(t, input))
	if err != nil {
		t.Fatalf("POST /events: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST /events status = %d, want 201", resp.StatusCode)
	}
	var result internalhttp.EventResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode EventResponse: %v", err)
	}
	return result
}

func deleteEvent(t *testing.T, id int64) {
	t.Helper()
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete,
		fmt.Sprintf("%s/events/%d", calendarURL, id), nil)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatalf("DELETE /events/%d: %v", id, err)
	}
	resp.Body.Close()
}

func pgDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", pgDSN)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("db.Ping: %v", err)
	}
	return db
}

// ── Calendar HTTP API ──────────────────────────────────────────────────────────

func TestIntegration_CalendarHealthCheck(t *testing.T) {
	resp, err := http.Get(calendarURL + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestIntegration_CreateEvent(t *testing.T) {
	input := internalhttp.EventInput{
		Title:       "Integration Create",
		DateTime:    time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second),
		EndDateTime: time.Now().UTC().Add(25 * time.Hour).Truncate(time.Second),
		UserID:      99,
	}
	ev := createEvent(t, input)
	defer deleteEvent(t, ev.ID)

	if ev.Title != input.Title {
		t.Errorf("title = %q, want %q", ev.Title, input.Title)
	}
	if ev.UserID != input.UserID {
		t.Errorf("user_id = %d, want %d", ev.UserID, input.UserID)
	}
}

func TestIntegration_UpdateEvent(t *testing.T) {
	dt := time.Now().UTC().Add(48 * time.Hour).Truncate(time.Second)
	ev := createEvent(t, internalhttp.EventInput{
		Title:       "Before Update",
		DateTime:    dt,
		EndDateTime: dt.Add(time.Hour),
		UserID:      1,
	})
	defer deleteEvent(t, ev.ID)

	updated := internalhttp.EventInput{
		Title:       "After Update",
		DateTime:    dt,
		EndDateTime: dt.Add(time.Hour),
		UserID:      1,
	}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut,
		fmt.Sprintf("%s/events/%d", calendarURL, ev.ID),
		mustJSON(t, updated))
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatalf("PUT /events/%d: %v", ev.ID, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	var result internalhttp.EventResponse
	json.NewDecoder(resp.Body).Decode(&result) //nolint:errcheck
	if result.Title != "After Update" {
		t.Errorf("title = %q, want 'After Update'", result.Title)
	}
}

func TestIntegration_DeleteEvent(t *testing.T) {
	dt := time.Now().UTC().Add(72 * time.Hour).Truncate(time.Second)
	ev := createEvent(t, internalhttp.EventInput{
		Title:       "To Delete",
		DateTime:    dt,
		EndDateTime: dt.Add(time.Hour),
		UserID:      1,
	})

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete,
		fmt.Sprintf("%s/events/%d", calendarURL, ev.ID), nil)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatalf("DELETE /events/%d: %v", ev.ID, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d, want 204", resp.StatusCode)
	}
}

func TestIntegration_ListEventsByDay(t *testing.T) {
	day := time.Date(2099, 1, 15, 10, 0, 0, 0, time.UTC)
	ev1 := createEvent(t, internalhttp.EventInput{
		Title: "Day1", DateTime: day, EndDateTime: day.Add(time.Hour), UserID: 1,
	})
	ev2 := createEvent(t, internalhttp.EventInput{
		Title: "Day2", DateTime: day.Add(2 * time.Hour), EndDateTime: day.Add(3 * time.Hour), UserID: 1,
	})
	evNext := createEvent(t, internalhttp.EventInput{
		Title: "NextDay", DateTime: day.Add(25 * time.Hour), EndDateTime: day.Add(26 * time.Hour), UserID: 1,
	})
	defer deleteEvent(t, ev1.ID)
	defer deleteEvent(t, ev2.ID)
	defer deleteEvent(t, evNext.ID)

	resp, err := http.Get(calendarURL + "/events/day?date=2099-01-15")
	if err != nil {
		t.Fatalf("GET /events/day: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var events []internalhttp.EventResponse
	json.NewDecoder(resp.Body).Decode(&events) //nolint:errcheck
	if len(events) != 2 {
		t.Errorf("len = %d, want 2", len(events))
	}
}

func TestIntegration_ListEventsByWeek(t *testing.T) {
	mon := time.Date(2099, 2, 3, 10, 0, 0, 0, time.UTC)
	fri := time.Date(2099, 2, 7, 10, 0, 0, 0, time.UTC)
	nxt := time.Date(2099, 2, 10, 10, 0, 0, 0, time.UTC)
	e1 := createEvent(t, internalhttp.EventInput{Title: "Mon", DateTime: mon, EndDateTime: mon.Add(time.Hour), UserID: 1})
	e2 := createEvent(t, internalhttp.EventInput{Title: "Fri", DateTime: fri, EndDateTime: fri.Add(time.Hour), UserID: 1})
	e3 := createEvent(t, internalhttp.EventInput{Title: "Nxt", DateTime: nxt, EndDateTime: nxt.Add(time.Hour), UserID: 1})
	defer deleteEvent(t, e1.ID)
	defer deleteEvent(t, e2.ID)
	defer deleteEvent(t, e3.ID)

	resp, err := http.Get(calendarURL + "/events/week?date=2099-02-03")
	if err != nil {
		t.Fatalf("GET /events/week: %v", err)
	}
	defer resp.Body.Close()
	var events []internalhttp.EventResponse
	json.NewDecoder(resp.Body).Decode(&events) //nolint:errcheck
	if len(events) != 2 {
		t.Errorf("len = %d, want 2", len(events))
	}
}

func TestIntegration_ListEventsByMonth(t *testing.T) {
	jun1 := time.Date(2099, 6, 1, 10, 0, 0, 0, time.UTC)
	jun30 := time.Date(2099, 6, 30, 10, 0, 0, 0, time.UTC)
	jul1 := time.Date(2099, 7, 1, 10, 0, 0, 0, time.UTC)
	e1 := createEvent(t, internalhttp.EventInput{Title: "Jun1", DateTime: jun1, EndDateTime: jun1.Add(time.Hour), UserID: 1})
	e2 := createEvent(t, internalhttp.EventInput{Title: "Jun30", DateTime: jun30, EndDateTime: jun30.Add(time.Hour), UserID: 1})
	e3 := createEvent(t, internalhttp.EventInput{Title: "Jul1", DateTime: jul1, EndDateTime: jul1.Add(time.Hour), UserID: 1})
	defer deleteEvent(t, e1.ID)
	defer deleteEvent(t, e2.ID)
	defer deleteEvent(t, e3.ID)

	resp, err := http.Get(calendarURL + "/events/month?date=2099-06-01")
	if err != nil {
		t.Fatalf("GET /events/month: %v", err)
	}
	defer resp.Body.Close()
	var events []internalhttp.EventResponse
	json.NewDecoder(resp.Body).Decode(&events) //nolint:errcheck
	if len(events) != 2 {
		t.Errorf("len = %d, want 2", len(events))
	}
}

func TestIntegration_InvalidDate(t *testing.T) {
	for _, path := range []string{"/events/day", "/events/week", "/events/month"} {
		resp, err := http.Get(calendarURL + path + "?date=notadate")
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("%s status = %d, want 400", path, resp.StatusCode)
		}
	}
}

// ── Scheduler → Kafka → Storer pipeline ────────────────────────────────────────

// TestIntegration_SchedulerStorerPipeline creates an event with notify_before
// set to fire immediately (event is in the future but within the notify window),
// then waits for the scheduler to pick it up, send it to Kafka, and the storer
// to persist it in the notification table.
func TestIntegration_SchedulerStorerPipeline(t *testing.T) {
	db := pgDB(t)
	defer db.Close()

	// Clean up any stale notifications from previous runs.
	db.Exec("DELETE FROM notification WHERE user_id = 777") //nolint:errcheck

	// Event fires in 30 seconds; notify_before = 60s → should be picked up immediately.
	notifyBefore := int64(60 * time.Second)
	eventTime := time.Now().UTC().Add(30 * time.Second).Truncate(time.Second)
	ev := createEvent(t, internalhttp.EventInput{
		Title:        "Pipeline Test Event",
		DateTime:     eventTime,
		EndDateTime:  eventTime.Add(time.Hour),
		UserID:       777,
		NotifyBefore: notifyBefore,
	})
	defer deleteEvent(t, ev.ID)
	defer db.Exec("DELETE FROM notification WHERE event_id = $1", ev.ID) //nolint:errcheck

	// Wait up to 30s for the notification to land in the DB
	// (scheduler scans every 10s, storer processes immediately).
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		var count int
		err := db.QueryRow(
			"SELECT COUNT(*) FROM notification WHERE event_id = $1", ev.ID,
		).Scan(&count)
		if err != nil {
			t.Fatalf("querying notification: %v", err)
		}
		if count > 0 {
			t.Logf("notification for event_id=%d found in DB after %.1fs",
				ev.ID, time.Until(deadline).Abs().Seconds())
			return
		}
		time.Sleep(2 * time.Second)
	}

	t.Errorf("notification for event_id=%d was NOT stored within 30s", ev.ID)
}

// TestIntegration_StoredNotificationFields verifies the notification fields
// written by the storer match the original event.
func TestIntegration_StoredNotificationFields(t *testing.T) {
	db := pgDB(t)
	defer db.Close()

	notifyBefore := int64(120 * time.Second)
	eventTime := time.Now().UTC().Add(60 * time.Second).Truncate(time.Second)
	ev := createEvent(t, internalhttp.EventInput{
		Title:        "Field Check Event",
		DateTime:     eventTime,
		EndDateTime:  eventTime.Add(time.Hour),
		UserID:       888,
		NotifyBefore: notifyBefore,
	})
	defer deleteEvent(t, ev.ID)
	defer db.Exec("DELETE FROM notification WHERE event_id = $1", ev.ID) //nolint:errcheck

	deadline := time.Now().Add(30 * time.Second)
	var (
		storedTitle  string
		storedUserID int64
		storedDate   time.Time
	)
	for time.Now().Before(deadline) {
		err := db.QueryRow(
			"SELECT title, user_id, event_date FROM notification WHERE event_id = $1", ev.ID,
		).Scan(&storedTitle, &storedUserID, &storedDate)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if storedTitle == "" {
		t.Fatalf("notification for event_id=%d not found within 30s", ev.ID)
	}
	if storedTitle != ev.Title {
		t.Errorf("title = %q, want %q", storedTitle, ev.Title)
	}
	if storedUserID != ev.UserID {
		t.Errorf("user_id = %d, want %d", storedUserID, ev.UserID)
	}
	if !storedDate.Equal(eventTime) {
		t.Errorf("event_date = %v, want %v", storedDate, eventTime)
	}
}
