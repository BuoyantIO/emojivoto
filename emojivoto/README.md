# Emoji.voto

A microservice application that allows users to vote for their favorite emoji,
and tracks votes received on a leaderboard. May the best emoji win.

The application is composed of the following 3 services:

* [emojivoto-web](emojivoto-web/): Web frontend and REST API
* [emojivoto-emoji-svc](emojivoto-emoji-svc/): gRPC API for finding and listing emoji
* [emojivoto-voting-svc](emojivoto-voting-svc/): gRPC API for voting and leaderboard

## Running

### In Minikube

Deploy the application to Minikube using the Conduit service mesh.

1. Install the Conduit CLI

```
curl https://run.conduit.io/install | sh
```

1. Install Conduit

```
conduit install | kubectl apply -f -
```

2. Inject, Deploy, and Enjoy

```
conduit inject emojivoto.yml --skip-inbound-ports=80 | kubectl apply -f -
```

3. Use the app!

```
minikube -n emojivoto service web-svc
```

4. View the dashboard!

```
conduit dashboard
```

### In docker-compose

It's also possible to run the app with docker-compose (without Conduit).

Build and run:

```
make deploy-to-docker-compose
```

The web app will be running on port 8080 of your docker host.
