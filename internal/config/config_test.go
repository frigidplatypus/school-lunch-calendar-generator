package config

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func writeConfig(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writeConfig: %v", err)
	}
	return path
}

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return path
}

func TestLoadWithPasswordHash(t *testing.T) {
	dir := t.TempDir()
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	path := writeConfig(t, dir, `password_hash: "`+string(hash)+`"
port: 9090
data_dir: /tmp/test
session_secret: "test-secret-1234"`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Port != 9090 {
		t.Fatalf("expected port 9090, got %d", cfg.Port)
	}
	if cfg.DataDir != "/tmp/test" {
		t.Fatalf("unexpected data_dir: %s", cfg.DataDir)
	}
	if cfg.SessionSecret != "test-secret-1234" {
		t.Fatalf("unexpected session_secret: %s", cfg.SessionSecret)
	}
	if cfg.PasswordExternal() {
		t.Fatal("expected PasswordExternal() = false")
	}
}

func TestLoadWithPasswordFile(t *testing.T) {
	dir := t.TempDir()
	pwFile := writeFile(t, dir, "pw.txt", "my-secret-pw\n")
	path := writeConfig(t, dir, `password_file: `+pwFile+`
data_dir: /tmp/test
session_secret: "test-secret-1234"`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.PasswordHash == "" {
		t.Fatal("expected PasswordHash to be set from file")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(cfg.PasswordHash), []byte("my-secret-pw")); err != nil {
		t.Fatalf("password hash does not match: %v", err)
	}
	if !cfg.PasswordExternal() {
		t.Fatal("expected PasswordExternal() = true")
	}
}

func TestLoadWithSessionSecretFile(t *testing.T) {
	dir := t.TempDir()
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	secFile := writeFile(t, dir, "secret.txt", "my-session-secret\n")
	path := writeConfig(t, dir, `password_hash: "`+string(hash)+`"
data_dir: /tmp/test
session_secret_file: `+secFile)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.SessionSecret != "my-session-secret" {
		t.Fatalf("expected session secret from file, got: %s", cfg.SessionSecret)
	}
}

func TestLoadPasswordFileEmpty(t *testing.T) {
	dir := t.TempDir()
	pwFile := writeFile(t, dir, "pw.txt", "  \n")
	path := writeConfig(t, dir, `password_file: `+pwFile+`
data_dir: /tmp/test
session_secret: "test-secret"`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty password_file")
	}
}

func TestLoadSessionSecretFileEmpty(t *testing.T) {
	dir := t.TempDir()
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	secFile := writeFile(t, dir, "secret.txt", "  \n")
	path := writeConfig(t, dir, `password_hash: "`+string(hash)+`"
data_dir: /tmp/test
session_secret_file: `+secFile)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty session_secret_file")
	}
}

func TestLoadNoPasswordSource(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, `data_dir: /tmp/test
session_secret: "test-secret"`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error when neither password_hash nor password_file is set")
	}
}

func TestLoadNoSessionSource(t *testing.T) {
	dir := t.TempDir()
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	path := writeConfig(t, dir, `password_hash: "`+string(hash)+`"
data_dir: /tmp/test`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error when neither session_secret nor session_secret_file is set")
	}
}

func TestLoadMissingPasswordFile(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, `password_file: /nonexistent/pw.txt
data_dir: /tmp/test
session_secret: "test-secret"`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing password_file")
	}
}

func TestLoadDefaults(t *testing.T) {
	dir := t.TempDir()
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	path := writeConfig(t, dir, `password_hash: "`+string(hash)+`"
session_secret: "test-secret"`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Port != 8080 {
		t.Fatalf("expected default port 8080, got %d", cfg.Port)
	}
	if cfg.DataDir != "./data" {
		t.Fatalf("expected default data_dir './data', got %s", cfg.DataDir)
	}
}

func TestLoadInvalidBcryptHash(t *testing.T) {
	dir := t.TempDir()
	path := writeConfig(t, dir, `password_hash: "not-a-valid-hash"
data_dir: /tmp/test
session_secret: "test-secret"`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid bcrypt hash")
	}
}

func TestUpdatePasswordHash(t *testing.T) {
	dir := t.TempDir()
	hash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	path := writeConfig(t, dir, `password_hash: "`+string(hash)+`"
port: 9090
data_dir: /tmp/test
session_secret: "test-secret"`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	newHash, _ := bcrypt.GenerateFromPassword([]byte("newpass"), bcrypt.DefaultCost)
	if err := cfg.UpdatePasswordHash(string(newHash)); err != nil {
		t.Fatalf("UpdatePasswordHash: %v", err)
	}

	cfg2, err := Load(path)
	if err != nil {
		t.Fatalf("Load (re): %v", err)
	}
	if cfg2.PasswordHash != string(newHash) {
		t.Fatal("password hash not persisted")
	}
	if cfg2.Port != 9090 {
		t.Fatalf("port changed after UpdatePasswordHash: %d", cfg2.Port)
	}
}

func TestPasswordExternal(t *testing.T) {
	cfg := &Config{PasswordFile: "/some/file"}
	if !cfg.PasswordExternal() {
		t.Fatal("expected true with PasswordFile set")
	}
	cfg2 := &Config{PasswordHash: "hash"}
	if cfg2.PasswordExternal() {
		t.Fatal("expected false with PasswordFile empty")
	}
}
