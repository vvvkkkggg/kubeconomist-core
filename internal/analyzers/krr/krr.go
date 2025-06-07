package analyzers

import (
	"context"

	"github.com/vvvkkkggg/kubeconomist-core/internal/model"

	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
)

var _ analyzers.Analyzer = &KrrAnalyzer{}

type KrrAnalyzer struct {
	billing   Billing
	collector *monitoring.Collector
}

func NewKrrAnalyzer(b Billing, collector *monitoring.Collector) *KrrAnalyzer {
	return &KrrAnalyzer{billing: b, collector: collector}
}

func (k *KrrAnalyzer) Run(ctx context.Context) {
	panic("implement me")
}

type ResourceOptimization struct {
	CPUReqOld model.CPUCount // e.g. 100m → 0.1
	CPUReqNew model.CPUCount // e.g. 50m  → 0.05
	RAMReqOld model.CPUCount // e.g. 512Mi → 0.5 (GiB-based or however your model interprets it)
	RAMReqNew model.CPUCount // e.g. 256Mi → 0.25
}

// CalculatePrice iterates over each container’s old vs. new requests,
// asks Billing for their ruble cost, and accumulates totals.
// Returns (currentTotal, optimizedTotal, gain).
func (k *KrrAnalyzer) CalculatePrice(rows []ResourceOptimization) (
	currentTotal model.PriceRUB,
	optimizedTotal model.PriceRUB,
	gain model.PriceRUB,
) {
	for _, r := range rows {
		// cost with “old” requests:
		curr := k.billing.GetPriceRUB(r.CPUReqOld, r.RAMReqOld)

		// cost with “new” (optimized) requests:
		opt := k.billing.GetPriceRUB(r.CPUReqNew, r.RAMReqNew)

		currentTotal += curr
		optimizedTotal += opt

		k.collector.AddResourceConsumption()
	}

	gain = currentTotal - optimizedTotal
	return
}
