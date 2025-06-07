package vpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/yandex"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/vpc/v1"
)

const (
	pricePerIPPerHourUSD = 0.005 // Установленная ставка за 1 IP/час
)

var _ analyzers.Analyzer = &VPCAnalyzer{}

type VPCAnalyzer struct {
	yandex *yandex.Client
}

type Address struct {
	CloudID    string
	FolderID   string
	IP         string
	IsUsed     bool
	IsReserved bool
}

func NewVPCAnalyzer(ya *yandex.Client) *VPCAnalyzer {
	return &VPCAnalyzer{
		yandex: ya,
	}
}

func (v *VPCAnalyzer) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("kubeconomist_vpc_ip_status", "Status of IPs", []string{"ip_address"}, nil)
}

func (v *VPCAnalyzer) Collect(ch chan<- prometheus.Metric) {
	addrs, err := v.GetAddresses(context.Background())
	if err != nil {
		return
	}

	for _, addr := range addrs {
		if !addr.IsReserved {
			continue
		}

		isUsed := 0.0
		if addr.IsUsed {
			isUsed = 1.0
		}

		desc := prometheus.NewDesc("kubeconomist_vpc_ip_status", "Status of IPs", []string{"ip_address"}, nil)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, isUsed, addr.IP)
	}
}

// func (v *VPCAnalyzer) GetMetric() *prometheus.GaugeVec {
// 	return prometheus.NewGaugeVec(
// 		prometheus.GaugeOpts{
// 			Namespace: "kubeconomist",
// 			Subsystem: "vpc",
// 			Name:      "ip_status",
// 			Help:      "Status of IP addresses in VPC",
// 		},
// 		[]string{"ip_address"},
// 	)
// }

func (v *VPCAnalyzer) GetAddresses(ctx context.Context) ([]Address, error) {
	clouds, err := v.yandex.GetClouds(ctx)
	if err != nil {
		return nil, err
	}

	addrs := make([]Address, 0)
	for _, cloud := range clouds {
		folders, err := v.yandex.GetFolders(ctx, cloud.Id)
		if err != nil {
			return nil, err
		}

		for _, folder := range folders {
			as, err := v.yandex.GetAddresses(ctx, folder.Id)
			if err != nil {
				return nil, err
			}

			for _, a := range as {
				if a.GetIpVersion() != vpc.Address_IPV4 {
					continue
				}

				addrs = append(addrs, Address{
					CloudID:    cloud.Id,
					FolderID:   folder.Id,
					IP:         a.GetExternalIpv4Address().GetAddress(),
					IsUsed:     a.GetUsed(),
					IsReserved: a.GetReserved(),
				})
			}
		}
	}

	return addrs, nil
}

func (v *VPCAnalyzer) Run(ctx context.Context) {}

func (v *VPCAnalyzer) GetCollectors() []prometheus.Collector {
	return nil
}
