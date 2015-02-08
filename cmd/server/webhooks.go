package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"

	"net/http"

	log "github.com/sirupsen/logrus"
)

type Webhook struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	CallbackURL string        `json:"callback_url" bson:"callback_url"`
	Filter      string        `json:"filter" bson:"filter"`
}

var repo *Repo

func init() {
	repo = ConnectRepo("172.17.0.3", "meathooks")
}

func WebhooksPost(w http.ResponseWriter, req *http.Request) {
	hook, err := ParseWebhookFromRequest(req)

	if failOnError(w, err) {
		return
	}

	log.Infof("POST /webhooks - %s", hook)
	err = repo.AddWebhook(&hook)

	if failOnError(w, err) {
		return
	}

	uri := fmt.Sprintf("/webhooks/%s", hook.Id)
	w.Header().Set("Location", uri)
	w.WriteHeader(201)
}

func WebhooksList(w http.ResponseWriter, req *http.Request) {
	hooks := repo.GetWebhooks()
	b, err := json.Marshal(hooks)

	if failOnError(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func WebhooksGet(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	log.Infof("GET /webhooks/%s", id)

	webhook := repo.GetWebhook(id)
	enc := json.NewEncoder(w)
	err := enc.Encode(webhook)

	if failOnError(w, err) {
		return
	}

	w.WriteHeader(200)
}

func failOnError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	fmt.Fprintf(w, "An error occurred: %s", err.Error())
	log.Errorf("An error occurred: %s", err.Error())
	w.WriteHeader(500)
	return true
}

func ParseWebhookFromRequest(req *http.Request) (Webhook, error) {
	var hook Webhook
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&hook)
	return hook, err
}
