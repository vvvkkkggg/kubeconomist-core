package service

import (
	"github.com/vvvkkkggg/kubeconomist-core/internal/model"
)

type Billing interface {
	GetPriceCPURUB(platform string, coreFraction string, cpuCount model.CPUCount) model.PriceRUB
	GetPriceRAMRUB(platform string, coreFraction string, cpuCount model.CPUCount) model.PriceRUB
}

type Service struct {
	billing Billing
}

func New(b *Billing) *Service {
	return nil
}
