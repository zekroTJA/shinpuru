# shinpuru main API
The shinpuru main REST API.

## Version: 1.0

## Etc
General root API functionalities.

### /allpermissions

#### GET
##### Summary

All Permissions

##### Description

Return a list of all available permissions.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ string ] |

### /healthcheck

#### GET
##### Summary

Healthcheck

##### Description

General system healthcheck.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ string ] |

### /me

#### GET
##### Summary

Me

##### Description

Returns the user object of the currently authenticated user.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.User](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsuser) |

### /privacyinfo

#### GET
##### Summary

Privacy Information

##### Description

Returns general global privacy information.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_models.Privacy](#github_com_zekrotja_shinpuru_internal_modelsprivacy) |

### /sysinfo

#### GET
##### Summary

System Information

##### Description

Returns general global system information.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.SystemInfo](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelssysteminfo) |

## Authorization
Authorization endpoints.

### /auth/accesstoken

#### POST
##### Summary

Access Token Exchange

##### Description

Exchanges the cookie-passed refresh token with a generated access token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.AccessTokenResponse](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsaccesstokenresponse) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /auth/check

#### GET
##### Summary

Authorization Check

##### Description

Returns OK if the request is authorized.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /auth/logout

#### POST
##### Summary

Logout

##### Description

Reovkes the currently used access token and clears the refresh token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |

### /auth/pushcode

#### POST
##### Summary

Pushcode

##### Description

Send a login push code resulting in a long-fetch request waiting for the code to be sent to shinpurus DMs.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| payload | body | The push code. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.PushCodeRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelspushcoderequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 410 | Gone | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |

## Channels
Channels specific endpoints.

### /channels/{guildid}

#### GET
##### Summary

Get Allowed Channels

##### Description

Returns a list of channels the user has access to.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| guildid | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Created | [discordgo.Message](#discordgomessage) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /channels/{guildid}/{id}

#### POST
##### Summary

Send Embed Message

##### Description

Send an Embed Message into a specified Channel.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| guildid | path | The ID of the guild. | Yes | string |
| id | path | The ID of the channel. | Yes | string |
| payload | body | The message embed object. | Yes | [discordgo.MessageEmbed](#discordgomessageembed) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Created | [discordgo.Message](#discordgomessage) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /channels/{guildid}/{id}/{msgid}

#### POST
##### Summary

Update Embed Message

##### Description

Update an Embed Message in a specified Channel with the given message ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| guildid | path | The ID of the guild. | Yes | string |
| id | path | The ID of the channel. | Yes | string |
| msgid | path | The ID of the message. | Yes | string |
| payload | body | The message embed object. | Yes | [discordgo.MessageEmbed](#discordgomessageembed) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [discordgo.Message](#discordgomessage) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

## Guilds
Guild specific endpoints.

### /guilds

#### GET
##### Summary

List Guilds

##### Description

Returns a list of guilds the authenticated user has in common with shinpuru.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.GuildReduced](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsguildreduced) ] |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}

#### GET
##### Summary

Get Guild

##### Description

Returns a single guild object by it's ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Guild](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsguild) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/antiraid/joinlog

#### GET
##### Summary

Get Antiraid Joinlog

##### Description

Returns a list of joined members during an antiraid trigger.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_models.JoinLogEntry](#github_com_zekrotja_shinpuru_internal_modelsjoinlogentry) ] |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### DELETE
##### Summary

Reset Antiraid Joinlog

##### Description

Deletes all entries of the antiraid joinlog.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/inviteblock

#### POST
##### Summary

Toggle Guild Inviteblock Enable

##### Description

Toggle enabled state of the guild invite block system.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The enable status payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.EnableStatus](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsenablestatus) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/permissions

#### GET
##### Summary

Get Guild Permission Settings

##### Description

Returns the specified guild permission settings.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.PermissionsMap](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelspermissionsmap) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Apply Guild Permission Rule

##### Description

Apply a new guild permission rule for a specified role.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The permission rule payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.PermissionsUpdate](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelspermissionsupdate) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.PermissionsMap](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelspermissionsmap) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/reports

#### GET
##### Summary

Get Guild Modlog

##### Description

Returns a list of guild modlog entries for the given guild.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| offset | query | The offset of returned entries | No | integer |
| limit | query | The amount of returned entries (0 = all) | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Report](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreport) ] |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/reports/count

#### GET
##### Summary

Get Guild Modlog Count

##### Description

Returns the total count of entries in the guild mod log.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Count](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelscount) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/scoreboard

#### GET
##### Summary

Get Guild Scoreboard

##### Description

Returns a list of scoreboard entries for the given guild.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| limit | query | Limit the amount of result values | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.GuildKarmaEntry](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsguildkarmaentry) ] |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/starboard

#### GET
##### Summary

Get Guild Starboard

##### Description

Returns a list of starboard entries for the given guild.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.StarboardEntryResponse](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstarboardentryresponse) ] |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/starboard/count

#### GET
##### Summary

Get Guild Starboard Count

##### Description

Returns the count of starboard entries for the given guild.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Count](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelscount) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/unbanrequests

#### GET
##### Summary

Get Guild Unbanrequests

##### Description

Returns the list of the guild unban requests.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListReponse | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.RichUnbanRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsrichunbanrequest) ] |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/unbanrequests/count

#### GET
##### Summary

Get Guild Unbanrequests Count

##### Description

Returns the total or filtered count of guild unban requests.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| state | query | Filter count by given state. | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Count](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelscount) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/unbanrequests/{requestid}

#### GET
##### Summary

Get Single Guild Unbanrequest

##### Description

Returns a single guild unban request by ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| requestid | path | The ID of the unbanrequest. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.RichUnbanRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsrichunbanrequest) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Process Guild Unbanrequest

##### Description

Process a guild unban request.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| requestid | path | The ID of the unbanrequest. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.RichUnbanRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsrichunbanrequest) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

## Guild Backups
Guild backup endpoints.

### /guilds/{id}/backups

#### GET
##### Summary

Get Guild Backups

##### Description

Returns a list of guild backups.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_services_backup_backupmodels.Entry](#github_com_zekrotja_shinpuru_internal_services_backup_backupmodelsentry) ] |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/backups/toggle

#### POST
##### Summary

Toggle Guild Backup Enable

##### Description

Toggle guild backup enable state.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | Enable state payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.EnableStatus](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsenablestatus) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/backups/{backupid}/download

