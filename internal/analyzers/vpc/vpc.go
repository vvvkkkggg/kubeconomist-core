package analyzers

import (
	"context"
	"log"
	"os"

	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/yandex"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/vpc/v1"
)

const (
	pricePerIPPerHourUSD = 0.005 // Установленная ставка за 1 IP/час
)

var _ analyzers.Analyzer = &VPCAnalyzer{}

type VPCAnalyzer struct {
}

func NewVPCAnalyzer() *VPCAnalyzer {
	panic("implement me")
}

func (v *VPCAnalyzer) Run(ctx context.Context) {
	yandex, err := yandex.New(ctx, os.Getenv("YANDEX_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create Yandex client: %v", err)
	}

	clouds, err := yandex.GetClouds(ctx)
	if err != nil {
		log.Fatalf("Failed to get clouds: %v", err)
	}

	unsedIPs := make([]*vpc.Address, 0)
	for _, cloud := range clouds {
		folders, err := yandex.GetFolders(ctx, cloud.Id)
		if err != nil {
			log.Fatalf("Failed to get folders for cloud %s: %v", cloud.Id, err)
		}

		for _, folder := range folders {
			addresses, err := yandex.GetAddresses(ctx, folder.Id)
			if err != nil {
				log.Fatalf("Failed to get addresses for folder %s: %v", folder.Id, err)
			}

			for _, address := range addresses {
				if address.GetIpVersion() != vpc.Address_IPV4 {
					continue
				}

				if address.GetUsed() {
					continue
				}

				if !address.GetReserved() {
					continue
				}

				unsedIPs = append(unsedIPs, address)
			}
		}
	}
}
