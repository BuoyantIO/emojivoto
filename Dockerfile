ARG PROJECT_ID

FROM gcr.io/$PROJECT_ID/emojivoto-base:v13 AS build
WORKDIR /app
COPY . .
RUN apt update && apt install -y unzip protoc-gen-go protoc-gen-go-grpc
RUN make build
RUN mkdir -p /target/usr/bin \
    && cp -p ./emojivoto-emoji-svc/target/emojivoto-emoji-svc /target/usr/bin/ \
    && cp -p ./emojivoto-voting-svc/target/emojivoto-voting-svc /target/usr/bin/ \
    && cp -p ./emojivoto-web/target/emojivoto-vote-bot /target/usr/bin/ \
    && cp -p ./emojivoto-web/target/emojivoto-web /target/usr/bin/

FROM debian:bookworm
COPY --from=build /target/ /
