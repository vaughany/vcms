#!/bin/bash

echo -e "\e[1mRunning 'go fmt ./...'\e[0m"
go fmt ./...

echo -e "\e[1mRunning 'go vet ./...'\e[0m"
go vet ./...

echo -e "\e[1mRunning 'golint ./...'\e[0m"
golint --min_confidence 0.2 ./...

echo -e "\e[1mRunning 'staticheck ./...'\e[0m"
staticcheck ./...

echo -e "\e[1mRunning 'go tool fix --diff .'\e[0m"
go tool fix --diff .

echo -e "\e[1mRunning 'revive ./...'\e[0m"
revive -config revive.toml ./...

echo -e "\e[1mRunning 'golangci-lint run'\e[0m"
# golangci-lint run --enable-all
golangci-lint run

echo -e "\e[1mDone.\e[0m"
