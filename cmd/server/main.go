package main

import (
	"log"
	"net/http"

	"skarmdump-backend/internal/config"
	"skarmdump-backend/internal/handler"
	"skarmdump-backend/internal/ratelimit"
	"skarmdump-backend/internal/service"
	"skarmdump-backend/internal/storage"
)

func main() {
	cfg := config.Load()

	local := storage.NewLocal(cfg.LocalPath)
	s3 := storage.NewS3(cfg)

	svc := service.New(local, s3, cfg.HashLength)

	rl := ratelimit.New()

	h := handler.New(svc, cfg, rl)

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", h.Upload)
	mux.HandleFunc("/", h.Get)

	log.Println("start:", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
