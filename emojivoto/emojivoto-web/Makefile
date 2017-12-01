svc_name := emojivoto-web

all: clean protoc dep webpack package-web test package

include ../common.mk

webpack:
	$(MAKE) -C webapp
	cp webapp/dist/index_bundle.js web/

package-web:
	cp web/index_bundle.js $(target_dir)
