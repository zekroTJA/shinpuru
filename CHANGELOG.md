[VERSION]

<!-- > **Attention**  
> This is a hotfix patch. If you want to see the changelog for release 1.30.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.30.0). -->

# New Login Method

On the main page, when clicking on `Login`, you can now choose between logging in to the web interface via OAauth2 using your Discord account or log in via sending the code, which is displayed in the web interface, to shinpuru via DM.

![](https://user-images.githubusercontent.com/16734205/154697491-b0aa34d3-ff79-40ee-9b49-ec77cfc23cee.gif)

# Deployment

The generated frontend files are now directly embedded into the binary of shinpuru. Therefore, downloading and providing the frontend files in the same directory of the binary is no more necessary.

# Bug Fixes

- The birtdhay command will now only send command responses to the user who invoked the command. [#354]
- The embed builder now only shows available text channels where the logged in user has read and write permissions. [#353]
- The embed preview now shows a placeholder title and description when empty. [#353]

# Acknowledgements

Big thanks to the following people who contributed to this release.

- @zordem

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
