#!/bin/bash

# go test ./... --cover
go test ./... --cover -v
# go test ./... --cover --bench=.
# go test ./... --cover -v --bench=.