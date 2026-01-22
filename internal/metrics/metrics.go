package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	ApiUp        *prometheus.GaugeVec
	LatencyGauge *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		ApiUp: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "api_up_info",
			Help: "Checks API health.",
		}, []string{"provider"}),

		LatencyGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "api_latency_seconds",
			Help: "The value of the latency for the API request.",
		}, []string{"provider"}),
	}

	reg.MustRegister(m.ApiUp, m.LatencyGauge)
	return m
}
