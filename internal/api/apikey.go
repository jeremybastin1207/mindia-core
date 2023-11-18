package api

import (
	"net/http"

	"github.com/jeremybastin1207/mindia-core/internal/apikey"
)

func apiKeyMiddleware(keyStorer apikey.Storer, masterKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apikey, err := bearerToken(r, "Authorization")
			if err != nil {
				writeError(w, http.StatusUnauthorized, err)
				return
			}
			if apikey != masterKey {
				_, err := keyStorer.GetByKey(apikey)
				if err != nil {
					writeError(w, http.StatusUnauthorized, err)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
