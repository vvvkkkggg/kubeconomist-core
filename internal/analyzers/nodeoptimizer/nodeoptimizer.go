package nodeoptimizer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/billing"
	"github.com/vvvkkkggg/kubeconomist-core/internal/yandex"
)

// Hardcode here, because yandex does not provide this information in API
var InstanceConfigurations = map[string]map[int][]struct {
	Cores  []int
	Memory []float64
}{
	"standard-v1": {
		5: {
			{Cores: []int{2, 4}, Memory: []float64{0.5, 1, 1.5, 2}},
		},
		20: {
			{Cores: []int{2, 4}, Memory: []float64{0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4}},
		},
		100: {
			{Cores: []int{2, 4, 6, 8, 10, 12, 14, 16, 20, 24, 28, 32}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8}},
		},
	},
	"standard-v2": {
		5: {
			{Cores: []int{2, 4}, Memory: []float64{0.25, 0.5, 1, 1.5, 2}},
		},
		20: {
			{Cores: []int{2, 4}, Memory: []float64{0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4}},
		},
		50: {
			{Cores: []int{2, 4}, Memory: []float64{0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4}},
		},
		100: {
			{Cores: []int{2, 4, 6, 8, 10, 12, 14, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 68, 72, 76, 80}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}},
		},
	},
	"standard-v3": {
		20: {
			{Cores: []int{2, 4}, Memory: []float64{0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4}},
		},
		50: {
			{Cores: []int{2, 4}, Memory: []float64{0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4}},
		},
		100: {
			{Cores: []int{2, 4, 6, 8, 10, 12, 14, 16, 20}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}},
			{Cores: []int{24}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}},
			{Cores: []int{28}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}},
			{Cores: []int{32}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{Cores: []int{36, 40}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}},
			{Cores: []int{44}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}},
			{Cores: []int{48}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}},
			{Cores: []int{52}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
			{Cores: []int{56}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}},
			{Cores: []int{60, 64}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{Cores: []int{68}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}},
			{Cores: []int{72, 76, 80}, Memory: []float64{1, 2, 3, 4, 5, 6, 7, 8}},
			{Cores: []int{84, 88}, Memory: []float64{1, 2, 3, 4, 5, 6, 7}},
			{Cores: []int{92, 96}, Memory: []float64{1, 2, 3, 4, 5, 6}},
		},
	},
}

type NodeOptimizer struct {
	yandex  *yandex.Client
	billing *billing.Billing

	nodeCostMetric   *prometheus.GaugeVec
	nodeCoresMetric  *prometheus.GaugeVec
	nodeMemoryMetric *prometheus.GaugeVec
}

func NewNodeOptimizer(yandex *yandex.Client, billing *billing.Billing) *NodeOptimizer {
	nodeCostMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kubeconomist",
			Subsystem: "node_optimizer",
			Name:      "node_cost",
			Help:      "Cost of the node",
		},
		[]string{"cloud_id", "folder_id", "instance_id", "status"},
	)

	nodeCoresMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kubeconomist",
			Subsystem: "node_optimizer",
			Name:      "node_cores",
			Help:      "Number of cores in the node",
		},
		[]string{"cloud_id", "folder_id", "instance_id", "status"},
	)

	nodeMemoryMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kubeconomist",
			Subsystem: "node_optimizer",
			Name:      "node_memory",
			Help:      "Memory of the node",
		},
		[]string{"cloud_id", "folder_id", "instance_id", "status"},
	)

	return &NodeOptimizer{
		yandex:           yandex,
		billing:          billing,
		nodeCostMetric:   nodeCostMetric,
		nodeCoresMetric:  nodeCoresMetric,
		nodeMemoryMetric: nodeMemoryMetric,
	}
}

func QueryPrometheus(promURL, query string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/query?query=%s", promURL, query))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Data struct {
			Result []struct {
				Value [2]interface{}
			}
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	if len(data.Data.Result) == 0 {
		return 0, fmt.Errorf("no data")
	}

	valueStr := data.Data.Result[0].Value[1].(string)
	var value float64
	fmt.Sscanf(valueStr, "%f", &value)
	return value, nil
}

