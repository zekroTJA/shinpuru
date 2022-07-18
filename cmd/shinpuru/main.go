package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekrotja/ken"

	"github.com/zekroTJA/shinpuru/internal/inits"
	"github.com/zekroTJA/shinpuru/internal/listeners"
	"github.com/zekroTJA/shinpuru/internal/services/backup"
	"github.com/zekroTJA/shinpuru/internal/services/birthday"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/karma"
	"github.com/zekroTJA/shinpuru/internal/services/kvcache"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/services/verification"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/auth"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/embedded"
	"github.com/zekroTJA/shinpuru/internal/util/startupmsg"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/argp"
	"github.com/zekroTJA/shinpuru/pkg/onetimeauth/v2"
	"github.com/zekroTJA/shinpuru/pkg/startuptime"

	"github.com/zekroTJA/shinpuru/pkg/angularservice"
)

const (
	envKeyProfile = "CPUPROFILE"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	snowflake.Epoch = static.DefEpoche
}

//////////////////////////////////////////////////////////////////////
//
//   SHINPURU
//   --------
//   This is the main initialization for shinpuru which initializes
//   all instances like the database middleware, the twitch notify
//   listener service, life cycle timer, storage middleware,
//   permission middleware, command handler and - finally -
//   initializes the discord session event loop.
//   shinpuru is configured via a configuration file which location
//   can be passed via the '-c' parameter.
//   When shinpuru is run in a Docker container, the '-docker' flag
//   should be passed to fix configuration values like the location
//   of the sqlite3 database (when the sqlite3 driver is used) or
//   the web server exposure port.
//
//////////////////////////////////////////////////////////////////////

