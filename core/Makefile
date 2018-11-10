GOTOOLS = github.com/golang/lint/golint \
          gopkg.in/alecthomas/gometalinter.v2

PKGS = $(shell go list ./...

TARGET_DIR = ../target
TARGET = iot-agent

all: clean lint test build

build: deps
	mkdir -p $(TARGET_DIR)
	@go build -o $(TARGET_DIR)/iot-agent main.go

clean:
	rm -rf $(TARGET_DIR)

deps:
	@go mod tidy

test: deps
	@go test -race ./...

tools:
	@go get $(GOTOOLS)
	@gometalinter.v2 --install > /dev/null

tools-update:
	@go get -u $(GOTOOLS)
	@gometalinter.v2 --install

lint: tools deps
	@gometalinter.v2 --config=.gometalinter.json $(PKGS)

lint-all: deps
	@gometalinter.v2 --config=.gometalinter.json --enable=interfacer --enable=gosimple $(PKGS)

.PHONY: all deps test lint lint-all tools tools-update build clean