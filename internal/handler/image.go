package handler

import (
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"skarmdump-backend/internal/config"
	"skarmdump-backend/internal/page"
	"skarmdump-backend/internal/ratelimit"
	"skarmdump-backend/internal/response"
	"skarmdump-backend/internal/service"
	"skarmdump-backend/pkg/utils"
)

type Handler struct {
	svc *service.Service
	cfg *config.Config
	rl  *ratelimit.RateLimiter
}

func New(s *service.Service, cfg *config.Config, rl *ratelimit.RateLimiter) *Handler {
	return &Handler{s, cfg, rl}
}

func ip(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func (h *Handler) authenticate(r *http.Request) (int, string, string) {
	if h.cfg.APIKey == "" {
		return http.StatusOK, "", "" // ключ не настроен — пропускаем
	}

	key := r.Header.Get("Authorization")
	if key == "" {
		return http.StatusUnauthorized, "auth_required", "authorization header is required"
	}
	if key != h.cfg.APIKey {
		return http.StatusForbidden, "invalid_key", "invalid API key"
	}
	return http.StatusOK, "", ""
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if status, code, msg := h.authenticate(r); status != http.StatusOK {
		response.Error(w, status, code, msg)
		return
	}

	if !h.rl.Allow(ip(r), time.Second) {
		response.Error(w, http.StatusTooManyRequests, "rate_limit", "rate limit exceeded")
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "missing_image", "image file is required")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "read_error", "cannot read uploaded image")
		return
	}

	if err := utils.ValidatePNG(data); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_image", "only PNG files are accepted")
		return
	}

	useS3 := r.FormValue("s3") == "true"

	hash, err := h.svc.Upload(data, useS3)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "storage_error", "failed to save image")
		return
	}

	response.Success(w, h.cfg.Domain+"/"+hash)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if !h.rl.Allow(ip(r), time.Millisecond*333) {
		response.Error(w, http.StatusTooManyRequests, "rate_limit", "rate limit exceeded")
		return
	}

	hash := strings.TrimPrefix(r.URL.Path, "/")
	hash = strings.TrimSuffix(hash, "/")

	if !h.svc.Exists(hash) {
		response.Error(w, http.StatusNotFound, "not_found", "screenshot not found")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	page.RenderOGPage(w, page.OGData{
		Title:    "Screenshot " + hash,
		ImgURL:   page.OGImgURL(h.cfg.Domain, hash),
		PageURL:  page.OGPageURL(h.cfg.Domain, hash),
		SiteName: h.cfg.SiteName,
	})
}
