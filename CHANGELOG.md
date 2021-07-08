1.17.0

## Major Changes

### New State Management

Even though this change is not primarily visible to users, I want to highlight it here because it is a huge step for shinpuru's development.

A Discord bot requires a so called `state` which caches retrieved objects from the Discord API to improve the user experience and to avoid running into rate limitations on the Discord API. Before, shinpuru relied on the internal state manager of discordgo. This was swapped out with a custom state mamnager called [`dgrs`](https://github.com/zekroTJA/dgrs) which stores state objects in Redis. [Here](https://github.com/zekroTJA/dgrs#dgrs-----) you can read more about the advantages of this implementation.

Also, this implementation allows to run shinpuru without the privileged presence intent. This strongly reduces incomming event traffic. Also, this might make it more easy to get shinpuru verified some time. :)

Because this implementation needs a lot more testing and might be still unstable, there is a new command for bot admins named [`maintenance`](https://github.com/zekroTJA/shinpuru/wiki/Commands#maintenance) which groups a collection of sub commands to flush the state cache, reconnect the Discord session or kill the bot instance to perform a hard-restart.

## Minor Changes

- Added metrics for REST API and Redis connection.

## Bugfixes

- Fixed issue that clear command removed one message less than expected. [#248]
- Fixed login oage to fit into new style.
- Fixed metrics endpoint server.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.17.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.17.0
```
