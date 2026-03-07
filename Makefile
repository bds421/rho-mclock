.PHONY: test test-race test-v bench vet clean all

test:
	go test -count=1 ./...

test-race:
	go test -race -count=1 ./...

test-v:
	go test -v -race -count=1 ./...

bench:
	go test -bench=. -benchmem ./...

vet:
	go vet ./...

clean:
	go clean -testcache

all: vet test-race bench
