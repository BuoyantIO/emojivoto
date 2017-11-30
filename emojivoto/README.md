# Emoji.voto

A demo app for the Conduit service mesh

## Demo Instructions (Minikube)

0. Pull images and install CLI (pre-release)

In the boron repo:

```
gcloud docker --authorize-only
bin/mkube bin/docker-pull latest
go install ./conduit
```

1. Install Conduit

```
conduit install | kubectl apply -f -
```

2. Build Votemoji images

```
eval $(minikube docker-env)
make build-base-docker-image build
```

3. Inject, Deploy, and Enjoy

```
conduit inject emojivoto.yml --skip-inbound-ports=80 | kubectl apply -f -
```

4. Use the app!

```
minikube -n emojivoto service web-svc
```

5. View the dashboard!

```
conduit dashboard
```

## Docker Instructions

To run the app locally with docker-compose:

```
make deploy-to-docker-compose
```

The web app will be running on port 8080 of your docker host.
