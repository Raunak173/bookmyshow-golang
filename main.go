package main

import (
	"github.com/gin-gonic/gin"
	"github.com/raunak173/bms-go/initializers"
	"github.com/raunak173/bms-go/routes"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDB()
	initializers.CreateAWSUploader()
}

var R = gin.Default()

func main() {
	routes.MovieRoutes(R)
	routes.UserRoutes(R)
	routes.VenueRoutes(R)
	routes.SeatRoutes(R)
	routes.OrderRoutes(R)
	R.Run()
}
