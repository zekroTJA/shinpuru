1.21.0

## Changes

### Added Global Search

By pressing `CTRL + F` anywhere in the web interface, you can now bring up a global fuzzy search which will search through all accessable guilds and their members. You can even navigate in the search only by using the keyboard! 😉

![image](https://user-images.githubusercontent.com/16734205/128348276-8a81ebf3-21eb-4da6-bac0-88ec3ff4bf78.png)


### Public Guild API

You can now enable a public API endpoint which exposes general information about your Guild via shinpuru's REST API.

![image](https://user-images.githubusercontent.com/16734205/128348885-e1e2dffc-6629-40db-b184-4fac8ac94e03.png)

The output of the endpoint will then look as following.

![image](https://user-images.githubusercontent.com/16734205/128349603-ebaa5bbf-6917-44f8-b296-05b05bf5be9e.png)


### Updated the Style of the Notifications

The design of the notifications now fits in better with the general design language of the web app. Also the space around the notification box was adjusted to fit under the new header.

<img src="https://user-images.githubusercontent.com/16734205/128346539-9dd58670-3b80-426a-9900-bd537e6be85c.png" height="400" />

### Updated the Style of some Icons

Also some of the used Icons did not fit in the new design anymore and have been adjusted. As an example, below you can see the old vs. the new drop down icon.

<img src="https://user-images.githubusercontent.com/16734205/128347326-1138f5c1-6bac-4887-9b2a-915370343dca.png" width="300" />


## Bugfixes

- **Vulnerability Fix**: OTA tokens now support scoping and scope validation so that they can only be used for the exact purpose they were originally issued for. [#264]
- Added message reaction tracking to [dgrs](https://github.com/zekrotja/dgrs) to fix starboard functionality.
- Starboard guild settings are now also properly saved to the database instead only to cache.
- Fix the heading content of the invite blocking toggle notification.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.21.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.21.0
```
