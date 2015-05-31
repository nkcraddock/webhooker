package mgmt

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	RegisterRoutes(r *mux.Router)
}

type MgmtServer struct {
	m        *mux.Router
	handlers []Handler
}

func NewMgmtServer(handlers []Handler) (*MgmtServer, error) {
	m := mux.NewRouter()

	// Let all the handlers register their routes
	for _, h := range handlers {
		h.RegisterRoutes(m.PathPrefix("/api").Subrouter())
	}

	return &MgmtServer{m, handlers}, nil
}

func (s *MgmtServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}
