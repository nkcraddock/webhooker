default: deps debug


deps:
	go get github.com/gorilla/mux
	go get github.com/sirupsen/logrus
	go get gopkg.in/mgo.v2

debug: deps
	go run cmd/server/*.go

build: deps
	mkdir -p ./build/
	go build -o ./build/server cmd/server/*.go

reset-rabbit:
	-docker stop meathook-rabbit && docker rm meathook-rabbit
	docker run -d -p 27017:27017 -p 28017:28017 \
		--name meathook-rabbit dockerfile/rabbitmq

reset-mongo:
	-docker stop meathook-mongo && docker rm meathook-mongo
	docker run -d -p 5672:5672 -p 15672:15672 \
		--name meathook-mongo dockerfile/mongodb


package: 
	docker build -t nkcraddock/meathooks .