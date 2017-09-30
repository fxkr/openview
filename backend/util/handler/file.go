package handler

import (
	"github.com/fxkr/openview/backend/util/safe"
	"net/http"
)

type FileHandler struct {
	Path safe.Path
}

// Statically assert that *FileHandler implements http.Handler.
var _ http.Handler = (*FileHandler)(nil)

func (h *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, h.Path.String())
}
