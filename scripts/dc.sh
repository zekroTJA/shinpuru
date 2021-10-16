#!/bin/sh

DEV_COMPOSE_FILE="docker-compose.dev.yml"

docker-compose -f $DEV_COMPOSE_FILE $@
