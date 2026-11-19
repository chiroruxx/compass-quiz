package main

import (
	"github.com/gin-gonic/gin"

	"compass-quiz/api"
)

func newRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	var options api.GinServerOptions

	options.ErrorHandler = func(ctx *gin.Context, _ error, status int) {
		// Return the generated router's error status with an empty response body.
		ctx.AbortWithStatus(status)
	}
	api.RegisterHandlersWithOptions(router, handler{useCase: useCase{}}, options)

	return router
}
