#!/bin/bash

env GOOS=linux GOARCH=386 go build -o release/historykv-$GOOS-$GOARCH
env GOOS=linux GOARCH=amd64 go build -o release/historykv-$GOOS-$GOARCH
env GOOS=darwin GOARCH=386 go build -o release/historykv-$GOOS-$GOARCH
env GOOS=darwin GOARCH=amd64 go build -o release/historykv-$GOOS-$GOARCH
env GOOS=openbsd GOARCH=386 go build -o release/historykv-$GOOS-$GOARCH
env GOOS=openbsd GOARCH=amd64 go build -o release/historykv-$GOOS-$GOARCH
