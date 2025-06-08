package dnsoptimizer

import (
	"context"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/billing"
	"github.com/vvvkkkggg/kubeconomist-core/internal/yandex"
)

var _ analyzers.Analyzer = &DNSOptimizer{}

type DNSOptimizer struct {
	yandex  *yandex.Client
	billing *billing.Billing

	metric *prometheus.GaugeVec
}

func NewDNSOptimizer(ya *yandex.Client, b *billing.Billing) *DNSOptimizer {
	// Константный прайс 38.88 рублей за 1 DNS зону в месяц
	m := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kubeconomist",
			Subsystem: "dns_optimizer",
			Name:      "dns_optimization_zone_not_empty",
			Help:      "Is DNS zone not empty",
		},
		[]string{"cloud_id", "folder_id", "zone_id"},
	)

	return &DNSOptimizer{
		yandex:  ya,
		billing: b,
		metric:  m,
	}
}

func (n *DNSOptimizer) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		n.metric,
	}
}

func (n *DNSOptimizer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Error("ctx done", slog.Any("err", ctx.Err()))
			return
		default:
		}

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
				zones, err := n.yandex.GetDNSZones(ctx, folder.Id)
				if err != nil {
					slog.Error("get dns zones err", slog.Any("err", err))
					return
				}

				for _, zone := range zones {
					isUsed, err := n.yandex.IsDNSUsed(ctx, zone.Id)
					if err != nil {
						slog.Error("check dns zone empty err", slog.Any("err", err))
						return
					}

					v := 0.0
					if isUsed {
						v = 1.0
					}

					n.metric.With(prometheus.Labels{
						"cloud_id":  cloud.Id,
						"folder_id": folder.Id,
						"zone_id":   zone.Id,
					}).Set(v)
				}
			}
		}
	}
}
