.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -o bin/liaradb cmd/liaradb/main.go

export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318

.PHONY: run
run:
	go run cmd/liaradb/main.go

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

.PHONY: title
title:
	@echo "\n---------------\n Test coverage\n---------------\n"

TEST_REGEX = s/^(ok)?[[:space:]]+github.com\/liaradb\/liaradb\/([^[:space:]]+)[[:space:]].*coverage:[[:space:]]+([^[:space:]]+).*/\2 \3/p
TEST_PARSE = sed -nE '$(TEST_REGEX)' | column -t

.PHONY: cover
cover: title
	@go test -cover ./... | $(TEST_PARSE)

.PHONY: cover-file
cover-file: title
	@go test -coverprofile=cover.out ./... | $(TEST_PARSE)

.PHONY: report
report: cover-file
	@go tool cover -html=cover.out

.PHONY: total
total: cover-file
	@COVERAGE=$$(go tool cover -func=cover.out | grep total | grep -Eo '[0-9]+\.[0-9]\%+'); \
	echo "\nCoverage: $$COVERAGE\n"

.PHONY: report-file
report-file: cover-file total
	@go tool cover -html=cover.out -o cover.html

.PHONY: doc
doc:
	go doc -http
