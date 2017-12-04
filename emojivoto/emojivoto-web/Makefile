svc_name := emojivoto-web

all: clean protoc dep webpack package-web test package

include ../common.mk

webpack:
	$(MAKE) -C webapp

package-web: webpack
	cp -a webapp/dist $(target_dir)
