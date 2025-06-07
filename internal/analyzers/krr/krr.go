package analyzers

import (
	"context"

	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
)

var _ analyzers.Analyzer = &KrrAnalyzer{}

type KrrAnalyzer struct {
}

func NewKrrAnalyzer() *KrrAnalyzer {
	panic("implement me")
}

func (k *KrrAnalyzer) Run(ctx context.Context) {
	panic("implement me")
}
