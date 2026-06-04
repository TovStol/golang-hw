package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/TovStol/hw12_13_14_15_calendar/internal/app"
	"github.com/TovStol/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/TovStol/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/TovStol/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/TovStol/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config",
		"/home/atovstolik/projects/golang-hw/hw12_13_14_15_16_calendar/configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	logg := logger.New(config.Level, config.Location)
	defer func() {
		if err := logg.Close(); err != nil {
			logg.Error("failed to close logger: " + err.Error())
		}
	}()
	var calendar *app.App
	if config.Storage == "inner" {
		storage := memorystorage.New()
		calendar = app.New(logg, storage)
	} else {
		storage := sqlstorage.New(config.DBDriverName, config.Dsn)
		// Retry connection with backoff
		var err error
		for i := 0; i < 10; i++ {
			err = storage.Connect()
			if err == nil {
				break
			}
			logg.Info("waiting for database... (attempt " + fmt.Sprintf("%d", i+1) + "/10)")
			time.Sleep(2 * time.Second)
		}
		if err != nil {
			logg.Error("failed to connect to database after retries: " + err.Error())
			return
		}
		calendar = app.New(logg, storage)
	}

	server := internalhttp.NewServer(logg, calendar, config.Host, config.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
