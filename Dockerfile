FROM golang:1.22.2-bullseye

WORKDIR /app

RUN apt update && apt install sshpass

COPY cmd ./cmd/

COPY go.mod ./

COPY go.sum ./

COPY pkg ./pkg/

RUN cd ./cmd && go build -o main

COPY config.yaml ./

ENTRYPOINT [ "./cmd/main" ]