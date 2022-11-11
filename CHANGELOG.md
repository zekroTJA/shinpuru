[VERSION]

<!-- > **Attention**  
> This is a hotfix patch. If you want to see the changelog for release 1.30.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.30.0). -->

# Role Selection [#363]

Added a new command [`/roleselect`](https://github.com/zekroTJA/shinpuru/wiki/Commands#roleselect) which you can use to create a role selection message. Alternatively, you can also attach role selection to a message which has been sent by shinpuru *(for example using [`/say`](https://github.com/zekroTJA/shinpuru/wiki/Commands#say))*.

https://user-images.githubusercontent.com/16734205/201224685-1393c46a-891e-4963-beea-b93ddbbba142.mp4

# Discord OAuth Handling

The login to shinpuru using the Discord OAuth2 flow is now taking advantage of passing a `state` parameter to the authentication redirect. This state is a JWT signed with a token randomly generated on startup. This secures the login-process against any type of cross-site request forgery attacks.

Also, this allows to pass additional information through the login process like redirection targets. Therefore, when you log in to the beta web interface, you are now also redirected back to the beta interface after the login instead of being redirected to the main interface. 

# Bug Fixes

- When you go to `shnp.de/beta`, you will now be redirected to the login page if not logged in. [#388]
- Fixed a typo in the anti-raid notification message. [#389]
- Fixed command manual generation. [#391]

# Beta Web Interface

- Guild Settings: Verification Route implemented
- Guild Settings: Code Execution Route implemented
- Guild Settings: Karma Route implemented

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
