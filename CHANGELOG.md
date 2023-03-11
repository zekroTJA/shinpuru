[VERSION]

<!-- > **Attention**  
> This is a hotfix patch. If you want to see the changelog for release 1.30.0, please look [**here**](https://github.com/zekroTJA/shinpuru/releases/tag/1.30.0). -->

## New Web Interface

This release finally brings a huge update to the web interface of shinpuru. Because the old web interface had no clear style concept while also growing with each new feature, it got more and more cluttered, unstructured, obscure and simply uglier. So I decided [almost a year ago](https://github.com/zekroTJA/shinpuru/issues/370) to rewrite the whole web interface, which has now come to the final stage. There is still a lot to do and - especially translation-wise - a lot missing, but the feature set is now 100% ported and so I decided to finally port it over.



## New Logger

To improve the logs of shinpuru both in visibility as well as in flexibility, I've created my own logging package called [rogu](https://github.com/zekroTJA/rogu). It allows colorful, human readable, taggable, strctured logging with a simple API to append multiple output writers.

![](https://user-images.githubusercontent.com/16734205/222913731-86c08d45-e769-49f2-96f1-a19adf1eda9e.png)

An additional [output writer](https://github.com/zekroTJA/shinpuru/tree/master/pkg/lokiwriter) has been written for pushing logs to [Grafana Loki](https://github.com/grafana/loki) which allows central log aggregation for multiple instances of shinpuru. Simply add the following config to your logging config to enable loki log pushing.

```yml
# Logging preferences
logging:
  # Set the log level of the logger
  # Log levels can be found here:
  # https://github.com/zekroTJA/rogu/blob/main/level/level.go
  loglevel: 4
  # Specify Grafana Loki configuration
  # for log aggregation
  loki:
    # Whether to enable sending logs to loki or not
    enabled: true
    # The address of the loki instance
    address: "https://loki.example.com"
    # The basic auth user name (leave empty if not used)
    username: "username"
    # The basic auth password (leave empty if not used)
    password: "2374n8er7nt8034675782345"
    # Additional labels set to all log entries.
    labels:
      # Some examples ...
      app: "shinpuru"
      instance: "main"
```

The provided [example Grafana Dashboard](config/grafana/example-dashboard.json) shows how aggregated logs can be visualized in Grafana.

![image](https://user-images.githubusercontent.com/16734205/222915283-41e6a6c7-6497-451e-8a83-a7eaa6a6bdd7.png)

## PushCode Login

Because there is a potential risk that the pushcode login system could be abused by attackers to phish login sessions, a confirmation promt has been added with a warning that you should **never** enter a login code to shinpuru's DMs which you have received from someone else (see issue #412).

![](https://user-images.githubusercontent.com/16734205/222915580-09db7f99-6a44-480d-bd5c-ea5905fca67b.png)


## API Changes

- New API Endpoint [`GET /allpermissions`](https://app.swaggerhub.com/apis-docs/zekroTJA/shinpuru-main-api/1.0#/Etc/get_allpermissions) which returns a list of all available permissions.
- New API Endpoint [`GET /healthcheck`](https://app.swaggerhub.com/apis-docs/zekroTJA/shinpuru-main-api/1.0#/Etc/get_healthcheck) which can be requested to get the health state of shinpuru services.
- New API Endpoint [`GET /guilds/{id}/starboard/count`](https://app.swaggerhub.com/apis-docs/zekroTJA/shinpuru-main-api/1.0#/Guilds/get_guilds__id__starboard_count) to retrieve the total count of starboard entries for a given guild.
- New API Endpoint [`GET /guilds/{id}/unbanrequests/count`](https://app.swaggerhub.com/apis-docs/zekroTJA/shinpuru-main-api/1.0#/Guilds/get_guilds__id__unbanrequests_count) to retrieve the total count of unbanrequests for a given guild.
- Update API Endpoint [`POST /guilds/{id}/permissions`](https://app.swaggerhub.com/apis-docs/zekroTJA/shinpuru-main-api/1.0#/Guilds/post_guilds__id__permissions) which now returns the resulting updated permissions map.

## Docker Image

The docker image now includes a healthcheck which shows and monitors the state of the shinpuru instance using the [`GET /healthcheck`](https://app.swaggerhub.com/apis-docs/zekroTJA/shinpuru-main-api/1.0#/Etc/get_healthcheck) API endpoint.

## Etc

- The state cache duration for users and members has now be increased from 30 days to 90 days for better performance.

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
