package mgmt

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nkcraddock/webhooker/webhooks"

	"github.com/gorilla/mux"
)

type filtersHandler struct {
	store webhooks.Store
	loc   func(f *webhooks.Filter) string
}

// construct a new filters handler
func NewFiltersHandler(store webhooks.Store) Handler {
	return &filtersHandler{
		store: store,
	}
}

// wire my routes up to the mux router
func (h *filtersHandler) RegisterRoutes(r *mux.Router) {
	get := r.HandleFunc("/hooks/{hook}/filters/{filter}", h.get).Methods("GET")
	// passing the route into this closure so we can use it later
	// to get resource URLs
	h.loc = func(f *webhooks.Filter) string {
		url, _ := get.URL("hook", f.Hook, "filter", f.Id)
		return url.String()
	}

	r.HandleFunc("/hooks/{hook}/filters", h.list).Methods("GET")
	r.HandleFunc("/hooks/{hook}/filters", h.save).Methods("POST")
}

func (h *filtersHandler) get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hookId := vars["hook"]
	filterId := vars["filter"]

	filter, err := h.store.GetFilter(hookId, filterId)
	if err != nil {
		respondErrorCode(w, http.StatusNotFound)
		return
	}

	respondJson(w, http.StatusOK, filter)
}

func (h *filtersHandler) list(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hookId := vars["hook"]

	filters, err := h.store.GetFilters(hookId)

	if err != nil {
		respondErrorCode(w, http.StatusNotFound)
		return
	}

	respondJson(w, http.StatusOK, filters)
}

func (h *filtersHandler) save(w http.ResponseWriter, r *http.Request) {
	hook := mux.Vars(r)["hook"]

	// Make sure the hook exists
	_, err := h.store.GetHook(hook)
	if err != nil {
		respondErrorCode(w, http.StatusNotFound)
		return
	}

	// Parse the filter
	f, err := loadFilterFromRequest(r)
	if err != nil {
		respondErrorCode(w, http.StatusBadRequest)
		return
	}

	// Make sure the hook is set correctly
	f.Hook = hook

	// Save the filter
	h.store.SaveFilter(f)

	// Write the location header
	w.Header().Set("Location", h.loc(f))

	// Write the entity
	respondJson(w, http.StatusCreated, f)
}

func loadFilterFromRequest(r *http.Request) (*webhooks.Filter, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body.Close()

	filter := new(webhooks.Filter)
	if err := json.Unmarshal(b, &filter); err != nil {
		return nil, err
	}
	return filter, nil
}
