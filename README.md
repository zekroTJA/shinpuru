<div align="center">
    <img src="https://zekro.de/src/shinpuru_avi_circle.png" height="300" />
    <h1>~ シンプル ~</h1>
    <strong>
        A siple multi purpose discord bot written in Go (discord.go)<br>
        with focus on stability and reliability
    </strong><br><br>
    <a href="https://dc.zekro.de"><img src="https://img.shields.io/discord/307084334198816769.svg?style=for-the-badge&logo=discord" /></a>&nbsp;
    <img src="https://forthebadge.com/images/badges/60-percent-of-the-time-works-every-time.svg" />&nbsp;
    <img src="https://forthebadge.com/images/badges/built-with-grammas-recipe.svg">
<br>
</div>

---

## Intro

シンプル (shinpuru), a simple *(as th name says)*, multi purpose Discord Bot written in Go, using bwmarrin's package [discord.go](https://github.com/bwmarrin/discordgo) as API and gateway wrapper. The focus on this bot is not to punch in as much features and commands as possible, just some commands and features I thought wich would be useful and which were the mostly used from my older Discord bots, like the [zekroBot 2](https://github.com/zekroTJA/zekroBot2), and more on making this bot as reliable and stable as possible.

Also, I want to use this project as chance for me, to get some deeper into Go and larger Go project structures. In a later development state, this bol will detach zekroBot 2.

---

# Development state

This project is in a verry early development state, so, currently, the bot is not available as 24/7 version. If you want to use this bot by yourself, pull the code by cloning the repository or download the assets as zip [here](https://github.com/zekroTJA/shinpuru/archive/master.zip).

Then, get all dependencies and build the binary. After that, generate a config by starting the bot, fill in your data and go on ;)

```
$ git clone https://github.com/zekroTJA/shinpuru.git
$ bash ./installdeps.sh
$ go build -o shinpuru ./src/
$ ./shinpuru -c yourconfig.yaml
```

And yes, as you may notice, this bot currently depends on a MySql database. Actually, I want to make this bot also compatible with SqLite or MongoDB, but this will take a while until this is most priority. ^^

---

# Third party dependencies

- [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)
- [go-yaml/yaml](https://github.com/go-yaml/yaml)
- [go-sql-driver/mysql](https://github.com/Go-SQL-Driver/MySQL/)

Avatar of [御中元　魔法少女詰め合わせ](https://www.pixiv.net/member_illust.php?mode=medium&illust_id=44692506) from [瑞希](https://www.pixiv.net/member.php?id=137253).

---

Copyright (c) 2018 zekro Developmenr (Ringo Hoffmann).  
Covered by MIT Licence.
