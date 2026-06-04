package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "calendar_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "calendar_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Event operation metrics
	eventsCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "calendar_events_created_total",
			Help: "Total number of events created",
		},
	)

	eventsUpdatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "calendar_events_updated_total",
			Help: "Total number of events updated",
		},
	)

	eventsDeletedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "calendar_events_deleted_total",
			Help: "Total number of events deleted",
		},
	)

	// Scheduler metrics
	schedulerRunsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "calendar_scheduler_runs_total",
			Help: "Total number of scheduler runs",
		},
	)

	schedulerNotificationsSent = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "calendar_scheduler_notifications_sent_total",
			Help: "Total number of notifications sent by scheduler",
		},
	)

	schedulerEventsProcessed = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "calendar_scheduler_events_processed_total",
			Help: "Total number of events processed by scheduler",
		},
	)

	schedulerErrorsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "calendar_scheduler_errors_total",
			Help: "Total number of scheduler errors",
		},
	)

	// Storer metrics
	storerNotificationsReceived = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "calendar_storer_notifications_received_total",
			Help: "Total number of notifications received by storer",
		},
	)

	storerNotificationsStored = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "calendar_storer_notifications_stored_total",
			Help: "Total number of notifications stored by storer",
		},
	)

	storerErrorsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "calendar_storer_errors_total",
			Help: "Total number of storer errors",
		},
	)
)

// HTTP metrics
func RecordHTTPRequest(method, path string, status int, duration float64) {
	httpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
	httpRequestDuration.WithLabelValues(method, path).Observe(duration)
}

// Event metrics
func RecordEventCreated() {
	eventsCreatedTotal.Inc()
}

func RecordEventUpdated() {
	eventsUpdatedTotal.Inc()
}

func RecordEventDeleted() {
	eventsDeletedTotal.Inc()
}

// Scheduler metrics
func RecordSchedulerRun() {
	schedulerRunsTotal.Inc()
}

func RecordNotificationSent() {
	schedulerNotificationsSent.Inc()
}

func RecordEventProcessed() {
	schedulerEventsProcessed.Inc()
}

func RecordSchedulerError() {
	schedulerErrorsTotal.Inc()
}

// Storer metrics
func RecordNotificationReceived() {
	storerNotificationsReceived.Inc()
}

func RecordNotificationStored() {
	storerNotificationsStored.Inc()
}

func RecordStorerError() {
	storerErrorsTotal.Inc()
}
