package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/inits"

	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/util"
)

var (
	flagConfigLocation = flag.String("c", "config.yml", "The location of the main config file")
	flagDocker         = flag.Bool("docker", false, "wether shinpuru is running in a docker container or not")
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

	conf := inits.InitConfig(*flagConfigLocation, new(config.YAMLConfigParser))

	if *flagDocker {
		if conf.Database.Sqlite == nil {
			conf.Database.Sqlite = new(config.DatabaseFile)
		}
		conf.Database.Sqlite.DBFile = "/etc/db/db.sqlite3"
		conf.WebServer.Addr = ":8080"
	}

	util.SetLogLevel(conf.Logging.LogLevel)

	database := inits.InitDatabase(conf.Database)
	defer func() {
		util.Log.Info("Shutting down database connection...")
		database.Close()
	}()

	tnw := inits.InitTwitchNotifyer(session, conf, database)

	lct := inits.InitLTCTimer()

	cmdHandler := inits.InitCommandHandler(session, conf, database, tnw, lct)
	inits.InitDiscordBotSession(session, conf, database, cmdHandler, lct)
	defer func() {
		util.Log.Info("Shutting down bot session...")
		session.Close()
	}()

	inits.InitWebServer(session, database, cmdHandler, conf)

	util.Log.Info("Started event loop. Stop with CTRL-C...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
