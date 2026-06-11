package handler

import (
	"net/http"

	"skarmdump-backend/internal/response"
)

func UserAgentMiddleware(cfgUserAgent string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cfgUserAgent != "" && r.UserAgent() != cfgUserAgent {
			response.Error(w, http.StatusForbidden, "bad_user_agent", "request from this User-Agent is not allowed")
			return
		}
		next.ServeHTTP(w, r)
	})
}
