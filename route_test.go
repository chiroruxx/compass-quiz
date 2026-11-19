package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthRoute(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(newRouter())
	t.Cleanup(server.Close)

	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL+"/healthz", nil)
	if err != nil {
		t.Fatalf("create health request: %v", err)
	}

	response, err := server.Client().Do(request)
	if err != nil {
		t.Fatalf("send health request: %v", err)
	}

	t.Cleanup(func() {
		err := response.Body.Close()
		if err != nil {
			t.Errorf("close response body: %v", err)
		}
	})

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read health response: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	if string(body) != "ok\n" {
		t.Errorf("body = %q, want %q", body, "ok\n")
	}
}
