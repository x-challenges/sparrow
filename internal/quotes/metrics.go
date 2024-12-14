package quotes

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// MetricName
const metricName = "sparrow.io/quotes"

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
	requestLatency, err = meter.Int64Histogram("quotes.request.latency",
		metric.WithUnit("ms"),
	)

	if err != nil {
		panic(err)
	}
}
