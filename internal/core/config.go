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

type ConfigLogging struct {
	CommandLogging bool
	LogLevel       int
}

type ConfigEtc struct {
	TwitchAppID string
}

type Config struct {
	Version  int `yaml:"configVersionPleaseDoNotChange"`
	Discord  *ConfigDiscord
	Database *ConfigDatabaseType
	Logging  *ConfigLogging
	Etc      *ConfigEtc
}

type ConfigParser interface {
	Decode(r io.Reader) (*Config, error)
	Encode(w io.Writer, c *Config) error
}

func NewDefaultConfig() *Config {
	return &Config{
		Version: 4,
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
		Logging: &ConfigLogging{
			CommandLogging: true,
			LogLevel:       4,
		},
		Etc: new(ConfigEtc),
	}
}
