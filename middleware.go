package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func errorMiddleware(ctx *gin.Context) {
	ctx.Next()

	lastError := ctx.Errors.Last()
	if lastError == nil || ctx.Writer.Written() {
		return
	}

	status := http.StatusInternalServerError
	if errors.Is(lastError.Err, errValidation) {
		status = http.StatusBadRequest
	}

	// Write only the status so error responses have no body.
	ctx.Status(status)
}
