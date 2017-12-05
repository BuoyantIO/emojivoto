svc_name := emojivoto-web

all: clean protoc dep webpack package-web test package

include ../common.mk

webpack:
	$(MAKE) -C webapp

package-web: webpack
	mkdir -p $(target_dir)/web
	cp web/favicon.ico $(target_dir)/web
	cp -a webapp/dist $(target_dir)
