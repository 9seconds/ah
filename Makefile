# Cross platorm build with docker

.PHONY: clean build

# ----------------------------------------------------------------------------

BUILD_PROG := ah
ROOT_DIR   := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
GOLANG_AH  := github.com/9seconds/ah

# ----------------------------------------------------------------------------

all: clean fix vet lint build

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

build: restore
	go build -o $(BUILD_PROG) $(GOLANG_AH)

clean:
	rm -f $(BUILD_PROG)
