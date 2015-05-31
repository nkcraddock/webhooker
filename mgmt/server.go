package mgmt

import "net/http"

type MgmtServer struct {
}

func NewMgmtServer() (http.Handler, error) {
	return &MgmtServer{}, nil
}

func (s *MgmtServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
