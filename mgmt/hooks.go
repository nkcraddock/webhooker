package mgmt

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nkcraddock/webhooker/webhooks"
)

type hooksHandler struct {
	store webhooks.Store
	loc   func(h *webhooks.Hook) string
}

// construct a new hooks handler
func NewHooksHandler(store webhooks.Store) Handler {
	return &hooksHandler{
		store: store,
	}
}

// wire my routes up to the mux router
func (h *hooksHandler) RegisterRoutes(r *mux.Router) {
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

// GET /hooks/:id - get a single hook by id
func (h *hooksHandler) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hook, err := h.store.GetHook(vars["hook"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}

	writeEntity(w, hook.Sanitize())
}

// GET /hooks - get a summary list of all hooks
func (h *hooksHandler) list(w http.ResponseWriter, r *http.Request) {
	hooks, _ := h.store.GetHooks("")
	cleanhooks := make([]interface{}, len(hooks))
	for i, h := range hooks {
		cleanhooks[i] = h.Sanitize()
	}
	writeEntity(w, cleanhooks)
}

// POST /hooks {hook} - add a new hook
func (h *hooksHandler) save(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	hook, _ := hookFromJson(body)
	hook.Id = ""

	h.store.SaveHook(hook)
	w.Header().Add("Location", h.loc(hook))
	w.WriteHeader(http.StatusCreated)
	writeEntity(w, hook)
}

// helper function - write an entity to the writer in json format
func writeEntity(w http.ResponseWriter, entity interface{}) error {
	body, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	return nil
}

func hookFromJson(data []byte) (*webhooks.Hook, error) {
	hook := new(webhooks.Hook)
	if err := json.Unmarshal(data, &hook); err != nil {
		return nil, err
	}
	return hook, nil
}
