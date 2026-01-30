package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	// API health metrics
	APIChecks *prometheus.CounterVec
	Latency   *prometheus.GaugeVec

	// Node status metrics
	Online   *prometheus.GaugeVec
	LastSeen *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		APIChecks: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "api_checks_total",
			Help: "Toyal number of API health checks perfomed.",
		}, []string{"provider", "status"}), //status tells you whether each check succeeded or failed

		Latency: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "api_latency_seconds",
			Help: "The value of the latency for the API request.",
		}, []string{"provider"}),

		Online: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "node_online",
			Help: "The online status of the node",
		}, []string{"provider", "node"}),

		LastSeen: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "node_last_seen_seconds",
			Help: "The timestamp for when the node was last online",
		}, []string{"provider", "node"}),
	}

	reg.MustRegister(m.APIChecks, m.Latency, m.LastSeen, m.Online)
	return m
}
