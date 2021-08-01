# shinpuru main API
The shinpuru main REST API.

## Version: 1.0

### /auth/accesstoken

#### POST
##### Summary

Access Token Exchange

##### Description

Exchanges the cookie-passed refresh token with a generated access token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.AccessTokenResponse](#modelsaccesstokenresponse) |
| 401 | Unauthorized | [models.Error](#modelserror) |

### /auth/check

#### GET
##### Summary

Authorization Check

##### Description

Returns OK if the request is authorized.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Status](#modelsstatus) |
| 401 | Unauthorized | [models.Error](#modelserror) |

### /auth/logout

#### POST
##### Summary

Logout

##### Description

Reovkes the currently used access token and clears the refresh token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Status](#modelsstatus) |

### /guilds

#### GET
##### Summary

List Guilds

##### Description

Returns a list of guilds the authenticated user has in common with shinpuru.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [models.GuildReduced](#modelsguildreduced) ] |
| 401 | Unauthorized | [models.Error](#modelserror) |

### /guilds/:id

#### GET
##### Summary

Get Guild

##### Description

Returns a single guild object by it's ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Guild](#modelsguild) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/antiraid/joinlog

#### GET
##### Summary

Get Antiraid Joinlog

##### Description

Returns a list of joined members during an antiraid trigger.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [models.JoinLogEntry](#modelsjoinlogentry) ] |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

#### DELETE
##### Summary

Reset Antiraid Joinlog

##### Description

Deletes all entries of the antiraid joinlog.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Status](#modelsstatus) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/backups

#### GET
##### Summary

Get Guild Backups

##### Description

Returns a list of guild backups.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [backupmodels.Entry](#backupmodelsentry) ] |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/backups/:backupid/download

#### GET
##### Summary

Download Backup File

##### Description

Download a single gziped backup file.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| backupid | path | The ID of the backup. | Yes | string |
| ota_token | query | The previously obtained OTA token to authorize the download. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | data |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 403 | Forbidden | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

#### POST
##### Summary

Obtain Backup Download OTA Key

##### Description

Returns an OTA key which is used to download a backup entry.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| backupid | path | The ID of the backup. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.AccessTokenResponse](#modelsaccesstokenresponse) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/backups/toggle

#### POST
##### Summary

Toggle Guild Backup Enable

##### Description

Toggle guild backup enable state.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| payload | body | Enable state payload. | Yes | [models.EnableStatus](#modelsenablestatus) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Status](#modelsstatus) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/inviteblock

#### POST
##### Summary

Toggle Guild Inviteblock Enable

##### Description

Toggle enabled state of the guild invite block system.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The enable status payload. | Yes | [models.EnableStatus](#modelsenablestatus) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Status](#modelsstatus) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/permissions

#### GET
##### Summary

Get Guild Permission Settings

##### Description

Returns the specified guild permission settings.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.PermissionsMap](#modelspermissionsmap) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

#### POST
##### Summary

Apply Guild Permission Rule

##### Description

Apply a new guild permission rule for a specified role.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The permission rule payload. | Yes | [models.PermissionsUpdate](#modelspermissionsupdate) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Status](#modelsstatus) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/reports

#### GET
##### Summary

Get Guild Modlog

##### Description

Returns a list of guild modlog entries for the given guild.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| offset | query | The offset of returned entries | No | integer |
| limit | query | The amount of returned entries (0 = all) | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [models.Report](#modelsreport) ] |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/reports/count

#### GET
##### Summary

Get Guild Modlog Count

##### Description

Returns the total count of entries in the guild mod log.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Count](#modelscount) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/scoreboard

#### GET
##### Summary

Get Guild Scoreboard

##### Description

Returns a list of scoreboard entries for the given guild.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| limit | query | Limit the amount of result values | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [models.GuildKarmaEntry](#modelsguildkarmaentry) ] |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/settings

#### GET
##### Summary

Get Guild Settings

##### Description

Returns the specified general guild settings.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.GuildSettings](#modelsguildsettings) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

#### POST
##### Summary

Get Guild Settings

##### Description

Returns the specified general guild settings.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| payload | body | Modified guild settings payload. | Yes | [models.GuildSettings](#modelsguildsettings) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Status](#modelsstatus) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/settings/antiraid

#### GET
##### Summary

Get Guild Antiraid Settings

##### Description

Returns the specified guild antiraid settings.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.AntiraidSettings](#modelsantiraidsettings) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

#### POST
##### Summary

Update Guild Antiraid Settings

##### Description

Update the guild antiraid settings specification.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The guild antiraid settings payload. | Yes | [models.AntiraidSettings](#modelsantiraidsettings) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Status](#modelsstatus) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/settings/flushguilddata

#### POST
##### Summary

Flush Guild Data

##### Description

Flushes all guild data from the database.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The guild flush payload. | Yes | [models.FlushGuildRequest](#modelsflushguildrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.State](#modelsstate) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/settings/karma

#### GET
##### Summary

Get Guild Karma Settings

##### Description

Returns the specified guild karma settings.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.KarmaSettings](#modelskarmasettings) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

#### POST
##### Summary

Update Guild Karma Settings

##### Description

Update the guild karma settings specification.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The guild karma settings payload. | Yes | [models.KarmaSettings](#modelskarmasettings) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Status](#modelsstatus) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/settings/karma/blocklist

#### GET
##### Summary

Get Guild Karma Blocklist

##### Description

Returns the specified guild karma blocklist entries.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [models.Member](#modelsmember) ] |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/settings/karma/blocklist/:memberid

#### PUT
##### Summary

Add Guild Karma Blocklist Entry

##### Description

Add a guild karma blocklist entry.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Status](#modelsstatus) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

#### DELETE
##### Summary

Remove Guild Karma Blocklist Entry

##### Description

Remove a guild karma blocklist entry.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Status](#modelsstatus) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/settings/karma/rules

#### GET
##### Summary

Get Guild Settings Karma Rules

##### Description

Returns a list of specified guild karma rules.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [models.KarmaRule](#modelskarmarule) ] |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

#### POST
##### Summary

Create Guild Settings Karma

##### Description

Create a guild karma rule.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The karma rule payload. | Yes | [models.KarmaRule](#modelskarmarule) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.KarmaRule](#modelskarmarule) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/settings/karma/rules/:ruleid

#### POST
##### Summary

Update Guild Settings Karma

##### Description

Update a karma rule by ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| ruleid | path | The ID of the rule. | Yes | string |
| payload | body | The karma rule update payload. | Yes | [models.KarmaRule](#modelskarmarule) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.KarmaRule](#modelskarmarule) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

#### DELETE
##### Summary

Remove Guild Settings Karma

##### Description

Remove a guild karma rule by ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| ruleid | path | The ID of the rule. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.State](#modelsstate) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/settings/logs

#### GET
##### Summary

Get Guild Log Count

##### Description

Returns the total or filtered count of guild log entries.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Count](#modelscount) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/settings/logs/state

#### GET
##### Summary

Get Guild Settings Log State

##### Description

Returns the enabled state of the guild log setting.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.State](#modelsstate) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

#### POST
##### Summary

Update Guild Settings Log State

##### Description

Update the enabled state of the log state guild setting.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The state payload. | Yes | [models.State](#modelsstate) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.State](#modelsstate) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/settings/logs/state/:entryid

#### DELETE
##### Summary

Delete Guild Log Entries

##### Description

Delete a single or all guild log entries.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| entryid | path | The ID of the entry to be deleted. | No | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.State](#modelsstate) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/starboard

#### GET
##### Summary

Get Guild Starboard

##### Description

Returns a list of starboard entries for the given guild.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [models.StarboardEntryResponse](#modelsstarboardentryresponse) ] |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/unbanrequests

#### GET
##### Summary

Get Guild Unbanrequests

##### Description

Returns the list of the guild unban requests.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListReponse | [ [models.UnbanRequest](#modelsunbanrequest) ] |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/unbanrequests/:requestid

#### GET
##### Summary

Get Single Guild Unbanrequest

##### Description

Returns a single guild unban request by ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| requestid | path | The ID of the unbanrequest. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.UnbanRequest](#modelsunbanrequest) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

#### POST
##### Summary

Process Guild Unbanrequest

##### Description

Process a guild unban request.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |
| requestid | path | The ID of the unbanrequest. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.UnbanRequest](#modelsunbanrequest) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /guilds/:id/unbanrequests/count

#### GET
##### Summary

Get Guild Unbanrequests Count

##### Description

Returns the total or filtered count of guild unban requests.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Count](#modelscount) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /me

#### GET
##### Summary

Me

##### Description

Returns the user object of the currently authenticated user.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.User](#modelsuser) |

### /ota

#### GET
##### Summary

OTA Login

##### Description

Logs in the current browser session by using the passed pre-obtained OTA token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 401 | Unauthorized | [models.Error](#modelserror) |

### /reports/:id

#### GET
##### Summary

Get Report

##### Description

Returns a single report object by its ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The report ID. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Report](#modelsreport) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /reports/:id/revoke

#### POST
##### Summary

Revoke Report

##### Description

Revokes a given report by ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | The report ID. | Yes | string |
| payload | body | The revoke reason payload. | Yes | [models.ReasonRequest](#modelsreasonrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Report](#modelsreport) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /settings/noguildinvite

#### GET
##### Summary

Get No Guild Invites Status

##### Description

Returns the settings status for the suggested guild invite when the logged in user is not on any guild with shinpuru.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.InviteSettingsResponse](#modelsinvitesettingsresponse) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 409 | Returned when no channel could be found to create invite for. | [models.Error](#modelserror) |

#### POST
##### Summary

Set No Guild Invites Status

##### Description

Set the status for the suggested guild invite when the logged in user is not on any guild with shinpuru.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| payload | body | Invite Settings Payload | Yes | [models.InviteSettingsRequest](#modelsinvitesettingsrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.APITokenResponse](#modelsapitokenresponse) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 409 | Returned when no channel could be found to create invite for. | [models.Error](#modelserror) |

### /settings/presence

#### GET
##### Summary

Get Presence

##### Description

Returns the bot's displayed presence status.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [presence.Presence](#presencepresence) |
| 401 | Unauthorized | [models.Error](#modelserror) |

#### POST
##### Summary

Set Presence

##### Description

Set the bot's displayed presence status.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| payload | body | Presence Payload | Yes | [presence.Presence](#presencepresence) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.APITokenResponse](#modelsapitokenresponse) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Is returned when no token was generated before. | [models.Error](#modelserror) |

### /sysinfo

#### GET
##### Summary

System Information

##### Description

Returns general global system information.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.SystemInfo](#modelssysteminfo) |

### /token

#### GET
##### Summary

API Token Info

##### Description

Returns general metadata information about a generated API token. The response does **not** contain the actual token!

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.APITokenResponse](#modelsapitokenresponse) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Is returned when no token was generated before. | [models.Error](#modelserror) |

#### POST
##### Summary

API Token Generation

##### Description

(Re-)Generates and returns general metadata information about an API token **including** the actual API token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.APITokenResponse](#modelsapitokenresponse) |
| 401 | Unauthorized | [models.Error](#modelserror) |

#### DELETE
##### Summary

API Token Deletion

##### Description

Invalidates the currently generated API token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.Status](#modelsstatus) |
| 401 | Unauthorized | [models.Error](#modelserror) |

### /unbanrequests

#### GET
##### Summary

Get Unban Requests

##### Description

Returns a list of unban requests created by the authenticated user.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [models.UnbanRequest](#modelsunbanrequest) ] |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

#### POST
##### Summary

Create Unban Requests

##### Description

Create an unban reuqest.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| payload | body | The unban request payload. | Yes | [models.UnbanRequest](#modelsunbanrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.UnbanRequest](#modelsunbanrequest) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /unbanrequests/bannedguilds

#### GET
##### Summary

Get Banned Guilds

##### Description

Returns a list of guilds where the currently authenticated user is banned.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [models.GuildReduced](#modelsguildreduced) ] |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /usersettings/ota

#### GET
##### Summary

Get OTA Usersettings State

##### Description

Returns the current state of the OTA user setting.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.UsersettingsOTA](#modelsusersettingsota) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

#### POST
##### Summary

Update OTA Usersettings State

##### Description

Update the OTA user settings state.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| payload | body | The OTA settings payload. | Yes | [models.UsersettingsOTA](#modelsusersettingsota) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.UsersettingsOTA](#modelsusersettingsota) |
| 400 | Bad Request | [models.Error](#modelserror) |
| 401 | Unauthorized | [models.Error](#modelserror) |
| 404 | Not Found | [models.Error](#modelserror) |

### /util/color/:hexcode

#### GET
##### Summary

Color Generator

##### Description

Produces a square image of the given color and size.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hexcode | path | Hex Code of the Color to produce | Yes | string |
| size | query | The dimension of the square image | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | data |

### /util/commands

#### GET
##### Summary

Command List

##### Description

Returns a list of registered commands and their description.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [models.CommandInfo](#modelscommandinfo) ] |

### /util/landingpageinfo

#### GET
##### Summary

Landing Page Info

##### Description

Returns general information for the landing page like the local invite parameters.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.LandingPageResponse](#modelslandingpageresponse) |

### Models

#### backupmodels.Entry

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| file_id | string |  | No |
| guild_id | string |  | No |
| timestamp | string |  | No |

#### discordgo.Channel

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| application_id | string | ApplicationID of the DM creator Zeroed if guild channel or not a bot user | No |
| bitrate | integer | The bitrate of the channel, if it is a voice channel. | No |
| guild_id | string | The ID of the guild to which the channel belongs, if it is in a guild. Else, this ID is empty (e.g. DM channels). | No |
| icon | string | Icon of the group DM channel. | No |
| id | string | The ID of the channel. | No |
| last_message_id | string | The ID of the last message sent in the channel. This is not guaranteed to be an ID of a valid message. | No |
| last_pin_timestamp | string | The timestamp of the last pinned message in the channel. Empty if the channel has no pinned messages. | No |
| name | string | The name of the channel. | No |
| nsfw | boolean | Whether the channel is marked as NSFW. | No |
| owner_id | string | ID of the DM creator Zeroed if guild channel | No |
| parent_id | string | The ID of the parent channel, if the channel is under a category | No |
| permission_overwrites | [ [discordgo.PermissionOverwrite](#discordgopermissionoverwrite) ] | A list of permission overwrites present for the channel. | No |
| position | integer | The position of the channel, used for sorting in client. | No |
| rate_limit_per_user | integer | Amount of seconds a user has to wait before sending another message (0-21600) bots, as well as users with the permission manage_messages or manage_channel, are unaffected | No |
| recipients | [ [discordgo.User](#discordgouser) ] | The recipients of the channel. This is only populated in DM channels. | No |
| topic | string | The topic of the channel. | No |
| type | integer | The type of the channel. | No |
| user_limit | integer | The user limit of the voice channel. | No |

#### discordgo.PermissionOverwrite

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| allow | string | _Example:_ `"0"` | No |
| deny | string | _Example:_ `"0"` | No |
| id | string |  | No |
| type | integer |  | No |

#### discordgo.Role

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| color | integer | The hex color of this role. | No |
| hoist | boolean | Whether this role is hoisted (shows up separately in member list). | No |
| id | string | The ID of the role. | No |
| managed | boolean | Whether this role is managed by an integration, and thus cannot be manually added to, or taken from, members. | No |
| mentionable | boolean | Whether this role is mentionable. | No |
| name | string | The name of the role. | No |
| permissions | string | The permissions of the role on the guild (doesn't include channel overrides). This is a combination of bit masks; the presence of a certain permission can be checked by performing a bitwise AND between this int and the permission.<br>_Example:_ `"0"` | No |
| position | integer | The position of this role in the guild's role hierarchy. | No |

#### discordgo.User

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| avatar | string | The hash of the user's avatar. Use Session.UserAvatar to retrieve the avatar itself. | No |
| bot | boolean | Whether the user is a bot. | No |
| discriminator | string | The discriminator of the user (4 numbers after name). | No |
| email | string | The email of the user. This is only present when the application possesses the email scope for the user. | No |
| flags | integer | The flags on a user's account. Only available when the request is authorized via a Bearer token. | No |
| id | string | The ID of the user. | No |
| locale | string | The user's chosen language option. | No |
| mfa_enabled | boolean | Whether the user has multi-factor authentication enabled. | No |
| premium_type | integer | The type of Nitro subscription on a user's account. Only available when the request is authorized via a Bearer token. | No |
| public_flags | integer | The public flags on a user's account. This is a combination of bit masks; the presence of a certain flag can be checked by performing a bitwise AND between this int and the flag. | No |
| system | boolean | Whether the user is an Official Discord System user (part of the urgent message system). | No |
| token | string | The token of the user. This is only present for the user represented by the current session. | No |
| username | string | The user's username. | No |
| verified | boolean | Whether the user's email is verified. | No |

#### models.APITokenResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created | string |  | No |
| expires | string |  | No |
| hits | integer |  | No |
| last_access | string |  | No |
| token | string |  | No |

#### models.AccessTokenResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| expires | string |  | No |
| token | string |  | No |

#### models.AntiraidSettings

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| burst | integer |  | No |
| regeneration_period | integer |  | No |
| state | boolean |  | No |

#### models.CommandInfo

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string |  | No |
| domain_name | string |  | No |
| group | string |  | No |
| help | string |  | No |
| invokes | [ string ] |  | No |
| is_executable_in_dm | boolean |  | No |
| sub_permission_rules | [ [shireikan.SubPermission](#shireikansubpermission) ] |  | No |

#### models.Count

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| count | integer |  | No |

#### models.EnableStatus

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| enabled | boolean |  | No |

#### models.Error

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | integer |  | No |
| context | string |  | No |
| error | string |  | No |

#### models.FlushGuildRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| leave_after | boolean |  | No |
| validation | string |  | No |

#### models.Guild

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| afk_channel_id | string |  | No |
| backups_enabled | boolean |  | No |
| banner | string |  | No |
| channels | [ [discordgo.Channel](#discordgochannel) ] |  | No |
| description | string |  | No |
| icon | string |  | No |
| icon_url | string |  | No |
| id | string |  | No |
| invite_block_enabled | boolean |  | No |
| joined_at | string |  | No |
| large | boolean |  | No |
| latest_backup_entry | string |  | No |
| member_count | integer |  | No |
| mfa_level | integer |  | No |
| name | string |  | No |
| owner_id | string |  | No |
| premium_subscription_count | integer |  | No |
| premium_tier | integer |  | No |
| region | string |  | No |
| roles | [ [discordgo.Role](#discordgorole) ] |  | No |
| self_member | [models.Member](#modelsmember) |  | No |
| splash | string |  | No |
| unavailable | boolean |  | No |
| verification_level | integer |  | No |

#### models.GuildKarmaEntry

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| member | [models.Member](#modelsmember) |  | No |
| value | integer |  | No |

#### models.GuildLogEntry

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| guildid | string |  | No |
| id | integer |  | No |
| message | string |  | No |
| module | string |  | No |
| severity | integer |  | No |
| timestamp | string |  | No |

#### models.GuildReduced

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| icon | string |  | No |
| icon_url | string |  | No |
| id | string |  | No |
| joined_at | string |  | No |
| member_count | integer |  | No |
| name | string |  | No |
| owner_id | string |  | No |
| region | string |  | No |

#### models.GuildSettings

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| autoroles | [ string ] |  | No |
| joinmessagechannel | string |  | No |
| joinmessagetext | string |  | No |
| leavemessagechannel | string |  | No |
| leavemessagetext | string |  | No |
| modlogchannel | string |  | No |
| perms | object |  | No |
| prefix | string |  | No |
| voicelogchannel | string |  | No |

#### models.InviteSettingsRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| guild_id | string |  | No |
| invite_code | string |  | No |
| message | string |  | No |

#### models.InviteSettingsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| guild | [models.Guild](#modelsguild) |  | No |
| invite_url | string |  | No |
| message | string |  | No |

#### models.JoinLogEntry

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| guild_id | string |  | No |
| tag | string |  | No |
| timestamp | string |  | No |
| user_id | string |  | No |

#### models.KarmaRule

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| action | string |  | No |
| argument | string |  | No |
| guildid | string |  | No |
| id | integer |  | No |
| trigger | integer |  | No |
| value | integer |  | No |

#### models.KarmaSettings

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| emotes_decrease | [ string ] |  | No |
| emotes_increase | [ string ] |  | No |
| penalty | boolean |  | No |
| state | boolean |  | No |
| tokens | integer |  | No |

#### models.LandingPageResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| localinvite | string |  | No |
| publiccaranyinvite | string |  | No |
| publicmaininvite | string |  | No |

#### models.Member

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| avatar_url | string |  | No |
| chat_muted | boolean |  | No |
| created_at | string |  | No |
| deaf | boolean | Whether the member is deafened at a guild level. | No |
| dominance | integer |  | No |
| guild_id | string | The guild ID on which the member exists. | No |
| joined_at | string | The time at which the member joined the guild, in ISO8601. | No |
| karma | integer |  | No |
| karma_total | integer |  | No |
| mute | boolean | Whether the member is muted at a guild level. | No |
| nick | string | The nickname of the member, if they have one. | No |
| pending | boolean | Is true while the member hasn't accepted the membership screen. | No |
| premium_since | string | When the user used their Nitro boost on the server | No |
| roles | [ string ] | A list of IDs of the roles which are possessed by the member. | No |
| user | [discordgo.User](#discordgouser) | The underlying user on which the member is based. | No |

#### models.PermissionsMap

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| models.PermissionsMap | object |  |  |

#### models.PermissionsUpdate

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| perm | string |  | No |
| role_ids | [ string ] |  | No |

#### models.ReasonRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| attachment | string |  | No |
| reason | string |  | No |
| timeout | string |  | No |

#### models.Report

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| attachment_url | string |  | No |
| created | string |  | No |
| executor_id | string |  | No |
| guild_id | string |  | No |
| id | integer |  | No |
| message | string |  | No |
| timeout | string |  | No |
| type | integer |  | No |
| type_name | string |  | No |
| victim_id | string |  | No |

#### models.StarboardEntryResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| author_avatar_url | string |  | No |
| author_id | string |  | No |
| author_username | string |  | No |
| channel_id | string |  | No |
| content | string |  | No |
| guild_id | string |  | No |
| media_urls | [ string ] |  | No |
| message_id | string |  | No |
| message_url | string |  | No |
| score | integer |  | No |
| starboard_id | string |  | No |

#### models.State

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| state | boolean |  | No |

#### models.Status

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | integer |  | No |

#### models.SystemInfo

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| arch | string |  | No |
| bot_invite | string |  | No |
| bot_user_id | string |  | No |
| build_date | string |  | No |
| commit_hash | string |  | No |
| cpus | integer |  | No |
| go_routines | integer |  | No |
| go_version | string |  | No |
| guilds | integer |  | No |
| heap_use | integer |  | No |
| heap_use_str | string |  | No |
| os | string |  | No |
| stack_use | integer |  | No |
| stack_use_str | string |  | No |
| uptime | integer |  | No |
| uptime_str | string |  | No |
| version | string |  | No |

#### models.UnbanRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created | string |  | No |
| guild_id | string |  | No |
| id | integer |  | No |
| message | string |  | No |
| processed | string |  | No |
| processed_by | string |  | No |
| processed_message | string |  | No |
| status | integer |  | No |
| user_id | string |  | No |
| user_tag | string |  | No |

#### models.User

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| avatar | string | The hash of the user's avatar. Use Session.UserAvatar to retrieve the avatar itself. | No |
| avatar_url | string |  | No |
| bot | boolean | Whether the user is a bot. | No |
| bot_owner | boolean |  | No |
| created_at | string |  | No |
| discriminator | string | The discriminator of the user (4 numbers after name). | No |
| email | string | The email of the user. This is only present when the application possesses the email scope for the user. | No |
| flags | integer | The flags on a user's account. Only available when the request is authorized via a Bearer token. | No |
| id | string | The ID of the user. | No |
| locale | string | The user's chosen language option. | No |
| mfa_enabled | boolean | Whether the user has multi-factor authentication enabled. | No |
| premium_type | integer | The type of Nitro subscription on a user's account. Only available when the request is authorized via a Bearer token. | No |
| public_flags | integer | The public flags on a user's account. This is a combination of bit masks; the presence of a certain flag can be checked by performing a bitwise AND between this int and the flag. | No |
| system | boolean | Whether the user is an Official Discord System user (part of the urgent message system). | No |
| token | string | The token of the user. This is only present for the user represented by the current session. | No |
| username | string | The user's username. | No |
| verified | boolean | Whether the user's email is verified. | No |

#### models.UsersettingsOTA

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| enabled | boolean |  | No |

#### presence.Presence

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| game | string |  | No |
| status | string |  | No |

#### shireikan.SubPermission

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string |  | No |
| explicit | boolean |  | No |
| term | string |  | No |
