package main

import (
	"flag"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/inits"

	"github.com/zekroTJA/shinpuru/internal/util"
)

var (
	flagConfigLocation = flag.String("c", "config.yml", "The location of the main config file")
)

func main() {
	flag.Parse()

	util.Log.Infof("シンプル (shinpuru) v.%s (commit %s)", util.AppVersion, util.AppCommit)
	util.Log.Info("© zekro Development (Ringo Hoffmann)")
	util.Log.Info("Covered by MIT Licence")
	util.Log.Info("Starting up...")

	session, err := discordgo.New("")
	if err != nil {
		util.Log.Fatal(err)
	}

	config := inits.InitConfig(*flagConfigLocation, new(core.YAMLConfigParser))

	util.SetLogLevel(config.Logging.LogLevel)

	database := inits.InitDatabase(config.Database)
	tnw := inits.InitTwitchNotifyer(session, config, database)
	cmdHandler := inits.InitCommandHandler(session, config, database, tnw)
	inits.InitDiscordBotSession(session, config, database, cmdHandler)
}
