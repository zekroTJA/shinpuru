# Privacy

This document explains which data is collected and stored by shinpuru, why it is collected and what you can do to delete your data from shinpuru.

## Which data is collected by shinpuru?

The following data is stored persistently which can be directly linked to your guild and/or user account.

- User IDs
- Guild IDs
- Channel IDs
- Role IDs
- Message IDs
- Message Contents*
- Message Attachments*

**Message contents and media attachments are only stored for specific purposes and not in a general way. No message contents or message attachments are logged in any way.*

## Why is this data stored?

First of all, Guild IDs are saved to store settings, preferences and permissions in combination with a guild. Channel IDs are stored for setting destinations for join/leave mesasges, voice log, mod log and the starboard. Role IDs need to be stored to map permission rules with assigned roles to users.

Message IDs are stored in combination with Channel and Guild IDs to link to vote messages or to the original messages of starboard messages. When a message is presented in the starboard, the message content is stored to the database to present it in the web interface.

User IDs are explicitly stored for API tokens, karma points, the karma blocklist, the antiraid joinlog, unban requests and refresh tokens *(for the web authentication system)*.

User IDs and message IDs are also stored temporarily for limiting and ensuring functionality of some features like the karma or starboard system. As mentioned, the storage is only temporarily and time limited.

## How and where is this data stored?

Data is stored in in 3 primary "levels" in the shinpuru stack.

The first level are caches inside the applications memory. Data in this level is only stored for a short time period and not persistent. This layer is mostly used for caching, rate limiting and scheduling tasks.

The second layer is another deeper cache layer which sits between the database and the application and mirrors parts of the database to reduce access times of highly frequently requested data from the database and also reducing the load on the database as well. This data is only stored during the runtime of the redis instance, so it is semi-persistent. Also, this cache is used for multiple instances of shinpuru when sharding is used.

The third layer is the central MariaDB database of shinpuru where all persistent data is stored permanently.

If you are interested in details where which data is stored in which way, take a look at the [database implementation of shinpuru](https://github.com/zekroTJA/shinpuru/tree/master/internal/services/database).

## How can I remove my data?

If you want to delete all guild data, this includes all reports and associated image data; all backups and associated backup files; all karma scores, settings, rules and blocklist; all starboard entries and configuration; all guild settings and permission specifications; tags; antiraid settings and joinlog and all unban requests; you can do it in the guild settings in the web interface at `/guilds/:guildid/guildadmin/data`.

To do so, you need the `sp.guild.admin.flushdata` permission.

![](https://i.imgur.com/vI2J0k9.png)

Otherwise, if you want to have specific user data removed, please contact me via mail to [privacy[at]zekro.de](mailto:privacy@zekro.de). Please send by any type of authentication that the account ID you privide is actually yours. That can be a screen-shot or I will contact you via Discord directly, if you want. Also, please describe as detailed as possible which data you want to have removed.

## Any questions?

If you have any questions or privacy concerns, feel free to contact me. 😊  
https://zekro.de/contact