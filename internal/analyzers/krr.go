package analyzers

import "context"

var _ Analyzer = &KrrAnalyzer{}

type KrrAnalyzer struct {
}

func NewKrrAnalyzer() *KrrAnalyzer {
	panic("implement me")
}

func (k *KrrAnalyzer) Run(ctx context.Context) {
	panic("implement me")
}
