include ./common.mk

.PHONY: web emoji-svc voting-svc integration-tests push

all: build integration-tests

build-base-docker-image:
	docker build . -f Dockerfile-base -t "datawire/emojivoto-svc-base:$(IMAGE_TAG)"

web:
	$(MAKE) -C emojivoto-web

web-app:
	$(MAKE) -C emojivoto-web-app

emoji-svc:
	$(MAKE) -C emojivoto-emoji-svc

voting-svc:
	$(MAKE) -C emojivoto-voting-svc

emoji-svc-local:
	$(MAKE) -C emojivoto-emoji-svc local

voting-svc-local:
	$(MAKE) -C emojivoto-voting-svc local

web-local:
	$(MAKE) -C emojivoto-web local

web-app-local:
	$(MAKE) -C emojivoto-web-app local

build: web web-app emoji-svc voting-svc

build-ui: build-container-ui

local: web-local web-app-local emoji-svc-local voting-svc-local

multi-arch:
	$(MAKE) -C emojivoto-web build-multi-arch
	$(MAKE) -C emojivoto-emoji-svc build-multi-arch
	$(MAKE) -C emojivoto-voting-svc build-multi-arch

deploy-to-minikube:
	$(MAKE) -C emojivoto-web build-container
	$(MAKE) -C emojivoto-emoji-svc build-container
	$(MAKE) -C emojivoto-voting-svc build-container
	kubectl delete -f emojivoto.yml || echo "ok"
	kubectl apply -f emojivoto.yml

deploy-to-docker-compose:
	docker-compose stop
	docker-compose rm -vf
	$(MAKE) -C emojivoto-web build-container
	$(MAKE) -C emojivoto-emoji-svc build-container
	$(MAKE) -C emojivoto-voting-svc build-container
	docker-compose -f ./docker-compose.yml up -d

push-%:
	docker push datawire/emojivoto-$*:$(IMAGE_TAG)

push: push-svc-base push-emoji-svc push-voting-svc push-web push-web-app
