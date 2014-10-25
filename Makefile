# Cross platorm build with docker

# ----------------------------------------------------------------------------

ROOT_DIR        := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_PROG      := $(ROOT_DIR)/ah
CROSS_BUILD_DIR := $(ROOT_DIR)/build
GOLANG_AH       := github.com/9seconds/ah

LINUX_ARCH      := amd64 386 arm
DARWIN_ARCH     := amd64 386
FREEBSD_ARCH    := amd64 386 arm

DOCKER_PROG     := docker
DOCKER_GOPATH   := /go
DOCKER_WORKDIR  := $(DOCKER_GOPATH)/src/$(GOLANG_AH)
DOCKER_IMAGE    := golang:1.3.3-cross

# ----------------------------------------------------------------------------

define crosscompile
	GOOS=$(1) GOARCH=$(2) go build -o $(CROSS_BUILD_DIR)/$(1)-$(2) $(GOLANG_AH)
endef

# ----------------------------------------------------------------------------

all: fix vet lint prog-build
cross: cross-linux cross-darwin cross-freebsd
clean: prog-clean cross-clean

# ----------------------------------------------------------------------------

fix:
	go fix $(GOLANG_AH)/...

vet:
	go vet $(GOLANG_AH)/...

lint:
	golint $(GOLANG_AH)/...

fmt:
	go fmt $(GOLANG_AH)/...

godep:
	go get github.com/tools/godep

save: godep
	godep save

restore: godep
	godep restore

prog-build: restore prog-clean
	go build -o $(BUILD_PROG) $(GOLANG_AH)

install: restore
	go install $(GOLANG_AH)

prog-clean:
	rm -f $(BUILD_PROG)

# ----------------------------------------------------------------------------

cross-linux: $(addprefix cross-linux-,$(LINUX_ARCH))
cross-freebsd: $(addprefix cross-freebsd-,$(FREEBSD_ARCH))
cross-darwin: $(addprefix cross-darwin-,$(DARWIN_ARCH))

cross-clean:
	rm -rf $(CROSS_BUILD_DIR)

cross-build-directory: cross-clean
	mkdir -p $(CROSS_BUILD_DIR)

cross-linux-%: restore cross-build-directory
	$(call crosscompile,linux,$*)

cross-darwin-%: restore cross-build-directory
	$(call crosscompile,darwin,$*)

cross-freebsd-%: restore cross-build-directory
	$(call crosscompile,freebsd,$*)

cross-docker:
	$(DOCKER_PROG) run --rm -i -t -v "$(ROOT_DIR)":$(DOCKER_WORKDIR) -w $(DOCKER_WORKDIR) $(DOCKER_IMAGE) \
	bash -i -c "make -j 4 cross"
