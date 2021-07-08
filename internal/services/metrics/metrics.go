package metrics

import (
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	DiscordEventTriggers = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "discord_eventtriggers_total",
		Help: "Total number of Discord events triggered.",
	}, []string{"event"})

	DiscordCommandsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "discord_commands_processed_total",
		Help: "Total number of chat commands processed.",
	}, []string{"command"})

	DiscordGatewayPing = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "discord_gatewayping",
		Help: "The ping time in milliseconds to the discord API gateay.",
	})

	RestapiRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "restapi_requests_total",
		Help: "Total number of HTTP requests processed.",
	}, []string{"method", "status"})

	RestapiRequestTimes = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "restapi_requests_duration_seconds",
		Help: "Duration of all HTTP requests by method and status.",
		Buckets: []float64{
			0.000000001, // 1ns
			0.000000002,
			0.000000005,
			0.00000001, // 10ns
			0.00000002,
			0.00000005,
			0.0000001, // 100ns
			0.0000002,
			0.0000005,
			0.000001, // 1µs
			0.000002,
			0.000005,
			0.00001, // 10µs
			0.00002,
			0.00005,
			0.0001, // 100µs
			0.0002,
			0.0005,
			0.001, // 1ms
			0.002,
			0.005,
			0.01, // 10ms
			0.02,
			0.05,
			0.1, // 100 ms
			0.2,
			0.5,
			1.0, // 1s
			2.0,
			5.0,
			10.0, // 10s
			15.0,
			20.0,
			30.0,
		},
	}, []string{"method", "status"})

	RedisKeyCount = promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "redis_key_count",
		Help: "Number of Redis keys.",
	}, func() float64 {
		return redisW.Get("key_count")
	})

	RedisMemoryUsed = promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "redis_memory_used_bytes",
		Help: "Redis memory usage in bytes.",
	}, func() float64 {
		return redisW.Get("used_memory")
	})

	RedisCommandsProcessed = promauto.NewCounterFunc(prometheus.CounterOpts{
		Name: "redis_commands_executed_total",
		Help: "Total count of redis commands executed",
	}, func() float64 {
		return redisW.Get("total_commands_processed")
	})
)

// MetricsServer wraps a simple HTTP server serving
// a prometheus metrics endpoint.
type MetricsServer struct {
	server *http.Server
}

var redisW *redisWatcher

// NewMetricsServer initializes a new MectricsServer
// instance with the given addr and registers all
// instruments.
func NewMetricsServer(addr string, redis redis.Cmdable) (ms *MetricsServer, err error) {
	_, err = startPingWatcher(30 * time.Second)
	if err != nil {
		return
	}

	if redis != nil {
		redisW = newRedisWatcher(redis)
	}

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
