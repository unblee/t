NAME       := t
VERSION    := v0.1.0
REVISION   := $(shell git rev-parse --short HEAD)
GO_VERSION := 1.7.4
LDFLAGS    := -w -s \
              -X main.Version=$(VERSION) \
              -X main.Revision=$(REVISION) \
              -X main.GoVersion=$(GO_VERSION) \
              -extldflags '-static'

DOCKER_CMD_OPT := -v $(PWD):/go/src/$(NAME) \
                  -w /go/src/$(NAME) \
                  -e NAME=$(NAME) \
                  -e VERSION=$(VERSION) \
                  -e LDFLAGS='$(LDFLAGS)' \
                  -e GITHUB_TOKEN='$(GITHUB_TOKEN)'
DOCKER_CMD     := docker run -it --rm $(DOCKER_CMD_OPT) golang:$(GO_VERSION)-alpine

.PHONY: deps
deps:
	glide install

.PHONY: build
build: deps
	go build -a -tags netgo -installsuffix netgo -ldflags "$(LDFLAGS)"

.PHONY: cross-build
cross-build: deps
	$(DOCKER_CMD) sh scripts/cross-build.sh

.PHONY: release
release: cross-build
	$(DOCKER_CMD) sh scripts/release.sh

.PHONY: clean
clean:
	rm -fr dest/
	rm -fr release/