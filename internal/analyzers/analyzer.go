package analyzers

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type Analyzer interface {
	Run(context.Context)
	GetCollectors() *prometheus.Collector
}
