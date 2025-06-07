package billing

import (
	"sync"
	"testing"

	"github.com/vvvkkkggg/kubeconomist-core/internal/model"
)

func TestGenCPUNameForGrep(t *testing.T) {
	tests := []struct {
		name         string
		platform     string
		coreFraction string
		want         string
	}{
		{
			name:         "Standard V1 with full core",
			platform:     "standard-v1",
			coreFraction: "100",
			want:         "Intel Broadwell. 100%",
		},
		{
			name:         "Standard V2 with half core",
			platform:     "standard-v2",
			coreFraction: "50",
			want:         "Intel Cascade Lake. 50%",
		},
		{
			name:         "Standard V3 with quarter core",
			platform:     "standard-v3",
			coreFraction: "25",
			want:         "Intel Ice Lake. 25%",
		},
		{
			name:         "High-frequency V3",
			platform:     "highfreq-v3",
			coreFraction: "100",
			want:         "Intel Ice Lake (Compute-Optimized). 100%",
		},
		{
			name:         "AMD platform",
			platform:     "amd-v1",
			coreFraction: "75",
			want:         "AMD Zen 3. 75%",
		},
		{
			name:         "Unknown platform",
			platform:     "unknown-platform",
			coreFraction: "100",
			want:         ". 100%",
		},
		{
			name:         "Empty core fraction",
			platform:     "standard-v1",
			coreFraction: "",
			want:         "Intel Broadwell. %",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := genCPUNameForGrep(tt.platform, tt.coreFraction)
			if got != tt.want {
				t.Errorf("genCPUNameForGrep(%q, %q) = %q, want %q", tt.platform, tt.coreFraction, got, tt.want)
			}
		})
	}
}

// Вспомогательная функция для создания тестового биллинга
func createTestBilling(skuList []SKU) *Billing {
	return &Billing{
		computeCloudPrices: skuList,
		mu:                 sync.RWMutex{},
	}
}

func TestGetPriceCPURUB_Found(t *testing.T) {
	b := createTestBilling([]SKU{
		{
			Name: "Intel Cascade Lake. 100%",
			PricingVersions: []PricingVersion{
				{
					PricingExpression: PricingExpression{
						Rates: []Rate{
							{UnitPrice: "1.5"},
						},
					},
				},
			},
		},
	})

	price, _ := b.GetPriceCPURUB("standard-v2", "100", model.CPUCount(4))
	expected := model.PriceRUB(6.0) // 1.5 * 4 = 6.0

	if price != expected {
		t.Errorf("Expected %.2f, got %.2f", expected, price)
	}
}

func TestGetPriceCPURUB_EmptyPrice(t *testing.T) {
	b := createTestBilling([]SKU{
		{
			Name: "Intel Broadwell. 50%",
			PricingVersions: []PricingVersion{
				{
					PricingExpression: PricingExpression{
						Rates: []Rate{
							{UnitPrice: ""}, // Пустая цена
						},
					},
				},
			},
		},
	})

	price, _ := b.GetPriceCPURUB("standard-v1", "50", model.CPUCount(3))
	expected := model.PriceRUB(0.0) // Ошибка парсинга

	if price != expected {
		t.Errorf("Expected %.2f, got %.2f", expected, price)
	}
}

func TestGetPriceCPURUB_InvalidPriceFormat(t *testing.T) {
	b := createTestBilling([]SKU{
		{
			Name: "AMD Zen 3. 75%",
			PricingVersions: []PricingVersion{
				{
					PricingExpression: PricingExpression{
						Rates: []Rate{
							{UnitPrice: "invalid"}, // Нечисловое значение
						},
					},
				},
			},
		},
	})

	price, _ := b.GetPriceCPURUB("amd-v1", "75", model.CPUCount(2))
	expected := model.PriceRUB(0.0) // Ошибка парсинга

	if price != expected {
		t.Errorf("Expected %.2f, got %.2f", expected, price)
	}
}

func TestGetPriceCPURUB_MultipleRates(t *testing.T) {
	b := createTestBilling([]SKU{
		{
			Name: "Intel Ice Lake. 100%",
			PricingVersions: []PricingVersion{
				{
					PricingExpression: PricingExpression{
						Rates: []Rate{
							{UnitPrice: "2.0"},
							{UnitPrice: "3.0"},
						},
					},
				},
			},
		},
	})

	price, _ := b.GetPriceCPURUB("standard-v3", "100", model.CPUCount(3))
	expected := model.PriceRUB(0.0) // 2.0 * 3 = 6.0

	if price != expected {
		t.Errorf("Expected %.2f, got %.2f", expected, price)
	}
}

func TestGetPriceCPURUB_CaseInsensitive(t *testing.T) {
	b := createTestBilling([]SKU{
		{
			Name: "intel cascade lake. 100%", // В нижнем регистре
			PricingVersions: []PricingVersion{
				{
					PricingExpression: PricingExpression{
						Rates: []Rate{
							{UnitPrice: "1.25"},
						},
					},
				},
			},
		},
	})

	price, _ := b.GetPriceCPURUB("standard-v2", "100", model.CPUCount(4))
	expected := model.PriceRUB(5.0) // 1.25 * 4 = 5.0

	if price != expected {
		t.Errorf("Expected %.2f, got %.2f", expected, price)
	}
}

func TestGetPriceCPURUB_NotFound(t *testing.T) {
	b := createTestBilling([]SKU{
		{
			Name: "Different SKU",
			PricingVersions: []PricingVersion{
				{
					PricingExpression: PricingExpression{
						Rates: []Rate{
							{UnitPrice: "2.0"},
						},
					},
				},
			},
		},
	})

	price, err := b.GetPriceCPURUB("standard-v1", "100", model.CPUCount(2))
	expected := model.PriceRUB(0.0)

	if price != expected {
		t.Errorf("Expected price %.2f, got %.2f", expected, price)
	}

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "unexpected length of pricing versions: 0"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}
