package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

func InitPermissionMiddleware(container di.Container) *middleware.PermissionsMiddleware {

	cfg := container.Get(static.DiConfig).(*config.Config)
	db := container.Get(static.DiDatabase).(database.Database)

	return middleware.NewPermissionMiddleware(db, cfg)
}
