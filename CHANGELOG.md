1.10.0

## New Features

### Starboard [#212]

shinpuru now has a starboard. With the command [`starboard`](https://github.com/zekroTJA/shinpuru/wiki/Commands#starboard), you can configure the properties for the starboard on your guild. After that, messages which receive the configured amount of ⭐ reactions will be highlighted in the starboard as well as in the web interface.

![](https://i.imgur.com/d5S9vm4.png)
![](https://i.imgur.com/ohJ9z6U.png)

### Voicelog Blocklist [#125]

The [`voicelog`](https://github.com/zekroTJA/shinpuru/wiki/Commands#voicelog) command has been extended by an option to block channels from being visible in the voice log channel.

## Minor Changes

- The issuer of the [`karma`](https://github.com/zekroTJA/shinpuru/wiki/Commands#starboard) command is now shown in the footer of the embed. [#211]
- When the guild join/leave message contains a mention, the mention is also contained in the message itself outside the embed so that the member is getting pinged. [#210]
- You can now specify a global command rate limiting in the config. [Here](https://github.com/zekroTJA/shinpuru/blob/6325dbfc9b042d5eb338fa2b80a0c2e75fd69ab0/config/config.example.yaml#L24-L29) you can find an example configuration for that.

## Bug Fixes

- The color reaction module does not anymore delete any reaction when not having permissions to execute color reactions. [#213]

## Backstage

- Package [onetimeauth](https://pkg.go.dev/github.com/zekroTJA/shinpuru/pkg/onetimeauth) is now publicly available.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.9.1
```