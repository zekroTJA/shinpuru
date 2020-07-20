# REST API Docs

When enabled, shinpuru exposes a RESTful HTTP(S) API which exposes all functionalities which are also available to the web frontend.

## Authentication

All requests to the API needs to be authenticated and authorized. To authenticate your requests, you need to generate an API token in shinpurus web interface.

![](https://i.zekro.de/brave_bUcxerLK1C.png)
![](https://i.zekro.de/brave_dnOB5DeuFy.png)
![](https://i.zekro.de/brave_Mt6HBD4PLe.png)

To authenticate your requests, you need to add an `Authentication` header to your request with the token as `Bearer` token.

```
> GET /api/me HTTP/1.1
> Host: sp.zekro.de
> Authorization: bearer eyJhbGciOiJIUzI1...
> Accept: */*
```

## Objects

The following are API models of objects returned from the API.

### List Response

Requests which produce a list as response are wrapped in the following model:

| Field | Type | Description |
|-------|------|-------------|
| `n` | `int` | Number of items in the list. |
| `data` | `object[]` | The list of items. |

Example: 
```json
{
  "n": 3,
  "data": [
    { ... },
    { ... },
    { ... }
  ]
}
```

### User

A Discord User object.

> The user objects has some more fields than listed below
> comming from the discordgo.User object which may not 
> contain valid data.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `string` | The snowflake ID of the user. |
| `username` | `string` | The username of the user. |
| `avatar` | `string` | The avatar hash of the user. |
| `discriminator` | `string` | The discriminator of the user. |
| `bot` | `boolean` | Whether the user is a bot. |
| `avatar_url` | `string` | Public url of the avatar image file. |
| `created_at` | `timestamp` | Timestamp of user account creation. |
| `bot_owner` | `boolean` | Whether the user is the specified bot owner. |

Example:
```json
{
  "id": "221905671296253953",
  "username": "zekro",
  "avatar": "a_752a15d01e68fb5f6f6ec83400461a6a",
  "discriminator": "0001",
  "bot": false,
  "avatar_url": "https://cdn.discordapp.com/avatars/221905671296253953/a_752a15d01e68fb5f6f6ec83400461a6a.gif",
  "created_at": "2016-09-04T08:38:26.976834845Z",
  "bot_owner": true
}
```

### Member

A Discord Guild Member object.

> The member objects has some more fields than listed below
> comming from the discordgo.Member object which may not 
> contain valid data.

| Field | Type | Description |
|-------|------|-------------|
| `guild_id` | `string` | Snowflake ID of the Guild. |
| `joined_at` | `timestamp` | The timestamp when the member has joined the guild. |
| `nick` | `string` | The nick name of the user on this guild. |
| `deaf` | `boolean` | Whether the member is deafed on the guild. |
| `mute` | `boolean` | Whether the member is muted on the guild. |
| `user` | `User` | User model of the member. |
| `roles` | `string[]` | Role IDs of the member. |
| `premium_since` | `timestamp` | Timestamp since member has started boosting the server. |
| `avatar_url` | `string` | Public url of the avatar image file. |
| `created_at` | `timestamp` | Timestamp of user account creation. |
| `dominance` | `int` | The permission dominance of the member:<br>`1` - Guild Admin<br>`2` - Guild Owner<br>`3` - Bot Owner |

Example: 
```json
{
  "guild_id": "362162947738566657",
  "joined_at": "2020-04-09T20:53:47.658000+00:00",
  "nick": "zekuro senpai",
  "deaf": false,
  "mute": false,
  "user": {
    "id": "221905671296253953",
    "username": "zekro",
    "avatar": "a_752a15d01e68fb5f6f6ec83400461a6a",
    "discriminator": "0001",
    "bot": false
  },
  "roles": [
    "362166557721362433",
    "362169804146081802"
  ],
  "premium_since": "2020-04-15T09:24:24.174000+00:00",
  "avatar_url": "https://cdn.discordapp.com/avatars/221905671296253953/a_752a15d01e68fb5f6f6ec83400461a6a.gif",
  "created_at": "2016-09-04T08:38:26.976834845Z",
  "dominance": 1
}
```

### Role

A Discord Guild Role object.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `string` | The snowflake ID of the role. |
| `name` | `string` | The name of the role. |
| `managed` | `boolean` | Whether the role is managed. |
| `mentionable` | `boolean` | Whether the role is mentionable. |
| `hoist` | `boolean` | Whether the role is hoisted. |
| `color` | `int` | The color value of the role. |
| `position` | `int` | The position of the role. |
| `permission` | `int` | The permissions flags of the role. |

Example:
```json
{
  "id": "362169804146081802",
  "name": "Pleb",
  "managed": false,
  "mentionable": true,
  "hoist": true,
  "color": 0,
  "position": 8,
  "permissions": 104193600
}
```

### Channel

A Discord Channel object.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `string` | The snowflake ID of the channel. |
| `guild_id` | `string` | The snowflake ID of the guild the channel belongs to. |
| `name` | `string` | The name of the channel. |
| `topic` | `string` | The topic of the channel. |
| `type` | `int` | The type of the channel:<br>`0` - text channel<br>`2` - voice channel<br>`4` - category<br>`5` - news channel<br>`6` - store channel |
| `nfsw` | `boolean` | Whether the channel is specified as NFSW. |
| `icon` | `string` | The icon hash of the channel. |
| `position` | `int` | The position of the channel. |
| `bitrate` | `int` | The bitrate of the channel *(only for voice channels)*. |
| `permission_overwrites` | `PermissionOverwrite[]` | List of permission overwrites. |
| `user_limit` | `int` | The user limit of the channel. |
| `parent_id` | `string` | The ID of an optional parent category channel. |

Example:
```json
{
  "id": "526073401794756619",
  "guild_id": "362162947738566657",
  "name": "Gaming Private",
  "topic": "",
  "type": 2,
  "nsfw": false,
  "icon": "",
  "position": 15,
  "bitrate": 96000,
  "permission_overwrites": [
    {
      "id": "362162947738566657",
      "type": "role",
      "deny": 1049600,
      "allow": 2097152
    },
    {
      "id": "362166741373288448",
      "type": "role",
      "deny": 0,
      "allow": 1049600
    }
  ],
  "user_limit": 0,
  "parent_id": "384716117069004812"
}
```

### Guild 

A Discord Guild object.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `string` | The snowflake ID of the guild. |
| `name` | `string` | The name of the guild. |
| `icon` | `string` | The icon hash of the guild. |
| `region` | `string` | The region of the guild. |
| `afk_channel_id` | `string` | The specified AFK channel of the guild. |
| `owner_id` | `string` | The snowflake ID of the owner of the guild. |
| `joined_at` | `string` | The timestamp the bot user has joined the guild. |
| `splash` | `string` | The splash of the guild. |
| `member_count` | `int` | The ammount of members on the guild. |
| `verification_level` | `int` | The required verification level of the guild. |
| `large` | `boolean` | Whether the guild is large. |
| `unavaliable` | `boolean` | Whether the guild is currently unavailable due to outage. |
| `mfa_enabled` | `boolean` | Whether the guild has MFA enabled for admins. |
| `description` | `string` | The description of the guild. |
| `banner` | `string` | The hash of the banner image of the guild. |
| `premium_tier` | `int` | The premium tier of the guild. |
| `premium_subscription_ammount` | `int` | The number of boosts the guild has. |
| `roles` | `Role[]` | List of roles of the guild. |
| `channels` | `Channel[]` | List of channels of the guild. |
| `self_member` | `Member` | The member object of the authenticated user on the guild. |
| `icon_url` | `string` | The resource URL of the guilds icon. |
| `backups_enabled` | `bool` | Whether backup generation is enabled on this guild or not. |
| `latest_backup_entry` | `timestamp` | Time of the latest backup created. |

Example:
```json
{
  "id": "362162947738566657",
  "name": "zekro's Privatbutze",
  "icon": "2bdf517d77a79b1d6ba60457bd00128e",
  "region": "europe",
  "afk_channel_id": "384315833104597005",
  "owner_id": "221905671296253953",
  "joined_at": "2019-01-21T18:59:09.405000+00:00",
  "splash": "",
  "member_count": 41,
  "verification_level": 4,
  "embed_enabled": false,
  "large": false,
  "unavailable": false,
  "mfa_level": 0,
  "description": "",
  "banner": "",
  "premium_tier": 1,
  "premium_subscription_count": 3,
  "roles": [
    {
      "id": "362162947738566657",
      "name": "@everyone",
      "managed": false,
      "mentionable": false,
      "hoist": false,
      "color": 0,
      "position": 0,
      "permissions": 37084224
    }
  ],
  "channels": [
    {
      "id": "596457051928920134",
      "guild_id": "362162947738566657",
      "name": "tft-stuff",
      "topic": "",
      "type": 0,
      "last_message_id": "598270124067127296",
      "last_pin_timestamp": "",
      "nsfw": false,
      "icon": "",
      "position": 14,
      "bitrate": 0,
      "recipients": null,
      "permission_overwrites": [
        {
          "id": "362162947738566657",
          "type": "role",
          "deny": 1024,
          "allow": 0
        }
      ],
      "user_limit": 0,
      "parent_id": "676181576249245697",
      "rate_limit_per_user": 0
    }
  ],
  "self_member": {
    "guild_id": "",
    "joined_at": "2020-04-09T20:53:47.658000+00:00",
    "nick": "",
    "deaf": false,
    "mute": false,
    "user": {
      "id": "221905671296253953",
      "email": "",
      "username": "zekro",
      "avatar": "a_752a15d01e68fb5f6f6ec83400461a6a",
      "locale": "",
      "discriminator": "0001",
      "token": "",
      "verified": false,
      "mfa_enabled": false,
      "bot": false
    },
    "roles": [
      "362166557721362433"
    ],
    "premium_since": "2020-04-15T09:24:24.174000+00:00",
    "avatar_url": "https://cdn.discordapp.com/avatars/221905671296253953/a_752a15d01e68fb5f6f6ec83400461a6a.gif",
    "created_at": "2016-09-04T08:38:26.976834845Z",
    "dominance": 1
  },
  "icon_url": "https://cdn.discordapp.com/icons/362162947738566657/2bdf517d77a79b1d6ba60457bd00128e.png",
  "backups_enabled": true,
  "latest_backup_entry": "2020-07-20T17:59:59+02:00"
}
```

### GuildReduced 

A Discord Guild object reduced to fewer necessary fields.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `string` | The snowflake ID of the guild. |
| `name` | `string` | The name of the guild. |
| `icon` | `string` | The icon hash of the guild. |
| `region` | `string` | The region of the guild. |
| `owner_id` | `string` | The snowflake ID of the owner of the guild. |
| `joined_at` | `string` | The timestamp the bot user has joined the guild. |
| `member_count` | `int` | The ammount of members on the guild. |
| `icon_url` | `string` | The resource URL of the guilds icon. |

Example:
```json
{
  "id": "362162947738566657",
  "name": "zekro's Privatbutze",
  "icon": "2bdf517d77a79b1d6ba60457bd00128e",
  "icon_url": "https://cdn.discordapp.com/icons/362162947738566657/2bdf517d77a79b1d6ba60457bd00128e.png",
  "region": "europe",
  "owner_id": "221905671296253953",
  "joined_at": "2019-01-21T18:59:09.405000+00:00",
  "member_count": 41
}
```

### Report

A shinpuru Report record object.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `string` | The snowflake ID of the report. |
| `type` | `int` | The type of the report:<br>`0` - kick<br>`1` - ban<br>`2` - mute<br>`4` - warn<br>`5` - advertisement |
| `guild_id` | `string` | The snowflake ID of the guild the report belongs to. |
| `executor_id` | `string` | The snowflake ID of the executor of the report. |
| `victim_id` | `string` | The snowflake ID of the victim of the report. |
| `message` | `string` | The message of the report. |
| `attachment_url` | `string` | The resource URL of an optional attachment. |
| `type_name` | `string` | The name of the type of the report. |
| `created` | `timestamp` | The creation timestamp of the report. |

Example:
```json
{
  "id": "6678266303259619328",
  "type": 3,
  "guild_id": "547762913876639754",
  "executor_id": "221905671296253953",
  "victim_id": "455819141245304832",
  "message": "Bad language",
  "attachment_url": "https://sp-canary.zekro.de/imagestore/6678266279931420672.png",
  "type_name": "WARN",
  "created": "2020-07-03T09:29:57Z"
}
```

### GuildSettings

A wrapper object for all guild specific preferences and settings.

| Field | Type | Description |
|-------|------|-------------|
| `prefix` | `string` | The guild specific prefix. |
| `perms` | `{ string: string[] }` | A map of role-sepcific permission rules. |
| `autorole` | `string` | The snowflake ID of the set autorole. |
| `modlogchannel` | `string` | The snowflake ID of the set modlog channel. |
| `voicelogchannel` | `string` | The snowflake ID of the set voicelog channel. |
| `joinmessagechannel` | `string` | The snowflake ID of the channel where join messages are sent to. |
| `joinmessagetext` | `string` | The text which is sent into the join message channel when a user joins the guild. |
| `leavemessagechannel` | `string` | The snowflake ID of the channel where leave messages are sent to. |
| `leavemessagetext` | `string` | The text which is sent into the leave message channel when a user leaves the guild. |

Example:
```json
{
  "prefix": "!",
  "perms": {
    "362166741373288448": [
      "+sp.guild.mod.*",
      "-sp.guild.mod.ban",
      "-sp.guild.mod.kick"
    ],
    "406891236407115777": [
      "+sp.guild.config.*"
    ]
  },
  "autorole": "362169804146081802",
  "modlogchannel": "529279471350710314",
  "voicelogchannel": "454618414258978839",
  "joinmessagechannel": "381119295632965655",
  "joinmessagetext": "Welcome, [user]!",
  "leavemessagechannel": "",
  "leavemessagetext": ""
}
```

### GlobalInviteSettings

An object which specifies the global invite settings.

| Field | Type | Description |
|-------|------|-------------|
| `invite_url` | `string` | The invite URL to the guild. |
| `message` | `string` | The displayed message. |
| `guild` | `Guild` | Guild object of the guild to be invited to. |


Example:
```json
{
  "invite_url": "https://discord.gg/sxUnqAn",
  "message": "Join the shinpuru Canary test guild to test the latest dev build!",
  "guild": {
    "id": "547762913876639754",
    "name": "shinpuru Canary Testing",
    "icon": "48684c3ad6e9b42793675fda9b09d64c",
    "region": "eu-central",
    "afk_channel_id": "",
    "owner_id": "221905671296253953",
    "joined_at": "",
    "splash": "",
    "member_count": 0,
    "verification_level": 0,
    "embed_enabled": false,
    "large": false,
    "unavailable": false,
    "mfa_level": 0,
    "description": "",
    "banner": "",
    "premium_tier": 0,
    "premium_subscription_count": 0,
    "roles": [
      {
        "id": "547762913876639754",
        "name": "@everyone",
        "managed": false,
        "mentionable": false,
        "hoist": false,
        "color": 0,
        "position": 0,
        "permissions": 104193601
      }
    ],
    "channels": null,
    "self_member": null,
    "icon_url": "https://cdn.discordapp.com/icons/547762913876639754/48684c3ad6e9b42793675fda9b09d64c.png"
  }
}
```

### GlobalPresence

An object which specifies the bot instance's presence.

| Field | Type | Description |
|-------|------|-------------|
| `game` | `string` | The game text. |
| `status` | `string` | The online status:<br>- `online`<br>- `away`<br>- `dnd`<br>- `invisible` |

Example:
```json
{
  "game": "sp-canary.zekro.de",
  "status": "online"
}
```

### SystemInfo

An object which wraps general information about the instance and the system where the instance is running on.

| Field | Type | Description |
|-------|------|-------------|
| `version` | `string` | The build version of shinpuru. |
| `commit_hash` | `string` | The commit hash of the build of shinpuru. |
| `build_date` | `timestamp` | The timestamp when the build was created. |
| `go_version` | `string` | The Go version used to compile the build. |
| `uptime` | `int` | Number of seconds since instance initialization. |
| `uptime_str` | `string` | Number of seconds since instance initialization as string. |
| `os` | `string` | The os identifier the instance is running on. |
| `arch` | `string` | The architecture the instance is running on. |
| `cpus` | `int` | Number of CPU threads allocated. |
| `go_routines` | `int` | Number of concurrently running goroutines. |
| `stack_use` | `int` | Number of bytes allocated on the stack. |
| `stack_use_str` | `string` | Number of bytes allocated on the stack as string. |
| `heap_use` | `int` | Number of bytes allocated on the heap. |
| `heap_use_str` | `string` | Number of bytes allocated on the heap as string. |
| `bot_user_id` | `string` | The ID of the bot account used. |
| `bot_invite` | `string` | The URL to invite the bot including needed permissions. |
| `guilds` | `int` | Number of guilds the instance is running on. |

Example:
```json
{
  "version": "0.17.0-99-g8c1ce6b",
  "commit_hash": "8c1ce6b9c07412067f955a654003c8e6f030353b",
  "build_date": "2020-07-03T09:23:40Z",
  "go_version": "go1.14.4",
  "uptime": 5308,
  "uptime_str": "5308",
  "os": "linux",
  "arch": "amd64",
  "cpus": 6,
  "go_routines": 29,
  "stack_use": 720896,
  "stack_use_str": "720896",
  "heap_use": 4866048,
  "heap_use_str": "4866048",
  "bot_user_id": "536916384026722314",
  "bot_invite": "https://discordapp.com/api/oauth2/authorize?client_id=536916384026722314&scope=bot&permissions=2080894065",
  "guilds": 9
}
```

### APIToken

An object representing information about an API token.

| Field | Type | Description |
|-------|------|-------------|
| `created` | `timestamp` | The creation timestamp of the token. |
| `expires` | `timestamp` | The expiration timestamp of the token. |
| `last_access` | `timestamp` | The creation timestamp of the token. |
| `hits` | `int` | The number of authentications processed with this token. |
| `token?` | `string` | The token string. **This data is only hydrated on token generation!** |


Example:
```json
{
  "created": "2020-07-03T10:59:06.100282743Z",
  "expires": "2021-07-03T10:59:06.100282743Z",
  "last_access": "0001-01-01T00:00:00Z",
  "hits": 0,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjUzMDk5NDYsImlhdCI6MTU5Mzc3Mzk0NiwiaXNzIjoic2hpbnB1cnUgdi4wLjE3LjAtOTktZzhjMWNlNmIiLCJuYmYiOjE1OTM3NzM5NDYsInN1YiI6IjIyMTkwNTY3MTI5NjI1Mzk1MyIsInNwX3NhbHQiOiJNd2RhUUxpcUJDNWZhNXFkaHdjdVpnPT0ifQ.2kifiXUHJTS-CNw-n8dUMSWD44Dwzb73EsvDtJjS8aE"
}
```

## Endpoints

### Get Self User

> ### `GET /api/me`

Shows general user information about the authenticated user.

**Example Response**

```json
{
  "id": "221905671296253953",
  "email": "",
  "username": "zekro",
  "avatar": "a_752a15d01e68fb5f6f6ec83400461a6a",
  "locale": "",
  "discriminator": "0001",
  "token": "",
  "verified": false,
  "mfa_enabled": false,
  "bot": false,
  "avatar_url": "https://cdn.discordapp.com/avatars/221905671296253953/a_752a15d01e68fb5f6f6ec83400461a6a.gif",
  "created_at": "2016-09-04T08:38:26.976834845Z",
  "bot_owner": true
}
```

### Get Guild List

> ### `GET /api/guilds`

Returns a list of guilds you and shinpuru are sharing.

**Example Response**

```json
{
  "n": 2,
  "data": [
    {
      "id": "362162947738566657",
      "name": "zekro's Privatbutze",
      "icon": "2bdf517d77a79b1d6ba60457bd00128e",
      "icon_url": "https://cdn.discordapp.com/icons/362162947738566657/2bdf517d77a79b1d6ba60457bd00128e.png",
      "region": "europe",
      "owner_id": "221905671296253953",
      "joined_at": "2019-01-21T18:59:09.405000+00:00",
      "member_count": 41
    },
    {
      "id": "526196711962705925",
      "name": "5w4gg3rn4ut_5t4t10n",
      "icon": "5ec6d60236376005794d684c6e9a209a",
      "icon_url": "https://cdn.discordapp.com/icons/526196711962705925/5ec6d60236376005794d684c6e9a209a.png",
      "region": "eu-central",
      "owner_id": "221905671296253953",
      "joined_at": "2019-07-15T13:53:34.795000+00:00",
      "member_count": 7
    }
  ]
}
```

### Get Guild

> ### `GET /api/guilds/:guildid`

Returns details of a guild.

**Example Response**

```json
{
  "id": "362162947738566657",
  "name": "zekro's Privatbutze",
  "icon": "2bdf517d77a79b1d6ba60457bd00128e",
  "region": "europe",
  "afk_channel_id": "384315833104597005",
  "owner_id": "221905671296253953",
  "joined_at": "2019-01-21T18:59:09.405000+00:00",
  "splash": "",
  "member_count": 41,
  "verification_level": 4,
  "embed_enabled": false,
  "large": false,
  "unavailable": false,
  "mfa_level": 0,
  "description": "",
  "banner": "",
  "premium_tier": 1,
  "premium_subscription_count": 3,
  "roles": [
    {
      "id": "362162947738566657",
      "name": "@everyone",
      "managed": false,
      "mentionable": false,
      "hoist": false,
      "color": 0,
      "position": 0,
      "permissions": 37084224
    },
    {
      "id": "362166557721362433",
      "name": "先生 (Sensei)",
      "managed": false,
      "mentionable": true,
      "hoist": true,
      "color": 16758016,
      "position": 15,
      "permissions": 512
    },
    {
      "id": "362166741373288448",
      "name": "先輩 (Senpai)",
      "managed": false,
      "mentionable": true,
      "hoist": true,
      "color": 3447003,
      "position": 10,
      "permissions": 2146958931
    }
  ],
  "channels": [
    {
      "id": "381119295632965655",
      "guild_id": "362162947738566657",
      "name": "general",
      "topic": "\n",
      "type": 0,
      "last_message_id": "720959238976962651",
      "last_pin_timestamp": "2019-06-20T04:05:01.229000+00:00",
      "nsfw": false,
      "icon": "",
      "position": 1,
      "bitrate": 0,
      "recipients": null,
      "permission_overwrites": [
        {
          "id": "362162947738566657",
          "type": "role",
          "deny": 1024,
          "allow": 0
        },
        {
          "id": "362169804146081802",
          "type": "role",
          "deny": 0,
          "allow": 1024
        }
      ],
      "user_limit": 0,
      "parent_id": "362162947738566658",
      "rate_limit_per_user": 0
    }
  ],
  "self_member": {
    "guild_id": "",
    "joined_at": "2020-04-09T20:53:47.658000+00:00",
    "nick": "",
    "deaf": false,
    "mute": false,
    "user": {
      "id": "221905671296253953",
      "email": "",
      "username": "zekro",
      "avatar": "a_752a15d01e68fb5f6f6ec83400461a6a",
      "locale": "",
      "discriminator": "0001",
      "token": "",
      "verified": false,
      "mfa_enabled": false,
      "bot": false
    },
    "roles": [
      "362166557721362433",
      "362169804146081802",
      "406891236407115777",
      "583411624514289675",
      "639385098021765131"
    ],
    "premium_since": "2020-04-15T09:24:24.174000+00:00",
    "avatar_url": "https://cdn.discordapp.com/avatars/221905671296253953/a_752a15d01e68fb5f6f6ec83400461a6a.gif",
    "created_at": "2016-09-04T08:38:26.976834845Z",
    "dominance": 1
  },
  "icon_url": "https://cdn.discordapp.com/icons/362162947738566657/2bdf517d77a79b1d6ba60457bd00128e.png"
}
```

### Get Guild Members List

> ### `GET /api/guilds/:guildid/members`

Returns a list of members of a guild. The upper limit of requestable items is `100`, so this request must be paginated.

**Query Parameters**

| Field | Type | Description |
|-------|------|-------------|
| `limit?` | `int` | Maximum ammount of items per request (`max` = `default` = `100`). |
| `after?` | `string` | Show items after the specified snowflake ID. |


**Example Response**

```json
{
  "n": 1,
  "data": [
    {
      "guild_id": "",
      "joined_at": "2020-04-09T20:53:47.658000+00:00",
      "nick": "",
      "deaf": false,
      "mute": false,
      "user": {
        "id": "221905671296253953",
        "email": "",
        "username": "zekro",
        "avatar": "a_752a15d01e68fb5f6f6ec83400461a6a",
        "locale": "",
        "discriminator": "0001",
        "token": "",
        "verified": false,
        "mfa_enabled": false,
        "bot": false
      },
      "roles": [
        "362166557721362433",
        "362169804146081802",
        "406891236407115777",
        "583411624514289675",
        "639385098021765131"
      ],
      "premium_since": "2020-04-15T09:24:24.174000+00:00",
      "avatar_url": "https://cdn.discordapp.com/avatars/221905671296253953/a_752a15d01e68fb5f6f6ec83400461a6a.gif",
      "created_at": "2016-09-04T08:38:26.976834845Z",
      "dominance": 0
    }
  ]
}
```

### Get Guild Member

> ### `GET /api/guilds/:guildid/:memberid`

Shows information about a specific memebr on the specified guild.

**Example Response**

```json
{
  "guild_id": "362162947738566657",
  "joined_at": "2020-04-09T20:53:47.658000+00:00",
  "nick": "",
  "deaf": false,
  "mute": false,
  "user": {
    "id": "221905671296253953",
    "email": "",
    "username": "zekro",
    "avatar": "a_752a15d01e68fb5f6f6ec83400461a6a",
    "locale": "",
    "discriminator": "0001",
    "token": "",
    "verified": false,
    "mfa_enabled": false,
    "bot": false
  },
  "roles": [
    "362166557721362433",
    "362169804146081802",
    "406891236407115777",
    "583411624514289675",
    "639385098021765131"
  ],
  "premium_since": "2020-04-15T09:24:24.174000+00:00",
  "avatar_url": "https://cdn.discordapp.com/avatars/221905671296253953/a_752a15d01e68fb5f6f6ec83400461a6a.gif",
  "created_at": "2016-09-04T08:38:26.976834845Z",
  "dominance": 1
}
```

### Get Guild Member Permissions

> ### `GET /api/guilds/:guildid/:memberid/permissions`

Returns the calculated permissions rules array of the specified user on the specified guild.

**Example Response**

```json
{
  "permissions": [
    "+sp.guild.config.*",
    "+sp.*",
    "+sp.guild.*",
    "+sp.etc.*",
    "+sp.chat.*"
  ]
}
```

### Get Guild Member Permissions Allowed

> ### `GET /api/guilds/:guildid/:memberid/permissions/allowed`

Returns a full list of all rule domains which are allowed for the specified user on the specified guild.

**Example Response**

```json
{
  "n": 33,
  "data": [
    "sp.etc.help",
    "sp.guild.config.prefix",
    "sp.guild.config.perms",
    "sp.guild.mod.clear",
    "sp.guild.mod.mvall",
    "sp.etc.info",
    "sp.chat.say",
    "sp.chat.quote",
    "sp.game",
    "sp.guild.config.autorole",
    "sp.guild.mod.report",
    "sp.guild.config.modlog",
    "sp.guild.mod.kick",
    "sp.guild.mod.ban",
    "sp.chat.vote",
    "sp.chat.profile",
    "sp.etc.id",
    "sp.guild.mod.mute",
    "sp.guild.mod.ment",
    "sp.chat.notify",
    "sp.guild.config.voicelog",
    "sp.etc.bug",
    "sp.etc.stats",
    "sp.chat.twitch",
    "sp.guild.mod.ghostping",
    "sp.chat.exec",
    "sp.guild.admin.backup",
    "sp.guild.mod.inviteblock",
    "sp.chat.tag",
    "sp.guild.config.joinmsg",
    "sp.guild.config.leavemsg",
    "sp.etc.snowflake",
    "sp.chat.chanstats"
  ]
}
```

### Get Guild Member Permissions

> ### `GET /api/guilds/:guildid/:memberid/permissions`

Returns the calculated permissions rules array of the specified user on the specified guild.

**Example Response**

```json
{
  "permissions": [
    "+sp.guild.config.*",
    "+sp.*",
    "+sp.guild.*",
    "+sp.etc.*",
    "+sp.chat.*"
  ]
}
```

### Get Guild Reports

> ### `GET /api/guilds/:guildid/reports`

Displays a list of reports on the specified guild. This request can also be paginated.

**Query Parameters**

| Field | Type | Description |
|-------|------|-------------|
| `limit?` | `int` | Maximum ammount of items per request. |
| `offset?` | `int` | Ammount of items to be skipped. |

**Example Response**

```json
{
  "n": 1,
  "data": [
    {
      "id": "6678266303259619328",
      "type": 3,
      "guild_id": "547762913876639754",
      "executor_id": "221905671296253953",
      "victim_id": "455819141245304832",
      "message": "Bad language",
      "attachment_url": "https://sp-canary.zekro.de/imagestore/6678266279931420672.png",
      "type_name": "WARN",
      "created": "2020-07-03T09:29:57Z"
    }
  ]
}
```

### Get Guild Reports Count

> ### `GET /api/guilds/:guildid/reports/counts`

Returns the accumulated ammount of reports on this guild.

**Example Response**

```json
{
  "count": 1
}
```

### Get Guild Member Reports

> ### `GET /api/guilds/:guildid/:memberid/reports`

Displays a list of reports for the specified member on the specified guild. This request can also be paginated.

**Query Parameters**

| Field | Type | Description |
|-------|------|-------------|
| `limit?` | `int` | Maximum ammount of items per request. |
| `offset?` | `int` | Ammount of items to be skipped. |

**Example Response**

```json
{
  "n": 1,
  "data": [
    {
      "id": "6678266303259619328",
      "type": 3,
      "guild_id": "547762913876639754",
      "executor_id": "221905671296253953",
      "victim_id": "455819141245304832",
      "message": "Bad language",
      "attachment_url": "https://sp-canary.zekro.de/imagestore/6678266279931420672.png",
      "type_name": "WARN",
      "created": "2020-07-03T09:29:57Z"
    }
  ]
}
```

### Get Guild Reports Count

> ### `GET /api/guilds/:guildid/:memberid/reports/counts`

Returns the accumulated ammount of reports on this guild for the specified member.

**Example Response**

```json
{
  "count": 1
}
```

### Get Guild Settings

> ### `GET /api/guilds/:guildid/settings`

Returns the settings and preferences for the specified guild.

**Example Response**

```json
{
  "prefix": "!",
  "perms": {
    "608577686528458762": [
      "+sp.guild.config.*"
    ],
    "608584907014537246": [
      "+sp.guild.mod.*"
    ]
  },
  "autorole": "547772764921004043",
  "modlogchannel": "547773902357790730",
  "voicelogchannel": "547774364968288296",
  "joinmessagechannel": "Welcome [ment]! :heart:",
  "joinmessagetext": "547762913876639757",
  "leavemessagechannel": "",
  "leavemessagetext": ""
}
```

### Set Guild Settings

> ### `POST /api/guilds/:guildid/settings`

Returns the settings and preferences for the specified guild.

**Requires Permissions**

- `sp.guild.config.autorole`
- `sp.guild.config.modlog`
- `sp.guild.config.prefix`
- `sp.guild.config.voicelog`
- `sp.guild.config.joinmsg`
- `sp.guild.config.leavemsg`

**Body Parameters**

Only set parameters will be updated. If you want to actively reset a value, you need to set the value to `__RESET__`.

| Field | Type | Description |
|-------|------|-------------|
| `prefix` | `string` | The guild specific prefix. |
| `autorole` | `string` | The snowflake ID of the set autorole. |
| `modlogchannel` | `string` | The snowflake ID of the set modlog channel. |
| `voicelogchannel` | `string` | The snowflake ID of the set voicelog channel. |
| `joinmessagechannel` | `string` | The snowflake ID of the channel where join messages are sent to. |
| `joinmessagetext` | `string` | The text which is sent into the join message channel when a user joins the guild. |
| `leavemessagechannel` | `string` | The snowflake ID of the channel where leave messages are sent to. |
| `leavemessagetext` | `string` | The text which is sent into the leave message channel when a user leaves the guild. |

**Example Response**

```json
{
  "code": 200,
  "message": "ok"
}
```

### Get Guild Permissions

> ### `GET /api/guilds/:guildid/permissions`

Returns the defined rule sets for roles on the guild.

**Example Response**

```json
{
  "608577686528458762": [
    "+sp.guild.config.*"
  ],
  "608584907014537246": [
    "+sp.guild.mod.*"
  ]
}
```

### Set Guild Permissions

> ### `POST /api/guilds/:guildid/permissions`

Defines a new rule for specified role IDs.

**Requires Permissions**

- `sp.guild.config.perms`

**Body Parameters**

| Field | Type | Description |
|-------|------|-------------|
| `perm` | `string` | The permission rule, for example:<br>`+sp.guild.mod.*` |
| `roles` | `string[]` | The list of roles IDs to set the rule for. |

**Example Response**

```json
{
  "code": 200,
  "message": "ok"
}
```

### Create Member Report

> ### `POST /api/guilds/:guildid/:memberid/reports`

Records a report for a member on the specified guild.

**Required Permissions**

- `sp.guild.mod.report`

**Body Parameters**

| Field | Type | Description |
|-------|------|-------------|
| `type` | `int` | Type of the report. Only these two are available for manual spec:<br>`3` - WARN<br>`4` - AD |
| `reason` | `string` | The issue reason. |
| `attachment?` | `string` | An image URL as report attachment. |

**Example Response**

```json
{
  "id": "6678266303259619328",
  "type": 3,
  "guild_id": "547762913876639754",
  "executor_id": "221905671296253953",
  "victim_id": "455819141245304832",
  "message": "Bad language",
  "attachment_url": "https://sp-canary.zekro.de/imagestore/6678266279931420672.png",
  "type_name": "WARN",
  "created": "2020-07-03T09:29:57Z"
}
```

### Create Member Kick Issue

> ### `POST /api/guilds/:guildid/:memberid/kick`

Issues a member kick which is recorded with a kick report.

**Required Permissions**

- `sp.guild.mod.kick`

**Body Parameters**

| Field | Type | Description |
|-------|------|-------------|
| `reason` | `string` | The issue reason. |
| `attachment?` | `string` | An image URL as report attachment. |

**Example Response**

```json
{
  "id": "6678266303259619328",
  "type": 0,
  "guild_id": "547762913876639754",
  "executor_id": "221905671296253953",
  "victim_id": "455819141245304832",
  "message": "Bad language",
  "attachment_url": "https://sp-canary.zekro.de/imagestore/6678266279931420672.png",
  "type_name": "KICK",
  "created": "2020-07-03T09:29:57Z"
}
```

### Create Member Ban Issue

> ### `POST /api/guilds/:guildid/:memberid/ban`

Issues a member ban which is recorded with a ban report.

**Required Permissions**

- `sp.guild.mod.ban`

**Body Parameters**

| Field | Type | Description |
|-------|------|-------------|
| `reason` | `string` | The issue reason. |
| `attachment?` | `string` | An image URL as report attachment. |

**Example Response**

```json
{
  "id": "6678266303259619328",
  "type": 1,
  "guild_id": "547762913876639754",
  "executor_id": "221905671296253953",
  "victim_id": "455819141245304832",
  "message": "Bad language",
  "attachment_url": "https://sp-canary.zekro.de/imagestore/6678266279931420672.png",
  "type_name": "BAN",
  "created": "2020-07-03T09:29:57Z"
}
```

### Get Report

> ### `GET /api/reports/:caseid`

Returns information about a report by case ID.

**Example Response**

```json
{
  "id": "6678266303259619328",
  "type": 1,
  "guild_id": "547762913876639754",
  "executor_id": "221905671296253953",
  "victim_id": "455819141245304832",
  "message": "Bad language",
  "attachment_url": "https://sp-canary.zekro.de/imagestore/6678266279931420672.png",
  "type_name": "BAN",
  "created": "2020-07-03T09:29:57Z"
}
```

### Get Bot Presence

> ### `GET /api/settings/presence`

Displays the bot instance set presence.

**Example Response**

```json
{
  "game": "sp-canary.zekro.de",
  "status": "online"
}
```

### Set Bot Presence

> ### `POST /api/settings/presence`

**Required Permissions**

- `sp.game`

Set the bot instance presence.

**Body Parameters**

| Field | Type | Description |
|-------|------|-------------|
| `game` | `string` | The presence game. |
| `status` | `string` | The online status:<br>- `online`<br>- `away`<br>- `dnd`<br>- `invisible` |

**Example Response**

```json
{
  "game": "sp-canary.zekro.de",
  "status": "online"
}
```

### Get No Guild Invite Setting

> ### `GET /api/settings/noguildinvite`

Displays the guild invite set which is displayed when soneone hits the web interface who does not share any guild with the bot instance account.

**Example Response**

```json
{
    "invite_url": "https://discord.gg/sxUnqAn",
  "message": "Join the shinpuru Canary test guild to test the latest dev build!",
  "guild": {
    "id": "547762913876639754",
    "name": "shinpuru Canary Testing",
    "icon": "48684c3ad6e9b42793675fda9b09d64c",
    "region": "eu-central",
    "afk_channel_id": "",
    "owner_id": "221905671296253953",
    "joined_at": "",
    "splash": "",
    "member_count": 0,
    "verification_level": 0,
    "embed_enabled": false,
    "large": false,
    "unavailable": false,
    "mfa_level": 0,
    "description": "",
    "banner": "",
    "premium_tier": 0,
    "premium_subscription_count": 0,
    "roles": [
      {
        "id": "547762913876639754",
        "name": "@everyone",
        "managed": false,
        "mentionable": false,
        "hoist": false,
        "color": 0,
        "position": 0,
        "permissions": 104193601
      }
    ],
    "channels": null,
    "self_member": null,
    "icon_url": "https://cdn.discordapp.com/icons/547762913876639754/48684c3ad6e9b42793675fda9b09d64c.png"
  }
}
```

### Set No Guild Invite Setting

> ### `POST /api/settings/noguildinvite`

Set the guild invite set which is displayed when soneone hits the web interface who does not share any guild with the bot instance account.

**Required Permissions**

- `sp.noguildinvite`

**Body Parameters**

| Field | Type | Description |
|-------|------|-------------|
| `guild_id` | `string` | The snowflake ID of the guild to be invited to. |
| `message` | `string` | The message displayed on the web interface page. |
| `invite_code?` | `string` | Set an already defined invite code. Otherwise, a new invite will be generated by shinpuru. |

**Example Response**

```json
{
  "code": 200,
  "message": "ok"
}
```

### Get System and Instance Information

> ### `POST /api/sysinfo`

Returns detailed information about the shinpuru instance and the system it is running on.

**Example Response**

```json
{
  "version": "0.17.0-100-g8faeb1a",
  "commit_hash": "8faeb1a5ab0aa574e896ba0f67372e096984b4e6",
  "build_date": "2020-07-03T11:04:01Z",
  "go_version": "go1.14.4",
  "uptime": 44384,
  "uptime_str": "44384",
  "os": "linux",
  "arch": "amd64",
  "cpus": 6,
  "go_routines": 20,
  "stack_use": 884736,
  "stack_use_str": "884736",
  "heap_use": 5136384,
  "heap_use_str": "5136384",
  "bot_user_id": "536916384026722314",
  "bot_invite": "https://discordapp.com/api/oauth2/authorize?client_id=536916384026722314&scope=bot&permissions=2080894065",
  "guilds": 9
}
```

### Get API Token

> ### `GET /api/token`

Returns information about a generated API token for the authenticated user. This does not contain the API token itself!

**Example Response**

```json
{
  "created": "2020-07-03T10:59:52.18640661Z",
  "expires": "2021-07-03T10:59:52.18640661Z",
  "last_access": "2020-07-05T09:43:40.923373646Z",
  "hits": 34
}
```

### Post API Token

> ### `POST /api/token`

Create a new API token for the authenticated user which is returned in the response. If a token was generated before, the old token will be rendered invalid.

**Example Response**

```json
{
  "created": "2020-07-03T10:59:06.100282743Z",
  "expires": "2021-07-03T10:59:06.100282743Z",
  "last_access": "0001-01-01T00:00:00Z",
  "hits": 0,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjUzMDk5NDYsImlhdCI6MTU5Mzc3Mzk0NiwiaXNzIjoic2hpbnB1cnUgdi4wLjE3LjAtOTktZzhjMWNlNmIiLCJuYmYiOjE1OTM3NzM5NDYsInN1YiI6IjIyMTkwNTY3MTI5NjI1Mzk1MyIsInNwX3NhbHQiOiJNd2RhUUxpcUJDNWZhNXFkaHdjdVpnPT0ifQ.2kifiXUHJTS-CNw-n8dUMSWD44Dwzb73EsvDtJjS8aE"
}
```

### Unset API Token

> ### `DELETE /api/token`

Deletes the generated API token for the authenticated user and makes it invalid.

**Example Response**

```json
{
  "code": 200,
  "message": "ok"
}
```