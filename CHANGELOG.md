1.11.0

Hey, zekro here. 👋  
I want to lose some quick words for this patch, because, on the first glance, it does not include that much rich updates for the end users, but **huge** changes were made in this update under the hood of shinpuru. If you are only interested by these changes, please click [here](#minor-changes).

Because I was really unsatisfied and actually almost annoyed working with the old dependency and service structure as well as the very badly designed web server structure up to the point that I was shied away from inplementing new features because of that, I decided to completely redesign the dependency management and web server implementation of shinpuru.

By the way, I have also written a [**Contribution Sheet**](https://github.com/zekroTJA/shinpuru/blob/master/CONTRIBUTING.md), where the whole structure — including the new implementations — are documented and explained in detail.

### The Package Structure

The old `internal/` package structure was kind of messy and inconsistent. So, I've re-ordered nearly all of the internal packages. The `models/` package, which contains any object data models used across services, is now directly under `internal/`. Also, I've more clearely separated `services/`, which are stateful instances used across the application and `util/`, which mainly contain stateless utility functions.

### The Dependency Injection System

The old DI system was based on simply passing service instances to service constructor functions. On the first hand, that required a specific initialization order in the `main()` function so that a service depending on another is not initialized before that service. Also, if I wanted to add a new dependency to a service, I needed to add it to the services constructor as well as wher the constructor was called, which was kind of unhandy.

Now, shinpuru uses [sarulabs/di](https://github.com/sarulabs/di) for that. It allows to use a service builder where you define *how* your services are built and what they need. Then, you can build a service container from that builder. Now, the service container cares about creating services when they are needed in the order they require, which makes service registration very easy. Also, you can define teardown functions which are called, when a service instance is shut down. Now, we only pass the whole service container to each service and each service can take whichever service it needs *(of course while respecting a tree shaped dependency pattern — it is also not magic to resolve cyclic dependencies!)*.

### The Web Server & REST API

The whole structure of the REST API And web server routes was horrible. At some point, I just started copy-pasting request handlers over and over again and repeating code instead of abstracting functionalities into middlewares or helper functions. Even tho I had helper functions, which were necessary becasue `fasthttp` and `fasthttp-routing` — which were both used for the web server and API routing — these were unhandy and annoying to use. Also, the whole handling of static files, "middlewares" and routes was really clunky and unmodular, which also made it really unmaintainable over time.

So, I have almost re-written the whole web server using much more modern patterns. The new API is based on the [fiber](https://github.com/gofiber/fiber) web framework, which has a way cleaner and more fun to use API.

The API can now be versioned using a `Router` implementation for each version. Endpoints are split up into `Controllers`, which hold the specific endpoint functions. Also, I heavily depend on middlewares which make permission control way more convenient now.

Also, I've thrown away the *"old"* sesion-cookie and anti-forgery-key dependent authentication system and replaced it with a way more modern, controllabe, persistent and performant access-refresh-token based authentication system.

### The new LifeCycleTimer

Another thing that I wanted to update was the `LifeCycleTimer` implementation, which is a central timer taking care of tasks like checking for expired votes, creating guild backups and cleaning up expired refresh tokens. The old impementation was based on delays between job execution. For example, guild backups are created every 12 hours. This timer starts after the bot has initialized. Now, when the instance had to be restarted for an update, for example, the timer was reset. That makes job execution very unreliable, which is especially a problem for crucial tasks like guild backup creation.

Now, a new implementation which uses the package [robfig/cron](https://github.com/robfig/cron) handles these tasks on a time schedule based system. The package allows creating shedules by using a cron like syntax, which also makes it very simple to pass these schedule values by the config. Also, tasks like backup creation can now be based on a fix time schedule which is not reset by restarts and is way more reliable and predictable.

## Minor Changes

- You can now configure a rate limit for all REST API endpoints in the config.  
  ```yaml
  webserver:
    # Ratelimit configuration
    ratelimit:
      # Whether or not to enable rate limiter
      enabled: true
      # Max requests in the given duration.
      # This value should not be that low, because first
      # connections to the API via the web interface might
      # require a lot of requests to be processed.
      max: 50
      # The reset duration until a rate limit exceeds.
      durationseconds: 3
  ```

- Also, you are now able to configure the time schedules of tasks like backup creation or expired refresh token cleanup. The schedule uses a cron like syntax. Read more about that [here](https://pkg.go.dev/github.com/robfig/cron/v3#hdr-Usage).  
  ```yaml
  schedules:
    # Guild backup schedule
    guildbackups:        '0 0 6,18 * * *'
    # Refresh token cleanup schedule
    refreshtokencleanup: '0 0 5 * * *'
  ```

- You can now also provide configuration via a JSON file with the same format as the YAML configuration.

## Bug Fixes

- Bots don't gain Karma points anymore when their messages get into the starboard. Also, karma counts are now hidden on bot accounts in the web interface. [#224]
- Moderation section is now hidden on the user page in the web interface when the auhtorized user is not permitted to take any moderation action.
- Explicit enable or disable of the color reactions using the [`color`](https://github.com/zekroTJA/shinpuru/wiki/Commands#color) command now properly applies the setting to the database. [#225]

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.11.0
```