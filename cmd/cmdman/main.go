package main

import (
	"flag"

	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/inits"
	"github.com/zekroTJA/shinpuru/internal/util"
)

var (
	flagExportFile = flag.String("o", "commandsManual.md", "output location of manual file")
)

func main() {
	flag.Parse()

	// Setting Release flag to true manually to prevent
	// registration of test command and exclude it in the
	// command manual.
	util.Release = "TRUE"

	config := &config.Config{
		Discord: &config.Discord{},
	}

	database := new(middleware.SqliteMiddleware)

	cmdHandler := inits.InitCommandHandler(nil, config, database, nil, nil, nil)
	if err := cmdHandler.ExportCommandManual(*flagExportFile); err != nil {
		util.Log.Fatal("Failed exporting command manual: ", err)
	}
	util.Log.Info("Successfully exported command manual file to " + *flagExportFile)
}
