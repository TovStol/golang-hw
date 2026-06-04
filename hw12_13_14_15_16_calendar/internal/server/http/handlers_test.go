package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/app"
	"github.com/TovStol/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/TovStol/hw12_13_14_15_calendar/internal/storage/memory"
)

func newTestServer() *httptest.Server {
	logg := logger.New("error", "stdout")
	storage := memorystorage.New()
	application := app.New(logg, storage)

	mux := http.NewServeMux()
	RegisterHandlers(mux, NewCalendarHandlers(application))

	return httptest.NewServer(mux)
}

func TestCreateEvent(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	input := EventInput{
		Title:       "Test event",
		DateTime:    time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC),
		EndDateTime: time.Date(2026, 6, 10, 11, 0, 0, 0, time.UTC),
		UserID:      1,
	}
	body, _ := json.Marshal(input)

	resp, err := http.Post(srv.URL+"/events", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST /events: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status = %d, want 201", resp.StatusCode)
	}

	var result EventResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.Title != input.Title {
		t.Errorf("title = %q, want %q", result.Title, input.Title)
	}
	if result.UserID != input.UserID {
		t.Errorf("user_id = %d, want %d", result.UserID, input.UserID)
	}
}

func TestCreateEvent_InvalidBody(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/events", "application/json", bytes.NewReader([]byte("not-json")))
	if err != nil {
		t.Fatalf("POST /events: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestUpdateEvent(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	input := EventInput{
		Title:       "Original",
		DateTime:    time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC),
		EndDateTime: time.Date(2026, 6, 10, 11, 0, 0, 0, time.UTC),
		UserID:      1,
	}
	body, _ := json.Marshal(input)
	resp, _ := http.Post(srv.URL+"/events", "application/json", bytes.NewReader(body))
	var created EventResponse
	json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()

	updated := EventInput{
		Title:       "Updated",
		DateTime:    time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC),
		EndDateTime: time.Date(2026, 6, 10, 11, 0, 0, 0, time.UTC),
		UserID:      1,
	}
	updBody, _ := json.Marshal(updated)

	req, _ := http.NewRequestWithContext(
		context.Background(),
		http.MethodPut,
		srv.URL+"/events/0",
		bytes.NewReader(updBody),
	)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		t.Fatalf("PUT /events/0: %v", err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", r.StatusCode)
	}

	var result EventResponse
	json.NewDecoder(r.Body).Decode(&result)
	if result.Title != "Updated" {
		t.Errorf("title = %q, want 'Updated'", result.Title)
	}
}

func TestDeleteEvent(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	input := EventInput{
		Title:       "To delete",
		DateTime:    time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC),
		EndDateTime: time.Date(2026, 6, 10, 11, 0, 0, 0, time.UTC),
		UserID:      1,
	}
	body, _ := json.Marshal(input)
	http.Post(srv.URL+"/events", "application/json", bytes.NewReader(body)) //nolint:errcheck

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, srv.URL+"/events/0", nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("DELETE /events/0: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d, want 204", resp.StatusCode)
	}
}

func TestListEventsByDay(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	dt := time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC)
	for _, title := range []string{"A", "B"} {
		input := EventInput{
			Title:       title,
			DateTime:    dt,
			EndDateTime: dt.Add(time.Hour),
			UserID:      1,
		}
		body, _ := json.Marshal(input)
		http.Post(srv.URL+"/events", "application/json", bytes.NewReader(body)) //nolint:errcheck
	}

	resp, err := http.Get(srv.URL + "/events/day?date=2026-06-10")
	if err != nil {
		t.Fatalf("GET /events/day: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}

	var events []EventResponse
	json.NewDecoder(resp.Body).Decode(&events)
	if len(events) != 2 {
		t.Errorf("len(events) = %d, want 2", len(events))
	}
}

func TestListEventsByDay_InvalidDate(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/events/day?date=bad-date")
	if err != nil {
		t.Fatalf("GET /events/day: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestListEventsByWeek(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	dt := time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC)
	input := EventInput{
		Title:       "Week event",
		DateTime:    dt,
		EndDateTime: dt.Add(time.Hour),
		UserID:      1,
	}
	body, _ := json.Marshal(input)
	http.Post(srv.URL+"/events", "application/json", bytes.NewReader(body)) //nolint:errcheck

	resp, err := http.Get(srv.URL + "/events/week?date=2026-06-08")
	if err != nil {
		t.Fatalf("GET /events/week: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}

	var events []EventResponse
	json.NewDecoder(resp.Body).Decode(&events)
	if len(events) != 1 {
		t.Errorf("len(events) = %d, want 1", len(events))
	}
}

func TestListEventsByMonth(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	dt := time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC)
	input := EventInput{
		Title:       "Month event",
		DateTime:    dt,
		EndDateTime: dt.Add(time.Hour),
		UserID:      1,
	}
	body, _ := json.Marshal(input)
	http.Post(srv.URL+"/events", "application/json", bytes.NewReader(body)) //nolint:errcheck

	resp, err := http.Get(srv.URL + "/events/month?date=2026-06-01")
	if err != nil {
		t.Fatalf("GET /events/month: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}

	var events []EventResponse
	json.NewDecoder(resp.Body).Decode(&events)
	if len(events) != 1 {
		t.Errorf("len(events) = %d, want 1", len(events))
	}
}
