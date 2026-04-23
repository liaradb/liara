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
	@echo "\n---------------\n Test coverage\n---------------\n"
	@go test -cover ./... \
		| sed -nE 's/^(ok)?[[:space:]]+github.com\/liaradb\/liaradb\/([^[:space:]]+)[[:space:]].*coverage:[[:space:]]+([^[:space:]]+).*/\2 \3/p' \
		| column -t

.PHONY: cover-file
cover-file:
	@echo "\n---------------\n Test coverage\n---------------\n"
	@go test -coverprofile=cover.out ./... \
		| sed -nE 's/^(ok)?[[:space:]]+github.com\/liaradb\/liaradb\/([^[:space:]]+)[[:space:]].*coverage:[[:space:]]+([^[:space:]]+).*/\2 \3/p' \
		| column -t

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
