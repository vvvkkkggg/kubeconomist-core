package clusterapi

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetNodeDescription(namespace, deploymentName string) error {
	config, err := getKubeConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	err = getNodeInfo(clientset, namespace, deploymentName)
	if err != nil {
		return err
	}

	return nil
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

func getNodeInfo(clientset *kubernetes.Clientset, namespace, deploymentName string) error {
	ctx := context.Background()

	// 1. Get the deployment
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %v", err)
	}

	// 2. Get the pods for this deployment
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
	})
	if err != nil {
		return fmt.Errorf("failed to get pods: %v", err)
	}

	if len(pods.Items) == 0 {
		return fmt.Errorf("no pods found for deployment %s", deploymentName)
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

		// Print node information
		fmt.Printf("\nPod: %s\n", pod.Name)
		fmt.Printf("Node Name: %s\n", node.Name)
		fmt.Printf("Node ID: %s\n", node.Status.NodeInfo.MachineID)
		fmt.Printf("Architecture: %s\n", node.Status.NodeInfo.Architecture)
		fmt.Printf("OS: %s\n", node.Status.NodeInfo.OperatingSystem)
		fmt.Printf("Kernel Version: %s\n", node.Status.NodeInfo.KernelVersion)
		fmt.Printf("Container Runtime: %s\n", node.Status.NodeInfo.ContainerRuntimeVersion)
		fmt.Printf("Kubelet Version: %s\n", node.Status.NodeInfo.KubeletVersion)
		fmt.Printf("CPU Capacity: %s\n", node.Status.Capacity.Cpu())
		fmt.Printf("Memory Capacity: %s\n", node.Status.Capacity.Memory())
		fmt.Printf("Labels: %v\n", node.Labels)
		fmt.Printf("Annotations: %v\n", node.Annotations)
	}

	return nil
}
