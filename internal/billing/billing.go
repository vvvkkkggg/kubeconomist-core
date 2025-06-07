package billing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vvvkkkggg/kubeconomist-core/internal/model"
)

const (
	baseURL             = "https://yandex.cloud/api/priceList/getPriceList"
	kubernetesServiceID = "dn2af04ph5otc5f23o1h"
	computeCloud        = "dn22pas77ftg9h3f2djj"
)

// Billing структура для работы с API биллинга Yandex Cloud
type Billing struct {
	client  *http.Client
	baseURL string

	computeCloudPrices []SKU

	mu sync.RWMutex
}

func New() *Billing {
	b := &Billing{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				IdleConnTimeout:     30 * time.Second,
				DisableCompression:  false,
				DisableKeepAlives:   false,
				MaxIdleConnsPerHost: 5,
			},
		},
		baseURL: baseURL,
	}

	return b
}

// PriceResponse структура ответа API
type PriceResponse struct {
	SKUs []SKU `json:"skus"`
}

// SKU представляет единицу хранения в системе ценообразования
type SKU struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	PricingUnit     string           `json:"pricingUnit"`
	ServiceID       string           `json:"serviceId"`
	UsageType       string           `json:"usageType"`
	Deprecated      bool             `json:"deprecated"`
	CreatedAt       int64            `json:"createdAt"`
	PricingVersions []PricingVersion `json:"pricingVersions"`
	EffectiveTime   int64            `json:"effectiveTime"`
}

// PricingVersion представляет версию ценообразования
type PricingVersion struct {
	ID                string            `json:"id"`
	PricingExpression PricingExpression `json:"pricingExpression"`
	EffectiveTime     int64             `json:"effectiveTime"`
}

// PricingExpression представляет выражение ценообразования
type PricingExpression struct {
	Quantum string `json:"quantum"`
	Rates   []Rate `json:"rates"`
}

// Rate представляет тарифную ставку
type Rate struct {
	StartPricingQuantity string `json:"startPricingQuantity"`
	UnitPrice            string `json:"unitPrice"`
}

// GetPrices получает текущие цены для указанного сервиса
func (b *Billing) GetPrices(ctx context.Context, serviceID string) ([]SKU, error) {
	params := url.Values{}
	params.Add("installationCode", "ru")
	params.Add("services[]", serviceID)
	params.Add("from", time.Now().Format("2006-01-02"))
	params.Add("to", time.Now().Format("2006-01-02"))
	//params.Add("pageSize", "50")
	params.Add("currency", "RUB")
	params.Add("lang", "ru")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s?%s", b.baseURL, params.Encode()),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var priceResponse PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&priceResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	filteredSKUs := make([]SKU, 0, len(priceResponse.SKUs))
	for _, sku := range priceResponse.SKUs {
		if !sku.Deprecated {
			filteredSKUs = append(filteredSKUs, sku)
		}
	}

	return filteredSKUs, nil
}

func (b *Billing) UpdatePricesCloudeCompute(ctx context.Context) {
	b.computeCloudPrices, _ = b.GetPrices(ctx, computeCloud)
}

// GetPricesForKubernetes получает цены для Kubernetes
func (b *Billing) GetPricesForComputeCloud() []SKU {
	return b.computeCloudPrices
}

// GetCurrentPrice возвращает текущую цену для SKU
func (sku *SKU) GetCurrentPrice() (string, error) {
	if len(sku.PricingVersions) == 0 {
		return "", fmt.Errorf("no pricing versions available")
	}

	latestVersion := sku.PricingVersions[0]
	if len(latestVersion.PricingExpression.Rates) == 0 {
		return "", fmt.Errorf("no rates available")
	}

	return latestVersion.PricingExpression.Rates[0].UnitPrice, nil
}

// GetEffectiveTime возвращает время вступления цены в силу
func (pv *PricingVersion) GetEffectiveTime() time.Time {
	return time.Unix(pv.EffectiveTime/1000, 0)
}

func platformMatcher(platform string) string {
	m := map[string]string{
		"standard-v1": "Intel Broadwell",
		"standard-v2": "Intel Cascade Lake",
		"standard-v3": "Intel Ice Lake",

		"highfreq-v3": "Intel Ice Lake (Compute-Optimized)",
		"amd-v1":      "AMD Zen 3",
	}

	return m[platform]
}

func genCPUNameForGrep(platform, coreFraction string) string {
	cpu := platformMatcher(platform)

	return fmt.Sprintf("%s. %s", cpu, coreFraction) + "%"
}

func (b *Billing) GetPriceCPURUB(platform string, coreFraction string, cpuCount model.CPUCount) (model.PriceRUB, error) {
	name := genCPUNameForGrep(platform, coreFraction)

	var foundedCPU SKU

	b.mu.RLock()

	for _, sku := range b.computeCloudPrices {
		ln := strings.ToLower(sku.Name)

		if strings.Contains(ln, strings.ToLower(name)) {
			foundedCPU = sku

			break
		}
	}

	b.mu.RUnlock()

	if len(foundedCPU.PricingVersions) != 1 {
		return 0, fmt.Errorf("unexpected length of pricing versions: %d",
			len(foundedCPU.PricingVersions))
	}

	if len(foundedCPU.PricingVersions[0].PricingExpression.Rates) != 1 {
		return 0, fmt.Errorf("unexpected length of pricing expressions rates: %d",
			len(foundedCPU.PricingVersions[0].PricingExpression.Rates))
	}

	price := foundedCPU.PricingVersions[0].PricingExpression.Rates[0].UnitPrice

	res, err := strconv.ParseFloat(price, 32)
	if err != nil {
		return 0, errors.New("failed to parse float")
	}

	return model.PriceRUB(res * float64(cpuCount)), nil
}

func (b *Billing) GetPriceRAMRUB(platform string, coreFraction string, ramCount model.RAMCount) (model.PriceRUB, error) {

	return 0, nil
}
