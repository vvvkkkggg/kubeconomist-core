package krr

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	labelResourceType      = "resource_type"
	labelConsumptionType   = "consumption_type"
	labelConsumptionStatus = "consumption_status"
)

type (
	Resource                   string
	ConsumptionMeasurementUnit string
	ConsumptionStatus          string
)

const (
	ResourceCPU Resource = "cpu"
	ResourceRAM Resource = "ram"

	ConsumptionMoney ConsumptionMeasurementUnit = "rub"
	ConsumptionReal  ConsumptionMeasurementUnit = "real"

	ConsumptionStatusCurrent     ConsumptionStatus = "current"
	ConsumptionStatusRecommended ConsumptionStatus = "recommended"
	ConsumptionStatusGain        ConsumptionStatus = "gain"
)

type Collector struct {
	resourceHistogram *prometheus.HistogramVec
}

func New(
	namespace, serviceName string,
) (*Collector, error) {
	return &Collector{
		resourceHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: serviceName,
				Name:      "resource_consumption",
				Help:      "A histogram of resource consumption by k8s cluster",
			},
			[]string{labelResourceType, labelConsumptionType, labelConsumptionStatus},
		),
	}, nil
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.resourceHistogram.Describe(ch)
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.resourceHistogram.Collect(ch)
}

func (c *Collector) AddResourceConsumption(
	resource Resource,
	unit ConsumptionMeasurementUnit,
	status ConsumptionStatus,
	amount float64,
) {
	c.resourceHistogram.
		WithLabelValues(string(resource), string(unit), string(status)).
		Observe(amount)
}
