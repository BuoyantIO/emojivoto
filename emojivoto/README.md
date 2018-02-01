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

2. Install Conduit

```
conduit install | kubectl apply -f -
```

3. View the dashboard!

```
conduit dashboard
```

4. Inject, Deploy, and Enjoy

```
conduit inject emojivoto.yml --skip-inbound-ports=80 | kubectl apply -f -
```

5. Use the app!

```
minikube -n emojivoto service web-svc
```

### In docker-compose

It's also possible to run the app with docker-compose (without Conduit).

Build and run:

```
make deploy-to-docker-compose
```

The web app will be running on port 8080 of your docker host.


### Generating some traffic

The `VoteBot` service can generate some traffic for you. It votes on emoji
"randomly" as follows:
- It votes for :doughnut: 15% of the time.
- It votes for :poop: 20% of the time.
- When not voting for :doughnut: or :poop:, it picks an emoji at random

If you're running the app using the instructions above,
the VoteBot will have been deployed and will
start sending traffic to the vote endpoint.

If you'd like to run the bot manually:
```
export WEB_HOST=localhost:8080 # replace with your web location
go run emojivoto-web/cmd/vote-bot/main.go
```
