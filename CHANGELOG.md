1.13.0

## Major Changes

### Web Inertace Redesign

Because the design of the web interface of shinpuru is kind of inconsistent and also kind of unappealing, I decided to re-design it step by step. The first results of this re-design offensive you can see in this release.

First of all, I changed the heading front to [`Cantarell`](https://fonts.google.com/specimen/Cantarell). It's clean and simple though it has some character to it which perfectly matches the general design of shinpuru. Also, it let's the headings stand out more and better highlights the important parts of the UI.

<img width="49%" src="https://i.imgur.com/zYhahPT.png"/><img width="49%" src="https://i.imgur.com/OeEWrAu.png"/>

The old design always gave me the feel of a rough, cluttered experience. So, as you can see, everything got a bit more round, smooth and spacy.

<img width="49%" src="https://i.imgur.com/sB3Skt6.png"/><img width="49%" src="https://i.imgur.com/DvIpHt7.png"/>

<img width="49%" src="https://i.imgur.com/JCl9RSg.png"/><img width="49%" src="https://i.imgur.com/Q6Uqj9O.png"/>

As well, I adobted the new design CI of Discord with the new color tones and the optimized logo.

![](https://i.imgur.com/yveJbqZ.png)

### Karma Rules [#231]

You can now specify karma rules which will be automatically applied depending on the karma levels of the members.

To set them up, navigate to the guild admin panel of your guild, then go to `Karma` and scroll down to `Rules`.

> Please be careful using this feature, especially with the kick and ban rules, because it did not pass through the full test period yet!

![](https://i.imgur.com/knRJ0n5.png)

<!-- ## Minor Changes -->


<!-- ## Bug Fixes -->


# Docker

[Here](https://hub.docker.com/r/zekro/shinpuru) you can find the docker hub page of shinpuru and [here](https://github.com/zekroTJA?tab=packages&repo_name=shinpuru) you can find Docker images released on the GHCR.

Pull the docker image of this release:

From DockerHub:

```
$ docker pull zekro/shinpuru:1.13.0
```

From GHCR:

```
$ docker pull ghcr.io/zekrotja/shinpuru:1.13.0
```
