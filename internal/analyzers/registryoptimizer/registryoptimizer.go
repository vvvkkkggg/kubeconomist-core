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

type RegistryOptimizer struct {
	Biling interface {
		GetRegistryCost() (float64, error)
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
		panic(err.Error())
	}

	for _, ns := range namespaces.Items {
		fmt.Printf("Namespace: %s\n", ns.Name)

		pods, err := clientset.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error getting pods: %v\n", err)
			continue
		}

		for _, pod := range pods.Items {
			fmt.Printf("  - Pod: %s\n", pod.Name)
			fmt.Printf("    Status: %s\n", pod.Status.Phase)
			fmt.Printf("    Node: %s\n", pod.Spec.NodeName)
			fmt.Printf("    IP: %s\n", pod.Status.PodIP)
			fmt.Println("    Containers:")
			for _, container := range pod.Spec.Containers {
				fmt.Printf("      - %s (Image: %s)\n", container.Name, container.Image)
			}
			fmt.Println("    --------------------")
		}
		fmt.Println("==============================")
	}

	return nil, nil
}

func (ro *RegistryOptimizer) Run(ctx context.Context) {
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
		log.Fatalf("Failed to create Yandex client: %v", err)
	}

	ycImages, err := getYandexImages(ctx, yandex)
	if err != nil {

	}

	k8sImages, err := getK8SImages(clientset)
	if err != nil {

	}

	registryCost, err := ro.Biling.GetRegistryCost()
	if err != nil {

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

	_ = totalCost

	return

}

func (ro *RegistryOptimizer) GetCollectors() *prometheus.Collector {
	return nil
}
