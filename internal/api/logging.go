package api

import (
	"fmt"
	"net/http"

	mlogging "github.com/jeremybastin1207/mindia-core/internal/logging"
)

func logging(logger *mlogging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Info(fmt.Sprintf("| %s | %s | %s | %s | %s", requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent()))
			}()
			next.ServeHTTP(w, r)
		})
	}
}
