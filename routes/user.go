package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/raunak173/bms-go/controllers"
)

func UserRoutes(c *gin.Engine) {
	User := c.Group("/user")
	{
		User.POST("/login", controllers.Login)
		User.POST("/signup", controllers.SignUp)
	}
}
