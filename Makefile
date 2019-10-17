all: test cli slackengine

CLI=labrat
SLACK_ENGINE=labrat-slack
BUILD_DIR=`pwd`/bin
CNTR_BUILD_DIR=/out/bin

build_image:
	podman build -t labrat_build .

test: build_image
	podman run -t --rm labrat_build go test ./cmd/$(CLI)
	podman run -t --rm labrat_build go test ./cmd/$(SLACK_ENGINE)
	podman run -t --rm labrat_build go test ./pkg/labrat


cli: build_image
	podman run -t --rm -v $(BUILD_DIR):$(CNTR_BUILD_DIR) labrat_build go build -o $(CNTR_BUILD_DIR)/$(CLI) ./cmd/labrat


slackengine: build_image
	podman run -t --rm -v $(BUILD_DIR):$(CNTR_BUILD_DIR) labrat_build go build -o $(CNTR_BUILD_DIR)/$(SLACK_ENGINE) ./cmd/labrat-slack


local: cli-local slackengine-local

cli-local:
	go build -o $(BUILD_DIR)/$(CLI) ./cmd/$(CLI)

slackengine-local:
	go build -o $(BUILD_DIR)/$(SLACK_ENGINE) ./cmd/$(SLACK_ENGINE)
