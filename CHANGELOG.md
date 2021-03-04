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

<!-- ## Bug Fixes

-  -->

## Backstage

- Update API of package [multierror](https://pkg.go.dev/github.com/zekroTJA/shinpuru/pkg/multierror).
- Refactored all around the report packages to now use a proper type "enum" and for a bit more logical package structure.  
<sub>(Tho, the general package structure of shinpuru is still horrible and really really really needs a complete refactor process...)</sub>
- Frontend updated to Angular 11.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.9.0
```