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
	"github.com/Kenji-Uema/bootstrap/internal/infra/log"
	"github.com/Kenji-Uema/bootstrap/internal/infra/mq"
	telemetry2 "github.com/Kenji-Uema/bootstrap/internal/infra/telemetry"
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

	telemetryShutdown, err := telemetry2.Init(baseCtx, configs.AppConfig, configs.TelemetryConfig)
	exitOnError("telemetry init", err)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(baseCtx, 5*time.Second)
		defer cancel()
		if err := telemetryShutdown(shutdownCtx); err != nil {
			slog.Error("telemetry shutdown", "error", err)
		}
	}()

	err = infra.BootstrapPhotos(configs.PhotosVolumeConfig)
	exitOnError("photos bootstrap", err)

	mongoDb, err := infra.NewMongoDB(baseCtx, configs.MongoConfig)
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

	rabbitMqClient, err := mq.NewRabbitMqConnection(baseCtx, configs.RabbitMqConfig)
	exitOnError("rabbitmq init", err)
	defer func() {
		if err := rabbitMqClient.Close(); err != nil {
			slog.Error("failed to close rabbitmq connection", "error", err)
		}
	}()

	cottageCollection := mongoDb.NewCollection("Cottage")
	guestCollection := mongoDb.NewCollection("Guest")
	mongoDb.NewCollection("Booking")

	rabbitmqProducer, err := mq.NewRabbitmqProducer(rabbitMqClient)
	exitOnError("rabbitmq producer init", err)

	err = infra.SetIndex(baseCtx, cottageCollection, "name")
	exitOnError("set cottage index", err)
	err = infra.SetIndex(baseCtx, guestCollection, "email")
	exitOnError("set guest index", err)

	err = infra.Seed[domain.Cottage](baseCtx, cottageCollection, "resources/cottages.json")
	exitOnError("seed cottage", err)

	exchanges := []config.ExchangeConfig{
		configs.CleaningExchangeConfig.Exchange,
		configs.TimeEventExchangeConfig.Exchange,
		configs.InvoiceExchangeConfig.Exchange,
		configs.PaymentExchangeConfig.Exchange,
		configs.NotificationExchangeConfig.Exchange,
	}

	for _, exchange := range exchanges {
		err = rabbitmqProducer.DeclareExchange(exchange)
		exitOnError("rabbitmq producer declare exchange", err)
	}

	err = rabbitmqProducer.CloseChannel()
	exitOnError("rabbitmq producer close channel", err)
	err = rabbitMqClient.Close()
	exitOnError("rabbitmq client close", err)

	jobDuration := time.Since(startTime)
	slog.Info("Bootstrap complete", "duration", jobDuration)
	telemetry2.JobDurationHistogram.Record(baseCtx, jobDuration.Milliseconds())
}
