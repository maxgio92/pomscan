.PHONY: pomscan
pomscan:
	@go build -v .

.PHONY: test
test:
	@go test  -v ./...

.PHONY: docs
docs:
	@go run docs/docs.go

