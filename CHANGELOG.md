[VERSION]

<!-- > **Attention**  
> This is a hotfix patch. If you want to see the changelog for release 1.30.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.30.0). -->

# New Web Interface Beta

Currently, I am in the process of re-implementing the whole web interface of shinpuru. If you want to know why and how far it already has been processed yet, please take a look at issue #370. With this release, you can check out the new interface by going to the `/beta` route in the web interface. But please, keep in mind that this is a very early state of development. A lot of features are still missing or only partly implemented and the user experience might be impaired.

![](https://user-images.githubusercontent.com/16734205/178149927-2b100aa8-f33e-403c-8b38-2d2d4afdca18.gif)

# Duration UX Improvements

All duration parsing of commands has now been exchanged from the [default Go Implementation](https://pkg.go.dev/time#ParseDuration) to a [custom implementation](https://github.com/zekroTJA/shinpuru/blob/dev/pkg/timeutil/timeutil.go#L49-L117) which allows more time units like days (`d`) or weeks (`w`). Also, it makes the parsing more flexible. For example, you can use spaces between units (example `1d 3h`) or you can even subtract values (example `1d-2h` == `22h`). This affects all commands which take a duration parameter like `/ban` or `/mute`.

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
