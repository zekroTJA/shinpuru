# `ci`

This directory contains scripts that are used in CI/CD pipelines.

The [`build.sh`](build.sh) script generates the statically linked binaries for shinpuru for different operating systems and architectures as well as the static files for the web frontend. Also, it bundles everything into a handy compressed package and also provides hash sums. All this stuff is then attached as artifacts to the release page.