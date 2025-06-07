package metrics

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vvvkkkggg/kubeconomist-core/internal/config"
)

func ListenAndServe(ctx context.Context, cfg config.MetricsConfig, metricsCollector ...prometheus.Collector) error {
	var (
		chErr    = make(chan error)
		mux      = http.NewServeMux()
		registry = prometheus.NewRegistry()
	)
	registry.MustRegister(metricsCollector...)

	mux.Handle("/metrics", promhttp.InstrumentMetricHandler(
		registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{})),
	)
	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: mux,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			chErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		ctxShutdown, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
		defer cancelFunc()
		return server.Shutdown(ctxShutdown) //nolint:contextcheck // deliberately created context for shutdown
	case err := <-chErr:
		return err
	}
}
