package inits

import (
	"strings"

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

func InitDatabase(databaseCfg *core.ConfigDatabaseType) core.Database {
	var database core.Database
	var err error

	switch strings.ToLower(databaseCfg.Type) {
	case "mysql", "mariadb":
		database = new(core.MySQL)
		err = database.Connect(databaseCfg.MySql)
	case "sqlite", "sqlite3":
		database = new(core.Sqlite)
		err = database.Connect(databaseCfg.Sqlite)
	}
	if err != nil {
		util.Log.Fatal("Failed connecting to database:", err)
	}
	util.Log.Info("Connected to database")

	return database
}
