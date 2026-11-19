package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

const (
	healthBody        = "ok\n"
	healthContentType = "text/plain; charset=utf-8"
)

func TestHealthHandler(t *testing.T) {
	t.Parallel()

	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	healthHandler{}.GetHealthz(ctx)

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
