package krr

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
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
	slog.Info("run krr analyzer")

	select {
	case <-ctx.Done():
		return
	default:
	}

	const fileName = "config/krr_config.json"
	// if err := k.runKRR(reportName); err != nil {
	// 	panic(err)
	// }

	krrStats, err := k.loadReport(fileName)
	if err != nil {
		panic(err)
	}

	k.calculatePrice(krrStats)
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

		if r.Object.Allocations.Requests.CPU != nil && *r.Object.Allocations.Requests.CPU > *r.Recommended.Requests.CPU.Value {
			price, err := k.billing.GetPriceCPURUB("standard-v1", "100", model.CPUCount(*r.Object.Allocations.Requests.CPU-*r.Recommended.Requests.CPU.Value))
			if err != nil {
				panic(err)
			}

			sendMetrics(
				*r.Object.Allocations.Requests.CPU, *r.Recommended.Requests.CPU.Value,
				ResourceCPU,
				ConsumptionReal,
			)

			sendMetrics(
				*r.Object.Allocations.Requests.CPU*float64(price)*float64(len(r.Object.Pods)),
				*r.Recommended.Requests.CPU.Value*float64(price)*float64(len(r.Object.Pods)),
				ResourceCPU,
				ConsumptionMoney,
			)
		}

		if r.Object.Allocations.Requests.Memory != nil && *r.Object.Allocations.Requests.Memory > *r.Recommended.Requests.Memory.Value {
			price, err := k.billing.GetPriceRAMRUB("standard-v1", model.RAMCount(*r.Object.Allocations.Requests.Memory-*r.Recommended.Requests.Memory.Value))
			if err != nil {
				panic(err)
			}

			sendMetrics(
				*r.Object.Allocations.Requests.Memory, *r.Recommended.Requests.Memory.Value,
				ResourceRAM,
				ConsumptionReal,
			)

			sendMetrics(
				*r.Object.Allocations.Requests.Memory*float64(price)*float64(len(r.Object.Pods)),
				*r.Recommended.Requests.Memory.Value*float64(price)*float64(len(r.Object.Pods)),
				ResourceRAM,
				ConsumptionMoney,
			)
		}
	}
}

func (k *KrrAnalyzer) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{k.resourceGauge}
}

func (k *KrrAnalyzer) runKRR(outputFile string) error {
	cmd := exec.Command(
		"krr", "simple",
		"-p", k.cfg.PrometheusURL,
		"--prometheus-auth-header", k.cfg.PrometheusAuthHeader,
		"--history-duration", k.cfg.HistoryDuration,
		"-f", "json",
		"--fileoutput", outputFile,
	)

	return cmd.Run()
}

func (k *KrrAnalyzer) loadReport(fileName string) ([]KrrOutput, error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	bytes = preprocessRawData(bytes)

	var krrResult Report
	if err := json.Unmarshal(bytes, &krrResult); err != nil {
		return nil, err
	}

	return krrResult.Scans, nil
}

func preprocessRawData(data []byte) []byte {
	return bytes.ReplaceAll(data, []byte(`"?"`), []byte(`null`))
}
