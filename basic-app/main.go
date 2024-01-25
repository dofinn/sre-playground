package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpReqs = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	},
	[]string{"code", "method"},
)

func hello(w http.ResponseWriter, req *http.Request) {
	httpReqs.WithLabelValues(fmt.Sprint(http.StatusOK), req.Method).Inc()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = http.StatusText(http.StatusOK)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	// fmt.Fprintf(w, "hello\n")
}

func serverErrors(w http.ResponseWriter, req *http.Request) {
	httpReqs.WithLabelValues(fmt.Sprint(http.StatusInternalServerError), req.Method).Inc()

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = http.StatusText(http.StatusInternalServerError)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func clientErrors(w http.ResponseWriter, req *http.Request) {
	httpReqs.WithLabelValues(fmt.Sprint(http.StatusUnsupportedMediaType), req.Method).Inc()

	w.WriteHeader(http.StatusUnsupportedMediaType)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = http.StatusText(http.StatusUnsupportedMediaType)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	prometheus.MustRegister(httpReqs)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/fives", serverErrors)
	http.HandleFunc("/fours", clientErrors)
	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":8090", nil)
}
