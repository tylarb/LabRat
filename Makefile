all: test cli slackengine

CLI=labrat
SLACK_ENGINE=labrat-slack
BUILD_DIR?=$(CURDIR)/bin
CNTR_BUILD_DIR=/out/bin
GOVER=1.13
GOFLAGS=-mod=vendor


test:
	podman run -t --rm -v $(CURDIR):/labrat --workdir /labrat -e GOFLAGS=$(GOFLAGS) golang:$(GOVER) make test-local


cli:
	podman run -t --rm -v $(CURDIR):/labrat -v $(BUILD_DIR):/out --workdir /labrat -e GOFLAGS=$(GOFLAGS) golang:$(GOVER) make -e BUILD_DIR=/out cli-local


slackengine:
	podman run -t --rm -v $(CURDIR):/labrat -v $(BUILD_DIR):/out --workdir /labrat -e GOFLAGS=$(GOFLAGS) golang:$(GOVER) make -e BUILD_DIR=/out slackengine-local


local: test-local cli-local slackengine-local

test-local:
	go test ./pkg/labrat
	go test ./cmd/labrat
	go test ./cmd/labrat-slack

cli-local:
	go build -o $(BUILD_DIR)/$(CLI) ./cmd/$(CLI)

slackengine-local:
	go build -o $(BUILD_DIR)/$(SLACK_ENGINE) ./cmd/$(SLACK_ENGINE)
