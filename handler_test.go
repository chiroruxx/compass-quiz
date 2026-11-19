package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthHandler(t *testing.T) {
	t.Parallel()

	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	healthHandler{}.GetHealthz(ctx)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("status = %d, want %d", status, http.StatusOK)
	}

	const wantContentType = "text/plain; charset=utf-8"
	if contentType := response.Header().Get("Content-Type"); contentType != wantContentType {
		t.Errorf("Content-Type = %q, want %q", contentType, wantContentType)
	}

	const wantBody = "ok\n"
	if body := response.Body.String(); body != wantBody {
		t.Errorf("body = %q, want %q", body, wantBody)
	}
}
