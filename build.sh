#!/bin/bash

echo -e "\e[1mBuilding Linux...'\e[0m"
env GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o bin/ ./cmd/...

echo -e "\e[1mBuilding Windows...'\e[0m"
env GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o bin/ ./cmd/...

echo -e "\e[1mBuilding Solaris...'\e[0m"
env GOOS=solaris GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o bin/collector-solaris ./cmd/collector/
env GOOS=solaris GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o bin/receiver-solaris ./cmd/receiver/

echo -e "\e[1mDone.\n\e[0m"
ls -hl bin/

echo
file bin/*