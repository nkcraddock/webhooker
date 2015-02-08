default: deps debug

deps:
	go get github.com/gorilla/mux
	go get github.com/sirupsen/logrus
	go get gopkg.in/mgo.v2

debug:
	go run cmd/server/*.go

build: deps
	mkdir -p ./build/
	go build -o ./build/server cmd/server/server.go

package: 
	docker build -t nkcraddock/meathooks .