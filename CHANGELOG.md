1.10.2

> This is a hotfix patch. [**Here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.10.0) you can find the full changelog of patch 1.10.0.

## Bug Fixes

- Media contained in messages stared from NSFW channels appearing in the starboard *(when it is not marked as NSWF itself)* are blurred out. [#217]   
  <img src="https://i.imgur.com/lJka5k9.png" height="350"/>
  <img src="https://i.imgur.com/p2VtKcb.png" height="350"/>

- Fix error which may occur on some guilds after setting the same starboard config again.

## Backstage

- Add package [`pkg/thumbnail`](https://github.com/zekroTJA/shinpuru/tree/master/pkg/thumbnail).

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.10.2
```