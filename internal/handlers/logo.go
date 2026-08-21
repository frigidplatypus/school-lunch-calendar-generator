package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type LogoHandler struct {
	logoDir string
}

func NewLogoHandler(logoDir string) *LogoHandler {
	return &LogoHandler{logoDir: logoDir}
}

func (h *LogoHandler) GetLogo(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(h.logoDir)
	if err != nil || len(entries) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	logoFile := ""
	for _, e := range entries {
		if !e.IsDir() {
			logoFile = e.Name()
			break
		}
	}

	if logoFile == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	path := filepath.Join(h.logoDir, logoFile)
	data, err := os.ReadFile(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ext := strings.ToLower(filepath.Ext(logoFile))
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	w.Write(data)
}

const maxLogoSize = 5 << 20

var allowedLogoExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".svg": true, ".webp": true,
}

func (h *LogoHandler) UploadLogo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxLogoSize+4096)
	r.ParseMultipartForm(maxLogoSize)

	file, header, err := r.FormFile("logo")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"no file uploaded or file too large"}`)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedLogoExts[ext] {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"unsupported file type"}`)
		return
	}

	entries, _ := os.ReadDir(h.logoDir)
	for _, e := range entries {
		if !e.IsDir() {
			os.Remove(filepath.Join(h.logoDir, e.Name()))
		}
	}

	filename := "logo" + ext
	path := filepath.Join(h.logoDir, filename)

	dst, err := os.Create(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"cannot save file"}`)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":"cannot save file"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","filename":"%s"}`, filename)
}
