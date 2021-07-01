#!/bin/bash

echo -e "\e[1mRunning 'go fmt ./...'\e[0m"
go fmt ./...

echo -e "\e[1mRunning 'go vet ./...'\e[0m"
go vet ./...

echo -e "\e[1mRunning 'golint ./...'\e[0m"
golint ./...

echo -e "\e[1mRunning 'staticheck ./...'\e[0m"
staticcheck ./...

echo -e "\e[1mRunning 'go tool fix --diff .'\e[0m"
go tool fix --diff .

echo -e "\e[1mDone.\e[0m"