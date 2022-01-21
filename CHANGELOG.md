[VERSION]

<!-- > **Attention**  
> This is a hotfix patch for issue #332. If you want to see the changelog for release 1.26.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.27.0). -->

# Code Inline Execution [#334]

https://user-images.githubusercontent.com/16734205/150545223-cba2e27b-9b98-4e80-83bd-b3e0d8063a90.mp4

Since [ranna](https://app.ranna.zekro.de) now supports inline code exectuion, this feature is now also implemented in shinpuru. When a language has been used which supports inline code execution, an extra emote is added to the message which invokes the inline code execution. 

Also, a help action has been added to explain all of this when clicked.

# Bug Fixes

- Fixed error handling and error bubbling when OAuth2 login fails. [#328]
- Fixed a bug that can crash shinpuru when a message's author is `nil` *(which even happens whyever?)* when using `/quote`. [#342]

# Codebase Changes

The codebase of shinpuru has now been upgarded to be compiled with Go version 1.18. Because this version is not release as stable, version `1.18beta1` is used until the official release of Go 1.18.

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
