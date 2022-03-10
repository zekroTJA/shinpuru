package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/zekroTJA/shinpuru/pkg/inline"
	"github.com/zekroTJA/shinpuru/pkg/random"
	. "github.com/zekroTJA/shinpuru/pkg/validators"
	"gopkg.in/yaml.v2"
)

//////////////////////////////////////////////////////////////////////
//
//   SETUP
//   -----
//   This CLI tool helps to create a pre-configured and ready to
//   deploy docker-compose.yml with simple question promts.
//
//////////////////////////////////////////////////////////////////////

const (
	version = "v1.0.0"
)

type DockerCompose struct {
	Version  string
	Volumes  map[string]any
	Services map[string]*DockerService
}

type DockerService struct {
	Image       string
	Command     []string          `yaml:",omitempty"`
	Ports       []string          `yaml:",omitempty"`
	Volumes     []string          `yaml:",omitempty"`
	Restart     string            `yaml:",omitempty"`
	Environment map[string]string `yaml:",omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
	Labels      map[string]string `yaml:",omitempty"`
}

var (
	dcConfig = &DockerCompose{
		Version: "3",
		Volumes: nil,
		Services: map[string]*DockerService{
			"traefik": {
				Image: "traefik:latest",
				Command: []string{
					"--providers.docker",
					"--providers.file.directory=/etc/traefik/dynamic_conf",
					"--entrypoints.http.address=:80",
					"--entrypoints.https.address=:443",
					"--providers.docker.exposedByDefault=false",
					// "--certificatesResolvers.le.acme.email=<your@email>", // fill in via cli
					"--certificatesResolvers.le.acme.storage=/etc/certstore/acme.json",
					"--certificatesResolvers.le.acme.httpChallenge.entryPoint=http",
				},
				Ports: []string{
					"80:80",
					"443:443",
				},
				Volumes: []string{
					"/var/run/docker.sock:/var/run/docker.sock",
					"./traefik/config:/etc/traefik/dynamic_conf",
					"/etc/cert:/etc/cert:ro",
				},
				Restart: "always",
			},
			"minio": {
				Image:   "minio/minio:latest",
				Volumes: []string{}, // fill in via cli
				Environment: map[string]string{
					// "MINIO_ACCESS_KEY":  "minio_access_key", // fill in via cli
					// "MINIO_SECRET_KEY":  "minio_secret_key", // fill in via cli
					"MINIO_REGION_NAME": "us-east-1",
				},
				Command: []string{"server", "/data"},
				Restart: "always",
			},
			"redis": {
				Image:   "redis:latest",
				Restart: "always",
			},
			"mysql": {
				Image: "mariadb:latest",
				Environment: map[string]string{
					// "MYSQL_ROOT_PASSWORD": "mysql_root_password" // finn in via cli
					"MYSQL_DATABASE": "shinpuru",
				},
				Volumes: []string{}, // fill in via cli
				Restart: "always",
			},
			"shinpuru": {
				// Image: "shinpuru:latest", // fill in via cli
				Volumes: []string{
					"./shinpuru/config:/etc/config",
					"/etc/cert:/etc/cert",
				},
				Environment: map[string]string{
					"SP_CONFIGVERSIONPLEASEDONOTCHANGE": "6",
					"SP_DISCORD_GENERALPREFIX":          "sp!",
					"SP_DATABASE_TYPE":                  "mysql",
					"SP_DATABASE_MYSQL_HOST":            "mysql",
					"SP_DATABASE_MYSQL_USER":            "root",
					"SP_DATABASE_MYSQL_DATABASE":        "shinpuru",
					"SP_CACHE_REDIS_ADDR":               "redis:6379",
					"SP_CACHE_REDIS_TYPE":               "0",
					"SP_CACHE_CACHEDATABASE":            "1",
					"SP_STORAGE_TYPE":                   "minio",
					"SP_STORAGE_MINIO_ENDPOINT":         "minio:9000",
					"SP_STORAGE_MINIO_LOCATION":         "us-east-1",
					"SP_STORAGE_MINIO_SECURE":           "0",
				}, // fill in via cli
				Restart: "always",
				DependsOn: []string{
					"mysql",
					"redis",
					"minio",
				},
				Labels: map[string]string{
					"traefik.enable": "true",
					"traefik.http.routers.shinpuru.entrypoints":      "https",
					"traefik.http.routers.shinpuru.tls":              "true",
					"traefik.http.routers.shinpuru.tls.certresolver": "le",
					// "traefik.http.routers.shinpuru.rule":             "Host(`<example.com>`)", // fill in via cli
				},
			},
		},
	}
)

func main() {
	prompts()
	hydrateCredentials()
	writeComposeFile()
}

func prompts() {
	fmt.Println(promptui.Styler(promptui.FGCyan)("shinpuru setup tool "), version)
	fmt.Println("\nℹ️ This tool sets up a pre-configured docker-compose.yml for hosting shinpuru.")

	spSvc := dcConfig.Services["shinpuru"]
	traefikSvc := dcConfig.Services["traefik"]

	// ----------------------
	// --- SHINPURU VERSION

	releaseChan, _ := mustVal2((&promptui.Select{
		Label: "First of all, which releace channel of shinpuru do you want to use?",
		Items: []string{
			"release",
			"canary",
		},
	}).Run())
	if releaseChan == 0 {
		spSvc.Image = "ghcr.io/zekrotja/shinpuru:latest"
	} else {
		spSvc.Image = "ghcr.io/zekrotja/shinpuru:canary"
	}

	// ----------------------
	// --- DISCORD CONFIG

	fmt.Println(promptui.Styler(promptui.FGMagenta)("\n1) Discord Credentials"))
	fmt.Println("Go to ",
		promptui.Styler(promptui.FGBlue)("https://discord.com/developers/applications"),
		" to create a Discord Bot application if not already done.")

	// token
	spSvc.Environment["SP_DISCORD_TOKEN"] = mustVal((&promptui.Prompt{
		Label:    "Bot Token",
		Mask:     '*',
		Validate: Length[string](50, 0),
	}).Run())

	// ownerid
	spSvc.Environment["SP_DISCORD_OWNERID"] = mustVal((&promptui.Prompt{
		Label:    "Owner ID (your Discord ID)",
		Validate: IsInteger(),
	}).Run())

	// Client ID
	spSvc.Environment["SP_DISCORD_CLIENTID"] = mustVal((&promptui.Prompt{
		Label:    "Bot Client ID",
		Validate: IsInteger(),
	}).Run())

	// Client Secret
	spSvc.Environment["SP_DISCORD_CLIENTSECRET"] = mustVal((&promptui.Prompt{
		Label:    "Bot Client Secret (not the Bot Token!)",
		Validate: Length[string](25, 0),
		Mask:     '*',
	}).Run())

	// Guilds Limit
	spSvc.Environment["SP_DISCORD_GUILDSLIMIT"] = mustVal((&promptui.Prompt{
		Label:    "Guilds Limit",
		Validate: InRange[string](0, 0),
		Default:  "0",
	}).Run())

	// Global Rate Limit
	globalRl, _ := mustVal2((&promptui.Select{
		Label: "Do you want to enable global rate limiting for slash commands?",
		Items: []string{
			"no",
			"yes",
		},
	}).Run())
	if globalRl == 1 {
		spSvc.Environment["SP_DISCORD_GLOBALCOMMANDRATELIMIT_ENABLED"] = "true"
		spSvc.Environment["SP_DISCORD_GLOBALCOMMANDRATELIMIT_BURST"] = mustVal((&promptui.Prompt{
			Label:    "Global Rate Limit: Burst",
			Validate: InRange[string](1, 0),
			Default:  "1",
		}).Run())
		spSvc.Environment["SP_DISCORD_GLOBALCOMMANDRATELIMIT_LIMITSECONDS"] = mustVal((&promptui.Prompt{
			Label:    "Global Rate Limit: Limit (Seconds)",
			Validate: InRange[string](1, 0),
			Default:  "10",
		}).Run())
	}

	// ----------------------
	// --- LOGGING

	fmt.Println(promptui.Styler(promptui.FGMagenta)("\n2) Logging"))

	// Command Logging
	commandLogging, _ := mustVal2((&promptui.Select{
		Label: "Do you want to enable command logging?",
		Items: []string{
			"no",
			"yes",
		},
	}).Run())
	spSvc.Environment["SP_LOGGING_COMMANDLOGGING"] = inline.II(commandLogging == 1, "true", "false")

	// Loglevel
	// Command Logging
	loglevel, _ := mustVal2((&promptui.Select{
		Label: "Do you want to enable command logging?",
		Items: []string{
			"panic",
			"fatal",
			"error",
			"warning",
			"info",
			"debug",
			"trace",
		},
	}).Run())
	spSvc.Environment["SP_LOGGING_LOGLEVEL"] = strconv.Itoa(loglevel)

	// ----------------------
	// --- WEBSERVER

	fmt.Println(promptui.Styler(promptui.FGMagenta)("\n3) Webserver"))
	fmt.Println("\n⚠️  When the web server is not enabled, a lot of features of shinpuru will not be available!")

	// Command Logging
	webserverEnable, _ := mustVal2((&promptui.Select{
		Label: "Do you want to enable the web server",
		Items: []string{
			"yes",
			"no",
		},
	}).Run())
	spSvc.Environment["SP_WEBSERVER_ENABLED"] = inline.II(webserverEnable == 0, "true", "false")

	// Domain
	publicDomain := mustVal((&promptui.Prompt{
		Label:    "Public Domain of the Server",
		Validate: IsDomain(),
	}).Run())
	spSvc.Environment["SP_WEBSERVER_PUBLICADDR"] = "https://" + publicDomain
	spSvc.Labels["traefik.http.routers.shinpuru.rule"] = fmt.Sprintf("Host(`%s`)", publicDomain)

	// E-Mail Address
	email := mustVal((&promptui.Prompt{
		Label:    "E-Mail Address (as contact for Lets Encrypt)",
		Validate: IsEmailAddress(),
	}).Run())
	traefikSvc.Command = append(traefikSvc.Command,
		"--certificatesResolvers.le.acme.email="+email)

	// Rate Limit
	wsRl, _ := mustVal2((&promptui.Select{
		Label: "Do you want to enable API rate limiting?",
		Items: []string{
			"no",
			"yes",
		},
	}).Run())
	if wsRl == 1 {
		spSvc.Environment["SP_WEBSERVER_RATELIMIT_ENABLED"] = "true"
		spSvc.Environment["SP_WEBSERVER_RATELIMIT_BURST"] = mustVal((&promptui.Prompt{
			Label:    "WS Rate Limit: Burst",
			Validate: InRange[string](1, 0),
			Default:  "50",
		}).Run())
		spSvc.Environment["SP_WEBSERVER_RATELIMIT_LIMITSECONDS"] = mustVal((&promptui.Prompt{
			Label:    "WS Rate Limit: Limit (Seconds)",
			Validate: InRange[string](1, 0),
			Default:  "3",
		}).Run())
	}

	// Captcha
	fmt.Println("\nℹ️ hCaptcha credentials can be obtained from",
		promptui.Styler(promptui.FGBlue)("https://dashboard.hcaptcha.com/overview"), ".")
	spSvc.Environment["SP_WEBSERVER_CAPTCHA_SITEKEY"] = mustVal((&promptui.Prompt{
		Label:   "hCaptcha Site Key",
		Default: "",
	}).Run())
	spSvc.Environment["SP_WEBSERVER_CAPTCHA_SECRETKEY"] = mustVal((&promptui.Prompt{
		Label:   "hCaptcha Secret Key",
		Default: "",
		Mask:    '*',
	}).Run())

	// ----------------------
	// --- TWITCH APP

	fmt.Println(promptui.Styler(promptui.FGMagenta)("\n4) Twitch"))

	fmt.Println("\nℹ️ Twitch API credentials can be obtained from",
		promptui.Styler(promptui.FGBlue)("https://glass.twitch.tv/console/apps"), ".")
	enableTwitch, _ := mustVal2((&promptui.Select{
		Label: "Do you want to enable Twitch integration?",
		Items: []string{
			"no",
			"yes",
		},
	}).Run())
	if enableTwitch == 1 {
		spSvc.Environment["SP_TWITCHAPP_CLIENTID"] = mustVal((&promptui.Prompt{
			Label:    "Twitch Client ID",
			Validate: Length[string](20, 0),
		}).Run())
		spSvc.Environment["SP_TWITCHAPP_CLIENTSECRET"] = mustVal((&promptui.Prompt{
			Label:    "Twitch Client Secret",
			Validate: Length[string](20, 0),
			Mask:     '*',
		}).Run())
	}

	// ----------------------
	// --- GIPHY

	fmt.Println(promptui.Styler(promptui.FGMagenta)("\n5) Giphy"))

	fmt.Println("\nℹ️ The Giphy API key can be obtained from",
		promptui.Styler(promptui.FGBlue)("https://developers.giphy.com/dashboard"), ".")
	fmt.Println("You can leave this empty if you don't need this right now.")
	spSvc.Environment["SP_GIPHY_APIKEY"] = mustVal((&promptui.Prompt{
		Label:    "Giphy API Key",
		Validate: Length[string](20, 0),
		Mask:     '*',
	}).Run())

	// ----------------------
	// --- CODEEXEC

	fmt.Println(promptui.Styler(promptui.FGMagenta)("\n6) Code Execution"))

	codeExecType, _ := mustVal2((&promptui.Select{
		Label: "Which type of code execution engine do you want to use?",
		Items: []string{
			"No code execution",
			"ranna (public)",
			"ranna (custom)",
			"jdoodle",
		},
	}).Run())
	switch codeExecType {
	case 1:
		spSvc.Environment["SP_CODEEXEC_TYPE"] = "ranna"
		spSvc.Environment["SP_CODEEXEC_RANNA_APIVERSION"] = "v1"
		spSvc.Environment["SP_CODEEXEC_RANNA_ENDPOINT"] = "https://public.ranna.dev"
	case 2:
		spSvc.Environment["SP_CODEEXEC_TYPE"] = "ranna"
		spSvc.Environment["SP_CODEEXEC_RANNA_ENDPOINT"] = mustVal((&promptui.Prompt{
			Label: "ranna API endpoint",
		}).Run())
		spSvc.Environment["SP_CODEEXEC_RANNA_APIVERSION"] = mustVal((&promptui.Prompt{
			Label:    "ranna API version",
			Validate: MatchesRegex(`^[vV]\d$`),
		}).Run())
		spSvc.Environment["SP_CODEEXEC_RANNA_TOKEN"] = mustVal((&promptui.Prompt{
			Label: "ranna API token",
			Mask:  '*',
		}).Run())
	case 3:
		spSvc.Environment["SP_CODEEXEC_TYPE"] = "jdoodle"
	}

	// Rate Limit
	codeExecRl, _ := mustVal2((&promptui.Select{
		Label: "Do you want to enable code exec rate limiting?",
		Items: []string{
			"yes",
			"no",
		},
	}).Run())
	if codeExecRl == 0 {
		spSvc.Environment["SP_CODEEXEC_RATELIMIT_ENABLED"] = "true"
		spSvc.Environment["SP_CODEEXEC_RATELIMIT_BURST"] = mustVal((&promptui.Prompt{
			Label:    "CE Rate Limit: Burst",
			Validate: InRange[string](1, 0),
			Default:  "5",
		}).Run())
		spSvc.Environment["SP_CODEEXEC_RATELIMIT_LIMITSECONDS"] = mustVal((&promptui.Prompt{
			Label:    "CE Rate Limit: Limit (Seconds)",
			Validate: InRange[string](1, 0),
			Default:  "60",
		}).Run())
	}

	// ----------------------
	// --- PRIVACY

	fmt.Println(promptui.Styler(promptui.FGMagenta)("\n7) Privacy"))

	privacyMode, _ := mustVal2((&promptui.Select{
		Label: "Do you want to run this instance of this bot publicly?",
		Items: []string{
			"Yes, I will run it publicly.",
			"No, just for me and my friends.",
		},
	}).Run())
	if privacyMode == 1 {
		spSvc.Environment["SP_PRIVACY_NOTICEURL"] = "https://github.com/zekroTJA/shinpuru/blob/master/PRIVACY.md"
		spSvc.Environment["SP_PRIVACY_CONTACT_0_TITLE"] = "private"
		spSvc.Environment["SP_PRIVACY_CONTACT_0_VALUE"] = "private"
	} else {
		fmt.Println("\nℹ️ When running shinpuru publicly, you must provide a privacy notice.")
		fmt.Println("Here you can find the privacy notice of the official instance of shinpuru:")
		fmt.Println(promptui.Styler(promptui.FGBlue)("https://github.com/zekroTJA/shinpuru/blob/master/PRIVACY.md"))
		fmt.Println(promptui.Styler(promptui.FGRed)("⚠️  You are not allowed to use this same notice for your instance!"))
		spSvc.Environment["SP_PRIVACY_NOTICEURL"] = mustVal((&promptui.Prompt{
			Label:    "Privacy Notice URL",
			Validate: IsSimpleUrl(),
		}).Run())
		spSvc.Environment["SP_PRIVACY_CONTACT_0_TITLE"] = "E-Mail"
		privacyMail := mustVal((&promptui.Prompt{
			Label:    "Contact E-Mail Address",
			Validate: IsEmailAddress(),
		}).Run())
		spSvc.Environment["SP_PRIVACY_CONTACT_0_VALUE"] = privacyMail
		spSvc.Environment["SP_PRIVACY_CONTACT_0_URL"] = "mailto:" + privacyMail
		fmt.Println("\nℹ️ You can add more privacy contact information later in your docker-compose.yml or shinpuru config.")
	}

	// ----------------------
	// --- DEPLOYMENT

	fmt.Println(promptui.Styler(promptui.FGMagenta)("\n8) Deployment"))
	fmt.Println("\nℹ️ Last of all, now some questions about your deployment preferences.")

	useVolumes, _ := mustVal2((&promptui.Select{
		Label: "Do you want to use local bindings or Docker volumes?",
		Items: []string{
			"Local Bindings (recommended)",
			"Docker Volumes",
		},
	}).Run())
	if useVolumes == 0 {
		dcConfig.Services["minio"].Volumes = []string{
			"./minio/data:/data",
		}
		dcConfig.Services["mysql"].Volumes = []string{
			"./mysql/cfg:/etc/mysql",
			"./mysql/lib:/var/lib/mysql",
		}
	} else {
		dcConfig.Volumes = map[string]any{
			"minio-data": struct{}{},
			"mysql-cfg":  struct{}{},
			"mysql-data": struct{}{},
		}
		dcConfig.Services["minio"].Volumes = []string{
			"minio-data:/data",
		}
		dcConfig.Services["mysql"].Volumes = []string{
			"mysql-cfg:/etc/mysql",
			"mysql-data:/var/lib/mysql",
		}
	}

	fmt.Println(" ")
}

func hydrateCredentials() {
	fmt.Println("🔑 ", promptui.Styler(promptui.FGGreen)("Generating random credentials for services ..."))

	minioSvc := dcConfig.Services["minio"]
	minioSvc.Environment["MINIO_ACCESS_KEY"] = mustGetRandomKey(32)
	minioSvc.Environment["MINIO_SECRET_KEY"] = mustGetRandomKey(32)

	msqlSvc := dcConfig.Services["mysql"]
	msqlSvc.Environment["MYSQL_ROOT_PASSWORD"] = mustGetRandomKey(32)

	spSvc := dcConfig.Services["shinpuru"]
	spSvc.Environment["SP_WEBSERVER_APITOKENKEY"] = mustGetRandomKey(64)
	spSvc.Environment["SP_WEBSERVER_ACCESSTOKEN_SECRET"] = mustGetRandomKey(64)
	spSvc.Environment["SP_STORAGE_MINIO_ACCESSKEY"] = minioSvc.Environment["MINIO_ACCESS_KEY"]
	spSvc.Environment["SP_STORAGE_MINIO_ACCESSSECRET"] = minioSvc.Environment["MINIO_SECRET_KEY"]
	spSvc.Environment["SP_DATABASE_MYSQL_PASSWORD"] = msqlSvc.Environment["MYSQL_ROOT_PASSWORD"]
}

func writeComposeFile() {
	fmt.Println("📃 ", promptui.Styler(promptui.FGYellow)("Creating docker-compose.yml file ..."))

	f := mustVal(os.Create("docker-compose.yml"))
	defer f.Close()
	must(yaml.NewEncoder(f).Encode(dcConfig))

	fmt.Println("✔️  ", promptui.Styler(promptui.FGGreen)("docker-compose.yml successfully set up!"))
}

func must(err error) {
	if err != nil {
		fmt.Println(
			promptui.Styler(promptui.BGRed, promptui.FGBlack)(" ERROR "),
			promptui.Styler(promptui.FGRed)(err.Error()))
		os.Exit(1)
	}
}

func mustVal[T any](v T, err error) T {
	must(err)
	return v
}

func mustVal2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	must(err)
	return v1, v2
}

func mustGetRandomKey(len int) (v string) {
	v, err := random.GetRandBase64Str(len)
	must(err)
	return
}
