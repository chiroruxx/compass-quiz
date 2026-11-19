package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"compass-quiz/api"
)

func newRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), errorMiddleware)

	var options api.GinServerOptions

	// The generated router calls ErrorHandler when parameter binding fails.
	// It then returns without calling the application handler.
	options.ErrorHandler = func(ctx *gin.Context, err error, _ int) {
		_ = ctx.Error(errors.Join(errValidation, err))
	}
	api.RegisterHandlersWithOptions(router, newHandler(newUseCase()), options)

	return router
}

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
