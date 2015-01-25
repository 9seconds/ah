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
DOCKER_IMAGE    := golang:1.4.1-cross

# ----------------------------------------------------------------------------

define crosscompile
	GOOS=$(1) GOARCH=$(2) go build -a -o $(CROSS_BUILD_DIR)/$(1)-$(2) $(GOLANG_AH)
endef

# ----------------------------------------------------------------------------

all: tools prog-build
tools: fix vet lint
cross: cross-linux cross-darwin cross-freebsd
clean: prog-clean cross-clean
ci: tools cross

# ----------------------------------------------------------------------------

fix:
	go fix $(GOLANG_AH)/...

vet: govet
	go vet $(GOLANG_AH)/...

lint: golint
	golint $(GOLANG_AH)/...

fmt:
	go fmt $(GOLANG_AH)/...

godep:
	go get github.com/tools/godep || true

govet:
	go get golang.org/x/tools/cmd/vet || true

golint:
	go get github.com/golang/lint/golint || true

save: godep
	godep save

restore: godep
	godep restore

prog-build: restore prog-clean
	go build -a -o $(BUILD_PROG) $(GOLANG_AH)

install: restore
	go install -a $(GOLANG_AH)

prog-clean:
	rm -f $(BUILD_PROG)

update:
	cat $(ROOT_DIR)/Godeps/Godeps.json \
		| grep ImportPath \
		| grep -v $(GOLANG_AH) \
		| awk '{print $$2}' \
		| sed 's/"//g; s/,$$//' \
		| xargs -n 1 godep update

upgrade_deps:
	cat $(ROOT_DIR)/Godeps/Godeps.json \
		| grep ImportPath \
		| grep -v $(GOLANG_AH) \
		| awk '{print $$2}' \
		| sed 's/"//g; s/,$$//' \
		| xargs -n 1 -P 4 go get -u

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
	sudo chown -R $(USER):$(USER) $(CROSS_BUILD_DIR)
