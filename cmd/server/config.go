package main

import (
	"flag"
)

type config struct {
	HostUrl   string
	MongoUrl  string
	MongoDb   string
	RabbitUrl string
	LogLevel  string
}

func loadConfig() config {
	host := flag.String("h", ":3001", "The host url to listen on")
	mongo := flag.String("m", "localhost", "The URL of the mongo server")
	rabbit := flag.String("r", "localhost", "The URL of the rabbit server")
	database := flag.String("db", "meathooks", "The name of the mongo database")
	logLevel := flag.String("l", "debug", "The log level to output to logs. (debug|info|warning|error|fatal)")

	flag.Parse()

	return config{
		HostUrl:   *host,
		MongoUrl:  *mongo,
		MongoDb:   *database,
		RabbitUrl: *rabbit,
		LogLevel:  *logLevel,
	}
}
