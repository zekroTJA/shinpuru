[VERSION]

<!-- > **Attention**  
> This is a hotfix patch. If you want to see the changelog for release 1.30.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.30.0). -->

# Changes

Most of the slash commands are now empehmeral. That means, that only you will see the response and call of the command in the chat.

<img width="414" alt="image" src="https://user-images.githubusercontent.com/16734205/210225426-0c6ed18e-ba79-46b9-941a-78302687b09b.png">

Therefore, also the `/login` command has been updated so that the response is directly sent to the chat. That is possible because only the sender can see the response. That makes the command available also to people who have disabled DMs from guild members.

<img width="635" alt="image" src="https://user-images.githubusercontent.com/16734205/210226175-c36d7aae-726f-46f4-8680-983700eae5dd.png">


# Bug Fixes

- Fixed a bug where the guild API token is reset when saving the guild API settings without a new token.
- Fixed misrepresentation of absolute sub command permissions in the `permissions/allowed` endpoint. [#398]
- Fixed some typos (and added some more).

# Beta Web Interface

- Added verification route.
- Guild Settings: Added guild data removal route.
- Guild Settings: Added guild API route.
- Added some more (english) explanation texts.

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
