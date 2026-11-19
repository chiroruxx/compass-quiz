package main

import (
	"errors"

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
	api.RegisterHandlersWithOptions(router, handler{useCase: useCase{}}, options)

	return router
}
