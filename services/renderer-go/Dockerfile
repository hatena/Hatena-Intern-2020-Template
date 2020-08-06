# syntax = docker/dockerfile:experimental
FROM golang:1.14-alpine AS builder

RUN apk --update add make

WORKDIR /services/renderer-go
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build make build


FROM alpine

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

COPY --from=builder /services/renderer-go/bin/server /services/renderer-go/bin/server

RUN adduser -D -u 1000 app
USER 1000

ENTRYPOINT ["/services/renderer-go/bin/server"]
