package service

import (
	"github.com/vvvkkkggg/kubeconomist-core/internal/model"
)

type Billing interface {
	GetPriceCPURUB(platform string, coreFraction string, cpuCount model.CPUCount) (model.PriceRUB, error)
	GetPriceRAMRUB(platform string, coreFraction string, ramCount model.RAMCount) (model.PriceRUB, error)
}

type Service struct {
	billing Billing
}

func New(b *Billing) *Service {
	return nil
}
