package registryoptimizer

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/yandex"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ analyzers.Analyzer = &RegistryOptimizer{}

type RegistryOptimizer struct {
}

func getYandexImages(ctx context.Context, yandex *yandex.Client) (map[string]struct{}, error) {
	clouds, err := yandex.GetClouds(ctx)
	if err != nil {
		log.Fatalf("Failed to get clouds: %v", err)
	}

	ycImages := make(map[string]struct{}, 0)
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
				ycImages[image.Name] = struct{}{}
			}
		}
	}

	return ycImages, err
}

func getK8SImages(clientset *kubernetes.Clientset) (map[string]struct{}, error) {
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

	_ = ycImages
	_ = k8sImages
}

func (ro *RegistryOptimizer) GetCollectors() *prometheus.Collector {
	return nil
}
