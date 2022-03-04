[VERSION]

<!-- > **Attention**  
> This is a hotfix patch. If you want to see the changelog for release 1.30.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.30.0). -->

# shinpuru is now finally verified! 🎉

Even though this does not relate to the actual release, I just want to say that the stable shinpuru instance ([`shinpuru#4878` / shnp.de](https://shnp.de)) is now verified for privileged intents. That means that the 100 guilds limit is now lifted and you can now invite shinpuru on howmany guidld you want! That also means that services which rely on member events and message content will still be operational in the future.

# Auto Voice Channel [#324, #357]

Thanks to the contribution of @TomRomeo, shinpuru now supports auto voice channels. You can simply specify a voice channel as auto voice channel using the `/autovc add` command. After joining the auto voice channel, a new channel is created and you are automatically moved into this channel. After that, all your friends can join this channel with you. After you leave the channel, it is automatically removed again.

https://user-images.githubusercontent.com/16734205/156553740-3628a6cf-5386-4b54-86ff-87c10cbf2cf9.mp4

# Update Information

shinpuru now also checks for updates by comparing the current build version with the latest tag of the GitHub repository. If a newer version is available, you will see this notification on startup in the log.

![](https://user-images.githubusercontent.com/16734205/156564557-6c406006-8ae0-4113-9ef4-b470e97e7cd8.png)

Also, when the bot owner logs in to the web interface, they will also receive a notification showing the available upgrade.

![](https://user-images.githubusercontent.com/16734205/156564782-081a8355-4033-4a83-8ec3-a67c1971e255.png)

# Guilds Limit

You can now specify a global guilds limit in the config on which your shinpuru instance can be member of. This might be useful to either hard-cap the amount of guilds to preserve resources or to avoid hitting the magical 100 guilds cap of unverified bot accounts.

When applied and the guild limit is exceeded after shinpuru joins a new guild, it will leave the guild and send an informative message to the owner of the guild.

![](https://user-images.githubusercontent.com/16734205/156731242-484bba6e-66dc-4105-9979-3e84855c21dc.png)

*Via config file*
```yml
discord:
  # Specify a maximum of guilds the bot can
  # be member of. When set to 0, there is
  # no limit applied.
  guildslimit: 90
```

*Or via environment variable*
```
SP_DISCORD_GUILDSLIMIT=90
```

# Bug Fixes

- Fixed some typos in autorole and twitch notify command. [#358]
- Fixed a bug that simultanious unauthorized requests produce multiple concurrent access token requests. [#361]
- Fixed error message in login screen when push codes are timing out.

# Acknowledgements

Big thanks to the following people who contributed to this release.

- @TomRomeo

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
