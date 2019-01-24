package inits

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/listeners"
	"github.com/zekroTJA/shinpuru/internal/util"
)

func InitConfig(configLocation string, cfgParser core.ConfigParser) *core.Config {
	cfgFile, err := os.Open(configLocation)
	if os.IsNotExist(err) {
		cfgFile, err = os.Create(configLocation)
		if err != nil {
			util.Log.Fatal("Config file was not found and failed creating default config:", err)
		}
		err = cfgParser.Encode(cfgFile, core.NewDefaultConfig())
		if err != nil {
			util.Log.Fatal("Config file was not found and failed writing to new config file:", err)
		}
		util.Log.Fatal("Config file was not found. Created default config file. Please open it and enter your configuration.")
	} else if err != nil {
		util.Log.Fatal("Failed opening config file:", err)
	}

	config, err := cfgParser.Decode(cfgFile)
	if err != nil {
		util.Log.Fatal("Failed decoding config file:", err)
	}

	if config.Version < util.ConfigVersion {
		util.Log.Fatalf("Config file structure is outdated and must be re-created. Just rename your config and start the bot to recreate the latest valid version of the config.")
	}

	if config.Discord.OwnerID == "" {
		util.Log.Warning("ATTENTION: Bot onwer ID is not set in config!",
			"You will not be identified as the owner of this bot so you will not have access to the owner-only commands!")
	}

	util.Log.Info("Config file loaded")

	return config
}

func InitDatabase(databaseCfg *core.ConfigDatabaseType) core.Database {
	var database core.Database
	var err error

	switch strings.ToLower(databaseCfg.Type) {
	case "mysql", "mariadb":
		database = new(core.MySql)
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

func InitCommandHandler(config *core.Config, database core.Database) *commands.CmdHandler {
	cmdHandler := commands.NewCmdHandler(database, config)

	cmdHandler.RegisterCommand(&commands.CmdHelp{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdPrefix{PermLvl: 10})
	cmdHandler.RegisterCommand(&commands.CmdPerms{PermLvl: 10})
	cmdHandler.RegisterCommand(&commands.CmdClear{PermLvl: 8})
	cmdHandler.RegisterCommand(&commands.CmdMvall{PermLvl: 5})
	cmdHandler.RegisterCommand(&commands.CmdInfo{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdSay{PermLvl: 3})
	cmdHandler.RegisterCommand(&commands.CmdQuote{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdGame{PermLvl: 999})
	cmdHandler.RegisterCommand(&commands.CmdAutorole{PermLvl: 9})
	cmdHandler.RegisterCommand(&commands.CmdReport{PermLvl: 5})
	cmdHandler.RegisterCommand(&commands.CmdModlog{PermLvl: 6})
	cmdHandler.RegisterCommand(&commands.CmdKick{PermLvl: 6})
	cmdHandler.RegisterCommand(&commands.CmdBan{PermLvl: 8})
	cmdHandler.RegisterCommand(&commands.CmdVote{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdProfile{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdId{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdMute{PermLvl: 4})
	cmdHandler.RegisterCommand(&commands.CmdMention{PermLvl: 4})
	cmdHandler.RegisterCommand(&commands.CmdNotify{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdVoicelog{PermLvl: 6})
	cmdHandler.RegisterCommand(&commands.CmdBug{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdStats{PermLvl: 0})

	if util.Release != "TRUE" {
		cmdHandler.RegisterCommand(&commands.CmdTest{})
	}

	if config.Permissions != nil {
		cmdHandler.UpdateCommandPermissions(config.Permissions.CustomCmdPermissions)
		if config.Permissions.BotOwnerLevel > 0 {
			util.PermLvlBotOwner = config.Permissions.BotOwnerLevel
		}
		if config.Permissions.GuildOwnerLevel > 0 {
			util.PermLvlGuildOwner = config.Permissions.GuildOwnerLevel
		}
	}

	util.Log.Infof("%d commands registered", cmdHandler.GetCommandListLen())

	return cmdHandler
}

func InitDiscordBotSession(config *core.Config, database core.Database, cmdHandler *commands.CmdHandler) {
	snowflake.Epoch = util.DefEpoche
	err := util.SetupSnowflakeNodes()
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
	session.AddHandler(listeners.NewListenerChannelCreate(database).Handler)
	session.AddHandler(listeners.NewListenerVoiceUpdate(database).Handler)

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
