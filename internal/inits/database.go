package inits

import (
	"strings"

	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/database/mysql"
	"github.com/zekroTJA/shinpuru/internal/services/database/redis"
	"github.com/zekroTJA/shinpuru/internal/services/database/sqlite"
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
	case "sqlite", "sqlite3":
		db = new(sqlite.SqliteMiddleware)
		err = db.Connect(cfg.Database.Sqlite)
		printSqliteWraning()
	}

	if m, ok := db.(database.Migration); ok {
		logrus.Info("Checking database for migrations and apply if needed...")
		if err = m.Migrate(); err != nil {
			logrus.WithError(err).Fatal("Database migration failed")
		}
	} else {
		logrus.Warning("Skip database migration: middleware does not support migrations")
	}

	if cfg.Database.Redis != nil && cfg.Database.Redis.Enable {
		db = redis.NewRedisMiddleware(cfg.Database.Redis, db)
		logrus.Info("Enabled redis as database cache")
	}

	if err != nil {
		logrus.WithError(err).Fatal("Failed connecting to database")
	}
	logrus.Info("Connected to database")

	return db
}

func printSqliteWraning() {
	logrus.Warning("You are currently using the SQLite Database Driver, which is marked " +
		"as DEPRECATED and will be removed in the upcoming version!")
}
