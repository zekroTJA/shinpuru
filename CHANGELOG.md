1.27.0

<!-- > **Attention**  
> This is a hotfix patch for issue #317. If you want to see the changelog for release 1.26.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.26.0). -->

# New Featrues

## User Verification (Beta) [#323]

You can now enable user verification in your guild settings. When enables, new users to the guild will be required to verify their account by solving a captcha. When this is done once, theit account is marked as `verified` across all guilds where shinpuru is member of. As long as a new user is not verified, they are timed out until they verify their account. Otherwise, they will be automatically kicked from the guild after 48 hours.

![](https://user-images.githubusercontent.com/16734205/148562832-11590501-e2be-481d-8fb5-87bed122b987.png)

You can also combine that with the antiraid system. When enabled, the user verification system will be enabled automatically when the antiraid alert triggers.

![](https://user-images.githubusercontent.com/16734205/148563904-532c7960-7a1f-47d9-957a-e80a5ce616d3.png)

As visible in the screen shots, this feature is still in beta phase and bugs as well as unexpected behaviour might occur. Also, the implementation is subject to change and extension in the future. If you encounter any issues or if you have ideas to improve the system, feel free to create an issue.

# Bug Fixes

- Fix `mute` legacy command parsing. [#325]

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.27.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.27.0
```
