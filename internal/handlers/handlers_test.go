package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"calendar-app/internal/database"

	"github.com/go-chi/chi/v5"
)

func newTestServer(t *testing.T) (*APIHandler, chi.Router) {
	t.Helper()
	db, err := database.Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	h := NewAPIHandler(db)
	r := chi.NewRouter()
	r.Get("/api/holidays", h.GetHolidays)
	r.Post("/api/holidays", h.AddHoliday)
	r.Delete("/api/holidays/{id}", h.DeleteHoliday)
	r.Get("/api/menu-items", h.GetMenuItems)
	r.Post("/api/menu-items", h.AddMenuItem)
	r.Delete("/api/menu-items/{id}", h.DeleteMenuItem)
	r.Get("/api/calendar/{year}/{month}", h.GetCalendarData)
	r.Post("/api/calendar/{year}/{month}", h.SaveCalendarData)
	r.Get("/api/export", h.ExportData)
	r.Post("/api/import", h.ImportData)
	r.Post("/api/import/holidays", h.ImportHolidays)
	r.Post("/api/import/menu-items", h.ImportMenuItems)
	r.Post("/api/import/calendar", h.ImportCalendarData)
	return h, r
}

func doJSON(t *testing.T, r chi.Router, method, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHolidaysCRUD(t *testing.T) {
	_, router := newTestServer(t)

	w := doJSON(t, router, "POST", "/api/holidays", map[string]interface{}{"month": 2, "date": 17, "name": "Test Day"})
	if w.Code != http.StatusCreated {
		t.Fatalf("AddHoliday status = %d, body = %s", w.Code, w.Body.String())
	}
	var created struct {
		ID int `json:"id"`
	}
	json.Unmarshal(w.Body.Bytes(), &created)
	if created.ID == 0 {
		t.Fatalf("expected non-zero id, got body %s", w.Body.String())
	}

	w = doJSON(t, router, "GET", "/api/holidays", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("GetHolidays status = %d", w.Code)
	}
	var holidays []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &holidays)
	if len(holidays) == 0 {
		t.Fatal("expected holidays (seeded or created)")
	}

	w = doJSON(t, router, "DELETE", fmt.Sprintf("/api/holidays/%d", created.ID), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("DeleteHoliday status = %d", w.Code)
	}

	w = doJSON(t, router, "DELETE", "/api/holidays/notanumber", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid id, got %d", w.Code)
	}
}

func TestAddMenuItemValidation(t *testing.T) {
	_, router := newTestServer(t)

	w := doJSON(t, router, "POST", "/api/menu-items", map[string]interface{}{"meal_type": "dinner", "item_name": "Steak"})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid meal_type, got %d", w.Code)
	}

	w = doJSON(t, router, "POST", "/api/menu-items", map[string]interface{}{"meal_type": "lunch", "item_name": "Pizza"})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCalendarSaveGet(t *testing.T) {
	_, router := newTestServer(t)

	w := doJSON(t, router, "POST", "/api/calendar/2026/7", map[string]interface{}{
		"school_name": "Test School",
		"menus":       []map[string]interface{}{},
	})
	if w.Code != http.StatusOK {
		t.Fatalf("SaveCalendarData status = %d: %s", w.Code, w.Body.String())
	}

	w = doJSON(t, router, "GET", "/api/calendar/2026/7", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("GetCalendarData status = %d", w.Code)
	}
	var cd map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &cd)
	if cd["school_name"] != "Test School" {
		t.Fatalf("unexpected school_name: %v", cd["school_name"])
	}

	w = doJSON(t, router, "GET", "/api/calendar/1999/0", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing calendar, got %d", w.Code)
	}
}

func TestExportImport(t *testing.T) {
	_, router := newTestServer(t)

	doJSON(t, router, "POST", "/api/menu-items", map[string]interface{}{"meal_type": "breakfast", "item_name": "Toast"})

	w := doJSON(t, router, "GET", "/api/export", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("ExportData status = %d", w.Code)
	}

	w = doJSON(t, router, "POST", "/api/import", json.RawMessage(w.Body.Bytes()))
	if w.Code != http.StatusOK {
		t.Fatalf("ImportData status = %d: %s", w.Code, w.Body.String())
	}
}

func TestSeparateImportsReplaceOnlySelectedData(t *testing.T) {
	_, router := newTestServer(t)

	w := doJSON(t, router, "POST", "/api/import/holidays", map[string]interface{}{
		"holidays": []map[string]interface{}{{"month": 9, "date": 1, "name": "Fall Break"}},
	})
	if w.Code != http.StatusOK {
		t.Fatalf("holiday import status = %d: %s", w.Code, w.Body.String())
	}
	w = doJSON(t, router, "POST", "/api/import/menu-items", map[string]interface{}{
		"items": []map[string]interface{}{{"meal_type": "lunch", "item_name": "Pizza"}},
	})
	if w.Code != http.StatusOK {
		t.Fatalf("menu import status = %d: %s", w.Code, w.Body.String())
	}
	w = doJSON(t, router, "GET", "/api/holidays", nil)
	var holidays []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &holidays)
	if len(holidays) != 1 || holidays[0]["name"] != "Fall Break" {
		t.Fatalf("unexpected imported holidays: %s", w.Body.String())
	}
	w = doJSON(t, router, "GET", "/api/menu-items", nil)
	var menus map[string][]map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &menus)
	if len(menus["lunch"]) != 1 || menus["lunch"][0]["item_name"] != "Pizza" {
		t.Fatalf("unexpected imported menu items: %s", w.Body.String())
	}
}

func TestSeparateImportValidation(t *testing.T) {
	_, router := newTestServer(t)
	w := doJSON(t, router, "POST", "/api/import/menu-items", map[string]interface{}{
		"items": []map[string]interface{}{{"meal_type": "dinner", "item_name": "Steak"}},
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid imported menu item, got %d", w.Code)
	}
}

func TestAuthMiddleware(t *testing.T) {
	mw := NewMiddleware("test-secret")

	handler := mw.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/holidays", nil)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for unauthenticated API request, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect for unauthenticated page request, got %d", w.Code)
	}
}

func TestUploadLogoRejectsBadExtension(t *testing.T) {
	logoDir := t.TempDir()
	h := NewLogoHandler(logoDir)

	body := &bytes.Buffer{}
	body.WriteString("--boundary\r\nContent-Disposition: form-data; name=\"logo\"; filename=\"evil.exe\"\r\n\r\ndata\r\n--boundary--\r\n")
	req := httptest.NewRequest("POST", "/api/logo", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")
	w := httptest.NewRecorder()
	h.UploadLogo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for .exe upload, got %d", w.Code)
	}

	entries, _ := os.ReadDir(logoDir)
	if len(entries) != 0 {
		t.Fatalf("expected empty logo dir, found %d files", len(entries))
	}
}

func TestUploadLogoAcceptsPNG(t *testing.T) {
	logoDir := t.TempDir()
	h := NewLogoHandler(logoDir)

	body := &bytes.Buffer{}
	body.WriteString("--boundary\r\nContent-Disposition: form-data; name=\"logo\"; filename=\"logo.png\"\r\n\r\ndata\r\n--boundary--\r\n")
	req := httptest.NewRequest("POST", "/api/logo", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")
	w := httptest.NewRecorder()
	h.UploadLogo(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for .png upload, got %d: %s", w.Code, w.Body.String())
	}
	if _, err := os.Stat(filepath.Join(logoDir, "logo.png")); err != nil {
		t.Fatal("logo.png not saved")
	}
}
