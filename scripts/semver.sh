#!/bin/bash

VERSION=$(git describe --tags --abbrev=0)
NCOMMS=$(git log $VERSION..HEAD --oneline | wc -l)
LASTCOMM=$(git rev-parse HEAD)

[ "$NCOMMS" == 0 ] \
    && echo "$VERSION" \
    || echo "$VERSION+$NCOMMS-${LASTCOMM:0:6}"