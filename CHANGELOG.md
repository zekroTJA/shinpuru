1.19.0

![](https://i.imgur.com/Dy0Xbr7.png)

## Minor Changes

- Updated the visual representation of the guild dropdown.  
  ![](https://i.imgur.com/zg6Sf2l.png)

- Use skeleton tiles as loading indicators instead of spinners.  
  ![](https://user-images.githubusercontent.com/16734205/126753381-224a6a62-33ec-4dd0-814e-ab71c0699fa3.gif)

## Bug Fixes

- Role position diffs are now properly checked on each ban/kick/mute/unmute report execution.

- Fixed guild settings button alignment.  
  ![](https://i.imgur.com/6rL1lKD.png)

- Fixed guild tiles alignment in home view.  
  ![](https://i.imgur.com/1nDlAr7.png)

- Fixed displaying animated guild icons.  
  ![](https://i.imgur.com/IbIhS7p.gif)

- The menu dropdown now properly closes when clicking on the user tile.  
  ![](https://i.imgur.com/0uYa7Qd.gif)

- Anonymous report creation will no more result in a `404 Member not found` error.

- Report timeouts can now also be defined when creating anonymous reports.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.19.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.19.0
```
