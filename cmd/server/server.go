package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	"github.com/justinas/alice"
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
	router := mux.NewRouter()
	webhooks.RegisterHandler(router, hooks)
	chain := alice.New(loggerHandler).Then(router)
	http.ListenAndServe(cfg.HostUrl, chain)
}

func loggerHandler(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}
