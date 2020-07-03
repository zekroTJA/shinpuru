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
  "icon_url": "https://cdn.discordapp.com/icons/362162947738566657/2bdf517d77a79b1d6ba60457bd00128e.png"
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

