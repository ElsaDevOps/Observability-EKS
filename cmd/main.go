package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ElsaDevOps/Observability-EKS/internal/config"
	"github.com/ElsaDevOps/Observability-EKS/internal/metrics"
	"github.com/ElsaDevOps/Observability-EKS/internal/provider"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var configPath = flag.String("config", "config.yaml", "path to config file")
var port = flag.Int("port", 8080, "port for server to listen on")

func startTickerLoop(m *metrics.Metrics, providers []provider.Provider, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		for _, p := range providers {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			healthy, latency, err := p.CheckAPI(ctx)
			cancel()

			var status string
			if healthy {
				status = "success"
			} else {
				status = "failure"
			}

			m.APIChecks.With(prometheus.Labels{"provider": p.Name(), "status": status}).Inc()

			if err == nil {
				m.Latency.With(prometheus.Labels{"provider": p.Name()}).Set(latency.Seconds())
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
				continue
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
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d providers, interval: %s", len(cfg.Providers), cfg.DefaultInterval)

	reg := prometheus.NewRegistry()
	m := metrics.NewMetrics(reg)

	var providers []provider.Provider

	for _, p := range cfg.Providers {
		factory, ok := provider.GetRegistry(p.Name)
		if !ok {
			log.Fatalf("unknown provider: %s", p.Name)
		}

		prov := factory(p)
		providers = append(providers, prov)
	}

	go startTickerLoop(m, providers, cfg.DefaultInterval)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
