# `.github`

This directory contains files and configurations necessary for the GitHub repository host.

The `ISSUE_TEMPLATE` directory contains all issue templates definitions you can find when [creating an issue](https://github.com/zekroTJA/shinpuru/issues/new/choose).

In `workflows`, you can find all CI/CD pipelines which are executed by [GitHub actions](https://github.com/zekroTJA/shinpuru/actions). These pipelines do various tasks fully automated like 
- [checking that the unit tests succeed on each commit](workflows/tests-ci.yml)
- [create Docker images for canary, latest and release](workflows/docker-cd.yml)
- [creating the release page on tag push](workflows/releases-cd.yml)
- [generating the command wiki page](workflows/wiki-ci.yml)
- [generate the bughunters page](workflows/bughunters.yml)