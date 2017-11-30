.PHONY: api web emoji-svc voting-svc integration-tests

all: build integration-tests

build-base-docker-image:
	docker build . -f Dockerfile-base -t "buoyantio/emojivoto-svc-base:v1"

api:
	$(MAKE) -C emojivoto-api

web:
	$(MAKE) -C emojivoto-web

emoji-svc:
	$(MAKE) -C emojivoto-emoji-svc

voting-svc:
	$(MAKE) -C emojivoto-voting-svc


build: api web emoji-svc voting-svc

deploy-to-minikube:
	$(MAKE) -C emojivoto-api build-container
	$(MAKE) -C emojivoto-web build-container
	$(MAKE) -C emojivoto-emoji-svc build-container
	$(MAKE) -C emojivoto-voting-svc build-container
	kubectl delete -f emojivoto.yml || echo "ok"
	kubectl apply -f emojivoto.yml


deploy-to-docker-compose:
	docker-compose stop
	docker-compose rm -vf
	$(MAKE) -C emojivoto-api build-container
	$(MAKE) -C emojivoto-web build-container
	$(MAKE) -C emojivoto-emoji-svc build-container
	$(MAKE) -C emojivoto-voting-svc build-container
	docker-compose -f ./docker-compose.yml up -d

integration-tests: deploy-to-docker-compose
	WEB_URL=http://localhost:8080 API_URL=http://localhost:9090 $(MAKE) -C integration_test test
	docker-compose stop
	docker-compose rm -vf
