svc_name := emojivoto-web-app
.PHONY: webpack

all: clean webpack package-ui
include ../common.mk
local: webpack

webpack:
	yarn install
	yarn webpack