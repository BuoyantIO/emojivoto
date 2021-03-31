svc_name := emojivoto-web

all: clean protoc webpack package-web test package

local: clean protoc webpack package-web test

include ../common.mk

webpack:
	$(MAKE) -C webapp

package-web: webpack
	mkdir -p $(target_dir)/web
	cp web/favicon.ico $(target_dir)/web
	cp -a webapp/dist $(target_dir)

compile-vote-bot:
	GOOS=linux go build -v -o $(target_dir)/emojivoto-vote-bot cmd/vote-bot/main.go

compile: compile-vote-bot
