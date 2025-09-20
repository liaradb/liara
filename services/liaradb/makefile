.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	go clean -testcache

.PHONY: cover
cover:
	go test -cover ./...

.PHONY: report
report:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out

.PHONY: report-file
report:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html
