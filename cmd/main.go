package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ElsaDevOps/Observability-EKS/internal/metrics"
	"github.com/ElsaDevOps/Observability-EKS/internal/provider"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func startTickerLoop(gauge *prometheus.GaugeVec, latencyGauge *prometheus.GaugeVec, providers []provider.Provider, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		for _, p := range providers {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			healthy, latency, err := p.CheckAPI(ctx)
			cancel()

			var value float64
			if healthy {
				value = 1.0
			} else {
				value = 0.0
			}

			gauge.With(prometheus.Labels{"provider": p.Name()}).Set(value)

			if err == nil {
				latencyGauge.With(prometheus.Labels{"provider": p.Name()}).Set(latency.Seconds())
			}

			if err != nil {
				log.Printf("Error for provider %s: %s", p.Name(), err)
			}
		}
	}
}

func main() {
	reg := prometheus.NewRegistry()
	m := metrics.NewMetrics(reg)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	h := provider.NewHeadscale("https://example.com/api/v1/node", "api-key")
	providers := []provider.Provider{h}

	go startTickerLoop(m.ApiUp, m.LatencyGauge, providers, 30*time.Second)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
