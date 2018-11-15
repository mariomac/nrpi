COVFILE  ?= cov.out
GOTOOLS  ?= github.com/GeertJohan/fgt \
			golang.org/x/tools/cmd/goimports \
			github.com/kisielk/errcheck \
			honnef.co/go/tools/cmd/staticcheck
PACKAGES = $(shell go list ./...)
SOURCES  = $(shell find . -name "*.go" | grep -v vendor/)
TARGET_DIR ?= ../target
TARGET ?= iot-agent

dep:
	go mod tidy

build: dep
	go build -o $(TARGET_DIR)/$(TARGET)

clean:
	rm -rf $(TARGET_DIR)

test: dep
	go test -race $(PACKAGES)

cov-gen: dep
	go test  $(PACKAGES) -coverprofile $(COVFILE)

cov: cov-gen
	go tool cover -html $(COVFILE)

lint:
	go get $(GOTOOLS)
	fgt go fmt $(PACKAGES)
	fgt goimports -w $(SOURCES)
	echo $(PACKAGES) | xargs -L1 fgt golint
	fgt go vet $(PACKAGES)
	fgt errcheck -ignore Close $(PACKAGES)
	staticcheck $(PACKAGES)
.SILENT: lint
