package prometheus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusClient is a struct that holds the prometheus client
type PrometheusClient struct{}

// NewPrometheusClient returns a new prometheus client
func NewPrometheusClient() *PrometheusClient {
	return &PrometheusClient{}
}

// prometheus metric initialization
var (
	DinRequestCount *prometheus.CounterVec
)

// RegisterMetrics registers the prometheus metrics
func RegisterMetrics() {
	// Register request count metric for inbound din http requests
	DinRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "din_http_request_count",
			Help: "Metric for counting din http requests with service, method, provider, and host_name labels",
		},
		[]string{"service", "method", "provider", "host_name", "res_status", "res_latency", "health_status", "block_number"},
	)
	prometheus.MustRegister(DinRequestCount)
}

type PromRequestMetricData struct {
	Method       string
	Service      string
	Provider     string
	HostName     string
	ResStatus    int
	ResLatency   time.Duration
	HealthStatus string
	BlockNumber  string
}

// handleRequestMetric increments prometheus metric based on request data passed in
func (p *PrometheusClient) HandleRequestMetric(reqBodyBytes []byte, data *PromRequestMetricData) {
	// First extract method data from body
	// define struct to hold request data
	var requestBody struct {
		Method string `json:"method,omitempty"`
	}
	err := json.Unmarshal(reqBodyBytes, &requestBody)
	if err != nil {
		fmt.Printf("Error decoding request body: %v", http.StatusBadRequest)
	}
	var method string
	if requestBody.Method != "" {
		method = requestBody.Method
	}

	service := strings.TrimPrefix(data.Service, "/")
	latency := strconv.FormatInt(data.ResLatency.Milliseconds(), 10)
	status := strconv.Itoa(data.ResStatus)

	// Increment prometheus metric based on request data
	DinRequestCount.WithLabelValues(service, method, data.Provider, data.HostName, status, latency, data.HealthStatus, data.BlockNumber).Inc()
}
