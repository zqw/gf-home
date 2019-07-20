#!/usr/bin/env bash

#go mod vendor

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o main main.go