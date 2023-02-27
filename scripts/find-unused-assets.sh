#!/bin/bash

OUTPUT_DIR=".unused-assets"

which rg > /dev/null || {
  echo "error: rg (ripgrep) is not installed"
  exit 1
}

[ -z "$1" ] && {
  echo "usage: $(basename $0) [path]"
  exit 1
}

for f in $(find web/src/assets -type f); do
  rg "'.*${f#web\/src\/*}'" "$1" > /dev/null || {
    echo "moving $f ..."
    BASE="${f#web\/src\/assets\/*}"
    DIR=$(dirname $BASE)
    mkdir -p "$OUTPUT_DIR/$DIR"
    mv "$f" "$OUTPUT_DIR/$BASE"
  }
done
