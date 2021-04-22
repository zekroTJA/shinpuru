package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/services/metrics"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

func InitMetrics(container di.Container) (ms *metrics.MetricsServer) {
	var err error

	cfg := container.Get(static.DiConfig).(*config.Config)

	if cfg.Metrics != nil && cfg.Metrics.Enable {
		if cfg.Metrics.Addr == "" {
			cfg.Metrics.Addr = ":9091"
		}

		ms, err = metrics.NewMetricsServer(cfg.Metrics.Addr)
		if err != nil {
			util.Log.Fatalf("failed initializing metrics server: %s", err.Error())
		}

		go func() {
			util.Log.Infof("Metrics server listening on %s...", cfg.Metrics.Addr)
			if err := ms.ListenAndServeBlocking(); err != nil {
				util.Log.Fatalf("failed setting up metrics server: %s", err.Error())
			}
		}()
	}

	return
}
