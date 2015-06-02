package mgmt

import (
	"encoding/json"
	"net/http"
)

func respondErrorCode(w http.ResponseWriter, code int) error {
	response := map[string]interface{}{"error": http.StatusText(code)}
	return respondJson(w, code, response)
}

func respondJson(w http.ResponseWriter, code int, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}
