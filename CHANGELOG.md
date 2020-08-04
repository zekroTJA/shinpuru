1.2.0-beta.1

> MAJOR PREVIEW PATCH

> This patch includes following minor patches:
> - [release 1.1.0](https://github.com/zekroTJA/shinpuru/releases/tag/1.1.0)

## Major

- **Karma system implementation.** [#134]  
  Karma is a value which shall provide a scale of the trustworthyness of a user on the guild. The system works similar to the karma system of Reddit or Stackoverflow, for example.  
  You can gain karma when other users react to your messages with `👍, 👌, ⭐, ✔` and you lose karma when users react with `👎, ❌` to your message.  
  The value of karma is shown in the profile command and in the web interface. Also, you can view a scoreboard of the members with most karma.  
  You can read the full proposal here in this issue: #134. 
![](https://i.imgur.com/xia2aeN.png)
![](https://i.imgur.com/u4SX0lW.png)

## Minor

- Add edit flag to [say command](https://github.com/zekroTJA/shinpuru/wiki/Commands#say). [#142]

## Fixes

- Fix typo in security cards in web interface. [#140]
- Fix self-member button in web interface. [#141]

## Backstage

- Database drivers are now moved to the internal package `internal/core/middleware` to make middleware drivers available for usage as database and general purpose cache or other use cases.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.2.0-beta.1
```