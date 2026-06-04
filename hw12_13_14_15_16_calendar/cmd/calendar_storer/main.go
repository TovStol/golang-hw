package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/kafka"
	"github.com/TovStol/hw12_13_14_15_calendar/internal/logger"
	"github.com/TovStol/hw12_13_14_15_calendar/internal/metrics"
	sqlstorage "github.com/TovStol/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/storer_config.toml", "Path to configuration file")
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

	consumer := kafka.NewConsumer(config.KafkaBrokers, config.KafkaTopic, config.KafkaGroupID)
	defer consumer.Close() //nolint:errcheck

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("calendar_storer is running...")

	for {
		n, err := consumer.ReadNotification(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, io.EOF) {
				break
			}
			logg.Error("ReadNotification error: " + err.Error())
			metrics.RecordStorerError()
			continue
		}

		metrics.RecordNotificationReceived()
		if err := storage.SaveNotification(n); err != nil {
			logg.Error("SaveNotification error: " + err.Error())
			metrics.RecordStorerError()
		} else {
			logg.Info("stored notification for event " + n.Title)
			metrics.RecordNotificationStored()
		}
	}

	logg.Info("calendar_storer is stopping...")
}
