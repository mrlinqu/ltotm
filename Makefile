APP_NAME=ltotm

export GO111MODULE=on
export GOPROXY=
export GOSUMDB=off

BIN_DIR:=./bin

SRC:=./cmd/${APP_NAME}
TRG:=${BIN_DIR}/${APP_NAME}
TMP_STORAGE:=${BIN_DIR}/storage

MINIFY_APP:=./cmd/minify
ASSETS_DIR:=./internal/webserver/assets

LDFLAGS:=

BUILD_ENVPARMS := CGO_ENABLED=0

.PHONY: generate
generate:
	go run $(MINIFY_APP) $(ASSETS_DIR)

.PHONY: build-win
build-win: generate
	$(info #Building win/amd64...)
	GOOS=windows GOARCH=amd64 $(BUILD_ENVPARMS) go build -ldflags "$(LDFLAGS)" -o $(TRG).exe $(SRC)

.PHONY: build-linux
build-linux: generate
	$(info #Building linux/amd64...)
	GOOS=linux GOARCH=amd64 $(BUILD_ENVPARMS) go build -ldflags "$(LDFLAGS)" -o $(TRG) $(SRC)

.PHONY: build
build: build-win

.PHONY: run
run: generate build
	$(info #Running...)
	mkdir -p ${TMP_STORAGE}
	${TRG}.exe --dir=${TMP_STORAGE}
	#$(BUILD_ENVPARMS) go run -ldflags "$(LDFLAGS)" $(SRC)

.PHONY: lint
lint:
	GOOS=windows GOARCH=amd64 go vet ./...

.PHONY: test
test:
	go test -v -count=1 ./...
