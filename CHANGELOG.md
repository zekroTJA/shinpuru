1.17.1

## Minor Changes

- Refine request logging *(only in debug log channel)*.

## Bugfixes

- Fix critical bug where occasionally, none or only a fraction of guilds were detected for member fetching on ready. [#252]
- Fix command logging which was not existent even when enabled by config. [#253]

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.17.1
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.17.1
```
