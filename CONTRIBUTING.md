# How To Contribute

First of all, everyone of you is welcome to contribute to the project. Whether small changes like typo fixes, simple bug fixes or large feature implementations, every contribution is a step further to make this project way nicer. 😄

Depending on the scale of the contribution, you might need some general understanding of the languages and frameworks used and of some simple development pattern which are applied. But also beginners are absolutely welcome.

## Used Languages, Frameworks and Techniques

Let me give you a quick overview over all used languages and frameworks of the project, so you know what you are working with.

### Discord Communication

The backend of shinpuru is completely written in [**Go**](https://go.dev/) *(golang)*. To communicate with Discord, the API wrapper [**discordgo**](https://github.com/bwmarrin/discordgo) is used. discordgo provides verry low level bindings to the Discord API with very little utilities around, therefore a lot of utility packages were created. These can be found in the `pkg/` directory. These are the main utility packages used in shinpuru:
- [acceptmsg](https://github.com/zekroTJA/shinpuru/tree/master/pkg/acceptmsg) creates an embed message with a ✔ and ❌ reaction added. Then, you can execute code depending on which reaction was clicked on.
- [discordutil](https://github.com/zekroTJA/shinpuru/tree/master/pkg/discordutil) provides general utility functions like getting message links, retrieving objects first from the discordgo cache and, when not available there, from the Discord API or checking if a user has admin previleges.
- [embedbuilder](https://github.com/zekroTJA/shinpuru/tree/master/pkg/embedbuilder) helps building embeds using the builder pattern.
- [fetch](https://github.com/zekroTJA/shinpuru/tree/master/pkg/fetch) is a widely used package in shinpuru used to get objects like users, members or roles by either ID, name or mention. This is designed to be as fuzzy as possible matching objects to provide a better experience to the user.

Take a look at the packages in the [pkg](https://github.com/zekroTJA/shinpuru/tree/master/pkg) yourself. All of them are as well documented as I was possible to and some also have unit tests where you can see some examples how to use them. 😉

Also, a lot of shared functionalities which require shinpuru speicific dependencies are located in the [internal/util](https://github.com/zekroTJA/shinpuru/tree/master/internal/util) directory. There you can find some utilities whcih can be used to access the imagestore, karma system, mectrics or votes.

### Database

First of all, you can find a [`Database`](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/database/database.go) interface at `internal/core/database`. This is mainly used to interact with the database. There, you can also find the specific database drivers available, which are currently [`mysql`](https://github.com/zekroTJA/shinpuru/tree/master/internal/core/database/mysql), [`sqlite`](https://github.com/zekroTJA/shinpuru/tree/master/internal/core/database/sqlite) and [`redis`](https://github.com/zekroTJA/shinpuru/tree/master/internal/core/database/redis).

shinpuru mainly uses MySQL/MariaDB as database. You *can* also use SQLite3 for development, but this is not tested anymore and may not reliable anymore. It is recommendet to set up a MariaDB instance on your server or dev system for development. Here you can find some resources how to set up MariaDB on mainly used systems:
- Windows: https://mid.as/kb/00197/install-configure-mariadb-on-windows
- Linux: https://opensource.com/article/20/10/mariadb-mysql-linux
- Docker: https://hub.docker.com/_/mariadb/

Redis is used as database cache. The [`RedisMiddleware`](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/database/redis/redis.go) generaly inherits functionalities from the specified database middleware instance and only overwrites using the specified functions. The database cache always keeps the cache as well as the database hot and always first tries to get objects from cache and, if not available there, from database.

![](https://i.imgur.com/TgkuhUY.png)

If you want to add functionalities to the database in your contributions, add the functions to the database interface as well as to the MySQL database driver and, if you need caching, the middleware functions to the redis caching middleware.

If you want to add a column to an existing table, take a look in the [`migrations`](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/database/mysql/migrations.go) implementation. There, you can add a migration function with the SQL statements which will be executed in order to migrate the database structure to the new state.

> The `MysqlMiddleware` is very "low level" and directly works with SQL statements instead of using an ORM or something like this. Don't be overwhelmed by the size of the middleware file. Its just because same functionalities are re-used over and over again, which is not very nice, but to be honest, the middleware is very old and I don't find the time to rewrite it and migrate the current database after that.