package inits

import (
	"strings"

	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
)

func InitDatabase(databaseCfg *config.DatabaseType) database.Database {
	var db database.Database
	var err error

	switch strings.ToLower(databaseCfg.Type) {
	case "mysql", "mariadb":
		db = new(database.MySQL)
		err = db.Connect(databaseCfg.MySql)
	case "sqlite", "sqlite3":
		db = new(database.Sqlite)
		err = db.Connect(databaseCfg.Sqlite)
	}
	if err != nil {
		util.Log.Fatal("Failed connecting to database:", err)
	}
	util.Log.Info("Connected to database")

	return db
}
