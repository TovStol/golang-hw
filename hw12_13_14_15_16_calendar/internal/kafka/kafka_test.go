package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/storage/models"
)

// mockProducer implements Producer for testing.
type mockProducer struct {
	sent []models.Notification
	err  error
}

func (m *mockProducer) SendNotification(_ context.Context, n models.Notification) error {
	if m.err != nil {
		return m.err
	}
	m.sent = append(m.sent, n)
	return nil
}

func (m *mockProducer) Close() error { return m.err }

// mockConsumer implements Consumer for testing.
type mockConsumer struct {
	queue []models.Notification
	err   error
}

func (m *mockConsumer) ReadNotification(_ context.Context) (models.Notification, error) {
	if m.err != nil {
		return models.Notification{}, m.err
	}
	if len(m.queue) == 0 {
		return models.Notification{}, errors.New("empty queue")
	}
	n := m.queue[0]
	m.queue = m.queue[1:]
	return n, nil
}

func (m *mockConsumer) Close() error { return m.err }

// --- Interface compliance ---

func TestProducerInterface(t *testing.T) {
	var _ Producer = &mockProducer{}
}

func TestConsumerInterface(t *testing.T) {
	var _ Consumer = &mockConsumer{}
}

// --- mockProducer tests ---

func TestMockProducer_SendNotification(t *testing.T) {
	ts := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	n := models.Notification{EventID: 1, Title: "Stand-up", EventDate: ts, UserID: 42}

	p := &mockProducer{}
	if err := p.SendNotification(context.Background(), n); err != nil {
		t.Fatalf("SendNotification() unexpected error: %v", err)
	}
	if len(p.sent) != 1 {
		t.Fatalf("sent len = %d, want 1", len(p.sent))
	}
	if p.sent[0] != n {
		t.Errorf("sent[0] = %+v, want %+v", p.sent[0], n)
	}
}

func TestMockProducer_SendNotification_Error(t *testing.T) {
	want := errors.New("broker unavailable")
	p := &mockProducer{err: want}
	err := p.SendNotification(context.Background(), models.Notification{})
	if !errors.Is(err, want) {
		t.Errorf("error = %v, want %v", err, want)
	}
}

func TestMockProducer_Close(t *testing.T) {
	p := &mockProducer{}
	if err := p.Close(); err != nil {
		t.Errorf("Close() unexpected error: %v", err)
	}
}

// --- mockConsumer tests ---

func TestMockConsumer_ReadNotification(t *testing.T) {
	ts := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	n := models.Notification{EventID: 7, Title: "Retro", EventDate: ts, UserID: 3}

	c := &mockConsumer{queue: []models.Notification{n}}
	got, err := c.ReadNotification(context.Background())
	if err != nil {
		t.Fatalf("ReadNotification() unexpected error: %v", err)
	}
	if got != n {
		t.Errorf("got = %+v, want %+v", got, n)
	}
}

func TestMockConsumer_ReadNotification_EmptyQueue(t *testing.T) {
	c := &mockConsumer{}
	_, err := c.ReadNotification(context.Background())
	if err == nil {
		t.Error("ReadNotification() expected error on empty queue, got nil")
	}
}

func TestMockConsumer_ReadNotification_Error(t *testing.T) {
	want := errors.New("read error")
	c := &mockConsumer{err: want}
	_, err := c.ReadNotification(context.Background())
	if !errors.Is(err, want) {
		t.Errorf("error = %v, want %v", err, want)
	}
}

func TestMockConsumer_Close(t *testing.T) {
	c := &mockConsumer{}
	if err := c.Close(); err != nil {
		t.Errorf("Close() unexpected error: %v", err)
	}
}

// --- JSON round-trip (mirrors SendNotification / ReadNotification serialization) ---

func TestNotification_JSONRoundTrip(t *testing.T) {
	ts := time.Date(2026, 6, 10, 15, 30, 0, 0, time.UTC)
	cases := []models.Notification{
		{EventID: 1, Title: "Meeting", EventDate: ts, UserID: 10},
		{EventID: 0, Title: "", EventDate: time.Time{}, UserID: 0},
		{EventID: 999, Title: "Long title with spaces", EventDate: ts, UserID: 1},
	}

	for _, orig := range cases {
		data, err := json.Marshal(orig)
		if err != nil {
			t.Fatalf("Marshal(%+v) error: %v", orig, err)
		}
		var got models.Notification
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if got.EventID != orig.EventID || got.Title != orig.Title || got.UserID != orig.UserID {
			t.Errorf("round-trip mismatch: got %+v, want %+v", got, orig)
		}
		if !got.EventDate.Equal(orig.EventDate) {
			t.Errorf("EventDate mismatch: got %v, want %v", got.EventDate, orig.EventDate)
		}
	}
}

func TestNotification_JSONUnmarshal_InvalidJSON(t *testing.T) {
	var n models.Notification
	if err := json.Unmarshal([]byte("not json"), &n); err == nil {
		t.Error("Unmarshal(invalid) expected error, got nil")
	}
}

// --- Producer/Consumer pipeline via mocks ---

func TestProducerConsumer_Pipeline(t *testing.T) {
	ts := time.Date(2026, 6, 10, 9, 0, 0, 0, time.UTC)
	notifications := []models.Notification{
		{EventID: 1, Title: "A", EventDate: ts, UserID: 1},
		{EventID: 2, Title: "B", EventDate: ts.Add(time.Hour), UserID: 2},
	}

	p := &mockProducer{}
	ctx := context.Background()

	for _, n := range notifications {
		if err := p.SendNotification(ctx, n); err != nil {
			t.Fatalf("SendNotification(%+v) error: %v", n, err)
		}
	}

	c := &mockConsumer{queue: append([]models.Notification{}, p.sent...)}

	for i, want := range notifications {
		got, err := c.ReadNotification(ctx)
		if err != nil {
			t.Fatalf("ReadNotification()[%d] error: %v", i, err)
		}
		if got != want {
			t.Errorf("[%d] got %+v, want %+v", i, got, want)
		}
	}
}
