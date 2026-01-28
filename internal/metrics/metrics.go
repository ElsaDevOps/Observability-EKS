package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	APIUp        *prometheus.GaugeVec
	LatencyGauge *prometheus.GaugeVec
	Online       *prometheus.GaugeVec
	LastSeen     *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		APIUp: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "api_up_info",
			Help: "Checks API health.",
		}, []string{"provider"}),

		LatencyGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
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

	reg.MustRegister(m.APIUp, m.LatencyGauge, m.LastSeen, m.Online)
	return m
}
