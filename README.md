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

<a href="https://discordapp.com/api/oauth2/authorize?client_id=524847123875889153&scope=bot&permissions=2080894065"><img src="https://img.shields.io/badge/%20-INVITE%20STABLE-0288D1.svg?style=for-the-badge&logo=discord" height="30" /></a>

<a href="https://discordapp.com/api/oauth2/authorize?client_id=536916384026722314&scope=bot&permissions=2080894065"><img src="https://img.shields.io/badge/%20-INVITE%20CANARY-FFA726.svg?style=for-the-badge&logo=discord" height="30" /></a>

# Intro

シンプル (shinpuru), a simple *(as the name says)*, multi-purpose Discord Bot written in Go, using bwmarrin's package [discord.go](https://github.com/bwmarrin/discordgo) as API and gateway wrapper and [shireikan](https://github.com/zekroTJA/shireikan) as command parser. The focus on this bot is to provide general purpose, administration and security tools while keeping stability, reliability and maintainability.

This bot is mainly used as administration and security tool on my [development discord](https://discord.zekro.de). Drop by to see shinpuru in action! 😉

---

# Features 

> In this [**wiki article**](https://github.com/zekroTJA/shinpuru/wiki/Commands), you can find an automatically generated list of all commands and their manuals.

Following, you will find a selected set of core features of shinpuru.

## Web Interface

shinpuru offers a web interface to view members profiles, reports, the guild mod log and also configure the guilds settings for shinpuru like mod log channel, voice log channel or join/leave messages and channels.

[**Demo**](https://i.imgur.com/hieLAua.gif)  
![](https://i.imgur.com/hieLAua.gif)

## Permission System

shinpuru has a fine grained and highly configurable permission system which uses "permission domains". You can specify permissions for whole groups of commands or for single commands for each role on your guild either by command or using the web interface.  
Please read [**this document**](https://github.com/zekroTJA/shinpuru/wiki/Permissions-Guide) about how the permission system exactly works and how to set it up correctly.

## Moderation

shinpuru brings general guild moderation features like clearing messages in text channels *(also user-specific, if required)*, reporting, muting, kicking and banning members. Those actions initiated with shinpurus moderation commands will be logged in a defined moderation text channel and in the database. So, all actions can be reviewed.

![](https://i.zekro.de/firefox_2019-02-22_14-54-59.png)
![](https://i.zekro.de/firefox_2019-02-22_14-57-37.png)

Also, there is a [`notify`](https://github.com/zekroTJA/shinpuru/wiki/Commands#notify) system, which creates a `@notify` role, which is as standard, not mentionable. Users can get or remove themselves this role by using the [`notify`](https://github.com/zekroTJA/shinpuru/wiki/Commands#notify) command. So, you can use this role as a replacement for `@everyone`, so it is like an "opt-in notification system".

You can combine that function with the [`ment`](https://github.com/zekroTJA/shinpuru/wiki/Commands#ment) command, which allows enabling or disabling mentionability of roles by command. If you enable the mentionability of `@notify`, for example, and after that, you mention this role in a message, the mentionability of this role will automatically be disabled.

Another feature is the [`autorole`](https://github.com/zekroTJA/shinpuru/wiki/Commands#autorole) system: You can specify a role, which will be added to every user joined the guild.

## Antiraid

![](https://i.imgur.com/vLMgrM9.png)  

Having trouble with Raids? shinpuru can help you with this by monitoring the rate of joining users and proceeding with security measures like increasing the servers security level.

You can also see the logs of the Antiraid to see if you got raided by Users. These logs can be downloaded or deleted in the web interface.

## Chat

Of course, there are some chats supporting commands like the [`say`](https://github.com/zekroTJA/shinpuru/wiki/Commands#say) command, where you can create embedded messages with the bot with custom colors, titles, footers, images, and so on. Also, it is possible to create embeds from raw json data (like documented in [Discords API docs](https://discordapp.com/developers/docs/resources/channel#embed-object)). For example, [here](https://github.com/dev-schueppchen/rules-and-docs/blob/master/embeds/welcome-msg.json) you can find the format of our development Discord guilds welcome message.

![](https://i.zekro.de/firefox_2019-02-22_15-16-46.png)

Another useful feature is the [`quote`](https://github.com/zekroTJA/shinpuru/wiki/Commands#quote) command, where you can quote messages from all text channels on a guild in any channel with *jump to* link. This can be generated by the ID or the URL of a message.

![](https://i.zekro.de/firefox_2019-02-22_15-19-32.png)

Time for some democracy? So, you can create reaction-interactive votes with the [`vote`](https://github.com/zekroTJA/shinpuru/wiki/Commands#vote) command.

![](https://i.zekro.de/firefox_2019-02-22_15-22-03.png)

Annoyed from ghost pings *(messages with mentions, which were deleted, so you only see a mention but no message)*? shinpuru has a system for detecting those [`ghost pings`](https://github.com/zekroTJA/shinpuru/wiki/Commands#ghost) and punish people doing so by exposing the message which was deleted actually. You can also specify a format of how the warn message should look like, if you do not want to expose the message content or ping the victim again, for example.

![](https://i.zekro.de/firefox_2019-02-22_15-26-56.png)
![](https://i.zekro.de/firefox_2019-02-22_15-27-39.png)

## Guild Backups

You want to be prepared for each emergency? Just enable the auto-backup system of shinpuru with the [`backup`](https://github.com/zekroTJA/shinpuru/wiki/Commands#backup) command. Then, a full backup of all of the guild roles, channels, members and guild settings will be saved every 12 hours. The last 10 backups are saved, so you have access to backups for the last 5 days. Of course, you can also automatically restore saved backups by using the `backup restore` command.

## Twitch Notifications

With the Twitch Notification System, you can stay up to date which channels are currently live on Twitch! Just enter the command [`!twitch <twitchUserName>`](https://github.com/zekroTJA/shinpuru/wiki/Commands#twitch) in a channel to set up the system. Then, every time, the streamer goes live, a message will be posted to this channel, which will be automatically removed when the channel goes offline on Twitch.  
*Because of API limitations, the delay until the bot notifies a status change can be up to one minute.*

![](https://i.zekro.de/firefox_2019-02-22_15-29-02.png)

## Code Execution

shinpuru is able to "compile" embedded code in messages on the fly, just by clicking a reaction under the message containing the code. The code will be sent to [jdoodle's](https://jdoodle.com) API, will be executed and the output will be displayed in the discord channel!

![](https://i.zekro.de/firefox_2019-02-22_15-36-36.png)

For setting up this system, use the [`exec setup`](https://github.com/zekroTJA/shinpuru/wiki/Commands#exec) command. Then, the bot will request your jdoodle's API credentials in DM *(because we don't want you to send your credentials into a public guilds text chat)*. Then, the system will be set up and enabled on your guild. Your credentials will only be used for your guild, so every guild is responsible for their credentials. That also means, if you have an advanced jdoodle plan, you can use this accounts credentials, of course, for your guild.

## Invite Link Blocking

By using the [`inv`](https://github.com/zekroTJA/shinpuru/wiki/Commands#inv) command, you can set up a guild-wide blocking for Discord Guild Invite Links. You can pass a minimum permission level users need to have to be allowed to send Invite Links. If users with a permission level below that, the messages including Invite Links will be deleted.

The system detects obvious invite links like `discord.gg/<InvID>` or `discordapp.com/invite/<InvID>`. Also, links which redirect to a Discord Invite link using the [`location` header](https://tools.ietf.org/html/rfc2616#section-14.30) or some sort of HTML redirection methods, like link shorteners do, for example, will be blocked.

## Voice Logging

Missing Teamspeak's voice activity log? Just specify a voice log channel with the [`voicelog`](https://github.com/zekroTJA/shinpuru/wiki/Commands#voicelog) command and every voice channel move will be logged in this channel.

![](https://i.zekro.de/firefox_2019-02-22_15-32-58.png)

---

# Docker

Read about how to host shinpuru using the provided Docker image in the [**wiki article**](https://github.com/zekroTJA/shinpuru/wiki/Docker).

---

# Compiling

Read about self-compiling in the [**wiki article**](https://github.com/zekroTJA/shinpuru/wiki/Self-Compiling).

---

# Public Packages

- [**`github.com/zekroTJA/shinpuru/pkg/acceptmsg`**](pkg/acceptmsg)  
  *Package acceptmsg provides a message model for discordgo which can be accepted or declined via message reactions.*

- [**`github.com/zekroTJA/shinpuru/pkg/angularservice`**](pkg/angularservice)  
  *Package angularservice provides bindings to start an Angular development server via the Angular CLI.*

- [**`github.com/zekroTJA/shinpuru/pkg/boolutil`**](pkg/boolutil)  
  *Package boolutil provides simple utility functions around booleans.*

- [**`github.com/zekroTJA/shinpuru/pkg/bytecount`**](pkg/bytecount)  
  *Package bytecount provides functionalities to format byte counts.*

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

- [**`github.com/zekroTJA/shinpuru/pkg/httpreq`**](pkg/httpreq)  
  *Package httpreq provides general utilities for around net/http requests for a simpler API and extra utilities for parsing JSON request and response boddies.*

- [**`github.com/zekroTJA/shinpuru/pkg/jdoodle`**](pkg/jdoodle)  
  *Package jdoodle provides an API wrapper for the jdoodle execute and credit-spent REST API.*

- [**`github.com/zekroTJA/shinpuru/pkg/lctimer`**](pkg/lctimer)  
  *Package lctimer provides a life cycle timer which calls registered callback handlers on timer elapse.*

- [**`github.com/zekroTJA/shinpuru/pkg/mimefix`**](pkg/mimefix)  
  *Package mimefix provides functionalities to bypass this issue with fasthttp on windows hosts*: https://github.com/golang/go/issues/32350*

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

- [**`github.com/zekroTJA/shinpuru/pkg/roleutil`**](pkg/roleutil)  
  *Package roleutil provides general purpose utilities for discordgo.Role objects and arrays.*

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

---

# Third party dependencies

### Back End

- [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)
- [bwmarrin/snowflake](https://github.com/bwmarrin/snowflake)
- [gabriel-vasile/mimetype](https://github.com/gabriel-vasile/mimetype)
- [dayvonjersen/vibrant](https://github.com/dayvonjersen/vibrant)
- [go-redis/redis](https://github.com/go-redis/redis)
- [go-sql-driver/mysql](https://github.com/Go-SQL-Driver/MySQL/)
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- [minio/minio-go](https://github.com/minio/minio-go)
- [op/go-logging](https://github.com/op/go-logging)
- [valyala/fasthttp ](https://github.com/valyala/fasthttp)
- [wcharczuk/go-chart](https://github.com/wcharczuk/go-chart)
- [zekroTJA/colorname](https://github.com/zekroTJA/colorname)
- [zekroTJA/ratelimit](https://github.com/zekroTJA/ratelimit)
- [zekroTJA/shireikan](https://github.com/zekroTJA/shireikan)
- [zekroTJA/timedmap](https://github.com/zekroTJA/timedmap)
- [gopkg.in/yaml.v2](https://gopkg.in/yaml.v2)

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
