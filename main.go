package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/snowflake"

	"github.com/zekroTJA/shinpuru/commands"
	"github.com/zekroTJA/shinpuru/core"
	"github.com/zekroTJA/shinpuru/listeners"
	"github.com/zekroTJA/shinpuru/util"

	"github.com/bwmarrin/discordgo"
)

var (
	configLocation = flag.String("c", "config.yml", "The location of the main config file")

	ldAppVersion = "TESTBUILD"
	ldAppCommit  = "TESTBUILD"
)

func main() {
	util.AppVersion = ldAppVersion
	util.AppCommit = ldAppCommit

	flag.Parse()
	util.Log.Infof("シンプル (shinpuru) v.%s (commit %s)", util.AppVersion, util.AppCommit)
	util.Log.Info("© zekro Development (Ringo Hoffmann)")
	util.Log.Info("Covered by MIT Licence")
	util.Log.Info("Starting up...")

	cfgParser := new(core.YAMLConfigParser)

	///////////////////////////
	// CONFIG INITIALIZATION //
	///////////////////////////

	cfgFile, err := os.Open(*configLocation)
	if os.IsNotExist(err) {
		cfgFile, err = os.Create(*configLocation)
		if err != nil {
			log.Fatal("Config file was not found and failed creating default config:", err)
		}
		err = cfgParser.Encode(cfgFile, core.NewDefaultConfig())
		if err != nil {
			log.Fatal("Config file was not found and failed writing to new config file:", err)
		}
		log.Fatal("Config file was not found. Created default config file. Please open it and enter your configuration.")
	} else if err != nil {
		log.Fatal("Failed opening config file:", err)
	}

	config, err := cfgParser.Decode(cfgFile)
	if err != nil {
		util.Log.Fatal("Failed decoding config file:", err)
	}

	if config.Discord.OwnerID == "" {
		util.Log.Warning("ATTENTION: Bot onwer ID is not set in config!",
			"You will not be identified as the owner of this bot so you will not have access to the owner-only commands!")
	}

	util.Log.Info("Config file loaded")

	////////////////////
	// DATABASE LOGIN //
	////////////////////

	database := new(core.MySql)
	if err := database.Connect(config.Database); err != nil {
		util.Log.Fatal("Failed connecting to database:", err)
	}
	util.Log.Info("Connected to database")

	//////////////////////////
	// COMMAND REGISTRATION //
	//////////////////////////

	cmdHandler := commands.NewCmdHandler(database, config)
	cmdHandler.RegisterCommand(new(commands.CmdTest))
	cmdHandler.RegisterCommand(new(commands.CmdHelp))
	cmdHandler.RegisterCommand(new(commands.CmdPrefix))
	cmdHandler.RegisterCommand(new(commands.CmdPerms))
	cmdHandler.RegisterCommand(new(commands.CmdClear))
	cmdHandler.RegisterCommand(new(commands.CmdMvall))
	cmdHandler.RegisterCommand(new(commands.CmdInfo))
	cmdHandler.RegisterCommand(new(commands.CmdSay))
	cmdHandler.RegisterCommand(new(commands.CmdQuote))
	cmdHandler.RegisterCommand(new(commands.CmdGame))
	cmdHandler.RegisterCommand(new(commands.CmdAutorole))
	cmdHandler.RegisterCommand(new(commands.CmdReport))
	cmdHandler.RegisterCommand(new(commands.CmdModlog))
	cmdHandler.RegisterCommand(new(commands.CmdKick))
	cmdHandler.RegisterCommand(new(commands.CmdBan))
	cmdHandler.RegisterCommand(new(commands.CmdVote))
	cmdHandler.RegisterCommand(new(commands.CmdProfile))
	cmdHandler.RegisterCommand(new(commands.CmdId))

	//////////////////////////
	// BOT SESSION CREATION //
	//////////////////////////

	snowflake.Epoch = util.DefEpoche
	err = util.SetupSnowflakeNodes()
	if err != nil {
		util.Log.Fatal("Failed setting up snowflake nodes: ", err)
	}

	session, err := discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		util.Log.Fatal("Failed creating Discord bot session:", err)
	}

	session.AddHandler(listeners.NewListenerReady(config, database).Handler)
	session.AddHandler(listeners.NewListenerCmd(config, database, cmdHandler).Handler)
	session.AddHandler(listeners.NewListenerGuildJoin(config).Handler)
	session.AddHandler(listeners.NewListenerMemberAdd(database).Handler)
	session.AddHandler(listeners.NewListenerVote(database).Handler)

	err = session.Open()
	if err != nil {
		util.Log.Fatal("Failed connecting Discord bot session:", err)
	}

	util.Log.Info("Started event loop. Stop with CTRL-C...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	util.Log.Info("Shutting down...")
	session.Close()
}
