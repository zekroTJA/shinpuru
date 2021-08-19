1.22.0

## Changes

### New Config Handling [#274]

shinpuru now uses [traefik/paerser](https://github.com/traefik/paerser) to parse the configuration from multiple sources at the same time. These are the available configuration sources, sorted by priority.

1. **Command Flags**  
   Configuration values passed via command flags. Example:
   ```
   ./shinpuru \
       --discord.token "..." \
       --discord.generalprefix "!"
   ```

2. **Environment Variables**  
   Configuration values passed via environment variables prefixed with `SP_`. This is especially useful when hosted via Docker. Example:
   ```
   SP_DISCORD_TOKEN="..."
   SP_DISCORD_GENERALPREFIX="!"
   ```

3. **Config File**  
   You can still pass configuration as usual via configuration file. Defaultly, the config is read from `./config.yml`, but you can pass another file location via the `-c` flag. You can also use other configuration formats like JSON or TOML. Example:
   ```
   ./shinpuru -c config/config.yml
   ```
   > config/config.yml
   ```yml
   discord:
     token: "..."
     generalprefix: "!"
   ```

You can combine all configuration sources listed above. Higher priorized configuration sources will overwrite values from less priorized sources.

### Embed Builder

You can now send and edit embed messages in guild channels using the `POST /api/v1/channels/{id}` and `POST /api/v1/channels/{id}/{messageid}` endpoints.

> This endpoint requires the `sp.chat.say` permission.

Here you can find the documentation:
https://github.com/zekroTJA/shinpuru/blob/master/docs/restapi/v1/restapi.md#channelsid

There is also an embed builder using these endpoints. But because this is still kind of beta, you can currently only access it directly via the following route in the web interface.

```
/guilds/{guildid}/utils/embeds
```

![](https://i.imgur.com/T9qEiyU.png)

## Bugfixes

- Fix report time representations. [#276]
- Fix report unmute reason propagation. [#277]
- Fix proper timezone handling on report expiration definition.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.22.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.22.0
```
