package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/logger"
	"github.com/TovStol/hw12_13_14_15_calendar/internal/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(logg *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writer := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(writer, r)

		duration := time.Since(start).Seconds()
		metrics.RecordHTTPRequest(r.Method, r.URL.Path, writer.statusCode, duration)

		logg.Info(fmt.Sprintf(
			"%s [%s] %s %s %s %d %s %q",
			r.RemoteAddr,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.RequestURI(),
			r.Proto,
			writer.statusCode,
			time.Since(start),
			r.UserAgent(),
		))
	})
}
