#!/bin/bash

MYSQLSCHEME=$(cat scripts/mysqlDbScheme.sql | base64 -w 0)
SQLITESCHEME=$(cat scripts/sqliteDbScheme.sql | base64 -w 0)

SQLLDFLAGS="\
    -X github.com/zekroTJA/shinpuru/core.MySqlDbSchemeB64=$MYSQLSCHEME \
    -X github.com/zekroTJA/shinpuru/core.SqliteDbSchemeB64=$SQLITESCHEME "