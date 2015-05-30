package main

import (
	"flag"
)

type config struct {
	HostUrl        string
	MongoUrl       string
	MongoDb        string
	RabbitUri      string
	RabbitUsername string
	RabbitPassword string
	LogLevel       string
}

func loadConfig() config {
	host := flag.String("h", ":3001", "The host url to listen on")
	mongo := flag.String("m", "localhost", "The URL of the mongo server")
	rabbitUri := flag.String("rabbitUri", "http://localhost:15672", "The URL of the rabbit management server")
	rabbitUsername := flag.String("rabbitUser", "guest", "The username to use for rabbit")
	rabbitPassword := flag.String("rabbitPass", "guest", "The password to use for rabbit")
	database := flag.String("db", "webhooker", "The name of the mongo database")
	logLevel := flag.String("l", "debug", "The log level to output to logs. (debug|info|warning|error|fatal)")

	flag.Parse()

	return config{
		HostUrl:        *host,
		MongoUrl:       *mongo,
		MongoDb:        *database,
		RabbitUri:      *rabbitUri,
		RabbitUsername: *rabbitUsername,
		RabbitPassword: *rabbitPassword,
		LogLevel:       *logLevel,
	}
}
