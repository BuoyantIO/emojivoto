svc_name := emojivoto-web

all: clean protoc package-web test package

local: clean protoc package-web test

include ../common.mk

package-web:
	mkdir -p $(target_dir)/web
	cp web/favicon.ico $(target_dir)/web

compile-vote-bot:
	GOOS=linux go build -v -o $(target_dir)/emojivoto-vote-bot cmd/vote-bot/main.go

compile: compile-vote-bot

.PHONY: build-releases
build-releases: clean protoc
	mkdir -p bin/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "-s -w" -o bin/emojivoto-server-amd64 cmd/server.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -ldflags "-s -w" -o bin/emojivoto-server-darwin-amd64 cmd/server.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags "-s -w" -o bin/emojivoto-server-darwin-arm64 cmd/server.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -ldflags "-s -w" -o bin/emojivoto-server.exe cmd/server.go
