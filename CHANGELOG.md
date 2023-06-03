[VERSION]

# New Features and Changes

## Mod Notifications [#423]

There is a new command [`/modnot`](https://github.com/zekroTJA/shinpuru/wiki/Commands#modnot) which can be used to define a Mod Notification channel.

![image](https://github.com/zekroTJA/shinpuru/assets/16734205/00cce5ff-8a83-401d-b61a-89d815bf5053)

Of course, you can also define the channel in the web interface on the General Guild Settings page.

![image](https://github.com/zekroTJA/shinpuru/assets/16734205/75ca7ae6-4027-40ec-8013-4361567658f3)

This channel is used to - as the name says - send notification important for the members moderating your guild. Currently, this only contains notifications when someone creates an unban request or an unban request has been processed. But other notifications will find their place here in the future like raid alerts or spam alerts.

![image](https://github.com/zekroTJA/shinpuru/assets/16734205/4c853e4d-8e6a-445c-a822-74abc9626798)

## German Translation

Thanks to the work of @voxain, the german translation of the web interface has improved in consistency, conciseness and understandability.

# Bug Fixes

- Ban reports can now be created when the executing user has only the `sp.guild.mod.ban` permission. [#426]
- Fixed the redirection to the dashboard after logging in to the web interface using the DM code.
- Fixed a bug where the unban request page (`/unbanme`) would not load and error when an unban request is still open on a guild where shinpuru has been removed from. [#427]
- Bans executed with shinpuru will no more create post ban prompts. [#424]
- Removed the gradient in the feature marquee on the start page because of incompatibilities across browsers. [#421]

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
