[VERSION]

> **Attention**  
> This is a hotfix patch for issue #332. If you want to see the changelog for release 1.26.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.27.0).

# Bug Fixes

- Fix timestamp parsing when creating starboard embed. [#332]
- Fix timestamp formatting in embeds.
- Replace `FATAL` log outputs with `ERROR` where shutdowns are not necessary.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:[VERSION]
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:[VERSION]
```
