# syntax = docker/dockerfile:experimental
FROM rust:1.45-stretch AS builder


RUN apt-get update && \
    apt-get install -y libssl-dev pkg-config musl musl-dev musl-tools

RUN ln -s /usr/include/x86_64-linux-gnu/asm /usr/include/x86_64-linux-musl/asm &&\
    ln -s /usr/include/asm-generic /usr/include/x86_64-linux-musl/asm-generic &&\
    ln -s /usr/include/linux /usr/include/x86_64-linux-musl/linux

RUN mkdir /musl

WORKDIR /root
RUN wget  https://github.com/openssl/openssl/archive/OpenSSL_1_1_1f.tar.gz
RUN tar zxvf OpenSSL_1_1_1f.tar.gz

WORKDIR /root/openssl-OpenSSL_1_1_1f/

RUN env CC="musl-gcc -fPIE -pie" ./Configure no-shared no-async --prefix=/musl --openssldir=/musl/ssl linux-x86_64
RUN make depend
RUN make -j$(nproc)
RUN make install
ENV PKG_CONFIG_ALLOW_CROSS 1
ENV OPENSSL_STATIC true
ENV OPENSSL_DIR /musl


WORKDIR /services/title-fetcher

RUN mkdir -p pb
COPY src/ src
COPY Cargo.toml .
COPY Cargo.lock .
COPY build.rs .
COPY ../../pb/title_fetcher.proto ./pb

RUN rustup target add x86_64-unknown-linux-musl
RUN rustup component add rustfmt
RUN --mount=type=cache,target=/services/title_fetcher/target cargo build --release --target=x86_64-unknown-linux-musl
RUN strip target/x86_64-unknown-linux-musl/release/title-fetcher

FROM alpine

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

COPY --from=builder /services/title-fetcher/target/x86_64-unknown-linux-musl/release/title-fetcher /services/title-fetcher/title-fetcher

RUN adduser -D -u 1000 app
USER 1000

ENTRYPOINT ["/services/title-fetcher/title-fetcher"]
