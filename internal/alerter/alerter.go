package alerter

import (
	"context"
	"fmt"
	"strings"

	"github.com/vvvkkkggg/kubeconomist-core/internal/metrics"
)

// Messenger интерфейс для отправки сообщений
type Messenger interface {
	SendMessage(ctx context.Context, text string) error
}

// Определения остальных структур (как в задании)

type Alerter struct {
	messenger Messenger
}

func NewAlerter(m Messenger) *Alerter {
	return &Alerter{messenger: m}
}

func (a *Alerter) SendReport(ctx context.Context, report metrics.OptimizerRecommendations) error {
	message := a.generateReport(report)
	return a.messenger.SendMessage(ctx, message)
}

func (a *Alerter) generateReport(r metrics.OptimizerRecommendations) string {
	var sb strings.Builder

	// Заголовок отчета
	sb.WriteString("🚀 *Отчет оптимизации Kubernetes кластера*\n\n")
	sb.WriteString("_Ресурсы, которые можно оптимизировать:_\n\n")

	used := false

	// Раздел VPC
	if len(r.VPCRecommendations) > 0 {
		used = true

		sb.WriteString("🔌 *Неиспользуемые IP-адреса:*\n")
		for _, vpc := range r.VPCRecommendations {
			status := "🟢 Используется"
			if !vpc.IsUsed {
				status = "🔴 Не используется"
			}
			sb.WriteString(fmt.Sprintf(
				"- `%s` (Cloud: `%s`, Folder: `%s`) %s\n",
				vpc.IPAddress,
				vpc.CloudID,
				vpc.FolderID,
				status,
			))
		}
		sb.WriteString("\n")
	}

	// Раздел DNS
	if len(r.DNSRecommendations) > 0 {
		used = true

		sb.WriteString("🌐 *Неиспользуемые DNS зоны:*\n")
		for _, dns := range r.DNSRecommendations {
			status := "🟢 Используется"
			if !dns.IsUsed {
				status = "🔴 Не используется"
			}
			sb.WriteString(fmt.Sprintf(
				"- `%s` (Cloud: `%s`, Folder: `%s`) %s\n",
				dns.ZoneId,
				dns.CloudID,
				dns.FolderID,
				status,
			))
		}
		sb.WriteString("\n")
	}

	// Раздел Node
	nodeSavings := 0.0
	if len(r.NodeRecommendations) > 0 {
		used = true

		sb.WriteString("🖥️ *Оптимизация нод:*\n")
		for _, node := range r.NodeRecommendations {
			saving := node.CurrentPrice - node.DesiredPrice
			nodeSavings += saving

			sb.WriteString(fmt.Sprintf(
				"* Instance: `%s`\n"+
					"  - CPU: %d → %d ядер\n"+
					"  - RAM: %d → %d MB\n"+
					"  - Экономия: *%.2f руб.*\n\n",
				node.InstanceId,
				node.CurrentCores,
				node.DesiredCores,
				node.CurrentMemory,
				node.DesiredMemory,
				saving,
			))
		}
	}

	// Раздел Platform
	platformSavings := 0.0
	if len(r.PlatformRecommendations) > 0 {
		used = true

		sb.WriteString("📦 *Оптимизация платформ:*\n")
		for _, platform := range r.PlatformRecommendations {
			platformSavings += platform.Savings
			sb.WriteString(fmt.Sprintf(
				"* NodeGroup: `%s`\n"+
					"  - Платформа: `%s` → `%s`\n"+
					"  - Экономия: *%.2f руб.*\n\n",
				platform.NodeGroupId,
				platform.CurrentPlatform,
				platform.DesiredPlatform,
				platform.Savings,
			))
		}
	}

	// Итоговая экономия
	totalSavings := nodeSavings + platformSavings
	if totalSavings > 0 {
		used = true

		sb.WriteString("💸 *Итоговая экономия:*\n")
		sb.WriteString(fmt.Sprintf(
			"- Ноды: *%.2f руб.*\n- Платформы: *%.2f руб.*\n"+
				"✨ *Всего: %.2f руб.*\n",
			nodeSavings,
			platformSavings,
			totalSavings,
		))
	}

	if !used {
		sb.WriteString("🎉 *Поздравляем! Оптимизация не требуется, кластер работает оптимально*")
	} else {
		sb.WriteString("\n_Оптимизация поможет снизить затраты и повысить эффективность кластера_")
	}

	return sb.String()
}
