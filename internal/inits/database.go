package inits

import (
	"strings"

	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/util"
)

func InitDatabase(databaseCfg *config.DatabaseType) database.Database {
	var db database.Database
	var err error

	switch strings.ToLower(databaseCfg.Type) {
	case "mysql", "mariadb":
		db = new(middleware.MySQLDriver)
		err = db.Connect(databaseCfg.MySql)
	case "sqlite", "sqlite3":
		db = new(middleware.SqliteDriver)
		err = db.Connect(databaseCfg.Sqlite)
		printSqliteWraning()
	}

	if databaseCfg.Redis != nil && databaseCfg.Redis.Enable {
		db = middleware.NewRedisMiddleware(databaseCfg.Redis, db)
		util.Log.Info("Enabled redis as database cache")
	}

	if err != nil {
		util.Log.Fatal("Failed connecting to database:", err)
	}
	util.Log.Info("Connected to database")

	return db
}

func printSqliteWraning() {
	util.Log.Warning("--------------------------[ ATTENTION ]--------------------------")
	util.Log.Warning("You are currently using SQLite as database driver. Please ONLY   ")
	util.Log.Warning("use SQLite during testing and debugging and NEVER use SQLite in a")
	util.Log.Warning("real production environment! Here you can read about why:        ")
	util.Log.Warning("https://github.com/zekroTJA/shinpuru/wiki/No-SQLite-in-production")
	util.Log.Warning("-----------------------------------------------------------------")
}
