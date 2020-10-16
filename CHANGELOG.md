1.3.0

> MAJOR PATCH

## Major Implementations

**Karma System Settings** [#146]

The Karma System is now configurable via the web interface. If you have the permission `sp.guild.config.karma`, you will see this entry in the guild settings of the web interface:  
![](https://i.imgur.com/HiaUhJO.png)
There, you can enter the advanced guild settings interface where you can define preferences like enable or disable the karma system, emojis used to increase or decrease karma or the ammount of tokens available per hour per user.  
![](https://i.imgur.com/CC7AgCV.png)

**Channel Locking** [#165]

Added the [`lock`](https://github.com/zekroTJA/shinpuru/wiki/Commands#lock) command to write-lock text channels. By either typing `sp!lock` into the channel which shall be locked or remotely by passing a channel resolvable (i.e. `sp!lock general`), you can write-lock channels. That means, that all roles in this channel (*below the role of the executor*) are explicitly disallowed to write messages in this channel. When hitting the same command again onto this channel, the permission state before the first execution of the `lock` command are restored.  
![](https://i.imgur.com/wLJMDmP.gif)

## Security

Because shinpuru is using cookies containing a singed JWT with session information to authenticate requests against the HTTP REST API, it was vulnerable to [XSRF (Cross-site request forgery) attacks](https://en.wikipedia.org/wiki/Cross-site_request_forgery). This is now fixed by generating session-bound anti-forgery tokens, which are set using the `XSRF-TOKEN` cookie, which is readable by JavaScript. Angular then reads the cookie and sets it as `X-XSRF-TOKEN` header for each following `POST`, `PUT` or `DELETE` request. API-Token based authentications do not need to send the `X-XSRF-TOKEN`, be cause they are already authenticated using headers.

## Bug Fixes

- Permission rules bound to `@everyone` are now correctly processed.
- Also non-command specific permission rules are now correctly listed by the `GET /api/guilds/:guildID/:userID/permissions/allowed` endpoint.

## Backstage

- Refactored the [SQLite3 database middleware](https://github.com/zekroTJA/shinpuru/blob/master/internal/core/middleware/sqlite.go) so that it inherits all bindings from the MySQL middleware which redured it from 952 to only 161 lines of code.
- API handlers are now split up in seperate files for better overview.
- The `GetMemberPermissions` function is now moved from the database middlewares to the permission middleware.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.3.0
```