[VERSION]

> **Note**  
> This is a hotfix patch. If you want to see the changelog for release 1.40.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.40.0).

# Bug Fixes and Minor Changes

- The tag command has been improved to make use of parameter autocompletion. [#439]
- A bug has been fixed which caused shinpuru to crash on fetching the Twitch API for Twitch notifications. [#445]
- Added better state caching to the voice log listener to try to fix issue #440. This might need more investigation though.
- Fixed permission check on routes which do no contain a `guild` URL path parameter. This fixes an issue where the presence page was not able to load in the web frontend because of the failing request.
- Some German translations have been improved thanks to the contributions made by @luxtracon.

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
