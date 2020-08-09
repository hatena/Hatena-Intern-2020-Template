# syntax = docker/dockerfile:experimental
FROM node:12-alpine AS builder

WORKDIR /services/renderer-ts
COPY package.json yarn.lock ./
RUN --mount=type=cache,id=yarn,target=/usr/local/share/.cache/yarn/v6 \
  yarn install --frozen-lockfile
COPY . .
RUN yarn build


FROM node:12-alpine

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

WORKDIR /services/renderer-ts
COPY --from=builder /services/renderer-ts/dist dist
COPY package.json yarn.lock ./
RUN --mount=type=cache,id=yarn,target=/usr/local/share/.cache/yarn/v6 \
  yarn install --frozen-lockfile --production true
COPY pb pb

USER 1000

ENTRYPOINT ["node", "dist/index.js"]
