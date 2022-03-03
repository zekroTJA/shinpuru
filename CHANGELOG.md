[VERSION]

<!-- > **Attention**  
> This is a hotfix patch. If you want to see the changelog for release 1.30.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.30.0). -->

# Auto Voice Channel [#324, #357]

Thanks to the contribution of @TomRomeo, shinpuru now supports auto voice channels. You can simply specify a voice channel as auto voice channel using the `/autovc add` command. After joining the auto voice channel, a new channel is created and you are automatically moved into this channel. After that, all your friends can join this channel with you. After you leave the channel, it is automatically removed again.

https://user-images.githubusercontent.com/16734205/156553740-3628a6cf-5386-4b54-86ff-87c10cbf2cf9.mp4

# Update Information

shinpuru now also checks for updates by comparing the current build version with the latest tag of the GitHub repository. If a newer version is available, you will see this notification on startup in the log.

![](https://user-images.githubusercontent.com/16734205/156564557-6c406006-8ae0-4113-9ef4-b470e97e7cd8.png)

Also, when the bot owner logs in to the web interface, they will also receive a notification showing the available upgrade.

![](https://user-images.githubusercontent.com/16734205/156564782-081a8355-4033-4a83-8ec3-a67c1971e255.png)

# Bug Fixes

- Fixed some typos in autorole and twitch notify command. [#358]
- Fixed a bug that simultanious unauthorized requests produce multiple concurrent access token requests. [#361]

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
