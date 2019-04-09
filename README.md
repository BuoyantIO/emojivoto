# Emoji.voto

A microservice application that allows users to vote for their favorite emoji,
and tracks votes received on a leaderboard. May the best emoji win.

The application is composed of the following 3 services:

* [emojivoto-web](emojivoto-web/): Web frontend and REST API
* [emojivoto-emoji-svc](emojivoto-emoji-svc/): gRPC API for finding and listing emoji
* [emojivoto-voting-svc](emojivoto-voting-svc/): gRPC API for voting and leaderboard

![Emojivoto Topology](assets/emojivoto-topology.png "Emojivoto Topology")

## Running

### In Minikube

Deploy the application to Minikube using the Linkerd2 service mesh.

1. Install the `linkerd` CLI

```
curl https://run.linkerd.io/install | sh
```

2. Install Linkerd2

```
linkerd install | kubectl apply -f -
```

3. View the dashboard!

```
linkerd dashboard
```

4. Inject, Deploy, and Enjoy

```
linkerd inject emojivoto.yml | kubectl apply -f -
```

5. Use the app!

```
minikube -n emojivoto service web-svc
```

### In docker-compose

It's also possible to run the app with docker-compose (without Linkerd2).

Build and run:

```
make deploy-to-docker-compose
```

The web app will be running on port 8080 of your docker host.

### Generating some traffic

The `VoteBot` service can generate some traffic for you. It votes on emoji
"randomly" as follows:
- It votes for :doughnut: 15% of the time.
- When not voting for :doughnut:, it picks an emoji at random

If you're running the app using the instructions above, the VoteBot will have
been deployed and will start sending traffic to the vote endpoint.

If you'd like to run the bot manually:
```
export WEB_HOST=localhost:8080 # replace with your web location
go run emojivoto-web/cmd/vote-bot/main.go
```

## Releasing a new version

To update the docker images:
1. Update the tag name in `common.mk`
2. Update the base image tags in `Makefile` and `Dockerfile`
3. Build base docker image `make build-base-docker-image`
4. Build docker images `make build`
5. Push the docker images to hub.docker.com
```bash
docker login
docker push buoyantio/emojivoto-svc-base:v8
docker push buoyantio/emojivoto-emoji-svc:v8
docker push buoyantio/emojivoto-voting-svc:v8
docker push buoyantio/emojivoto-web:v8
```
6. Update `emojivoto.yml`, `docker-compose.yml`


## Local Development

### Emojivoto webapp

This app is written with React and bundled with webpack.
Use the following to run the emojivoto go services and develop on the frontend.

Set up proto files, build apps
```
make build
```

Start the voting service
```
GRPC_PORT=8081 go run emojivoto-voting-svc/cmd/server.go
```

[In a separate terminal window] Start the emoji service
```
GRPC_PORT=8082 go run emojivoto-emoji-svc/cmd/server.go
```

[In a separate terminal window] Bundle the frontend assets
```
cd emojivoto-web/webapp
yarn install
yarn webpack # one time asset-bundling OR
yarn webpack-dev-server --port 8083 # bundle/serve reloading assets
```

[In a separate terminal window] Start the web service
```
export WEB_PORT=8080
export VOTINGSVC_HOST=localhost:8081
export EMOJISVC_HOST=localhost:8082

# if you ran yarn webpack
export INDEX_BUNDLE=emojivoto-web/webapp/dist/index_bundle.js

# if you ran yarn webpack-dev-server
export WEBPACK_DEV_SERVER=http://localhost:8083

# start the webserver
go run emojivoto-web/cmd/server.go
```

[Optional] Start the vote bot for automatic traffic generation.
```
export WEB_HOST=localhost:8080
go run emojivoto-web/cmd/vote-bot/main.go
```

View emojivoto
```
open http://localhost:8080
```
