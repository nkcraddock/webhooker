package main

import (
	"net/http"

	"github.com/nkcraddock/meathooks/db"
	"github.com/nkcraddock/meathooks/webhooks"

	"github.com/gorilla/mux"
)

var (
	conn  *db.Connection
	cfg   config
	hooks webhooks.Store
)

func init() {
	cfg = loadConfig()
	init_mongo(cfg.MongoUrl, cfg.MongoDb)
}

func init_mongo(url string, database string) {
	var err error
	conn, err = db.Dial(url, database)

	if err != nil {
		panic("Failed to connect to mongo")
	}

	hooks = webhooks.NewMongoStore(conn)
}

func main() {
	r := mux.NewRouter()

	webhooks.RegisterHandler(r, hooks)

	http.Handle("/", r)
	http.ListenAndServe(cfg.HostUrl, nil)
}
