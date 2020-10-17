1.3.1

> MINOR PATCH

## Minor Improvements

- The `lock` command now has a process visualization for better clarity.

## Bug Fixes

- A critical flaw in the permissions middleware which would practically bypass the whole permission system is now fixed. [#169]

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.3.1
```