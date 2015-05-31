VENDOR=$(CURDIR)/_vendor
GOPATH=$(VENDOR):$(realpath ../../../..)
SERVER_FILES := $(shell find cmd/server -type f -name "*.go" ! -name "*_test.go")

default: build test

clean: 
	rm -rf build

vendor:
	GOPATH=$(VENDOR)
	mkdir -p $(VENDOR)
	go get -d github.com/gorilla/mux
	go get -d github.com/onsi/ginkgo
	go get -d github.com/onsi/gomega
	go get -d github.com/michaelklishin/rabbit-hole
	go get -d github.com/nu7hatch/gouuid
	go get -d gopkg.in/redis.v3
	find $(VENDOR) -type d -name '.git' | xargs rm -rf

run:
	cd $(CURDIR)
	go run $(SERVER_FILES)

build: vendor clean
	mkdir -p ./build/
	CGO_ENABLED=0 go build -a -installsuffix cgo -o build/mervis --ldflags '-s' $(SERVER_FILES)

test: 
	go test -v ./...

