package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/raunak173/bms-go/controllers"
	"github.com/raunak173/bms-go/middleware"
)

func OrderRoutes(c *gin.Engine) {
	Order := c.Group("/orders")
	{
		Order.GET("/", middleware.RequireAuth, controllers.GetOrders)
	}
}
