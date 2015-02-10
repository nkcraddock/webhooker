package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

var (
	config Config
	repo   *Repo
)

func init() {
	config = LoadConfig()
	log.SetOutput(os.Stdout)
	log.ParseLevel(config.LogLevel)

	log.Infof("Initializing with config %s", config)

	init_mongo(config.MongoUrl, config.MongoDb)
}

func init_mongo(url string, db string) {
	var err error
	repo, err = ConnectRepo(url, db)

	if err != nil {
		log.Fatalf("Failed to connect to mongo: %s", err.Error())
		panic("Failed to connect to mongo")
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/webhooks", WebhooksPost).Methods("POST")
	r.HandleFunc("/webhooks/{id:[0-9a-fA-F]{24}}", WebhooksGet).Methods("GET")
	r.HandleFunc("/webhooks", WebhooksList).Methods("GET")

	http.Handle("/", r)
	log.Infof("Listening on %s", config.HostUrl)
	http.ListenAndServe(config.HostUrl, nil)
}
