svc_name := emojivoto-web

all: clean protoc dep package-web test package 

include ../common.mk

package-web:
	cp web/index_bundle.js $(target_dir)
