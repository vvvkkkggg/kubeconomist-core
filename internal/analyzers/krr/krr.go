package krr

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/config"
)

var _ analyzers.Analyzer = &KrrAnalyzer{}

type KrrAnalyzer struct {
	billing analyzers.Billing
	cfg     config.KrrAnalyzerConfig

	resourceGauge *prometheus.GaugeVec
}

func NewKrrAnalyzer(
	b analyzers.Billing,
	cfg config.KrrAnalyzerConfig,
) *KrrAnalyzer {
	return &KrrAnalyzer{
		billing: b,
		cfg:     cfg,

		resourceGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "kubeconomist",
				Subsystem: "krr",
				Name:      "resource_consumption",
				Help:      "A histogram of resource consumption by k8s cluster",
			},
			[]string{labelResourceType, labelConsumptionType, labelConsumptionStatus},
		),
	}
}

func (k *KrrAnalyzer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		krrStats, err := k.callKRR()
		if err != nil {
			panic(err)
		}

		k.CalculatePrice(krrStats)
	}
}

// CalculatePrice iterates over each container’s old vs. new requests,
// asks Billing for their ruble cost, and accumulates totals.
// Returns (currentTotal, optimizedTotal, gain).
func (k *KrrAnalyzer) CalculatePrice(rows []KrrOutput) {
	for _, r := range rows {
		sendMetrics := func(prev, new float64, resource Resource, unit ConsumptionMeasurementUnit) {
			k.writeConsumptionToGauge(
				resource, unit, ConsumptionStatusGain,
				prev-new,
			)

			k.writeConsumptionToGauge(
				resource, unit, ConsumptionStatusCurrent,
				prev,
			)

			k.writeConsumptionToGauge(
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

func (k *KrrAnalyzer) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{k.resourceGauge}
}

func (k *KrrAnalyzer) callKRR() ([]KrrOutput, error) {
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
		return nil, err
	}

	resultJSON := stdout.Bytes()

	var krrResult []KrrOutput
	if err := json.Unmarshal(resultJSON, &krrResult); err != nil {
		return nil, err
	}

	return krrResult, nil
}
