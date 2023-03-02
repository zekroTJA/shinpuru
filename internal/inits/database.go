package inits

import (
	"strings"

	goredis "github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/database/mysql"
	"github.com/zekroTJA/shinpuru/internal/services/database/redis"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/rogu/log"
)

func InitDatabase(container di.Container) database.Database {
	var db database.Database
	var err error

	cfg := container.Get(static.DiConfig).(config.Provider)

	log := log.Tagged("Database")

	drv := strings.ToLower(cfg.Config().Database.Type)
	log.Info().Field("driver", drv).Msg("Initializing database ...")

	switch drv {
	case "mysql", "mariadb":
		db = mysql.New()
		err = db.Connect(cfg.Config().Database.MySql)
	default:
		log.Fatal().Field("driver", drv).Msg("Unsupported database driver")
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Failed connecting to database")
	}

	if m, ok := db.(database.Migration); ok {
		log.Info().Msg("Checking database for migrations and apply if needed...")
		if err = m.Migrate(); err != nil {
			log.Fatal().Err(err).Msg("Database migration failed")
		}
	} else {
		log.Warn().Msg("Skip database migration: middleware does not support migrations")
	}

	// Redis Database Cache
	if cfg.Config().Cache.CacheDatabase {
		rd := container.Get(static.DiRedis).(*goredis.Client)
		db = redis.NewRedisMiddleware(db, rd)
		log.Info().Msg("Enabled Redis as database cache")
	} else {
		log.Warn().Msg("Database cache is disabled! You can enbale it in the config (.cache.cachedatabase).")
	}

	log.Info().Msg("Connected to database")

	return db
}
