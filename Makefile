LOCAL_BIN := $(CURDIR)/bin

ifndef BIN_DIR
BIN_DIR = $(LOCAL_BIN)
endif

BUILD_DIR := $(CURDIR)/cmd/ltotm
MAIN_BUILD_DIR := $(CURDIR)/cmd/ltotm
MINIFY_BUILD_DIR := $(CURDIR)/cmd/minify
ASSETS_DIR := $(CURDIR)/internal/webserver/assets

GO_MAJOR_VERSION := $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION := $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
GO_PATCH_VERSION := $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f3)

GO_VERSION_SHORT := $(GO_MAJOR_VERSION).$(GO_MINOR_VERSION).$(GO_PATCH_VERSION)

GO_MIN_SUPPORTED_MAJOR_VERSION := 1
GO_MIN_SUPPORTED_MINOR_VERSION := 20

APP_VERSION := $(if $(CI_COMMIT_REF_SLUG),$(CI_COMMIT_REF_SLUG),$(if $(GIT_BRANCH),$(GIT_BRANCH),$(GIT_HASH)))

CI_PROJECT_ID ?= ltotm

GIT_HASH := $(shell git log --format="%h" --max-count 1 2> /dev/null)

GIT_LOG := $(shell git log --decorate --oneline --max-count 1 2> /dev/null | base64 | tr -d '\n')

GIT_BRANCH := $(shell git branch --show-current)

BUILD_TS := $(shell date +%FT%T%z)

LDFLAGS = \
    -X 'github.com/mrlinqu/ltotm/internal/config.ProjectID=$(CI_PROJECT_ID)' \
    -X 'github.com/mrlinqu/ltotm/internal/config.Version=$(APP_VERSION)' \
    -X 'github.com/mrlinqu/ltotm/internal/config.GoVersion=$(GO_VERSION_SHORT)' \
    -X 'github.com/mrlinqu/ltotm/internal/config.BuildDate=$(BUILD_TS)' \
    -X 'github.com/mrlinqu/ltotm/internal/config.GitLog=$(GIT_LOG)' \
    -X 'github.com/mrlinqu/ltotm/internal/config.GitHash=$(GIT_HASH)' \
    -X 'github.com/mrlinqu/ltotm/internal/config.GitBranch=$(GIT_BRANCH)'

.build: generate $(BUILD_DIR)
$(BUILD_DIR):
	$(eval APP_NAME = $(shell basename $@))
	$(eval LDFLAGS += -X 'github.com/mrlinqu/ltotm/internal/config.Name=$(APP_NAME)')
	$(BUILD_ENVPARMS) go build -o="$(BIN_DIR)/$(APP_NAME)" -ldflags "$(LDFLAGS)" $@

build: .build

.run:
	go run \
		-ldflags "$(LDFLAGS)" \
		$(MAIN_BUILD_DIR)

run: .run

GOLANGCI_BIN := $(LOCAL_BIN)/golangci-lint

GOLANGCI_TAG ?= 1.55.2

ifneq ($(shell command -v $(GOLANGCI_BIN)),)
GOLANGCI_BIN_VERSION := $(shell $(GOLANGCI_BIN) --version 2> /dev/null)

ifneq ($(GOLANGCI_BIN_VERSION),)
GOLANGCI_BIN_VERSION_SHORT := $(shell echo "$(GOLANGCI_BIN_VERSION)" | sed -En 's/.*([0-9]+\.[0-9]+\.[0-9]+) built.*/\1/p' | sed -E 's/v(.*)/\1/g')
else
GOLANGCI_BIN_VERSION_SHORT := 0
endif

ifneq ("$(GOLANGCI_TAG)", "$(word 1, $(sort $(GOLANGCI_TAG) $(GOLANGCI_BIN_VERSION_SHORT)))")
GOLANGCI_BIN :=
endif

endif

ifneq ($(shell command -v golangci-lint 2> /dev/null),)
GOLANGCI_VERSION := $(shell golangci-lint --version 2> /dev/null)

ifneq ($(GOLANGCI_VERSION),)
GOLANGCI_VERSION_SHORT := $(shell echo "$(GOLANGCI_VERSION)" | sed -En 's/.*([0-9]+\.[0-9]+\.[0-9]+) built.*/\1/p' | sed -E 's/v(.*)/\1/g')
else
GOLANGCI_VERSION_SHORT := 0
endif

ifeq ("$(GOLANGCI_TAG)", "$(word 1, $(sort $(GOLANGCI_TAG) $(GOLANGCI_VERSION_SHORT)))")
GOLANGCI_BIN := $(shell command -v golangci-lint 2> /dev/null)
endif

endif

install-lint: export GOBIN := $(LOCAL_BIN)
install-lint:
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info Installing golangci-lint v$(GOLANGCI_TAG))
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_TAG)
GOLANGCI_BIN := $(LOCAL_BIN)/golangci-lint
else
	$(info Golangci-lint is already installed to $(GOLANGCI_BIN))
endif

.lint: install-lint
	$(info Running lint against changed files...)
	$(GOLANGCI_BIN) run \
		--new-from-rev=origin/master \
		--config=.golangci.pipeline.yaml \
		./...

lint: .lint

.lint-full: install-lint
	$(info Running lint against all project files...)
	$(GOLANGCI_BIN) run \
		--config=.golangci.pipeline.yaml \
		./...

lint-full: .lint-full

.deps:
	$(info Installing dependencies...)
	go mod download

deps: .deps

generate:
	go run $(MINIFY_BUILD_DIR) $(ASSETS_DIR)

.PHONY: \
	all \
	.install-lint \
	install-lint \
	.lint \
	lint \
	.lint-full \
	lint-full \
	.deps \
	deps \
	.test \
	test \
	.build \
	build \
	.run \
	run \
	generate \
	$(BUILD_DIR)
