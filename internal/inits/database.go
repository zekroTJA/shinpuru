package inits

import (
	"strings"

	goredis "github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/database/mysql"
	"github.com/zekroTJA/shinpuru/internal/services/database/redis"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

func InitDatabase(container di.Container) database.Database {
	var db database.Database
	var err error

	cfg := container.Get(static.DiConfig).(*config.Config)

	switch strings.ToLower(cfg.Database.Type) {
	case "mysql", "mariadb":
		db = new(mysql.MysqlMiddleware)
		err = db.Connect(cfg.Database.MySql)
	}

	if m, ok := db.(database.Migration); ok {
		logrus.Info("Checking database for migrations and apply if needed...")
		if err = m.Migrate(); err != nil {
			logrus.WithError(err).Fatal("Database migration failed")
		}
	} else {
		logrus.Warning("Skip database migration: middleware does not support migrations")
	}

	if cfg.Cache.CacheDatabase {
		rd := container.Get(static.DiRedis).(*goredis.Client)
		db = redis.NewRedisMiddleware(db, rd)
		logrus.Info("Enabled redis as database cache")
	} else {
		logrus.Warning("Database cache is disabled! You can enbale it in the config (.cache.cachedatabase).")
	}

	if err != nil {
		logrus.WithError(err).Fatal("Failed connecting to database")
	}
	logrus.Info("Connected to database")

	return db
}
