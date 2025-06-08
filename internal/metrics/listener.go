package metrics

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vvvkkkggg/kubeconomist-core/internal/config"
)

type OptimizerRecommendations struct {
	VPCRecommendations      []VPCOptimizerRecommendations      `json:"vpcRecommendations"`
	DNSRecommendations      []DNSOptimizerRecommendations      `json:"dnsRecommendations"`
	NodeRecommendations     []NodeOptimizerRecommendations     `json:"nodeRecommendations"`
	PlatformRecommendations []PlatformOptimizerRecommendations `json:"platformRecommendations"`
}

type VPCOptimizerRecommendations struct {
	Id         string `json:"id"`
	CloudID    string `json:"cloudId"`
	FolderID   string `json:"folderId"`
	IPAddress  string `json:"ipAddress"`
	IsUsed     bool   `json:"isUsed"`
	IsReserved bool   `json:"isReserved"`
}

type DNSOptimizerRecommendations struct {
	Id       string `json:"id"`
	CloudID  string `json:"cloudId"`
	FolderID string `json:"folderId"`
	ZoneId   string `json:"zoneId"`
	IsUsed   bool   `json:"isUsed"`
}

type NodeOptimizerRecommendations struct {
	Id            string  `json:"id"`
	CloudID       string  `json:"cloudId"`
	FolderID      string  `json:"folderId"`
	InstanceId    string  `json:"instanceId"`
	CurrentCores  int     `json:"currentCores"`
	DesiredCores  int     `json:"desiredCores"`
	CurrentMemory int     `json:"currentMemory"`
	DesiredMemory int     `json:"desiredMemory"`
	CurrentPrice  float64 `json:"currentPrice"`
	DesiredPrice  float64 `json:"desiredPrice"`
}

type PlatformOptimizerRecommendations struct {
	Id              string  `json:"id"`
	CloudID         string  `json:"cloudId"`
	FolderID        string  `json:"folderId"`
	NodeGroupId     string  `json:"nodeGroupId"`
	CurrentPlatform string  `json:"currentPlatform"`
	DesiredPlatform string  `json:"desiredPlatform"`
	CurrentPrice    float64 `json:"currentPrice"`
	DesiredPrice    float64 `json:"desiredPrice"`
	Savings         float64 `json:"savings"`
}