#### GET
##### Summary

Download Backup File

##### Description

Download a single gziped backup file.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| backupid | path | The ID of the backup. | Yes | string |
| ota_token | query | The previously obtained OTA token to authorize the download. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | file |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 403 | Forbidden | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Obtain Backup Download OTA Key

##### Description

Returns an OTA key which is used to download a backup entry.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| backupid | path | The ID of the backup. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.AccessTokenResponse](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsaccesstokenresponse) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

## Members
Members specific endpoints.

### /guilds/{id}/members

#### GET
##### Summary

Get Guild Member List

##### Description

Returns a list of guild members.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| after | query | Request members after the given member ID. | No | string |
| limit | query | The amount of results returned. | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wraped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Member](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsmember) ] |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/{memberid}

#### GET
##### Summary

Get Guild Member

##### Description

Returns a single guild member by ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the member. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Member](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsmember) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/{memberid}/permissions

#### GET
##### Summary

Get Guild Member Permissions

##### Description

Returns the permission array of the given user.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the member. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.PermissionsResponse](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelspermissionsresponse) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/{memberid}/permissions/allowed

#### GET
##### Summary

Get Guild Member Allowed Permissions

##### Description

Returns all detailed permission DNS which the member is alloed to perform.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the member. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ string ] |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/{memberid}/reports

#### GET
##### Summary

Get Guild Member Reports

##### Description

Returns a list of reports of the given member.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the member. | Yes | string |
| limit | query | The amount of results returned. | No | integer |
| offset | query | The amount of results to be skipped. | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Report](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreport) ] |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Create A Member Report

##### Description

Creates a member report.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the victim member. | Yes | string |
| payload | body | The report payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.ReportRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreportrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Report](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreport) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/{memberid}/reports/count

#### GET
##### Summary

Get Guild Member Reports Count

##### Description

Returns the total count of reports of the given user.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the member. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Count](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelscount) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/{memberid}/unbanrequests

#### GET
##### Summary

Get Guild Member Unban Requests

##### Description

Returns the list of unban requests of the given member

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the member. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_models.UnbanRequest](#github_com_zekrotja_shinpuru_internal_modelsunbanrequest) ] |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/{memberid}/unbanrequests/count

#### GET
##### Summary

Get Guild Member Unban Requests Count

##### Description

Returns the total or filtered count of unban requests of the given member.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the member. | Yes | string |
| state | query | Filter unban requests by state. | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Count](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelscount) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

## Guild Settings
Guild specific settings endpoints.

### /guilds/{id}/settings

#### GET
##### Summary

Get Guild Settings

##### Description

Returns the specified general guild settings.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.GuildSettings](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsguildsettings) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Get Guild Settings

##### Description

Returns the specified general guild settings.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | Modified guild settings payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.GuildSettings](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsguildsettings) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/antiraid

#### GET
##### Summary

Get Guild Antiraid Settings

##### Description

Returns the specified guild antiraid settings.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.AntiraidSettings](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsantiraidsettings) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Update Guild Antiraid Settings

##### Description

Update the guild antiraid settings specification.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The guild antiraid settings payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.AntiraidSettings](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsantiraidsettings) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/antiraid/action

#### POST
##### Summary

Guild Antiraid Bulk Action

##### Description

Execute a specific action on antiraid listed users

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The antiraid action payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.AntiraidAction](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsantiraidaction) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/api

#### GET
##### Summary

Get Guild Settings API State

##### Description

Returns the settings state of the Guild API.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_models.GuildAPISettings](#github_com_zekrotja_shinpuru_internal_modelsguildapisettings) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Set Guild Settings API State

##### Description

Set the settings state of the Guild API.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The guild API settings payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.GuildAPISettingsRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsguildapisettingsrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_models.GuildAPISettings](#github_com_zekrotja_shinpuru_internal_modelsguildapisettings) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/codeexec

#### GET
##### Summary

Get Guild Settings Code Exec State

##### Description

Returns the settings state of the Guild Code Exec.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.EnableStatus](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsenablestatus) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Set Guild Settings Code Exec State

##### Description

Set the settings state of the Guild Code Exec.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The guild API settings payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.EnableStatus](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsenablestatus) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.EnableStatus](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsenablestatus) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/flushguilddata

#### POST
##### Summary

Flush Guild Data

##### Description

Flushes all guild data from the database.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The guild flush payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.FlushGuildRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsflushguildrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.State](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstate) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/karma

#### GET
##### Summary

Get Guild Karma Settings

##### Description

Returns the specified guild karma settings.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.KarmaSettings](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelskarmasettings) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Update Guild Karma Settings

##### Description

Update the guild karma settings specification.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The guild karma settings payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.KarmaSettings](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelskarmasettings) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/karma/blocklist

#### GET
##### Summary

Get Guild Karma Blocklist

##### Description

Returns the specified guild karma blocklist entries.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Member](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsmember) ] |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/karma/blocklist/{memberid}

#### PUT
##### Summary

Add Guild Karma Blocklist Entry

##### Description

Add a guild karma blocklist entry.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Member](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsmember) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### DELETE
##### Summary

Remove Guild Karma Blocklist Entry

##### Description

Remove a guild karma blocklist entry.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/karma/rules

#### GET
##### Summary

Get Guild Settings Karma Rules

##### Description

Returns a list of specified guild karma rules.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_models.KarmaRule](#github_com_zekrotja_shinpuru_internal_modelskarmarule) ] |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Create Guild Settings Karma

##### Description

Create a guild karma rule.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The karma rule payload. | Yes | [github_com_zekroTJA_shinpuru_internal_models.KarmaRule](#github_com_zekrotja_shinpuru_internal_modelskarmarule) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_models.KarmaRule](#github_com_zekrotja_shinpuru_internal_modelskarmarule) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/karma/rules/{ruleid}

#### POST
##### Summary

Update Guild Settings Karma

##### Description

Update a karma rule by ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| ruleid | path | The ID of the rule. | Yes | string |
| payload | body | The karma rule update payload. | Yes | [github_com_zekroTJA_shinpuru_internal_models.KarmaRule](#github_com_zekrotja_shinpuru_internal_modelskarmarule) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_models.KarmaRule](#github_com_zekrotja_shinpuru_internal_modelskarmarule) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### DELETE
##### Summary

Remove Guild Settings Karma

##### Description

Remove a guild karma rule by ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| ruleid | path | The ID of the rule. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.State](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstate) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/logs

#### GET
##### Summary

Get Guild Log Count

##### Description

Returns the total or filtered count of guild log entries.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| severity | query | Filter by log severity. | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Count](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelscount) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### DELETE
##### Summary

Delete Guild Log Entries

##### Description

Delete all guild log entries.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.State](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstate) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/logs/state

