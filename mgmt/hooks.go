package mgmt

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type hooks struct {
}

func NewHooksHandler() Handler {
	return new(hooks)
}

func (h *hooks) RegisterRoutes(r *mux.Route) {
	r.Path("/hooks").HandlerFunc(h.list).
		Methods("GET")
}

func (h *hooks) list(w http.ResponseWriter, r *http.Request) {
	log.Println("LIST HOOKS")
}
