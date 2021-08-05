#!/bin/bash

PYTHON=py
which $PYTHON > /dev/null 2>&1 || PYTHON=python3

$PYTHON scripts/gen-package-descriptions.py
$PYTHON scripts/parse-gomod.py
$PYTHON scripts/md-replace.py -i README.md docs/public-packages.md docs/requirements.md
