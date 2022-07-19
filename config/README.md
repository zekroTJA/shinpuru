# `config`

This directory contains everything that is necessary for configuring shinpuru and adjacent services like [prometheus](prometheus/) or [grafana](grafana/).

[`config.example.yml`](config.example.yml) displays an example configuration of shinpuru with detailed documentation about each configuration key.

You can use [`my.private.config.yml`](my.private.config.yml) to set up a development configuration. Just copy the file to `private.config.yaml` and enter your credentials. This file is then automatically ignored by Git so that you do not accidentally leak your credentials.

The [`coder.private.config.yml`](coder.private.config.yml) can be used as same as the `my.private.config.yml` when you set up the coder workspace for shinpuru.