package inits

import (
	"os"

	"github.com/zekroTJA/shinpuru/internal/core"
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
