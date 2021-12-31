1.26.0

# Antiraid Joinlog Account Age [#283]

The antiraid joinlog table now also shows the account age and creation date which is very useful to validate if an account might be boted.

![](https://user-images.githubusercontent.com/16734205/147818785-6cd358d7-fa48-475f-9b4b-15defd70c730.png)

# `userinfo` user command [#308]

The `/user` slash command can now be directly executed on members using the `userinfo` user command. This is available when right-clicking on a user, going to `Apps` and clicking `userinfo`. Then, the user info card will be dropped into your currently selected text channel.

![](https://user-images.githubusercontent.com/16734205/147783482-0b3dc68c-2f07-4bed-b26c-421c0a6ddb17.png)

# `quotemessage` message command

The `/quote` slash command is now also usable directly via the `quotemessage` app command when right-clicking the message to be quoted. The quote message will appear in the currently selected channel as well.

![](https://user-images.githubusercontent.com/16734205/147783769-d7b80e68-ba5a-4649-aff6-0571bb99b132.png)

# Mute rework [#315]

Because Discord recently added the `timeout` feature, the usage of specific mute roles which disallow sending messages in all channels is no more necessary. Instead, the mute/unmute command and web interface hook utilizes the timeout implementation of Discord. So, you do not need to setup and maintain a muterole anymore and you can directly use the timeout integration of Discord with the advantages of shinpurus modlog! 🤯


# Bug Fixes

- Fixed a bug where the pro tip card in the web interface prevents clicking underlying elements after dismissal. [#316]
- The rework of the mute command also obsoletes the mute role issue. [#312]

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.26.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.26.0
```
