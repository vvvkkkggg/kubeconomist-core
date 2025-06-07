package registryoptimizer

import (
	"testing"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
)

func TestComputeCost_AllImagesInUse(t *testing.T) {
	ycImages := map[string]*compute.Image{
		"nginx": {StorageSize: 1073741824}, // 1GB
		"redis": {StorageSize: 5368709120}, // 5GB
	}

	k8sImages := []string{"docker.io/nginx:latest", "docker.io/redis:alpine"}
	registryCost := 0.5 // $0.5/GB-hour

	cost := computeCost(ycImages, k8sImages, registryCost)
	if cost != 0.0 {
		t.Errorf("Expected cost 0, got %.2f", cost)
	}
}

func TestComputeCost_SomeUnusedImages(t *testing.T) {
	ycImages := map[string]*compute.Image{
		"nginx": {StorageSize: 1073741824}, // 1GB
		"redis": {StorageSize: 2147483648}, // 2GB (unused)
		"mysql": {StorageSize: 3221225472}, // 3GB (unused)
	}

	k8sImages := []string{"docker.io/nginx:latest"}
	registryCost := 0.2 // $0.2/GB-hour

	expected := (2 + 3) * registryCost
	cost := computeCost(ycImages, k8sImages, registryCost)
	if cost != expected {
		t.Errorf("Expected %.2f, got %.2f", expected, cost)
	}
}

func TestComputeCost_NameMatchingVariations(t *testing.T) {
	ycImages := map[string]*compute.Image{
		"app":      {StorageSize: 1073741824}, // 1GB
		"app-beta": {StorageSize: 1073741824}, // 1GB
	}

	tests := []struct {
		name     string
		k8sImage string
		expected float64
	}{
		{"Exact match", "app", 0.3},
		{"Partial match", "my-app:v1", 0.3},
		{"Different registry", "registry.example.com/app", 0.3},
		{"No match", "database", 2 * 0.3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := computeCost(ycImages, []string{tt.k8sImage}, 0.3)
			if cost != tt.expected {
				t.Errorf("Expected %.2f, got %.2f", tt.expected, cost)
			}
		})
	}
}

func TestComputeCost_EmptyInputs(t *testing.T) {
	tests := []struct {
		name      string
		ycImages  map[string]*compute.Image
		k8sImages []string
	}{
		{"No YC images", map[string]*compute.Image{}, []string{"nginx"}},
		{"No k8s images", map[string]*compute.Image{"nginx": {}}, nil},
		{"Both empty", map[string]*compute.Image{}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := computeCost(tt.ycImages, tt.k8sImages, 0.5)
			if cost != 0.0 {
				t.Errorf("Expected 0, got %.2f", cost)
			}
		})
	}
}

func TestComputeCost_StorageSizeEdgeCases(t *testing.T) {
	ycImages := map[string]*compute.Image{
		"tiny":  {StorageSize: 1},             // 1 byte
		"large": {StorageSize: 1099511627776}, // 1TB (1024GB)
	}
	k8sImages := []string{} // all unused
	registryCost := 0.01    // $0.01/GB-hour

	expected := (1.0/1073741824 + 1024) * registryCost
	cost := computeCost(ycImages, k8sImages, registryCost)

	// Allow small floating-point error
	if diff := cost - expected; diff > 0.000001 || diff < -0.000001 {
		t.Errorf("Expected ~%.10f, got %.10f", expected, cost)
	}
}

func TestComputeCost_RegistryCostVariations(t *testing.T) {
	ycImages := map[string]*compute.Image{
		"app": {StorageSize: 1073741824}, // 1GB
	}
	k8sImages := []string{} // unused

	tests := []struct {
		cost     float64
		expected float64
	}{
		{0.0, 0.0},
		{0.5, 0.5},
		{1.23, 1.23},
		{100, 100},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := computeCost(ycImages, k8sImages, tt.cost)
			if result != tt.expected {
				t.Errorf("Expected %.2f, got %.2f", tt.expected, result)
			}
		})
	}
}
