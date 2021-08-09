# `config`

This directory contains everything that is necessary for configuring shinpuru and adjacent services like [prometheus](prometheus/) or [grafana](grafana/).

[`config.example.yaml`](config.example.yaml) displays an example configuration of shinpuru with detailed documentation about each configuration key.

You can use [`my.private.config.yaml`](my.private.config.yaml) to set up a development configuration. Just copy the file to `private.config.yaml` and enter your credentials. This file is then automatically ignored by Git so that you do not accidentally leak your credentials.