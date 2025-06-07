package krrstub

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

func (k *KrrAnalyzer) writeConsumptionToGauge(
	resource Resource,
	unit ConsumptionMeasurementUnit,
	status ConsumptionStatus,
	amount float64,
) {
	k.resourceGauge.
		WithLabelValues(string(resource), string(unit), string(status)).Set(amount)
}
