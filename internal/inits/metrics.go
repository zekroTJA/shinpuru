package inits

import (
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/metrics"
)

func InitMetrics(cfg *config.Config) {
	if cfg.Metrics != nil && cfg.Metrics.Enable {
		if cfg.Metrics.Addr == "" {
			cfg.Metrics.Addr = ":9091"
		}

		ms := metrics.NewMetricsServer(cfg.Metrics.Addr)

		go func() {
			util.Log.Infof("Metrics server listening on %s...", cfg.Metrics.Addr)
			if err := ms.ListenAndServeBlocking(); err != nil {
				util.Log.Fatalf("failed setting up metrics server: %s", err.Error())
			}
		}()
	}
}
