package main

import (
	"github.com/gin-gonic/gin"

	"compass-quiz/api"
)

func newRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	api.RegisterHandlers(router, healthHandler{})

	return router
}
