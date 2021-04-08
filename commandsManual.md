> Auto generated command manual | Thu, 08 Apr 2021 17:06:46 CEST

# Explicit Sub Commands

The commands below have sub command permissions which must be set explicitly and can not be applied by wildcards (`*`). So here you have them if you want to allow them for specific roles:

**vote**

- **`sp.chat.vote.close`** - Allows closing votes also from other users

**notify**

- **`sp.chat.notify.setup`** - Allows setting up the notify role for this guild.

**exec**


**inv**

- **`sp.guild.mod.inviteblock.send`** - Allows sending invites even if invite block is enabled

**tag**

- **`sp.chat.tag.create`** - Allows creating tags
- **`sp.chat.tag.edit`** - Allows editing tags (of every user)
- **`sp.chat.tag.delete`** - Allows deleting tags (of every user)

**color**



# Command List

## CHAT
- [chanstats](#chanstats)
- [exec](#exec)
- [guild](#guild)
- [karma](#karma)
- [quote](#quote)
- [say](#say)
- [tag](#tag)
- [twitch](#twitch)
- [user](#user)
- [vote](#vote)

## ETC
- [bug](#bug)
- [id](#id)
- [login](#login)
- [snowflake](#snowflake)
- [stats](#stats)

## GENERAL
- [help](#help)
- [info](#info)

## GLOBAL ADMIN
- [game](#game)

## GUILD ADMIN
- [backup](#backup)

## GUILD CONFIG
- [autorole](#autorole)
- [color](#color)
- [joinmsg](#joinmsg)
- [leavemsg](#leavemsg)
- [modlog](#modlog)
- [perms](#perms)
- [prefix](#prefix)
- [starboard](#starboard)
- [voicelog](#voicelog)

## MODERATION
- [ban](#ban)
- [clear](#clear)
- [ghost](#ghost)
- [inv](#inv)
- [kick](#kick)
- [lock](#lock)
- [ment](#ment)
- [mute](#mute)
- [mvall](#mvall)
- [notify](#notify)
- [report](#report)

# Command Details

## CHAT

### chanstats

> Get channel contribution statistics.

| | |
|---|---|
| Domain Name | sp.chat.chanstats |
| Group | CHAT |
| Aliases | cstats |
| DM Capable | No |

**Usage**  
`chanstats (<ChannelIdentifier>) (limit:<nLimit>)` - get channel stats  
`chanstats msgs (<ChannelIdentifier>) (limit:<nLimit>)` - get channel stats by messages  
`chanstats att (<ChannelIdentifier>) (limit:<nLimit>)` - get channel stats by attachments

### exec

> Setup code execution of code embeds.

| | |
|---|---|
| Domain Name | sp.chat.exec |
| Group | CHAT |
| Aliases | ex, execute, jdoodle |
| DM Capable | No |

**Usage**  
`exec setup` - enter jdoodle setup  
`exec reset` - disable and delete token from database  
`exec check` - retrurns the number of tokens consumed this day  



**Sub Permission Rules**
- **`sp.chat.exec.exec`** - Allows activating a code execution in chat via reaction


### guild

> Outputs info about the current guild.

| | |
|---|---|
| Domain Name | sp.chat.guild |
| Group | CHAT |
| Aliases | guildinfo, g, gi |
| DM Capable | No |

**Usage**  
`guild` - Prints guild info

### karma

> Display users karma count or the guilds karma scoreboard.

| | |
|---|---|
| Domain Name | sp.chat.karma |
| Group | CHAT |
| Aliases | scoreboard, leaderboard, lb, sb, top |
| DM Capable | No |

**Usage**  
`karma` - Display karma scoreboard  
`karma <userResolvable>` - Display karma count of this user  


### quote

> Quote a message from any chat.

| | |
|---|---|
| Domain Name | sp.chat.quote |
| Group | CHAT |
| Aliases | q |
| DM Capable | No |

**Usage**  
`quote <msgID/msgURL> (<comment>)`

### say

> Send an embedded message with the bot.

| | |
|---|---|
| Domain Name | sp.chat.say |
| Group | CHAT |
| Aliases | msg |
| DM Capable | Yes |

**Usage**  
`say [flags] <message>`  
  
**Flags:**   
```  
-c string  
      color (default "orange")  
-e string  
      Message Link or [ChannelID/]MessageID of the message to be edited  
-f string  
      footer  
-raw  
      parses following content as raw embed from json (see https://discord.com/developers/docs/resources/channel#embed-object)  
-t string  
      title  
  
```  
**Colors:**  
red, pink, green, white, violet, blue, cyan, yellow, orange, black

### tag

> Set texts as tags which can be fastly re-posted later.

| | |
|---|---|
| Domain Name | sp.chat.tag |
| Group | CHAT |
| Aliases | t, note, tags |
| DM Capable | No |

**Usage**  
`tag` - Display all created tags on the current guild  
`tag create <identifier> <content>` - Create a tag  
`tag edit <identifier|ID> <content>` - Edit a tag  
`tag delete <identifier|ID>` - Delete a tag  
`tag raw <identifier|ID>` - Display tags content as raw markdown  
`tag <identifier|ID>` - Display tag


**Sub Permission Rules**
- **`sp.chat.tag.create`** `[EXPLICIT]` - Allows creating tags
- **`sp.chat.tag.edit`** `[EXPLICIT]` - Allows editing tags (of every user)
- **`sp.chat.tag.delete`** `[EXPLICIT]` - Allows deleting tags (of every user)


### twitch

> Get notifications in channels when someone goes live on Twitch.

| | |
|---|---|
| Domain Name | sp.chat.twitch |
| Group | CHAT |
| Aliases | tn, twitchnotify |
| DM Capable | No |

**Usage**  
`twitch` - list all currently monitored twitch channels  
`twitch <twitchUsername>` - get notified in the current channel when the streamer goes online  
`twitch remove <twitchUsername>` - remove monitor

### user

> Get information about a user.

| | |
|---|---|
| Domain Name | sp.chat.profile |
| Group | CHAT |
| Aliases | u, profile |
| DM Capable | No |

**Usage**  
`profile (<userResolvable>)` - get user info

### vote

> Create and manage polls.

| | |
|---|---|
| Domain Name | sp.chat.vote |
| Group | CHAT |
| Aliases | poll |
| DM Capable | No |

**Usage**  
`vote <description> | <possibility1> | <possibility2> (| <possibility3> ...)` - create vote  
`vote list` - display currentltly running votes  
`vote expire <duration> (<voteID>)` - set expire to last created (or specified) vote  
`vote close (<VoteID>|all) (nochart|nc)` - close your last vote, a vote by ID or all your open votes


**Sub Permission Rules**
- **`sp.chat.vote.close`** `[EXPLICIT]` - Allows closing votes also from other users


## ETC

### bug

> Get information how to submit a bug report or feature request.

| | |
|---|---|
| Domain Name | sp.etc.bug |
| Group | ETC |
| Aliases | bugreport, issue, suggestion |
| DM Capable | Yes |

**Usage**  
`bug`

### id

> Get the discord ID(s) by resolvable.

| | |
|---|---|
| Domain Name | sp.etc.id |
| Group | ETC |
| Aliases | ids |
| DM Capable | No |

**Usage**  
`id (<resolvable>)`

### login

> Get a link via DM to log into the shinpuru web interface.

| | |
|---|---|
| Domain Name | sp.etc.login |
| Group | ETC |
| Aliases | weblogin, token |
| DM Capable | Yes |

**Usage**  
`login`

### snowflake

> Calculate information about a Discord or Shinpuru snowflake.

| | |
|---|---|
| Domain Name | sp.etc.snowflake |
| Group | ETC |
| Aliases | sf |
| DM Capable | Yes |

**Usage**  
`snowflake <snowflake> (dc/sp)` - get snowflake information  
If you attach `dc` (Discord) or `sp` (shinpuru), you will force the calculation mode for the snowflake. With nothing given, the mode will be chosen automatically.

### stats

> Display some stats like uptime or guilds/user count.

| | |
|---|---|
| Domain Name | sp.etc.stats |
| Group | ETC |
| Aliases | uptime, numbers |
| DM Capable | Yes |

**Usage**  
`stats`

## GENERAL

### help

> Display list of command or get help for a specific command.

| | |
|---|---|
| Domain Name | sp.etc.help |
| Group | GENERAL |
| Aliases | h, ?, man |
| DM Capable | Yes |

**Usage**  
`help` - display command list  
`help <command>` - display help of specific command

### info

> Display some information about this bot.

| | |
|---|---|
| Domain Name | sp.etc.info |
| Group | GENERAL |
| Aliases | information, description, credits, version, invite |
| DM Capable | Yes |

**Usage**  
`info`

## GLOBAL ADMIN

### game

> Set the presence of the bot.

| | |
|---|---|
| Domain Name | sp.game |
| Group | GLOBAL ADMIN |
| Aliases | presence, botmsg |
| DM Capable | Yes |

**Usage**  
`game msg <displayMessage>` - set the presence game text  
`game status <online|dnd|idle>` - set the status

## GUILD ADMIN

### backup

> Enable, disable and manage guild backups.

| | |
|---|---|
| Domain Name | sp.guild.admin.backup |
| Group | GUILD ADMIN |
| Aliases | backups, bckp, guildbackup |
| DM Capable | No |

**Usage**  
`backup <enable|disable>` - enable or disable backups for your guild  
`backup (list)` - list all saved backups  
`backup restore <id>` - restore a backup  
`backup purge` - delete all backups of the guild

## GUILD CONFIG

### autorole

> Set the autorole for the current guild.

| | |
|---|---|
| Domain Name | sp.guild.config.autorole |
| Group | GUILD CONFIG |
| Aliases | arole |
| DM Capable | No |

**Usage**  
`autorole` - display currently set autorole  
`autorole <roleResolvable>` - set an auto role for the current guild  
`autorole reset` - disable autorole

### color

> Toggle color reactions enable or disable.

| | |
|---|---|
| Domain Name | sp.guild.config.color |
| Group | GUILD CONFIG |
| Aliases | clr, colorreaction |
| DM Capable | No |

**Usage**  
`color` - toggle enable or disable  
`color (enable|disable)` - set enabled or disabled


**Sub Permission Rules**
- **`sp.chat.colorreactions`** - Allows executing color reactions in chat by reaction


### joinmsg

> Set a message which will be sent into the defined channel when a member joins.

| | |
|---|---|
| Domain Name | sp.guild.config.joinmsg |
| Group | GUILD CONFIG |
| Aliases | joinmessage |
| DM Capable | No |

**Usage**  
`joinmsg msg <message>` - Set the message of the join message.  
`joinmsg channel <ChannelIdentifier>` - Set the channel where the message will be sent into.  
`joinmsg reset` - Reset and disable join messages.  
  
`[user]` will be replaced with the user name and `[ment]` will be replaced with the users mention when used in message text.

### leavemsg

> Set a message which will be sent into the defined channel when a member leaves.

| | |
|---|---|
| Domain Name | sp.guild.config.leavemsg |
| Group | GUILD CONFIG |
| Aliases | leavemessage |
| DM Capable | No |

**Usage**  
`leavemsg msg <message>` - Set the message of the leave message.  
`leavemsg channel <ChannelIdentifier>` - Set the channel where the message will be sent into.  
`leavemsg reset` - Reset and disable leave messages.  
  
`[user]` will be replaced with the user name and `[ment]` will be replaced with the users mention when used in message text.

### modlog

> Set the mod log channel for a guild.

| | |
|---|---|
| Domain Name | sp.guild.config.modlog |
| Group | GUILD CONFIG |
| Aliases | setmodlog, modlogchan, ml |
| DM Capable | No |

**Usage**  
`modlog` - set this channel as modlog channel  
`modlog <chanResolvable>` - set any text channel as mod log channel  
`modlog reset` - reset mod log channel

### perms

> Set the permission for specific groups on your server.

| | |
|---|---|
| Domain Name | sp.guild.config.perms |
| Group | GUILD CONFIG |
| Aliases | perm, permlvl, plvl |
| DM Capable | No |

**Usage**  
`perms` - get current permission settings  
`perms <PDNS> <RoleResolvable> (<RoleResolvable> ...)` - set permission for specific roles  
  
PDNS (permission domain name specifier) is used to define permissions to groups by domains. This specifier consists of two parts:  
The allow (`+`) / disallow (`-`) part and the domain name (`sp.guilds.config.*` for example).  
  
For example, if you want to allow all guild moderation commands for moderators use `+sp.guild.mod.*`. If you want to disallow a role to use a specific command like `sp!ban`, you can do this by disallowing the specific domain name `-sp.guild.mod.ban`.  
  
Keep in mind:  
`-` and `+` of the same domain always results in a disallow.  
Higher level rules (like `sp.guild.config.*`) always override lower level rules (like `sp.guild.*`).  
  
[**Here**](https://github.com/zekroTJA/shinpuru/blob/master/docs/permissions-guide.md) you can find further information about the permission system.

### prefix

> Set a custom prefix for your guild.

| | |
|---|---|
| Domain Name | sp.guild.config.prefix |
| Group | GUILD CONFIG |
| Aliases | pre, guildpre, guildprefix |
| DM Capable | No |

**Usage**  
`prefix` - display current guilds prefix  
`prefix <newPrefix>` - set the current guilds prefix

### starboard

> Set guild starboard settings.

| | |
|---|---|
| Domain Name | sp.guild.config.stats |
| Group | GUILD CONFIG |
| Aliases | star, stb |
| DM Capable | No |

**Usage**  
`starboard channel (<channelResolvable>)` - define a starboard channel  
`starboard threshold <int>` - define a threshold for reaction count  
`starboard emote <emoteName>` - define an emote to be used as starboard reaction  
`starboard karma <int>` - define the amount of karma gained  
`starboard disable` - disable starboard

### voicelog

> Set the mod log channel for a guild.

| | |
|---|---|
| Domain Name | sp.guild.config.voicelog |
| Group | GUILD CONFIG |
| Aliases | setvoicelog, voicelogchan, vl |
| DM Capable | No |

**Usage**  
`voicelog` - set this channel as voicelog channel  
`voicelog <chanResolvable>` - set any text channel as voicelog channel  
`voicelog reset` - reset voice log channel  
`voicelog ignore <chanResolvable>` - add voice channel to ignore list  
`voicelog unignore <chanResolvable> - removes a voice channel from the ignore list  
``voicelog ignorelist` - display ignored voice channels

## MODERATION

### ban

> Ban users with creating a report entry.

| | |
|---|---|
| Domain Name | sp.guild.mod.ban |
| Group | MODERATION |
| Aliases | userban |
| DM Capable | No |

**Usage**  
`ban <UserResolvable> <Reason>`

### clear

> Clear messages in a channel.

| | |
|---|---|
| Domain Name | sp.guild.mod.clear |
| Group | MODERATION |
| Aliases | c, purge |
| DM Capable | Yes |

**Usage**  
`clear` - delete last message  
`clear <n>` - clear an ammount of messages  
`clear <n> <userResolvable>` - clear an ammount of messages by a specific user

### ghost

> Send a message when someone ghost pinged a member.

| | |
|---|---|
| Domain Name | sp.guild.mod.ghostping |
| Group | MODERATION |
| Aliases | gp, ghostping, gping |
| DM Capable | No |

**Usage**  
`ghost` - display current ghost ping settings  
`ghost set (<msgPattern>)` - Set a ghost ping message pattern. If no 2nd argument is provided, the default pattern will be used.  
`ghost reset` - reset message and disable ghost ping warnings  
  
Usable variables in message pattern:  
- `{@pinger}` - mention of the user sent the ghost ping  
- `{pinger}` - username#discriminator of the user sent the ghost ping  
- `{@pinged}` - mention of the user got ghost pinged  
- `{pinged}` - username#discriminator of the user got ghost pinged  
- `{msg}` - the content of the message which ghost pinged  
  
Default message pattern:  
```  
{pinger} ghost pinged {pinged} with message:  
  
{msg}  
```

### inv

> Manage Discord invite blocking in chat.

| | |
|---|---|
| Domain Name | sp.guild.mod.inviteblock |
| Group | MODERATION |
| Aliases | invblock |
| DM Capable | No |

**Usage**  
`inv enable` - enable invite link blocking  
`inv disable` - disable link blocking


**Sub Permission Rules**
- **`sp.guild.mod.inviteblock.send`** `[EXPLICIT]` - Allows sending invites even if invite block is enabled


### kick

> Kick users with creating a report entry.

| | |
|---|---|
| Domain Name | sp.guild.mod.kick |
| Group | MODERATION |
| Aliases | userkick |
| DM Capable | No |

**Usage**  
`kick <UserResolvable> <Reason>`

### lock

> Locks the channel so that no one can write there anymore until unlocked.

| | |
|---|---|
| Domain Name | sp.guild.mod.lock |
| Group | MODERATION |
| Aliases | unlock, lockchan, unlockchan, readonly, ro, chatlock |
| DM Capable | No |

**Usage**  
`lock (<channelResolvable>)` - locks or unlocks either the current or the passed channel  


### ment

> Toggle the mentionability of a role.

| | |
|---|---|
| Domain Name | sp.guild.mod.ment |
| Group | MODERATION |
| Aliases | mnt, mention, mentions |
| DM Capable | No |

**Usage**  
`ment` - display currently mentionable roles  
`ment <roleResolvable> (g)` - make role mentioanble until you mention the role in a message on the guild. By attaching the parameter `g`, the role will be mentionable until this command will be exeuted on the role again.

### mute

> Mute members in text channels.

| | |
|---|---|
| Domain Name | sp.guild.mod.mute |
| Group | MODERATION |
| Aliases | m, silence, unmute, um, unsilence |
| DM Capable | No |

**Usage**  
`mute setup (<roleResolvable>)` - creates (or uses given) mute role and sets this role in every channel as muted  
`mute <userResolvable>` - mute/unmute a user  
`mute list` - display muted users on this guild  
`mute` - display currently set mute role

### mvall

> Move all members in your current voice channel into another one.

| | |
|---|---|
| Domain Name | sp.guild.mod.mvall |
| Group | MODERATION |
| Aliases | mva |
| DM Capable | No |

**Usage**  
`mvall <otherChanResolvable>`

### notify

> Get, remove or setup the notify rule.

| | |
|---|---|
| Domain Name | sp.chat.notify |
| Group | MODERATION |
| Aliases | n |
| DM Capable | No |

**Usage**  
`notify setup (<roleName>)` - creates the notify role and registers it for this command  
`notify` - get or remove the role


**Sub Permission Rules**
- **`sp.chat.notify.setup`** `[EXPLICIT]` - Allows setting up the notify role for this guild.


### report

> Report a user.

| | |
|---|---|
| Domain Name | sp.guild.mod.report |
| Group | MODERATION |
| Aliases | rep, warn |
| DM Capable | No |

**Usage**  
`report <userResolvable>` - list all reports of a user  
`report <userResolvable> [<type>] <reason>` - report a user *(if type is empty, its defaultly 0 = warn)*  
`report revoke <caseID> <reason>` - revoke a report  
  
**TYPES:**  
`0` - KICK  
`1` - BAN  
`2` - MUTE  
`3` - WARN  
`4` - AD  
Types `BAN`, `KICK` and `MUTE` are reserved for bans and kicks executed with this bot.

