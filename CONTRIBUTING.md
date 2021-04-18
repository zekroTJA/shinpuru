# How To Contribute

First of all, everyone of you is welcome to contribute to the project. Whether small changes like typo fixes, simple bug fixes or large feature implementations, every contribution is a step further to make this project way nicer. 😄

Depending on the scale of the contribution, you might need some general understanding of the languages and frameworks used and of some simple development pattern which are applied. But also beginners are absolutely welcome.

## Used Languages/Frameworks and general Structure

Let me give you a quick overview over all modules and structures and used languages and frameworks of the project, so you know what you are working with.

### Config

shinpuru is configured using either a YAML or JSON config file which can be passed by `-c` command line parameter.

Take a look at the [**example configuration file**](https://github.com/zekroTJA/shinpuru/blob/master/config/config.example.yaml) which holds rich documentation about each configuration key.

For development, you can take the provided [`my.private.config.yml`](https://github.com/zekroTJA/shinpuru/blob/master/config/my.private.config.yaml) and enter your credentials. Then, rename it to `private.config.yaml`, so your secrets will not be commited to the repository by accident.

### Discord Communication

The backend of shinpuru is completely written in [**Go**](https://go.dev/) *(golang)*. To communicate with Discord, the API wrapper [**discordgo**](https://github.com/bwmarrin/discordgo) is used. discordgo provides very low level bindings to the Discord API with very little utilities around, therefore a lot of utility packages were created. These can be found in the `pkg/` directory. These are the main utility packages used in shinpuru:
- [acceptmsg](https://github.com/zekroTJA/shinpuru/tree/master/pkg/acceptmsg) creates an embed message with a ✔ and ❌ reaction added. Then, you can execute code depending on which reaction was clicked on.
- [discordutil](https://github.com/zekroTJA/shinpuru/tree/master/pkg/discordutil) provides general utility functions like getting message links, retrieving objects first from the discordgo cache and, when not available there, from the Discord API or checking if a user has admin privileges.
- [embedbuilder](https://github.com/zekroTJA/shinpuru/tree/master/pkg/embedbuilder) helps building embeds using the builder pattern.
- [fetch](https://github.com/zekroTJA/shinpuru/tree/master/pkg/fetch) is a widely used package in shinpuru used to get objects like users, members or roles by either ID, name or mention. This is designed to be as fuzzy as possible matching objects to provide a better experience to the user.

Take a look at the packages in the [pkg](https://github.com/zekroTJA/shinpuru/tree/master/pkg) yourself. All of them are as well documented as I was possible to and some also have unit tests where you can see some examples how to use them. 😉

Also, a lot of shared functionalities which require shinpuru specific dependencies are located in the [internal/util](https://github.com/zekroTJA/shinpuru/tree/master/internal/util) directory. There you can find some utilities which can be used to access the imagestore, karma system, metrics or votes.

For command handling, shinpuru uses [**shireikan**](https://github.com/zekroTJA/shireikan). Take a look there and in the examples. Just like that, commands are handled and defined in shinpuru. All command definitions can be found in the [`internal/commands`](https://github.com/zekroTJA/shinpuru/tree/master/internal/commands) directory. If you want to add a command, just implement shireikans [`Command`](https://github.com/zekroTJA/shireikan/blob/master/command.go) interface and take a look how the other commands are implemented to match the conventions applied in the other commands. After that, register the command in the [`cmdhandler`](https://github.com/zekroTJA/shinpuru/blob/master/internal/inits/cmdhandler.go) `InitCommandHandler()` function using the `cmdHandler.RegisterCommand(&commands.YourCmd{})` method.

Discord event handlers and listeners can be found in the [`listeners`](https://github.com/zekroTJA/shinpuru/tree/master/internal/core/listeners) package. A listener is a struct which exposes one or more event handler methods. Listeners must be registered [`botsession`](https://github.com/zekroTJA/shinpuru/blob/master/internal/inits/cmdhandler.go) `InitDiscordBotSession()` function using the `session.AddHandler(listeners.NewYourListener(container).Handler)` method.

### Database

First of all, you can find a [`Database`](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/database/database.go) interface at `internal/core/database`. This is mainly used to interact with the database. There, you can also find the specific database drivers available, which are currently [`mysql`](https://github.com/zekroTJA/shinpuru/tree/master/internal/core/database/mysql), [`sqlite`](https://github.com/zekroTJA/shinpuru/tree/master/internal/core/database/sqlite) and [`redis`](https://github.com/zekroTJA/shinpuru/tree/master/internal/core/database/redis).

shinpuru mainly uses MySQL/MariaDB as database. You *can* also use SQLite3 for development, but this is not tested anymore and may not be reliable anymore. It is recommended to set up a MariaDB instance on your server or dev system for development. Here you can find some resources how to set up MariaDB on mainly used systems:
- Windows: https://mid.as/kb/00197/install-configure-mariadb-on-windows
- Linux: https://opensource.com/article/20/10/mariadb-mysql-linux
- Docker: https://hub.docker.com/_/mariadb/

Redis is used as database cache. The [`RedisMiddleware`](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/database/redis/redis.go) generaly inherits functionalities from the specified database middleware instance and only overwrites using the specified functions. The database cache always keeps the cache as well as the database hot and always first tries to get objects from cache and, if not available there, from database.

![](https://i.imgur.com/TgkuhUY.png)

If you want to add functionalities to the database in your contributions, add the functions to the database interface as well as to the MySQL database driver and, if you need caching, the middleware functions to the redis caching middleware.

If you want to add a column to an existing table, take a look in the [`migrations`](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/database/mysql/migrations.go) implementation. There, you can add a migration function with the SQL statements which will be executed in order to migrate the database structure to the new state. If you add an entirely new table, you don't need to add a migration function. Just add the table definition in the `setup()` method in the [`mysql`](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/database/mysql/mysql.go) driver.

> The `MysqlMiddleware` is very "low level" and directly works with SQL statements instead of using an ORM or something like this. Don't be overwhelmed by the size of the middleware file. Its just because same functionalities are re-used over and over again, which is not very nice, but to be honest, the middleware is very old and I don't find the time to rewrite it and migrate the current database after that.

### Storage

shinpuru utilizes a simple object storage for storing images and backup files, described by the [`Storage`](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/storage/storage.go) interface in [`internal/core/storage`](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/storage). Currently, shinpuru implements two storage drivers: A firect [file storage](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/storage/file.go) driver and a [minio object storage](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/storage/minio.go) driver, which can also connect to other object storages like Amazon S3 or Google cloud storage.

### REST API

The web interface communicates with the shinpuru backend over a RESTful HTTP API. Therefore, [**fiber**](https://gofiber.io/) is used as HTTP framework. Most of the code of the web server is in the [`internal/core/webserver`](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/webserver) directory. The web server is split up in `Router`'s and `Controller`'s. Routers are for versioning the API *(e.g. `/api/v1`, `/api/v2`, ...)* and [`controllers`](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/webserver/v1/controllers) split up the endpoints in different logical sections *(e.g. `/guilds`, `/backups`, `/guilds/:id/members`, ...)*. Also, there are [`models`](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/webserver/v1/models), which define the object structure of request and response objects as well as some transformation functions, for example to transform a `discordgo.Guild` object to a `models.Guild` object.

If you want to add API endpoints, just add the endpoints to one of the controllers *(don't forget to register the endpoint in the controller's `Setup` method!)*, or create a new entire controller, which then needs to be registered in the API `Route`. If you need service dependencies in your controller, just add it to the controllers struct and get it from the passed `di.Container` *(more explained below)* in the `Setup` method.

Also, fiber works a lot with middlewares, which can be chained anywhere into the fiber route chain. In shinpurus implementation, there are three main types of middlewares.
1. The high level middlewares like the rate limiter, CORS or file system middleware, which are set before all incomming requests.
2. Controller specific middlewares which are defined in the router. Mainly, this is used for the authorizeation middleware, which checks for auth tokens in the requests. This middleware is required by some controllers and not required for others.
3. Endpoint specific middlewares which are defined for specific endpoints only. Mainly, this is used for the permission middleware which checks for required user permissions to execute specific endpoints.

Here you can see a simple overview over the routing structure of the shinpuru webserver.
![](https://i.imgur.com/VFuU7rj.png)

### Dependency Injection

> If you are unfamiliar with the concepts of dependency injection, please read this [**blog post**](https://blog.zekro.de/dependency-injection) I have recently written about DI, also with examples in Go. 😉

shinpuru widely uses DI *(dependency injection)* to share service instances using the package [**di**](https://github.com/sarulabs/di) from sarulabs. It's a really straight forward implementation of a DI container which does not take use of reflection, which makes it quite simple and fast. Also, the `di` cares about constructing service instances when they are needed and tearing them down when they are no more needed.

The whole service specification happens in the main function of shinpuru in the [`cmd/shinpuru/main`](https://github.com/zekroTJA/shinpuru/blob/master/cmd/shinpuru/main.go) file. For example, the database service initialization looks like following:

```go
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
```

As you can see, all service identifiers are registered in the [`internal/util/static/di`](https://github.com/zekroTJA/shinpuru/blob/master/internal/util/static/di.go) file.

After building the `diBuilder`, you will have a `di.Container` to work with where you can get any service registered. Because all services are registered in the `App` scope, once they are initialized, all requests are getting the same instance of the service. This makes service development very easy, because every service is getting passed the same service container and every service can grab the instance of any other registered service instance.

When you want to use a service, just take it from the passed service conatiner by the specified identifier. Let's take a look at the [`starboard` listener](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/listeners/starboard.go) , for example:

```go
func NewListenerStarboard(container di.Container) *ListenerStarboard {
	cfg := container.Get(static.DiConfig).(*config.Config)
	var publicAddr string
	if cfg.WebServer != nil {
		publicAddr = cfg.WebServer.PublicAddr
	}

	return &ListenerStarboard{
		db:         container.Get(static.DiDatabase).(database.Database),
		st:         container.Get(static.DiObjectStorage).(storage.Storage),
		publicAddr: publicAddr,
	}
}
```

As you can see, the `NewListenerStarboard` function is getting passed the `di.Container` from somewhere above. Then, the config is taken from the container to resolve the public web server address, if specified. Also, the database as well as the storage service instance is retrieved.

The only thing important to keep in mind is that you should always build your service dependency structure like a tree, and not like a circle. That means, when service `A` needs service `B` to be built, service `B` can not depend on service `A` on construction.

![](https://i.imgur.com/8hTVWC3.png)

### Web Frontend

The shinpuru web frontend is a compiled [**Angular**](https://angular.io) SPA, which is directly hosted form the shinpuru web server. Stylesheets are written in [**SCSS**](https://sass-lang.com/documentation/syntax) because SCSS has huge advantages to default CSS like nesting, mixins and variables, which are widely used in stylesheets.


## Preparing a Development Environment

First of all, create a fork of this repository.  
![](https://i.imgur.com/V0uP5lu.png)

Then, clone the repository to your PC either using HTTPS, SSH or the Git CLI.  
![](https://i.imgur.com/xWOovnk.png)

Of course, you need to download and install the **Go compiler toolchain**. Please follow [**these**](https://golang.org/doc/install) instructions to do so.

Also, you need some sort of a C compiler like `gcc`. If you are using linux, just [install C build tools](https://www.ubuntupit.com/how-to-install-and-use-gcc-compiler-on-linux-system/) using your package manager. When you are using windows, you can download and install [TDM-GCC](https://sourceforge.net/projects/tdm-gcc/).

Also, to compile the web frontend, you need to install NodeJS and NPM. Please follow [**these**](https://nodejs.org/en/) instructions to do so.

This repository also provides a [`Makefile`](https://github.com/zekroTJA/shinpuru/blob/master/Makefile) with a lot of useful recipies for development. Just enter `make help` to get a quick overview over all make recipes.

> Read [this](https://www.cyberciti.biz/faq/howto-installing-gnu-c-compiler-development-environment-on-ubuntu/) to install GNU Make on Linux and [this](http://gnuwin32.sourceforge.net/packages/make.htm) to install it on Windows.

## Any Questions?

If you have any questions, please hit me on my [**Dev Discord**](https://discord.zekro.de) (`zekro#0001`) or on [**Twitter**](https://twitter.com/zekrotja). You can also simply send me an [e-mail](mailto:contact@zekro.de). 😉