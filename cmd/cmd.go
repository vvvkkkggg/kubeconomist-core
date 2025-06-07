package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vvvkkkggg/kubeconomist-core/internal/billing"
)

func Run() {
	billingClient := billing.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	skus, err := billingClient.GetPricesForKubernetes(ctx)
	if err != nil {
		log.Fatalf("Failed to get prices: %v", err)
	}

	fmt.Printf("Found %d active SKUs for Kubernetes:\n", len(skus))
	for _, sku := range skus {
		price, err := sku.GetCurrentPrice()
		if err != nil {
			log.Printf("Warning: couldn't get price for SKU %s: %v", sku.ID, err)
			continue
		}

		fmt.Printf("\n[%s] %s\n", sku.ID, sku.Name)
		fmt.Printf("Price: %s RUB/%s\n", price, sku.PricingUnit)

		if len(sku.PricingVersions) > 0 {
			effectiveTime := sku.PricingVersions[0].GetEffectiveTime()
			if !effectiveTime.IsZero() {
				fmt.Printf("Effective since: %s\n", effectiveTime.Format("2006-01-02"))
			}
		}
	}
}
