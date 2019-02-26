package main

import (
	"flag"

	"github.com/zekroTJA/shinpuru/internal/util"

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/inits"
)

var (
	flagExportFile = flag.String("o", "commandsManual.md", "output location of manual file")
)

func main() {
	flag.Parse()

	config := new(core.Config)
	database := new(core.MySql)
	cmdHandler := inits.InitCommandHandler(nil, config, database, nil)
	if err := cmdHandler.ExportCommandManual(*flagExportFile); err != nil {
		util.Log.Fatal("Failed exporting command manual: ", err)
	}
	util.Log.Info("Successfully exported command manual file to " + *flagExportFile)
}
