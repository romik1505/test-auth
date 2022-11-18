package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/romik1505/auth/internal/services/auth"
)

var (
	ErrorBadRequest = errors.New("bad request")
	ErrValidate     = errors.New("validate error")
)

func ErrorWrapper(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			log.Printf("Error: %v", err)
			if errors.Is(err, ErrorBadRequest) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if errors.Is(err, ErrValidate) {
				http.Error(w, err.Error(), http.StatusOK)
				return
			}
			if errors.Is(err, auth.ErrUserAlreadyExist) {
				http.Error(w, err.Error(), http.StatusOK)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
