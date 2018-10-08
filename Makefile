

EXAMPLES=rng password hackernews nena hodor pocli discordia

.PHONY: build test examples clean fmt

all: build test

build:
	go build

test: build
	go test ./...

examples: build
	for e in $(EXAMPLES) ; do (cd examples/$$e && echo $$e: && go run *.go ) ; done

clean:
	go clean
	for e in $(EXAMPLES) ; do (cd examples/$$e && go clean ) ; done

fmt:
	go fmt ./...

.PHONY: report
report:
	-go get -u github.com/client9/misspell/cmd/misspell
	-go get -u github.com/fzipp/gocyclo
	-misspell *.go lib
	-gocyclo -top 15 -avg .
	-go tool vet .
