## Install library for production
deps:
	go get -u ./...
.PHONY: deps

## Install library for development
devel-deps: deps
	GO111MODULE=off go get \
		golang.org/x/lint/golint \
		honnef.co/go/tools/staticcheck \
		github.com/kisielk/errcheck \
		golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow \
		github.com/securego/gosec/cmd/gosec \
		github.com/motemen/gobump/cmd/gobump \
		github.com/Songmu/make2help/cmd/make2help
.PHONY: devel-deps

## Execute unit test
test: deps
	go test -v -count=1 -cover ./...
.PHONY: test

## Output test coverage
cov:
	go test -count=1 -coverprofile=cover.out ./...
	go tool cover -html=cover.out
.PHONY: cov

## Clean up artifact
clean:
	rm -rf ./artifact/*
.PHONY: clean

## Lint
lint: devel-deps
	go vet ./...
	staticcheck ./...
	errcheck ./...
	gosec -quiet ./...
	golint -set_exit_status ./...
	shadow ./...
	gosec -quiet ./...
.PHONY: lint

## Build downloader in file
build-downloader-in-file:
	go build -o dl1 ./handlers/in_file
.PHONY: build-downloader-in-file

## Build downloader in memory
build-downloader-in-memory:
	go build -o dl2 ./handlers/in_memory
.PHONY: build-downloader-in-memory

## Help
help:
	@make2help $(MAKEFILE_LIST)
.PHONY: help
