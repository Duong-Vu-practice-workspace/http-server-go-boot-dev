package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogoServed(t *testing.T) {
	// If your code registers handlers in main(), call that setup function here.
	// Example: registerHandlers()
	// Alternatively, create file server directly:
	handler := http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))

	req := httptest.NewRequest("GET", "/assets/logo.png", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", rr.Code, http.StatusOK)
	}
	ct := rr.Header().Get("Content-Type")
	if ct == "" || (ct != "image/png" && ct != "image/png; charset=utf-8") {
		t.Fatalf("Content-Type = %q; want contain image/png", ct)
	}
}
