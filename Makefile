default: deps debug

debug:
	go run webapp/main.go

deps:
	go get github.com/sirupsen/logrus