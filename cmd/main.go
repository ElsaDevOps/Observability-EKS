package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ElsaDevOps/Observability-EKS/internal/provider"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	api_up *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		api_up: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "api_up_info",
			Help: "Checks API health.",
		},
			[]string{"provider"}),
	}

	reg.MustRegister(m.api_up)
	return m
}

func main() {
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)
	m.api_up.WithLabelValues("headscale").Set(1)
	m.api_up.WithLabelValues("tailscale").Set(0)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	h := provider.NewHeadscale("https://example.com/api/v1/node", "api-key")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	healthy, latency, err := h.CheckAPI(ctx)
	fmt.Printf("Healthy: %v, Latency: %v, Error: %v\n", healthy, latency, err)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
