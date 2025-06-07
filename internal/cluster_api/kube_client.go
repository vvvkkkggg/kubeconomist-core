package clusterapi

import (
	"context"
	"fmt"
	"os"

	k8s "github.com/yandex-cloud/go-genproto/yandex/cloud/k8s/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetNodePlatform(ctx context.Context, namespace, deploymentName string) (string, error) {
	config, err := getKubeConfig()
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	node, err := getNode(ctx, clientset, namespace, deploymentName)
	if err != nil {
		return "", err
	}

	nodeGroup, err := ysdk.Kubernetes().NodeGroup().Get(ctx, &k8s.GetNodeGroupRequest{
		NodeGroupId: node.Labels["yandex.cloud/node-group-id"],
	})
	if err != nil {
		return "", err
	}

	return nodeGroup.NodeTemplate.PlatformId, nil
}

func getKubeConfig() (*rest.Config, error) {
	// Try in-cluster config first
	if config, err := rest.InClusterConfig(); err == nil {
		return config, nil
	}

	// Fall back to kubeconfig file
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func getNode(ctx context.Context, clientset *kubernetes.Clientset, namespace, deploymentName string) (*corev1.Node, error) {
	// 1. Get the deployment
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %v", err)
	}

	// 2. Get the pods for this deployment
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %v", err)
	}

	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no pods found for deployment %s", deploymentName)
	}

	// 3. Get node information for each pod
	for _, pod := range pods.Items {
		if pod.Spec.NodeName == "" {
			fmt.Printf("Pod %s is not assigned to a node yet\n", pod.Name)
			continue
		}

		// Get the node details
		node, err := clientset.CoreV1().Nodes().Get(ctx, pod.Spec.NodeName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Failed to get node %s: %v\n", pod.Spec.NodeName, err)
			continue
		}
		return node, nil

		// // Print node information
		// fmt.Printf("\nPod: %s\n", pod.Name)
		// fmt.Printf("Node Name: %s\n", node.Name)
		// fmt.Printf("Node ID: %s\n", node.Status.NodeInfo.MachineID)
		// fmt.Printf("Architecture: %s\n", node.Status.NodeInfo.Architecture)
		// fmt.Printf("OS: %s\n", node.Status.NodeInfo.OperatingSystem)
		// fmt.Printf("Kernel Version: %s\n", node.Status.NodeInfo.KernelVersion)
		// fmt.Printf("Container Runtime: %s\n", node.Status.NodeInfo.ContainerRuntimeVersion)
		// fmt.Printf("Kubelet Version: %s\n", node.Status.NodeInfo.KubeletVersion)
		// fmt.Printf("CPU Capacity: %s\n", node.Status.Capacity.Cpu())
		// fmt.Printf("Memory Capacity: %s\n", node.Status.Capacity.Memory())
		// fmt.Printf("Labels: %v\n", node.Labels)
		// fmt.Printf("Annotations: %v\n", node.Annotations)
	}
	return nil, fmt.Errorf("Failed to find any node")
}