#### GET
##### Summary

Get Guild Settings Log State

##### Description

Returns the enabled state of the guild log setting.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.State](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstate) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Update Guild Settings Log State

##### Description

Update the enabled state of the log state guild setting.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The state payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.State](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstate) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.State](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstate) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/logs/{entryid}

#### DELETE
##### Summary

Delete Guild Log Entries

##### Description

Delete a single log entry.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| entryid | path | The ID of the entry to be deleted. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.State](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstate) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/settings/verification

#### GET
##### Summary

Get Guild Settings Verification State

##### Description

Returns the settings state of the Guild Verification.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.EnableStatus](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsenablestatus) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Set Guild Settings Verification State

##### Description

Set the settings state of the Guild Verification.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | The guild API settings payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.EnableStatus](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsenablestatus) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.EnableStatus](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsenablestatus) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

## Member Reporting
Member reporting endpoints.

### /guilds/{id}/{memberid}/ban

#### POST
##### Summary

Create A Member Ban Report

##### Description

Creates a member ban report.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the victim member. | Yes | string |
| payload | body | The report payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.ReasonRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreasonrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Report](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreport) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/{memberid}/kick

#### POST
##### Summary

Create A Member Kick Report

##### Description

Creates a member kick report.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the victim member. | Yes | string |
| payload | body | The report payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.ReasonRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreasonrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Report](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreport) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/{memberid}/mute

#### POST
##### Summary

Unmute A Member

##### Description

Unmute a muted member.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the victim member. | Yes | string |
| payload | body | The unmute payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.ReasonRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreasonrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/{memberid}/reports

#### GET
##### Summary

Get Guild Member Reports

##### Description

Returns a list of reports of the given member.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the member. | Yes | string |
| limit | query | The amount of results returned. | No | integer |
| offset | query | The amount of results to be skipped. | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Report](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreport) ] |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Create A Member Report

##### Description

Creates a member report.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the victim member. | Yes | string |
| payload | body | The report payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.ReportRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreportrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Report](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreport) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

## OTA
One Time Auth token endpoints.

### /ota

#### GET
##### Summary

OTA Login

##### Description

Logs in the current browser session by using the passed pre-obtained OTA token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK |  |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

## Public
Public API endpoints.

### /public/guilds/{id}

#### GET
##### Summary

Get Public Guild

##### Description

Returns public guild information, if enabled by guild config.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The Guild ID. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.GuildReduced](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsguildreduced) |

## Reports
General reports endpoints.

### /reports/{id}

#### GET
##### Summary

Get Report

##### Description

Returns a single report object by its ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The report ID. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Report](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreport) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /reports/{id}/revoke

#### POST
##### Summary

Revoke Report

##### Description

Revokes a given report by ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The report ID. | Yes | string |
| payload | body | The revoke reason payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.ReasonRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreasonrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Report](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsreport) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

## Search
Search endpoints.

### /search

#### GET
##### Summary

Global Search

##### Description

Search through guilds and members by ID, name or displayname.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| query | query | The search query (either ID, name or displayname). | Yes | string |
| limit | query | The maximum amount of result items (per group). | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.SearchResult](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelssearchresult) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

## Global Settings
Global bot settings endpoints.

### /settings/noguildinvite

#### GET
##### Summary

Get No Guild Invites Status

##### Description

Returns the settings status for the suggested guild invite when the logged in user is not on any guild with shinpuru.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.InviteSettingsResponse](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsinvitesettingsresponse) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 409 | Returned when no channel could be found to create invite for. | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Set No Guild Invites Status

##### Description

Set the status for the suggested guild invite when the logged in user is not on any guild with shinpuru.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| payload | body | Invite Settings Payload | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.InviteSettingsRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsinvitesettingsrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.APITokenResponse](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsapitokenresponse) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 409 | Returned when no channel could be found to create invite for. | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /settings/presence

#### GET
##### Summary

Get Presence

##### Description

Returns the bot's displayed presence status.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_util_presence.Presence](#github_com_zekrotja_shinpuru_internal_util_presencepresence) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Set Presence

##### Description

Set the bot's displayed presence status.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| payload | body | Presence Payload | Yes | [github_com_zekroTJA_shinpuru_internal_util_presence.Presence](#github_com_zekrotja_shinpuru_internal_util_presencepresence) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.APITokenResponse](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsapitokenresponse) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Is returned when no token was generated before. | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

## Tokens
API token endpoints.

### /token

#### GET
##### Summary

API Token Info

##### Description

Returns general metadata information about a generated API token. The response does **not** contain the actual token!

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.APITokenResponse](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsapitokenresponse) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Is returned when no token was generated before. | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

API Token Generation

##### Description

(Re-)Generates and returns general metadata information about an API token **including** the actual API token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.APITokenResponse](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsapitokenresponse) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### DELETE
##### Summary

API Token Deletion

##### Description

Invalidates the currently generated API token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

## Unban Requests
Unban requests endpoints.

### /unbanrequests

#### GET
##### Summary

Get Unban Requests

##### Description

Returns a list of unban requests created by the authenticated user.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.RichUnbanRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsrichunbanrequest) ] |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Create Unban Requests

##### Description

Create an unban reuqest.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| payload | body | The unban request payload. | Yes | [github_com_zekroTJA_shinpuru_internal_models.UnbanRequest](#github_com_zekrotja_shinpuru_internal_modelsunbanrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.RichUnbanRequest](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsrichunbanrequest) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /unbanrequests/bannedguilds

#### GET
##### Summary

Get Banned Guilds

##### Description

Returns a list of guilds where the currently authenticated user is banned.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.GuildReduced](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsguildreduced) ] |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

## default

### /users/{id}

#### GET
##### Summary

User

##### Description

Returns the information of a user by ID.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.User](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsuser) |

## User Settings
User specific settings endpoints.

### /usersettings/flush

#### POST
##### Summary

FLush all user data

##### Description

Flush all user data.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.UsersettingsOTA](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsusersettingsota) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /usersettings/ota

#### GET
##### Summary

Get OTA Usersettings State

##### Description

Returns the current state of the OTA user setting.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.UsersettingsOTA](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsusersettingsota) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Update OTA Usersettings State

