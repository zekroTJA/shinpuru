1.9.0

## Minor Updates

- **Web Interface Responsiveness** [#172, #192]  
  The web interface is now more web responsive and for better usability on mobile devices.

- **Login Command** [#206]  
  A new command [`sp!login`](https://github.com/zekroTJA/shinpuru/wiki/Commands#login) was added which
  sends a message to the executor via DM *(if enabeld by the user)* containing a Link with a one time
  authentication token which is valid for 1 minute which logs you in to the web interface without being
  logged in in the browser, which is especially useful on mobile devices.  
  ![](https://i.imgur.com/BrpZcOY.png)

## Security Fixes

- The download of guild backups now requires the `sp.guild.admin.backup` permission, because the backup file
  contains confidential information like hidden channel details or guild settings. [#208]

## Backstage

- Add package [timerstack](https://pkg.go.dev/github.com/zekroTJA/shinpuru/pkg/timerstack).

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.9.0
```