package webserver

import (
	"strings"

	routing "github.com/qiangxue/fasthttp-routing"
)

func (ws *WebServer) handleMetrics(ctx *routing.Context) error {
	// method := string(ctx.Method())
	// endpoint := getUnparameterizedPath(ctx)

	// metrics.RestapiRequest

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
