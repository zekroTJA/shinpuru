1.16.0

## Major Changes

### SQLite Deprecation

The SQLite Database driver, which is neither continuously maintained, nor tested, nor recommendet to be used in a production scenario, is now marked as **deprecated** and will be removed in the upcoming version!

For development environments, please use the provided [`docker-compose.dev.yml`](docker-compose.dev.yml) to easily set up a development environment for shinpuru.

### Web Frontend

This release contains further design changes of the web frontend.

First of all, I've finally changed the background color scheme. Originally, shinpurus web frontend should have followed the original Discord desogn language. Thats why still most colors are derived from them. Though, my view has changed a bit on the blue-gray background colors. To be honest, I really started hating them so much so that I've changed them to a more 'classical' dark-gray.

Also, as you might have noticed, the header is now split up into two "floating" parts. In my opinion, the central part of the header was just a useless waste of space, so now, it's a bit more condensed and fits better into the general design.

![](https://i.imgur.com/78XgnyZ.png)

Also the guild settings page got a massive redesign. First of all, the navigation bar is now also a floating navigation menu. The page content is now centered to be consistent with other content orientations.

![](https://i.imgur.com/HpLJ1mH.gif)

Also, I've done a lot to make the guild settings more responsive to mobile devices. It's not perfect though, but better than before. ^^

![](https://i.imgur.com/hOckQ0J.gif)

<!-- ## Minor Changes -->

## Bug Fixes

- Code execution can now only be triggered by the author of the code message. [#244]

- The route `/invite` now redirects to the bot's invite link again. [#245]

## Backstage

- All report utilities are now summarized in a report module which is registered in the DI container. This makes further modifications and implementations with the report system more easy.

- shinpuru now uses a new version of [`timedmap`](https://github.com/zekroTJA/timedmap) which should bring some performance improvements.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.16.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.16.0
```
