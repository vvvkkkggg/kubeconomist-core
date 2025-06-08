package krr

const (
	labelResourceType      = "resource_type"   // cpu, ram
	labelConsumptionType   = "requests_type"   // rub, real
	labelConsumptionStatus = "requests_status" // current, recommended, gain
	labelCluster           = "cluster"         // cluster name
	labelPodName           = "pod_name"
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

func (k *KrrAnalyzer) writeConsumptionToGauge(
	podname string,
	cluster string,
	resource Resource,
	unit ConsumptionMeasurementUnit,
	status ConsumptionStatus,
	amount float64,
) {
	k.resourceGauge.
		WithLabelValues(podname, cluster, string(resource), string(unit), string(status)).Set(amount)
}
