package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	failedRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "istio_ca_shim_step_failed_requests_total",
		Help: "The total number of failed requests.",
	})
	requestsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "istio_ca_shim_step_processed_requests_total",
		Help: "The total number of processed requests.",
	})
)

func IncFailedRequests() {
	failedRequests.Inc()
}

func IncProcessedRequests() {
	requestsProcessed.Inc()
}

func ServeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		_ = http.ListenAndServe(":2112", nil)
	}()
}
