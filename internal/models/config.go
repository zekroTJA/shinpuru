package models

import (
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/random"
)

var DefaultConfig = Config{
	Version: 6,
	Discord: Discord{
		GeneralPrefix: "sp!",
		GlobalCommandRateLimit: Ratelimit{
			Enabled:      true,
			Burst:        5,
			LimitSeconds: 3,
		},
	},
	Permissions: Permissions{
		DefaultUserRules:  static.DefaultUserRules,
		DefaultAdminRules: static.DefaultAdminRules,
	},
	Database: DatabaseType{
		Type:  "mysql",
		MySql: DatabaseCreds{},
	},
	Cache: Cache{
		Redis: CacheRedis{
			Addr:     "localhost:6379",
			Password: "",
			Type:     0,
		},
		CacheDatabase: true,
	},
	Logging: Logging{
		CommandLogging: true,
		LogLevel:       4,
	},
	TwitchApp: TwitchApp{},
	Storage: StorageType{
		Type: "file",
		File: StorageFile{
			Location: "./data",
		},
		Minio: StorageMinio{
			Location: "us-east-1",
			Secure:   true,
		},
	},
	WebServer: WebServer{
		Enabled:     true,
		Addr:        ":8080",
		APITokenKey: random.MustGetRandBase64Str(32),
		PublicAddr:  "https://example.com:8080",
		TLS: WebServerTLS{
			Enabled: false,
		},
		AccessToken: AccessToken{
			Secret:          random.MustGetRandBase64Str(64),
			LifetimeSeconds: 10 * 60,
		},
		LandingPage: LandingPage{
			ShowLocalInvite:   true,
			ShowPublicInvites: true,
		},
		RateLimit: Ratelimit{
			Enabled:      false,
			Burst:        30,
			LimitSeconds: 3,
		},
	},
	Metrics: Metrics{
		Enable: false,
		Addr:   ":9091",
	},
	Schedules: Schedules{
		GuildBackups:        "0 0 6,18 * * *",
		RefreshTokenCleanup: "0 0 5 * * *",
		ReportsExpiration:   "@every 5m",
		VerificationKick:    "@every 1h",
	},
	CodeExec: CodeExec{
		Type: "jdoodle",
		Ranna: CodeExecRanna{
			ApiVersion: "v1",
		},
		RateLimit: Ratelimit{
			Enabled:      true,
			Burst:        5,
			LimitSeconds: 60,
		},
	},
}

// Discord holds general configurations to connect
// to the Discord API application and using the
// OAuth2 workflow for web frontend authorization.
type Discord struct {
	Token                  string    `json:"token"`
	GeneralPrefix          string    `json:"generalprefix"`
	OwnerID                string    `json:"ownerid"`
	ClientID               string    `json:"clientid"`
	ClientSecret           string    `json:"clientsecret"`
	GuildBackupLoc         string    `json:"guildbackuploc"`
	GlobalCommandRateLimit Ratelimit `json:"globalcommandratelimit"`
	DisabledCommands       []string  `json:"disabledcommands"`
	Sharding               Sharding  `json:"sharding"`
	GuildsLimit            int       `json:"guildslimit"`
}

// Sharding holds configuration for guild event sharding.
type Sharding struct {
	AutoID bool `json:"autoid"`
	ID     int  `json:"id"`
	Pool   int  `json:"pool"`
	Total  int  `json:"total"`
}

// DatabaseCreds holds credentials to connect to
// a generic database.
type DatabaseCreds struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// CacheRedis holds credentials and settings
// to connect to a Redis instance.
type CacheRedis struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Type     int    `json:"type"`
	// Deprecated. Just here for downwards compatibility.
	Enable bool `json:"enable"`
}

// DatabaseType holds the preference for which
// database module to be used and the seperate
// "slots" for database configurations.
type DatabaseType struct {
	Type  string        `json:"type"`
	MySql DatabaseCreds `json:"mysql"`
	Redis CacheRedis    `json:"redis"`
}

// Cache holds the preferences for caching
// services.
type Cache struct {
	Redis         CacheRedis `json:"redis"`
	CacheDatabase bool       `json:"cachedatabase"`
}

