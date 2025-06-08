package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	dnsoptimizer "github.com/vvvkkkggg/kubeconomist-core/internal/analyzers/dns"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers/krr"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers/nodeoptimizer"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers/platformoptimizer"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers/registryoptimizer"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers/storageoptimizer"
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

	// FIXME: ВОТ ТУТ ОТКЛЮЧАТЬ АНАЛАЙЗЕРЫ ДЛЯ ДЕБАГА
	analyzerList := []analyzers.Analyzer{
		krr.NewKrrAnalyzer(billing, cfg.Analyzers.KRR),
		vpc.NewVPCAnalyzer(yandexClient, cfg.Analyzers.CloudID, cfg.Analyzers.FolderID),
		platformoptimizer.NewPlatformOptimizer(yandexClient, billing, cfg.Analyzers.CloudID, cfg.Analyzers.FolderID),
		registryoptimizer.NewRegistryOptimizer(billing, cfg),
		nodeoptimizer.NewNodeOptimizer(yandexClient, billing, cfg.Analyzers.CloudID, cfg.Analyzers.FolderID),
		dnsoptimizer.NewDNSOptimizer(yandexClient, billing, cfg.Analyzers.CloudID, cfg.Analyzers.FolderID),
		storageoptimizer.NewStorageOptimizer(yandexClient, billing, cfg.Analyzers.CloudID, cfg.Analyzers.FolderID),
	}

	var collectors []prometheus.Collector

	for _, a := range analyzerList {
		collectors = append(collectors, a.GetCollectors()...)
	}

	for _, a := range analyzerList {
		go func() {
			for {
				slog.Info("running analyzer", slog.String("analyzer", reflect.TypeOf(a).String()))

				select {
				case <-ctx.Done():
					return
				default:
				}

				a.Run(ctx)

				time.Sleep(10 * time.Second)
			}

		}()
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
