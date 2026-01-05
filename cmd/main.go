package main

import (
	"net/http"
	"log"
	"io"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)


func main() {
	h1 := func(w http.ResponseWriter, r *http.Request){
		io.WriteString (w, "hello")

	}

	http.HandleFunc("/metrics", h1)
	http.Handle("/metrics", promhttp.Handler())


	log.Fatal(http.ListenAndServe(":8080", nil))
}