FROM golang:1.17-alpine as build

COPY . /src
WORKDIR /src

RUN go build -v -a -tags netgo ./cmd/docker-buildx -o docker-buildx

FROM docker:20.10-dind

ARG BUILDX_VERSION

# renovate: datasource=github-releases depName=docker/buildx
ENV BUILDX_VERSION="${BUILDX_VERSION:-v0.6.3}"

ENV DOCKER_HOST=unix:///var/run/docker.sock

RUN \
    apk --update add --virtual .build-deps curl && \
    mkdir -p /usr/lib/docker/cli-plugins/ && \
    curl -SsL -o /usr/lib/docker/cli-plugins/docker-buildx "https://github.com/docker/buildx/releases/download/v${BUILDX_VERSION##v}/buildx-v${BUILDX_VERSION##v}.linux-amd64" && \
    chmod 755 /usr/lib/docker/cli-plugins/docker-buildx && \
    apk del .build-deps && \
    rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

COPY --from=build /src/docker-buildx /bin/docker-buildx

ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "docker-buildx"]
