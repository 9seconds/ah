# Cross platorm build with docker

DOCKER_PROG := docker
DOCKER_IMAGE := golang:1.3.3-cross
ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_PROG := ah
BUILD_FILE := ah.go

define compile
	$(DOCKER_PROG) run --rm -i -t -v "$(ROOT_DIR)":/usr/src/app -w /usr/src/app \
		-e GOOS=$(1) -e GOARCH=$(2) $(DOCKER_IMAGE) \
		bash -c "go get -d -v; go build -v -o build/$(BUILD_PROG)-$(1)-$(2) $(BUILD_FILE)";
endef

all: linux darwin freebsd

build_directory: clean
	mkdir -p ./build

linux: linux-386 linux-amd64 linux-arm
darwin: darwin-386 darwin-amd64
freebsd: freebsd-386 freebsd-amd64 freebsd-arm

linux-386: build_directory
	$(call compile,linux,386)

linux-amd64: build_directory
	$(call compile,linux,amd64)

linux-arm: build_directory
	$(call compile,linux,arm)

darwin-386: build_directory
	$(call compile,darwin,386)

darwin-amd64: build_directory
	$(call compile,darwin,amd64)

freebsd-386: build_directory
	$(call compile,freebsd,386)

freebsd-amd64: build_directory
	$(call compile,freebsd,amd64)

freebsd-arm: build_directory
	$(call compile,freebsd,arm)

clean:
	rm -rf ./build
