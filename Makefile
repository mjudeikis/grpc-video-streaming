APP_REPO ?= quay.io/synpse/video-streaming
TAG_NAME = $(shell git rev-parse --short=7 HEAD)$(shell [[ $$(git status --porcelain) = "" ]] || echo -dirty)
OUTPUT_BIN_APP ?= release
LDFLAGS		+= -s -w
GOARCH ?= arm64

.PHONY: image
image:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t ${APP_REPO}:latest --push -f Dockerfile .

.PHONY: app-arm64, app-amd64
app-arm64:
	CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o ${OUTPUT_BIN_APP}/linux/arm64/app ./cmd/ui

app-amd64:
	CGO_ENABLED=1  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ${OUTPUT_BIN_APP}/linux/amd64/app ./cmd/ui

app: app-arm app-arm64 app-amd64