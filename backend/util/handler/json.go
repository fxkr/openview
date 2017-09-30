package handler

import (
	"encoding/json"
	"net/http"
)

type JSONHandler struct {
	Data interface{}
}

// Statically assert that *JSONHandler implements http.Handler.
var _ http.Handler = (*JSONHandler)(nil)

func (h *JSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.Data)
}