##### Description

Update the OTA user settings state.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| payload | body | The OTA settings payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.UsersettingsOTA](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsusersettingsota) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.UsersettingsOTA](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsusersettingsota) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### /usersettings/privacy

#### GET
##### Summary

Get Privacy Usersettings

##### Description

Returns the current Privacy user settinga.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.UsersettingsPrivacy](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsusersettingsprivacy) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Update Privacy Usersettings

##### Description

Update the Privacy user settings.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| payload | body | The privacy settings payload. | Yes | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.UsersettingsPrivacy](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsusersettingsprivacy) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.UsersettingsPrivacy](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsusersettingsprivacy) |
| 400 | Bad Request | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

## Utilities
General utility functionalities.

### /util/color/{hexcode}

#### GET
##### Summary

Color Generator

##### Description

Produces a square image of the given color and size.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| hexcode | path | Hex Code of the Color to produce | Yes | string |
| size | query | The dimension of the square image | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | file |

### /util/landingpageinfo

#### GET
##### Summary

Landing Page Info

##### Description

Returns general information for the landing page like the local invite parameters.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.LandingPageResponse](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelslandingpageresponse) |

### /util/slashcommands

#### GET
##### Summary

Slash Command List

##### Description

Returns a list of registered slash commands and their description.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.SlashCommandInfo](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsslashcommandinfo) ] |

### /util/updateinfo

#### GET
##### Summary

Update Information

##### Description

Returns update information.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Update info response | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.UpdateInfoResponse](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsupdateinforesponse) |

## Verification
User verification endpoints.

### /verification/sitekey

#### GET
##### Summary

Sitekey

##### Description

Returns the sitekey for the captcha

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.CaptchaSiteKey](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelscaptchasitekey) |

### /verification/verify

#### POST
##### Summary

Verify

##### Description

Verify a returned verification token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsstatus) |
| 403 | Forbidden | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelserror) |

### Models

