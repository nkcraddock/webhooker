package webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nkcraddock/rabbit-hole"
)

type HttpHandler struct {
	hooks  Store
	rabbit *rabbitFarm
}

func RegisterHandler(r *mux.Router, store Store, rabbit *rabbithole.Client) *HttpHandler {
	handler := HttpHandler{
		hooks:  store,
		rabbit: newRabbitFarm(rabbit),
	}

	r.HandleFunc("/webhooks", handler.Post).Methods("POST")
	r.HandleFunc("/webhooks/{id:[0-9a-fA-F]{24}}", handler.Delete).Methods("DELETE")
	r.HandleFunc("/webhooks/{id:[0-9a-fA-F]{24}}", handler.Get).Methods("GET")
	r.HandleFunc("/webhooks", handler.List).Methods("GET")

	return &handler
}

// Post handles POST /webhooks
func (h *HttpHandler) Post(w http.ResponseWriter, req *http.Request) {
	hook, err := ParseWebhookFromRequest(req)

	if failOnError(w, err) {
		return
	}

	err = h.hooks.Add(&hook)

	if failOnError(w, err) {
		return
	}

	h.rabbit.SaveUrlQueue(hook.Id.Hex())

	uri := fmt.Sprintf("/webhooks/%s", hook.Id.Hex())
	w.Header().Set("Location", uri)
	w.WriteHeader(201)
}

func (h *HttpHandler) List(w http.ResponseWriter, req *http.Request) {
	hooks := h.hooks.All()
	b, err := json.Marshal(hooks)

	if failOnError(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *HttpHandler) Get(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	hook := h.hooks.GetById(id)

	if len(hook.Id) == 0 {
		http.NotFound(w, req)
		return
	}

	enc := json.NewEncoder(w)
	err := enc.Encode(hook)

	if failOnError(w, err) {
		return
	}
}

func (h *HttpHandler) Delete(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	err := h.hooks.Delete(id)

	if failOnError(w, err) {
		return
	}

	h.rabbit.DeleteUrlQueue(id)

	w.WriteHeader(200)
}

func failOnError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	msg := fmt.Sprintf("An error occurred: %s", err.Error())
	http.Error(w, msg, 500)
	return true
}

func ParseWebhookFromRequest(req *http.Request) (hook Webhook, err error) {
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&hook)
	return
}
