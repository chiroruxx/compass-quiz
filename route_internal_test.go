package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	healthPath        = "/healthz"
	healthBody        = "ok\n"
	healthContentType = "text/plain; charset=utf-8"
)

// TestHealthRoute checks the generated Gin route and its response without a database.
func TestHealthRoute(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, healthPath, nil)
	response := httptest.NewRecorder()
	newRouter().ServeHTTP(response, request)

	assertHealthz(t, response.Code, response.Header().Get("Content-Type"), response.Body.String())
}

func assertHealthz(t *testing.T, status int, contentType, body string) {
	t.Helper()

	if status != http.StatusOK {
		t.Errorf("status = %d, want %d", status, http.StatusOK)
	}

	if contentType != healthContentType {
		t.Errorf("Content-Type = %q, want %q", contentType, healthContentType)
	}

	if body != healthBody {
		t.Errorf("body = %q, want %q", body, healthBody)
	}
}
