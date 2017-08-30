

EXAMPLES=rng

.PHONY: build test examples clean fmt

all: build test

build:
	go build

test: build
	cd test && go test -v

examples: build
	for e in $(EXAMPLES) ; do cd examples/$$e && go run *.go; done

clean:
	go clean
	for e in $(EXAMPLES) ; do cd examples/$$e && go clean ; done

fmt:
	go fmt
	cd test && go fmt
	for e in $(EXAMPLES) ; do cd examples/$$e && go fmt; done