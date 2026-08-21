package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"calendar-app/internal/config"
	"calendar-app/internal/database"
	"calendar-app/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed frontend/*
var frontendFS embed.FS

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		log.Fatalf("Cannot create data directory: %v", err)
	}
	if err := os.MkdirAll(cfg.LogDir(), 0755); err != nil {
		log.Fatalf("Cannot create logos directory: %v", err)
	}

	db, err := database.Open(cfg.DBPath())
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	mw := handlers.NewMiddleware(cfg.SessionSecret)
	authHandler := handlers.NewAuthHandler(cfg, mw.Store())
	apiHandler := handlers.NewAPIHandler(db)
	logoHandler := handlers.NewLogoHandler(cfg.LogDir())

	frontendSub, err := fs.Sub(frontendFS, "frontend")
	if err != nil {
		log.Fatalf("Cannot setup frontend filesystem: %v", err)
	}
	fileServer := http.FileServer(http.FS(frontendSub))

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		data, _ := fs.ReadFile(frontendSub, "login.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})
	r.Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)

	r.Get("/api/logo", logoHandler.GetLogo)

	r.Group(func(r chi.Router) {
		r.Use(mw.Auth)

		r.Post("/api/change-password", authHandler.ChangePassword)

		r.Get("/api/holidays", apiHandler.GetHolidays)
		r.Post("/api/holidays", apiHandler.AddHoliday)
		r.Delete("/api/holidays/{id}", apiHandler.DeleteHoliday)

		r.Get("/api/menu-items", apiHandler.GetMenuItems)
		r.Post("/api/menu-items", apiHandler.AddMenuItem)
		r.Delete("/api/menu-items/{id}", apiHandler.DeleteMenuItem)

		r.Get("/api/calendar/{year}/{month}", apiHandler.GetCalendarData)
		r.Post("/api/calendar/{year}/{month}", apiHandler.SaveCalendarData)

		r.Post("/api/logo", logoHandler.UploadLogo)

		r.Get("/api/export", apiHandler.ExportData)
		r.Post("/api/import", apiHandler.ImportData)

		r.Handle("/static/*", fileServer)
		r.Handle("/css/*", fileServer)
		r.Handle("/js/*", fileServer)
		r.Handle("/fonts/*", fileServer)

		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			data, _ := fs.ReadFile(frontendSub, "index.html")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
		})
	})

	srv := &http.Server{
		Addr:          fmt.Sprintf(":%d", cfg.Port),
		Handler:       r,
		ReadTimeout:   10 * time.Second,
		WriteTimeout:  30 * time.Second,
		IdleTimeout:   60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Server starting on http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
	if err := db.Close(); err != nil {
		log.Printf("DB close error: %v", err)
	}
	log.Println("Server stopped")
}
