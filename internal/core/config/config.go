package config

import (
	"io"

	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type Discord struct {
	Token          string
	GeneralPrefix  string
	OwnerID        string
	ClientID       string
	ClientSecret   string
	GuildBackupLoc string
}

type DatabaseCreds struct {
	Host     string
	User     string
	Password string
	Database string
}

type DatabaseFile struct {
	DBFile string
}

type DatabaseRedis struct {
	Enable   bool
	Addr     string
	Password string
	Type     int
}

type DatabaseType struct {
	Type   string
	MySql  *DatabaseCreds
	Sqlite *DatabaseFile
	Redis  *DatabaseRedis
}

type Logging struct {
	CommandLogging bool
	LogLevel       int
}

type Etc struct {
	TwitchAppID string
}

type WS struct {
	Enabled         bool   `json:"enabled"`
	Addr            string `json:"addr"`
	TLS             *WSTLS `json:"tls"`
	PublicAddr      string `json:"publicaddr"`
	DebugPublicAddr string `json:"debugpublicaddr,omitempty"`
}

type WSTLS struct {
	Enabled bool   `json:"enabled"`
	Cert    string `json:"certfile"`
	Key     string `json:"keyfile"`
}

type Permissions struct {
	DefaultUserRules  []string `json:"defaultuserrules"`
	DefaultAdminRules []string `json:"defaultadminrules"`
}

type StorageMinio struct {
	Endpoint     string `json:"endpoint"`
	AccessKey    string `json:"accesskey"`
	AccessSecret string `json:"accesssecret"`
	Secure       bool   `json:"secure"`
}

type StorageFile struct {
	Location string `json:"location"`
}

type StorageType struct {
	Type  string        `json:"type"`
	Minio *StorageMinio `json:"minio"`
	File  *StorageFile  `json:"file"`
}

type Config struct {
	Version     int `yaml:"configVersionPleaseDoNotChange"`
	Discord     *Discord
	Permissions *Permissions
	Database    *DatabaseType
	Logging     *Logging
	Etc         *Etc
	Storage     *StorageType
	WebServer   *WS
}

type Parser interface {
	Decode(r io.Reader) (*Config, error)
	Encode(w io.Writer, c *Config) error
}

func NewDefaultConfig() *Config {
	return &Config{
		Version: 6,
		Discord: &Discord{
			GeneralPrefix: "sp!",
		},
		Permissions: &Permissions{
			DefaultUserRules:  static.DefaultUserRules,
			DefaultAdminRules: static.DefaultAdminRules,
		},
		Database: &DatabaseType{
			Type:  "sqlite",
			MySql: &DatabaseCreds{},
			Sqlite: &DatabaseFile{
				DBFile: "shinpuru.sqlite3.db",
			},
			Redis: &DatabaseRedis{
				Enable:   false,
				Addr:     "localhost:6379",
				Password: "",
				Type:     0,
			},
		},
		Logging: &Logging{
			CommandLogging: true,
			LogLevel:       4,
		},
		Etc: &Etc{},
		Storage: &StorageType{
			Type: "file",
			File: &StorageFile{
				Location: "./data",
			},
			Minio: &StorageMinio{},
		},
		WebServer: &WS{
			Enabled:    true,
			Addr:       ":8080",
			PublicAddr: "https://example.com:8080",
			TLS: &WSTLS{
				Enabled: false,
			},
		},
	}
}
