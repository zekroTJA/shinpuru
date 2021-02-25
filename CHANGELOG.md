1.8.0

## Minor Updates

- **Karma Blocklist [#201]**  
  You can now block users from the karma system in the guild karma settings. Blocked users are not able to give or remove karma of other members and are also unable to gain karma by other users.  
  ![](https://i.imgur.com/ZokXOze.gif)

- **Web Interface: Permission Input Autocomplete [#203]**  
  The permission input field in the guild settings of the web interface now has auto complete. 🎉  
  ![](https://user-images.githubusercontent.com/16734205/109003709-13944700-76a8-11eb-92db-1eff56d1b520.gif)

- **Ban message now contains unban link [#204]**  
  The ban message, which is sent to a user banned from a guild with shinpuru's ban command, now contains a mention to the link to the unban request form in the web interface to submit unban requests.  
  ![](https://i.imgur.com/XilUjgV.png)

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
$ docker pull zekro/shinpuru:1.8.0
```