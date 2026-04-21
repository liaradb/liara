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
	@echo "Test coverage"
	@go test -cover ./...

.PHONY: cover-file
cover-file:
	@go test -coverprofile=cover.out ./...

.PHONY: report
report: cover-file
	@go tool cover -html=cover.out

.PHONY: total
total: cover-file
	@COVERAGE=$$(go tool cover -func=cover.out | grep total | grep -Eo '[0-9]+\.[0-9]\%+'); \
	echo "Coverage: $$COVERAGE"

.PHONY: report-file
report-file: cover-file total
	@go tool cover -html=cover.out -o cover.html

.PHONY: doc
doc:
	go doc -http
