> Auto generated command manual | Thu, 24 Jan 2019 09:36:22 CET

# Command List

## CHAT
- [say](#say)
- [quote](#quote)
- [vote](#vote)

## ETC
- [user](#user)
- [id](#id)
- [bug](#bug)
- [stats](#stats)
- [test](#test)

## GENERAL
- [help](#help)
- [info](#info)
- [ment](#ment)

## GLOBAL ADMIN
- [game](#game)

## GUILD CONFIG
- [prefix](#prefix)
- [perms](#perms)
- [autorole](#autorole)
- [modlog](#modlog)
- [voicelog](#voicelog)

## MODERATION
- [clear](#clear)
- [mvall](#mvall)
- [report](#report)
- [kick](#kick)
- [ban](#ban)
- [mute](#mute)
- [notify](#notify)

# Command Details

## CHAT

### say

> send an embeded message with the bot

| | |
|---|---|
| Permission | 3 |
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
black, cyan, orange, violet, blue, green, yellow, white, red, pink

### quote

> quote a message from any chat

| | |
|---|---|
| Permission | 0 |
| Group | CHAT |
| Aliases | q |

**Usage**  
`quote <msgID/msgURL>`

### vote

> create and manage polls

| | |
|---|---|
| Permission | 0 |
| Group | CHAT |
| Aliases | poll |

**Usage**  
`vote <description> | <possibility1> | <possibility2> (| <possibility3> ...)` - create vote  
`vote close (<VoteID>|all)` - close your last vote, a vote by ID or all your open votes

## ETC

### user

> Get information about a user

| | |
|---|---|
| Permission | 0 |
| Group | ETC |
| Aliases | u, profile |

**Usage**  
`profile (<userResolvable>)` - get user info

### id

> Get the discord ID(s) by resolvable

| | |
|---|---|
| Permission | 0 |
| Group | ETC |
| Aliases | ids |

**Usage**  
`id (<resolvable>)`

### bug

> Get information how to submit a bug report or feature request

| | |
|---|---|
| Permission | 0 |
| Group | ETC |
| Aliases | bugreport, issue, suggestion |

**Usage**  
`bug`

### stats

> dispaly some stats like uptime or guilds/user count

| | |
|---|---|
| Permission | 0 |
| Group | ETC |
| Aliases | uptime, numbers |

**Usage**  
`stats`

### test

> just for testing purposes

| | |
|---|---|
| Permission | 999 |
| Group | ETC |
| Aliases |  |

**Usage**  


## GENERAL

### help

> dispaly list of command or get help for a specific command

| | |
|---|---|
| Permission | 0 |
| Group | GENERAL |
| Aliases | h, ?, man |

**Usage**  
`help` - display command list  
`help <command>` - display help of specific command

### info

> display some information about this bot

| | |
|---|---|
| Permission | 0 |
| Group | GENERAL |
| Aliases | information, description, credits, version, invite |

**Usage**  
`info`

### ment

> toggle the mentionability of a role

| | |
|---|---|
| Permission | 4 |
| Group | GENERAL |
| Aliases | mnt, mention, mentions |

**Usage**  
`ment` - display currently mentionable roles  
`ment <roleResolvable> (g)` - make role mentioanble until you mention the role in a message on the guild. By attaching the parameter `g`, the role will be mentionable until this command will be exeuted on the role again.

## GLOBAL ADMIN

### game

> set the presence of the bot

| | |
|---|---|
| Permission | 999 |
| Group | GLOBAL ADMIN |
| Aliases | presence, botmsg |

**Usage**  
`game msg <displayMessage>` - set the presence game text  
`game status <online|dnd|idle>` - set the status

## GUILD CONFIG

### prefix

> set a custom prefix for your guild

| | |
|---|---|
| Permission | 10 |
| Group | GUILD CONFIG |
| Aliases | pre, guildpre, guildprefix |

**Usage**  
`prefix` - display current guilds prefix  
`prefix <newPrefix>` - set the current guilds prefix

### perms

> Set the permission for specific groups on your server

| | |
|---|---|
| Permission | 10 |
| Group | GUILD CONFIG |
| Aliases | perm, permlvl, plvl |

**Usage**  
`perms` - get current permission settings  
`perms <LvL> <RoleResolvable> (<RoleResolvable> ...)` - set permission level for specific roles

### autorole

> set the autorole for the current guild

| | |
|---|---|
| Permission | 9 |
| Group | GUILD CONFIG |
| Aliases | arole |

**Usage**  
`autorole` - display currently set autorole  
`autorole <roleResolvable>` - set an auto role for the current guild

### modlog

> set the mod log channel for a guild

| | |
|---|---|
| Permission | 6 |
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
| Permission | 6 |
| Group | GUILD CONFIG |
| Aliases | setvoicelog, voicelogchan, vl |

**Usage**  
`voicelog` - set this channel as voicelog channel  
`voicelog <chanResolvable>` - set any text channel as voicelog channel  
`voicelog reset` - reset voice log channel

## MODERATION

### clear

> clear messages in a channel

| | |
|---|---|
| Permission | 8 |
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
| Permission | 5 |
| Group | MODERATION |
| Aliases | mva |

**Usage**  
`mvall <otherChanResolvable>`

### report

> report a user

| | |
|---|---|
| Permission | 5 |
| Group | MODERATION |
| Aliases | rep, warn |

**Usage**  
`report <userResolvable> [<type>] <reason>` - report a user *(if type is empty, its defaultly 0 = warn)*  
  
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
| Permission | 6 |
| Group | MODERATION |
| Aliases | userkick |

**Usage**  
`kick <UserResolvable> <Reason>`

### ban

> ban users with creating a report entry

| | |
|---|---|
| Permission | 8 |
| Group | MODERATION |
| Aliases | userban |

**Usage**  
`ban <UserResolvable> <Reason>`

### mute

> Mute members in text channels

| | |
|---|---|
| Permission | 4 |
| Group | MODERATION |
| Aliases | m, silence |

**Usage**  
`mute setup` - creates mute role and sets this role in every channel as muted  
`mute <userResolvable>` - mute/unmute a user  
`mute list` - display muted users on this guild

### notify

> get, remove or setup the notify rule

| | |
|---|---|
| Permission | 0 |
| Group | MODERATION |
| Aliases | n |

**Usage**  
`notify setup (<roleName>)` - creates the notify role and registers it for this command  
`notify` - get or remove the role

