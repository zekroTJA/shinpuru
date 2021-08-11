1.21.1

## Changes

### Configuration Update

The `.database.redis` configuration has now been moved to `.cache.redis`, because the redis instance is also used for state caching and not only for database caching anymore. Also, `.database.redis` has now been marked as **deprecated and will be completely removed in upcoming patches**.

Also, a new settings key `.cache.cachedatabase` (type: `boolean`, default: `true`) has been added which enables or disables database request caching in Redis.

From [config/config.example.yaml](https://github.com/zekroTJA/shinpuru/blob/master/config/config.example.yaml)
```yaml
# Caching prefrences.
cache:
  # Redis connection configuration.
  redis:
    # Redis host address
    addr: "localhost:6379"
    # Redis password
    password: "myredispassword"
    # Database type
    type: 0
  # If enabled, most frequently used database
  # requests are automatically cached in redis
  # to minimize load on the database as well as
  # request times.
  # It is recomendet to leave this enabled. If
  # you want to disable it for whatever reason,
  # you can do it here.
  cachedatabase: true
```

## Bugfixes

- In the guild settings page, the navbar now only shows the sections the current user actually has permission access to. [#268]
- "Anonymous Reports" have now been renamed to "Ghost Reports" to prevent missunderstandings. Also added a feature explanation to the Ghost Report modal. [#270, #271]
- The global search checks now for login status before being shown. [#273]

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.21.1
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.21.1
```
