> Auto generated command manual | Wed, 19 Feb 2020 15:59:50 CET

# Explicit Sub Commands

The commands below have sub command permissions which must be set explicitly and can not be applied by wildcards (`*`). So here you have them if you want to allow them for specific roles:

### vote

- **`sp.chat.vote.close`** - Allows closing votes also from other users
### inv

- **`sp.guild.mod.inviteblock.send`** - Allows sending invites even if invite block is enabled
### tag

- **`sp.chat.tag.create`** - Allows creating tags
- **`sp.chat.tag.delete`** - Allows deleting tags (of every user)

# Command List

## CHAT
- [say](#say)
- [quote](#quote)
- [vote](#vote)
- [user](#user)
- [twitch](#twitch)
- [exec](#exec)
- [tag](#tag)

## ETC
- [id](#id)
- [bug](#bug)
- [stats](#stats)
- [test](#test)

## GENERAL
- [help](#help)
- [info](#info)

## GLOBAL ADMIN
- [game](#game)

## GUILD ADMIN
- [backup](#backup)

## GUILD CONFIG
- [prefix](#prefix)
- [perms](#perms)
- [autorole](#autorole)
- [modlog](#modlog)
- [voicelog](#voicelog)
- [joinmsg](#joinmsg)
- [leavemsg](#leavemsg)

## MODERATION
- [clear](#clear)
- [mvall](#mvall)
- [report](#report)
- [kick](#kick)
- [ban](#ban)
- [mute](#mute)
- [ment](#ment)
- [notify](#notify)
- [ghost](#ghost)
- [inv](#inv)

# Command Details

## CHAT

### say

> send an embeded message with the bot

| | |
|---|---|
| Domain Name | sp.chat.say |
| Group | CHAT |
| Aliases | msg |

**Usage**  
`say [flags] <message>`  
  
**Flags:**   
```  
-c string  
	color (default "orange")  
-f string  
	footer  
-raw string  
	raw embed from json (see https://discordapp.com/developers/docs/resources/channel#embed-object)  
-t string  
	title  
```  
**Colors:**  
red, pink, blue, yellow, orange, black, violet, cyan, green, white

### quote

> quote a message from any chat

| | |
|---|---|
| Domain Name | sp.chat.quote |
| Group | CHAT |
| Aliases | q |

**Usage**  
`quote <msgID/msgURL>`

### vote

> create and manage polls

| | |
|---|---|
| Domain Name | sp.chat.vote |
| Group | CHAT |
| Aliases | poll |

**Usage**  
`vote <description> | <possibility1> | <possibility2> (| <possibility3> ...)` - create vote  
`vote list` - display currentltly running votes  
`vote expire <duration> (<voteID>)` - set expire to last created (or specified) vote  
`vote close (<VoteID>|all)` - close your last vote, a vote by ID or all your open votes

### user

> Get information about a user

| | |
|---|---|
| Domain Name | sp.chat.profile |
| Group | CHAT |
| Aliases | u, profile |

**Usage**  
`profile (<userResolvable>)` - get user info

### twitch

> Get notifications in channels when someone goes live on twitch

| | |
|---|---|
| Domain Name | sp.chat.twitch |
| Group | CHAT |
| Aliases | tn, twitchnotify |

**Usage**  
`twitch` - list all currently monitored twitch channels  
`twitch <twitchUsername>` - get notified in the current channel when the streamer goes online  
`twitch remove <twitchUsername>` - remove monitor

### exec

> setup code execution of code embeds

| | |
|---|---|
| Domain Name | sp.chat.exec |
| Group | CHAT |
| Aliases | ex, execute, jdoodle |

**Usage**  
`exec setup` - enter jdoodle setup  
`exec reset` - disable and delete token from database  


### tag

> set texts as tags which can be fastly re-posted later

| | |
|---|---|
| Domain Name | sp.chat.tag |
| Group | CHAT |
| Aliases | t, note |

**Usage**  
`tag` - Display all created tags on the current guild  
`tag create <identifier> <content>` - Create a tag  
`tag edit <identifier|ID> <content>` - Edit a tag  
`tag delete <identifier|ID>` - Delete a tag  
`tag raw <identifier|ID>` - Display tags content as raw markdown  
`tag <identifier|ID>` - Display tag

## ETC

### id

> Get the discord ID(s) by resolvable

| | |
|---|---|
| Domain Name | sp.etc.id |
| Group | ETC |
| Aliases | ids |

**Usage**  
`id (<resolvable>)`

### bug

> Get information how to submit a bug report or feature request

| | |
|---|---|
| Domain Name | sp.etc.bug |
| Group | ETC |
| Aliases | bugreport, issue, suggestion |

**Usage**  
`bug`

### stats

> display some stats like uptime or guilds/user count

| | |
|---|---|
| Domain Name | sp.etc.stats |
| Group | ETC |
| Aliases | uptime, numbers |

**Usage**  
`stats`

### test

> just for testing purposes

| | |
|---|---|
| Domain Name | sp.test |
| Group | ETC |
| Aliases |  |

**Usage**  


## GENERAL

### help

> display list of command or get help for a specific command

| | |
|---|---|
| Domain Name | sp.etc.help |
| Group | GENERAL |
| Aliases | h, ?, man |

**Usage**  
`help` - display command list  
`help <command>` - display help of specific command

### info

> display some information about this bot

| | |
|---|---|
| Domain Name | sp.etc.info |
| Group | GENERAL |
| Aliases | information, description, credits, version, invite |

**Usage**  
`info`

## GLOBAL ADMIN

### game

> set the presence of the bot

| | |
|---|---|
| Domain Name | sp.game |
| Group | GLOBAL ADMIN |
| Aliases | presence, botmsg |

**Usage**  
`game msg <displayMessage>` - set the presence game text  
`game status <online|dnd|idle>` - set the status

## GUILD ADMIN

### backup

> enable, disable and manage guild backups

| | |
|---|---|
| Domain Name | sp.guild.admin.backup |
| Group | GUILD ADMIN |
| Aliases | backups, bckp, guildbackup |

**Usage**  
`backup <enable|disable>` - enable or disable backups for your guild  
`backup (list)` - list all saved backups  
`backup restore <id>` - restore a backup

## GUILD CONFIG

### prefix

> set a custom prefix for your guild

| | |
|---|---|
| Domain Name | sp.guild.config.prefix |
| Group | GUILD CONFIG |
| Aliases | pre, guildpre, guildprefix |

**Usage**  
`prefix` - display current guilds prefix  
`prefix <newPrefix>` - set the current guilds prefix

### perms

> Set the permission for specific groups on your server

| | |
|---|---|
| Domain Name | sp.guild.config.perms |
| Group | GUILD CONFIG |
| Aliases | perm, permlvl, plvl |

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

### autorole

> set the autorole for the current guild

| | |
|---|---|
| Domain Name | sp.guild.config.autorole |
| Group | GUILD CONFIG |
| Aliases | arole |

**Usage**  
`autorole` - display currently set autorole  
`autorole <roleResolvable>` - set an auto role for the current guild  
`autorole reset` - disable autorole

### modlog

> set the mod log channel for a guild

| | |
|---|---|
| Domain Name | sp.guild.config.modlog |
| Group | GUILD CONFIG |
| Aliases | setmodlog, modlogchan, ml |

**Usage**  
`modlog` - set this channel as modlog channel  
`modlog <chanResolvable>` - set any text channel as mod log channel  
`modlog reset` - reset mod log channel

### voicelog

> set the mod log channel for a guild

| | |
|---|---|
| Domain Name | sp.guild.config.voicelog |
| Group | GUILD CONFIG |
| Aliases | setvoicelog, voicelogchan, vl |

**Usage**  
`voicelog` - set this channel as voicelog channel  
`voicelog <chanResolvable>` - set any text channel as voicelog channel  
`voicelog reset` - reset voice log channel

### joinmsg

> Set a message which will be sent into the defined channel when a member joins.

| | |
|---|---|
| Domain Name | sp.guild.config.joinmsg |
| Group | GUILD CONFIG |
| Aliases | joinmessage |

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

**Usage**  
`leavemsg msg <message>` - Set the message of the leave message.  
`leavemsg channel <ChannelIdentifier>` - Set the channel where the message will be sent into.  
`leavemsg reset` - Reset and disable leave messages.  
  
`[user]` will be replaced with the user name and `[ment]` will be replaced with the users mention when used in message text.

## MODERATION

### clear

> clear messages in a channel

| | |
|---|---|
| Domain Name | sp.guild.mod.clear |
| Group | MODERATION |
| Aliases | c, purge |

**Usage**  
`clear` - delete last message  
`clear <n>` - clear an ammount of messages  
`clear <n> <userResolvable>` - clear an ammount of messages by a specific user

### mvall

> move all members in your current voice channel into another one

| | |
|---|---|
| Domain Name | sp.guild.mod.mvall |
| Group | MODERATION |
| Aliases | mva |

**Usage**  
`mvall <otherChanResolvable>`

### report

> report a user

| | |
|---|---|
| Domain Name | sp.guild.mod.report |
| Group | MODERATION |
| Aliases | rep, warn |

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
Types `BAN`, `KICK` and `MUTE` are reserved for bands and kicks executed with this bot.

### kick

> kick users with creating a report entry

| | |
|---|---|
| Domain Name | sp.guild.mod.kick |
| Group | MODERATION |
| Aliases | userkick |

**Usage**  
`kick <UserResolvable> <Reason>`

### ban

> ban users with creating a report entry

| | |
|---|---|
| Domain Name | sp.guild.mod.ban |
| Group | MODERATION |
| Aliases | userban |

**Usage**  
`ban <UserResolvable> <Reason>`

### mute

> Mute members in text channels

| | |
|---|---|
| Domain Name | sp.guild.mod.mute |
| Group | MODERATION |
| Aliases | m, silence |

**Usage**  
`mute setup (<roleResolvable>)` - creates (or uses given) mute role and sets this role in every channel as muted  
`mute <userResolvable>` - mute/unmute a user  
`mute list` - display muted users on this guild  
`mute` - display currently set mute role

### ment

> toggle the mentionability of a role

| | |
|---|---|
| Domain Name | sp.guild.mod.ment |
| Group | MODERATION |
| Aliases | mnt, mention, mentions |

**Usage**  
`ment` - display currently mentionable roles  
`ment <roleResolvable> (g)` - make role mentioanble until you mention the role in a message on the guild. By attaching the parameter `g`, the role will be mentionable until this command will be exeuted on the role again.

### notify

> get, remove or setup the notify rule

| | |
|---|---|
| Domain Name | sp.chat.notify |
| Group | MODERATION |
| Aliases | n |

**Usage**  
`notify setup (<roleName>)` - creates the notify role and registers it for this command  
`notify` - get or remove the role

### ghost

> Send a message when someone ghost pinged a member

| | |
|---|---|
| Domain Name | sp.guild.mod.ghostping |
| Group | MODERATION |
| Aliases | gp, ghostping, gping |

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

> manage Discord invite blocking in chat

| | |
|---|---|
| Domain Name | sp.guild.mod.inviteblock |
| Group | MODERATION |
| Aliases | invblock |

**Usage**  
`inv enable` - enable invite link blocking  
`inv disable` - disable link blocking

