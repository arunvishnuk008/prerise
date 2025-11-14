package middleware

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

func PanicRecovery(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				err := recover()
				if err != nil {
					logger.Error("Error", err)
					jsonError, _ := json.Marshal(map[string]string{
						"error": "There was an internal server error",
					})

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, err := w.Write(jsonError)
					if err != nil {
						log.Print("Failed to write panic as json response")
					}
				}

			}()

			next.ServeHTTP(w, r)

		})
	}
}
