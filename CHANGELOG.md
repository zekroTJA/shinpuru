1.23.0

# Changes

## Slash Command Implementation [#287]

> If you want to read the full story, please take a look into this issue: #287.

TLDR: Since Discord will require a privileged intent for message content access after April 2022 (see [this article](https://support-dev.discord.com/hc/en-us/articles/4404772028055-Message-Content-Access-Deprecation-for-Verified-Bots) for more information), the current command system will not work anymore after that. So [ken](https://github.com/zekrotja/ken) was created as a Discordgo slash command framework and all commands were ported to the new system.

And because slash commands are so well integrated with the Discord chat, this also really improves the user experience when interacting with shinpuru.

Due to the changes made, following adjustments were made:

- The `report`, `kick` and `ban` command are now combined into the single slash command `/report` where you can specify the `type` as a parameter.
- The `game` command is now renamed to `/presence` as slash command.
- The `joinmsg` and `leavemsg` commands are now combined into the single slash command `/announcements` where you can specify the `type` as a parameter.

Also, some commands were not ported to the new command system and will be removed subsequently.
- The `ment` command is now obsolete because admins can now mention roles even if they are marked as not mentionable.
- The `prefix` command is now obsolete due to slash commands do not require nor allow custom guild prefixes.
- The `help` command is now obsolete because command information is directly displayed in the Discord chat when using slash commands.

The legacy command system is from now marked as **deprecated** and will be fully removed in following updates. To be able to use slash commands, you must kick shinpuru from your guild and re-invite the bot. This is due to a new OAuth2 scope which is required for a bot to be able to register slash commands on your guild.

## ⚠️ Permission Adjustments

Due to the slash command implementation, some permission domain names have changed to maintain consistency. **This will require guild administrators to adjust your permission settings accordingly.**

- `sp.guild.config.joinmsg` and `sp.guild.config.leavemsg` is now combined into **`sp.guild.config.announcements`**.
- `sp.chat.exec` has been changed to **`sp.guild.config.exec`**.
- `sp.guild.config.stats` has been changed to **`sp.guild.config.starboard`**.
- `sp.game` has been changed to **`sp.presence`**.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.23.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.23.0
```
