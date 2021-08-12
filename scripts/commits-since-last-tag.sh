#!/bin/bash

git log $(git rev-parse $(git describe --tags --abbrev=0))..HEAD --oneline
