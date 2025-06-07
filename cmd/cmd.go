package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers/krr"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers/vpc"
	"github.com/vvvkkkggg/kubeconomist-core/internal/billing"
	"github.com/vvvkkkggg/kubeconomist-core/internal/config"
	"github.com/vvvkkkggg/kubeconomist-core/internal/metrics"
	"github.com/vvvkkkggg/kubeconomist-core/internal/yandex"
)

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		return err
	}

	billing := billing.New()
	if err := billing.UpdatePricesCloudeCompute(ctx); err != nil {
		return fmt.Errorf("update price cloud compute compute price: %v", err)
	}
	if err := billing.UpdateObjectStoragePrice(ctx); err != nil {
		return fmt.Errorf("update price object storage: %v", err)
	}
	if err := billing.UpdatePricesContainerRegistry(ctx); err != nil {
		return fmt.Errorf("update price continaer registry: %v", err)
	}

	yandexClient, err := yandex.New(ctx, cfg.Analyzers.VPC.YCToken)
	if err != nil {
		return err
	}

	analyzerList := []analyzers.Analyzer{
		krr.NewKrrAnalyzer(billing, cfg.Analyzers.KRR),
		vpc.NewVPCAnalyzer(yandexClient),
		//registryoptimizer.NewRegistryOptimizer(billing),
	}

	var collectors []prometheus.Collector

	for _, a := range analyzerList {
		collectors = append(collectors, a.GetCollectors()...)
	}

	for _, a := range analyzerList {
		go a.Run(ctx)
	}

	slog.Info("analyzers are running")

	if err := metrics.ListenAndServe(
		ctx,
		cfg.Metrics,
		collectors...,
	); err != nil {
		return err
	}

	return nil
}
