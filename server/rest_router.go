package server

// adapted from https://github.com/vsouza/go-gin-boilerplate

import (
	"github.com/Lachstec/digsinet-ng/server/controllers"
	"github.com/Lachstec/digsinet-ng/server/rest_middlewares"
	"github.com/gin-gonic/gin"
)

func NewRESTRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)
	router.Use(rest_middlewares.AuthMiddleware())

	v1 := router.Group("v1")
	{
		siblingGroup := v1.Group("sibling")
		{
			sibling := new(controllers.SiblingController)
			siblingGroup.GET("", sibling.GetSiblings)
			siblingGroup.GET("/:id", sibling.GetSiblingByID)

			siblingGroup.POST("", sibling.CreateSibling)

			siblingGroup.POST("/:id/start", sibling.StartSiblingByID)
			siblingGroup.POST("/:id/stop", sibling.StopSiblingByID)

			siblingGroup.DELETE("", sibling.DeleteSiblings)
			siblingGroup.DELETE("/:id", sibling.DeleteSiblingByID)
		}
	}
	return router
}
