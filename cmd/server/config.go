package main

import (
	"flag"
)

type config struct {
	HostUrl        string
	RabbitUri      string
	RabbitUsername string
	RabbitPassword string
	LogLevel       string
	ClientRoot     string
}

func loadConfig() config {
	host := flag.String("h", ":3001", "The host url to listen on")
	clientroot := flag.String("c", "", "The relative path to the client files (default will use bindata)")
	rabbitUri := flag.String("rabbitUri", "http://localhost:15672", "The URL of the rabbit management server")
	rabbitUsername := flag.String("rabbitUser", "guest", "The username to use for rabbit")
	rabbitPassword := flag.String("rabbitPass", "guest", "The password to use for rabbit")
	logLevel := flag.String("l", "debug", "The log level to output to logs. (debug|info|warning|error|fatal)")

	flag.Parse()

	return config{
		HostUrl:        *host,
		RabbitUri:      *rabbitUri,
		RabbitUsername: *rabbitUsername,
		RabbitPassword: *rabbitPassword,
		LogLevel:       *logLevel,
		ClientRoot:     *clientroot,
	}
}
