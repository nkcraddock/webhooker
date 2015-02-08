package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gorilla/mux"

	"net/http"

	log "github.com/sirupsen/logrus"
)

type Webhook struct {
	CallbackURL string `json:"callback_url"`
	Filter      string `json:"filter"`
}

var webhooks []Webhook

func init() {
	webhooks = []Webhook{
		Webhook{CallbackURL: "http://localhost:3002/callback", Filter: "*"},
		Webhook{CallbackURL: "http://localhost:3002/callback2", Filter: "*"},
		Webhook{CallbackURL: "http://localhost:3002/callback3", Filter: "*"},
	}
}

func WebhooksPost(w http.ResponseWriter, req *http.Request) {
	hook, err := ParseWebhookFromRequest(req)

	if failOnError(w, err) {
		return
	}

	log.Infof("POST /webhooks - %s", hook)
	ix := len(webhooks)
	webhooks = append(webhooks, hook)
	uri := fmt.Sprintf("/webhooks/%d", ix)
	w.Header().Set("Location", uri)
	w.WriteHeader(201)
}

func WebhooksList(w http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(webhooks)

	if failOnError(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func WebhooksGet(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if failOnError(w, err) {
		return
	}
	log.Infof("GET /webhooks/%d", id)

	enc := json.NewEncoder(w)
	err = enc.Encode(webhooks[id])

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
