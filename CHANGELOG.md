[VERSION]

<!-- > **Attention**  
> This is a hotfix patch. If you want to see the changelog for release 1.30.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.30.0). -->

# Bug Fixes

- Fixed the endlessly loading account verification captcha. [#392]
- Fixed a crash when shinpuru has no access on the guild audit log after detecting a ban (postban system).
- Fixed a crash when shinpuru has no access on the guild mod log channel after detecting a ban (postban system).
- Fixed a bug that not postban message is sent when no reason is specified to the ban.

# Beta Web Interface

- Guild Settings: Guild Log implemented

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
