[VERSION]

> **Note**  
> This is a hotfix patch. If you want to see the changelog for release 1.30.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.39.0).

# Bug Fixes

- A bug has been fixed which resulted in a faulty permission check on revoking reports in the web interface. [#418]
- When the report message can not be sent via DM to the target user, the error will no more be reported. Also, when the report fails to be sent in the mod log channel, the reported error is now more concise. [#419]

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
