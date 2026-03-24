package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kenji-Uema/bootstrap/internal/config"
	"github.com/Kenji-Uema/bootstrap/internal/infra/log"
	"github.com/Kenji-Uema/bootstrap/internal/infra/mdb"
	"github.com/Kenji-Uema/bootstrap/internal/infra/mq"
	"github.com/Kenji-Uema/bootstrap/internal/infra/telemetry"
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
	configs, err := config.LoadConfigs()
	exitOnError("config load", err)

	telemetryShutdown, err := telemetry.Init(baseCtx, configs.AppConfig, configs.TelemetryConfig)
	exitOnError("telemetry init", err)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
		defer cancel()
		if err := telemetryShutdown(shutdownCtx); err != nil {
			slog.Error("telemetry shutdown", "error", err)
		}
	}()

	err = mdb.BootstrapPhotos(configs.AppConfig)
	exitOnError("photos bootstrap", err)

	closeMongoConn, err := mdb.BootstrapMongodb(baseCtx, configs)
	exitOnError("mongodb bootstrap", err)
	defer closeMongoConn()

	closeRabbitChannel, closeRabbitConn, err := mq.BootstrapRabbitmq(baseCtx, configs)
	exitOnError("rabbitmq bootstrap", err)
	defer closeRabbitConn()
	defer closeRabbitChannel()

	jobDuration := time.Since(startTime)
	slog.Info("Bootstrap complete", "duration", jobDuration)
	telemetry.JobDurationHistogram.Record(baseCtx, jobDuration.Milliseconds())
}
