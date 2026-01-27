package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	meter                = otel.Meter("clock-emulator")
	JobDurationHistogram metric.Int64Histogram
)

func initMetrics() error {
	var err error
	JobDurationHistogram, err = meter.Int64Histogram(
		"app.startup.duration",
		metric.WithUnit("ms"),
		metric.WithDescription("Time from process start to ready state"),
	)
	if err != nil {
		return err
	}

	return nil
}
