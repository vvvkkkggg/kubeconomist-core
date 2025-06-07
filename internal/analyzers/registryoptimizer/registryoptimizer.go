package registryoptimizer

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/yandex"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ analyzers.Analyzer = &RegistryOptimizer{}

type Biling interface {
	GetRegistryCost() (float64, error)
}

type RegistryOptimizer struct {
	billing   Biling
	clientset *kubernetes.Clientset
	yandex    *yandex.Client

	resourceGauge *prometheus.GaugeVec
}

func NewRegistryOptimizer(billing Biling) *RegistryOptimizer {
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

	yandex, err := yandex.New(ctx, os.Getenv("YANDEX_TOKEN"))
	if err != nil {
		panic(err.Error())
	}

	return &RegistryOptimizer{
		billing:   billing,
		clientset: clientset,
		yandex:    yandex,

		resourceGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "ycr_registry_storage_cost_per_hour",
				Help: "Current hourly storage cost for registry images in RUB",
			},
			[]string{},
		),
	}
}

func getYandexImages(ctx context.Context, yandex *yandex.Client) (map[string]*compute.Image, error) {
	clouds, err := yandex.GetClouds(ctx)
	if err != nil {
		log.Fatalf("Failed to get clouds: %v", err)
	}

	ycImages := make(map[string]*compute.Image, 0)
	for _, cloud := range clouds {
		folders, err := yandex.GetFolders(ctx, cloud.Id)
		if err != nil {
			log.Fatalf("Failed to get folders for cloud %s: %v", cloud.Id, err)
		}

		for _, folder := range folders {
			images, err := yandex.GetImages(ctx, folder.Id)
			if err != nil {
				log.Fatalf("Failed to get addresses for folder %s: %v", folder.Id, err)
			}

			for _, image := range images.Images {
				ycImages[image.Name] = image
			}
		}
	}

	return ycImages, err
}

func getK8SImages(clientset *kubernetes.Clientset) ([]string, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %w", err)
	}

	imageSet := make(map[string]struct{})

	for _, ns := range namespaces.Items {
		pods, err := clientset.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
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

func (ro *RegistryOptimizer) Run(ctx context.Context) {
	ycImages, err := getYandexImages(ctx, ro.yandex)
	if err != nil {
		panic(err.Error())
	}

	k8sImages, err := getK8SImages(ro.clientset)
	if err != nil {
		panic(err.Error())
	}

	registryCost, err := ro.billing.GetRegistryCost()
	if err != nil {
		panic(err.Error())
	}

	inUse := func(ycName string) bool {
		for _, imageName := range k8sImages {
			if strings.Contains(imageName, ycName) {
				return true
			}
		}

		return false
	}

	imagesToDelete := make([]string, 0, len(ycImages))

	// если есть в Yandex Cloud, но нет в кубе
	for name := range ycImages {
		if inUse(name) {
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
	totalCost := sumGB * registryCost

	ro.resourceGauge.WithLabelValues().Set(totalCost)
}

func (ro *RegistryOptimizer) GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{ro.resourceGauge}
}
