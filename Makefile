.PHONY: build clean
#PKGS := $(shell go list ./... | grep -v /vendor/)
VERSION := 0.1
COMMIT := unknown #VERSION := $(shell git describe --always)
GOOS ?= darwin
GOARCH ?= amd64
GOPATH ?= /Users/sven/go/

build:
	@echo "Compiling source for $(GOOS) $(GOARCH)"
	@mkdir -p build
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -ldflags "-X main.version=$(VERSION) -X main.build=$(COMMIT)" -o build/lora-logger-$(GOOS)-$(GOARCH)$(BINEXT) main.go

build-arm:
	@echo "Compiling source for linux arm-5"
	@mkdir -p build
	cd .docker/; ./pre-build.sh
	mv build .build
	@GOPATH=$(GOPATH) xgo -image=svenagn/multitech-libpcap -out .build/lora-logger --targets=linux/arm-5 .
	mv .build build
	cd .docker/; ./post-build.sh

clean:
	@echo "Cleaning up workspace"
	@rm -rf build
	@rm -rf lora.log
