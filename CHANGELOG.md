1.18.0

## Major Changes

### Ban/Mute Timeout [#246]

You can now specify a timeout duration or date to mute and ban reports. After this time, the victim is automatically unbanned or unmuted, but the report will not be removed.

![](https://i.imgur.com/bs8meYt.png)

Here you can read on how to use timeouts in the mute and ban command.
- [`mute`](https://github.com/zekroTJA/shinpuru/wiki/Commands#mute)
- [`ban`](https://github.com/zekroTJA/shinpuru/wiki/Commands#ban)

The expiration of reports is defaultly checked every 5 minutes, but you can specify an other schedule in the config.
```yml
schedules:
  reportsexpiration: "@every 5m"
```

### SQLite3 Deprecation

The SQLite3 driver was marked as deprecated in patch [1.16.0](https://github.com/zekroTJA/shinpuru/releases/tag/1.16.0) and has now fully been removed. If you need information why this step was taken and how to switch to MariaDB for development, please read [this document](https://github.com/zekroTJA/shinpuru/wiki/SQLIte-Deprecation). 

## Minor Changes

- Add loading indicator in the web interface.  
  ![](https://i.imgur.com/HK8fXUJ.gif)

- [dgrs](https://github.com/zekroTJA/dgrs) was now updated to v0.3.0 which allows dehydration of removed objects and caching of User→GuildIDs relationships.

- Hence, the load time of the guild list in the web interface should be significantly faster. [#257]

- Messages are now only cached for 14 Days in state.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.18.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.18.0
```
