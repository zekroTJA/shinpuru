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
| size | query | The dimension of the square image (default: 24) | No | integer |

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

#### models.Error

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | integer |  | No |
| context | string |  | No |
| error | string |  | No |

#### models.LandingPageResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| localinvite | string |  | No |
| publiccaranyinvite | string |  | No |
| publicmaininvite | string |  | No |

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

#### shireikan.SubPermission

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string |  | No |
| explicit | boolean |  | No |
| term | string |  | No |
