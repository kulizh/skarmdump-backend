package handler

import (
	"net/http"
)

func (h *Handler) checkUserAgent(r *http.Request) (int, string, string) {
	if h.cfg.UserAgent == "" {
		return http.StatusOK, "", ""
	}
	if r.UserAgent() != h.cfg.UserAgent {
		return http.StatusForbidden, "bad_user_agent", "request from this User-Agent is not allowed"
	}
	return http.StatusOK, "", ""
}
