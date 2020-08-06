1.2.0-beta.2

> MAJOR PREVIEW PATCH

> This is a prerelease issue and the changelog displays all changes since the last release [`1.1.1`](https://github.com/zekroTJA/shinpuru/releases/tag/1.1.1).

## Major

- **Karma system implementation.** [#134, #145]  
  Karma is a value which shall provide a scale of the trustworthyness of a user on the guild. The system works similar to the karma system of Reddit or Stackoverflow, for example.  
  You can gain karma when other users react to your messages with `👍, 👌, ⭐, ✔` and you lose karma when users react with `👎, ❌` to your message.  
  The value of karma is shown in the profile command and in the web interface. Also, you can view a scoreboard of the members with most karma.  
  You can read the full proposal here in this issue: #134. 
![](https://i.imgur.com/xia2aeN.png)
![](https://i.imgur.com/9sROCVi.png)

## Minor

- Add edit flag to [say command](https://github.com/zekroTJA/shinpuru/wiki/Commands#say). [#142]
- Update Header in Web Interface which is now static at the top of the window and has a drop shadow for better visual seperation.
- Optimize permission role input in web interface. [#148]

## Fixes

- Fix typo in security cards in web interface. [#140]
- Fix self-member button in web interface. [#141]

## Backstage

- Database drivers are now moved to the internal package `internal/core/middleware` to make middleware drivers available for usage as database and general purpose cache or other use cases.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.2.0-beta.2
```