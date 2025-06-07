package service

import "github.com/vvvkkkggg/kubeconomist-core/internal/model"

type Billing interface {
	GetPriceRUB(cpuCount model.CPUCount, ramCount model.CPUCount) model.PriceRUB
}

type Service struct {
	billing Billing
}

func New(b *Billing) *Service {
	return nil
}
