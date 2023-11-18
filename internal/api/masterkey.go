package api

import (
	"net/http"
)

func masterKeyMiddleware(masterKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key, err := bearerToken(r, "Authorization")
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			if key != masterKey {
				writeError(w, http.StatusUnauthorized, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

}
