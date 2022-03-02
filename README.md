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

Most features of shinpuru available via slash commands are also accessible in the web interface with additional visualization and information provided.

https://user-images.githubusercontent.com/16734205/138589590-87301377-463d-43c3-8441-98ec84a1304c.mp4

## REST API

The web interface simply connects to the REST API exposed by the web server of shinpuru. You can also acquire an API token linked to your account to access the REST API directly, if you want.

[**Here**](https://github.com/zekroTJA/shinpuru/wiki/REST-API-Docs) you can read more about how to connect to shinpuru's REST API and which endpoints are available.

![](https://user-images.githubusercontent.com/16734205/138591104-a08890f8-52b8-44ee-b0fa-40123e3b84ba.png)


## Chat Features

### Code Execution

When someone posts code inside a code block, shinpuru can extract the code and language and execute it outputting the result into chat.

The code is picked up and sent to a code execution engine, which safely executes the code and sends back the result via a REST API. Therefore, you can chose between [ranna](https://github.com/ranna-go) or [JDoodle](https://www.jdoodle.com/) in the config.

![](https://user-images.githubusercontent.com/16734205/138688386-620119ac-659e-4903-8de8-5a6f0098666b.gif)

### Karma

shinpuru features a Karma system which is inspired by Reddit. You can define specific emotes which, when attached to a message, increase or reduce the karma points of a member. You can also specify the amount of "tokens" which can be spent each hour as well as a penalty for giving negative karma, which also takes karma from the executor to prevent downvote spam.

It is also possible to execute actions when passing specific amounts of karma. For example, you can add or remove roles, send messages or even kick/ban members depending on their karma points.

![](https://user-images.githubusercontent.com/16734205/138691018-385ef4a9-6997-46be-a8a0-880a1427d015.png)

### Color Reactions

Another unique feature of shinpuru are color reactions. When enabled (see [`/colorreaction`](https://github.com/zekroTJA/shinpuru/wiki/Commands#colorreaction)), shinpuru can fetch colors from chat messages and display them into a reaction. When clicked on the reaction, more information about the color is then posted into chat.

![](https://user-images.githubusercontent.com/16734205/138690308-457ac50b-3f3c-4782-82f9-c6c95a937efa.gif)

### Votes

You can simply create votes using the [`/vote`](https://github.com/zekroTJA/shinpuru/wiki/Commands#vote) slash command where users can [*pseudo anonymously*](https://github.com/zekroTJA/shinpuru/wiki/Why-are-Votes-%22pseudo-anonymous%22%3F) vote using reactions.

<img height="230px" src="https://user-images.githubusercontent.com/16734205/138737192-600c0385-74ce-44ab-bec8-32433b73d5ff.png" /><img height="230px" src="https://user-images.githubusercontent.com/16734205/138737253-001d3a50-5cfb-4f48-b5e6-eb7d169ef052.png" /><img height="230px" src="https://user-images.githubusercontent.com/16734205/138737282-9829833b-28a8-4338-8dad-0817e4d5669a.png" />

### Twitch Notifications

You can add names of twitch streamers to a watchlist (see [`/twitchnotify`](https://github.com/zekroTJA/shinpuru/wiki/Commands#twitchnotify)) and when they go live, a notification message in sent into the specified channel.

![](https://user-images.githubusercontent.com/16734205/138742312-202fe2de-b99d-4606-81c8-980748813939.png)

### Quote Messages

You can use the [`/quote`](https://github.com/zekroTJA/shinpuru/wiki/Commands#quote) command to quote messages by ID or link (even cross-channel).

![image](https://user-images.githubusercontent.com/16734205/138743500-cf16c25b-68c0-4d99-bb2c-a93d4619c8ac.png)

### Starboard

*As literally any other bot,* shinpuru also features a starboard! You can even specify an amount of karma members get when their message get into the starboard.

![](https://user-images.githubusercontent.com/16734205/138848039-771248f3-3f67-49a6-9256-3f14c4bb12fb.png)

### Channel Statistics

You can use the [`/channelstats`]([`/quote`](https://github.com/zekroTJA/shinpuru/wiki/Commands#channelstats)) command to analyze contribution statistics for specific text channels.

![](https://user-images.githubusercontent.com/16734205/138848776-a5a5446a-3e5f-4a3c-8b45-4ac832822c9f.png)

## Guild Security & Moderation

### Report System

shinpuru features a deeply integrated reporting and moderation system. You can create reports for members who violate guild rules which then are posted into a modlog channel (if specified). Also, all reports of a member can be viewed on their user profiles as well as in the web interface.

Of course, you can also kick and ban members with shinpuru, which also creates a report record in the modlog. It is even possible to create so called "ghost reports". It allows to report or ban members by ID which are no more part of the server.

![](https://user-images.githubusercontent.com/16734205/138853312-9bbfdb68-6875-41c4-b7ba-6febf27638f8.png)

![](https://user-images.githubusercontent.com/16734205/138854271-ecec133a-70eb-4d12-8105-87be93994138.gif)

When a member wants to request an unban, this can be done via the web interface when navigating to `<webAddress>/unbanme`.

https://user-images.githubusercontent.com/16734205/140642193-a89e90c5-f38d-40cd-82bb-d25b65aa3dc7.mp4

### Guild Backups

When enabled, shinpuru will create a backup of your guild's infrastructure every 12 hours. This includes guild settings, channels (names, positions and groups), roles (names, positions and permissions) and members (nicknames and applied roles).

When your guild gets raided or an admin goes rouge, you can simply choose one of the created backups and reset the guilds state.

The last 10 backups are stored and can be reviewed in the web interface.

![](https://user-images.githubusercontent.com/16734205/140642616-20dae0d7-d2d7-421a-9d41-5a717cb4ca78.png)

### Raid Alerting

This system allows you to set a threshold of new user ingress rate. When this rate exceeds, for example when a lot of (bot) accounts flush in to your guild (aka `raiding`), all admins of the guild will be alerted via DM. Also, the guilds moderation setting will be raised to `Highest` so that only users with roles or a valid phone number can chat.

![image](https://user-images.githubusercontent.com/16734205/140644018-9652d8c9-2716-43ae-bf5b-c1b2c17f895a.png)

![](https://user-images.githubusercontent.com/16734205/140643905-32c9e258-4971-4054-b99f-ec27c8fcd33a.png)

Additionally, all joined users after the event triggered are logged in a list which can be viewed in the web interface. You can also use this list to bulk kick or ban users captured in the antiraid join list.

![](https://user-images.githubusercontent.com/16734205/140643988-d1b857e5-8c62-4a3e-b0ba-409dc839f46e.png)

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

- [**`github.com/zekroTJA/shinpuru/pkg/giphy`**](pkg/giphy)  
  *Package giphy provides a crappy and inclomplete - but at least bloat free - Giphy API client.*

- [**`github.com/zekroTJA/shinpuru/pkg/hammertime`**](pkg/hammertime)  
  *Package hammertime provides functionailities to format a time.Time into a Discord timestamp mention. The name was used after the very useful web app hammertime.djdavid98.art.*

- [**`github.com/zekroTJA/shinpuru/pkg/hashutil`**](pkg/hashutil)  
  *Package hashutil provides general utility functionalities to generate simple and fast hashes with salt and pepper.*

- [**`github.com/zekroTJA/shinpuru/pkg/httpreq`**](pkg/httpreq)  
  *Package httpreq provides general utilities for around net/http requests for a simpler API and extra utilities for parsing JSON request and response boddies.*

- [**`github.com/zekroTJA/shinpuru/pkg/inline`**](pkg/inline)  
  *Package inline provides general inline operation functions like inline if or null coalescence.*

- [**`github.com/zekroTJA/shinpuru/pkg/intutil`**](pkg/intutil)  
  *Package intutil provides some utility functionalities for integers.*

- [**`github.com/zekroTJA/shinpuru/pkg/jdoodle`**](pkg/jdoodle)  
  *Package jdoodle provides an API wrapper for the jdoodle execute and credit-spent REST API.*

- [**`github.com/zekroTJA/shinpuru/pkg/lctimer`**](pkg/lctimer)  
  *Package lctimer provides a life cycle timer which calls registered callback handlers on timer elapse. This package is a huge buggy piece of crap, please don't use it. :)*

- [**`github.com/zekroTJA/shinpuru/pkg/limiter`**](pkg/limiter)  
  *Package limiter provides a fiber middleware for a bucket based request rate limiter.*

- [**`github.com/zekroTJA/shinpuru/pkg/logmsg`**](pkg/logmsg)  
  *No package description.*

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

- [**`github.com/zekroTJA/shinpuru/pkg/versioncheck`**](pkg/versioncheck)  
  *Package versioncheck provides endpoints to retrieve version information via different providers and utilities to compare versions.*

- [**`github.com/zekroTJA/shinpuru/pkg/voidbuffer`**](pkg/voidbuffer)  
  *Package voidbuffer provides a simple, concurrency proof push buffer with a fixed size which "removes" firstly pushed values when fully filled.*

<!-- end:PUBLIC_PACKAGES -->

---

# Third party dependencies

### Back End

<!-- start:REQUIREMENTS -->
- [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo) `(v0.23.3-0.20220216202327-6015eed9333e)`
- [bwmarrin/snowflake](https://github.com/bwmarrin/snowflake) `(v0.3.0)`
- [esimov/stackblur-go](https://github.com/esimov/stackblur-go) `(v1.0.2)`
- [gabriel-vasile/mimetype](https://github.com/gabriel-vasile/mimetype) `(v1.4.0)`
- [generaltso/vibrant](https://github.com/generaltso/vibrant) `(v0.0.0-20200703055536-90f922bee78c)`
- [go-ping/ping](https://github.com/go-ping/ping) `(v0.0.0-20211130115550-779d1e919534)`
- [redis/v8](https://github.com/go-redis/redis/v8) `(v8.11.4)`
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) `(v1.6.0)`
- [fiber/v2](https://github.com/gofiber/fiber/v2) `(v2.24.0)`
- [jwt/v4](https://github.com/golang-jwt/jwt/v4) `(v4.2.0)`
- [kataras/hcaptcha](https://github.com/kataras/hcaptcha) `(v0.0.0-20200711031247-2927d4faf32f)`
- [makeworld-the-better-one/go-isemoji](https://github.com/makeworld-the-better-one/go-isemoji) `(v1.3.0)`
- [minio/minio-go](https://github.com/minio/minio-go) `(v6.0.14+incompatible)`
- [narqo/go-badge](https://github.com/narqo/go-badge) `(v0.0.0-20210814192603-33684e887a6d)`
- [prometheus/client_golang](https://github.com/prometheus/client_golang) `(v1.11.0)`
- [qiangxue/fasthttp-routing](https://github.com/qiangxue/fasthttp-routing) `(v0.0.0-20160225050629-6ccdc2a18d87)`
- [ranna-go/ranna](https://github.com/ranna-go/ranna) `(v0.2.1-0.20220121125614-50d5fceefc7c)`
- [cron/v3](https://github.com/robfig/cron/v3) `(v3.0.1)`
- [sahilm/fuzzy](https://github.com/sahilm/fuzzy) `(v0.1.0)`
- [di/v2](https://github.com/sarulabs/di/v2) `(v2.4.2)`
- [sirupsen/logrus](https://github.com/sirupsen/logrus) `(v1.8.1)`
- [stretchr/testify](https://github.com/stretchr/testify) `(v1.7.0)`
- [traefik/paerser](https://github.com/traefik/paerser) `(v0.1.4)`
- [valyala/fasthttp](https://github.com/valyala/fasthttp) `(v1.32.0)`
- [wcharczuk/go-chart](https://github.com/wcharczuk/go-chart) `(v2.0.1+incompatible)`
- [zekroTJA/colorname](https://github.com/zekroTJA/colorname) `(v1.0.0)`
- [zekroTJA/ratelimit](https://github.com/zekroTJA/ratelimit) `(v1.0.0)`
- [zekroTJA/shireikan](https://github.com/zekroTJA/shireikan) `(v0.7.0)`
- [zekroTJA/timedmap](https://github.com/zekroTJA/timedmap) `(v1.4.0)`
- [zekrotja/dgrs](https://github.com/zekrotja/dgrs) `(v0.5.1)`
- [zekrotja/ken](https://github.com/zekrotja/ken) `(v0.12.1)`
- [zekrotja/sop](https://github.com/zekrotja/sop) `(v0.3.1)`
- [x/image](https://golang.org/x/image) `(v0.0.0-20211028202545-6944b10bf410)`
- [x/sys](https://golang.org/x/sys) `(v0.0.0-20220209214540-3681064d5158)`
- [x/time](https://golang.org/x/time) `(v0.0.0-20211116232009-f0f3c7e86c11)`
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
