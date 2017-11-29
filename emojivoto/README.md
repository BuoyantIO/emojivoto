# Emoji.voto

A demo app for the Conduit service mesh

## Demo Instructions (pre-release)

1. Pull and install Conduit

```
gcloud docker --authorize-only
bin/mkube bin/docker-pull latest
bin/go-run ./conduit/main.go install | kubectl apply -f -
```

2. Build Votemoji images

```
cd demos/emojivoto
bin/mkube make build-base-docker-image build
```

3. Inject, Deploy, and Enjoy

```
bin/go-run ./conduit/main.go inject demos/emojivoto/emojivoto.yml --skip-inbound-ports=80 | kubectl apply -f -
```

4. Use the app!

```
minikube -n emojivoto service web-svc
```

5. View the dashboard!

```
bin/go-run ./conduit/main.go dashboard
```
