package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"

	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/listeners"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/inits"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/onetimeauth"

	"github.com/zekroTJA/shinpuru/pkg/angularservice"
)

var (
	flagConfigLocation = flag.String("c", "config.yml", "The location of the main config file")
	flagDocker         = flag.Bool("docker", false, "wether shinpuru is running in a docker container or not")
	flagDevMode        = flag.Bool("devmode", false, "start in development mode")
	flagProfile        = flag.String("cpuprofile", "", "Records a CPU profile to the desired location")
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

	// Initialize dependency injection builder
	diBuilder, _ := di.NewBuilder()

	// Setup config parser
	diBuilder.Add(di.Def{
		Name: static.DiConfigParser,
		Build: func(ctn di.Container) (interface{}, error) {
			return new(config.YAMLConfigParser), nil
		},
	})

	// Initialize config
	diBuilder.Add(di.Def{
		Name: static.DiConfig,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitConfig(*flagConfigLocation, ctn), nil
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
			return inits.InitTwitchNotifyListener(ctn), nil
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
			return inits.InitLTCTimer(), nil
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
			return inits.InitDiscordBotSession(ctn), nil
		},
		Close: func(obj interface{}) error {
			session := obj.(*discordgo.Session)
			util.Log.Info("Shutting down bot session...")
			session.Close()
			return nil
		},
	})

	// Initialize OTA generator
	diBuilder.Add(di.Def{
		Name: static.DiOneTimeAuth,
		Build: func(ctn di.Container) (interface{}, error) {
			return onetimeauth.New(&onetimeauth.Options{
				Issuer: "shinpuru v." + util.AppVersion,
			})
		},
	})

	// Initialize backup handler
	diBuilder.Add(di.Def{
		Name: static.DiBackupHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitBackupHandler(ctn), nil
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

	// Get Discord Session to initialize Discord
	// session and connection
	ctn.Get(static.DiDiscordSession)
	// Get Web WebServer instance to start web
	// server listener
	ctn.Get(static.DiWebserver)

	// Block main go routine until one of the following
	// specified exit syscalls occure.
	util.Log.Info("Started event loop. Stop with CTRL-C...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Tear down dependency instances
	ctn.DeleteWithSubContainers()
}

func setupDevMode() {
	if util.Release == "TRUE" {
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
