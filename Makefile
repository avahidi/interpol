

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

clean:
	go clean
	for e in $(EXAMPLES) ; do (cd examples/$$e && go clean ) ; done
	rm -rf parts prime stage snap police_*.snap
	rm -f *.snap.xdelta3

SNAPTARGETS = amd64 arm64 armhf # ppc64 i686
.PHONY: snap
snap:
	make clean
	make build
	make test
	for t in $(SNAPTARGETS) ; do set -e ; snapcraft clean ; snapcraft  --debug --target-arch $$t ; done
	snapcraft login
	for s in $$(ls *.snap) ; do set -e ; snapcraft push  --release edge,beta  $$s ; done
	snapcraft logout

snap-test:
	make clean
	snapcraft clean
	snapcraft --debug

snap-setup:
	# i386
	# sudo apt install -y gcc-i686-linux-gnu
	# armhf
	sudo apt install -y binutils-arm-linux-gnueabihf cpp-7-arm-linux-gnueabihf cpp-arm-linux-gnueabihf gcc-7-arm-linux-gnueabihf gcc-7-arm-linux-gnueabihf-base gcc-arm-linux-gnueabihf libasan4-armhf-cross libatomic1-armhf-cross libc6-armhf-cross libc6-dev-armhf-cross libcilkrts5-armhf-cross libgcc-7-dev-armhf-cross libgcc1-armhf-cross libgomp1-armhf-cross libstdc++6-armhf-cross libubsan0-armhf-cross linux-libc-dev-armhf-cross
	# arm64
	# ??

fmt:
	go fmt ./...

.PHONY: report
report:
	-go get -u github.com/client9/misspell/cmd/misspell
	-go get -u github.com/fzipp/gocyclo
	-misspell *.go lib
	-gocyclo -top 15 -avg .
	-go tool vet .
