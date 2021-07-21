#!/bin/bash

WEB_DIR=$PWD/web

for f in $(find $WEB_DIR -name '*.sass'); do 
    cleanname=${f%.*}
    echo Converting $f ...
    sass2scss -c -p -p < $f > ${cleanname}.scss
    rm $f
done