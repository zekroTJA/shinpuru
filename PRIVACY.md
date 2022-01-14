# Privacy Notice

This document explains which data is collected and stored by shinpuru, why it is collected and what you can do to delete your data from shinpuru.

The main instances of shinpuru and it's databases are running on a VPS hosted by [contabo.de](https://contabo.de/) in Germany. The server is hosted by me, Ringo Hoffmann aka. zekro (contact details below). 

**This privacy notice only refers to the main instances of shinpuru (`shinpuru#4878` / https://shnp.de and `shinpuru Canary#3479` / https://c.shnp.de) and excludes all 3rd party hosted instances and forks of shinpuru.**

## Contact details

Ringo Hoffmann

Mail: contact@zekro.de  
Twitter: https://twitter.com/zekrotja

You can also contact me on my Discord: `zekro#0001` @ https://discord.zekro.de

If you need my address, please refer to my Imprint: https://www.zekro.de/imprint

## Which data is collected and for what purpose?

All data collected by shinpuru either originates from direct user input or via the Discord REST and Gateway API. 

There are two main purposes data is collected and stored in shinpuru:

### A) Persistence

Some data is stored **persistently (permanently)** in the database of shinpuru to ensure a good user experience as well as maintaining state for security features implemented in shinpuru.

This includes the following data:
- User IDs
- Message content¹
- User input text
- User uploaded media

¹Message content is stored only when a message has been voted into the starboard. The starboard is also available in the web interface and therefore, the content of the original message is stored to display it in the web interface without overloading the Discord API. The users are abloe to globally opt-out from the starboard via the privacy settings in the unser settings in the web interface.

User IDs are used to link stored data to Discord users like:
- API and refresh tokens for authentication against shinpuru's API (for example via the web interface)
- User specific settings
- User verification states
- Reports created by or against users
- Karma scores per guild
- Ownership of a created tag or vote
- Starboard entries

User input text and uploaded media is stored to be used in the following cases:
- Reason and proof for a created report
- Description of vote contents

shinpuru also has a feature called `Antiraid` which watches the influx rate of users to a guild. When shinpuru detects an anomaly, the system triggers and all following users joining are logged. **This data is automatically removed after 48 hours**. These logs include the following data:
- User ID
- Username and Discriminator
- Account created date *(calculated from the ID)*
- Guild join timestamp

### B) Performance

shinpuru **temporarily** caches a lot of data from the Discord API to improve performance and reduce load on the Discord API (primarily to avoid rate limit timeouts). All sensitive data stored in the cache is no longer stored than 30 days.

This includes the following data:
- User IDs
- Usernames
- Discriminators
- Avatar IDs
- Nicknames
- Guild join timestamps
- Message content / embeds

The main purpose of storing this data is to requesting it from the cache, if available, instead of from the Discord API. This massively reduces latencies and avoids frequent calls to the Discord API, which will result in rate limit timeouts.

## How can I remove my personal data from shinpuru?

There are two options in place to directly remove personal data in the web interface.

### A) Guild data removal

You are able to remove all stored data linked to your guild via the web interface. Therefore, open the web interface, select the desired guild, click on the ⚙️ Button to go to the guild settings, select `Data` in the navigation menu and click on `DELETE ALL GUILD DATA PERMANENTLY`.

https://user-images.githubusercontent.com/16734205/149119858-f41fe4e6-8239-45ec-b23d-c8a666e2d128.mp4

This removes **all** data from the database and cache which is linked to the guild.

### B) User data removal

You can remove personal data stored and linked to your Discord account directly via the web interface. Therefore, log in to the web interface, navigate to `User Settings` and press `DELETE USER DATA` in the `Privacy` section.

https://user-images.githubusercontent.com/16734205/149120171-1d698e73-def1-42e9-a653-d7c2caed5bd6.mp4

This action removes the following data linked to your Discord account:
- Starboard entries 
- API tokens
- Refresh tokens
- Tags created by you
- Unban requests created or processed by you
- Your user settings
- Your verification state
- User and member data in the cache

Because of measures to avoid abusement of the security systems implemented into shinpuru (like the report and karma system), not all data is removed which is linked to the user. The following data will not be removed by this action:
- Reports created against you
- Reports created by you are not removed but your ID will be anonymized
- Karma entries when your karma score is below 0
- User settings temporarily stored in the cache

If you want this data to be removed, please contact me (zekro) directly. You can find my contact information [at the top of this document](#contact-details).

---

shinpuru Privacy Notice v1.0.0.  
Last Edit: 2022/01/13.
