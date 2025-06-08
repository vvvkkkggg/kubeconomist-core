package vpc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/yandex"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/vpc/v1"
)

var _ analyzers.Analyzer = &VPCAnalyzer{}

type VPCAnalyzer struct {
	yandex *yandex.Client
	m      *prometheus.GaugeVec

	cloudID  string
	folderID string
}

type Address struct {
	CloudID    string
	FolderID   string
	IP         string
	IsUsed     bool
	IsReserved bool
}

func NewVPCAnalyzer(ya *yandex.Client, cloudID, folderID string) *VPCAnalyzer {
	// 241.05 рублей за 1 неиспользуемый IP адрес в месяц
	m := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kubeconomist",
			Subsystem: "vpc",
			Name:      "ip_info",
			Help:      "IP info",
		},
		[]string{"cloud_id", "folder_id", "ip_address", "is_used", "is_reserved"},
	)

	return &VPCAnalyzer{
		yandex: ya,
		m:      m,
	}
}

func (v *VPCAnalyzer) GetAddresses(ctx context.Context) ([]Address, error) {
	folders, err := v.yandex.GetAllFolders(ctx, v.cloudID, v.folderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all folders: %w", err)
	}

	addrs := make([]Address, 0)

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
				CloudID:    folder.CloudId,
				FolderID:   folder.Id,
				IP:         a.GetExternalIpv4Address().GetAddress(),
				IsUsed:     a.GetUsed(),
				IsReserved: a.GetReserved(),
			})
		}
	}

	return addrs, nil
}

func (v *VPCAnalyzer) Run(ctx context.Context) {
	addrs, err := v.GetAddresses(context.Background())
	if err != nil {
		return
	}

	for _, addr := range addrs {
		if !addr.IsReserved {
			continue
		}

		v.m.With(prometheus.Labels{
			"cloud_id":    addr.CloudID,
			"folder_id":   addr.FolderID,
			"ip_address":  addr.IP,
			"is_used":     strconv.Itoa(boolToInt(addr.IsUsed)),
			"is_reserved": strconv.Itoa(boolToInt(addr.IsReserved)),
		}).Set(1.0)
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (v *VPCAnalyzer) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		v.m,
	}
}
