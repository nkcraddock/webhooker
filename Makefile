VENDOR=$(CURDIR)/_vendor
GOPATH=$(VENDOR):$(realpath ../../../..)
default: build test

clean: 
	rm -rf build

vendor:
	GOPATH=$(VENDOR)
	mkdir -p $(VENDOR)
	go get -d github.com/onsi/ginkgo
	go get -d github.com/onsi/gomega
	go get -d github.com/michaelklishin/rabbit-hole
	go get -d github.com/emicklei/go-restful
	go get -d github.com/emicklei/go-restful/swagger
	go get -d github.com/nu7hatch/gouuid
	go get -d gopkg.in/redis.v3
	find $(VENDOR) -type d -name '.git' | xargs rm -rf

debug: vendor
	go run cmd/server/*.go

build: vendor clean
	mkdir -p ./build/
	go build -o ./build/server cmd/server/*.go  

test: 
	go test -v ./...

reset-rabbit:
	-docker stop webhooker-rabbit && docker rm webhooker-rabbit
	docker run -d -p 5672:5672 -p 15672:15672 \
		--name webhooker-rabbit dockerfile/rabbitmq

reset-mongo:
	-docker stop webhooker-mongo && docker rm webhooker-mongo
	docker run -d -p 27017:27017 -p 28017:28017 \
		--name webhooker-mongo dockerfile/mongodb

package: 
	docker build -t nkcraddock/webhooker .
