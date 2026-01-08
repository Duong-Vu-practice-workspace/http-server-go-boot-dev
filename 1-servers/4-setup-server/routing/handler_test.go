package routing

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupMux(t *testing.T) *http.ServeMux {
	t.Helper()
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "index.html"), []byte("Welcome to Chirpy"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &ApiConfig{}
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(tmp))
	// adjust middleware name if different in your code
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app/", fs)))
	mux.HandleFunc("/metrics", cfg.HandlerMetrics)
	mux.HandleFunc("/reset", cfg.HandlerReset)
	return mux
}

func TestMetricsAndReset(t *testing.T) {
	mux := setupMux(t)

	// ensure reset starts at 0
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reset", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("/reset status = %d; want %d", rec.Code, http.StatusOK)
	}

	// single hit
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/app/", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("/app/ status = %d; want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "Welcome to Chirpy") {
		t.Fatalf("/app/ body missing welcome text")
	}

	// metrics should show 1
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("/metrics status = %d; want %d", rec.Code, http.StatusOK)
	}
	body, _ := io.ReadAll(rec.Result().Body)
	if strings.TrimSpace(string(body)) != "Hits: 1" {
		t.Fatalf("/metrics = %q; want %q", string(body), "Hits: 1")
	}

	// make three more hits -> total 4
	for i := 0; i < 3; i++ {
		rec = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodGet, "/app/", nil)
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("/app/ status = %d; want %d", rec.Code, http.StatusOK)
		}
	}

	// metrics should show 4
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	mux.ServeHTTP(rec, req)
	body, _ = io.ReadAll(rec.Result().Body)
	if strings.TrimSpace(string(body)) != "Hits: 4" {
		t.Fatalf("/metrics = %q; want %q", string(body), "Hits: 4")
	}

	// reset and verify 0
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/reset", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("/reset status = %d; want %d", rec.Code, http.StatusOK)
	}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	mux.ServeHTTP(rec, req)
	body, _ = io.ReadAll(rec.Result().Body)
	if strings.TrimSpace(string(body)) != "Hits: 0" {
		t.Fatalf("/metrics after reset = %q; want %q", string(body), "Hits: 0")
	}
}
