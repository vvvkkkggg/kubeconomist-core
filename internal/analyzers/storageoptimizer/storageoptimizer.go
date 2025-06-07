package storageoptimizer

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/billing"
)

type S3Client interface {
	ListBuckets(ctx context.Context) ([]string, error)
	GetBucketLocation(ctx context.Context, bucket string) (string, error)
	GetBucketStorageClass(ctx context.Context, bucket string) (string, error)
}

type Billing interface {
	GetObjectStoragePricesRUB(ctx context.Context) (*billing.ObjectStoragePrices, error)
}

type StorageOptimizer struct {
	profitGauge *prometheus.GaugeVec
	s3Client    S3Client
	billing     Billing
}

func New() *StorageOptimizer {
	profitGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "s3_bucket_optimization_profit",
			Help: "Profit from optimizing S3 bucket storage class",
		},
		[]string{"bucket"},
	)

	return &StorageOptimizer{
		profitGauge: profitGauge,
	}
}

func (so *StorageOptimizer) Run(ctx context.Context) {
	_, err := so.billing.GetObjectStoragePricesRUB(ctx)

	buckets, err := so.s3Client.ListBuckets(ctx)
	if err != nil {
		// Обработка ошибки
		return
	}

	for _, bucket := range buckets {
		_ = bucket
	}
}

func (so *StorageOptimizer) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{}
}
