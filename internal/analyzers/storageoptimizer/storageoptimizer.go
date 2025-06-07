package storageoptimizer

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	rpsThresholdHot = 200
)

type StorageOptimizer struct {
}

func New() *StorageOptimizer {
	return &StorageOptimizer{}
}

func (so *StorageOptimizer) Run(ctx context.Context) {

}

func (so *StorageOptimizer) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{}
}
