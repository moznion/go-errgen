.PHONY: test

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
