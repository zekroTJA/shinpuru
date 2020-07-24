1.1.1

> MINOR PATCH  

> This is only a hotfix release. To read the full changelog, see [release `1.1.0`](https://github.com/zekroTJA/shinpuru/releases/tag/1.1.0).

## Fixes

- Fixed [button bindings](https://i.zekro.de/f9pyAAhYbF.gif) on security cards in web interface.
- Updated `Session#GetGuild` functions with `discordutil.GetGuild` function in web server handlers, which provides issues with unhydrated data for guild channels or members.
- Fixed invite page when a user does not share any guilds with shinpuru.
- Fixed no guild invite API endpoint.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.1.1
```