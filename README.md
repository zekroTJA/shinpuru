<div align="center">
    <img src="https://zekro.de/src/shinpuru_avi_circle.png" height="300" />
    <h1>~ シンプル ~</h1>
    <strong>
        A simple multi purpose discord bot written in Go (discord.go)<br>
        with focus on stability and reliability
    </strong><br><br>
    <a href="https://dc.zekro.de"><img height="28" src="https://img.shields.io/discord/307084334198816769.svg?style=for-the-badge&logo=discord" /></a>&nbsp;
    <a href="https://github.com/zekroTJA/shinpuru/releases"><img height="28" src="https://img.shields.io/github/tag/zekroTJA/shinpuru.svg?style=for-the-badge"/></a>&nbsp;
    <a href="https://cloud.docker.com/u/zekro/repository/docker/zekro/shinpuru"><img alt="Docker Cloud Automated build" src="https://img.shields.io/docker/cloud/automated/zekro/shinpuru.svg?color=cyan&logo=docker&logoColor=cyan&style=for-the-badge"></a>&nbsp;
    <img height="28" src="https://forthebadge.com/images/badges/built-with-grammas-recipe.svg">
<br>
</div>

---

| Branch | Build |
|--------|-------|
| master | <a href="https://travis-ci.org/zekroTJA/shinpuru"><img src="https://travis-ci.org/zekroTJA/shinpuru.svg?branch=master" /></a> |
| dev | <a href="https://travis-ci.org/zekroTJA/shinpuru"><img src="https://travis-ci.org/zekroTJA/shinpuru.svg?branch=dev" /></a> |
| docker | <a href="https://cloud.docker.com/u/zekro/repository/docker/zekro/shinpuru"><img alt="Docker Cloud Build Status" src="https://img.shields.io/docker/cloud/build/zekro/shinpuru.svg"></a> |

---

# Invite

Here you can choose between the stable or canary version of shinpuru:

<a href="https://discordapp.com/api/oauth2/authorize?client_id=524847123875889153&scope=bot&permissions=2080894065"><img src="https://img.shields.io/badge/%20-INVITE%20STABLE-0288D1.svg?style=for-the-badge&logo=discord" height="30" /></a>

<a href="https://discordapp.com/api/oauth2/authorize?client_id=536916384026722314&scope=bot&permissions=2080894065"><img src="https://img.shields.io/badge/%20-INVITE%20CANARY-FFA726.svg?style=for-the-badge&logo=discord" height="30" /></a>

> **Attention**<br>The canary version runs on the latest build pushed to the dev branch and can contain bugs! Also, the canary version is running on a separate database which is not included in my daily database backup.

# Intro

シンプル (shinpuru), a simple *(as the name says)*, multi-purpose Discord Bot written in Go, using bwmarrin's package [discord.go](https://github.com/bwmarrin/discordgo) as API and gateway wrapper. The focus on this bot is not to punch in as many features and commands as possible, just some commands and features which I thought would be useful and which were the most used with my older Discord bots, like [zekroBot 2](https://github.com/zekroTJA/zekroBot2), and more on making this bot as reliable and stable as possible.

Also, I want to use this project as a chance for me, to get some deeper into Go and larger Go project structures. In a later development state, this bot will detach zekroBot 2.

---

# Features 

