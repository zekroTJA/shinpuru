1.7.0

> MAJOR PATCH

## Major Implementations

### Unban Requests [#196]

Finally, shinpuru now has the ability that banned users can create unban requests when they were banned from a guild where the shinpuru moderation system is used.

Banned users can login to the shinpuru web interface and then navigate to `/unbanme`. There, they can select the guild(s) where they are banned and can submit an unban request.

![](https://i.imgur.com/9XxnUgH.gif)

Members with the permission `sp.guild.mod.unbanrequests` can review and process pending unbanreqeusts. If a request is accepted, the user is automatically unbanned from the guild.

![](https://i.imgur.com/kEJ6ETu.gif)

## Minor Updates

- **Database Migration**  
  Database modules can now implement the [`Migration`](internal/core/database/migration.go) interface which allows automatic database model migration on startup.  
  For example of the `mysql` database model: A `migrations` table is created which holds latest applied migrations. For each database update, a migration function can be supplied which can be applied one-by-one on the startup of shinpuru. See the [`mysql`](internal/core/database/mysql) module for more details.

## Bug Fixes

- Fixed output of `ToUnix` function of [timeutil](pkg/timeutil) package
- Fixed typo in report command description [#195]

## Backstage

- Added unit tests for various public packages

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.7.0
```