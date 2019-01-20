package core

import "io"

type ConfigDiscord struct {
	Token         string
	GeneralPrefix string
	OwnerID       string
}

type ConfigDatabaseCreds struct {
	Host     string
	User     string
	Password string
	Database string
}

type ConfigDatabaseFile struct {
	DBFile string
}

type ConfigDatabaseType struct {
	Type   string
	MySql  *ConfigDatabaseCreds
	Sqlite *ConfigDatabaseFile
}

type ConfigPermissions struct {
	BotOwnerLevel        int
	GuildOwnerLevel      int
	CustomCmdPermissions map[string]int
}

type Config struct {
	Version        int `yaml:"configVersionPleaseDoNotChange"`
	Discord        *ConfigDiscord
	Database       *ConfigDatabaseType
	Permissions    *ConfigPermissions
	CommandLogging bool
}

type ConfigParser interface {
	Decode(r io.Reader) (*Config, error)
	Encode(w io.Writer, c *Config) error
}

func NewDefaultConfig() *Config {
	return &Config{
		Version: 3,
		Discord: &ConfigDiscord{
			Token:         "",
			GeneralPrefix: "sp!",
			OwnerID:       "",
		},
		Database: &ConfigDatabaseType{
			Type:  "sqlite",
			MySql: new(ConfigDatabaseCreds),
			Sqlite: &ConfigDatabaseFile{
				DBFile: "shinpuru.sqlite3.db",
			},
		},
		Permissions: &ConfigPermissions{
			BotOwnerLevel:   1000,
			GuildOwnerLevel: 10,
			CustomCmdPermissions: map[string]int{
				"cmdinvoke": 0,
			},
		},
		CommandLogging: true,
	}
}
