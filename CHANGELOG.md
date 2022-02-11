[VERSION]

<!-- > **Attention**  
> This is a hotfix patch for issue #332. If you want to see the changelog for release 1.26.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.27.0). -->

# Birthday Notifications [#349]

shinpuru now features birthday notifications!

![](https://user-images.githubusercontent.com/16734205/153576590-28b51ce9-e11f-4aa1-b86b-41fc7d6d6a31.gif)

Simply, set your birthday using the `/birthday set` slash command. The date must be in the format `YYYY-MM-DD[+-]TIMEZONE`. You can also use `/` or `.` as separator. The timezone must be set to your local time zone in offset hours from UTC. So, for example, for `CET` this would be `+1`, for `EST`, this would be `-5` and so on. 

It is also possible to just set a month and day if you don't want to store your birth year. If you want to show your age in the notification, you can enable it by setting `show-year` to `True`. Otherwise, your age will be hidden in the notification.

![](https://user-images.githubusercontent.com/16734205/153576393-73616bd2-b21d-4813-bdb8-ad146cecc542.png)

*With age hidden, it would look like this.*  
![](https://user-images.githubusercontent.com/16734205/153580431-e5a0f4b9-1b51-473f-bfa4-40c3fa650c89.gif)

You can, of course, also unset your birthday at any time.

![](https://user-images.githubusercontent.com/16734205/153577075-e191e2d8-1e39-4ab4-9a24-31e5939d887f.png)

To enable birthday notifications in your guild, you need to specify a birthday channel. This requires the permission `sp.guild.config.birthday`

![](https://user-images.githubusercontent.com/16734205/153576456-9c184ae0-4408-4ded-9dce-9f69ee36f5e9.png)

To unset this setting, simply use the `/birthday unset-cannel` command. 

To enable Gifs in the birthday message, the bot needs a Giphy API key. You can get one by creating a Giphy acount and going to the [Developer Dashboard](https://developers.giphy.com/dashboard/). There you can create an app. After that, copy the API key and add it ti shinpuru's config.
```yaml
giphy:
  apikey: dWl3ZXF6diBzZGR0NnczNDg5NTZuZG4w
```

# Sharding [#238]

If you are running shinpuru of a lot of Guilds, you might want to split up your single instance into multiple instances to split up the load. This is now possible using Discord's Gateway sharding.

> I **strongly** recommend taking a look into [Discord's Documentation](https://discord.com/developers/docs/topics/gateway#sharding) about sharding when you want to split up your instance.

You can simply spin up multiple instances of shinpuru behind a load balancer which all connect to the same database and redis instance. This distributes a common synced persistent state between all instances.

If you want to set up sharding and load balancing, you can find more information on how to set up and configure shinpuru [**here**](https://github.com/zekroTJA/shinpuru/tree/dev/docs/sharding). There you can also find an example deploymet using docker swarm.

# Bug Fixes

- Fix a bug where guild settings were not saved to database.
- Properly bubble up errors when setting guild settings.

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
