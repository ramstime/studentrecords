####### Makefile ###########
LINTERVERSION=v1.52.0
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
MODULE = $(shell go list -m)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif
GO      = go
BIN      = $(CURDIR)/bin
export GOBIN := $(GOBIN)
export PATH := $(CURDIR)/bin:$(PATH)
IMG ?= artifactory/studentrecords:latest

# place all tools in the pkg/tools/tools.go to ensure we version lock
$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) installing $(PACKAGE)…)
	GOBIN=$(BIN) $(GO) install -mod=vendor $(PACKAGE)
MODULE = $(shell go list -m)
MOCKERY = $(BIN)/mockery
$(BIN)/mockery: PACKAGE=github.com/vektra/mockery/cmd/mockery

# Interfaces for mock code generation using the mockery framwork.
# Each entry should be of the format <path>/<interface-name>.mock
# The mock files are output in the ./mocks folder.
MOCKSRCS = \
	controllers/ \
	controllers/books.mock \

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

#all: deps fmt build lint
all: deps fmt build

deps:
# go clean --modcache
#go mod edit -go=$(GOMODVERSION)
	go mod tidy
	GOPRIVATE= go mod vendor

##@ Development
# lint checks
GOLANGCI-LINT = $(BIN)/golangci-lint
$(BIN)/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(LINTERVERSION)

GOIMPORTS =$(BIN)/goimports
$(BIN)/goimports: PACKAGE=golang.org/x/tools/cmd/goimports

format: | $(GOIMPORTS); $(info $(M) Runing goimports formatter on all files except vendor)
	@$(GOIMPORTS) -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

lint: $(BIN)/golangci-lint
	$(GOLANGCI-LINT) run ./...

#
# Targets for mock file generation. The package targets are included from mw-common
#
.PHONY: mockery mockgen
mockgen: mockery ${MOCKSRCS}

%.mock:
	@echo "Generating mock for Interface=$(basename $(@F))"
	mockery --dir $(@D) --name $(basename $(@F)) --output $(@D)/mocks

# Run go generate
go-generate: $(MOCKERY) ; $(info $(M) running go generate…) @ ## Run go generate on all source files
	$Q $(GO) generate -mod=vendor ./...

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

test: lint 
	go test ./... -coverprofile cover.out fmt
	go tool cover -html=cover.out -o cover.html

##@ Build

build: deps | $(BIN) ; $(info $(M) building executable…) @ ## Build program binary
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build \
		-mod=vendor \
		-o $(BIN)/studentrecords main.go
# mkdir -p studentrecords/bin/

run: fmt vet ## Run a controller from your host.
	go run main.go

IMGTAG ?= latest
export IMGTAG
#
# Build the docker image
# To override the IMGTAG, use: IMGTAG=<value> make <target>
#

docker-build:
	docker build . --build-arg GOVERSION=$(GOVERSION) -t ${IMG}

#
# Push the docker image to the MW artifactory (hub)
#

docker-push:
	docker push ${IMG}
##@ Deployment

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
