package main

import (
	
	"log"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	api_up *prometheus.GaugeVec
}



func NewMetrics( reg prometheus.Registerer) *metrics {
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

	log.Fatal(http.ListenAndServe(":8080", nil))
}
