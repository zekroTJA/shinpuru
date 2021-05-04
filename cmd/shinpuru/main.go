package main

import (
	"errors"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"

	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/inits"
	"github.com/zekroTJA/shinpuru/internal/listeners"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	"github.com/zekroTJA/shinpuru/internal/services/backup"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/karma"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/auth"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/startupmsg"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/onetimeauth/v2"
	"github.com/zekroTJA/shinpuru/pkg/startuptime"
	"github.com/zekroTJA/shireikan"

	"github.com/zekroTJA/shinpuru/pkg/angularservice"
)

var (
	flagConfig  = flag.String("c", "config.yml", "The location of the main config file")
	flagDocker  = flag.Bool("docker", false, "wether shinpuru is running in a docker container or not")
	flagDevMode = flag.Bool("devmode", false, "start in development mode")
	flagProfile = flag.String("cpuprofile", "", "Records a CPU profile to the desired location")
	flagQuiet   = flag.Bool("quiet", false, "Dont print startup message")
)

const (
	envKeyProfile = "CPUPROFILE"
)

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
	flag.Parse()

	if !*flagQuiet {
		startupmsg.Output(os.Stdout)
	}

	// Initialize dependency injection builder
	diBuilder, _ := di.NewBuilder()

	// Setup config parser
	diBuilder.Add(di.Def{
		Name: static.DiConfigParser,
		Build: func(ctn di.Container) (p interface{}, err error) {
			ext := strings.ToLower(filepath.Ext(*flagConfig))
			switch ext {
			case ".yml", ".yaml":
				p = new(config.YAMLConfigParser)
			case ".json":
				p = new(config.JSONConfigParser)
			default:
				err = errors.New("unsupported configuration file")
			}
			return
		},
	})

	// Initialize config
	diBuilder.Add(di.Def{
		Name: static.DiConfig,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitConfig(*flagConfig, ctn), nil
		},
	})

	// Initialize metrics server
	diBuilder.Add(di.Def{
		Name: static.DiMetrics,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitMetrics(ctn), nil
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
			util.Log.Info("Shutting down database connection...")
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
			util.Log.Info("Shutting down twitch notify listener...")
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
		Name: static.DiLifecycleTimer,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitLTCTimer(ctn), nil
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
		Name: static.DiPermissionMiddleware,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitPermissionMiddleware(ctn), nil
		},
	})

	// Initialize ghost ping ignore command handler middleware
	diBuilder.Add(di.Def{
		Name: static.DiGhostpingIgnoreMiddleware,
		Build: func(ctn di.Container) (interface{}, error) {
			return middleware.NewGhostPingIgnoreMiddleware(), nil
		},
	})

	// Initialize discord bot session and shutdown routine
	diBuilder.Add(di.Def{
		Name: static.DiDiscordSession,
		Build: func(ctn di.Container) (interface{}, error) {
			return discordgo.New()
		},
		Close: func(obj interface{}) error {
			session := obj.(*discordgo.Session)
			util.Log.Info("Shutting down bot session...")
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
				Issuer: "shinpuru v." + util.AppVersion,
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
			return inits.InitCommandHandler(ctn), nil
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

	// Build dependency injection container
	ctn := diBuilder.Build()

	// Setting log level from config
	cfg := ctn.Get(static.DiConfig).(*config.Config)
	util.SetLogLevel(cfg.Logging.LogLevel)

	// Initial log output
	util.Log.Infof("シンプル (shinpuru) v.%s (commit %s)", util.AppVersion, util.AppCommit)
	util.Log.Info("© zekro Development (Ringo Hoffmann)")
	util.Log.Info("Covered by MIT Licence")
	util.Log.Info("Starting up...")

	if profLoc := util.GetEnv(envKeyProfile, *flagProfile); profLoc != "" {
		setupProfiler(profLoc)
	}

	if *flagDevMode {
		setupDevMode()
	}

	// Initialize discord session and event
	// handlers
	inits.InitDiscordBotSession(ctn)

	// This is currently the really hacky workaround
	// to bypass the di.Container when trying to get
	// the Command handler instance inside a command
	// context, because the handler can not resolve
	// itself on build, so it is bypassed here using
	// shireikans object map. Maybe I find a better
	// solution for that at some time.
	handler := ctn.Get(static.DiCommandHandler).(shireikan.Handler)
	handler.SetObject(static.DiCommandHandler, handler)

	// Get Web WebServer instance to start web
	// server listener
	ctn.Get(static.DiWebserver)
	// Get Backup Handler to ensure backup
	// timer is running.
	ctn.Get(static.DiBackupHandler)

	// Block main go routine until one of the following
	// specified exit syscalls occure.
	util.Log.Info("Started event loop. Stop with CTRL-C...")

	util.Log.Infof("Initialization finished - took %s", startuptime.Took().String())
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Tear down dependency instances
	ctn.DeleteWithSubContainers()
}

func setupDevMode() {
	if util.IsRelease() {
		util.Log.Fatal("development mode is not available in production builds")
	}

	util.DevModeEnabled = true

	// Angular dev server
	angServ := angularservice.New(angularservice.Options{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Cd:     "web",
		Port:   8081,
	})
	util.Log.Info("Starting Angular dev server...")
	if err := angServ.Start(); err != nil {
		util.Log.Fatalf("Failed starting Angular dev server: %s", err.Error())
	}
	defer func() {
		util.Log.Info("Shutting down Angular dev server...")
		angServ.Stop()
	}()
}

func setupProfiler(profLoc string) {
	f, err := os.Create(profLoc)
	if err != nil {
		util.Log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	util.Log.Warningf("CPU profiling is active (loc: %s)", profLoc)
	defer pprof.StopCPUProfile()
}
