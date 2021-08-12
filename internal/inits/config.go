package inits

import (
	"os"

	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

func InitConfig(configLocation string, container di.Container) *config.Config {
	defaultConfig := config.GetDefaultConfig()

	cfgParser := container.Get(static.DiConfigParser).(config.Parser)

	cfgFile, err := os.Open(configLocation)
	defer cfgFile.Close()
	if os.IsNotExist(err) {
		cfgFile, err = os.Create(configLocation)
		if err != nil {
			logrus.WithError(err).Fatal("Config file was not found and failed creating default config")
		}
		err = cfgParser.Encode(cfgFile, defaultConfig)
		if err != nil {
			logrus.WithError(err).Fatal("Config file was not found and failed writing to new config file")
		}
		logrus.Fatal("Config file was not found. Created default config file. Please open it and enter your configuration.")
	} else if err != nil {
		logrus.WithError(err).Fatal("Failed opening config file")
	}

	cfg, err := cfgParser.Decode(cfgFile)
	if err != nil {
		logrus.WithError(err).Fatal("Failed decoding config file")
	}

	if cfg.Version < static.ConfigVersion {
		logrus.Fatal("Config file structure is outdated and must be re-created. Just rename your config and start the bot to recreate the latest valid version of the config.")
	}

	if cfg.Discord.OwnerID == "" {
		logrus.Warn("Bot onwer ID is not set in config!",
			"You will not be identified as the owner of this bot so you will not have access to the owner-only commands!")
	}

	if cfg.Cache == nil {
		cfg.Cache = &config.Cache{}
	}
	if cfg.Cache.Redis == nil {
		if cfg.Database.Redis == nil {
			logrus.Fatal("A redis connection must be defined! (.cache.redis)")
		}
		logrus.Warn("The '.database.redis' config key is deprecated and is moved to '.cache.redis'. " +
			"Please move your existing redis configuration to '.cache.redis'. " +
			"Support for '.database.redis' will be removed in an upcoming version!")
		cfg.Cache.Redis = cfg.Database.Redis
	}

	logrus.Info("Config file loaded")

	cfg.Defaults = defaultConfig
	return cfg
}
