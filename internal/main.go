package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kenji-Uema/bootstrap/internal/config"
	"github.com/Kenji-Uema/bootstrap/internal/domain"
	"github.com/Kenji-Uema/bootstrap/internal/infra"
	"github.com/Kenji-Uema/bootstrap/internal/tooling/log"
	"github.com/Kenji-Uema/bootstrap/internal/tooling/telemetry"
)

func exitOnError(errMsg string, err error) {
	if err != nil {
		slog.Error(errMsg, "error", err)
		os.Exit(1)
	}
}

func main() {
	startTime := time.Now()
	slog.SetDefault(log.NewLogger())
	baseCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("Bootstrap starting")

	appConfig, err := config.LoadConfig[config.AppConfig]()
	exitOnError("load app config", err)
	telemetryConfig, err := config.LoadConfig[config.TelemetryConfig]()
	exitOnError("load telemetry config", err)
	photosVolumeConfig, err := config.LoadConfig[config.PhotosVolumeConfig]()
	exitOnError("load photos volume config", err)
	mongoConfig, err := config.LoadConfig[config.MongoConfig]()
	exitOnError("load mongo config", err)

	telemetryShutdown, err := telemetry.Init(baseCtx, appConfig, telemetryConfig)
	exitOnError("telemetry init", err)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
		defer cancel()
		if err := telemetryShutdown(shutdownCtx); err != nil {
			slog.Error("telemetry shutdown", "error", err)
		}
	}()

	err = infra.BootstrapPhotos(photosVolumeConfig)
	exitOnError("photos bootstrap", err)

	mongoDb, err := infra.NewMongoDB(baseCtx, mongoConfig)
	exitOnError("mongo init", err)
	defer func() {
		closeCtx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
		defer cancel()
		if err := mongoDb.Close(closeCtx); err != nil {
			slog.Error("failed to close mongo connection", "error", err)
		}
	}()

	err = mongoDb.DropAll(baseCtx)
	exitOnError("mongo drop", err)

	cottageCollection := mongoDb.NewCollection("Cottage")
	guestCollection := mongoDb.NewCollection("Guest")
	mongoDb.NewCollection("Booking")

	err = infra.SetIndex(baseCtx, cottageCollection, "name")
	exitOnError("set cottage index", err)
	err = infra.SetIndex(baseCtx, guestCollection, "email")
	exitOnError("set guest index", err)

	err = infra.Seed[domain.Cottage](baseCtx, cottageCollection, "resources/cottages.json")
	exitOnError("seed cottage", err)

	jobDuration := time.Since(startTime)
	telemetry.JobDurationHistogram.Record(baseCtx, jobDuration.Milliseconds())
}
