package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/kafka"
	"github.com/TovStol/hw12_13_14_15_calendar/internal/logger"
	"github.com/TovStol/hw12_13_14_15_calendar/internal/metrics"
	"github.com/TovStol/hw12_13_14_15_calendar/internal/storage/models"
	sqlstorage "github.com/TovStol/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/scheduler_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := NewConfig(configFile)
	logg := logger.New(config.Level, config.LogLocation)
	defer logg.Close() //nolint:errcheck

	storage := sqlstorage.New(config.DbDriverName, config.Dsn)
	if err := storage.Connect(); err != nil {
		logg.Error("failed to connect to DB: " + err.Error())
		os.Exit(1)
	}
	defer storage.Close() //nolint:errcheck

	producer := kafka.NewProducer(config.KafkaBrokers, config.KafkaTopic)
	defer producer.Close() //nolint:errcheck

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("calendar_scheduler is running...")

	ticker := time.NewTicker(config.ScanInterval.Duration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logg.Info("calendar_scheduler is stopping...")
			return
		case <-ticker.C:
			runScan(ctx, logg, storage, producer)
		}
	}
}

func runScan(ctx context.Context, logg *logger.Logger, storage *sqlstorage.SQLStorage, producer kafka.Producer) {
	metrics.RecordSchedulerRun()
	now := time.Now()

	events, err := storage.FindEventsToNotify(now)
	if err != nil {
		logg.Error("FindEventsToNotify error: " + err.Error())
		metrics.RecordSchedulerError()
	} else {
		for _, e := range events {
			metrics.RecordEventProcessed()
			n := models.Notification{
				EventID:   e.ID,
				Title:     e.Title,
				EventDate: e.DateTime,
				UserID:    e.UserID,
			}
			if err := producer.SendNotification(ctx, n); err != nil {
				logg.Error("SendNotification error: " + err.Error())
				metrics.RecordSchedulerError()
			} else {
				logg.Info("sent notification for event " + e.Title)
				metrics.RecordNotificationSent()
			}
		}
	}

	oneYearAgo := now.AddDate(-1, 0, 0)
	if err := storage.DeleteOldEvents(oneYearAgo); err != nil {
		logg.Error("DeleteOldEvents error: " + err.Error())
		metrics.RecordSchedulerError()
	}
}
