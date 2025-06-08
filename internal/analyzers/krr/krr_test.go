package krr

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vvvkkkggg/kubeconomist-core/internal/config"
)

const reportName = "testdata/test_report.json"

var testConfig = config.KrrAnalyzerConfig{
	PrometheusURL:        "https://monitoring.api.cloud.yandex.net/prometheus/workspaces/monjnbprk79m451udi83",
	PrometheusAuthHeader: "Bearer AQVN2t0W3J0xEoc0cVe0iieK5VcDSWrPZRIIZwI5",
	HistoryDuration:      "1",
}

func setupKrrAnalyzer() KrrAnalyzer {
	return *NewKrrAnalyzer(nil, testConfig)
}

func TestLoadReport(t *testing.T) {
	krr := setupKrrAnalyzer()
	result, err := krr.loadReport(reportName)
	require.NoError(t, err)
	// assert.Equal(t, nil, result)
	_ = result
}
