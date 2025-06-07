package cmd

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers/krr"
	"github.com/vvvkkkggg/kubeconomist-core/internal/billing"
	"github.com/vvvkkkggg/kubeconomist-core/internal/config"
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
		krr.NewKrrAnalyzer(billing, cfg.Analyzers.KRR),
	}

	var collectors []prometheus.Collector

	for _, a := range analyzerList {
		collectors = append(collectors, a.GetCollectors()...)
	}

	for _, a := range analyzerList {
		go a.Run(ctx)
	}

	return nil
}
