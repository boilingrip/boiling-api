#!/usr/bin/env bash

cd $GOPATH/src/github.com/boilingrip/boiling-api/cmd/boilingtest

while true; do
    go run main.go -config ~/boiling.yaml
    sleep 2
done