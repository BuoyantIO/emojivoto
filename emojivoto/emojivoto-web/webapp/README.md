# Votemoji web app

This webapp is built using React, webpack and yarn.

## Running
```
cd demos/emojivoto/v1-monolith/app/votemoji
# Install dependencies, build assets
yarn && yarn webpack

cd /demos/emojivoto/v1-monolith/app
# Run webserver
WEB_PORT=8080 API_PORT=9090 go run cmd/server.go
open http://localhost:9090
```

You can also run `yarn webpack-dev-server` to rebuild assets on file change for development.
You'll need to point the webserver to this address, which you can do by using this in the webserver:

```
<script type="text/javascript" src="http://localhost:8080/index_bundle.js" async></script>
```
(TODO: make this configurable via a flag)