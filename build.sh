#!/bin/sh -u

env GOOS=windows GOARCH=amd64 go build -o ./bin/win64.exe
env GOOS=windows GOARCH=386 go build -o ./bin/win32.exe
env GOOS=darwin GOARCH=amd64 go build -o ./bin/mac64
