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

type Config struct {
	Discord  *ConfigDiscord
	Database *ConfigDatabase
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
	}
}
