1.2.1

> MINOR PATCH

## Minor Improvements

- Discord assets in web interface are now requested with a proper resolution. [#163]
- Report revoke link in web interface is now only shown if the user has the permission to do so.

## Bug Fixes

- Fix permission update handling when adding the same permission rule as existent. [#161]
- Fix dropdown style in web interface ([old](https://i.imgur.com/m5uQZdq.png) vs [new](https://i.imgur.com/PWet0kD.png)).
- Fix color reaction spam on message edit. [#162]

## Backstage

- Moved the [permissions](https://github.com/zekroTJA/shinpuru/tree/master/pkg/permissions) and the [twitchnotify](https://github.com/zekroTJA/shinpuru/tree/master/pkg/twitchnotify) packages to the `pkg` public package domain.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.2.1
```