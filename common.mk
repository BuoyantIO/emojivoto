export IMAGE_TAG := v11.5

.PHONY: package protoc test

target_dir := target

clean:
	rm -rf gen
	rm -rf $(target_dir)
	mkdir -p $(target_dir)
	mkdir -p gen

protoc: 
	protoc -I .. ../proto/*.proto --go_out=paths=source_relative:./gen --go-grpc_out=paths=source_relative:./gen

package: protoc compile build-container

package-ui: build-container-ui

build-container:
	docker build .. --platform linux/amd64 -t "datawire/$(svc_name):$(IMAGE_TAG)" --build-arg svc_name=$(svc_name)

build-multi-arch:
	docker buildx build .. -t "datawire/$(svc_name):$(IMAGE_TAG)" --build-arg svc_name=$(svc_name) \
		-f ../Dockerfile-multi-arch --platform linux/amd64,linux/arm64 --push

build-container-ui:
	docker build .. -t "datawire/emojivoto-web-app:$(IMAGE_TAG)" --build-arg svc_name=emojivoto-web-app -f ../Dockerfile-ui

compile:
	GOOS=linux GOARCH=amd64 go build -v -o $(target_dir)/$(svc_name) cmd/server.go

test:
	go test ./...

run:
	go run cmd/server.go