FROM golang:1.22.2-bullseye

WORKDIR /app

COPY . .

RUN cd ./cmd && go build -o main

ENTRYPOINT [ "./cmd/main" ]