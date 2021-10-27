<div align="center">
    <img src=".media/rendered/sp-banner-slim.png" width="100%" />
    <hr>
    <h1>~ シンプル ~</h1>
    <strong>
        A simple multi purpose discord bot written in Go (discord.go)<br>
        with focus on stability and reliability
    </strong><br><br>
    <a href="https://dc.zekro.de"><img height="28" src="https://img.shields.io/discord/307084334198816769.svg?style=for-the-badge&logo=discord" /></a>&nbsp;
    <a href="https://github.com/zekroTJA/shinpuru/releases"><img height="28" src="https://img.shields.io/github/tag/zekroTJA/shinpuru.svg?style=for-the-badge"/></a>&nbsp;
    <a href="https://hub.docker.com/r/zekro/shinpuru"><img alt="Docker Cloud Automated build" src="https://img.shields.io/docker/cloud/automated/zekro/shinpuru.svg?color=cyan&logo=docker&logoColor=cyan&style=for-the-badge"></a>&nbsp;
    <img height="28" src="https://forthebadge.com/images/badges/built-with-grammas-recipe.svg">
<br>
</div>

---

| Branch | Main CI | Docker CD | Releases CD |
|--------|---------|-----------|-------------|
| stable | [![](https://github.com/zekroTJA/shinpuru/workflows/Main%20CI/badge.svg?branch=master)](https://github.com/zekroTJA/shinpuru/actions?query=workflow%3A%22Main+CI%22+branch%3Amaster) | [![](https://github.com/zekroTJA/shinpuru/workflows/Docker%20CD%20Latest/badge.svg)](https://github.com/zekroTJA/shinpuru/actions?query=workflow%3A%22Docker+CD%22+branch%3Amaster) | [![](https://github.com/zekroTJA/shinpuru/workflows/Releases%20CD/badge.svg?branch=master)](https://github.com/zekroTJA/shinpuru/actions?query=workflow%3A%22Releases+CD%22+branch%3Amaster)
| canary    | [![](https://github.com/zekroTJA/shinpuru/workflows/Main%20CI/badge.svg?branch=dev)](https://github.com/zekroTJA/shinpuru/actions?query=workflow%3A%22Main+CI%22+branch%3Adev) | [![](https://github.com/zekroTJA/shinpuru/workflows/Docker%20CD%20Canary/badge.svg)](https://github.com/zekroTJA/shinpuru/actions?query=workflow%3A%22Docker+CD%22+branch%dev) | |

---

# Invite

Here you can choose between the stable or canary version of shinpuru:

<a href="https://shnp.de/invite"><img src="https://img.shields.io/badge/%20-INVITE%20STABLE-0288D1.svg?style=for-the-badge&logo=discord" height="30" /></a>

<a href="https://c.shnp.de/invite"><img src="https://img.shields.io/badge/%20-INVITE%20CANARY-FFA726.svg?style=for-the-badge&logo=discord" height="30" /></a>

# Intro

シンプル (shinpuru), a simple *(as the name says)*, multi-purpose Discord Bot written in Go, using bwmarrin's package [discord.go](https://github.com/bwmarrin/discordgo) as API and gateway wrapper and [ken](https://github.com/zekroTJA/ken) as slash command framework. The focus on this bot is to provide general purpose, administration and security tools while keeping stability, reliability and maintainability.

This bot is mainly used as administration and security tool on my [development discord](https://discord.zekro.de). Drop by to see shinpuru in action! 😉

---

# Features 

## Slash Commands

shinpuru mainly uses slash commands to interact with the bot. In the [**wiki**](https://github.com/zekroTJA/shinpuru/wiki/Commands), you can find an automatically generated list of commands, their descriptions and how to use them.

You can also find a searchable list in the [**web interface**](https://shnp.de/commands) of shinpuru.

https://user-images.githubusercontent.com/16734205/138589141-1cc18316-0d07-4526-b86a-be5aa91bbc5a.mp4

## Web Interface

If you are sick of using chat commands, you can also use the web interface of shinpuru. Simply log in with your Discord Account (alternatively, you can also use the [`/login`](https://github.com/zekroTJA/shinpuru/wiki/Commands#login) command).

https://user-images.githubusercontent.com/16734205/138589590-87301377-463d-43c3-8441-98ec84a1304c.mp4

---

# Docker

Read about how to host shinpuru using the provided Docker image in the [**wiki article**](https://github.com/zekroTJA/shinpuru/wiki/Docker).

---

# Compiling

Read about self-compiling in the [**wiki article**](https://github.com/zekroTJA/shinpuru/wiki/Self-Compiling).

---
<!-- start:PUBLIC_PACKAGES -->
# Public Packages

- [**`github.com/zekroTJA/shinpuru/pkg/acceptmsg`**](pkg/acceptmsg)  
  *Package acceptmsg provides a message model for discordgo which can be accepted or declined via message reactions.*

- [**`github.com/zekroTJA/shinpuru/pkg/angularservice`**](pkg/angularservice)  
  *Package angularservice provides bindings to start an Angular development server via the Angular CLI.*

- [**`github.com/zekroTJA/shinpuru/pkg/argp`**](pkg/argp)  
  *Package argp is a stupid simple flag (argument) parser which allows to parse flags without panicing when non-registered flags are passed.*

- [**`github.com/zekroTJA/shinpuru/pkg/boolutil`**](pkg/boolutil)  
  *Package boolutil provides simple utility functions around booleans.*

- [**`github.com/zekroTJA/shinpuru/pkg/bytecount`**](pkg/bytecount)  
  *Package bytecount provides functionalities to format byte counts.*

- [**`github.com/zekroTJA/shinpuru/pkg/checksum`**](pkg/checksum)  
  *Package checksum provides functions to generate a hash sum from any given object.*

- [**`github.com/zekroTJA/shinpuru/pkg/colors`**](pkg/colors)  
  *Package color provides general utilities for image/color objects and color codes.*

- [**`github.com/zekroTJA/shinpuru/pkg/ctypes`**](pkg/ctypes)  
  *Package ctype provides some custom types with useful function extensions.*

- [**`github.com/zekroTJA/shinpuru/pkg/discordoauth`**](pkg/discordoauth)  
  *package discordoauth provides fasthttp handlers to authenticate with via the Discord OAuth2 endpoint.*

- [**`github.com/zekroTJA/shinpuru/pkg/discordutil`**](pkg/discordutil)  
  *Package discordutil provides general purpose extensuion functionalities for discordgo.*

- [**`github.com/zekroTJA/shinpuru/pkg/embedbuilder`**](pkg/embedbuilder)  
  *Package embedbuilder provides a builder pattern to create discordgo message embeds.*

- [**`github.com/zekroTJA/shinpuru/pkg/etag`**](pkg/etag)  
  *Package etag implements generation functionalities for the ETag specification of RFC7273 2.3. https://tools.ietf.org/html/rfc7232#section-2.3.1*

- [**`github.com/zekroTJA/shinpuru/pkg/fetch`**](pkg/fetch)  
  *Package fetch provides functionalities to fetch roles, channels, members and users by so called resolavbles. That means, these functions try to match a member, role or channel by their names, displaynames, IDs or mentions as greedy as prossible.*

- [**`github.com/zekroTJA/shinpuru/pkg/hammertime`**](pkg/hammertime)  
  *Package hammertime provides functionailities to format a time.Time into a Discord timestamp mention. The name was used after the very useful web app hammertime.djdavid98.art.*

- [**`github.com/zekroTJA/shinpuru/pkg/hashutil`**](pkg/hashutil)  
  *Package hashutil provides general utility functionalities to generate simple and fast hashes with salt and pepper.*

- [**`github.com/zekroTJA/shinpuru/pkg/httpreq`**](pkg/httpreq)  
  *Package httpreq provides general utilities for around net/http requests for a simpler API and extra utilities for parsing JSON request and response boddies.*

- [**`github.com/zekroTJA/shinpuru/pkg/intutil`**](pkg/intutil)  
  *Package intutil provides some utility functionalities for integers.*

- [**`github.com/zekroTJA/shinpuru/pkg/jdoodle`**](pkg/jdoodle)  
  *Package jdoodle provides an API wrapper for the jdoodle execute and credit-spent REST API.*

- [**`github.com/zekroTJA/shinpuru/pkg/lctimer`**](pkg/lctimer)  
  *Package lctimer provides a life cycle timer which calls registered callback handlers on timer elapse. This package is a huge buggy piece of crap, please don't use it. :)*

- [**`github.com/zekroTJA/shinpuru/pkg/limiter`**](pkg/limiter)  
  *Package limiter provides a fiber middleware for a bucket based request rate limiter.*

- [**`github.com/zekroTJA/shinpuru/pkg/mimefix`**](pkg/mimefix)  
  *Package mimefix provides functionalities to bypass this issue with fasthttp on windows hosts*: https://github.com/golang/go/issues/32350*

- [**`github.com/zekroTJA/shinpuru/pkg/mody`**](pkg/mody)  
  *Package mody allows to modify fields in an object.*

- [**`github.com/zekroTJA/shinpuru/pkg/msgcollector`**](pkg/msgcollector)  
  *Package msgcollector provides functionalities to collect messages in a channel in conect of a single command request.*

- [**`github.com/zekroTJA/shinpuru/pkg/multierror`**](pkg/multierror)  
  *Package multierror impements handling multiple errors as one error object.*

- [**`github.com/zekroTJA/shinpuru/pkg/onetimeauth`**](pkg/onetimeauth)  
  *Package onetimeout provides short duration valid JWT tokens which are only valid exactly once.*

- [**`github.com/zekroTJA/shinpuru/pkg/permissions`**](pkg/permissions)  
  *Package permissions provides functionalities to calculate, update and merge arrays of permission domain rules. Read this to get more information about how permission domains and rules are working: https://github.com/zekroTJA/shinpuru/wiki/Permissions-Guide*

- [**`github.com/zekroTJA/shinpuru/pkg/random`**](pkg/random)  
  *Package random provides some general purpose cryptographically pseudo-random utilities.*

- [**`github.com/zekroTJA/shinpuru/pkg/rediscmdstore`**](pkg/rediscmdstore)  
  *Package rediscmdstore provides an implementation of github.com/zekrotja/ken/store.CommandStore using a redis client to store the command cache.*

- [**`github.com/zekroTJA/shinpuru/pkg/roleutil`**](pkg/roleutil)  
  *Package roleutil provides general purpose utilities for discordgo.Role objects and arrays.*

- [**`github.com/zekroTJA/shinpuru/pkg/startuptime`**](pkg/startuptime)  
  *Package startuptime provides simple functionalities to measure the startup time of an application.*

- [**`github.com/zekroTJA/shinpuru/pkg/stringutil`**](pkg/stringutil)  
  *Package stringutil provides generl string utility functions.*

- [**`github.com/zekroTJA/shinpuru/pkg/thumbnail`**](pkg/thumbnail)  
  *Package thumbnail provides simple functionalities to generate thumbnails from images with a max witdh or height.*

- [**`github.com/zekroTJA/shinpuru/pkg/timerstack`**](pkg/timerstack)  
  *Package timerstack provides a timer which can execute multiple delayed functions one after one.*

- [**`github.com/zekroTJA/shinpuru/pkg/timeutil`**](pkg/timeutil)  
  *Package timeutil provides some general purpose functionalities around the time package.*

- [**`github.com/zekroTJA/shinpuru/pkg/twitchnotify`**](pkg/twitchnotify)  
  *Package twitchnotify provides functionalities to watch the state of twitch streams and notifying changes by polling the twitch REST API.*

- [**`github.com/zekroTJA/shinpuru/pkg/voidbuffer`**](pkg/voidbuffer)  
  *Package voidbuffer provides a simple, concurrency proof push buffer with a fixed size which "removes" firstly pushed values when fully filled.*

<!-- end:PUBLIC_PACKAGES -->

---

# Third party dependencies

### Back End

<!-- start:REQUIREMENTS -->
- [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo) `(v0.23.3-0.20210821175000-0fad116c6c2a)`
- [bwmarrin/snowflake](https://github.com/bwmarrin/snowflake) `(v0.3.0)`
- [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) `(v3.2.0+incompatible)`
- [esimov/stackblur-go](https://github.com/esimov/stackblur-go) `(v1.0.0)`
- [gabriel-vasile/mimetype](https://github.com/gabriel-vasile/mimetype) `(v1.4.0)`
- [generaltso/vibrant](https://github.com/generaltso/vibrant) `(v0.0.0-20200703055536-90f922bee78c)`
- [go-ping/ping](https://github.com/go-ping/ping) `(v0.0.0-20211014180314-6e2b003bffdd)`
- [redis/v8](https://github.com/go-redis/redis/v8) `(v8.11.4)`
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) `(v1.6.0)`
- [fiber/v2](https://github.com/gofiber/fiber/v2) `(v2.20.2)`
- [makeworld-the-better-one/go-isemoji](https://github.com/makeworld-the-better-one/go-isemoji) `(v1.2.0)`
- [minio/minio-go](https://github.com/minio/minio-go) `(v6.0.14+incompatible)`
- [narqo/go-badge](https://github.com/narqo/go-badge) `(v0.0.0-20210814192603-33684e887a6d)`
- [prometheus/client_golang](https://github.com/prometheus/client_golang) `(v1.11.0)`
- [qiangxue/fasthttp-routing](https://github.com/qiangxue/fasthttp-routing) `(v0.0.0-20160225050629-6ccdc2a18d87)`
- [ranna-go/ranna](https://github.com/ranna-go/ranna) `(v0.1.0)`
- [cron/v3](https://github.com/robfig/cron/v3) `(v3.0.1)`
- [sahilm/fuzzy](https://github.com/sahilm/fuzzy) `(v0.1.0)`
- [di/v2](https://github.com/sarulabs/di/v2) `(v2.4.2)`
- [sirupsen/logrus](https://github.com/sirupsen/logrus) `(v1.8.1)`
- [stretchr/testify](https://github.com/stretchr/testify) `(v1.7.0)`
- [traefik/paerser](https://github.com/traefik/paerser) `(v0.1.4)`
- [valyala/fasthttp](https://github.com/valyala/fasthttp) `(v1.31.0)`
- [wcharczuk/go-chart](https://github.com/wcharczuk/go-chart) `(v2.0.1+incompatible)`
- [zekroTJA/colorname](https://github.com/zekroTJA/colorname) `(v1.0.0)`
- [zekroTJA/ratelimit](https://github.com/zekroTJA/ratelimit) `(v1.0.0)`
- [zekroTJA/shireikan](https://github.com/zekroTJA/shireikan) `(v0.7.0)`
- [zekroTJA/timedmap](https://github.com/zekroTJA/timedmap) `(v1.4.0)`
- [zekrotja/dgrs](https://github.com/zekrotja/dgrs) `(v0.3.2)`
- [zekrotja/ken](https://github.com/zekrotja/ken) `(v0.10.0)`
- [x/image](https://golang.org/x/image) `(v0.0.0-20210628002857-a66eb6448b8d)`
- [x/sys](https://golang.org/x/sys) `(v0.0.0-20211015200801-69063c4bb744)`
- [x/time](https://golang.org/x/time) `(v0.0.0-20210723032227-1f47c861a9ac)`
<!-- end:REQUIREMENTS -->

### Web Front End

- [Angular 9](https://angular.io)
- [Bootstrap](https://ng-bootstrap.github.io)
- [dateformat](https://www.npmjs.com/package/dateformat)
- [core-js](https://www.npmjs.com/package/core-js)
- [rxjs](https://www.npmjs.com/package/rxjs)
- [tslib](https://www.npmjs.com/package/tslib)
- [zone.js](https://www.npmjs.com/package/zone.js)

### Assets

- Avatar used from album [御中元 魔法少女詰め合わせ](https://www.pixiv.net/member_illust.php?mode=medium&illust_id=44692506) made by [瑞希](https://www.pixiv.net/member.php?id=137253)
- Icons uded from [Material Icons Set](https://material.io/resources/icons/?style=baseline)
- Discord Icon used from [Discord's Branding Resources](https://discord.com/new/branding)

---

Copyright © 2018-2021 zekro Development (Ringo Hoffmann).  
Covered by MIT License.
