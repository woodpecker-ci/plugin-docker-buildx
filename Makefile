TARGETOS ?= linux
TARGETARCH ?= amd64
LDFLAGS := -s -w -extldflags "-static"

build:
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags '${LDFLAGS}' -v -a -tags netgo -o docker-buildx ./cmd/docker-buildx
