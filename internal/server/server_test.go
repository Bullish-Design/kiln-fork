// @feature:dev-server Tests for custom 404 page serving and notFoundRecorder.
package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestServe404_CustomPage(t *testing.T) {
	dir := t.TempDir()
	notFoundPage := filepath.Join(dir, "404.html")
	content := "<html><body>Custom 404</body></html>"
	if err := os.WriteFile(notFoundPage, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	serve404(w, notFoundPage)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("expected text/html content type, got %q", w.Header().Get("Content-Type"))
	}
	if w.Body.String() != content {
		t.Errorf("expected custom 404 body, got %q", w.Body.String())
	}
}

func TestServe404_MissingPage(t *testing.T) {
	w := httptest.NewRecorder()
	serve404(w, "/nonexistent/404.html")

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
	if w.Body.String() == "" {
		t.Error("expected fallback body, got empty")
	}
}

func TestNotFoundRecorder_Suppresses404(t *testing.T) {
	w := httptest.NewRecorder()
	rec := &notFoundRecorder{ResponseWriter: w}

	rec.WriteHeader(http.StatusNotFound)
	n, err := rec.Write([]byte("default 404 body"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 16 {
		t.Errorf("expected Write to report 16 bytes, got %d", n)
	}
	if rec.status != http.StatusNotFound {
		t.Errorf("expected recorded status 404, got %d", rec.status)
	}
	if w.Body.Len() != 0 {
		t.Errorf("expected suppressed body, got %q", w.Body.String())
	}
}

func TestNotFoundRecorder_PassesThrough(t *testing.T) {
	w := httptest.NewRecorder()
	rec := &notFoundRecorder{ResponseWriter: w}

	rec.WriteHeader(http.StatusOK)
	rec.Write([]byte("hello"))

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "hello" {
		t.Errorf("expected body 'hello', got %q", w.Body.String())
	}
}