#### discordgo.ApplicationCommandOption

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| autocomplete | boolean | NOTE: mutually exclusive with Choices. | No |
| channel_types | [ [discordgo.ChannelType](#discordgochanneltype) ] |  | No |
| choices | [ [discordgo.ApplicationCommandOptionChoice](#discordgoapplicationcommandoptionchoice) ] |  | No |
| description | string |  | No |
| description_localizations | object |  | No |
| max_length | integer | Maximum length of string option. | No |
| max_value | number | Maximum value of number/integer option. | No |
| min_length | integer | Minimum length of string option. | No |
| min_value | number | Minimal value of number/integer option. | No |
| name | string |  | No |
| name_localizations | object |  | No |
| options | [ [discordgo.ApplicationCommandOption](#discordgoapplicationcommandoption) ] |  | No |
| required | boolean |  | No |
| type | [discordgo.ApplicationCommandOptionType](#discordgoapplicationcommandoptiontype) |  | No |

#### discordgo.ApplicationCommandOptionChoice

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string |  | No |
| name_localizations | object |  | No |
| value |  |  | No |

#### discordgo.ApplicationCommandOptionType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.ApplicationCommandOptionType | integer |  |  |

#### discordgo.Channel

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| application_id | string | ApplicationID of the DM creator Zeroed if guild channel or not a bot user | No |
| applied_tags | [ string ] | The IDs of the set of tags that have been applied to a thread in a forum channel. | No |
| available_tags | [ [discordgo.ForumTag](#discordgoforumtag) ] | The set of tags that can be used in a forum channel. | No |
| bitrate | integer | The bitrate of the channel, if it is a voice channel. | No |
| default_forum_layout | [discordgo.ForumLayout](#discordgoforumlayout) | The default forum layout view used to display posts in forum channels. Defaults to ForumLayoutNotSet, which indicates a layout view has not been set by a channel admin. | No |
| default_reaction_emoji | [discordgo.ForumDefaultReaction](#discordgoforumdefaultreaction) | Emoji to use as the default reaction to a forum post. | No |
| default_sort_order | [discordgo.ForumSortOrderType](#discordgoforumsortordertype) | The default sort order type used to order posts in forum channels. Defaults to null, which indicates a preferred sort order hasn't been set by a channel admin. | No |
| default_thread_rate_limit_per_user | integer | The initial RateLimitPerUser to set on newly created threads in a channel. This field is copied to the thread at creation time and does not live update. | No |
| flags | [discordgo.ChannelFlags](#discordgochannelflags) | Channel flags. | No |
| guild_id | string | The ID of the guild to which the channel belongs, if it is in a guild. Else, this ID is empty (e.g. DM channels). | No |
| icon | string | Icon of the group DM channel. | No |
| id | string | The ID of the channel. | No |
| last_message_id | string | The ID of the last message sent in the channel. This is not guaranteed to be an ID of a valid message. | No |
| last_pin_timestamp | string | The timestamp of the last pinned message in the channel. nil if the channel has no pinned messages. | No |
| member_count | integer | An approximate count of users in a thread, stops counting at 50 | No |
| message_count | integer | An approximate count of messages in a thread, stops counting at 50 | No |
| name | string | The name of the channel. | No |
| nsfw | boolean | Whether the channel is marked as NSFW. | No |
| owner_id | string | ID of the creator of the group DM or thread | No |
| parent_id | string | The ID of the parent channel, if the channel is under a category. For threads - id of the channel thread was created in. | No |
| permission_overwrites | [ [discordgo.PermissionOverwrite](#discordgopermissionoverwrite) ] | A list of permission overwrites present for the channel. | No |
| position | integer | The position of the channel, used for sorting in client. | No |
| rate_limit_per_user | integer | Amount of seconds a user has to wait before sending another message or creating another thread (0-21600) bots, as well as users with the permission manage_messages or manage_channel, are unaffected | No |
| recipients | [ [discordgo.User](#discordgouser) ] | The recipients of the channel. This is only populated in DM channels. | No |
| thread_member | [discordgo.ThreadMember](#discordgothreadmember) | Thread member object for the current user, if they have joined the thread, only included on certain API endpoints | No |
| thread_metadata | [discordgo.ThreadMetadata](#discordgothreadmetadata) | Thread-specific fields not needed by other channels | No |
| topic | string | The topic of the channel. | No |
| type | [discordgo.ChannelType](#discordgochanneltype) | The type of the channel. | No |
| user_limit | integer | The user limit of the voice channel. | No |

#### discordgo.ChannelFlags

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.ChannelFlags | integer |  |  |

#### discordgo.ChannelType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.ChannelType | integer |  |  |

#### discordgo.EmbedType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.EmbedType | string |  |  |

#### discordgo.Emoji

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| animated | boolean |  | No |
| available | boolean |  | No |
| id | string |  | No |
| managed | boolean |  | No |
| name | string |  | No |
| require_colons | boolean |  | No |
| roles | [ string ] |  | No |
| user | [discordgo.User](#discordgouser) |  | No |

#### discordgo.ForumDefaultReaction

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| emoji_id | string | The id of a guild's custom emoji. | No |
| emoji_name | string | The unicode character of the emoji. | No |

#### discordgo.ForumLayout

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.ForumLayout | integer |  |  |

#### discordgo.ForumSortOrderType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.ForumSortOrderType | integer |  |  |

#### discordgo.ForumTag

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| emoji_id | string |  | No |
| emoji_name | string |  | No |
| id | string |  | No |
| moderated | boolean |  | No |
| name | string |  | No |

#### discordgo.InteractionType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.InteractionType | integer |  |  |

#### discordgo.Member

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| avatar | string | The hash of the avatar for the guild member, if any. | No |
| communication_disabled_until | string | The time at which the member's timeout will expire. Time in the past or nil if the user is not timed out. | No |
| deaf | boolean | Whether the member is deafened at a guild level. | No |
| guild_id | string | The guild ID on which the member exists. | No |
| joined_at | string | The time at which the member joined the guild. | No |
| mute | boolean | Whether the member is muted at a guild level. | No |
| nick | string | The nickname of the member, if they have one. | No |
| pending | boolean | Is true while the member hasn't accepted the membership screen. | No |
| permissions | string | Total permissions of the member in the channel, including overrides, returned when in the interaction object.<br>*Example:* `"0"` | No |
| premium_since | string | When the user used their Nitro boost on the server | No |
| roles | [ string ] | A list of IDs of the roles which are possessed by the member. | No |
| user | [discordgo.User](#discordgouser) | The underlying user on which the member is based. | No |

#### discordgo.Message

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| activity | [discordgo.MessageActivity](#discordgomessageactivity) | Is sent with Rich Presence-related chat embeds | No |
| application | [discordgo.MessageApplication](#discordgomessageapplication) | Is sent with Rich Presence-related chat embeds | No |
| attachments | [ [discordgo.MessageAttachment](#discordgomessageattachment) ] | A list of attachments present in the message. | No |
| author | [discordgo.User](#discordgouser) | The author of the message. This is not guaranteed to be a valid user (webhook-sent messages do not possess a full author). | No |
| channel_id | string | The ID of the channel in which the message was sent. | No |
| content | string | The content of the message. | No |
| edited_timestamp | string | The time at which the last edit of the message occurred, if it has been edited. | No |
| embeds | [ [discordgo.MessageEmbed](#discordgomessageembed) ] | A list of embeds present in the message. | No |
| flags | [discordgo.MessageFlags](#discordgomessageflags) | The flags of the message, which describe extra features of a message. This is a combination of bit masks; the presence of a certain permission can be checked by performing a bitwise AND between this int and the flag. | No |
| guild_id | string | The ID of the guild in which the message was sent. | No |
| id | string | The ID of the message. | No |
| interaction | [discordgo.MessageInteraction](#discordgomessageinteraction) | Is sent when the message is a response to an Interaction, without an existing message. This means responses to message component interactions do not include this property, instead including a MessageReference, as components exist on preexisting messages. | No |
| member | [discordgo.Member](#discordgomember) | Member properties for this message's author, contains only partial information | No |
| mention_channels | [ [discordgo.Channel](#discordgochannel) ] | Channels specifically mentioned in this message Not all channel mentions in a message will appear in mention_channels. Only textual channels that are visible to everyone in a lurkable guild will ever be included. Only crossposted messages (via Channel Following) currently include mention_channels at all. If no mentions in the message meet these requirements, this field will not be sent. | No |
| mention_everyone | boolean | Whether the message mentions everyone. | No |
| mention_roles | [ string ] | The roles mentioned in the message. | No |
| mentions | [ [discordgo.User](#discordgouser) ] | A list of users mentioned in the message. | No |
| message_reference | [discordgo.MessageReference](#discordgomessagereference) | MessageReference contains reference data sent with crossposted or reply messages. This does not contain the reference *to* this message; this is for when *this* message references another. To generate a reference to this message, use (*Message).Reference(). | No |
| pinned | boolean | Whether the message is pinned or not. | No |
| reactions | [ [discordgo.MessageReactions](#discordgomessagereactions) ] | A list of reactions to the message. | No |
| referenced_message | [discordgo.Message](#discordgomessage) | The message associated with the message_reference NOTE: This field is only returned for messages with a type of 19 (REPLY) or 21 (THREAD_STARTER_MESSAGE). If the message is a reply but the referenced_message field is not present, the backend did not attempt to fetch the message that was being replied to, so its state is unknown. If the field exists but is null, the referenced message was deleted. | No |
| sticker_items | [ [discordgo.Sticker](#discordgosticker) ] | An array of Sticker objects, if any were sent. | No |
| thread | [discordgo.Channel](#discordgochannel) | The thread that was started from this message, includes thread member object | No |
| timestamp | string | The time at which the messsage was sent. CAUTION: this field may be removed in a future API version; it is safer to calculate the creation time via the ID. | No |
| tts | boolean | Whether the message is text-to-speech. | No |
| type | [discordgo.MessageType](#discordgomessagetype) | The type of the message. | No |
| webhook_id | string | The webhook ID of the message, if it was generated by a webhook | No |

#### discordgo.MessageActivity

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| party_id | string |  | No |
| type | [discordgo.MessageActivityType](#discordgomessageactivitytype) |  | No |

#### discordgo.MessageActivityType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.MessageActivityType | integer |  |  |

#### discordgo.MessageApplication

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cover_image | string |  | No |
| description | string |  | No |
| icon | string |  | No |
| id | string |  | No |
| name | string |  | No |

#### discordgo.MessageAttachment

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| content_type | string |  | No |
| ephemeral | boolean |  | No |
| filename | string |  | No |
| height | integer |  | No |
| id | string |  | No |
| proxy_url | string |  | No |
| size | integer |  | No |
| url | string |  | No |
| width | integer |  | No |

#### discordgo.MessageEmbed

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| author | [discordgo.MessageEmbedAuthor](#discordgomessageembedauthor) |  | No |
| color | integer |  | No |
| description | string |  | No |
| fields | [ [discordgo.MessageEmbedField](#discordgomessageembedfield) ] |  | No |
| footer | [discordgo.MessageEmbedFooter](#discordgomessageembedfooter) |  | No |
| image | [discordgo.MessageEmbedImage](#discordgomessageembedimage) |  | No |
| provider | [discordgo.MessageEmbedProvider](#discordgomessageembedprovider) |  | No |
| thumbnail | [discordgo.MessageEmbedThumbnail](#discordgomessageembedthumbnail) |  | No |
| timestamp | string |  | No |
| title | string |  | No |
| type | [discordgo.EmbedType](#discordgoembedtype) |  | No |
| url | string |  | No |
| video | [discordgo.MessageEmbedVideo](#discordgomessageembedvideo) |  | No |

#### discordgo.MessageEmbedAuthor

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| icon_url | string |  | No |
| name | string |  | No |
| proxy_icon_url | string |  | No |
| url | string |  | No |

#### discordgo.MessageEmbedField

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| inline | boolean |  | No |
| name | string |  | No |
| value | string |  | No |

#### discordgo.MessageEmbedFooter

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| icon_url | string |  | No |
| proxy_icon_url | string |  | No |
| text | string |  | No |

#### discordgo.MessageEmbedImage

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| height | integer |  | No |
| proxy_url | string |  | No |
| url | string |  | No |
| width | integer |  | No |

#### discordgo.MessageEmbedProvider

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string |  | No |
| url | string |  | No |

#### discordgo.MessageEmbedThumbnail

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| height | integer |  | No |
| proxy_url | string |  | No |
| url | string |  | No |
| width | integer |  | No |

#### discordgo.MessageEmbedVideo

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| height | integer |  | No |
| url | string |  | No |
| width | integer |  | No |

#### discordgo.MessageFlags

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.MessageFlags | integer |  |  |

#### discordgo.MessageInteraction

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string |  | No |
| member | [discordgo.Member](#discordgomember) | Member is only present when the interaction is from a guild. | No |
| name | string |  | No |
| type | [discordgo.InteractionType](#discordgointeractiontype) |  | No |
| user | [discordgo.User](#discordgouser) |  | No |

#### discordgo.MessageReactions

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| count | integer |  | No |
| emoji | [discordgo.Emoji](#discordgoemoji) |  | No |
| me | boolean |  | No |

#### discordgo.MessageReference

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| channel_id | string |  | No |
| guild_id | string |  | No |
| message_id | string |  | No |

#### discordgo.MessageType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.MessageType | integer |  |  |

#### discordgo.MfaLevel

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.MfaLevel | integer |  |  |

#### discordgo.PermissionOverwrite

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| allow | string | *Example:* `"0"` | No |
| deny | string | *Example:* `"0"` | No |
| id | string |  | No |
| type | [discordgo.PermissionOverwriteType](#discordgopermissionoverwritetype) |  | No |

#### discordgo.PermissionOverwriteType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.PermissionOverwriteType | integer |  |  |

#### discordgo.PremiumTier

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.PremiumTier | integer |  |  |

#### discordgo.Role

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| color | integer | The hex color of this role. | No |
| hoist | boolean | Whether this role is hoisted (shows up separately in member list). | No |
| id | string | The ID of the role. | No |
| managed | boolean | Whether this role is managed by an integration, and thus cannot be manually added to, or taken from, members. | No |
| mentionable | boolean | Whether this role is mentionable. | No |
| name | string | The name of the role. | No |
| permissions | string | The permissions of the role on the guild (doesn't include channel overrides). This is a combination of bit masks; the presence of a certain permission can be checked by performing a bitwise AND between this int and the permission.<br>*Example:* `"0"` | No |
| position | integer | The position of this role in the guild's role hierarchy. | No |

#### discordgo.Sticker

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| available | boolean |  | No |
| description | string |  | No |
| format_type | [discordgo.StickerFormat](#discordgostickerformat) |  | No |
| guild_id | string |  | No |
| id | string |  | No |
| name | string |  | No |
| pack_id | string |  | No |
| sort_value | integer |  | No |
| tags | string |  | No |
| type | [discordgo.StickerType](#discordgostickertype) |  | No |
| user | [discordgo.User](#discordgouser) |  | No |

#### discordgo.StickerFormat

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.StickerFormat | integer |  |  |

#### discordgo.StickerType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.StickerType | integer |  |  |

#### discordgo.ThreadMember

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| flags | integer | Any user-thread settings, currently only used for notifications | No |
| id | string | The id of the thread | No |
| join_timestamp | string | The time the current user last joined the thread | No |
| user_id | string | The id of the user | No |

#### discordgo.ThreadMetadata

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| archive_timestamp | string | Timestamp when the thread's archive status was last changed, used for calculating recent activity | No |
| archived | boolean | Whether the thread is archived | No |
| auto_archive_duration | integer | Duration in minutes to automatically archive the thread after recent activity, can be set to: 60, 1440, 4320, 10080 | No |
| invitable | boolean | Whether non-moderators can add other non-moderators to a thread; only available on private threads | No |
| locked | boolean | Whether the thread is locked; when a thread is locked, only users with MANAGE_THREADS can unarchive it | No |

#### discordgo.User

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| accent_color | integer | User's banner color, encoded as an integer representation of hexadecimal color code | No |
| avatar | string | The hash of the user's avatar. Use Session.UserAvatar to retrieve the avatar itself. | No |
| banner | string | The hash of the user's banner image. | No |
| bot | boolean | Whether the user is a bot. | No |
| discriminator | string | The discriminator of the user (4 numbers after name). | No |
| email | string | The email of the user. This is only present when the application possesses the email scope for the user. | No |
| flags | integer | The flags on a user's account. Only available when the request is authorized via a Bearer token. | No |
| id | string | The ID of the user. | No |
| locale | string | The user's chosen language option. | No |
| mfa_enabled | boolean | Whether the user has multi-factor authentication enabled. | No |
| premium_type | integer | The type of Nitro subscription on a user's account. Only available when the request is authorized via a Bearer token. | No |
| public_flags | [discordgo.UserFlags](#discordgouserflags) | The public flags on a user's account. This is a combination of bit masks; the presence of a certain flag can be checked by performing a bitwise AND between this int and the flag. | No |
| system | boolean | Whether the user is an Official Discord System user (part of the urgent message system). | No |
| token | string | The token of the user. This is only present for the user represented by the current session. | No |
| username | string | The user's username. | No |
| verified | boolean | Whether the user's email is verified. | No |

#### discordgo.UserFlags

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.UserFlags | integer |  |  |

#### discordgo.VerificationLevel

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.VerificationLevel | integer |  |  |

#### github_com_zekroTJA_shinpuru_internal_models.Contact

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| title | string |  | No |
| url | string |  | No |
| value | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_models.GuildAPISettings

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| allowed_origins | string |  | No |
| enabled | boolean |  | No |
| protected | boolean |  | No |
| token_hash | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_models.GuildLogEntry

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| guildid | string |  | No |
| id | integer |  | No |
| message | string |  | No |
| module | string |  | No |
| severity | [github_com_zekroTJA_shinpuru_internal_models.GuildLogSeverity](#github_com_zekrotja_shinpuru_internal_modelsguildlogseverity) |  | No |
| timestamp | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_models.GuildLogSeverity

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| github_com_zekroTJA_shinpuru_internal_models.GuildLogSeverity | integer |  |  |

#### github_com_zekroTJA_shinpuru_internal_models.JoinLogEntry

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| account_created | string |  | No |
| guild_id | string |  | No |
| tag | string |  | No |
| timestamp | string |  | No |
| user_id | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_models.KarmaAction

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| github_com_zekroTJA_shinpuru_internal_models.KarmaAction | string |  |  |

#### github_com_zekroTJA_shinpuru_internal_models.KarmaRule

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| action | [github_com_zekroTJA_shinpuru_internal_models.KarmaAction](#github_com_zekrotja_shinpuru_internal_modelskarmaaction) |  | No |
| argument | string |  | No |
| guildid | string |  | No |
| id | integer |  | No |
| trigger | [github_com_zekroTJA_shinpuru_internal_models.KarmaTriggerType](#github_com_zekrotja_shinpuru_internal_modelskarmatriggertype) |  | No |
| value | integer |  | No |

#### github_com_zekroTJA_shinpuru_internal_models.KarmaTriggerType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| github_com_zekroTJA_shinpuru_internal_models.KarmaTriggerType | integer |  |  |

#### github_com_zekroTJA_shinpuru_internal_models.Privacy

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| contact | [ [github_com_zekroTJA_shinpuru_internal_models.Contact](#github_com_zekrotja_shinpuru_internal_modelscontact) ] |  | No |
| noticeurl | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_models.ReportType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| github_com_zekroTJA_shinpuru_internal_models.ReportType | integer |  |  |

#### github_com_zekroTJA_shinpuru_internal_models.UnbanRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created | string |  | No |
| guild_id | string |  | No |
| id | integer |  | No |
| message | string |  | No |
| processed | string |  | No |
| processed_by | string |  | No |
| processed_message | string |  | No |
| status | [github_com_zekroTJA_shinpuru_internal_models.UnbanRequestState](#github_com_zekrotja_shinpuru_internal_modelsunbanrequeststate) |  | No |
| user_id | string |  | No |
| user_tag | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_models.UnbanRequestState

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| github_com_zekroTJA_shinpuru_internal_models.UnbanRequestState | integer |  |  |

#### github_com_zekroTJA_shinpuru_internal_services_backup_backupmodels.Entry

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| file_id | string |  | No |
| guild_id | string |  | No |
| timestamp | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_permissions.SubPermission

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string |  | No |
| explicit | boolean |  | No |
| term | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.APITokenResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created | string |  | No |
| expires | string |  | No |
| hits | integer |  | No |
| last_access | string |  | No |
| token | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.AccessTokenResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| expires | string |  | No |
| token | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.AntiraidAction

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| ids | [ string ] |  | No |
| type | integer |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.AntiraidSettings

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| burst | integer |  | No |
| regeneration_period | integer |  | No |
| state | boolean |  | No |
| verification | boolean |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.CaptchaSiteKey

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| sitekey | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Count

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| count | integer |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.EnableStatus

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| enabled | boolean |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Error

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | integer |  | No |
| context | string |  | No |
| error | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.FlatUser

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| avatar_url | string |  | No |
| bot | boolean |  | No |
| discriminator | string |  | No |
| id | string |  | No |
| username | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.FlushGuildRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| leave_after | boolean |  | No |
| validation | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Guild

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
| mfa_level | [discordgo.MfaLevel](#discordgomfalevel) |  | No |
| name | string |  | No |
| owner_id | string |  | No |
| premium_subscription_count | integer |  | No |
| premium_tier | [discordgo.PremiumTier](#discordgopremiumtier) |  | No |
| region | string |  | No |
| roles | [ [discordgo.Role](#discordgorole) ] |  | No |
| self_member | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Member](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsmember) |  | No |
| splash | string |  | No |
| unavailable | boolean |  | No |
| verification_level | [discordgo.VerificationLevel](#discordgoverificationlevel) |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.GuildAPISettingsRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| allowed_origins | string |  | No |
| enabled | boolean |  | No |
| protected | boolean |  | No |
| reset_token | boolean |  | No |
| token | string |  | No |
| token_hash | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.GuildKarmaEntry

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| member | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Member](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsmember) |  | No |
| value | integer |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.GuildReduced

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| icon | string |  | No |
| icon_url | string |  | No |
| id | string |  | No |
| joined_at | string |  | No |
| member_count | integer |  | No |
| name | string |  | No |
| online_member_count | integer |  | No |
| owner_id | string |  | No |
| region | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.GuildSettings

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

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.InviteSettingsRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| guild_id | string |  | No |
| invite_code | string |  | No |
| message | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.InviteSettingsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| guild | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Guild](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsguild) |  | No |
| invite_url | string |  | No |
| message | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.KarmaSettings

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| emotes_decrease | [ string ] |  | No |
| emotes_increase | [ string ] |  | No |
| penalty | boolean |  | No |
| state | boolean |  | No |
| tokens | integer |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.LandingPageResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| localinvite | string |  | No |
| publiccaranyinvite | string |  | No |
| publicmaininvite | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Member

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| avatar | string | The hash of the avatar for the guild member, if any. | No |
| avatar_url | string |  | No |
| chat_muted | boolean |  | No |
| communication_disabled_until | string | The time at which the member's timeout will expire. Time in the past or nil if the user is not timed out. | No |
| created_at | string |  | No |
| deaf | boolean | Whether the member is deafened at a guild level. | No |
| dominance | integer |  | No |
| guild_id | string | The guild ID on which the member exists. | No |
| guild_name | string |  | No |
| joined_at | string | The time at which the member joined the guild. | No |
| karma | integer |  | No |
| karma_total | integer |  | No |
| mute | boolean | Whether the member is muted at a guild level. | No |
| nick | string | The nickname of the member, if they have one. | No |
| pending | boolean | Is true while the member hasn't accepted the membership screen. | No |
| permissions | string | Total permissions of the member in the channel, including overrides, returned when in the interaction object.<br>*Example:* `"0"` | No |
| premium_since | string | When the user used their Nitro boost on the server | No |
| roles | [ string ] | A list of IDs of the roles which are possessed by the member. | No |
| user | [discordgo.User](#discordgouser) | The underlying user on which the member is based. | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.PermissionsMap

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.PermissionsMap | object |  |  |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.PermissionsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| permissions | [ string ] |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.PermissionsUpdate

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| override | boolean |  | No |
| perm | string |  | No |
| role_ids | [ string ] |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.PushCodeRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.ReasonRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| attachment | string |  | No |
| attachment_data | string |  | No |
| reason | string |  | No |
| timeout | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Report

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| attachment_url | string |  | No |
| created | string |  | No |
| executor | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.FlatUser](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsflatuser) |  | No |
| executor_id | string |  | No |
| guild_id | string |  | No |
| id | integer |  | No |
| message | string |  | No |
| timeout | string |  | No |
| type | [github_com_zekroTJA_shinpuru_internal_models.ReportType](#github_com_zekrotja_shinpuru_internal_modelsreporttype) |  | No |
| type_name | string |  | No |
| victim | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.FlatUser](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsflatuser) |  | No |
| victim_id | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.ReportRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| attachment | string |  | No |
| attachment_data | string |  | No |
| reason | string |  | No |
| timeout | string |  | No |
| type | [github_com_zekroTJA_shinpuru_internal_models.ReportType](#github_com_zekrotja_shinpuru_internal_modelsreporttype) |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.RichUnbanRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created | string |  | No |
| creator | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.FlatUser](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsflatuser) |  | No |
| guild_id | string |  | No |
| id | integer |  | No |
| message | string |  | No |
| processed | string |  | No |
| processed_by | string |  | No |
| processed_message | string |  | No |
| processor | [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.FlatUser](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsflatuser) |  | No |
| status | [github_com_zekroTJA_shinpuru_internal_models.UnbanRequestState](#github_com_zekrotja_shinpuru_internal_modelsunbanrequeststate) |  | No |
| user_id | string |  | No |
| user_tag | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.SearchResult

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| guilds | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.GuildReduced](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsguildreduced) ] |  | No |
| members | [ [github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Member](#github_com_zekrotja_shinpuru_internal_services_webserver_v1_modelsmember) ] |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.SlashCommandInfo

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string |  | No |
| dm_capable | boolean |  | No |
| domain | string |  | No |
| group | string |  | No |
| name | string |  | No |
| options | [ [discordgo.ApplicationCommandOption](#discordgoapplicationcommandoption) ] |  | No |
| subdomains | [ [github_com_zekroTJA_shinpuru_internal_services_permissions.SubPermission](#github_com_zekrotja_shinpuru_internal_services_permissionssubpermission) ] |  | No |
| version | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.StarboardEntryResponse

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

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.State

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| state | boolean |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.Status

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | integer |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.SystemInfo

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

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.UpdateInfoResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| current | [github_com_zekroTJA_shinpuru_pkg_versioncheck.Semver](#github_com_zekrotja_shinpuru_pkg_versionchecksemver) |  | No |
| current_str | string |  | No |
| isold | boolean |  | No |
| latest | [github_com_zekroTJA_shinpuru_pkg_versioncheck.Semver](#github_com_zekrotja_shinpuru_pkg_versionchecksemver) |  | No |
| latest_str | string |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.User

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| accent_color | integer | User's banner color, encoded as an integer representation of hexadecimal color code | No |
| avatar | string | The hash of the user's avatar. Use Session.UserAvatar to retrieve the avatar itself. | No |
| avatar_url | string |  | No |
| banner | string | The hash of the user's banner image. | No |
| bot | boolean | Whether the user is a bot. | No |
| bot_owner | boolean |  | No |
| captcha_verified | boolean |  | No |
| created_at | string |  | No |
| discriminator | string | The discriminator of the user (4 numbers after name). | No |
| email | string | The email of the user. This is only present when the application possesses the email scope for the user. | No |
| flags | integer | The flags on a user's account. Only available when the request is authorized via a Bearer token. | No |
| id | string | The ID of the user. | No |
| locale | string | The user's chosen language option. | No |
| mfa_enabled | boolean | Whether the user has multi-factor authentication enabled. | No |
| premium_type | integer | The type of Nitro subscription on a user's account. Only available when the request is authorized via a Bearer token. | No |
| public_flags | [discordgo.UserFlags](#discordgouserflags) | The public flags on a user's account. This is a combination of bit masks; the presence of a certain flag can be checked by performing a bitwise AND between this int and the flag. | No |
| system | boolean | Whether the user is an Official Discord System user (part of the urgent message system). | No |
| token | string | The token of the user. This is only present for the user represented by the current session. | No |
| username | string | The user's username. | No |
| verified | boolean | Whether the user's email is verified. | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.UsersettingsOTA

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| enabled | boolean |  | No |

#### github_com_zekroTJA_shinpuru_internal_services_webserver_v1_models.UsersettingsPrivacy

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| starboard_optout | boolean |  | No |

#### github_com_zekroTJA_shinpuru_internal_util_presence.Presence

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| game | string |  | No |
| status | [github_com_zekroTJA_shinpuru_internal_util_presence.Status](#github_com_zekrotja_shinpuru_internal_util_presencestatus) |  | No |

#### github_com_zekroTJA_shinpuru_internal_util_presence.Status

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| github_com_zekroTJA_shinpuru_internal_util_presence.Status | string |  |  |

#### github_com_zekroTJA_shinpuru_pkg_versioncheck.Semver

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| attachment | string |  | No |
| major | integer |  | No |
| minor | integer |  | No |
| patch | integer |  | No |
