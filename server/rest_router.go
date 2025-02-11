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

	// TODO: send back JSON with error message or result for each endpoint
	v1 := router.Group("v1")
	{
		siblingGroup := v1.Group("sibling")
		{
			sibling := new(controllers.SiblingController)
			siblingGroup.GET("", sibling.GetSiblings)
			siblingGroup.GET("/:id", sibling.GetSibling)

			siblingGroup.POST("", sibling.CreateSibling)

			siblingGroup.POST("/:id/start", sibling.StartSibling)
			siblingGroup.POST("/:id/stop", sibling.StopSibling)

			siblingGroup.POST("/:id/:node/start-node-iface", sibling.StartNodeIface)
			siblingGroup.POST("/:id/:node/stop-node-iface", sibling.StopNodeIface)

			siblingGroup.DELETE("", sibling.DeleteSiblings)
			siblingGroup.DELETE("/:id", sibling.DeleteSibling)
		}
	}
	return router
}
