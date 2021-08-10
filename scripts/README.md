# `scripts`

This directory is a collection of various scripts to maintain this repository.

> The `db-updates` directory contains SQL patches which must have been applied manually in older versions of shinpuru after some patches before auto migration was implemented. So you can perfectly ignore this directory, but I still want to keep it if some people might need it for some reason. 

## Python Scripts

- [`bughunters.py`](bughunters.py): A script which uses the GitHub GraphQL API to get the number of issues and pull requests per member to create a scoreboard.
- [`gen-package-descriptions.py`](gen-package-descriptions.py): Detects all packages in the `/pkg` directory, obtains the package description from the source files and lists them in a markdown file.
- [`md-replace.py`](md-replace.py): Is used to inject the content of other markdown files into another markdown file by using comment markers. This script is primarily used to inject the `/docs/requirements.md` and `/docs/public-packages.md` into the main `README.md`.
- [`parse-gomod.py`](parse-gomod.py): Reads the `go.mod` file and lists all required packages from there into a markdown file.

## Bash Scripts

- [`commits-since-last-tag.sh`](commits-since-last-tag.sh): A script which lists all commits created after a the last tag. I use this to create the shinpuru changelogs.
- [`sass2scss.sh`](sass2scss.sh): This script has been used to convert the old SASS style sheets to SCSS style sheets.
- [`semver.sh`](semver.sh): Outputs a [semver](https://semver.org/) compliant version string using Git.
- [`update-fe-packages.sh`](update-fe-packages.sh): Automatically removes `node_modules/` and `package-lock.json` and executes `npm install` after to update all web app packages.
- [`update-readme.sh`](update-readme.sh): Executes `gen-package-descriptions.py` and `parse-gomod.py` and then uses `md-replace.py` to inject the generated markdown files into the main `README.md`.