package alerter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/vvvkkkggg/kubeconomist-core/internal/metrics"
)

func TestGenerateReport_EmptyRecommendations(t *testing.T) {
	a := NewAlerter(nil)
	report := a.generateReport(metrics.OptimizerRecommendations{})

	if !strings.Contains(report, "🎉 *Поздравляем! Оптимизация не требуется") {
		t.Error("Отчет не содержит сообщение об отсутствии оптимизации")
	}
}

func TestGenerateReport_FullRecommendations(t *testing.T) {
	a := NewAlerter(nil)

	report := metrics.OptimizerRecommendations{
		VPCRecommendations: []metrics.VPCOptimizerRecommendations{
			{
				IPAddress: "192.168.1.10",
				CloudID:   "aws",
				FolderID:  "prod",
				IsUsed:    false,
			},
		},
		DNSRecommendations: []metrics.DNSOptimizerRecommendations{
			{
				ZoneId:   "zone123",
				CloudID:  "gcp",
				FolderID: "dev",
				IsUsed:   false,
			},
		},
		NodeRecommendations: []metrics.NodeOptimizerRecommendations{
			{
				InstanceId:    "i-123456",
				CurrentCores:  8,
				DesiredCores:  4,
				CurrentMemory: 16384,
				DesiredMemory: 8192,
				CurrentPrice:  5000.50,
				DesiredPrice:  3750.00,
			},
		},
		PlatformRecommendations: []metrics.PlatformOptimizerRecommendations{
			{
				NodeGroupId:     "ng-prod",
				CurrentPlatform: "n1-standard",
				DesiredPlatform: "e2-medium",
				Savings:         1250.75,
			},
		},
	}

	result := a.generateReport(report)

	// Проверка секций
	tests := []struct {
		name     string
		expected string
	}{
		{"VPC section", "🔌 *Неиспользуемые IP-адреса:*"},
		{"DNS section", "🌐 *Неиспользуемые DNS зоны:*"},
		{"Node section", "🖥️ *Оптимизация нод:*"},
		{"Platform section", "📦 *Оптимизация платформ:*"},
		{"Total savings", "💸 *Итоговая экономия:*"},
		{"VPC entry", "- `192.168.1.10` (Cloud: `aws`, Folder: `prod`) 🔴 Не используется"},
		{"DNS entry", "- `zone123` (Cloud: `gcp`, Folder: `dev`) 🔴 Не используется"},
		{"Node entry", "  - CPU: 8 → 4 ядер"},
		{"Platform entry", "  - Платформа: `n1-standard` → `e2-medium`"},
		{"Total calculation", "✨ *Всего: 2501.25 руб.*"},
	}

	for _, tt := range tests {
		if !strings.Contains(result, tt.expected) {
			t.Errorf("Отчет не содержит ожидаемый элемент: %s", tt.name)
		}
	}
}

func TestGenerateReport_OnlyUnusedResources(t *testing.T) {
	a := NewAlerter(nil)
	report := metrics.OptimizerRecommendations{
		VPCRecommendations: []metrics.VPCOptimizerRecommendations{
			{IPAddress: "10.0.0.5", IsUsed: false},
		},
		DNSRecommendations: []metrics.DNSOptimizerRecommendations{
			{ZoneId: "internal-zone", IsUsed: false},
		},
	}

	result := a.generateReport(report)

	// Проверяем что нет секций с экономией
	if strings.Contains(result, "💸 *Итоговая экономия:*") {
		t.Error("Отчет не должен содержать секцию экономии при отсутствии рекомендаций по нодам/платформам")
	}

	// Проверяем наличие ресурсов
	if !strings.Contains(result, "10.0.0.5") || !strings.Contains(result, "internal-zone") {
		t.Error("Отчет не содержит все неиспользуемые ресурсы")
	}
}

