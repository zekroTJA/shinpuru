[VERSION]

<!-- > **Attention**  
> This is a hotfix patch for issue #332. If you want to see the changelog for release 1.26.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.27.0). -->

# Privacy

A lot of stuff was added and has been changed to improve the privacy of the users of shinpuru as well as making the data usage and removal more easy and transparent for the users.

## New Privacy Policy

For this update, I've completely re-written the current privacy policy of shinpuru to make it more transparent what data is stored for how long and for what specific purpose it is stored.

[**Here**](https://github.com/zekroTJA/shinpuru/blob/master/PRIVACY.md) you can read the new privacy policy.

## User Data Removal [#335]

You can now request a user data removal directly from the web interface which invokes a deletion of the data linked to your Discord User ID stored in the database as well as in the cache *(see screenshot below)*. Of course, this excludes entries in the database which are essential to the security systems like report entries created against you or negative karma scores.

## Starboard Opt-Out [#338]

In the user privacy settings, you are also now able to globally opt-out of the starboard system. When enabled, your messages will not be presented in a guilds starboard channel or web interface panel.

![](https://user-images.githubusercontent.com/16734205/149481912-0ebc7d0a-432a-4820-a8a6-678c5654435f.png)

## Minor privacy related changes

- It is now ensured that the state cache duration for privacy critical data is capped at 30 days.
- When flusing guild data, this now also includes data stored in the state cache.
- Open votes are now also cleared from a guild when guild data is flushed, which also removes the vote entries from the database.
- Sometimes, after a restart, the antiraid joinlist would not be cleared after 48 hours. This is now enforced by a cleanup lifecycle routine.
- Privacy information and contact details are now shown in the `/info` command as well as in the `/info` view in the web interface.

## ⚠️ Important for Self Hosters

You now **must** provide a privacy notice as well as contact information in the config of shinpuru. Otherwise, the startup will be prohibited. You can do this by adding the following to your configuration.

```yaml
# Privacy information and contact details
# which are shown in the /info command as well
# as in the web interface.
privacy:
  # URL to your privacy notice.
  # DO NOT USE THE NOTICE BELOW BECAUSE IT IS ONLY
  # VALID FOR THE OFFICAL HOST OF SHINPURU!
  noticeurl: https://github.com/zekroTJA/shinpuru/blob/master/PRIVACY.md
  # Contact details.
  contact:
      # Title of the contact type
    - title: E-Mail
      # The displayed value
      value: contact@example.de
      # An optional link URL
      url: "mailto:contact@example.de"
```

# Code Execution

## Rate Limit Configuration [#329]

You can now configure the rate limiting of the code execution in the configuration file of shinpuru.

```yml
codeexec:
  type: ranna
  ranna:
    apiversion: v1
    endpoint: 'https://public.ranna.zekro.de'
    token: ''
  ratelimit:
    # Whether or not to enable the rate limiting.
    enabled: true
    # The burst rate of the limiter.
    burst: 5
    # The time in seconds between regeneration
    # of rate limiter tokens.
    limitseconds: 60
```

## Configure Code Execution via the Web Inetrface [#330, #265]

You are now also able to configure code execution in the web interface. That also allows to disable code execution on a per-guild basis.

*When using ranna as CEE:*  
![](https://user-images.githubusercontent.com/16734205/149314012-92915c7c-eaf6-4686-b141-3be491d9088b.png)

*When using JDoodle as CEE:*  
![](https://user-images.githubusercontent.com/16734205/149314312-f0637a85-dbb8-42c1-a917-20e5eacc6de7.png)

As you can see above, this now also allows to set and reset JDoodle credentials directly form the web interface so you don't need to use `/exec setup` therefore anymore!

# Legacy Command Deprecation Warning

When you use a legacy ("non-slash command"), you will now be presented with a deprecation warning (not every time, but once per day). If you want to read more about legacy command deprecation, please read [this wiki article](https://github.com/zekroTJA/shinpuru/wiki/Legacy-Command-Deprecation).

![](https://user-images.githubusercontent.com/16734205/149315269-cd0f65ec-0235-4bc8-89d3-5baf7c053ed2.png)

# Bug Fixes

- Fixed a bug where the bot crashes when using the `/quote` or `/user` slash command. [#334]
- Switched to using `string` instead of `int` parameter types for snowflake input in slash command because `int` types will show as invalid when too long.
- Fixed command info type mapping. [#337]
- Fixed a bug where the command info view does not load.

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
