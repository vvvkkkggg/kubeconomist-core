package billing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vvvkkkggg/kubeconomist-core/internal/analyzers"
	"github.com/vvvkkkggg/kubeconomist-core/internal/model"
)

const (
	baseURL             = "https://yandex.cloud/api/priceList/getPriceList"
	kubernetesServiceID = "dn2af04ph5otc5f23o1h"
	computeCloud        = "dn22pas77ftg9h3f2djj"
	imageRegistry       = "dn2tng436tjcn7cjudv1"
	objectStorage       = "dn2li5qddoc5cad2n6br"
)

var _ analyzers.Billing = &Billing{}

// Billing структура для работы с API биллинга Yandex Cloud
type Billing struct {
	client  *http.Client
	baseURL string

	computeCloudPrices  []SKU
	imageRegistryPrices []SKU
	objectStoragePrices []SKU

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

func (b *Billing) UpdatePricesCloudeCompute(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	var err error

	b.computeCloudPrices, err = b.GetPrices(ctx, computeCloud)
	if err != nil {
		return err
	}

	return nil
}

func (b *Billing) UpdatePricesContainerRegistry(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	var err error

	b.imageRegistryPrices, err = b.GetPrices(ctx, imageRegistry)
	if err != nil {
		return err
	}

	return nil
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

func genRAMNameForGrep(platform string) string {
	cpu := platformMatcher(platform)

	return fmt.Sprintf("%s. RAM", cpu)
}

func (b *Billing) GetPriceRAMRUB(platform string, ramCount model.RAMCount) (model.PriceRUB, error) {
	name := genRAMNameForGrep(platform)

	var foundedRAM SKU

	b.mu.RLock()

	for _, sku := range b.computeCloudPrices {
		ln := strings.ToLower(sku.Name)

		if ln == strings.ToLower(name) {
			foundedRAM = sku

			break
		}
	}

	b.mu.RUnlock()

	if len(foundedRAM.PricingVersions) != 1 {
		return 0, fmt.Errorf("unexpected length of pricing versions: %d",
			len(foundedRAM.PricingVersions))
	}

	if len(foundedRAM.PricingVersions[0].PricingExpression.Rates) != 1 {
		return 0, fmt.Errorf("unexpected length of pricing expressions rates: %d",
			len(foundedRAM.PricingVersions[0].PricingExpression.Rates))
	}

	price := foundedRAM.PricingVersions[0].PricingExpression.Rates[0].UnitPrice

	res, err := strconv.ParseFloat(price, 32)
	if err != nil {
		return 0, errors.New("failed to parse float")
	}

	return model.PriceRUB(res * float64(ramCount)), nil
}

func (b *Billing) GetContainerRegistryPriceRUB() (model.PriceRUB, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, h := range b.imageRegistryPrices {
		if h.Name == "Container Registry – занятое место в хранилище" {
			if len(h.PricingVersions) != 1 {
				return 0, fmt.Errorf("unexpected length of pricing versions: %d",
					len(h.PricingVersions))
			}

			if len(h.PricingVersions[0].PricingExpression.Rates) != 1 {
				return 0, fmt.Errorf("unexpected length of pricing expressions rates: %d",
					len(h.PricingVersions[0].PricingExpression.Rates))
			}

			price := h.PricingVersions[0].PricingExpression.Rates[0].UnitPrice

			res, err := strconv.ParseFloat(price, 32)
			if err != nil {
				return 0, errors.New("failed to parse float")
			}

			return model.PriceRUB(res), nil
		}
	}

	return 0.0, errors.New("needed block is not found")
}

type ObjectStoragePrices struct {
	StoragePrices struct {
		Ice      float64 `json:"icePricePerGBHour"`
		Standard float64 `json:"standardPricePerGBHour"`
		Cold     float64 `json:"coldPricePerGBHour"`
	} `json:"storagePrices"`

	IceOperations struct {
		Get        float64 `json:"getPerRequest"`
		Head       float64 `json:"headPerRequest"`
		List       float64 `json:"listPerRequest"`
		Options    float64 `json:"optionsPerRequest"`
		Patch      float64 `json:"patchPerRequest"`
		Post       float64 `json:"postPerRequest"`
		Put        float64 `json:"putPerRequest"`
		Transition float64 `json:"transitionPerRequest"`
	} `json:"iceOperations"`

	StandardOperations struct {
		Get     float64 `json:"getPerRequest"`
		Head    float64 `json:"headPerRequest"`
		List    float64 `json:"listPerRequest"`
		Options float64 `json:"optionsPerRequest"`
		Patch   float64 `json:"patchPerRequest"`
		Post    float64 `json:"postPerRequest"`
		Put     float64 `json:"putPerRequest"`
	} `json:"standardOperations"`

	ColdOperations struct {
		Get        float64 `json:"getPerRequest"`
		Head       float64 `json:"headPerRequest"`
		List       float64 `json:"listPerRequest"`
		Options    float64 `json:"optionsPerRequest"`
		Patch      float64 `json:"patchPerRequest"`
		Post       float64 `json:"postPerRequest"`
		Put        float64 `json:"putPerRequest"`
		Transition float64 `json:"transitionPerRequest"`
	} `json:"coldOperations"`
}

func (b *Billing) UpdateObjectStoragePrice(ctx context.Context) error {
	skus, err := b.GetPrices(ctx, objectStorage)
	if err != nil {
		return err
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.objectStoragePrices = skus

	return nil
}

func (b *Billing) GetObjectStoragePricesRUB(ctx context.Context) (*ObjectStoragePrices, error) {
	prices := &ObjectStoragePrices{}

	b.mu.RLock()
	defer b.mu.RUnlock()

	skus := b.objectStoragePrices

	for _, sku := range skus {
		nameLower := strings.ToLower(sku.Name)

		// Обработка цен на хранение
		if strings.Contains(nameLower, "занятое место") {
			price, err := extractPrice(sku)
			if err != nil {
				return nil, err
			}

			switch {
			case strings.Contains(nameLower, "ледяном"):
				prices.StoragePrices.Ice = price
			case strings.Contains(nameLower, "стандартном"):
				prices.StoragePrices.Standard = price
			case strings.Contains(nameLower, "холодном"):
				prices.StoragePrices.Cold = price
			}
		}

		if strings.Contains(nameLower, "операции") || strings.Contains(nameLower, "удаление") {
			price, err := extractPrice(sku)
			if err != nil {
				return nil, err
			}

			switch {
			case strings.Contains(nameLower, "ледяное") || strings.Contains(nameLower, "ледяном"):
				switch {
				case strings.Contains(nameLower, "get"):
					prices.IceOperations.Get = price
				case strings.Contains(nameLower, "head"):
					prices.IceOperations.Head = price
				case strings.Contains(nameLower, "list"):
					prices.IceOperations.List = price
				case strings.Contains(nameLower, "options"):
					prices.IceOperations.Options = price
				case strings.Contains(nameLower, "patch"):
					prices.IceOperations.Patch = price
				case strings.Contains(nameLower, "post"):
					prices.IceOperations.Post = price
				case strings.Contains(nameLower, "put"):
					prices.IceOperations.Put = price
				case strings.Contains(nameLower, "transition"):
					prices.IceOperations.Transition = price
				}

			case strings.Contains(nameLower, "стандартное") || strings.Contains(nameLower, "стандартном"):
				switch {
				case strings.Contains(nameLower, "get"):
					prices.StandardOperations.Get = price
				case strings.Contains(nameLower, "head"):
					prices.StandardOperations.Head = price
				case strings.Contains(nameLower, "list"):
					prices.StandardOperations.List = price
				case strings.Contains(nameLower, "options"):
					prices.StandardOperations.Options = price
				case strings.Contains(nameLower, "patch"):
					prices.StandardOperations.Patch = price
				case strings.Contains(nameLower, "post"):
					prices.StandardOperations.Post = price
				case strings.Contains(nameLower, "put"):
					prices.StandardOperations.Put = price
				}

			case strings.Contains(nameLower, "холодное") || strings.Contains(nameLower, "холодном"):
				switch {
				case strings.Contains(nameLower, "get"):
					prices.ColdOperations.Get = price
				case strings.Contains(nameLower, "head"):
					prices.ColdOperations.Head = price
				case strings.Contains(nameLower, "list"):
					prices.ColdOperations.List = price
				case strings.Contains(nameLower, "options"):
					prices.ColdOperations.Options = price
				case strings.Contains(nameLower, "patch"):
					prices.ColdOperations.Patch = price
				case strings.Contains(nameLower, "post"):
					prices.ColdOperations.Post = price
				case strings.Contains(nameLower, "put"):
					prices.ColdOperations.Put = price
				case strings.Contains(nameLower, "transition"):
					prices.ColdOperations.Transition = price
				}
			}
		}
	}

	return prices, nil
}

// extractPrice извлекает и нормализует цену из SKU
func extractPrice(sku SKU) (float64, error) {
	if len(sku.PricingVersions) == 0 {
		return 0, nil
	}

	// Выбираем самую актуальную версию ценообразования
	latest := sku.PricingVersions[0]
	for _, v := range sku.PricingVersions {
		if v.EffectiveTime > latest.EffectiveTime {
			latest = v
		}
	}

	// Сортируем ставки по возрастанию лимита
	sort.Slice(latest.PricingExpression.Rates, func(i, j int) bool {
		q1, _ := strconv.ParseFloat(latest.PricingExpression.Rates[i].StartPricingQuantity, 64)
		q2, _ := strconv.ParseFloat(latest.PricingExpression.Rates[j].StartPricingQuantity, 64)
		return q1 < q2
	})

	// Ищем первую ненулевую ставку
	for _, rate := range latest.PricingExpression.Rates {
		price, err := normalizePrice(rate.UnitPrice, sku.PricingUnit)
		if err != nil {
			continue
		}
		if price > 0 {
			return price, nil
		}
	}

	return 0, nil
}

// normalizePrice преобразует строку цены в float64 с учетом единиц измерения
func normalizePrice(priceStr, unit string) (float64, error) {
	// Обработка диапазонов цен (берем максимальное значение)
	if strings.Contains(priceStr, "—") {
		parts := strings.Split(priceStr, "—")
		if len(parts) > 1 {
			priceStr = parts[1] // Берем вторую (максимальную) часть
		}
	}

	// Удаляем все нечисловые символы кроме точек и минусов
	cleaned := strings.Builder{}
	decimalSeparator := false
	for _, r := range priceStr {
		switch {
		case r >= '0' && r <= '9':
			cleaned.WriteRune(r)
		case r == ',' || r == '.':
			if !decimalSeparator {
				cleaned.WriteRune('.')
				decimalSeparator = true
			}
		}
	}

	result, err := strconv.ParseFloat(cleaned.String(), 64)
	if err != nil {
		return 0, err
	}

	// Нормализация по единицам измерения
	switch {
	case strings.Contains(unit, "10k"):
		result /= 10000.0 // Для операций за 10k запросов
	case strings.Contains(unit, "1k"):
		result /= 1000.0 // Для операций за 1k запросов
		// Для gbyte*hour нормализация не требуется
	}

	return result, nil
}
