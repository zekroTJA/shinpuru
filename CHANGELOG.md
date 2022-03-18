[VERSION]

<!-- > **Attention**  
> This is a hotfix patch. If you want to see the changelog for release 1.30.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.30.0). -->

# Interactive Setup

Recently, a lot of questions about the setup of shinpuru received me. shinpuru has a lot of configuration values and ways to configure it, so I decided to create a small tool where you can enter all your settings and credentials and which creates a pre-configured `docker-compose.yml` which you can simply use to set up shinpuru on your server. 

https://user-images.githubusercontent.com/16734205/158845896-50eb4869-fa2a-42f0-887c-0b946ddcecd0.mp4

Simply download the tool from the releases below, execute it in your terminal, enter your settings and credentials and then you will have a ready to go Docker Compose set up to be deployed to your server.

# Report QoL Changes

Previosuly, you were able to report, kick and ban users using the `type` parameter of the report command. On the one side, this is highly unintuitive because people more likely expect a ban and kick command for each action separately. Also, this is very inconsistent with the `mute` command, which is actually also just a report type with extra actions like banning an kicking as well at it has its own command.

So, with report you can now only create, list and revoke reports and banning as well as kicking members is done with both separate commands `ban` and `kick`.

This change also results in some changes to the permissions.
- `sp.guild.mod.report.kick` is now `sp.guild.mod.kick`
- `sp.guild.mod.report.ban` is now `sp.guild.mod.ban`

# Starboard Changes [#369 <small>*nice*</small>]

When posing message with links to images or videos, these are extracted from the message and put into the embed itself when voted into the starboard.

Also, videos are now properly displayed in the web interface. And yes, autoplay is disabled and the volume is muted by default. 😉

![](https://user-images.githubusercontent.com/16734205/158994685-eae81863-77c0-4e82-b05e-a01b273a1577.gif)

# Minor Changes

- The guild log order in the web inetrface is now descending starting with the latest entries. This is way more intuitive and easier to look for recent issues. [#366]

- shinpuru now reads a `.env` file in the directory of the executable and applies it to the environment vaiables.

# Bug Fixes

- Fixed a bug which results in a `nil` author on a message retrieved from the cache which lead to a panic when voting such a message into the starboard. [#366]

- Not available or deleted starboard channels are now properly cleared from the settings when voting messages into it.

# Code Base

So that go 1.18 has now officially being released, you can now compile shinpuru with the latest version of the Go toolchain instead of using a pre-release build.

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
