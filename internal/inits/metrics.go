package inits

import (
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/metrics"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

func InitMetrics(container di.Container) (ms *metrics.MetricsServer) {
	var err error

	cfg := container.Get(static.DiConfig).(config.Provider)

	if cfg.Config().Metrics.Enable {
		if cfg.Config().Metrics.Addr == "" {
			cfg.Config().Metrics.Addr = ":9091"
		}

		redis := container.Get(static.DiRedis).(redis.Cmdable)
		ms, err = metrics.NewMetricsServer(cfg.Config().Metrics.Addr, redis)
		if err != nil {
			logrus.WithError(err).Fatal("failed initializing metrics server")
		}

		go func() {
			logrus.WithField("addr", cfg.Config().Metrics.Addr).Info("Metrics server started")
			if err := ms.ListenAndServeBlocking(); err != nil {
				logrus.WithError(err).Fatal("failed setting up metrics server")
			}
		}()
	}

	return
}
