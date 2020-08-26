# syntax = docker/dockerfile:experimental
FROM rust:1.45-stretch AS builder

WORKDIR /services/title-fetcher

RUN apt-get update && \
    apt-get install -y libssl-dev pkg-config

RUN mkdir -p pb
COPY . .
COPY ../../pb/title_fetcher.proto ./pb

RUN rustup component add rustfmt
RUN --mount=type=cache,target=/services/title-fetcher/target cargo build --release
RUN strip target/release/title-fetcher

FROM debian:latest

RUN apt-get update && \
    apt-get install -y wget
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

COPY --from=builder /services/title-fetcher/target/release/title-fetcher /services/title-fetcher/title-fetcher

RUN adduser -u 1000 app
USER 1000

ENTRYPOINT ["/services/title-fetcher/title-fetcher"]
