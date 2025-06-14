package registryoptimizer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/config"
	"github.com/vvvkkkggg/kubeconomist-core/internal/model"
	"github.com/vvvkkkggg/kubeconomist-core/internal/yandex"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ analyzers.Analyzer = &RegistryOptimizer{}

type Biling interface {
	GetContainerRegistryPriceRUB() (model.PriceRUB, error)
}

type RegistryOptimizer struct {
	billing   Biling
	clientset *kubernetes.Clientset
	yandex    *yandex.Client

	resourceGauge *prometheus.GaugeVec

	CloudID  string
	FolderID string
}

func NewRegistryOptimizer(billing Biling, cfg *config.Config) *RegistryOptimizer {
	ctx := context.Background()

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	yandex, err := yandex.New(ctx, cfg.Analyzers.VPC.YCToken)
	if err != nil {
		panic(err.Error())
	}

	return &RegistryOptimizer{
		billing:   billing,
		clientset: clientset,
		yandex:    yandex,

		resourceGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "kubeconomist",
				Subsystem: "node_optimizer",
				Name:      "unused_images_storage_cost_per_hour",
				Help:      "Current hourly storage cost for registry images in RUB",
			},
			[]string{},
		),
	}
}

func getYandexImages(ctx context.Context, yandex *yandex.Client, cloudID, folderID string) (map[string]*compute.Image, error) {
	folders, err := yandex.GetAllFolders(ctx, cloudID, folderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get folders: %w", err)
	}

	ycImages := make(map[string]*compute.Image, 0)

	for _, folder := range folders {
		images, err := yandex.GetImages(ctx, folder.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get addresses for folder %s: %w", folder.Id, err)
		}

		for _, image := range images.Images {
			ycImages[image.Name] = image
		}
	}

	return ycImages, nil
}

func getK8SImages(ctx context.Context, clientset *kubernetes.Clientset) ([]string, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %w", err)
	}

	imageSet := make(map[string]struct{})

	for _, ns := range namespaces.Items {
		pods, err := clientset.CoreV1().Pods(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error getting pods in namespace %s: %v\n", ns.Name, err)

			continue
		}

		for _, pod := range pods.Items {
			for _, container := range pod.Spec.Containers {
				imageSet[container.Image] = struct{}{}
			}

			for _, container := range pod.Spec.InitContainers {
				imageSet[container.Image] = struct{}{}
			}

			for _, container := range pod.Spec.EphemeralContainers {
				imageSet[container.Image] = struct{}{}
			}
		}
	}

	images := make([]string, 0, len(imageSet))
	for image := range imageSet {
		images = append(images, image)
	}

	return images, nil
}

func computeCost(ycImages map[string]*compute.Image, k8sImages []string, registryCost float64) float64 {
	inUse := func(ycName string) bool {
		for _, imageName := range k8sImages {
			// MAY BE NOT GOOD BUT IDK
			if strings.Contains(imageName, ycName) {
				return true
			}
		}

		return false
	}

	imagesToDelete := make([]string, 0, len(ycImages))

	// если есть в Yandex Cloud, но нет в кубе
	for name := range ycImages {
		if !inUse(name) {
			imagesToDelete = append(imagesToDelete, name)
		}
	}

	sumGB := 0.0

	for _, name := range imagesToDelete {
		if image, exists := ycImages[name]; exists {
			sizeGB := float64(image.StorageSize) / (1024 * 1024 * 1024)
			sumGB += sizeGB
		}
	}

	// сколько сэкономим в час
	return sumGB * registryCost
}

func (ro *RegistryOptimizer) Run(ctx context.Context) {
	ycImages, err := getYandexImages(ctx, ro.yandex, ro.CloudID, ro.FolderID)
	if err != nil {
		panic(err.Error())
	}

	k8sImages, err := getK8SImages(ctx, ro.clientset)
	if err != nil {
		panic(err.Error())
	}

	registryCost, err := ro.billing.GetContainerRegistryPriceRUB()
	if err != nil {
		panic(err.Error())
	}

	totalCost := computeCost(ycImages, k8sImages, float64(registryCost))

	ro.resourceGauge.WithLabelValues().Set(totalCost)
}

func (ro *RegistryOptimizer) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{ro.resourceGauge}
}
