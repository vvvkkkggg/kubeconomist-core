package service

import "github.com/vvvkkkggg/kubeconomist-core/internal/model"

type ResourceOptimization struct {
	CPUReqOld model.CPUCount // e.g. 100m → 0.1
	CPUReqNew model.CPUCount // e.g. 50m  → 0.05
	RAMReqOld model.CPUCount // e.g. 512Mi → 0.5 (GiB-based or however your model interprets it)
	RAMReqNew model.CPUCount // e.g. 256Mi → 0.25
}

type Billing interface {
	GetPriceRUB(cpuCount model.CPUCount, ramCount model.CPUCount) model.PriceRUB
}

type Service struct {
	billing Billing
}

func New(b *Billing) *Service {
	return nil
}

// CalculatePrice iterates over each container’s old vs. new requests,
// asks Billing for their ruble cost, and accumulates totals.
// Returns (currentTotal, optimizedTotal, gain).
func (s *Service) CalculatePrice(rows []ResourceOptimization) (
	currentTotal model.PriceRUB,
	optimizedTotal model.PriceRUB,
	gain model.PriceRUB,
) {
	for _, r := range rows {
		// cost with “old” requests:
		curr := s.billing.GetPriceRUB(r.CPUReqOld, r.RAMReqOld)

		// cost with “new” (optimized) requests:
		opt := s.billing.GetPriceRUB(r.CPUReqNew, r.RAMReqNew)

		currentTotal += curr
		optimizedTotal += opt
	}

	gain = currentTotal - optimizedTotal
	return
}
