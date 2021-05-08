1.12.0

## Major Changes

### ranna implementaiton [#232]

shinpuru now supports code executing using [**ranna**](https://github.com/ranna-go). You can add the folowing to your shinpuru config to connect to a ranna service instance.
```yaml
codeexec:
  type: "ranna"
  ranna:
    # The ranna instance endpoint.
    endpoint: "https://private.ranna.zekro.de"
    # The ranna instance API version.
    apiversion: "v1"
    # Here you can pass a token if the ranna instance requires one.
    # This is the exact value set as 'Authentication' header. If
    # you do not specify a token type ('bearer ...', for example),
    # the token is automatically prefixed with 'basic '.
    token: ""
```

![](https://i.imgur.com/r2l5gaa.png)

If you do not specify any code executor, Jdoodle is still used. Also, if using Jdoodle, required credentials are still required to be set on a per-guild-basis using the [`exec`](https://github.com/zekroTJA/shinpuru/wiki/Commands#exec) command. When using ranna as code execution engine, this command is disabled and code execution is enabled always enabled.

## Minor Changes

- [logrus](https://github.com/sirupsen/logrus) is now used as logger to provide a more rich log output and make logging easier in general.

## Bug Fixes

- shinpuru will no more crash when not providing rate limit configuration. [#230]
- Invite links from the same guild as where the message was sent from are no more blocked by the guild invite block system.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:
```
$ docker pull zekro/shinpuru:1.12.0
```

From GHCR:
```
$ docker pull ghcr.io/zekrotja/shinpuru:1.12.0
```