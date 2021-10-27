1.24.0

# Changes

## Antiraid Joinlog Bulk Actions [#282]

You can now select users from the antiraid joinlog to bulk kick or ban them directly from the web interface.

![](https://user-images.githubusercontent.com/16734205/139068595-4eac2654-a51b-496b-a2da-7ab218cef240.png)

## Help Command Improvements

You can now get more detailed help information about a specific command using the `/help` slash command.

![](https://user-images.githubusercontent.com/16734205/139070186-01926e33-043a-4783-8b6e-0e887b31cadd.png)

## Vote "anonymization" [#281]

Before this patch, vote ticks were saved in the database with clear user IDs. That entails the risk that user votes can be backtraced from database dumps. These user IDs are now hashed so that it is more difficult for potential attackers to backtrace user votings.

If you want to read more about this, please read [this wiki article](https://github.com/zekroTJA/shinpuru/wiki/Why-are-Votes-%22pseudo-anonymous%22%3F).

# Bug Fixes

- The [`/commands`](https://shnp.de/commands) web interface route can now also be accessed when not being logged in. [#301]
- Fixed User ID resolution of reports and unban requests in the web interface. [#304]
- Fix `/login` command domain name.
- Fix database initialization and add database connection check.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.24.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.24.0
```
