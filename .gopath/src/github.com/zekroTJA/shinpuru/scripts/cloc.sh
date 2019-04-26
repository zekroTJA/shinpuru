#!/bin/bash

cloc . --md > ./cloc.md
git add .
git commit -m "updated cloc.md"
git push
