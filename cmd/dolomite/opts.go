package main

import (
	"flag"
	"strconv"
)

type options struct {
	RabbitUri   string
	RabbitVhost string
	RabbitUser  string
	RabbitPass  string
	RedisUri    string
	RedisDb     int
}

func getOpts() options {
	rabbitUri := flag.String("qurl", "http://localhost:15672", "The URL of the rabbit host")
	rabbitVhost := flag.String("qvhost", "webhooker", "The name of the rabbit virtual host")
	rabbitUser := flag.String("quser", "guest", "The username to use for rabbit")
	rabbitPass := flag.String("qpass", "guest", "The password to use for rabbit")
	redisHost := flag.String("rhost", "localhost:6379", "The redis host")
	redisDb := flag.String("rdb", "1", "The redis db number (1-16)")

	flag.Parse()

	db, err := strconv.Atoi(*redisDb)
	if err != nil {
		panic(err)
	}

	return options{
		RabbitUri:   *rabbitUri,
		RabbitVhost: *rabbitVhost,
		RabbitUser:  *rabbitUser,
		RabbitPass:  *rabbitPass,
		RedisUri:    *redisHost,
		RedisDb:     db,
	}
}
