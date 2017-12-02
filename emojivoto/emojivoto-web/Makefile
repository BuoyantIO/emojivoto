svc_name := emojivoto-web

all: clean protoc dep webpack package-web test package

include ../common.mk

# TODO: figure out ./web/ vs. ./webapp/
# currently non-Docker local builds are broken because /dist only lives in
# ./webapp/ but the go code loads assets from ./web/

webpack:
	$(MAKE) -C webapp
	cp webapp/dist/index_bundle.js web/

package-web:
	cp web/index_bundle.js $(target_dir)
	cp -a webapp/dist $(target_dir)
