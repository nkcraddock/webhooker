package webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type HttpHandler struct {
	hooks Store
}

func NewHttpHandler(store Store) *HttpHandler {
	return &HttpHandler{
		hooks: store,
	}
}

func (h *HttpHandler) Post(w http.ResponseWriter, req *http.Request) {
	hook, err := ParseWebhookFromRequest(req)

	if failOnError(w, err) {
		return
	}

	err = h.hooks.Add(&hook)

	if failOnError(w, err) {
		return
	}

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

	w.WriteHeader(200)
}

func failOnError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	fmt.Fprintf(w, "An error occurred: %s", err.Error())
	w.WriteHeader(500)
	return true
}

func ParseWebhookFromRequest(req *http.Request) (hook Webhook, err error) {
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&hook)
	return
}
