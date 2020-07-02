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
