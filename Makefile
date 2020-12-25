

EXAMPLES=hackernews nena hodor discordia

.PHONY: build test examples clean fmt

all: build test

build:
	go build
	cd cmd/police && go build
	# go build ./cmd/...

test: build
	go test ./... --cover


examples: build
	for e in $(EXAMPLES) ; do (cd examples/$$e && echo $$e: && go run *.go ) ; done

	# include examples from police so we know its working...
	go run cmd/police/main.go -sep ", " "Hello" "{{set sep=' ' data='Kitty World Dolly goodbye'}}!"
	go run cmd/police/main.go -lsep ":" "{{random min=0 max=255 count=8 format=%02x}}"


clean:
	go clean
	for e in $(EXAMPLES) ; do (cd examples/$$e && go clean ) ; done
	rm -rf parts prime stage snap police_*.snap
	rm -f *.snap.xdelta3

fmt:
	go fmt ./...

.PHONY: report
report:
	-go get -u github.com/client9/misspell/cmd/misspell
	-go get -u github.com/fzipp/gocyclo
	-misspell *.go lib
	-gocyclo -top 15 -avg .
	-go tool vet .


# ---- snaps ----

SNAPTARGETS = amd64 arm64 armhf # ppc64 i686
SNAPCRAFT=snapcraft --use-lxd --debug

.PHONY: snap
snap:
	make clean
	make build
	make test
	
	$(SNAPCRAFT)
	$(SNAPCRAFT) login
	$(SNAPCRAFT) push --release edge,beta  *.snap
	$(SNAPCRAFT) logout

	snapcraft clean # remove new garbage
	-lxc delete -f snapcraft-police


snap-setup:
	sudo snap install snapcraft #multipass
#	sudo apt install -y binutils-arm-linux-gnueabihf gcc-arm-linux-gnueabihf

