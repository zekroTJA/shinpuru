package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	DiscordEventTriggers = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "discord_eventtriggers_total",
		Help: "Total number of discord events triggered",
	}, []string{"event"})
)

// MetricsServer wraps a simple HTTP server serving
// a prometheus metrics endpoint.
type MetricsServer struct {
	server *http.Server
}

// NewMetricsServer initializes a new MectricsServer
// instance with the given addr and registers all
// instruments.
func NewMetricsServer(addr string) (ms *MetricsServer) {
	prometheus.MustRegister(
		DiscordEventTriggers)

	ms = new(MetricsServer)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	ms.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return
}

// ListenAndServeBlocking starts the listening loop of
// the web server which blocks the current goroutine.
func (ms *MetricsServer) ListenAndServeBlocking() error {
	return ms.server.ListenAndServe()
}
