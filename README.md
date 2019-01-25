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

# Invite

Here you can choose between the stable or canary version of shinpuru:

<a href="https://discordapp.com/api/oauth2/authorize?client_id=524847123875889153&scope=bot&permissions=2080894065"><img src="https://img.shields.io/badge/%20-INVITE%20STABLE-0288D1.svg?style=for-the-badge&logo=discord" height="30" /></a>

<a href="https://discordapp.com/api/oauth2/authorize?client_id=536916384026722314&scope=bot&permissions=2080894065"><img src="https://img.shields.io/badge/%20-INVITE%20CANARY-FFA726.svg?style=for-the-badge&logo=discord" height="30" /></a>

> **Attention**<br>The canary version runs on the latest build pushed to the dev branch and can contain bugs! Also the canary version is running on a seperate database which is not included in my daily database backup.

# Intro

シンプル (shinpuru), a simple *(as the name says)*, multi purpose Discord Bot written in Go, using bwmarrin's package [discord.go](https://github.com/bwmarrin/discordgo) as API and gateway wrapper. The focus on this bot is not to punch in as much features and commands as possible, just some commands and features which I thought would be useful and which were the most used with my older Discord bots, like [zekroBot 2](https://github.com/zekroTJA/zekroBot2), and more on making this bot as reliable and stable as possible.

Also, I want to use this project as chance for me, to get some deeper into Go and larger Go project structures. In a later development state, this bot will detach zekroBot 2.

---

# Command list 

In this [**wiki article**](https://github.com/zekroTJA/shinpuru/wiki/Commands), you can find a automatically generated list of all commands and their manuals.

---

# Development state

This project is in a very early development state, so, currently, the bot is not available as 24/7 version. If you want to use this bot by yourself, pull the code by cloning the repository or download the assets as zip [here](https://github.com/zekroTJA/shinpuru/archive/master.zip).

Then, get all dependencies and build the binary. After that, generate a config by starting the bot, fill in your data and go on ;)

# Compiling

Read about self-compiling in the [**wiki article**](https://github.com/zekroTJA/shinpuru/wiki/Self-Compiling).

---

# Third party dependencies

- [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)
- [go-yaml/yaml](https://github.com/go-yaml/yaml)
- [go-sql-driver/mysql](https://github.com/Go-SQL-Driver/MySQL/)
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- [op/go-logging](https://github.com/op/go-logging)

Avatar of [御中元 魔法少女詰め合わせ](https://www.pixiv.net/member_illust.php?mode=medium&illust_id=44692506) from [瑞希](https://www.pixiv.net/member.php?id=137253).

---

Copyright (c) 2018-2019 zekro Developmenr (Ringo Hoffmann).  
Covered by MIT Licence.
