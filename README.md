<div align="center">
    <img src="https://zekro.de/src/shinpuru_avi_circle.png" height="300" />
    <h1>~ シンプル ~</h1>
    <strong>
        A simple multi purpose discord bot written in Go (discord.go)<br>
        with focus on stability and reliability
    </strong><br><br>
    <a href="https://dc.zekro.de"><img height="28" src="https://img.shields.io/discord/307084334198816769.svg?style=for-the-badge&logo=discord" /></a>&nbsp;
    <a href="https://github.com/zekroTJA/shinpuru/releases"><img height="28" src="https://img.shields.io/github/tag/zekroTJA/shinpuru.svg?style=for-the-badge"/></a>&nbsp;
    <img height="28" src="https://forthebadge.com/images/badges/60-percent-of-the-time-works-every-time.svg" />&nbsp;
    <img height="28" src="https://forthebadge.com/images/badges/built-with-grammas-recipe.svg">
<br>
</div>

---

| Branch | Build |
|--------|-------|
| master | <a href="https://travis-ci.org/zekroTJA/shinpuru"><img src="https://travis-ci.org/zekroTJA/shinpuru.svg?branch=master" /></a> |
| dev | <a href="https://travis-ci.org/zekroTJA/shinpuru"><img src="https://travis-ci.org/zekroTJA/shinpuru.svg?branch=dev" /></a> |

---

## Intro

シンプル (shinpuru), a simple *(as the name says)*, multi purpose Discord Bot written in Go, using bwmarrin's package [discord.go](https://github.com/bwmarrin/discordgo) as API and gateway wrapper. The focus on this bot is not to punch in as much features and commands as possible, just some commands and features which I thought would be useful and which were the most used with my older Discord bots, like [zekroBot 2](https://github.com/zekroTJA/zekroBot2), and more on making this bot as reliable and stable as possible.

Also, I want to use this project as chance for me, to get some deeper into Go and larger Go project structures. In a later development state, this bot will detach zekroBot 2.

---

# Development state

This project is in a very early development state, so, currently, the bot is not available as 24/7 version. If you want to use this bot by yourself, pull the code by cloning the repository or download the assets as zip [here](https://github.com/zekroTJA/shinpuru/archive/master.zip).

Then, get all dependencies and build the binary. After that, generate a config by starting the bot, fill in your data and go on ;)

## Compiling

For compiling, you will need:
- git
- go
- gcc (if you are on windows, use the [TDM-GCC toolchain](https://sourceforge.net/projects/tdm-gcc/))

```
$ git clone https://github.com/zekroTJA/shinpuru.git src/github.com/zekroTJA/shinpuru
$ export GOPATH=$PWD
$ cd src/github.com/zekroTJA/shinpuru
$ bash scripts/build.sh
$ ./shinpuru -c yourconfig.yaml
```

**Important:** For getting shinpuru working properly, you will need to use the bild script. If you are on windows, execute it in the git bash or with WSL.

The bot currently supports MySql and SQLite as database.

---

# Third party dependencies

- [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)
- [go-yaml/yaml](https://github.com/go-yaml/yaml)
- [go-sql-driver/mysql](https://github.com/Go-SQL-Driver/MySQL/)
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- [op/go-logging](https://github.com/op/go-logging)

Avatar of [御中元 魔法少女詰め合わせ](https://www.pixiv.net/member_illust.php?mode=medium&illust_id=44692506) from [瑞希](https://www.pixiv.net/member.php?id=137253).

---

Copyright (c) 2018 zekro Developmenr (Ringo Hoffmann).  
Covered by MIT Licence.
