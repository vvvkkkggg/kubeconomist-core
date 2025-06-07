package krr

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"

	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/config"
	"github.com/vvvkkkggg/kubeconomist-core/internal/model"
)

var _ analyzers.Analyzer = &KrrAnalyzer{}

type KrrAnalyzer struct {
	billing   analyzers.Billing
	collector *Collector
	cfg       config.KrrAnalyzerConfig
}

func NewKrrAnalyzer(
	b analyzers.Billing,
	collector *Collector,
	cfg config.KrrAnalyzerConfig,
) *KrrAnalyzer {
	return &KrrAnalyzer{
		billing:   b,
		collector: collector,
		cfg:       cfg,
	}
}

func (k *KrrAnalyzer) Run(ctx context.Context) {
	panic("implement me")
}

type ResourceOptimization struct {
	Cluster   string
	Namespace string
	PodName   string
	PodCount  uint
	PodType   string
	Container string
	CPUReqOld *model.CPUCount // e.g. 100m → 0.1
	CPUReqNew *model.CPUCount // e.g. 50m  → 0.05
	RAMReqOld *model.CPUCount // e.g. 512Mi → 0.5 (GiB-based or however your model interprets it)
	RAMReqNew *model.CPUCount // e.g. 256Mi → 0.25
	CPULimOld *model.CPUCount // e.g. 100m → 0.1
	CPULimNew *model.CPUCount // e.g. 50m  → 0.05
	RAMLimOld *model.CPUCount // e.g. 512Mi → 0.5 (GiB-based or however your model interprets it)
	RAMLimNew *model.CPUCount // e.g. 256Mi → 0.25
}

// CalculatePrice iterates over each container’s old vs. new requests,
// asks Billing for their ruble cost, and accumulates totals.
// Returns (currentTotal, optimizedTotal, gain).
func (k *KrrAnalyzer) CalculatePrice(rows []krrOutput) {
	for _, r := range rows {
		sendMetrics := func(prev, new float64, resource Resource, unit ConsumptionMeasurementUnit) {
			k.collector.AddResourceConsumption(
				resource, unit, ConsumptionStatusGain,
				prev-new,
			)

			k.collector.AddResourceConsumption(
				resource, unit, ConsumptionStatusCurrent,
				prev,
			)

			k.collector.AddResourceConsumption(
				resource, unit, ConsumptionStatusRecommended,
				new,
			)
		}

		if r.Object.Requests.CPU != nil && *r.Object.Requests.CPU > *r.Recommended.Requests.CPU {
			sendMetrics(
				*r.Object.Requests.CPU, *r.Recommended.Requests.CPU,
				ResourceCPU,
				ConsumptionReal,
			)

			sendMetrics(
				*r.Object.Requests.CPU, *r.Recommended.Requests.CPU, // todo: add money multiplier
				ResourceCPU,
				ConsumptionMoney,
			)
		}

		if r.Object.Requests.Memory != nil && *r.Object.Requests.Memory > *r.Recommended.Requests.Memory {
			sendMetrics(
				*r.Object.Requests.Memory, *r.Recommended.Requests.Memory,
				ResourceRAM,
				ConsumptionReal,
			)

			sendMetrics(
				*r.Object.Requests.Memory, *r.Recommended.Requests.Memory, // todo: add money multiplier
				ResourceRAM,
				ConsumptionMoney,
			)
		}
	}

	return
}

func (k *KrrAnalyzer) callKRR() (krrOutput, error) {
	cmd := exec.Command(
		"krr", "simple",
		"-p", k.cfg.PrometheusURL,
		"--prometheus-auth-header", k.cfg.PrometheusAuthHeader,
		"--history-duration", k.cfg.HistoryDuration,
		"-f", "json")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return krrOutput{}, err
	}

	resultJSON := stdout.Bytes()

	var krrResult krrOutput
	if err := json.Unmarshal(resultJSON, &krrResult); err != nil {
		return krrOutput{}, err
	}

	return krrResult, nil
}
