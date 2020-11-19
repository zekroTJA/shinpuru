1.5.0

> MAJOR PATCH

## Major Implementations

### Landing Page

shinpurus `/login` page is now decorated by a landing page which shows off some of shinpurus core features, some invite links and options to self host as well as some links to dive deeper.

> This page is still in a kind of **beta state**. A lot of stuff is still missing like proper support for mobile devices as well as further feature spotlights.

![](https://i.imgur.com/4V6VVab.gif)

## Minor Updates

- Two imporvements of the color feature:  
  1. A name of the color which is closest to the specified color is now displayed. This is provided by the [`zekroTJA/colorname`](https://github.com/zekroTJA/colorname) package.
  2. The name of the embed executor is now displayed in the embed footer. [#183]

  ![](https://i.imgur.com/4dzBN8z.png)

- The web server endpoint `/invite` now redirects to the invite link of the current shinpuru instance (e.g. https://shnp.de/invite).

## Bug Fixes

- Fix hex notation of color reaction embeds.

## Backstage

- Move `stringutils` package to `pkg/stringutils`.

# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru.

Pull the docker image of this release:
```
$ docker pull zekro/shinpuru:1.4.0
```