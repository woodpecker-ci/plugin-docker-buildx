TARGETOS ?= linux
TARGETARCH ?= amd64
LDFLAGS := -s -w -extldflags "-static"

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags '${LDFLAGS}' -v -a -tags netgo -o plugin-docker-buildx ./cmd/docker-buildx

format: install-tools
	gofumpt -extra -w .

formatcheck: install-tools
	@([ -z "$(shell gofumpt -d . | head)" ]) || (echo "Source is unformatted"; exit 1)

install-tools: ## Install development tools
	@hash gofumpt > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go install mvdan.cc/gofumpt@latest; \
	fi