// Loging holds configuration values for the
// main logger.
type Logging struct {
	CommandLogging bool `json:"commandlogging"`
	LogLevel       int  `json:"loglevel"`
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
	Enabled         bool         `json:"enabled"`
	Addr            string       `json:"addr"`
	TLS             WebServerTLS `json:"tls"`
	APITokenKey     string       `json:"apitokenkey"`
	PublicAddr      string       `json:"publicaddr"`
	LandingPage     LandingPage  `json:"landingpage"`
	DebugPublicAddr string       `json:"debugpublicaddr,omitempty"`
	RateLimit       Ratelimit    `json:"ratelimit"`
	Captcha         Captcha      `json:"captcha"`
	AccessToken     AccessToken  `json:"accesstoken"`
}

// AccessToken holds the secret and lifetime for
// JWT access token signature.
type AccessToken struct {
	Secret          string `json:"secret"`
	LifetimeSeconds int    `json:"lifetimeseconds"`
}

// WebServerTLS wraps preferences for the TLS
// configuration of the web server.
type WebServerTLS struct {
	Enabled bool   `json:"enabled"`
	Cert    string `json:"certfile"`
	Key     string `json:"keyfile"`
}

// Ratelimit wraps generic rate limit
// configuration.
type Ratelimit struct {
	Enabled      bool `json:"enabled"`
	Burst        int  `json:"burst"`
	LimitSeconds int  `json:"limitseconds"`
}

// LandingPage wraps the settings for the web
// interfaces landing page.
type LandingPage struct {
	ShowPublicInvites bool `json:"showpublicinvites"`
	ShowLocalInvite   bool `json:"showlocalinvite"`
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
	Type  string       `json:"type"`
	Minio StorageMinio `json:"minio"`
	File  StorageFile  `json:"file"`
}

// Metrics holds the settings for the prometheus
// metrics server.
type Metrics struct {
	Enable bool   `json:"enable"`
	Addr   string `json:"addr"`
}

// Schedules holds cron-like job schedule
// specifications for continuously running
// jobs.
type Schedules struct {
	GuildBackups        string `json:"guildbackups"`
	RefreshTokenCleanup string `json:"refreshtokencleanup"`
	ReportsExpiration   string `json:"reportsexpiration"`
	VerificationKick    string `json:"verificationkick"`
}

// CodeExec wraps configurations for the
// code execution API used.
type CodeExec struct {
	Type      string        `json:"type"`
	Ranna     CodeExecRanna `json:"ranna"`
	RateLimit Ratelimit     `json:"ratelimit"`
}

// CodeExecRanna holds configuration values
// for ranna as code execution engine.
type CodeExecRanna struct {
	Token      string `json:"token"`
	Endpoint   string `json:"endpoint"`
	ApiVersion string `json:"apiversion"`
}

// Captcha holds the configuration for a
// captcha verification.
type Captcha struct {
	SiteKey   string `json:"sitekey"`
	SecretKey string `json:"secretkey"`
}

// Contact holds contact information.
type Contact struct {
	Title string `json:"title"`
	Value string `json:"value"`
	URL   string `json:"url,omitempty"`
}

// Privacy holds privacy and contact
// information shown in shinpuru.
type Privacy struct {
	NoticeURL string    `json:"noticeurl"`
	Contact   []Contact `json:"contact"`
}

// Giphy holds credentials and configuration
// to connect to the Giphy.com API.
type Giphy struct {
	APIKey string `json:"apikey"`
}

// Config wraps the whole configuration structure
// including a version, which must not be changed
// by users to identify the integrity of config
// files over version updates.
type Config struct {
	Version     int          `json:"configVersionPleaseDoNotChange"`
	Discord     Discord      `json:"discord"`
	Permissions Permissions  `json:"permissions"`
	Database    DatabaseType `json:"database"`
	Cache       Cache        `json:"cache"`
	Logging     Logging      `json:"logging"`
	TwitchApp   TwitchApp    `json:"twitchapp"`
	Storage     StorageType  `json:"storage"`
	WebServer   WebServer    `json:"webserver"`
	Metrics     Metrics      `json:"metrics"`
	Schedules   Schedules    `json:"schedules"`
	CodeExec    CodeExec     `json:"codeexec"`
	Giphy       Giphy        `json:"giphy"`
	Privacy     Privacy      `json:"privacy"`
}
