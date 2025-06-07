package analyzers

import "github.com/vvvkkkggg/kubeconomist-core/internal/model"

type Billing interface {
	GetPriceCPURUB(platform string, coreFraction string, cpuCount model.CPUCount) model.PriceRUB
	GetPriceRAMRUB(platform string, ramCount model.RAMCount) model.PriceRUB
}
