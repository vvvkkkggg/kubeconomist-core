package cmd

import (
	"context"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers/krrstub"
	"github.com/vvvkkkggg/kubeconomist-core/internal/billing"
	"github.com/vvvkkkggg/kubeconomist-core/internal/config"
	"github.com/vvvkkkggg/kubeconomist-core/internal/metrics"
)

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		return err
	}

	billing := billing.New()

	analyzerList := []analyzers.Analyzer{
		krrstub.NewKrrAnalyzer(billing, cfg.Analyzers.KRR),
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
