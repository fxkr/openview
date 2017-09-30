package handler

import (
	"bufio"
	"net/http"
)

type ByteHandler struct {
	Bytes       []byte
	ContentType string
}

// Statically assert that *ByteHandler implements http.Handler.
var _ http.Handler = (*ByteHandler)(nil)

func (h *ByteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", h.ContentType)
	w.WriteHeader(http.StatusOK)
	bufio.NewWriter(w).Write(h.Bytes)
}
