include ./common.mk

.PHONY: web emoji-svc voting-svc integration-tests push

all: build integration-tests

build-base-docker-image:
	docker build . -f Dockerfile-base -t "buoyantio/emojivoto-svc-base:$(IMAGE_TAG)"

web:
	$(MAKE) -C emojivoto-web

emoji-svc:
	$(MAKE) -C emojivoto-emoji-svc

voting-svc:
	$(MAKE) -C emojivoto-voting-svc

build: web emoji-svc voting-svc

containers: build build-base-docker-image
	$(MAKE) -C emojivoto-web build-container
	$(MAKE) -C emojivoto-emoji-svc build-container
	$(MAKE) -C emojivoto-voting-svc build-container

deploy-to-minikube: containers
	kubectl delete -f emojivoto.yml || echo "ok"
	kubectl apply -f emojivoto.yml

stop-docker-compose:
	docker-compose stop
	docker-compose rm -vf

deploy-to-docker-compose: stop-docker-compose containers
	docker-compose -f ./docker-compose.yml up -d

integration-tests: deploy-to-docker-compose
	WEB_URL=http://localhost:8080 $(MAKE) -C integration_test test
	docker-compose stop
	docker-compose rm -vf

push-%:
	docker push buoyantio/emojivoto-$*:$(IMAGE_TAG)

push: push-svc-base push-emoji-svc push-voting-svc push-web
