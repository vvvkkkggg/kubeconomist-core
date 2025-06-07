package nodeoptimizer

import (
	"context"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/billing"
	"github.com/vvvkkkggg/kubeconomist-core/internal/model"
	"github.com/vvvkkkggg/kubeconomist-core/internal/yandex"
)

type NodeOptimizer struct {
	yandex  *yandex.Client
	billing *billing.Billing

	metric *prometheus.GaugeVec
}

func NewNodeOptimizer(ya *yandex.Client, b *billing.Billing) *NodeOptimizer {
	m := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kubeconomist",
			Subsystem: "node_optimizer",
			Name:      "node_optimization_status",
			Help:      "Status of node optimization",
		},
		[]string{"cloud_id", "folder_id", "node_group_id", "platform_id", "status"},
	)

	return &NodeOptimizer{
		yandex:  ya,
		billing: b,
		metric:  m,
	}
}

// TODO: Return metric to expose via Prometheus

func (n *NodeOptimizer) calculatePrice(platformID string, coreFraction, cores, memory int64) (float64, error) {
	currentPriceCPU, err := n.billing.GetPriceCPURUB(platformID, strconv.Itoa(int(coreFraction)), model.CPUCount(float64(cores)))
	if err != nil {
		return 0, err
	}

	currentPriceMemory, err := n.billing.GetPriceRAMRUB(platformID, model.RAMCount(float64(memory)))
	if err != nil {
		return 0, err
	}

	return float64(currentPriceCPU + currentPriceMemory), nil
}

func (n *NodeOptimizer) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		clouds, err := n.yandex.GetClouds(ctx)
		if err != nil {
			// TODO: handle error properly
			return err
		}

		for _, cloud := range clouds {
			folders, err := n.yandex.GetFolders(ctx, cloud.Id)
			if err != nil {
				// TODO: handle error properly
				return err
			}

			for _, folder := range folders {
				nodeGroups, err := n.yandex.GetNodeGroups(ctx, folder.Id)
				if err != nil {
					// TODO: handle error properly
					return err
				}

				for _, nodeGroup := range nodeGroups {
					coreFraction := nodeGroup.GetNodeTemplate().GetResourcesSpec().GetCoreFraction()
					cores := nodeGroup.GetNodeTemplate().GetResourcesSpec().GetCores()
					memory := nodeGroup.GetNodeTemplate().GetResourcesSpec().GetMemory()
					platformID := nodeGroup.GetNodeTemplate().GetPlatformId()

					currentPrice, err := n.calculatePrice(platformID, coreFraction, cores, memory)
					if err != nil {
						// TODO: handle error properly
						return err
					}

					cheapestPrice := currentPrice
					cheapestPlatformID := platformID
					// TODO: Get platfrom IDs from billing package
					for _, p := range []string{"standard-v1", "standard-v2", "standard-v3"} {
						if p == platformID {
							continue
						}

						price, err := n.calculatePrice(p, coreFraction, cores, memory)
						if err != nil {
							// TODO: handle error properly
							return err
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
}