func FindClosestCPU(platformID string, coreFraction int, targetCPU float64) int {
	configs, ok := InstanceConfigurations[platformID][coreFraction]
	if !ok {
		return -1
	}

	min := -1
	for i := len(configs) - 1; i >= 0; i-- {
		for j := len(configs[i].Cores) - 1; j >= 0; j-- {
			if float64(configs[i].Cores[j]) > targetCPU {
				min = configs[i].Cores[j]
			}
		}
	}

	return min
}

func (n *NodeOptimizer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			slog.Error("ctx done", slog.Any("err", ctx.Err()))
			return
		default:
		}

		clouds, err := n.yandex.GetClouds(ctx)
		if err != nil {
			slog.Error("get clouds err", slog.Any("err", err))
			return
		}

		for _, cloud := range clouds {
			folders, err := n.yandex.GetFolders(ctx, cloud.Id)
			if err != nil {
				slog.Error("get folders err", slog.Any("err", err))
				return
			}

			for _, folder := range folders {
				instances, err := n.yandex.GetInstances(ctx, folder.Id)
				if err != nil {
					slog.Error("get instances err", slog.Any("err", err))
					return
				}

				for _, instance := range instances {
					coreFraction := instance.GetResources().GetCoreFraction()
					cores := instance.GetResources().GetCores()
					memory := instance.GetResources().GetMemory()
					platformID := instance.GetPlatformId()

					// TODO: yandex monitoring does not export RAM usage
					currentCPU, err := QueryPrometheus("http://localhost:8428", fmt.Sprintf("rate(cpu_usage{resource_id=\"%s\"}[1h])", instance.GetName()))
					if err != nil {
						fmt.Println(instance.GetName())
						slog.Error("query prometheus err", slog.Any("err", err))
						continue
					}

					minCPU := FindClosestCPU(platformID, int(coreFraction), currentCPU)
					if minCPU == -1 {
						slog.Error("find closest CPU err", slog.String("platform_id", instance.GetPlatformId()), slog.Int("core_fraction", int(instance.GetResources().GetCoreFraction())))
						continue
					}

					currentPrice, err := n.billing.CalculatePrice(platformID, coreFraction, cores, memory)
					if err != nil {
						slog.Error("calculate price err", slog.Any("err", err))
						continue
					}

					desiredPrice, err := n.billing.CalculatePrice(platformID, coreFraction, int64(minCPU), memory)
					if err != nil {
						slog.Error("calculate desired price err", slog.Any("err", err))
						continue
					}

					n.nodeCostMetric.WithLabelValues(
						cloud.Id,
						folder.Id,
						instance.GetName(),
						"current",
					).Set(currentPrice)

					n.nodeCostMetric.WithLabelValues(
						cloud.Id,
						folder.Id,
						instance.GetName(),
						"desired",
					).Set(desiredPrice)

					n.nodeCoresMetric.WithLabelValues(
						cloud.Id,
						folder.Id,
						instance.GetName(),
						"current",
					).Set(float64(cores))

					n.nodeCoresMetric.WithLabelValues(
						cloud.Id,
						folder.Id,
						instance.GetName(),
						"desired",
					).Set(float64(minCPU))

					// Convert memory from bytes to gigabytes
					memoryGB := float64(memory) / (1024 * 1024 * 1024)

					n.nodeMemoryMetric.WithLabelValues(
						cloud.Id,
						folder.Id,
						instance.GetName(),
						"current",
					).Set(memoryGB)

					n.nodeMemoryMetric.WithLabelValues(
						cloud.Id,
						folder.Id,
						instance.GetName(),
						"desired",
					).Set(memoryGB)
				}
			}
		}
	}
}

func (n *NodeOptimizer) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		n.nodeCostMetric,
		n.nodeCoresMetric,
		n.nodeMemoryMetric,
	}
}
