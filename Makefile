default: build test

clean: 
	rm -rf build

deps:
	go get gopkg.in/mgo.v2
	go get github.com/justinas/alice
	go get github.com/michaelklishin/rabbit-hole
	go get github.com/emicklei/go-restful
	go get github.com/emicklei/go-restful/swagger

debug: deps
	go run cmd/server/*.go

build: deps clean
	mkdir -p ./build/
	go build -o ./build/server cmd/server/*.go  

test: deps
	go test -v ./...

reset-rabbit:
	-docker stop meathook-rabbit && docker rm meathook-rabbit
	docker run -d -p 5672:5672 -p 15672:15672 \
		--name meathook-rabbit dockerfile/rabbitmq

reset-mongo:
	-docker stop meathook-mongo && docker rm meathook-mongo
	docker run -d -p 27017:27017 -p 28017:28017 \
		--name meathook-mongo dockerfile/mongodb

package: 
	docker build -t nkcraddock/meathooks .
