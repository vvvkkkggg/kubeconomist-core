package krr

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vvvkkkggg/kubeconomist-core/internal/config"
)

var testConfig = config.KrrAnalyzerConfig{
	PrometheusURL:        "https://monitoring.api.cloud.yandex.net/prometheus/workspaces/monjnbprk79m451udi83",
	PrometheusAuthHeader: "Bearer AQVN2t0W3J0xEoc0cVe0iieK5VcDSWrPZRIIZwI5",
	HistoryDuration:      "1",
}

func setupKrrAnalyzer() KrrAnalyzer {
	return *NewKrrAnalyzer(nil, testConfig)
}

func TestCallKRR(t *testing.T) {
	krr := setupKrrAnalyzer()
	result, err := krr.callKRR()
	assert.NoError(t, err)
	assert.Equal(t, result, "")
}
