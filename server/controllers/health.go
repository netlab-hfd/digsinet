package controllers

// adapted from https://github.com/vsouza/go-gin-boilerplate

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthController struct{}

func (h HealthController) Status(c *gin.Context) {
	c.String(http.StatusOK, "DigSiNet REST API is running!")
}