func TestGenerateReport_NodeOptimizationsOnly(t *testing.T) {
	a := NewAlerter(nil)
	report := metrics.OptimizerRecommendations{
		NodeRecommendations: []metrics.NodeOptimizerRecommendations{
			{
				InstanceId:    "node-1",
				CurrentCores:  16,
				DesiredCores:  8,
				CurrentMemory: 32768,
				DesiredMemory: 16384,
				CurrentPrice:  10000,
				DesiredPrice:  6000,
			},
			{
				InstanceId:    "node-2",
				CurrentCores:  4,
				DesiredCores:  2,
				CurrentMemory: 8192,
				DesiredMemory: 4096,
				CurrentPrice:  2500,
				DesiredPrice:  1500,
			},
		},
	}

	result := a.generateReport(report)

	// Проверка экономии (10000-6000 + 2500-1500 = 5000)
	if !strings.Contains(result, "✨ *Всего: 5000.00 руб.*") {
		t.Errorf("Неправильный расчет экономии для нод: %s", result)
	}

	if !strings.Contains(result, "32768 → 16384 MB") {
		t.Error("Неправильное форматирование памяти")
	}
}

func TestGenerateReport_MixedSavings(t *testing.T) {
	a := NewAlerter(nil)
	report := metrics.OptimizerRecommendations{
		NodeRecommendations: []metrics.NodeOptimizerRecommendations{
			{CurrentPrice: 5000, DesiredPrice: 4000},
		},
		PlatformRecommendations: []metrics.PlatformOptimizerRecommendations{
			{Savings: 750.25},
			{Savings: 1250.50},
		},
	}

	result := a.generateReport(report)

	// Проверка суммарной экономии (1000 + 750.25 + 1250.50 = 3000.75)
	if !strings.Contains(result, "✨ *Всего: 3000.75 руб.*") {
		t.Errorf("Неправильный расчет суммарной экономии: %s", result)
	}

	// Проверка разбивки по категориям
	if !strings.Contains(result, "- Ноды: *1000.00 руб.*") ||
		!strings.Contains(result, "- Платформы: *2000.75 руб.*") {
		t.Error("Неправильное отображение экономии по категориям")
	}
}

func TestGenerateReport_NoUnusedResources(t *testing.T) {
	a := NewAlerter(nil)
	report := metrics.OptimizerRecommendations{
		VPCRecommendations: []metrics.VPCOptimizerRecommendations{
			{IsUsed: true},
		},
		DNSRecommendations: []metrics.DNSOptimizerRecommendations{
			{IsUsed: true},
		},
	}

	result := a.generateReport(report)

	// Должны отсутствовать заголовки неиспользуемых ресурсов
	if strings.Contains(result, "Неиспользуемые IP-адреса") ||
		strings.Contains(result, "Неиспользуемые DNS зоны") {
		t.Error("Отчет не должен содержать секции неиспользуемых ресурсов")
	}

	// Должно быть сообщение об отсутствии оптимизации
	if !strings.Contains(result, "🎉 *Поздравляем!") {
		t.Error("Отсутствует сообщение об отсутствии оптимизации")
	}
}

func TestGenerateReport_PlatformOptimizationsOnly(t *testing.T) {
	a := NewAlerter(nil)
	report := metrics.OptimizerRecommendations{
		PlatformRecommendations: []metrics.PlatformOptimizerRecommendations{
			{
				NodeGroupId:     "ng-1",
				CurrentPlatform: "n1-highcpu",
				DesiredPlatform: "e2-custom",
				Savings:         1500.50,
			},
		},
	}

	result := a.generateReport(report)

	fmt.Println(result)

	if !strings.Contains(result, "`n1-highcpu` → `e2-custom`") {
		t.Error("Неправильное отображение платформ")
	}

	if !strings.Contains(result, "💸 *Итоговая экономия:*") ||
		!strings.Contains(result, "✨ *Всего: 1500.50 руб.*") {
		t.Error("Неправильное отображение экономии платформ")
	}
}
