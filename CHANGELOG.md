1.6.0

> MAJOR PATCH

## Major Implementations

### "Anonymous" Reports [#189]

Now, it is possible to report members by ID which are not actually on the guild. These reports are handled as same as "normal" reports without actual user validation.  
This is especially useful to report members which violate the rules of your guild but leave the guild before you are able to report them.

Also, you are actually able to ban members which are not on the guild so you are able to ban members which left the guild or never were a member of the guild.

![](https://i.imgur.com/MNqwsCR.png)
![](https://i.imgur.com/3trvWHm.png)

### `guildinfo` Command [#191]

shinpuru now has a new command: [`guild`](https://github.com/zekroTJA/shinpuru/wiki/Commands#guild).  
It simply outputs some general information about the guild where the command was executed on.

![](https://i.imgur.com/8RmrDT7.png)

## Minor Updates

- Some features like the vote command or color embed system now take advantage of the new [reply feature](https://support.discord.com/hc/en-us/articles/360057382374-Replies-FAQ) of Discord.  
![](https://i.imgur.com/wOjcqyv.png)

- The header of the web interface now uses the new logo of shinpuru. Also some spacings issues are fixed now.  
![](https://i.imgur.com/vEU7PJv.png)

## Bug Fixes

- Fix guild tile titles in web interface [#190]
- Fix API token route layout

## Backstage

- Add package [embedbuilder](https://github.com/zekroTJA/shinpuru/tree/master/pkg/embedbuilder) to `/pkg`

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.6.0
```