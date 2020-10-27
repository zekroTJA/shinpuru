1.4.0

> MAJOR PATCH

## Major Implementations

### Antiraid System [#159]

A new guild security feature has been added to shinpuru: The Antiraid System.

> **What is a "Raid"?**  
> A guild raid is mostly refered to a large, burst-like ammount of accounts joining the guild in a short period of time. This is mostly caused by a single user or a group of users which utilize bot-created or hijacked accounts to flood a guild.

To counteract this, the antiraid system constantly checks the rate of users joining your guild. If the rate increases over a certain threshold, the antiraid system triggers. Following, the guilds security level is raised to `verry high` and for the following 24 hours, all users joining the guild are logged in a list which is accessable via the web interface. Also, all admins of the guild will be informed about the incident.

Of course, the antiraid system can be toggled and the trigger threshold values can be managed in the web interface *(if you have the `sp.guild.config.antiraid` permission)*.  
![](https://i.imgur.com/vLMgrM9.png)

### Metrics Monitoring [#170]

You are now able to monitor core metrics of shinpuru using Prometheus and Grafana.

You can enable the prometheus scraping endpoint by adding this to your shinpuru config:
```yml
metrics:
  enable: true
  addr: ":9091"
```

[Here](https://github.com/zekroTJA/shinpuru/blob/master/config/prometheus/prometheus.yml) you can find an example Prometheus configuration and [here](https://github.com/zekroTJA/shinpuru/blob/master/config/grafana/example-dashboard.json) you can find an example grafana dashboard to monitor shinpuru's metrics.  

*Example dashboard. Data from shinpuru Canary instance.*  
![](https://i.imgur.com/fEkV7fe.png)


## Minor Updates

- Add aliases to `karma` command: `leaderboard`, `lb`, `sb` and `top`. [#181]
- The `karma` command now shows the karma points of a user when specified as argument. [#179]

## Bug Fixes

- The web frontend route `/guilds/:guildid/guildadmin` now redirects to `/guilds/:guildid/guildadmin/antiraid` instead of firing errors. [#180]


# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.4.0
```