func main() {
	// Parse command line flags
	flagConfig, _ := argp.String("-c", "config.yml", "Optional config file location.")
	flagDevMode, _ := argp.Bool("-devmode", false, "Enable development mode.")
	flagProfile, _ := argp.String("-cpuprofile", "", "CPU profile output location.")
	flagQuiet, _ := argp.Bool("-quiet", false, "Hide startup message.")
	_, _ = argp.Bool("-docker", false, "Docker mode (deprecated)")

	if flagHelp, _ := argp.Bool("-h", false, "Display help."); flagHelp {
		fmt.Println("Usage:\n" + argp.Help())
		return
	}

	if !flagQuiet {
		startupmsg.Output(os.Stdout)
	}

	// Initialize dependency injection builder
	diBuilder, _ := di.NewBuilder()

	// Initialize time provider
	diBuilder.Add(di.Def{
		Name: static.DiTimeProvider,
		Build: func(ctn di.Container) (interface{}, error) {
			return timeprovider.Time{}, nil
		},
	})

	// Initialize config
	diBuilder.Add(di.Def{
		Name: static.DiConfig,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.NewPaerser(argp.Args(), flagConfig), nil
		},
	})

	// Initialize metrics server
	diBuilder.Add(di.Def{
		Name: static.DiMetrics,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitMetrics(ctn), nil
		},
	})

	// Initialize redis client
	diBuilder.Add(di.Def{
		Name: static.DiRedis,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(static.DiConfig).(config.Provider)
			return redis.NewClient(&redis.Options{
				Addr:     cfg.Config().Cache.Redis.Addr,
				Password: cfg.Config().Cache.Redis.Password,
				DB:       cfg.Config().Cache.Redis.Type,
			}), nil
		},
	})

	// Initialize database middleware and shutdown routine
	diBuilder.Add(di.Def{
		Name: static.DiDatabase,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitDatabase(ctn), nil
		},
		Close: func(obj interface{}) error {
			database := obj.(database.Database)
			logrus.Info("Shutting down database connection...")
			database.Close()
			return nil
		},
	})

	// Initialize twitch notification listener
	diBuilder.Add(di.Def{
		Name: static.DiTwitchNotifyListener,
		Build: func(ctn di.Container) (interface{}, error) {
			return listeners.NewListenerTwitchNotify(ctn), nil
		},
		Close: func(obj interface{}) error {
			listener := obj.(*listeners.ListenerTwitchNotify)
			logrus.Info("Shutting down twitch notify listener...")
			listener.TearDown()
			return nil
		},
	})

	// Initialize twitch notification worker
	diBuilder.Add(di.Def{
		Name: static.DiTwitchNotifyWorker,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitTwitchNotifyWorker(ctn), nil
		},
	})

	// Initialize life cycle timer
	diBuilder.Add(di.Def{
		Name: static.DiScheduler,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitScheduler(ctn), nil
		},
	})

	// Initialize storage middleware
	diBuilder.Add(di.Def{
		Name: static.DiObjectStorage,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitStorage(ctn), nil
		},
	})

	// Initialize permissions command handler middleware
	diBuilder.Add(di.Def{
		Name: static.DiPermissions,
		Build: func(ctn di.Container) (interface{}, error) {
			return permissions.NewPermissions(ctn), nil
		},
	})

	// Initialize discord bot session and shutdown routine
	diBuilder.Add(di.Def{
		Name: static.DiDiscordSession,
		Build: func(ctn di.Container) (interface{}, error) {
			return discordgo.New("")
		},
		Close: func(obj interface{}) error {
			session := obj.(*discordgo.Session)
			logrus.Info("Shutting down bot session...")
			session.Close()
			return nil
		},
	})

	// Initialize Discord OAuth Module
	diBuilder.Add(di.Def{
		Name: static.DiDiscordOAuthModule,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitDiscordOAuth(ctn), nil
		},
	})

	// Initialize auth refresh token handler
	diBuilder.Add(di.Def{
		Name: static.DiAuthRefreshTokenHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewDatabaseRefreshTokenHandler(ctn), nil
		},
	})

	// Initialize auth access token handler
	diBuilder.Add(di.Def{
		Name: static.DiAuthAccessTokenHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewJWTAccessTokenHandler(ctn)
		},
	})

	// Initialize auth API token handler
	diBuilder.Add(di.Def{
		Name: static.DiAuthAPITokenHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewDatabaseAPITokenHandler(ctn)
		},
	})

	// Initialize OAuth API handler implementation
	diBuilder.Add(di.Def{
		Name: static.DiOAuthHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewRefreshTokenRequestHandler(ctn), nil
		},
	})

	// Initialize access token authorization middleware
	diBuilder.Add(di.Def{
		Name: static.DiAuthMiddleware,
		Build: func(ctn di.Container) (interface{}, error) {
			return auth.NewAccessTokenMiddleware(ctn), nil
		},
	})

	// Initialize OTA generator
	diBuilder.Add(di.Def{
		Name: static.DiOneTimeAuth,
		Build: func(ctn di.Container) (interface{}, error) {
			return onetimeauth.NewJwt(&onetimeauth.JwtOptions{
				Issuer: "shinpuru v." + embedded.AppVersion,
			})
		},
	})

	// Initialize backup handler
	diBuilder.Add(di.Def{
		Name: static.DiBackupHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return backup.New(ctn), nil
		},
	})

	// Initialize command handler
	diBuilder.Add(di.Def{
		Name: static.DiCommandHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitCommandHandler(ctn)
		},
		Close: func(obj interface{}) error {
			logrus.Info("Unegister commands ...")
			return obj.(*ken.Ken).Unregister()
		},
	})

	// Initialize web server
	diBuilder.Add(di.Def{
		Name: static.DiWebserver,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitWebServer(ctn), nil
		},
	})

	// Initialize code execution factroy
	diBuilder.Add(di.Def{
		Name: static.DiCodeExecFactory,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitCodeExec(ctn), nil
		},
	})

	// Initialize karma service
	diBuilder.Add(di.Def{
		Name: static.DiKarma,
		Build: func(ctn di.Container) (interface{}, error) {
			return karma.NewKarmaService(ctn), nil
		},
	})

	// Initialize report service
	diBuilder.Add(di.Def{
		Name: static.DiReport,
		Build: func(ctn di.Container) (interface{}, error) {
			return report.New(ctn)
		},
	})

	// Initialize guild logger
	diBuilder.Add(di.Def{
		Name: static.DiGuildLog,
		Build: func(ctn di.Container) (interface{}, error) {
			return guildlog.New(ctn), nil
		},
	})

	// Initialize KV cache
	diBuilder.Add(di.Def{
		Name: static.DiKVCache,
		Build: func(ctn di.Container) (interface{}, error) {
			return kvcache.NewTimedmapCache(10 * time.Minute), nil
		},
	})

	// Initialize State
	diBuilder.Add(di.Def{
		Name: static.DiState,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitState(ctn)
		},
	})

	// Initialize verification service
	diBuilder.Add(di.Def{
		Name: static.DiVerification,
		Build: func(ctn di.Container) (interface{}, error) {
			return verification.New(ctn), nil
		},
	})

	diBuilder.Add(di.Def{
		Name: static.DiBirthday,
		Build: func(ctn di.Container) (interface{}, error) {
			return birthday.New(ctn), nil
		},
	})

	// Build dependency injection container
	ctn := diBuilder.Build()
	// Tear down dependency instances
	defer ctn.DeleteWithSubContainers()

	// Setting log level from config
	cfg := ctn.Get(static.DiConfig).(config.Provider)
	if err := cfg.Parse(); err != nil {
		logrus.WithError(err).Fatal("Failed to parse config")
	}
	logrus.SetLevel(logrus.Level(cfg.Config().Logging.LogLevel))
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 15:04:05 MST",
	})
	if err := config.Validate(cfg); err != nil {
		entry := logrus.NewEntry(logrus.StandardLogger())
		if validationError, ok := err.(config.ValidationError); ok {
			entry = entry.WithField("key", validationError.Key())
		}
		entry.WithError(err).Fatal("Invalid config")
	}

	// Initial log output
	logrus.Info("Starting up...")

	if old, curr, latest := util.CheckForUpdate(); old {
		logrus.
			WithField("current", curr.String()).
			WithField("latest", latest.String()).
			Warn("Update available")
	}

	if profLoc := util.GetEnv(envKeyProfile, flagProfile); profLoc != "" {
		setupProfiler(profLoc)
	}

	if flagDevMode {
		setupDevMode()
	}

	ctn.Get(static.DiCommandHandler)

	// Initialize discord session and event
	// handlers
	releaseShard := inits.InitDiscordBotSession(ctn)
	defer releaseShard()

	// Get Web WebServer instance to start web
	// server listener
	ctn.Get(static.DiWebserver)
	// Get Backup Handler to ensure backup
	// timer is running.
	ctn.Get(static.DiBackupHandler)
	// Get Metrics Server to start metrics
	// endpoint.
	ctn.Get(static.DiMetrics)

	// Block main go routine until one of the following
	// specified exit syscalls occure.
	logrus.Info("Started event loop. Stop with CTRL-C...")

	logrus.WithField("took", startuptime.Took().String()).Info("Initialization finished")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func setupDevMode() {
	if embedded.IsRelease() {
		logrus.Fatal("development mode is not available in production builds")
	}

	util.DevModeEnabled = true

	// Angular dev server
	angServ := angularservice.New(angularservice.Options{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Cd:     "web",
		Port:   8081,
	})
	logrus.Info("Starting Angular dev server...")
	if err := angServ.Start(); err != nil {
		logrus.WithError(err).Fatal("Failed starting Angular dev server")
	}
	defer func() {
		logrus.Info("Shutting down Angular dev server...")
		angServ.Stop()
	}()
}

func setupProfiler(profLoc string) {
	f, err := os.Create(profLoc)
	if err != nil {
		logrus.WithError(err).Fatal("failed starting profiler")
	}
	pprof.StartCPUProfile(f)
	logrus.WithField("location", profLoc).Warning("CPU profiling is active")
	defer pprof.StopCPUProfile()
}
