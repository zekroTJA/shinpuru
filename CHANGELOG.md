1.23.1

![](https://i.imgur.com/epbx2Cw.jpeg)

# Changes

## Config Management

With the [`/maintenance`](https://github.com/zekroTJA/shinpuru/wiki/Commands#maintenance) command has been extended by 2 sub-commands.

- `reload-config` allows to reload the config during runtime from the given config sources.

- `set-config-value` allows to set config values during runtime by field name and value JSON representation.  
  *Please use this function with caution because it can hardly impair the functionaility of shinpuru.*

Also, keep in mind that some config changes may only take effect after a restart. So both commands may have no effect.

## Bug Fixes

- Set requirement on `/ghostping setup message:<string>` argument so that shinpuru will not crash anymore if the parameter was not specified.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.23.1
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.23.1
```
