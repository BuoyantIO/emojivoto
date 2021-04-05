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
