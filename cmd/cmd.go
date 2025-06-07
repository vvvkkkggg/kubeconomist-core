package cmd

import (
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers/krr"
	"github.com/vvvkkkggg/kubeconomist-core/internal/billing"
	"github.com/vvvkkkggg/kubeconomist-core/internal/config"
)

func Run() error {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		return err
	}

	var analyzerList []analyzers.Analyzer

	_ = cfg

	billing := billing.New()

	collector, err := krr.NewCollector("krr")
	if err != nil {
		return err
	}

	krr.NewKrrAnalyzer(billing, collector)

	panic("implement me")

	return nil
}
