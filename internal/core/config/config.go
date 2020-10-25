package config

import (
	"io"

	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/random"
)

// Discord holds general configurations to connect
// to the Discord API application and using the
// OAuth2 workflow for web frontend authorization.
type Discord struct {
	Token          string
	GeneralPrefix  string
	OwnerID        string
	ClientID       string
	ClientSecret   string
	GuildBackupLoc string
}

// DatabaseCreds holds credentials to connect to
// a generic database.
type DatabaseCreds struct {
	Host     string
	User     string
	Password string
	Database string
}

// DatabaseFile holds information to use a file
// database like SQLite.
type DatabaseFile struct {
	DBFile string
}

// DatabaseRedis holds credentials and settings
// to connect to a Redis database.
type DatabaseRedis struct {
	Enable   bool
	Addr     string
	Password string
	Type     int
}

// DatabaseType holds the preference for which
// database module to be used and the seperate
// "slots" for database configurations.
type DatabaseType struct {
	Type   string
	MySql  *DatabaseCreds
	Sqlite *DatabaseFile
	Redis  *DatabaseRedis
}

// Loging holds configuration values for the
// main logger.
type Logging struct {
	CommandLogging bool
	LogLevel       int
}

// TwitchApp holds credentials to connect to
// a Twitch API application.
type TwitchApp struct {
	ClientID     string `json:"clientid"`
	ClientSecret string `json:"clientsecret"`
}

// WebServer holds general configurations for
// the exposed web server.
type WebServer struct {
	Enabled         bool          `json:"enabled"`
	Addr            string        `json:"addr"`
	TLS             *WebServerTLS `json:"tls"`
	APITokenKey     string        `json:"apitokenkey"`
	PublicAddr      string        `json:"publicaddr"`
	DebugPublicAddr string        `json:"debugpublicaddr,omitempty"`
}

// WebServerTLS wraps preferences for the TLS
// configuration of the web server.
type WebServerTLS struct {
	Enabled bool   `json:"enabled"`
	Cert    string `json:"certfile"`
	Key     string `json:"keyfile"`
}

// Permissions wrap standard rulesets for specific
// user groups like guild admins and default users
// with no special previleges.
type Permissions struct {
	DefaultUserRules  []string `json:"defaultuserrules"`
	DefaultAdminRules []string `json:"defaultadminrules"`
}

// StorageMinio holds connection preferences to
// conenct to a storage provider like MinIO,
// Amazon S3 or Google Cloud.
type StorageMinio struct {
	Endpoint     string `json:"endpoint"`
	AccessKey    string `json:"accesskey"`
	AccessSecret string `json:"accesssecret"`
	Location     string `json:"location"`
	Secure       bool   `json:"secure"`
}

// StorageFile holds preferences for a local
// file storage provider.
type StorageFile struct {
	Location string `json:"location"`
}

// StorageType holds the preferences for which
// storage type is to be used and "slots" for
// the specific configuration of them.
type StorageType struct {
	Type  string        `json:"type"`
	Minio *StorageMinio `json:"minio"`
	File  *StorageFile  `json:"file"`
}

// Metrics holds the settings for the prometheus
// metrics server.
type Metrics struct {
	Enable bool   `json:"enable"`
	Addr   string `json:"addr"`
}

// Config wraps the whole configuration structure
// including a version, which must not be changed
// by users to identify the integrity of config
// files over version updates.
type Config struct {
	Version     int `yaml:"configVersionPleaseDoNotChange"`
	Discord     *Discord
	Permissions *Permissions
	Database    *DatabaseType
	Logging     *Logging
	TwitchApp   *TwitchApp
	Storage     *StorageType
	WebServer   *WebServer
	Metrics     *Metrics
}

// Parser describes a general configuration parser
// to decode and encode a Config from or to file.
type Parser interface {
	// Decode deserializes a Config instance
	// from the passed data reader and returns
	// the Config instance and errors occured
	// during deserialization.
	Decode(r io.Reader) (*Config, error)
	// Encode serializes a data stream to the
	// passed stream writer from the passed
	// Config instance and returns errors
	// during serialization.
	Encode(w io.Writer, c *Config) error
}

// GetDefaultConfig returns a Config instance with
// default values.
func GetDefaultConfig() *Config {
	apiTokenKey, _ := random.GetRandBase64Str(32)

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
		TwitchApp: &TwitchApp{},
		Storage: &StorageType{
			Type: "file",
			File: &StorageFile{
				Location: "./data",
			},
			Minio: &StorageMinio{
				Location: "us-east-1",
				Secure:   true,
			},
		},
		WebServer: &WebServer{
			Enabled:     true,
			Addr:        ":8080",
			APITokenKey: apiTokenKey,
			PublicAddr:  "https://example.com:8080",
			TLS: &WebServerTLS{
				Enabled: false,
			},
		},
		Metrics: &Metrics{
			Enable: false,
			Addr:   ":9091",
		},
	}
}
