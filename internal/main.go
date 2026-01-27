package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Kenji-Uema/mongodbBootstrap/internal/domain"
	"github.com/Kenji-Uema/mongodbBootstrap/internal/infra"
	"github.com/Kenji-Uema/mongodbBootstrap/internal/tooling/log"
	"github.com/Kenji-Uema/mongodbBootstrap/internal/tooling/telemetry"
)

func main() {
	if err := run(); err != nil {
		slog.Error("bootstrap failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	startTime := time.Now()
	slog.SetDefault(log.NewLogger())
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info("Clock Emulation Starting")

	telemetryShutdown, err := telemetry.Init(ctx)
	if err != nil {
		return fmt.Errorf("telemetry init: %w", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := telemetryShutdown(shutdownCtx); err != nil {
			slog.Error("telemetry shutdown", "error", err)
		}
	}()

	wg := sync.WaitGroup{}
	errCh := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := infra.BootStrapPhotos(); err != nil {
			errCh <- fmt.Errorf("bootstrap photos: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		mongoDb, err := infra.NewMongoDb(ctx)
		if err != nil {
			errCh <- fmt.Errorf("mongo connection: %w", err)
			return
		}
		defer func() {
			closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := mongoDb.Close(closeCtx); err != nil {
				slog.Error("failed to close mongo connection", "error", err)
			}
		}()

		if err := mongoDb.DropAll(ctx); err != nil {
			errCh <- fmt.Errorf("drop database: %w", err)
			return
		}

		cottageCollection := mongoDb.NewCollection("Cottage")
		guestCollection := mongoDb.NewCollection("Guest")
		mongoDb.NewCollection("Booking")

		if err := infra.SetIndex(cottageCollection, "name"); err != nil {
			errCh <- fmt.Errorf("set cottage index: %w", err)
			return
		}
		if err := infra.SetIndex(guestCollection, "email"); err != nil {
			errCh <- fmt.Errorf("set guest index: %w", err)
			return
		}

		if err := infra.Seed[domain.Cottage](ctx, cottageCollection, "resources/cottages.json"); err != nil {
			errCh <- fmt.Errorf("seed cottages: %w", err)
			return
		}
	}()

	wg.Wait()
	close(errCh)

	jobDuration := time.Since(startTime)
	telemetry.JobDurationHistogram.Record(ctx, jobDuration.Milliseconds())

	var runErr error
	for err := range errCh {
		if err != nil {
			runErr = errors.Join(runErr, err)
		}
	}
	if runErr != nil {
		return runErr
	}
	return nil
}
