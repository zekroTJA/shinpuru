1.14.0

# Major Changes

## Guild Logs [#229]

Errors occuring in scope of guilds like, for example, failing backups, errors on setting autoroles, fails on sending messages in channels due to missing permissions, and so on, are now collected in the database and are visible in the guild admin panel.

Also, it will help spotting bugs in shinpuru by giving occuring errors more exposure.

You can also delete single or all entries of the error logs if there are any confidential information you do not want to be displayed in there.

![](https://i.imgur.com/z8NT1Vw.gif)

> Also, there are plans to display guild settings audit information in there like who changed when which guild setting.

## Guild Settings Relocation [#194]

General guild settings like guild prefix, autorole or join and leave messages are now moved to the guild settings route. 

The `Guild Settings` dropdown is now removed and you can now find a settings button next to the guild heading *(if you have permissional access to any guild settings)*, where you can now get to the guild settings route.

![](https://i.imgur.com/BimDz17.gif)

## Privacy Considerations [#229]

Because I try to take privacy very serious myself, I want to take a step forward to also do so with shinpuru.

There is now a new tab in the guild settings called `Data`. There, you are now able to file a database flush of all data correlated to the guild. 

This includes all reports and associated image data; all backups and associated backup files; all karma scores, settings, rules and blocklist; all starboard entries and configuration; all guild settings and permission specifications; tags; antiraid settings and joinlog and all unban requests.

![](https://i.imgur.com/savo6kH.png)

Also, you can now find a [Privacy Statement](https://github.com/zekroTJA/shinpuru/blob/master/PRIVACY.md) in shinpurus repository where I tried to point down as much details about which personal data is stroed in the shinpuru services as well as why they are stored and how they are stroed. Also I linked some contact information if you want to have data removed which is linked to your identity.

## Minor Changes

- Updated and unified some more frontend designs like capitalization of headings, for example.

## Bug Fixes

- Error messages in the web interface now actually give information about the error.
- Karma rules are now length-capped in the web interface so they will not overflow anymore.
- Fix the style of the pro tip overlays.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.14.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.14.0
```
