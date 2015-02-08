default: deps debug

deps:
	go get github.com/sirupsen/logrus

debug:
	go run cmd/server/server.go

build:
	mkdir -p ./build/
	go build -o ./build/server cmd/server/server.go