In this [**wiki article**](https://github.com/zekroTJA/shinpuru/wiki/Commands), you can find an automatically generated list of all commands and their manuals.

## Moderation

shinpuru brings general guild moderation features like clearing messages in text channels *(also user-specific, if required)*, reporting, muting, kicking and banning members. Those actions initiated with shinpurus moderation commands will be logged in a defined moderation text channel and in the database. So, all actions can be reviewed.

![](https://i.zekro.de/firefox_2019-02-22_14-54-59.png)
![](https://i.zekro.de/firefox_2019-02-22_14-57-37.png)

Also, there is a [`notify`](https://github.com/zekroTJA/shinpuru/wiki/Commands#notify) system, which creates a `@notify` role, which is as standard, not mentionable. Users can get or remove themselves this role by using the [`notify`](https://github.com/zekroTJA/shinpuru/wiki/Commands#notify) command. So, you can use this role as a replacement for `@everyone`, so it is like an "opt-in notification system".

You can combine that function with the [`ment`](https://github.com/zekroTJA/shinpuru/wiki/Commands#ment) command, which allows enabling or disabling mentionability of roles by command. If you enable the mentionability of `@notify`, for example, and after that, you mention this role in a message, the mentionability of this role will automatically be disabled.

Another feature is the [`autorole`](https://github.com/zekroTJA/shinpuru/wiki/Commands#autorole) system: You can specify a role, which will be added to every user joined the guild.

## Chat

Of course, there are some chats supporting commands like the [`say`](https://github.com/zekroTJA/shinpuru/wiki/Commands#say) command, where you can create embedded messages with the bot with custom colors, titles, footers, images, and so on. Also, it is possible to create embeds from raw json data (like documented in [Discords API docs](https://discordapp.com/developers/docs/resources/channel#embed-object)). For example, [here](https://github.com/dev-schueppchen/rules-and-docs/blob/master/embeds/welcome-msg.json) you can find the format of our development Discord guilds welcome message.

![](https://i.zekro.de/firefox_2019-02-22_15-16-46.png)

Another useful feature is the [`quote`](https://github.com/zekroTJA/shinpuru/wiki/Commands#quote) command, where you can quote messages from all text channels on a guild in any channel with *jump to* link. This can be generated by the ID or the URL of a message.

![](https://i.zekro.de/firefox_2019-02-22_15-19-32.png)

Time for some democracy? So, you cna create reaction-interactive votes with the [`vote`](https://github.com/zekroTJA/shinpuru/wiki/Commands#vote) command.

![](https://i.zekro.de/firefox_2019-02-22_15-22-03.png)

Annoyed from ghost pings *(messages with mentions, which were deleted, so you only see a mention but no message)*? shinpuru has a system for detecting those [`ghost pings`](https://github.com/zekroTJA/shinpuru/wiki/Commands#ghost) and punish people doing so by exposing the message which was deleted actually. You can also specify a format of how the warn message should look like, if you do not want to expose the message content or ping the victim again, for example.

![](https://i.zekro.de/firefox_2019-02-22_15-26-56.png)
![](https://i.zekro.de/firefox_2019-02-22_15-27-39.png)

## Guild Backups

You want to be prepared for each emergency? Just enable the auto-backup system of shinpuru with the [`backup`](https://github.com/zekroTJA/shinpuru/wiki/Commands#backup) command. Then, a full backup of all of the guilds roles, channels, members and guild settings will be saved every 12 hours. The last 10 backups are saved, so you have access to backups for the last 5 days. Of course, you can also automatically restore saved backups by using the `backup restore` command.

## Twitch Notifications

With the Twitch Notification System, you can stay up to date which channels are currently live on Twitch! Just enter the command [`!twitch <twitchUserName>`](https://github.com/zekroTJA/shinpuru/wiki/Commands#twitch) in a channel to set up the system. Then, every time, the channel goes live, a message will be posted to this channel, which will be automatically removed when the channel goes offline on Twitch.  
*Because of API limitations, the delay until the bot notifies a status change can be up to 3 minutes.*

![](https://i.zekro.de/firefox_2019-02-22_15-29-02.png)

## Code Execution

shinpuru is able to compile embedded code in messages on the fly, just by clicking a reaction under the message containing the code. The code will be sent to [jdoodle's](https://jdoodle.com) API, will be executed and the output will be displayed in the discord channel!

![](https://i.zekro.de/firefox_2019-02-22_15-36-36.png)

For setting up this system, use the [`exec setup`](https://github.com/zekroTJA/shinpuru/wiki/Commands#exec) command. Then, the bot will request your jdoodle's API credentials in DM *(because we don't want you to send your credentials into a public guilds text chat)*. Then, the system will be set up and enabled on your guild. Your credentials will only be used for your guild, so every guild is responsible for their credentials. That also means, if you have an advanced jdoodle plan, you can use this accounts credentials, of course, for your guild.

## Invite Link Blocking

By using the [`inv`](https://github.com/zekroTJA/shinpuru/wiki/Commands#inv) command, you can set up a guild-wide blocking for Discord Guild Invite Links. You can pass a minimum permission level users need to have to be allowed to send Invite Links. If users with a permission level below that, the messages including Invite Links will be deleted.

The system detects obvious invite links like `discord.gg/<InvID>` or `discordapp.com/invite/<InvID>`. Also, links which redirect to a Discord Invite link using the [`location` header](https://tools.ietf.org/html/rfc2616#section-14.30) or some sort of HTML redirection methods, like link shorteners do, for example, will be blocked.

## Voice Logging

Missing Teamspeaks voice activity log? Just specify a voice log channel with the [`voicelog`](https://github.com/zekroTJA/shinpuru/wiki/Commands#voicelog) command and every voice channel move will be logged in this channel.

![](https://i.zekro.de/firefox_2019-02-22_15-32-58.png)

---

# Docker

Read about how to self-host shinpuru using the provided Docker image in the [**wiki article**](https://github.com/zekroTJA/shinpuru/wiki/Docker).

---

# Compiling

Read about self-compiling in the [**wiki article**](https://github.com/zekroTJA/shinpuru/wiki/Self-Compiling).

---

# Third party dependencies

- [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)
- [go-yaml/yaml](https://github.com/go-yaml/yaml)
- [go-sql-driver/mysql](https://github.com/Go-SQL-Driver/MySQL/)
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- [op/go-logging](https://github.com/op/go-logging)
- [dayvonjersen/vibrant](https://github.com/dayvonjersen/vibrant)
- [bwmarrin/snowflake](https://github.com/bwmarrin/snowflake)
- [gopkg.in/yaml.v2](https://gopkg.in/yaml.v2)

Avatar of [御中元 魔法少女詰め合わせ](https://www.pixiv.net/member_illust.php?mode=medium&illust_id=44692506) from [瑞希](https://www.pixiv.net/member.php?id=137253).

---

Copyright © 2018-2019 zekro Development (Ringo Hoffmann).  
Covered by MIT Licence.
