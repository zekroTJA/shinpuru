1.15.0

## Minor Changes

- If you are not present on the karma scoreboard in the web interface, your current karma score for that guild is now displayed above the scoreboard. [#240]  
![](https://i.imgur.com/tA4dpC0.png)

- In the karma preferences, you can now enable a Karma penalty. When enabled and when a user decreases the karma of another user by giving them a downvote, the executor of the downvote pays with 1 Karma point from their own karma account. So you give someone -1 Karma and you will also get -1 Karma. This is introduced to reduce karma trolling and uncontrolled downvoting of members.    
![](https://i.imgur.com/Ert3Tdd.png)

## Bug Fixes

- Votes should now be saved properly in the database. [#242]

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.15.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.15.0
```
