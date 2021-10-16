> Auto generated command manual | Sun, 17 Oct 2021 01:16:14 CEST

# Command List

## CHAT

- [channelstats](#channelstats)
- [karma](#karma)
- [guild](#guild)
- [say](#say)
- [user](#user)
- [twitchnotify](#twitchnotify)
- [vote](#vote)
- [notify](#notify)
- [quote](#quote)
- [tag](#tag)

## GUILD CONFIG

- [starboard](#starboard)
- [modlog](#modlog)
- [voicelog](#voicelog)
- [announcements](#announcements)
- [autorole](#autorole)
- [colorreaction](#colorreaction)
- [perms](#perms)
- [exec](#exec)

## ETC

- [id](#id)
- [info](#info)
- [bug](#bug)
- [login](#login)
- [snowflake](#snowflake)
- [stats](#stats)

## 

- [presence](#presence)
- [maintenance](#maintenance)

## GUILD ADMIN

- [backup](#backup)

## GUILD MOD

- [clear](#clear)
- [mute](#mute)
- [ghostping](#ghostping)
- [inviteblock](#inviteblock)
- [moveall](#moveall)
- [report](#report)
- [lock](#lock)

# Command Details

## GUILD MOD

### clear

Clear messages in a channel.

| | |
|--|--|
| Domain Name | sp.guild.mod.clear |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Commands

##### `last`
Clears the last message

##### `amount`
Clear a specified amount of messages
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| amount | `Integer` | `true` | Amount of messages to clear |  || user | `Integer` | `false` | Clear messages send by this User |  |
##### `selected`
Removes either messages selected with ❌ emote by you or all messages below the 🔻 emote by you

### mute

Mute members or setup mute.

| | |
|--|--|
| Domain Name | sp.guild.mod.mute |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Commands

##### `toggle`
Toggle mute/unmute state of a member.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| user | `User` | `true` | The user to be muted/unmuted. |  || reason | `String` | `false` | The mute reason. |  || imageurl | `String` | `false` | Image attachment URL. |  || expire | `String` | `false` | Expiration time. |  |
##### `list`
List muted members.

##### `setup`
Setup mute role.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| role | `Role` | `false` | The role used to mute members (new one will be created if not specified). |  |
### ghostping

Setup the ghost ping system.

| | |
|--|--|
| Domain Name | sp.guild.mod.ghostping |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Commands

##### `status`
Display the current status of the ghost ping settings.

##### `setup`
Setup ghostping messages.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| message | `String` | `false` | The ghost ping message pattern. Use `/ghostping help` to get more info. |  |
##### `disable`
Disable ghostping messages.

##### `help`
Display help about the message pattern which can be used.

### inviteblock

Enable, disable or show state of invite blocking.

| | |
|--|--|
| Domain Name | sp.guild.mod.inviteblock |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Permission Rules

- **`sp.guild.mod.inviteblock.send`** `[EXPLICIT]` - Allows sending invites even if invite block is enabled


#### Arguments

| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| enable | `Boolean` | `false` | Set state to enabled or disabled. |  |### moveall

Move all members of the current voice channel to another one.

| | |
|--|--|
| Domain Name | sp.guild.mod.mvall |
| Version | 1.0.0 |
| DM Capable | false |


#### Arguments

| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| channel | `Channel` | `true` | Voice channel to move to. |  |### report

Create, revoke or list user reports.

| | |
|--|--|
| Domain Name | sp.guild.mod.report |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Permission Rules

- **`sp.guild.mod.report.list`** - List a users reports.
- **`sp.guild.mod.report.warn`** - Warn a member.
- **`sp.guild.mod.report.kick`** - Kick a member.
- **`sp.guild.mod.report.ban`** - Ban a member.
- **`sp.guild.mod.report.revoke`** - Revoke a report.


#### Sub Commands

##### `create`
File a new report.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| type | `Integer` | `true` | The type of report. | - `warn` (`3`)</br>- `ad` (`4`)</br>- `kick` (`0`)</br>- `ban` (`1`)</br> || user | `User` | `true` | The user. |  || reason | `String` | `true` | A short and concise report reason. |  || imageurl | `String` | `false` | An image url embedded into the report. |  || expire | `String` | `false` | Expire report after given time. |  |
##### `revoke`
Revoke a report.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| id | `Integer` | `true` | ID of the report to be revoked. |  || reason | `String` | `true` | Reason of the revoke. |  |
##### `list`
List the reports of a user.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| user | `User` | `true` | User to list reports of. |  |
### lock

Lock or unlock a channel so that no messages can be sent anymore.

| | |
|--|--|
| Domain Name | sp.guild.mod.lock |
| Version | 1.0.0 |
| DM Capable | false |


#### Arguments

| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| channel | `Channel` | `false` | The channel to be locked or unlocked (selects current channel if not passed). |  |
---
## CHAT

### channelstats

Get channel contribution statistics.

| | |
|--|--|
| Domain Name | sp.chat.chanstats |
| Version | 1.0.0 |
| DM Capable | false |


#### Arguments

| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| mode | `String` | `true` | The analysis mode. | - `messages` (`messages`)</br>- `attachments` (`attachments`)</br> || channel | `Channel` | `false` | The channel to be analyzed (defaultly current channel). |  || limit | `Integer` | `false` | The maximum amount of messages analyzed. |  |### karma

Display users karma count or the guilds karma scoreboard.

| | |
|--|--|
| Domain Name | sp.chat.karma |
| Version | 1.0.0 |
| DM Capable | false |


#### Arguments

| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| user | `User` | `false` | Display karma stats of a specific user. |  |### guild

Displays information about the current guild.

| | |
|--|--|
| Domain Name | sp.chat.guild |
| Version | 1.0.0 |
| DM Capable | false |


### say

Send an embedded message with the bot.

| | |
|--|--|
| Domain Name | sp.chat.say |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Commands

##### `embed`
Send an embed message.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| message | `String` | `true` | The message content. |  || color | `Integer` | `false` | The color. | - `default` (`16761095`)</br>- `cyan` (`48340`)</br>- `red` (`13840175`)</br>- `gray` (`11583173`)</br>- `green` (`9159498`)</br>- `lime` (`9159498`)</br>- `orange` (`16485376`)</br>- `violett` (`6953882`)</br>- `yellow` (`16761095`)</br> || title | `String` | `false` | The title content. |  || footer | `String` | `false` | The footer content. |  || channel | `Channel` | `false` | The channel to send the message into (or to edit a message in). |  || editmessage | `Integer` | `false` | The ID of the message to be edited. |  |
##### `raw`
Send raw embed message.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| json | `String` | `true` | The raw JSON data of the embed to be sent. |  || channel | `Channel` | `false` | The channel to send the message into (or to edit a message in). |  || editmessage | `Integer` | `false` | The ID of the message to be edited. |  |
### user

Get information about a user.

| | |
|--|--|
| Domain Name | sp.chat.profile |
| Version | 1.0.0 |
| DM Capable | false |


#### Arguments

| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| user | `User` | `true` | The user to be displayed. |  |### twitchnotify

Get notifications in channels when someone goes live on Twitch.

| | |
|--|--|
| Domain Name | sp.chat.twitch |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Commands

##### `list`
List alöl registered notifies for the guild.

##### `add`
Add a twitch user to be watched.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| twitchname | `String` | `true` | The username of the twitch user. |  || channel | `Channel` | `false` | The channel where the notifications are sent into (defaultly current channel). |  |
##### `remove`
Remove a twitch user from the watch list.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| twitchname | `String` | `true` | The username of the twitch user. |  |
### vote

Create and manage votes.

| | |
|--|--|
| Domain Name | sp.chat.vote |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Permission Rules

- **`sp.chat.vote.close`** `[EXPLICIT]` - Allows closing votes also from other users


#### Sub Commands

##### `create`
Create a new vote.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| body | `String` | `true` | The vote body content. |  || choises | `String` | `true` | The choises - split by `,`. |  || imageurl | `String` | `false` | An optional image URL. |  || channel | `Channel` | `false` | The channel to create the vote in (defaultly the current channel). |  || timeout | `String` | `false` | Timeout of the vote (i.e. `1h`, `30m`, ...) |  |
##### `list`
List currently running votes.

##### `expire`
Set the expiration of a running vote.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| id | `String` | `true` | The ID of the vote or `all` if you want to close all. |  || timeout | `String` | `true` | Timeout of the vote (i.e. `1h`, `30m`, ...) |  |
##### `close`
Close a running vote.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| id | `String` | `true` | The ID of the vote or `all` if you want to close all. |  || chart | `Boolean` | `false` | Display chart (default `true`). |  |
### notify

Get, remove or setup the notify role.

| | |
|--|--|
| Domain Name | sp.chat.notify |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Permission Rules

- **`sp.chat.notify.setup`** `[EXPLICIT]` - Allows setting up the notify role for this guild.


#### Sub Commands

##### `toggle`
Get or remove notify role.

##### `setup`
Setup notify role.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| role | `Role` | `false` | The role to be used as notify role (will be created if not specified). |  |
### quote

Quote a message from any chat.

| | |
|--|--|
| Domain Name | sp.chat.quote |
| Version | 1.0.0 |
| DM Capable | false |


#### Arguments

| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| id | `String` | `true` | The message ID or URL to be quoted. |  || comment | `String` | `false` | Add a comment directly to the quote. |  |### tag

Set texts as tags which can be fastly re-posted later.

| | |
|--|--|
| Domain Name | sp.chat.tag |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Permission Rules

- **`sp.chat.tag.create`** `[EXPLICIT]` - Allows creating tags
- **`sp.chat.tag.edit`** `[EXPLICIT]` - Allows editing tags (of every user)
- **`sp.chat.tag.delete`** `[EXPLICIT]` - Allows deleting tags (of every user)


#### Sub Commands

##### `show`
Show the content of a tag.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| name | `String` | `true` | The name of the Tag. |  |
##### `list`
List created tags.

##### `set`
Create or update a tag.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| name | `String` | `true` | The name of the Tag. |  || content | `String` | `true` | The content of the tag. You can use markdown as well as `\n` for line breaks. |  |
##### `delete`
Delete a tag.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| name | `String` | `true` | The name of the Tag. |  |
##### `raw`
Show a raw tag.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| name | `String` | `true` | The name of the Tag. |  |

---
## GUILD CONFIG

### starboard

Set guild starboard settings.

| | |
|--|--|
| Domain Name | sp.guild.config.starboard |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Commands

##### `set`
Set starboard settings.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| channel | `Channel` | `false` | The channel where the starboard messages will appear. |  || threshold | `Integer` | `false` | The minimum number of emote votes until a message gets into the starboard. |  || emote | `String` | `false` | The name or emote of the emote to be used for staring messages. |  || karma | `Integer` | `false` | The amount of karma gain when a users message gets into the starboard. |  |
##### `disable`
Disable the starboard.

### modlog

Set the mod log channel for a guild.

| | |
|--|--|
| Domain Name | sp.guild.config.modlog |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Commands

##### `set`
Set this or a specified channel as mod log channel.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| channel | `Channel` | `false` | A channel to be set as mod log (current channel if not specified). |  |
##### `disable`
Disable modlog.

### voicelog

Set the voice log channel for a guild.

| | |
|--|--|
| Domain Name | sp.guild.config.voicelog |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Commands

##### `set`
Set this or a specified channel as voice log channel.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| channel | `Channel` | `false` | A channel to be set as voice log (current channel if not specified). |  |
##### `disable`
Disable voicelog.

##### `ignore`
Add a voice channel to the ignorelist.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| channel | `Channel` | `true` | A voice channel to be ignored. |  |
##### `unignore`
Remove a voice channel from the ignorelist.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| channel | `Channel` | `true` | A voice channel to be unset from the ignore list. |  |
##### `ignorelist`
Show all ignored voice channels.

### announcements

Set a message which will show up when a user joins or leaves the guild.

| | |
|--|--|
| Domain Name | sp.guild.config.announcements |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Commands

##### `set`
Set a message and channel for an announcement message.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| type | `String` | `true` | The announcement type. | - `join` (`join`)</br>- `leave` (`leave`)</br> || message | `String` | `false` | The message. [user] will be replaced with the username and [ment] with the mention. |  || channel | `Channel` | `false` | A channel to be set. |  |
##### `disable`
Disable announcements.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| type | `String` | `true` | The announcement type. | - `join` (`join`)</br>- `leave` (`leave`)</br> |
### autorole

Manage guild autoroles.

| | |
|--|--|
| Domain Name | sp.guild.config.autorole |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Commands

##### `show`
Display the currently set autorole.

##### `add`
Add a role as autorole.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| role | `Role` | `true` | The autorole to be set. |  |
##### `remove`
Remove a role as autorole.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| role | `Role` | `true` | The autorole to be set. |  |
##### `purge`
Unset all autoroles.

### colorreaction

Toggle color reactions enable or disable.

| | |
|--|--|
| Domain Name | sp.guild.config.color |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Permission Rules

- **`sp.chat.colorreactions`** - Allows executing color reactions in chat by reaction


#### Arguments

| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| enable | `Boolean` | `false` | Set the enabled state of color reactions. |  |### perms

Set the permissions for groups on your guild.

| | |
|--|--|
| Domain Name | sp.guild.config.perms |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Commands

##### `list`
List the current permission definitions.

##### `set`
Set a permission rule for specific roles.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| mode | `String` | `true` | Set the permission as allow or disallow. | - `allow` (`+`)</br>- `disallow` (`-`)</br> || dns | `String` | `true` | Permission Domain Name Specifier |  || role | `Role` | `true` | The role to apply the permission to. |  || role2 | `Role` | `false` | Additional role to apply the permission to. |  || role3 | `Role` | `false` | Additional role to apply the permission to. |  || role4 | `Role` | `false` | Additional role to apply the permission to. |  || role5 | `Role` | `false` | Additional role to apply the permission to. |  |
##### `help`
Display help information for this command.

### exec

Setup code execution of code embeds.

| | |
|--|--|
| Domain Name | sp.guild.config.exec |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Commands

##### `setup`
Setup code execution.

##### `reset`
Disable code execution and remove stored credentials.

##### `check`
Show the status of the current code execution setup.


---
## ETC

### id

Get the discord ID(s) by resolvable.

| | |
|--|--|
| Domain Name | sp.etc.id |
| Version | 1.0.0 |
| DM Capable | false |


#### Arguments

| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| resolvable | `String` | `true` | The name of a discord object. |  |### info

Display some information about this bot.

| | |
|--|--|
| Domain Name | sp.etc.info |
| Version | 1.0.0 |
| DM Capable | true |


### bug

Get information how to submit a bug report or feature request.

| | |
|--|--|
| Domain Name | sp.etc.bug |
| Version | 1.0.0 |
| DM Capable | false |


### login

Receive a link via DM to log into the shinpuru web interface.

| | |
|--|--|
| Domain Name | sp.etc.Login |
| Version | 1.0.0 |
| DM Capable | true |


### snowflake

Calculate information about a Discord or Shinpuru snowflake.

| | |
|--|--|
| Domain Name | sp.etc.snowflake |
| Version | 1.0.0 |
| DM Capable | true |


#### Arguments

| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| snowflake | `Integer` | `true` | The snowflake ID. |  || type | `Integer` | `false` | The type of snowflake (will be determindes if not specified). | - `discord` (`0`)</br>- `shinpuru` (`1`)</br> |### stats

Display some stats like uptime or guilds/user count.

| | |
|--|--|
| Domain Name | sp.etc.stats |
| Version | 1.0.0 |
| DM Capable | true |



---
## 

### presence

Get information how to submit a bug report or feature request.

| | |
|--|--|
| Domain Name | sp.presence |
| Version | 1.0.0 |
| DM Capable | false |


#### Arguments

| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| message | `String` | `false` | The presence message. |  || status | `String` | `false` | The presence status. | - `online` (`online`)</br>- `idle` (`idle`)</br>- `dnd` (`dnd`)</br>- `invisible` (`invisible`)</br> |### maintenance

Maintenance utilities.

| | |
|--|--|
| Domain Name | sp.maintenance |
| Version | 1.0.0 |
| DM Capable | true |


#### Sub Commands

##### `flush-state`
Flush dgrs state.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| reconnect | `Boolean` | `false` | Disconnect and reconnect session after flush. |  || subkeys | `String` | `false` | The cache sub keys (comma seperated). |  |
##### `kill`
Kill the bot process.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| exitcode | `Integer` | `false` | The exit code. |  |
##### `reconnect`
Reconnects the Discord session.


---
## GUILD ADMIN

### backup

Manage guild backups.

| | |
|--|--|
| Domain Name | sp.guild.admin.backup |
| Version | 1.0.0 |
| DM Capable | false |


#### Sub Commands

##### `state`
Enable or disable the backup system.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| state | `Boolean` | `false` | Dispaly or set the backup state to enabled or disabled |  |
##### `list`
List all stored backups.

##### `restore`
Restore a backup.
**Arguments**
| Name | Type | Required | Description | Choises |
|------|------|----------|-------------|---------|
| index | `Integer` | `true` | The index of the backup to be restored. |  |
##### `purge`
Delete all stored backups.


---
