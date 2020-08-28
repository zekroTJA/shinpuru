1.2.0

> MAJOR PATCH

## Major

- **Karma system implementation.** [#134, #145]  
  Karma is a value which shall provide a scale of the trustworthyness of a user on the guild. The system works similar to the karma system of Reddit or Stackoverflow, for example.  
  You can gain karma when other users react to your messages with `👍, 👌, ⭐, ✔` and you lose karma when users react with `👎, ❌` to your message.  
  The value of karma is shown in the profile command and in the web interface. Also, you can view a scoreboard of the members with most karma.  
  You can read the full proposal here in this issue: #134.  
  ![](https://i.imgur.com/xia2aeN.png)
  ![](https://i.imgur.com/9sROCVi.png)

- **shireikan implementation.** [#152]  
  [shireikan](https://github.com/zekroTJA/shireikan) is a command handler package which replaces the internal command handler of shinpuru. Read more about the advantages of this implementation in issue #152. In order of this implementation, a lot of commands and modules needed to be refactored.

- **Web frontend changes.** [#151, #153]  
  The layout of the front end is now kept at a max width and centered which leads to a way more clear and nice-looking design. Also, it is now possible to revoke reports via the web interface.  
  ![](https://i.imgur.com/7DXTeXL.png)

- **Add color reactions.** [#155]  
  Color reactions are a system which, when enabled by the [`color`](https://github.com/zekroTJA/shinpuru/wiki/Commands#color) command, scrapes messages for hexadecimal color codes. Then, a reaction is added which shows the color. After clicking the reaction, more information about the color is shown.  
  ![](https://i.imgur.com/VICm9BV.gif)

- **Command overview in web interface.** [#158]  
  Add a command list in the web interface where you have a clear overview over all commands of shinpuru and how they are used.  
  ![](https://i.imgur.com/sTHzdEN.gif)

## Minor

- Add edit flag to [say command](https://github.com/zekroTJA/shinpuru/wiki/Commands#say). [#142]
- Update Header in Web Interface which is now static at the top of the window and has a drop shadow for better visual seperation.
- Optimize permission role input in web interface. [#148]
- Add fuzzy search for help command. [#157]

## Fixes

- Fix typo in security cards in web interface. [#140]
- Fix self-member button in web interface. [#141]
- Unify command descriptions.

## Backstage

- Database drivers are now moved to the internal package `internal/core/middleware` to make middleware drivers available for usage as database and general purpose cache or other use cases.
- Add service which starts the Angular dev server alongside with shinpuru when passing the `-devmode` flag on start.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.2.0-rc.1
```