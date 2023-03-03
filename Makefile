.PHONY: default
default: test

.PHONY: test
test:
	go test ./...

.PHONY: bench
bench:
	go test -bench . ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: doc
doc:
	godoc -http=:6060
