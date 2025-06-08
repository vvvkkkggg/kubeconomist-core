package storageoptimizer

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/billing"
	"github.com/vvvkkkggg/kubeconomist-core/internal/prom"
	"github.com/vvvkkkggg/kubeconomist-core/internal/yandex"
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
	price        *prometheus.GaugeVec
	storageClass *prometheus.GaugeVec
	yandex       *yandex.Client
	billing      *billing.Billing
}

func NewStorageOptimizer(yandex *yandex.Client, billing *billing.Billing) *StorageOptimizer {
	price := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kubeconomist",
			Subsystem: "storage_optimizer",
			Name:      "storage_optimization_price",
			Help:      "Price of the storage",
		},
		[]string{"cloud_id", "folder_id", "bucket_name", "status"},
	)

	storageClass := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kubeconomist",
			Subsystem: "storage_optimizer",
			Name:      "storage_optimization_class_is_optimal",
			Help:      "Desired storage class for the bucket",
		},
		[]string{"cloud_id", "folder_id", "bucket_name", "storage_class"},
	)

	return &StorageOptimizer{
		price:        price,
		storageClass: storageClass,
		yandex:       yandex,
		billing:      billing,
	}
}

func (so *StorageOptimizer) Run(ctx context.Context) {
	prices, err := so.billing.GetObjectStoragePricesRUB(ctx)
	if err != nil {
		slog.Error("get object storage prices err", slog.Any("err", err))
		return
	}

	for {
		select {
		case <-ctx.Done():
			slog.Error("ctx done", slog.Any("err", ctx.Err()))
			return
		default:
		}

		clouds, err := so.yandex.GetClouds(ctx)
		if err != nil {
			slog.Error("get clouds err", slog.Any("err", err))
			return
		}

		for _, cloud := range clouds {
			folders, err := so.yandex.GetFolders(ctx, cloud.Id)
			if err != nil {
				slog.Error("get folders err", slog.Any("err", err))
				return
			}

			for _, folder := range folders {
				buckets, err := so.yandex.GetBuckets(ctx, folder.Id)
				if err != nil {
					slog.Error("get buckets err", slog.Any("err", err))
					return
				}

				for _, bucket := range buckets {
					// TODO: hardcode localhost:8428, should be configurable
					queriesCount, err := prom.QueryValue("http://localhost:8428", fmt.Sprintf("sum(increase(rps{resource_id=\"%s\"}[24h]))", bucket.GetName()))
					if err != nil {
						slog.Error("failed to query Prometheus for bucket RPS", slog.Any("err", err), slog.String("bucket", bucket.GetName()))
						continue
					}

					spaceUsage, err := prom.QueryValue("http://localhost:8428", fmt.Sprintf("space_usage{resource_id=\"%s\",resource_type=\"bucket\",storage_type=\"All\"}", bucket.GetName()))
					if err != nil {
						slog.Error("failed to query Prometheus for bucket space usage", slog.Any("err", err), slog.String("bucket", bucket.GetName()))
						continue
					}

					spaceUsageGiB := spaceUsage / 1024 / 1024 / 1024

					storageClass := bucket.GetDefaultStorageClass()

					currentPrice := -1.0
					minPrice := math.MaxFloat64
					minStorageClass := storageClass
					// TODO: do not hardcode
					for _, sc := range []string{"STANDARD_IA", "COLD", "ICE"} {
						// TODO: queries count shoud be divided 10k or 1k. But skip for now.
						// TODO: queries should count for each operation type, lets get only POST for now.
						var price float64
						switch sc {
						case "STANDARD_IA":
							price = prices.StoragePrices.Standard*spaceUsageGiB + (queriesCount/1000.0)*prices.StandardOperations.Post
						case "COLD":
							price = prices.StoragePrices.Cold*spaceUsageGiB + (queriesCount/1000.0)*prices.ColdOperations.Post
						case "ICE":
							price = prices.StoragePrices.Ice*spaceUsageGiB + (queriesCount/1000.0)*prices.IceOperations.Post
						default:
							panic("здорово ты это придумал")
						}

						if price < minPrice {
							minPrice = price
							minStorageClass = sc
						}

						if sc == storageClass {
							currentPrice = price
						}
					}

					so.price.With(prometheus.Labels{
						"cloud_id":    cloud.Id,
						"folder_id":   folder.Id,
						"bucket_name": bucket.GetName(),
						"status":      "current",
					}).Set(currentPrice)

					so.price.With(prometheus.Labels{
						"cloud_id":    cloud.Id,
						"folder_id":   folder.Id,
						"bucket_name": bucket.GetName(),
						"status":      "desired",
					}).Set(minPrice)

					for _, storageClass := range []string{"STANDARD_IA", "COLD", "ICE"} {
						indicator := 0.0
						if storageClass == minStorageClass {
							indicator = 1.0
						}

						so.storageClass.With(prometheus.Labels{
							"cloud_id":      cloud.Id,
							"folder_id":     folder.Id,
							"bucket_name":   bucket.GetName(),
							"storage_class": storageClass,
						}).Set(indicator)
					}
				}
			}
		}
	}
}

func (so *StorageOptimizer) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		so.price,
		so.storageClass,
	}
}
