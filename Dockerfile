ARG BUILDX_VERSION=0.14.0
ARG DOCKER_VERSION=26.0.2-dind
ARG GOLANG_VERSION=1.22

FROM --platform=$BUILDPLATFORM golang:${GOLANG_VERSION} as build

COPY . /src
WORKDIR /src

ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    make build

FROM docker/buildx-bin:${BUILDX_VERSION} as buildx-bin
FROM docker:${DOCKER_VERSION}

RUN apk --update --no-cache add coredns git

COPY --from=build /src/Corefile /etc/coredns/Corefile
COPY --from=buildx-bin /buildx /usr/libexec/docker/cli-plugins/docker-buildx
COPY --from=build /src/plugin-docker-buildx /bin/plugin-docker-buildx

ENV DOCKER_HOST=unix:///var/run/docker.sock

ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "plugin-docker-buildx"]
