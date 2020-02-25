package main

import (
	"flag"

	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"

	"github.com/zekroTJA/shinpuru/internal/inits"
)

var (
	flagExportFile = flag.String("o", "commandsManual.md", "output location of manual file")
)

func main() {
	flag.Parse()

	config := &config.Config{
		Discord: &config.Discord{},
	}

	database := new(database.MySQL)

	cmdHandler := inits.InitCommandHandler(nil, config, database, nil, nil)
	if err := cmdHandler.ExportCommandManual(*flagExportFile); err != nil {
		util.Log.Fatal("Failed exporting command manual: ", err)
	}
	util.Log.Info("Successfully exported command manual file to " + *flagExportFile)
}