func GetOptimizerRecommendations() (OptimizerRecommendations, error) {
	metricFamilies, err := registry.Gather()
	if err != nil {
		return OptimizerRecommendations{}, err
	}

	vpcOpt := make([]VPCOptimizerRecommendations, 0)
	nodeOpt := make(map[string]NodeOptimizerRecommendations)
	dnsOpt := make([]DNSOptimizerRecommendations, 0)
	platformOpt := make(map[string]PlatformOptimizerRecommendations)

	for _, mf := range metricFamilies {
		metricsName := mf.GetName()
		for _, m := range mf.GetMetric() {
			switch metricsName {
			case "kubeconomist_vpc_ip_info":
				v := VPCOptimizerRecommendations{}

				for _, label := range m.GetLabel() {
					switch label.GetName() {
					case "cloud_id":
						v.CloudID = label.GetValue()
					case "folder_id":
						v.FolderID = label.GetValue()
					case "ip_address":
						v.IPAddress = label.GetValue()
					case "is_used":
						v.IsUsed = label.GetValue() == "true"
					case "is_reserved":
						v.IsReserved = label.GetValue() == "true"
					}
				}

				vpcOpt = append(vpcOpt, v)
			case "kubeconomist_dns_optimizer_dns_optimization_zone":
				v := DNSOptimizerRecommendations{}

				for _, label := range m.GetLabel() {
					switch label.GetName() {
					case "cloud_id":
						v.CloudID = label.GetValue()
					case "folder_id":
						v.FolderID = label.GetValue()
					case "zone_id":
						v.ZoneId = label.GetValue()
					case "is_used":
						v.IsUsed = label.GetValue() == "true"
					}
				}
				dnsOpt = append(dnsOpt, v)

			case "kubeconomist_node_optimizer_node_optimization_cores":
				cloudID, folderID, instanceID, status := "", "", "", ""

				for _, label := range m.GetLabel() {
					switch label.GetName() {
					case "cloud_id":
						cloudID = label.GetValue()
					case "folder_id":
						folderID = label.GetValue()
					case "instance_id":
						instanceID = label.GetValue()
					case "status":
						status = label.GetValue()
					}
				}

				v, ok := nodeOpt[instanceID]
				if !ok {
					v = NodeOptimizerRecommendations{}
				}

				v.CloudID = cloudID
				v.FolderID = folderID
				v.InstanceId = instanceID

				if status == "current" {
					v.CurrentCores = int(m.GetGauge().GetValue())
				} else {
					v.DesiredCores = int(m.GetGauge().GetValue())
				}

				nodeOpt[instanceID] = v

			case "kubeconomist_node_optimizer_node_optimization_memory":
				cloudID, folderID, instanceID, status := "", "", "", ""

				for _, label := range m.GetLabel() {
					switch label.GetName() {
					case "cloud_id":
						cloudID = label.GetValue()
					case "folder_id":
						folderID = label.GetValue()
					case "instance_id":
						instanceID = label.GetValue()
					case "status":
						status = label.GetValue()
					}
				}

				v, ok := nodeOpt[instanceID]
				if !ok {
					v = NodeOptimizerRecommendations{}
				}

				v.CloudID = cloudID
				v.FolderID = folderID
				v.InstanceId = instanceID

				if status == "current" {
					v.CurrentMemory = int(m.GetGauge().GetValue())
				} else {
					v.DesiredMemory = int(m.GetGauge().GetValue())
				}

				nodeOpt[instanceID] = v

			case "kubeconomist_node_optimizer_node_optimization_price":
				cloudID, folderID, instanceID, status := "", "", "", ""

				for _, label := range m.GetLabel() {
					switch label.GetName() {
					case "cloud_id":
						cloudID = label.GetValue()
					case "folder_id":
						folderID = label.GetValue()
					case "instance_id":
						instanceID = label.GetValue()
					case "status":
						status = label.GetValue()
					}
				}

				v, ok := nodeOpt[instanceID]
				if !ok {
					v = NodeOptimizerRecommendations{}
				}

				v.CloudID = cloudID
				v.FolderID = folderID
				v.InstanceId = instanceID

				if status == "current" {
					v.CurrentPrice = m.GetGauge().GetValue()
				} else {
					v.DesiredPrice = m.GetGauge().GetValue()
				}

				nodeOpt[instanceID] = v

			case "kubeconomist_platform_optimizer_platform_optimizer_price":
				cloudID, folderID, nodeGroupId, platformID, status := "", "", "", "", ""

				for _, label := range m.GetLabel() {
					switch label.GetName() {
					case "cloud_id":
						cloudID = label.GetValue()
					case "folder_id":
						folderID = label.GetValue()
					case "node_group_id":
						nodeGroupId = label.GetValue()
					case "platform_id":
						platformID = label.GetValue()
					case "status":
						status = label.GetValue()
					}
				}

				v, ok := platformOpt[nodeGroupId]
				if !ok {
					v = PlatformOptimizerRecommendations{}
				}

				v.CloudID = cloudID
				v.FolderID = folderID
				v.NodeGroupId = nodeGroupId

				if status == "current" {
					v.CurrentPlatform = platformID
					v.CurrentPrice = m.GetGauge().GetValue()
				} else {
					v.DesiredPlatform = platformID
					v.DesiredPrice = m.GetGauge().GetValue()
				}

				platformOpt[nodeGroupId] = v
				// case "kubeconomist_storage_optimizer_storage_optimization_class_is_optimal":
				// case "kubeconomist_storage_optimizer_storage_optimization_price":
			}
		}
	}

	v := OptimizerRecommendations{
		VPCRecommendations:      vpcOpt,
		DNSRecommendations:      dnsOpt,
		NodeRecommendations:     make([]NodeOptimizerRecommendations, 0, len(nodeOpt)),
		PlatformRecommendations: make([]PlatformOptimizerRecommendations, 0, len(platformOpt)),
	}
	for _, n := range nodeOpt {
		v.NodeRecommendations = append(v.NodeRecommendations, n)
	}
	for _, p := range platformOpt {
		v.PlatformRecommendations = append(v.PlatformRecommendations, p)
	}

	return v, nil
}

var registry = prometheus.NewRegistry()

func ListenAndServe(ctx context.Context, cfg config.MetricsConfig, metricsCollector ...prometheus.Collector) error {
	var (
		chErr    = make(chan error)
		mux      = http.NewServeMux()
		registry = registry
	)
	registry.MustRegister(metricsCollector...)

	mux.Handle("/metrics", promhttp.InstrumentMetricHandler(
		registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{})),
	)
	mux.HandleFunc("/metrics-json", func(w http.ResponseWriter, r *http.Request) {
		v, err := GetOptimizerRecommendations()
		if err != nil {
			http.Error(w, "Failed to gather metrics", http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(v); err != nil {
			http.Error(w, "Failed to encode metrics", http.StatusInternalServerError)
		}
	})

	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: mux,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			chErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		ctxShutdown, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
		defer cancelFunc()
		return server.Shutdown(ctxShutdown) //nolint:contextcheck // deliberately created context for shutdown
	case err := <-chErr:
		return err
	}
}
