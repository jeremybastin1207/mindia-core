package api

import (
	"net/http"

	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func apiHandler(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if err, ok := err.(*mindiaerr.Error); ok {
				switch err.ErrCode {
				case mindiaerr.ErrCodeMediaNotFound:
					w.WriteHeader(http.StatusNotFound)
				case mindiaerr.ErrCodeMimeTypeNotSupported:
					writeError(w, http.StatusBadRequest, err)
				case mindiaerr.ErrCodeUnauthorizedRequest:
					writeError(w, http.StatusUnauthorized, err)
				case mindiaerr.ErrCodeServiceUnavailable:
					writeError(w, http.StatusInternalServerError, err)
				case mindiaerr.ErrCodeNamedTransformationNotFound:
					writeError(w, http.StatusBadRequest, err)
				default:
					writeError(w, http.StatusInternalServerError, err)
				}
			}
		}
	}
}
