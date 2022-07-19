# [`internal`](internal/)

This directory contains packages which are internally used in shinpuru.

- [`inits`](inits/): Functions used to wrap initialization procedures of services so that they can be used on service registration in a simple and clean way.
- [`listeners`](listeners/): All registered listeners for Discord events.
- [`messagecommands`](messagecommands/): The application command implementations for message commands.
- [`middleware`](middleware/): Command middlewares registered by [shireikan](https://github.com/zekrotja/shireikan) to check permission, rate limit and log on command execution.
- [`models`](models/): Global data models used across services.
- [`services`](services/): Service definitions like the database binding, web server, storage binding, karma service, report service and more.
- [`slashcommands`](slashcommands/): The application command implementations for slash commands.
- [`usercommands`](usercommands/): The application command implementations for user commands.
- [`util`](util/): Collection of utility functionalities which utilize other internal packages and models.