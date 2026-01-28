package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ElsaDevOps/Observability-EKS/internal/config"
	"github.com/ElsaDevOps/Observability-EKS/internal/metrics"
	"github.com/ElsaDevOps/Observability-EKS/internal/provider"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func startTickerLoop(m *metrics.Metrics, providers []provider.Provider, interval time.Duration) {
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

			m.APIUp.With(prometheus.Labels{"provider": p.Name()}).Set(value)

			if err == nil {
				m.LatencyGauge.With(prometheus.Labels{"provider": p.Name()}).Set(latency.Seconds())
			}

			if err != nil {
				log.Printf("Error for provider %s: %s", p.Name(), err)
				continue
			}

			ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
			nodes, err := p.ListNodes(ctx2)
			cancel2()

			if err != nil {
				log.Printf("Error for provider %s: %s", p.Name(), err)
			}

			for _, node := range nodes {
				var value float64
				if node.Online {
					value = 1.0
				} else {
					value = 0.0
				}

				m.Online.With(prometheus.Labels{"provider": p.Name(), "node": node.Name}).Set(value)
				m.LastSeen.With(prometheus.Labels{"provider": p.Name(), "node": node.Name}).Set(float64(node.LastSeen.Unix()))
			}

		}
	}
}

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d providers, interval: %s", len(cfg.Providers), cfg.Interval)

	reg := prometheus.NewRegistry()
	m := metrics.NewMetrics(reg)

	var providers []provider.Provider

	for _, p := range cfg.Providers {
		switch p.Name {
		case "headscale":
			h := provider.NewHeadscale(p.URL, p.APIKey)
			providers = append(providers, h)
		case "tailscale":
			t := provider.NewTailscale(p.TailnetID, p.APIKey)
			providers = append(providers, t)
		default:
			log.Printf("unknown provider: %s", p.Name)
		}
	}

	go startTickerLoop(m, providers, cfg.Interval)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
