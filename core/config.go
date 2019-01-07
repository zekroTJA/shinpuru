package core

import "io"

type ConfigDiscord struct {
	Token         string
	GeneralPrefix string
	OwnerID       string
}

type ConfigDatabase struct {
	Host     string
	User     string
	Password string
	Database string
}

type ConfigPermissions struct {
	BotOwnerLevel        int
	GuildOwnerLevel      int
	CustomCmdPermissions map[string]int
}

type Config struct {
	Discord        *ConfigDiscord
	Database       *ConfigDatabase
	Permissions    *ConfigPermissions
	CommandLogging bool
}

type ConfigParser interface {
	Decode(r io.Reader) (*Config, error)
	Encode(w io.Writer, c *Config) error
}

func NewDefaultConfig() *Config {
	return &Config{
		Discord: &ConfigDiscord{
			Token:         "",
			GeneralPrefix: "sp!",
			OwnerID:       "",
		},
		Database: new(ConfigDatabase),
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
