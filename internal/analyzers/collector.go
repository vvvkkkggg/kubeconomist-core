package analyzers

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	labelResourceType      = "resource_type"
	labelParameterType     = "parameter_type"
	labelConsumptionType   = "consumption_type"
	labelConsumptionStatus = "consumption_status"
)

type (
	Resource                   string
	ConsumptionMeasurementUnit string
	ParameterType              string
	ConsumptionStatus          string
)

const (
	ResourceCPU Resource = "cpu"
	ResourceRAM Resource = "ram"

	ParameterTypeRequest ParameterType = "requests"
	ParameterTypeLimit   ParameterType = "limits"

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
			[]string{labelResourceType, labelParameterType, labelConsumptionType, labelConsumptionStatus},
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
	parameterType ParameterType,
	unit ConsumptionMeasurementUnit,
	status ConsumptionStatus,
	amount float64,
) {
	c.resourceHistogram.
		WithLabelValues(string(resource), string(parameterType), string(unit), string(status)).
		Observe(amount)
}
