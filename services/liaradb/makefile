.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -o bin/liaradb cmd/liaradb/main.go

.PHONY: clean-all
clean-all: clean clean-data

.PHONY: clean
clean: clean-build clean-test

.PHONY: clean-build
clean-build:
	rm -rf bin
	go clean

.PHONY: clean-test
clean:
	go clean -testcache

.PHONY: clean-data
clean-data:
	rm -rf .dbdata

.PHONY: cover
cover:
	go test -cover ./...

.PHONY: report
report:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out

.PHONY: report-file
report-file:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: doc
doc:
	go doc -http
