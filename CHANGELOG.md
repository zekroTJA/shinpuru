1.25.0

# Minor Changes

- All slash commands which take a channel argument now specify the `channel type` specification so that only accepted channel types are selectable.  
  ![image](https://user-images.githubusercontent.com/16734205/142614849-ce795a13-fe59-4096-862b-83078dbb072f.png)

- The buttons `Kick` and `Ban` in the antiraid joinlog are now disabled when no entry is selected. [#307]

- In the embed builder, now only channels which the user has write permission to are shown.

- Also removed the `#` prefix of channel names in the embed builder for better searchability in the select input.

# Bug Fixes

- Bot owners now will get full permissions regardless of guild specification.
- Fix antiraid joinlog table structure. [#311]
- The report revocation of bans now properly unbans the banned user. [#303]
- Fix permission check on sending messgaes using the embed builder. [#309]
- Fix default time formatting in web interface. [#306]

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.25.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.25.0
```
