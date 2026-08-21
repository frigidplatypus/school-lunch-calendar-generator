package handlers

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

const SessionName = "calendar-session"
const SessionKey = "authenticated"

type Middleware struct {
	store *sessions.CookieStore
}

func NewMiddleware(sessionSecret string) *Middleware {
	store := sessions.NewCookieStore([]byte(sessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 365,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return &Middleware{store: store}
}

func (m *Middleware) Store() *sessions.CookieStore {
	return m.store
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := m.store.Get(r, SessionName)
		if auth, ok := session.Values[SessionKey].(bool); !ok || !auth {
			if r.Header.Get("Accept") == "application/json" || strings.HasPrefix(r.URL.Path, "/api/") {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
