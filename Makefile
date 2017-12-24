.PHONY: all build build-arm patch unpatch clean test
VERSION := 0.2
COMMIT := $(shell git describe --always)
GOOS ?= darwin
GOARCH ?= amd64
GOPATH ?= $(HOME)/go/

all: clean build build-arm

build:
	@echo "Compiling source for $(GOOS) $(GOARCH)"
	@mkdir -p build
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -ldflags "-X main.version=$(VERSION) -X main.build=$(COMMIT)" -o build/lora-logger-$(GOOS)-$(GOARCH)$(BINEXT) main.go

build-arm:
	@echo "Compiling source for linux arm-5"
	@mkdir -p build
	@mv build .build
	@$(MAKE) patch
	@GOPATH=$(GOPATH) xgo -image=svenagn/multitech-libpcap -ldflags "-X main.version=$(VERSION) -X main.build=$(COMMIT)" -out .build/lora-logger --targets=linux/arm-5 .
	@$(MAKE) unpatch
	@mv .build build

patch:
	@-patch -p1 -N -f < .docker/gopacket_pcap.patch

unpatch:
	@-patch -p1 -R -f < .docker/gopacket_pcap.patch

clean:
	@echo "Cleaning up workspace"
	@$(MAKE) unpatch
	@rm -rf build
	@rm -rf .build
	@rm -rf lora.log
