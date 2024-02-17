package main

import (
	// "github.com/gin-gonic/gin"
	"github.com/sb/simple-backend/api/routes"
)

func main() {
	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// r.Run("0.0.0.0:9090") // listen and serve on 0.0.0.0:8080

	routes.InitializeRouter()
}
