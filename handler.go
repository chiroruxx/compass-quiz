package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthHandler struct{}

func (healthHandler) GetHealthz(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok\n")
}
