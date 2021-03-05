1.9.1

## New Features

### One Time Authentication (OTA) [#206]

A new command [`sp!login`](https://github.com/zekroTJA/shinpuru/wiki/Commands#login) was added which
sends a message to the executor via DM *(if enabeld by the user)* containing a link with a one time
authentication token which is valid for 60 seconds to log you in to the web interface without being
logged in in the browser, which is especially useful on mobile devices.  
![](https://i.imgur.com/BrpZcOY.png)

[Here](https://github.com/zekroTJA/shinpuru/wiki/One-Time-Authentication-(OTA)) you can find more details 
on how the OTA system is implemented.

OTA is disabled by default. You need to enable OTA via the new 
User Settings Hub in the shinpuru web interface.  
![](https://i.imgur.com/DfxX7ql.png)

Also, as you can see, the user API key page is now moved there.

### Pro Tips

Pro tips are little informational cards which can be displayed everywhere in the web interface to inform users about new stuff or features most users might not know about *(like OTA)*.  
![](https://i.imgur.com/MSm8zrs.png)  
Dissmissed messages are stored to local sotrage, so they do not show again after dismission.  
![](https://i.imgur.com/jocKeTl.png)

## Minor Updates

- **Web Interface Responsiveness** [#172, #192]  
  The web interface is now more web responsive and for better usability on mobile devices.

- **Login Redirection** [#209]  
  When you are not logged in while opening *(for example, the user settings page)* in the web interface, you are automatically redirected to the login page. After login, you are now redirected back to the origin destination page *(the user settings page)*.

## Security Fixes

- The download of guild backups now requires the `sp.guild.admin.backup` permission, because the backup file
  contains confidential information like hidden channel details or guild settings. [#208]

## Bug Fixes

- The favicon of the web interface is now properly requestable.

## Backstage

- Add package [timerstack](https://pkg.go.dev/github.com/zekroTJA/shinpuru/pkg/timerstack).

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.9.0
```