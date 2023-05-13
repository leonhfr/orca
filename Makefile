.PHONY: default
default: test

.PHONY: gen
gen:
	go generate ./...

.PHONY: test
test:
	go test ./...

.PHONY: bench
bench:
	go test -bench . ./... -benchmem -run=^# | tee ./docs/benchmarks.txt

.PHONY: lint
lint:
	golangci-lint run

.PHONY: doc
doc:
	godoc -http=:6060
