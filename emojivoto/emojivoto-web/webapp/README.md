# Emojivoto webapp

This app is written with React and bundled with webpack.

## Local development

Use the following to run the emojivoto go services and develop on the frontend.

Start the emoji service
```
cd ../../emojivoto-emoji-svc
export GRPC_PORT=8082
go run cmd/server.go
```

Start the voting service
```
cd emojivoto-voting-svc
export GRPC_PORT=8081
go run cmd/server.go
```

Bundle the frontend assets and start the web service
```
yarn install
yarn webpack # one time asset-bundling OR
yarn webpack-dev-server # bundle/serve reloading assets

cd ..
export WEB_PORT=8080
export EMOJISVC_HOST=localhost:8082
export VOTINGSVC_HOST=localhost:8081
export INDEX_BUNDLE=webapp/dist/index_bundle.js

# if you want reloading assets
# export WEBPACK_DEV_SERVER=http://localhost:8083
go run cmd/server.go
```

[Optional] Start the vote bot for automatic traffic generation.
```
cd ../../emojivoto-web/cmd/vote-bot
export WEB_HOST=localhost:8080
go run main.go
```
