package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	collectors []prometheus.Collector
}

func NewCollector(collectors []prometheus.Collector) (*Collector, error) {
	return &Collector{collectors: collectors}, nil
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range c.collectors {
		collector.Describe(ch)
	}
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	for _, collector := range c.collectors {
		collector.Collect(ch)
	}
}
