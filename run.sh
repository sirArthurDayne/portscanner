#!/bin/bash

if [ -f ./portscanner ]; then
    echo "removing old exe"
    rm portscanner
fi

go build -race $1
./portscanner

