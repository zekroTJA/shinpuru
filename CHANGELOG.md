[VERSION]

<!-- > **Attention**  
> This is a hotfix patch. If you want to see the changelog for release 1.30.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.30.0). -->

# Message Components [#320]

With release [v0.17.0](https://github.com/zekroTJA/ken/releases/tag/v0.17.0) and [v0.17.1](https://github.com/zekroTJA/ken/releases/tag/v0.17.1), message components have now been added to shinpuru.

The [`acceptmsg` package](pkg/acceptmsg) now uses message components to accept and decline the message.

![](https://user-images.githubusercontent.com/16734205/187892843-22cc6e0e-a838-40cc-a24d-2d957fd5d4d7.png)

Step by step, more commands and services will be adapted to message components and modals. One of the first commands using message components and modals is the `/backup` slash command.

<!-- TODO: Add GIF. -->

# "Postban" System [#383]

A new "postban" system has been implemented which detects when a user has been banned manually (directly using the Discord utilities) and sends a notice message into the mod log channel afterwards. Then, you are able to import the ban into the shinpuru report system. You are also able to edit the reason specified in the ban and also add an attachment via an URL.

<!-- TODO: Add GIF. -->

# Minor Changes

- The mime fix utility has now been removed because the issue [has been fixed in the go language](https://go-review.googlesource.com/c/go/+/406894/).
- Add code execution guild settings to the new beta web interface.

# Code Base

The `Makefile` in this repository - which is primarily there to simplify common prep, build, test and deploy tasks - will be swaped out in favor of [`Taskfile`](https://taskfile.dev) which has a lot of advantages over make - especially in this environment.

If you don't have taskfile installed, please head to [taskfile.dev/installation](https://taskfile.dev/installation/) to find out how to install taskfile on your system.

Currently, the `Makefile` will be kept in the repository for compatibility reasons but will be removed in upcoming updates.

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
