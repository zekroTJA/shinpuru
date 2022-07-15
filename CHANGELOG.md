[VERSION]

<!-- > **Attention**  
> This is a hotfix patch. If you want to see the changelog for release 1.30.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.30.0). -->

# Legacy Commands Removal

The until now marked as deprecated "legacy commands" *(meaning commands executed directly in the chat instead of using the new slash commands)* have now finally been removed from shinpuru. All commands now have to be executed using the slash commands. Therefore, corresponding middlewares and listeners have also meen removed.

Please refer to [the Command Wiki](https://github.com/zekroTJA/shinpuru/wiki/Commands) to get a brief introduction into all available slash commands and how to use them.

**If you dont see any slash commands on your guild, you need to kick shinpuru and re-invite the bot using the [provided invite link](https://shnp.de/invite)!** For explaination: The Bot must be invited with the scope `applications.commands` so that shinpuru can register applications on the guild. Because shinpuru did not use this scope before - and probably before you have invited it to your guild - the bot will not be able to create the application commands on your guild until you re-invite it with the new invite link containing the scope.

So, lets have a little moment of silence to remember the faithful legacy commands. 🙏

# Bug Fixes

- Fixed a bug where shinpuru would send the warn DM message to each guild admin every time a member joins after the triggering of the antiraid system. [#375]
- Fixed the download of the antiraid join list export. [#376]
- Karma of members on the karma block list is no longer altered when reacting to their messages with karma emotes. 
- Bots are now ignored by the verification system. [#377]
- When a member leaves during the process of creating a ban, the ban will be set to `Anonymous` implicitly instead of failing. [#378]
- You are now able to ban users with higher roles than you when being admin.
- *Finally* fixed this stupid error message popping up right after clicking on login in the web interface.

# Code Base

In the last couple of weeks, I started adding more and more unit tests around the core listeners and services of shinpuru. Generally, unit and integration tests allow way better and more consistent and continuous testing of important components of shinpuru. Also, it helps by going through all of shinpurus components and re-evaluating the implementation when constructing and implementing the tests. Actually, several bugs listed above have been discovered just by writing tests for the components with expected behaviours which were not met by the implementation. Even though it is a very hard and tideous process to implement "good" unit tests for each service, I think this will bring much improvement to the overall reliability and resilience of shinpuru - also over time when implementing more features.

Another important point is the value for potential contributors. For example, changes in implementations in pull requests which might cause unexpected behaviour somewhere else can be verified way quicker when catched by a failing unit test. Also, contributors (or those willing to) can use the unit tests to get an exemplary view of the usage of some services and functionalities.

So, I will try to implement unit tests for the most parts of the codebase of shinpuru and also do some changes to make the code better testable.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:[VERSION]
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:[VERSION]
```
