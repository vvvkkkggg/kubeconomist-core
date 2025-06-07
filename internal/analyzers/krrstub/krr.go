package krrstub

import (
	"context"
	"log/slog"
	"math/rand/v2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/config"
)

var _ analyzers.Analyzer = &KrrAnalyzer{}

type KrrAnalyzer struct {
	billing analyzers.Billing
	cfg     config.KrrAnalyzerConfig

	resourceGauge *prometheus.GaugeVec
	krrStubOutput []KrrOutput
}

func NewKrrAnalyzer(
	b analyzers.Billing,
	cfg config.KrrAnalyzerConfig,
) *KrrAnalyzer {

	return &KrrAnalyzer{
		krrStubOutput: []KrrOutput{generateRandomKrrOutput()},
		billing:       b,
		cfg:           cfg,

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
	return k.krrStubOutput, nil
}

func generateRandomKrrOutput() KrrOutput {
	cluster := "cluster-" + randomString(5)
	namespace := "ns-" + randomString(3)
	container := "container-" + randomString(4)
	podCount := rand.Int()%3 + 1 // 1-3 пода

	pods := make([]Pod, podCount)
	for i := 0; i < podCount; i++ {
		pods[i] = Pod{Name: "pod-" + randomString(5)}
	}

	// Генерируем текущие значения (object)
	objectCPU := rand.Float64()*2 + 0.5       // 0.5-2.5 CPU
	objectMemory := float64(rand.Int()%8 + 2) // 2-10 GB

	// Генерируем рекомендованные значения (обычно на 20-50% ниже)
	recommendedCPU := objectCPU * (0.5 + rand.Float64()*0.3)       // 50-80% от текущего
	recommendedMemory := objectMemory * (0.5 + rand.Float64()*0.3) // 50-80% от текущего

	// Иногда (10% случаев) рекомендации могут быть такими же
	if rand.Float64() < 0.1 {
		recommendedCPU = objectCPU
		recommendedMemory = objectMemory
	}

	// Иногда (10% случаев) одно из значений может быть nil
	var cpuPtr, memPtr *float64
	var recCPUPtr, recMemPtr *float64

	if rand.Float64() > 0.1 {
		cpuPtr = &objectCPU
		recCPUPtr = &recommendedCPU
	}
	if rand.Float64() > 0.1 {
		memPtr = &objectMemory
		recMemPtr = &recommendedMemory
	}

	return KrrOutput{
		Object: Parameters{
			Cluster:   cluster,
			Namespace: namespace,
			Pods:      pods,
			Container: container,
			Requests: Resources{
				CPU:    cpuPtr,
				Memory: memPtr,
			},
			Limits: Resources{
				CPU:    cpuPtr,
				Memory: memPtr,
			},
		},
		Recommended: Parameters{
			Cluster:   cluster,
			Namespace: namespace,
			Pods:      pods,
			Container: container,
			Requests: Resources{
				CPU:    recCPUPtr,
				Memory: recMemPtr,
			},
			Limits: Resources{
				CPU:    recCPUPtr,
				Memory: recMemPtr,
			},
		},
	}
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Int()%len(charset)]
	}
	return string(b)
}
