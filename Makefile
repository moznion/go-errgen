.PHONY: test

RELEASE_DIR=bin
PKGS := $(shell go list ./...)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
CURRENT_VERSION := $(shell git tag | sort -r | head -1)

check: test lint vet fmt-check

test: test-clean test-build test-gen
	go test -v $(PKGS)

test-clean:
	rm -f ./test/*_errmsg_gen.go

test-build:
	go build -o ./test/bin/errgen cmd/errgen/errgen.go

test-gen:
	PATH="$$(pwd)/test/bin/:$$PATH" go generate ./...

lint:
	golint $(PKGS)

vet:
	go vet $(PKGS)

fmt-check:
	gofmt -l -s *.go | grep [^*][.]go$$; \
	EXIT_CODE=$$?; \
	if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi; \
	goimports -l *.go | grep [^*][.]go$$; \
	EXIT_CODE=$$?; \
	if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi \

fmt:
	gofmt -w -s *.go
	goimports -w *.go

installdeps:
	GO111MODULE=on go mod vendor

clean:
	find bin -type f | grep -v .gitkeep | xargs rm -f

all: installdeps build-linux-amd64 build-linux-386 build-darwin-amd64 build-darwin-386 build-windows-amd64 build-windows-386

build: $(RELEASE_DIR)/errgen_$(GOOS)_$(GOARCH)

build-linux-amd64:
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-386:
	@$(MAKE) build GOOS=linux GOARCH=386

build-darwin-amd64:
	@$(MAKE) build GOOS=darwin GOARCH=amd64

build-darwin-386:
	@$(MAKE) build GOOS=darwin GOARCH=386

build-windows-amd64:
	@$(MAKE) build GOOS=windows GOARCH=amd64

build-windows-386:
	@$(MAKE) build GOOS=windows GOARCH=386

$(RELEASE_DIR)/errgen_$(GOOS)_$(GOARCH):
	go build \
		-ldflags "-X main.revision=$(CURRENT_REVISION) -X main.version=$(CURRENT_VERSION)" \
		-o $(RELEASE_DIR)/errgen_$(GOOS)_$(GOARCH)_$(CURRENT_VERSION) \
		./cmd/errgen

