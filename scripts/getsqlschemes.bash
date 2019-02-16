#!/bin/bash

MYSQLSCHEME=$(cat scripts/mysqlDbScheme.sql | base64 -w 0)
SQLITESCHEME=$(cat scripts/sqliteDbScheme.sql | base64 -w 0)

echo "\
    -X github.com/zekroTJA/shinpuru/internal/core.MySqlDbSchemeB64=$MYSQLSCHEME \
    -X github.com/zekroTJA/shinpuru/internal/core.SqliteDbSchemeB64=$SQLITESCHEME "