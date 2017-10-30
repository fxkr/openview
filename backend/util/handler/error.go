package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ErrorHandler struct {
	status int
	cause  error
}

// Statically assert that *ErrorHandler implements http.Handler.
var _ http.Handler = (*ErrorHandler)(nil)

func Error(cause error) *ErrorHandler {
	return getWebError(http.StatusInternalServerError, cause)
}

func StatusError(status int, cause error) *ErrorHandler {
	return &ErrorHandler{status, cause}
}

func Status(status int) *ErrorHandler {
	return &ErrorHandler{status, errors.New(http.StatusText(status))}
}

// Cause implements github.com/pkg/errors.causer.
func (e *ErrorHandler) Cause() error {
	return e.cause
}

// Error implements error.error.
func (e *ErrorHandler) Error() string {
	return http.StatusText(e.status)
}

// MarshalJSON implements json.Marshaler.
func (e *ErrorHandler) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Status int
		Error  string
	}{
		e.status,
		fmt.Sprintf("%+v", e.cause),
	})
}

// For http.Handler.
func (e *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.WithError(e).WithFields(log.Fields{
		"status": e.status,
		"path":   r.URL.Path,
	}).Error("Request failed")

	// Show pretty-printed response for manual requests
	if !strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.Header().Set("X-Content-bType-Options", "nosniff")
		w.WriteHeader(e.status)
		fmt.Fprintf(w, "%d %s\n\n%+v\n\n", e.status, http.StatusText(e.status), e.cause)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(e.status)

	enc := json.NewEncoder(w)
	enc.Encode(e)
}

func getWebError(defaultStatus int, err error) (webErr *ErrorHandler) {
	type causer interface {
		Cause() error
	}

	for err != nil {
		webErr, found := err.(*ErrorHandler)
		if found {
			return webErr
		}

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}

	return &ErrorHandler{defaultStatus, err}
}
