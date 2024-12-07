# stage - base image
#
FROM golang:1.23.3 AS base

RUN apt-get -y install git make wget ca-certificates tzdata


# stage - golang build image
#
FROM base AS golang

ENV GO111MODULE=on

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/github.com/x-challenges/sparrow/

COPY go.* ./

RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download -x

COPY . ./


# stage - build app
#
FROM golang AS build

RUN --mount=target=. \
	--mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	go build -v -a -o /tmp/server


# stage - final, build app image
#
FROM alpine:latest

RUN apk add --update --no-cache ca-certificates tzdata && \
	rm -rf /var/cache/apk/*

WORKDIR /app

COPY configs /app/configs

COPY --from=build /tmp/server ./

RUN chmod a+x /app/server

CMD ["/app/server"]
