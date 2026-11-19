package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

//go:generate go tool oapi-codegen -config oapi-codegen.yaml openapi.yaml

const readHeaderTimeout = 5 * time.Second

func main() {
	server := new(http.Server)
	server.Addr = ":8080"
	server.Handler = NewRouter()
	server.ReadHeaderTimeout = readHeaderTimeout

	err := server.ListenAndServe()
	if err != nil {
		slog.Error("HTTP server stopped", "error", err)
		os.Exit(1)
	}
}

// NewRouter registers the routes generated from the OpenAPI definition.
func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	RegisterHandlers(router, healthHandler{})

	return router
}
