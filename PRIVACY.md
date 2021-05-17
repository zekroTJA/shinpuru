# Privacy

This document explains which data is collected and stored by shinpuru, why it is collected and what you can do to delete your data from shinpuru.

## Wich data is collected by shinpuru?

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

Data is stored in in 3 primary levels in the shinpuru stack.

The first level are chaches inside the applications memory. Data in this level is only stored for a short time period and not persistent.