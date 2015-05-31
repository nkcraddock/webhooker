package mgmt

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nkcraddock/webhooker/webhooks"
)

type hooks struct {
	store webhooks.Store
	loc   func(h *webhooks.Hook) string
}

func NewHooksHandler(store webhooks.Store) Handler {
	return &hooks{
		store: store,
	}
}

func (h *hooks) RegisterRoutes(r *mux.Router) {
	get := r.HandleFunc("/hooks/{hook}", h.get).Methods("GET")
	// passing the route into this closure so we can use it later
	// to get resource URLs
	h.loc = func(h *webhooks.Hook) string {
		url, _ := get.URL("hook", h.Id)
		return url.String()
	}

	r.HandleFunc("/hooks", h.list).Methods("GET")
	r.HandleFunc("/hooks", h.save).Methods("POST")
}

func (h *hooks) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hook, err := h.store.GetHook(vars["hook"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}
	writeEntity(w, hook)
}

func (h *hooks) list(w http.ResponseWriter, r *http.Request) {
	hooks, _ := h.store.GetHooks("")
	writeEntity(w, hooks)
}

func (h *hooks) save(w http.ResponseWriter, r *http.Request) {
	hook := webhooks.NewHook("", 0)
	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	json.Unmarshal(body, hook)
	h.store.SaveHook(hook)
	w.Header().Add("Location", h.loc(hook))
	w.WriteHeader(http.StatusCreated)
	writeEntity(w, hook)
}

func writeEntity(w http.ResponseWriter, entity interface{}) error {
	body, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	return nil
}
