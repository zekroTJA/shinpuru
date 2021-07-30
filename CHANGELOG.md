1.20.0

## Changes

### Using Reactions to control `clear` command [#249]

You can now either select single messages to be deleted by adding the ❌ reaction to the messages or you can delete all messages after (and including) the message marked with the 🔻 reaction. The delete action is then performed on entering the [`clear selected`](https://github.com/zekroTJA/shinpuru/wiki/Commands#clear) command.

https://user-images.githubusercontent.com/16734205/126893039-d23dbe44-8bdd-4ab0-b03e-793168fbf620.mov

https://user-images.githubusercontent.com/16734205/126893059-7e54c886-9c95-4d13-ab2d-3b458614d723.mov

### Setting multiple autoroles [#147]

You can now set multiple autoroles in the web interface as well as via the [`autorole`](https://github.com/zekroTJA/shinpuru/wiki/Commands#autorole) command.

![](https://i.imgur.com/JDO30Uf.gif)

https://user-images.githubusercontent.com/16734205/127679480-fec63b3b-11e4-4ba7-a5f8-4e62c657e612.mov

### Invite Badge

As you might know, when the web interface is enabled, you can use `<address>/invite` as redirect link to the invite link of the bot instance (see https://shnp.de/invite, for example). Additional to that, an endpoint was added which generates a badge with the current guild count of the instance.

Examples:

- ![](https://shnp.de/invite/badge.svg) `https://shnp.de/invite/badge.svg`
- ![](https://c.shnp.de/invite/badge.svg?title=invite%20(canary)&color=orange) `https://c.shnp.de/invite/badge.svg?title=invite%20(canary)&color=orange`

## Bugfixes

- shireikan now also uses dgrs state when passed.
- Hydrated guild states are now obtained where required.
- The tagxinput component has now proper item alignment.
- Setting `@everyone` as autorole is no more possible.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.20.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.20.0
```
