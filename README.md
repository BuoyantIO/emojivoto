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

    ```bash
    curl https://run.linkerd.io/install | sh
    ```

1. Install Linkerd2

    ```bash
    linkerd install | kubectl apply -f -
    ```

1. View the dashboard!

    ```bash
    linkerd dashboard
    ```

1. Inject, Deploy, and Enjoy

    ```bash
    kubectl kustomize kustomize/deployment | \
        linkerd inject - | \
        kubectl apply -f -
    ```

1. Use the app!

    ```bash
    minikube -n emojivoto service web-svc
    ```

### In docker-compose

It's also possible to run the app with docker-compose (without Linkerd2).

Build and run:

```bash
make deploy-to-docker-compose
```

The web app will be running on port 8080 of your docker host.

### Via URL

To deploy standalone to an existing cluster:

```bash
kubectl apply -k github.com/BuoyantIO/emojivoto/kustomize/deployment
```

### Generating some traffic

The `VoteBot` service can generate some traffic for you. It votes on emoji
"randomly" as follows:

- It votes for :doughnut: 15% of the time.
- When not voting for :doughnut:, it picks an emoji at random

If you're running the app using the instructions above, the VoteBot will have
been deployed and will start sending traffic to the vote endpoint.

If you'd like to run the bot manually:

```bash
export WEB_HOST=localhost:8080 # replace with your web location
go run emojivoto-web/cmd/vote-bot/main.go
```

## Releasing a new version

To build and push multi-arch docker images:

1. Update the tag name in `common.mk`
1. Create the Buildx builder instance

    ```bash
    docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
    docker buildx create --name=multiarch-builder --driver=docker-container --use
    docker buildx inspect multiarch-builder --bootstrap
    ```

1. Build & push the multi-arch docker images to hub.docker.com

    ```bash
    docker login
    make multi-arch
    ```

1. Update:
    - `docker-compose.yml`
    - `kustomize/deployment/emoji.yml`
    - `kustomize/deployment/vote-bot.yml`
    - `kustomize/deployment/voting.yml`
    - `kustomize/deployment/web.yml`

1. Distribute to the Linkerd website repo

    ```bash
    kubectl kustomize kustomize/deployment  > ../website/run.linkerd.io/public/emojivoto.yml
    kubectl kustomize kustomize/daemonset   > ../website/run.linkerd.io/public/emojivoto-daemonset.yml
    kubectl kustomize kustomize/statefulset > ../website/run.linkerd.io/public/emojivoto-statefulset.yml
    ```

## Prometheus Metrics

By default the voting service exposes Prometheus metrics about current vote count on port `8801`.

This can be disabled by unsetting the `PROM_PORT` environment variable.

## Local Development

### Emojivoto webapp

This app is written with React and bundled with webpack.
Use the following to run the emojivoto go services and develop on the frontend.

Set up proto files, build apps

```bash
make build
```

Start the voting service

```bash
GRPC_PORT=8081 go run emojivoto-voting-svc/cmd/server.go
```

[In a separate terminal window] Start the emoji service

```bash
GRPC_PORT=8082 go run emojivoto-emoji-svc/cmd/server.go
```

[In a separate terminal window] Bundle the frontend assets

```bash
cd emojivoto-web/webapp
yarn install
yarn webpack # one time asset-bundling OR
yarn webpack-dev-server --port 8083 # bundle/serve reloading assets
```

[In a separate terminal window] Start the web service

```bash
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

```bash
export WEB_HOST=localhost:8080
go run emojivoto-web/cmd/vote-bot/main.go
```

View emojivoto

```bash
open http://localhost:8080
```

### Testing Linkerd Service Profiles

[Service Profiles](https://linkerd.io/2/features/service-profiles/) are a
feature of Linkerd that provide per-route functionality such as telemetry,
timeouts, and retries. The Emojivoto application is designed to showcase
Service Profiles by following the instructions below.

#### Generate the ServiceProfile definitions from the `.proto` files

The `emoji` and `voting` services are [gRPC](https://grpc.io/) applications
which have [Protocol Buffers (protobuf)](https://developers.google.com/protocol-buffers)
definition files. These `.proto` files can be used as input to the `linkerd
profile` command in order to create the `ServiceProfile` definition yaml files.
The [Linkerd Service Profile documentation](https://linkerd.io/2/tasks/setting-up-service-profiles/#protobuf)
outlines the steps necessary to create the yaml files, and these are the
commands you can use from the root of this repository:

```
linkerd profile --proto proto/Emoji.proto emoji-svc -n emojivoto
```
```
linkerd profile --proto proto/Voting.proto voting-svc -n emojivoto
```

Each of these commands will output yaml that you can write to a file or pipe
directly to `kubectl apply`. For example:

- To write to a file:
```
linkerd profile --proto proto/Emoji.proto emoji-svc -n emojivoto > emoji
-sp.yaml
```

- To apply directly:
```
linkerd profile --proto proto/Voting.proto voting-svc -n emojivoto | \
kubectl apply -f -
```

#### Generate the ServiceProfile definition for the Web deployment

The `web-svc` deployment of emojivoto is a React application that is hosted by a
Go server. We can use [`linkerd profile auto creation`](https://linkerd.io/2/tasks/setting-up-service-profiles/#auto-creation)
to generate the `ServiceProfile` resource for the web-svc with this command:

```bash
linkerd profile -n emojivoto web-svc --tap deploy/web --tap-duration 10s | \
   kubectl apply -f -
```

Now that the service profiles are generated for all the services, you can
observe the per-route metrics for each service on the [Linkerd Dashboard](https://linkerd.io/2/features/dashboard/)
or with the `linkerd routes` command

```bash
linkerd -n emojivoto routes deploy/web-svc --to svc/emoji-svc
```
## License

Copyright 2020 Buoyant, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
these files except in compliance with the License. You may obtain a copy of the
License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
