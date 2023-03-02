package inits

import (
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/metrics"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/rogu/log"
)

func InitMetrics(container di.Container) (ms *metrics.MetricsServer) {
	var err error

	cfg := container.Get(static.DiConfig).(config.Provider)

	log := log.Tagged("METRICS")

	if cfg.Config().Metrics.Enable {
		log.Info().Msg("Initializing metrics server ...")

		if cfg.Config().Metrics.Addr == "" {
			cfg.Config().Metrics.Addr = ":9091"
		}

		redis := container.Get(static.DiRedis).(redis.Cmdable)
		ms, err = metrics.NewMetricsServer(cfg.Config().Metrics.Addr, redis)
		if err != nil {
			log.Fatal().Err(err).Msg("failed initializing metrics server")
		}

		go func() {
			log.Info().Field("addr", cfg.Config().Metrics.Addr).Msg("Metrics server running")
			if err := ms.ListenAndServeBlocking(); err != nil {
				log.Fatal().Err(err).Msg("Failed setting up metrics server")
			}
		}()
	}

	return
}
