.PHONY: default
default: build

.PHONY: build
build:
	go build .

.PHONY: run
run:
	go run .

.PHONY: perft
perft:
	perftree ./test/perft/perft.sh

.PHONY: gen
gen:
	go generate ./...

.PHONY: test
test:
	go test ./...

.PHONY: test-race
test-race:
	go test -race ./...

.PHONY: bench
bench:
	go test -bench . ./... -benchmem -run=^# | tee ./docs/benchmarks.txt

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	golangci-lint fmt

.PHONY: coverage-html
coverage-html: coverage
	go tool cover -html=coverage.out

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out -coverpkg=github.com/leonhfr/orca/... ./...

.PHONY: doc
doc:
	godoc -http=:6060

.PHONY: release
release:
	goreleaser release --snapshot --clean
