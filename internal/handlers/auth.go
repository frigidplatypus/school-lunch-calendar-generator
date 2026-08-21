package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"calendar-app/internal/config"
	"calendar-app/internal/models"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	cfg    *config.Config
	store  *sessions.CookieStore
}

func NewAuthHandler(cfg *config.Config, store *sessions.CookieStore) *AuthHandler {
	return &AuthHandler{cfg: cfg, store: store}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid request"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(h.cfg.PasswordHash), []byte(req.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid password"})
		return
	}

	session, _ := h.store.Get(r, SessionName)
	session.Values[SessionKey] = true
	if err := session.Save(r, w); err != nil {
		log.Printf("session save error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, SessionName)
	session.Values[SessionKey] = false
	session.Options.MaxAge = -1
	session.Save(r, w)

	if r.Header.Get("Accept") == "application/json" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.cfg.PasswordExternal() {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "password is managed externally and cannot be changed via the UI"})
		return
	}

	var req models.PasswordChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid request"})
		return
	}

	if len(req.NewPassword) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "password must be at least 4 characters"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "failed to hash password"})
		return
	}

	if err := h.cfg.UpdatePasswordHash(string(hash)); err != nil {
		log.Printf("password update error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "failed to update password"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
