VENDOR=$(CURDIR)/_vendor
GOPATH=$(VENDOR):$(realpath ../../../..)
SERVER_FILES := $(shell find cmd/server -type f -name "*.go" ! -name "*_test.go")
MGMTCLIENT=$(CURDIR)/mgmt/client

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

run-server:
	cd $(CURDIR)
	go run $(SERVER_FILES) -c mgmt/client/web/build/

run:
	go run cmd/dolomite/*.go

build: vendor clean
	mkdir -p ./build/
	CGO_ENABLED=0 go build -a -installsuffix cgo -o build/mgmt-service --ldflags '-s' $(SERVER_FILES)

test: 
	go test -v -cover ./...

mgmt-client-deps:
	if [ ! -d "$(MGMTCLIENT)/web/node_modules" ]; then \
		cd $(MGMTCLIENT)/web; \
		npm install; \
		bower install; \
		fi;

mgmt-client: mgmt-client-deps
	GOPATH=$(VENDOR)
	mkdir -p $(VENDOR)
	go get github.com/jteeuwen/go-bindata/...
	grunt --gruntfile $(MGMTCLIENT)/web/Gruntfile.js package
	$(VENDOR)/bin/go-bindata -o "$(MGMTCLIENT)/resources.go" -pkg="client" -prefix="$(MGMTCLIENT)/web/build/" $(MGMTCLIENT)/web/build/...
