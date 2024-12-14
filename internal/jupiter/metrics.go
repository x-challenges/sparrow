package jupiter

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// MetricName
const metricName = "sparrow.io/jupiter"

// Providers
var (
	meter = otel.Meter(metricName)
)

var (
	// RequestLatency
	requestLatency metric.Int64Histogram
)

func init() {
	var err error

	// RequestLatency
	requestLatency, err = meter.Int64Histogram("jupiter.request.latency",
		metric.WithUnit("ms"),
	)

	if err != nil {
		panic(err)
	}
}
