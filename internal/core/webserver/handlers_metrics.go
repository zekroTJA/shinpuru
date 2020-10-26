package webserver

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/zekroTJA/shinpuru/internal/util/metrics"
)

func (ws *WebServer) handleMetrics(ctx *routing.Context) error {
	method := strings.ToUpper(string(ctx.Method()))
	endpoint := strings.ToLower(getUnparameterizedPath(ctx))

	metrics.RestapiRequests.
		With(prometheus.Labels{"endpoint": endpoint, "method": method}).
		Add(1)

	return nil
}

var urlParameterNames = []string{"id", "hexcode", "guildid", "backupid", "memberid"}

func getUnparameterizedPath(ctx *routing.Context) string {
	path := string(ctx.Path())

	for _, paramName := range urlParameterNames {
		paramValue := ctx.Param(paramName)
		if paramValue == "" {
			continue
		}

		strings.ReplaceAll(path, paramValue, ":"+paramName)
	}

	return path
}
