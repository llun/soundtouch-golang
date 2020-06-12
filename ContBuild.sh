#!/bin/sh
CompileDaemon -exclude-dir=.git -exclude-dir=docs --build="go build ./server/cmd/soundtouch-r-e-s-tful-json-server-server/main.go" --command="./main -i en0 -n 7 -l info"
