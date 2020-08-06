# syntax = docker/dockerfile:experimental
FROM golang:1.14-alpine AS builder

RUN apk --update add make

WORKDIR /services/blog
COPY go.mod go.sum ./
RUN go mod download
COPY Makefile ./
RUN make setup
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build make build


FROM alpine

COPY --from=builder /services/blog/bin/server /services/blog/bin/server

RUN adduser -D -u 1000 app
USER 1000

ENTRYPOINT ["/services/blog/bin/server"]
