package instruments

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// MetricName
const metricName = "sparrow.io/instruments"

// Providers
var (
	meter = otel.Meter(metricName)
)

// Metrics
var (
	instrumentsCounter metric.Int64Counter
)

func init() {
	var err error

	// init instruments counter
	instrumentsCounter, err = meter.Int64Counter("instruments",
		metric.WithDescription("The number of instruments"),
		metric.WithUnit("{instrument}"),
	)

	if err != nil {
		panic(err)
	}
}
