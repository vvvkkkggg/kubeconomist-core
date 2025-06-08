package platformoptimizer

import (
	"context"
	"errors"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/billing"
	"github.com/vvvkkkggg/kubeconomist-core/internal/yandex"
)

var _ analyzers.Analyzer = &PlatformOptimizer{}

type PlatformOptimizer struct {
	yandex  *yandex.Client
	billing *billing.Billing

	metric *prometheus.GaugeVec
}

func NewPlatformOptimizer(ya *yandex.Client, b *billing.Billing) *PlatformOptimizer {
	m := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kubeconomist",
			Subsystem: "platform_optimizer",
			Name:      "platform_optimizer_price",
			Help:      "Price of the node",
		},
		[]string{"cloud_id", "folder_id", "node_group_id", "platform_id", "status"},
	)

	return &PlatformOptimizer{
		yandex:  ya,
		billing: b,
		metric:  m,
	}
}

func (n *PlatformOptimizer) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		n.metric,
	}
}

func (n *PlatformOptimizer) Run(ctx context.Context) {
	clouds, err := n.yandex.GetClouds(ctx)
	if err != nil {
		slog.Error("get clouds err", slog.Any("err", err))
		return
	}

	for _, cloud := range clouds {
		folders, err := n.yandex.GetFolders(ctx, cloud.Id)
		if err != nil {
			slog.Error("get folders err", slog.Any("err", err))
			return
		}

		for _, folder := range folders {
			nodeGroups, err := n.yandex.GetNodeGroups(ctx, folder.Id)
			if err != nil {
				slog.Error("get node groups err", slog.Any("err", err))
				return
			}

			for _, nodeGroup := range nodeGroups {
				coreFraction := nodeGroup.GetNodeTemplate().GetResourcesSpec().GetCoreFraction()
				cores := nodeGroup.GetNodeTemplate().GetResourcesSpec().GetCores()
				memory := nodeGroup.GetNodeTemplate().GetResourcesSpec().GetMemory()
				platformID := nodeGroup.GetNodeTemplate().GetPlatformId()

				currentPrice, err := n.billing.CalculatePrice(platformID, coreFraction, cores, memory)
				if err != nil {
					slog.Error("calculate price err", slog.Any("err", err))
					return
				}

				cheapestPrice := currentPrice
				cheapestPlatformID := platformID
				// TODO: Get platfrom IDs from billing package
				for _, p := range []string{"standard-v1", "standard-v2", "standard-v3"} {
					if p == platformID {
						continue
					}

					price, err := n.billing.CalculatePrice(p, coreFraction, cores, memory)
					if err != nil {
						if errors.Is(err, billing.ErrFlavourNotFound) {
							continue
						}
						slog.Error("calculate price err", slog.Any("err", err))
						return
					}

					if price < currentPrice {
						cheapestPrice = price
						cheapestPlatformID = p
					}
				}

				n.metric.With(prometheus.Labels{
					"cloud_id":      cloud.Id,
					"folder_id":     folder.Id,
					"node_group_id": nodeGroup.Id,
					"platform_id":   platformID,
					"status":        "current",
				}).Set(float64(currentPrice))

				n.metric.With(prometheus.Labels{
					"cloud_id":      cloud.Id,
					"folder_id":     folder.Id,
					"node_group_id": nodeGroup.Id,
					"platform_id":   cheapestPlatformID,
					"status":        "desired",
				}).Set(float64(cheapestPrice))
			}
		}
	}
}
