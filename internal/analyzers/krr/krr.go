package krr

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/config"
	"github.com/vvvkkkggg/kubeconomist-core/internal/model"
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
				Name:      "resource_recommendation",
				Help:      "A recommendation of resource by k8s cluster",
			},
			[]string{labelPodName, labelCluster, labelResourceType, labelConsumptionType, labelConsumptionStatus},
		),
	}
}

func (k *KrrAnalyzer) Run(ctx context.Context) {
	for {
		slog.Info("run krr analyzer")

		select {
		case <-ctx.Done():
			return
		default:
		}

		krrStats, err := k.callKRR()
		if err != nil {
			panic(err)
		}

		k.calculatePrice(krrStats)
	}
}

// calculatePrice iterates over each container’s old vs. new requests,
// asks Billing for their ruble cost, and accumulates totals.
// Returns (currentTotal, optimizedTotal, gain).
func (k *KrrAnalyzer) calculatePrice(rows []KrrOutput) {
	for _, r := range rows {
		sendMetrics := func(prev, new float64, resource Resource, unit ConsumptionMeasurementUnit) {
			k.writeConsumptionToGauge(
				r.Object.Pods[0].Name,
				r.Object.Cluster,
				resource, unit, ConsumptionStatusGain,
				prev-new,
			)

			k.writeConsumptionToGauge(
				r.Object.Pods[0].Name,
				r.Object.Cluster,
				resource, unit, ConsumptionStatusCurrent,
				prev,
			)

			k.writeConsumptionToGauge(
				r.Object.Pods[0].Name,
				r.Object.Cluster,
				resource, unit, ConsumptionStatusRecommended,
				new,
			)
		}

		if r.Object.Requests.CPU != nil && *r.Object.Requests.CPU > *r.Recommended.Requests.CPU {
			// todo: remove hardcoded platform
			price, err := k.billing.GetPriceCPURUB("standard-v1", "100", model.CPUCount(*r.Object.Requests.CPU-*r.Recommended.Requests.CPU))
			if err != nil {
				panic(err)
			}

			sendMetrics(
				*r.Object.Requests.CPU, *r.Recommended.Requests.CPU,
				ResourceCPU,
				ConsumptionReal,
			)

			sendMetrics(
				*r.Object.Requests.CPU*float64(price), *r.Recommended.Requests.CPU*float64(price), // todo: add money multiplier
				ResourceCPU,
				ConsumptionMoney,
			)
		}

		if r.Object.Requests.Memory != nil && *r.Object.Requests.Memory > *r.Recommended.Requests.Memory {
			// todo: remove hardcoded platform
			price, err := k.billing.GetPriceRAMRUB("standard-v1", model.RAMCount(*r.Object.Requests.Memory-*r.Recommended.Requests.Memory))
			if err != nil {
				panic(err)
			}

			sendMetrics(
				*r.Object.Requests.Memory, *r.Recommended.Requests.Memory,
				ResourceRAM,
				ConsumptionReal,
			)

			sendMetrics(
				*r.Object.Requests.Memory*float64(price), *r.Recommended.Requests.Memory*float64(price), // todo: add money multiplier
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
